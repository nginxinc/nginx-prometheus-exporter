# Changelog

### 0.10.0

A list of changes can be found on Github at: [Release v0.10.0](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.10.0)

### 0.9.0

A list of changes can be found on Github at: [Release v0.9.0](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.9.0)

### 0.8.0

CHANGES:
* [103](https://github.com/nginxinc/nginx-prometheus-exporter/pull/103): Switch to `gcr.io/distroless/static` image. Use a non-root user to run the exporter process by default. Thanks to [Alex SZAKALY](https://github.com/alex1989hu).
* Update Go version to 1.14

BUGFIXES:
* [99](https://github.com/nginxinc/nginx-prometheus-exporter/pull/99): Fix link to metrics path. Thanks to [Yoan Blanc](https://github.com/greut).
* [101](https://github.com/nginxinc/nginx-prometheus-exporter/pull/101): docs: fix dockerfile link. Thanks to [Eric Carboni](https://github.com/eric-hc).

UPGRADE:
* Use the 0.8.0 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.8.0`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.8.0).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R19 or newer.

### 0.7.0

FEATURES:
* [86](https://github.com/nginxinc/nginx-prometheus-exporter/pull/86): Implemented TLS client certificate authentication. Thanks to [Fabian Lüpke](https://github.com/Fluepke).

BUGFIXES:
* [96](https://github.com/nginxinc/nginx-prometheus-exporter/pull/96): Add const labels to upMetric. Thanks to [Robert Toth](https://github.com/robert-toth).

UPGRADE:
* Use the 0.7.0 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.7.0`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.7.0).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R19 or newer.

### 0.6.0

FEATURES:
* [77](https://github.com/nginxinc/nginx-prometheus-exporter/pull/77): Add constLabels support via cli arg/env variable.

CHANGES:
* Update alpine image.

UPGRADE:
* Use the 0.6.0 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.6.0`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.6.0).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R19 or newer.

### 0.5.0

FEATURES:
* [70](https://github.com/nginxinc/nginx-prometheus-exporter/pull/70): Set user agent on scrape requests to nginx.
* [68](https://github.com/nginxinc/nginx-prometheus-exporter/pull/68): Add ability to scrape and listen on unix domain sockets.
* [64](https://github.com/nginxinc/nginx-prometheus-exporter/pull/64): Add location zone and resolver metric support.

BUGFIXES:
* [73](https://github.com/nginxinc/nginx-prometheus-exporter/pull/73): Fix typo in stream_zone_sync_status_nodes_online metric description.
* [71](https://github.com/nginxinc/nginx-prometheus-exporter/pull/71): Do not assume default datasource in Grafana panels.
* [62](https://github.com/nginxinc/nginx-prometheus-exporter/pull/62): Set correct nginx_up query and instance variable expression.

UPGRADE:
* Use the 0.5.0 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.5.0`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.5.0).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R19 or newer.

### 0.4.2

BUGFIXES:
* [60](https://github.com/nginxinc/nginx-prometheus-exporter/pull/60): *Fix session metrics for stream server zones*. Session metrics with a status of `4xx` or `5xx` are now correctly reported. Previously they were always reported as `0`.

UPGRADE:
* Use the 0.4.2 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.4.2`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.4.2).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R18 or newer.

### 0.4.1

BUGFIXES:
* [55](https://github.com/nginxinc/nginx-prometheus-exporter/pull/55): Do not export zone sync metrics if they are not reported by NGINX Plus. Previously, in such case, the metrics were exported with zero values.

UPGRADE:
* Use the 0.4.1 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.4.1`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.4.1).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R18 or newer.

### 0.4.0

FEATURES:
* [50](https://github.com/nginxinc/nginx-prometheus-exporter/pull/50): Add zone sync metrics support.
* [37](https://github.com/nginxinc/nginx-prometheus-exporter/pull/37): Implement a way to retry connection to NGINX if it is unreachable. Add -nginx.retries for setting the number of retries and -nginx.retry-interval for setting the interval between retries, both as cli-arguments.

UPGRADE:
* Use the 0.4.0 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.4.0`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.4.0).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R18 or newer.

### 0.3.0

FEATURES:
* [32](https://github.com/nginxinc/nginx-prometheus-exporter/pull/32): Add nginxexporter_build_info metric.
* [31](https://github.com/nginxinc/nginx-prometheus-exporter/pull/31): Implement nginx_up and nginxplus_up metrics. Add -nginx.timeout cli argument for setting a timeout for scrapping metrics from NGINX or NGINX Plus.

UPGRADE:
* Use the 0.3.0 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.3.0`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.3.0).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R14 or newer.

## 0.2.0

FEATURES:
* [16](https://github.com/nginxinc/nginx-prometheus-exporter/pull/16): Add stream metrics support.
* [13](https://github.com/nginxinc/nginx-prometheus-exporter/pull/13): Add a flag for controlling SSL verification of NGINX stub_status/API endpoint. Thanks to [Raza Jhaveri](https://github.com/razaj92).
* [3](https://github.com/nginxinc/nginx-prometheus-exporter/pull/3): Support for environment variables.

UPGRADE:
* Use the 0.2.0 image from our DockerHub: `nginx/nginx-prometheus-exporter:0.2.0`
* Download the latest binaries from [GitHub releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases/tag/v0.2.0).

COMPATIBILITY:
* NGINX 0.1.18 or newer.
* NGINX Plus R14 or newer.

## 0.1.0

* Initial release.
