# Changelog

## [2.2.1](https://github.com/szpaczyna/speedtest-exporter/compare/v2.2.0...v2.2.1) (2026-11-28)


### Bug Fixes

* cache the speedtest result between scrapes ([77618c8](https://github.com/szpaczyna/speedtest-exporter/commit/77618c8fe904702470540e4d597731c376b3e205))
* cache the speedtest result between scrapes ([71eb019](https://github.com/szpaczyna/speedtest-exporter/commit/71eb01948f334b415a0761b88363dff38a18c77f))


### Miscellaneous

* **deps:** update grafana/grafana docker tag to v13.3.2 ([447cb4f](https://github.com/szpaczyna/speedtest-exporter/commit/447cb4f6560aadde78df1d1fcad8b15453efd800))
* **deps:** update grafana/grafana docker tag to v13.3.2 ([b589cad](https://github.com/szpaczyna/speedtest-exporter/commit/b589cadba34960ad6c7c8ea39bf355fc1c03fd61))


### CI/CD

* scope release permissions per job and drop QEMU ([6ff747e](https://github.com/szpaczyna/speedtest-exporter/commit/6ff747ebfd91aff6c6d4f2af8e82058f641dd5a0))
* scope release permissions per job and drop QEMU ([4a177c2](https://github.com/szpaczyna/speedtest-exporter/commit/4a177c2e88c6cefcbf50659ee753fb993f5cace3))

## [2.2.0](https://github.com/szpaczyna/speedtest-exporter/compare/v2.1.1...v2.2.0) (2026-11-23)


### Features

* log startup and speedtest results ([eb68051](https://github.com/szpaczyna/speedtest-exporter/commit/eb6805154a8a6afbef7594a04a8b58ac3c3f96e7))
* log startup and speedtest results ([c25d955](https://github.com/szpaczyna/speedtest-exporter/commit/c25d955073bf4428b28ede923622fb4d6d9ae050))

## [2.1.1](https://github.com/szpaczyna/speedtest-exporter/compare/v2.1.0...v2.1.1) (2026-11-22)


### Bug Fixes

* **docker:** verify Ookla CLI tarball checksum ([1bf66dc](https://github.com/szpaczyna/speedtest-exporter/commit/1bf66dc2e9ac94e65b30639ac7264aebdc237b2d))
* reject speedtest output that is not a result ([6974f80](https://github.com/szpaczyna/speedtest-exporter/commit/6974f80861724c4ee9a4a85ef20615f4d3e2403c))
* validate speedtest output and verify Ookla CLI checksum ([afdc250](https://github.com/szpaczyna/speedtest-exporter/commit/afdc250d603562235319cd39da4209958ca80d7f))


### CI/CD

* drop unused metadata step from PR image build ([3c0ca06](https://github.com/szpaczyna/speedtest-exporter/commit/3c0ca06f3b554248a206a582fe79c9417a14bec2))
* remove no-op sed from release version step ([c686dde](https://github.com/szpaczyna/speedtest-exporter/commit/c686ddec795ef91c7c795d740eef402b54ffb39c))


### Documentation

* bind docker run example to localhost ([109db0e](https://github.com/szpaczyna/speedtest-exporter/commit/109db0ed26f960a726f5726818ce1454d7987e6f))

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
