# Changelog

## [2.1.0](https://github.com/szpaczyna/speedtest-exporter/compare/v2.0.0...v2.1.0) (2026-11-21)


### Features

* **dashboard:** add download histogram and panel colors ([9d1b92d](https://github.com/szpaczyna/speedtest-exporter/commit/9d1b92df111b65c75a29c80d8462e2acf6c9401e))
* **dashboard:** add download histogram and panel colors ([abeb1bd](https://github.com/szpaczyna/speedtest-exporter/commit/abeb1bdbde9e0a900a80ba679508d66f101685c4))

## [2.0.0](https://github.com/szpaczyna/speedtest-exporter/compare/v1.2.0...v2.0.0) (2026-11-21)


### ⚠ BREAKING CHANGES

* the metrics path moved from /probe?script=speedtest to /metrics. speedtest_latency_seconds now reports seconds instead of milliseconds, speedtest_jittter_seconds is renamed to speedtest_jitter_seconds, speedtest_download_bytes and speedtest_upload_bytes are renamed to speedtest_download_speed_Bps and speedtest_upload_speed_Bps, and speedtest_downloadedbytes_bytes and speedtest_uploadedbytes_bytes are removed.

### Features

* **dashboard:** add gradient bars to the last result table ([b5b2395](https://github.com/szpaczyna/speedtest-exporter/commit/b5b2395f94d9c7c138edf18cee478a2db89144fd))
* rewrite exporter in Go on the Ookla CLI ([ad3ee92](https://github.com/szpaczyna/speedtest-exporter/commit/ad3ee9205cb11b480cbafa8b9bbd6e99fe20d0c9))


### Bug Fixes

* **dashboard:** show throughput in bits per second ([49e12f2](https://github.com/szpaczyna/speedtest-exporter/commit/49e12f21b0a4c6c189785f110f3639e09dd7cc88))
* log speedtest CLI stderr on failure ([35f2c5d](https://github.com/szpaczyna/speedtest-exporter/commit/35f2c5d0283ffca6a828fd3a522ef89677232e94))

## [1.2.0](https://github.com/szpaczyna/speedtest-exporter/compare/v1.1.0...v1.2.0) (2026-11-21)


### Features

* add multi-arch speedtest exporter image and release workflows ([216e6e9](https://github.com/szpaczyna/speedtest-exporter/commit/216e6e96004d88afc0035a2035c3662f039ffe37))
