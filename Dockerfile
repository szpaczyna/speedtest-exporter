FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /speedtest-exporter .

FROM alpine:3.24 AS speedtest

ARG SPEEDTEST_VERSION=1.2.0
ARG TARGETARCH

SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

RUN case "$TARGETARCH" in \
      amd64) arch=x86_64 ;; \
      arm64) arch=aarch64 ;; \
      *) echo "unsupported architecture: $TARGETARCH" >&2; exit 1 ;; \
    esac && \
    wget -qO- "https://install.speedtest.net/app/cli/ookla-speedtest-${SPEEDTEST_VERSION}-linux-${arch}.tgz" \
      | tar xz -C /tmp speedtest

FROM alpine:3.24

COPY --from=speedtest /tmp/speedtest /usr/local/bin/speedtest
COPY --from=build /speedtest-exporter /usr/local/bin/speedtest-exporter

# Ookla writes its license acceptance under $HOME.
ENV HOME=/tmp

USER 65535:65535

EXPOSE 9469

ENTRYPOINT ["/usr/local/bin/speedtest-exporter"]
