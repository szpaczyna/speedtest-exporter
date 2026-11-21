# speedtest-exporter

Prometheus exporter that runs the [Ookla Speedtest CLI](https://www.speedtest.net/apps/cli) through [script_exporter](https://github.com/ricoberger/script_exporter) and exposes the results on port 9469.

## Images

- `docker.io/szpaczyn/speedtest-exporter`
- `ghcr.io/szpaczyna/speedtest-exporter`

Both are built for `linux/amd64` and `linux/arm64`.

## Run

```bash
docker run -d --name speedtest-exporter -p 9469:9469 szpaczyn/speedtest-exporter:latest
curl 'http://localhost:9469/probe?script=speedtest'
```

A probe runs a full speed test, so one request takes about 30 seconds.

## Metrics

`speedtest_latency_seconds`, `speedtest_jittter_seconds`, `speedtest_download_bytes`, `speedtest_upload_bytes`, `speedtest_downloadedbytes_bytes`, `speedtest_uploadedbytes_bytes`.

## Build locally

```bash
docker buildx build --platform linux/arm64 --load -t speedtest-exporter:dev .
```

## Release

Release Please opens a release PR from Conventional Commits on `main`. Merging that PR tags the release, and the `Release Please` workflow pushes the image to Docker Hub and GHCR.

The workflow needs the `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` repository secrets. GHCR uses `GITHUB_TOKEN`.

## Deployment

The Helm chart lives in [rpi-k3s](https://github.com/szpaczyna/rpi-k3s/tree/main/cluster/helm/speedtest).
