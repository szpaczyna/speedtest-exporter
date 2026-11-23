package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func fixture(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/speedtest.json")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

const resultLabels = `{isp="Play",server_country="Poland",server_id="37249",server_location="Warszawa",server_name="Play"}`

func TestCollect(t *testing.T) {
	data := fixture(t)
	c := &collector{run: func() ([]byte, error) { return data, nil }}
	want := strings.ReplaceAll(`
# HELP speedtest_up Whether the last speedtest succeeded.
# TYPE speedtest_up gauge
speedtest_up 1
# HELP speedtest_latency_seconds Idle latency.
# TYPE speedtest_latency_seconds gauge
speedtest_latency_seconds{L} 0.015576000000000001
# HELP speedtest_jitter_seconds Idle jitter.
# TYPE speedtest_jitter_seconds gauge
speedtest_jitter_seconds{L} 0.002833
# HELP speedtest_packet_loss_ratio Packet loss ratio.
# TYPE speedtest_packet_loss_ratio gauge
speedtest_packet_loss_ratio{L} 0.005
# HELP speedtest_download_speed_Bps Download speed in bytes per second.
# TYPE speedtest_download_speed_Bps gauge
speedtest_download_speed_Bps{L} 8.8264522e+07
# HELP speedtest_upload_speed_Bps Upload speed in bytes per second.
# TYPE speedtest_upload_speed_Bps gauge
speedtest_upload_speed_Bps{L} 1.194281e+07
`, "{L}", resultLabels)
	err := testutil.CollectAndCompare(c, strings.NewReader(want),
		"speedtest_up", "speedtest_latency_seconds", "speedtest_jitter_seconds",
		"speedtest_packet_loss_ratio", "speedtest_download_speed_Bps", "speedtest_upload_speed_Bps")
	if err != nil {
		t.Fatal(err)
	}
}

func assertDown(t *testing.T, c *collector) {
	t.Helper()
	want := `
# HELP speedtest_up Whether the last speedtest succeeded.
# TYPE speedtest_up gauge
speedtest_up 0
`
	if err := testutil.CollectAndCompare(c, strings.NewReader(want), "speedtest_up"); err != nil {
		t.Fatal(err)
	}
	if n := testutil.CollectAndCount(c); n != 2 {
		t.Fatalf("got %d metrics, want up and scrape_duration only", n)
	}
}

func TestCollectError(t *testing.T) {
	assertDown(t, &collector{run: func() ([]byte, error) { return nil, errors.New("boom") }})
}

func TestCollectNotResult(t *testing.T) {
	assertDown(t, &collector{run: func() ([]byte, error) { return []byte("{}"), nil }})
}

func TestCollectInvalidJSON(t *testing.T) {
	assertDown(t, &collector{run: func() ([]byte, error) { return []byte("x"), nil }})
}

func TestCollectWithoutPacketLoss(t *testing.T) {
	var m map[string]any
	if err := json.Unmarshal(fixture(t), &m); err != nil {
		t.Fatal(err)
	}
	delete(m, "packetLoss")
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	c := &collector{run: func() ([]byte, error) { return data, nil }}
	if n := testutil.CollectAndCount(c, "speedtest_packet_loss_ratio"); n != 0 {
		t.Fatalf("got %d packet loss metrics, want 0", n)
	}
	if n := testutil.CollectAndCount(c, "speedtest_latency_seconds"); n != 1 {
		t.Fatalf("got %d latency metrics, want 1", n)
	}
}

func TestCollectConcurrent(t *testing.T) {
	data := fixture(t)
	var calls atomic.Int32
	started := make(chan struct{})
	c := &collector{run: func() ([]byte, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		time.Sleep(200 * time.Millisecond)
		return data, nil
	}}
	var wg sync.WaitGroup
	wg.Go(func() { testutil.CollectAndCount(c) })
	<-started
	wg.Go(func() { testutil.CollectAndCount(c) })
	wg.Wait()
	if n := calls.Load(); n != 1 {
		t.Fatalf("run called %d times, want 1", n)
	}
}

func TestWithStderr(t *testing.T) {
	_, err := exec.Command("sh", "-c", "echo 'Limit reached' >&2; exit 173").Output()
	got := withStderr(err).Error()
	if got != "exit status 173: Limit reached" {
		t.Fatalf("got %q", got)
	}
}

func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

func TestLogResult(t *testing.T) {
	buf := captureLog(t)
	data := fixture(t)
	testutil.CollectAndCount(&collector{run: func() ([]byte, error) { return data, nil }})
	for _, want := range []string{`msg="speedtest finished"`, "server_id=37249", "download_mbps=706.116176", "upload_mbps=95.54248"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("log missing %q:\n%s", want, buf)
		}
	}
}

func TestLogError(t *testing.T) {
	buf := captureLog(t)
	testutil.CollectAndCount(&collector{run: func() ([]byte, error) { return nil, errors.New("boom") }})
	if !strings.Contains(buf.String(), `level=ERROR msg="speedtest failed" err=boom`) {
		t.Errorf("unexpected log:\n%s", buf)
	}
}
