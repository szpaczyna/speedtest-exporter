package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	labels = []string{"server_id", "server_name", "server_location", "server_country", "isp"}

	upDesc         = prometheus.NewDesc("speedtest_up", "Whether the last speedtest succeeded.", nil, nil)
	durationDesc   = prometheus.NewDesc("speedtest_scrape_duration_seconds", "Duration of the speedtest run.", nil, nil)
	latencyDesc    = prometheus.NewDesc("speedtest_latency_seconds", "Idle latency.", labels, nil)
	jitterDesc     = prometheus.NewDesc("speedtest_jitter_seconds", "Idle jitter.", labels, nil)
	packetLossDesc = prometheus.NewDesc("speedtest_packet_loss_ratio", "Packet loss ratio.", labels, nil)
	downloadDesc   = prometheus.NewDesc("speedtest_download_speed_Bps", "Download speed in bytes per second.", labels, nil)
	uploadDesc     = prometheus.NewDesc("speedtest_upload_speed_Bps", "Upload speed in bytes per second.", labels, nil)
)

type result struct {
	Type string `json:"type"`
	Ping struct {
		Latency float64 `json:"latency"`
		Jitter  float64 `json:"jitter"`
	} `json:"ping"`
	PacketLoss *float64 `json:"packetLoss"`
	ISP        string   `json:"isp"`
	Server     struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
	} `json:"server"`
	Download struct {
		Bandwidth float64 `json:"bandwidth"`
	} `json:"download"`
	Upload struct {
		Bandwidth float64 `json:"bandwidth"`
	} `json:"upload"`
}

type collector struct {
	run func() ([]byte, error)

	mu       sync.Mutex
	finished time.Time
	metrics  []prometheus.Metric
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range []*prometheus.Desc{upDesc, durationDesc, latencyDesc, jitterDesc, packetLossDesc, downloadDesc, uploadDesc} {
		ch <- d
	}
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	arrived := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	// A test that finished while we waited for the lock is fresh enough to reuse.
	if !c.finished.After(arrived) {
		c.metrics = c.test()
		c.finished = time.Now()
	}
	for _, m := range c.metrics {
		ch <- m
	}
}

func (c *collector) test() []prometheus.Metric {
	slog.Info("speedtest started")
	start := time.Now()
	ms, r, err := c.parse()
	duration := time.Since(start)
	up := 1.0
	if err != nil {
		slog.Error("speedtest failed", "err", err, "duration", duration)
		up = 0
	} else {
		slog.Info("speedtest finished",
			"server_id", r.Server.ID,
			"server_name", r.Server.Name,
			"isp", r.ISP,
			"download_mbps", r.Download.Bandwidth*8/1e6,
			"upload_mbps", r.Upload.Bandwidth*8/1e6,
			"latency_ms", r.Ping.Latency,
			"duration", duration)
	}
	return append(ms,
		prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, up),
		prometheus.MustNewConstMetric(durationDesc, prometheus.GaugeValue, duration.Seconds()),
	)
}

func (c *collector) parse() ([]prometheus.Metric, *result, error) {
	out, err := c.run()
	if err != nil {
		return nil, nil, err
	}
	var r result
	if err := json.Unmarshal(out, &r); err != nil {
		return nil, nil, err
	}
	if r.Type != "result" {
		return nil, nil, fmt.Errorf("unexpected result type %q", r.Type)
	}
	lv := []string{strconv.Itoa(r.Server.ID), r.Server.Name, r.Server.Location, r.Server.Country, r.ISP}
	gauge := func(d *prometheus.Desc, v float64) prometheus.Metric {
		return prometheus.MustNewConstMetric(d, prometheus.GaugeValue, v, lv...)
	}
	ms := []prometheus.Metric{
		gauge(latencyDesc, r.Ping.Latency/1000),
		gauge(jitterDesc, r.Ping.Jitter/1000),
		gauge(downloadDesc, r.Download.Bandwidth),
		gauge(uploadDesc, r.Upload.Bandwidth),
	}
	if r.PacketLoss != nil {
		ms = append(ms, gauge(packetLossDesc, *r.PacketLoss/100))
	}
	return ms, &r, nil
}

func runSpeedtest() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "speedtest", "--accept-license", "--accept-gdpr", "-f", "json").Output()
	return out, withStderr(err)
}

// withStderr adds the CLI's stderr, which holds the reason (e.g. Ookla rate limit), to the bare exit status.
func withStderr(err error) error {
	var ee *exec.ExitError
	if errors.As(err, &ee) && len(ee.Stderr) > 0 {
		return fmt.Errorf("%w: %s", err, bytes.TrimSpace(ee.Stderr))
	}
	return err
}

func main() {
	prometheus.MustRegister(&collector{run: runSpeedtest})
	http.Handle("/metrics", promhttp.Handler())
	addr := ":9469"
	slog.Info("listening", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
