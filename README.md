# speedtest-exporter

Prometheus exporter in Go that runs the official [Ookla Speedtest CLI](https://www.speedtest.net/apps/cli) when scraped and serves the results on port 9469 at `/metrics`.

## Images

- `docker.io/szpaczyn/speedtest-exporter`
- `ghcr.io/szpaczyna/speedtest-exporter`

Both are built for `linux/amd64` and `linux/arm64`.

## Run

```bash
docker run -d --name speedtest-exporter -p 127.0.0.1:9469:9469 szpaczyn/speedtest-exporter:latest
curl http://localhost:9469/metrics
```

A scrape runs a full speed test, which takes about 30 seconds, unless the last test finished less than 10 minutes ago. In that case the exporter returns the last result, so extra requests to `/metrics` do not hit the Ookla rate limit. The exporter stops the CLI after 80 seconds, so set `scrape_timeout` to 90s or more and use a `scrape_interval` longer than 10 minutes:

```yaml
scrape_configs:
  - job_name: speedtest
    scrape_interval: 15m
    scrape_timeout: 90s
    static_configs:
      - targets: ["speedtest-exporter:9469"]
```

## Metrics

| Metric | Unit | Labels |
|---|---|---|
| `speedtest_up` | 1 if the last test succeeded, 0 otherwise | none |
| `speedtest_scrape_duration_seconds` | seconds | none |
| `speedtest_latency_seconds` | seconds | result labels |
| `speedtest_jitter_seconds` | seconds | result labels |
| `speedtest_download_speed_Bps` | bytes per second | result labels |
| `speedtest_upload_speed_Bps` | bytes per second | result labels |
| `speedtest_packet_loss_ratio` | ratio, 0 to 1 | result labels |

The result labels are `server_id`, `server_name`, `server_location`, `server_country` and `isp`. The exporter emits `speedtest_packet_loss_ratio` only when Ookla reports packet loss for the test.

## Local stack

```bash
docker compose up -d
```

The compose file starts the exporter, Prometheus and Grafana. Open Grafana at http://localhost:3000 (login `admin`/`admin`) and find the provisioned "Speedtest" dashboard. Prometheus scrapes every 15 minutes, so the first data point can take up to 15 minutes to appear.

To use the dashboard in another Grafana, import `dashboards/speedtest.json` and pick your Prometheus datasource in the `datasource` variable. The dashboard shows throughput in bits per second, like ISP plans and speedtest.net; the metrics stay in bytes per second.

## Development

```bash
go test ./...
```

GitHub Actions builds the images: `docker-build.yml` on pull requests and `release.yml` on releases.

## Migration from 1.x

The metrics path changed from `/probe?script=speedtest` to `/metrics`. Update the scrape config and any queries that use the old names:

| 1.x | 2.x |
|---|---|
| `speedtest_latency_seconds` (milliseconds despite the name) | `speedtest_latency_seconds` (seconds) |
| `speedtest_jittter_seconds` | `speedtest_jitter_seconds` |
| `speedtest_download_bytes` | `speedtest_download_speed_Bps` |
| `speedtest_upload_bytes` | `speedtest_upload_speed_Bps` |
| `speedtest_downloadedbytes_bytes` | removed |
| `speedtest_uploadedbytes_bytes` | removed |

New metrics: `speedtest_up`, `speedtest_scrape_duration_seconds` and `speedtest_packet_loss_ratio`.

## Release

Release Please opens a release PR from Conventional Commits on `main`. Merging that PR tags the release, and the `Release Please` workflow pushes the image to Docker Hub and GHCR.

The workflow needs the `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` repository secrets. GHCR uses `GITHUB_TOKEN`.

## Deployment

The Helm chart lives in [rpi-k3s](https://github.com/szpaczyna/rpi-k3s/tree/main/cluster/helm/speedtest).
