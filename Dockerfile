FROM alpine:3.23

ARG SPEEDTEST_VERSION=1.2.0
ARG SCRIPT_EXPORTER_VERSION=v2.1.2

# hadolint ignore=DL3018
RUN apk add --no-cache bash ca-certificates curl tar

RUN ARCH="$(apk info --print-arch)" && \
    curl -fsSL -o /tmp/ookla-speedtest.tgz \
      "https://install.speedtest.net/app/cli/ookla-speedtest-${SPEEDTEST_VERSION}-linux-${ARCH}.tgz" && \
    tar xzf /tmp/ookla-speedtest.tgz -C /usr/local/bin speedtest && \
    rm /tmp/ookla-speedtest.tgz

RUN case "$(apk info --print-arch)" in \
      x86_64) _arch=amd64 ;; \
      aarch64) _arch=arm64 ;; \
      *) echo "unsupported architecture" >&2; exit 1 ;; \
    esac && \
    curl -fsSL -o /usr/local/bin/script_exporter \
      "https://github.com/ricoberger/script_exporter/releases/download/${SCRIPT_EXPORTER_VERSION}/script_exporter-linux-${_arch}" && \
    chmod +x /usr/local/bin/script_exporter

COPY config.yaml /config.yaml
COPY speedtest-exporter.sh /usr/local/bin/speedtest-exporter.sh

EXPOSE 9469

ENTRYPOINT ["/usr/local/bin/script_exporter"]
