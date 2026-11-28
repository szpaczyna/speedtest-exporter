FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /speedtest-exporter .

FROM --platform=$BUILDPLATFORM alpine:3.24 AS speedtest

ARG SPEEDTEST_VERSION=1.2.0
ARG TARGETARCH

SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

RUN case "$TARGETARCH" in \
      amd64) arch=x86_64 sha256=5690596c54ff9bed63fa3732f818a05dbc2db19ad36ed68f21ca5f64d5cfeeb7 ;; \
      arm64) arch=aarch64 sha256=3953d231da3783e2bf8904b6dd72767c5c6e533e163d3742fd0437affa431bd3 ;; \
      *) echo "unsupported architecture: $TARGETARCH" >&2; exit 1 ;; \
    esac && \
    wget -qO /tmp/speedtest.tgz "https://install.speedtest.net/app/cli/ookla-speedtest-${SPEEDTEST_VERSION}-linux-${arch}.tgz" && \
    echo "${sha256}  /tmp/speedtest.tgz" | sha256sum -c - && \
    tar xzf /tmp/speedtest.tgz -C /tmp speedtest

FROM alpine:3.24

COPY --from=speedtest /tmp/speedtest /usr/local/bin/speedtest
COPY --from=build /speedtest-exporter /usr/local/bin/speedtest-exporter

# Ookla writes its license acceptance under $HOME.
ENV HOME=/tmp

USER 65535:65535

EXPOSE 9469

ENTRYPOINT ["/usr/local/bin/speedtest-exporter"]
