<!-- markdownlint-disable-next-line first-line-h1 -->
[![OpenSSFScorecard](https://api.securityscorecards.dev/projects/github.com/nginx/nginx-prometheus-exporter/badge)](https://scorecard.dev/viewer/?uri=github.com/nginx/nginx-prometheus-exporter)
[![CI](https://github.com/nginx/nginx-prometheus-exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/nginx/nginx-prometheus-exporter/actions/workflows/ci.yml)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B5618%2Fgithub.com%2Fnginx%2Fnginx-prometheus-exporter.svg?type=shield)](https://app.fossa.com/projects/custom%2B5618%2Fgithub.com%2Fnginx%2Fnginx-prometheus-exporter?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/nginx/nginx-prometheus-exporter)](https://goreportcard.com/report/github.com/nginx/nginx-prometheus-exporter)
[![codecov](https://codecov.io/gh/nginx/nginx-prometheus-exporter/graph/badge.svg?token=J6Oz10LWy3)](https://codecov.io/gh/nginx/nginx-prometheus-exporter)
![GitHub all releases](https://img.shields.io/github/downloads/nginx/nginx-prometheus-exporter/total?logo=github)
![GitHub release (latest by SemVer)](https://img.shields.io/github/downloads/nginx/nginx-prometheus-exporter/latest/total?sort=semver&logo=github)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nginx/nginx-prometheus-exporter?logo=github&sort=semver)](https://github.com/nginx/nginx-prometheus-exporter/releases/latest)
[![nginx-prometheus-exporter](https://snapcraft.io/nginx-prometheus-exporter/badge.svg)](https://snapcraft.io/nginx-prometheus-exporter)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nginx/nginx-prometheus-exporter?logo=go)
[![Docker Pulls](https://img.shields.io/docker/pulls/nginx/nginx-prometheus-exporter?logo=docker&logoColor=white)](https://hub.docker.com/r/nginx/nginx-prometheus-exporter)
![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/nginx/nginx-prometheus-exporter?logo=docker&logoColor=white&sort=semver)
[![Slack](https://img.shields.io/badge/slack-%23nginx--prometheus--exporter-green?logo=slack)](https://nginxcommunity.slack.com/channels/nginx-prometheus-exporter)
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

# NGINX Prometheus Exporter

NGINX Prometheus exporter makes it possible to monitor NGINX or NGINX Plus using Prometheus.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [A Note about NGINX Ingress Controller](#a-note-about-nginx-ingress-controller)
  - [Prerequisites](#prerequisites)
  - [Running the Exporter in a Docker Container](#running-the-exporter-in-a-docker-container)
  - [Running the Exporter Binary](#running-the-exporter-binary)
- [Usage](#usage)
  - [Command-line Arguments](#command-line-arguments)
- [Exported Metrics](#exported-metrics)
  - [Common metrics](#common-metrics)
  - [Metrics for NGINX OSS](#metrics-for-nginx-oss)
    - [Stub status metrics](#stub-status-metrics)
  - [Metrics for NGINX Plus](#metrics-for-nginx-plus)
    - [Connections](#connections)
    - [HTTP](#http)
    - [SSL](#ssl)
    - [HTTP Server Zones](#http-server-zones)
    - [Stream Server Zones](#stream-server-zones)
    - [HTTP Upstreams](#http-upstreams)
    - [Stream Upstreams](#stream-upstreams)
    - [Stream Zone Sync](#stream-zone-sync)
    - [Location Zones](#location-zones)
    - [Resolver](#resolver)
    - [HTTP Requests Rate Limiting](#http-requests-rate-limiting)
    - [HTTP Connections Limiting](#http-connections-limiting)
    - [Stream Connections Limiting](#stream-connections-limiting)
    - [Cache](#cache)
    - [Worker](#worker)
- [Troubleshooting](#troubleshooting)
- [Releases](#releases)
  - [Docker images](#docker-images)
  - [Binaries](#binaries)
  - [Homebrew](#homebrew)
  - [Snap](#snap)
  - [Scoop](#scoop)
  - [Nix](#nix)
- [Building the Exporter](#building-the-exporter)
  - [Building the Docker Image](#building-the-docker-image)
  - [Building the Binary](#building-the-binary)
- [Grafana Dashboard](#grafana-dashboard)
- [SBOM (Software Bill of Materials)](#sbom-software-bill-of-materials)
  - [Binaries](#binaries-1)
  - [Docker Image](#docker-image)
- [Provenance](#provenance)
- [Contacts](#contacts)
- [Contributing](#contributing)
- [Support](#support)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Overview

[NGINX](https://nginx.org) exposes a handful of metrics via the [stub_status
page](https://nginx.org/en/docs/http/ngx_http_stub_status_module.html#stub_status). [NGINX
Plus](https://www.nginx.com/products/nginx/) provides a richer set of metrics via the
[API](https://nginx.org/en/docs/http/ngx_http_api_module.html) and the [monitoring
dashboard](https://docs.nginx.com/nginx/admin-guide/monitoring/live-activity-monitoring/). NGINX Prometheus exporter
fetches the metrics from a single NGINX or NGINX Plus, converts the metrics into appropriate Prometheus metrics types
and finally exposes them via an HTTP server to be collected by [Prometheus](https://prometheus.io/).

## Getting Started

In this section, we show how to quickly run NGINX Prometheus Exporter for NGINX or NGINX Plus.

### A Note about NGINX Ingress Controller

If you’d like to use the NGINX Prometheus Exporter with [NGINX Ingress
Controller](https://github.com/nginx/kubernetes-ingress/) for Kubernetes, see [this
doc](https://docs.nginx.com/nginx-ingress-controller/logging-and-monitoring/prometheus/) for the installation
instructions.

### Prerequisites

We assume that you have already installed Prometheus and NGINX or NGINX Plus. Additionally, you need to:

- Expose the built-in metrics in NGINX/NGINX Plus:
  - For NGINX, expose the [stub_status
    page](https://nginx.org/en/docs/http/ngx_http_stub_status_module.html#stub_status) at `/stub_status` on port `8080`.
  - For NGINX Plus, expose the [API](https://nginx.org/en/docs/http/ngx_http_api_module.html#api) at `/api` on port
    `8080`.
- Configure Prometheus to scrape metrics from the server with the exporter. Note that the default scrape port of the
  exporter is `9113` and the default metrics path -- `/metrics`.

### Running the Exporter in a Docker Container

To start the exporter we use the [docker run](https://docs.docker.com/engine/reference/run/) command.

- To export NGINX metrics, run:

  ```console
  docker run -p 9113:9113 nginx/nginx-prometheus-exporter:1.4.0 --nginx.scrape-uri=http://<nginx>:8080/stub_status
  ```

  where `<nginx>` is the IP address/DNS name, through which NGINX is available.

- To export NGINX Plus metrics, run:

  ```console
  docker run -p 9113:9113 nginx/nginx-prometheus-exporter:1.4.0 --nginx.plus --nginx.scrape-uri=http://<nginx-plus>:8080/api
  ```

  where `<nginx-plus>` is the IP address/DNS name, through which NGINX Plus is available.

### Running the Exporter Binary

- To export NGINX metrics, run:

  ```console
  nginx-prometheus-exporter --nginx.scrape-uri=http://<nginx>:8080/stub_status
  ```

  where `<nginx>` is the IP address/DNS name, through which NGINX is available.

- To export NGINX Plus metrics:

  ```console
  nginx-prometheus-exporter --nginx.plus --nginx.scrape-uri=http://<nginx-plus>:8080/api
  ```

  where `<nginx-plus>` is the IP address/DNS name, through which NGINX Plus is available.

- To scrape NGINX metrics with unix domain sockets, run:

  ```console
  nginx-prometheus-exporter --nginx.scrape-uri=unix:<nginx>:/stub_status
  ```

  where `<nginx>` is the path to unix domain socket, through which NGINX stub status is available.

**Note**. The `nginx-prometheus-exporter` is not a daemon. To run the exporter as a system service (daemon), you can
follow the example in [examples/systemd](./examples/systemd/README.md). Alternatively, you can run the exporter
in a Docker container.

## Usage

### Command-line Arguments

```console
usage: nginx-prometheus-exporter [<flags>]


Flags:
  -h, --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
      --[no-]web.systemd-socket  Use systemd socket activation listeners instead
                                 of port listeners (Linux only). ($SYSTEMD_SOCKET)
      --web.listen-address=:9113 ...
                                 Addresses on which to expose metrics and web interface. Repeatable for multiple addresses. ($LISTEN_ADDRESS)
      --web.config.file=""       Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md ($CONFIG_FILE)
      --web.telemetry-path="/metrics"
                                 Path under which to expose metrics. ($TELEMETRY_PATH)
      --[no-]nginx.plus          Start the exporter for NGINX Plus. By default, the exporter is started for NGINX. ($NGINX_PLUS)
      --nginx.scrape-uri=http://127.0.0.1:8080/stub_status ...
                                 A URI or unix domain socket path for scraping NGINX or NGINX Plus metrics. For NGINX, the stub_status page must be available through the URI. For NGINX Plus -- the API. Repeatable for multiple URIs. ($SCRAPE_URI)
      --[no-]nginx.ssl-verify    Perform SSL certificate verification. ($SSL_VERIFY)
      --nginx.ssl-ca-cert=""     Path to the PEM encoded CA certificate file used to validate the servers SSL certificate. ($SSL_CA_CERT)
      --nginx.ssl-client-cert=""
                                 Path to the PEM encoded client certificate file to use when connecting to the server. ($SSL_CLIENT_CERT)
      --nginx.ssl-client-key=""  Path to the PEM encoded client certificate key file to use when connecting to the server. ($SSL_CLIENT_KEY)
      --nginx.timeout=5s         A timeout for scraping metrics from NGINX or NGINX Plus. ($TIMEOUT)
      --prometheus.const-label=PROMETHEUS.CONST-LABEL ...
                                 Label that will be used in every metric. Format is label=value. It can be repeated multiple times. ($CONST_LABELS)
      --log.level=info           Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt        Output format of log messages. One of: [logfmt, json]
      --[no-]version             Show application version.
```

## Exported Metrics

### Common metrics

| Name                                         | Type     | Description                                  | Labels                                                                    |
| -------------------------------------------- | -------- | -------------------------------------------- | ------------------------------------------------------------------------- |
| `nginx_exporter_build_info`                  | Gauge    | Shows the exporter build information.        | `branch`, `goarch`, `goos`, `goversion`, `revision`, `tags` and `version` |
| `promhttp_metric_handler_requests_total`     | Counter  | Total number of scrapes by HTTP status code. | `code` (the HTTP status code)                                             |
| `promhttp_metric_handler_requests_in_flight` | Gauge    | Current number of scrapes being served.      | []                                                                        |
| `go_*`                                       | Multiple | Go runtime metrics.                          | []                                                                        |

### Metrics for NGINX OSS

| Name       | Type  | Description                                                                                      | Labels |
| ---------- | ----- | ------------------------------------------------------------------------------------------------ | ------ |
| `nginx_up` | Gauge | Shows the status of the last metric scrape: `1` for a successful scrape and `0` for a failed one | []     |

#### [Stub status metrics](https://nginx.org/en/docs/http/ngx_http_stub_status_module.html)

| Name                         | Type    | Description                                                         | Labels |
| ---------------------------- | ------- | ------------------------------------------------------------------- | ------ |
| `nginx_connections_accepted` | Counter | Accepted client connections.                                        | []     |
| `nginx_connections_active`   | Gauge   | Active client connections.                                          | []     |
| `nginx_connections_handled`  | Counter | Handled client connections.                                         | []     |
| `nginx_connections_reading`  | Gauge   | Connections where NGINX is reading the request header.              | []     |
| `nginx_connections_waiting`  | Gauge   | Idle client connections.                                            | []     |
| `nginx_connections_writing`  | Gauge   | Connections where NGINX is writing the response back to the client. | []     |
| `nginx_http_requests_total`  | Counter | Total http requests.                                                | []     |

### Metrics for NGINX Plus

| Name           | Type  | Description                                                                                      | Labels |
| -------------- | ----- | ------------------------------------------------------------------------------------------------ | ------ |
| `nginxplus_up` | Gauge | Shows the status of the last metric scrape: `1` for a successful scrape and `0` for a failed one | []     |

#### [Connections](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_connections)

| Name                             | Type    | Description                        | Labels |
| -------------------------------- | ------- | ---------------------------------- | ------ |
| `nginxplus_connections_accepted` | Counter | Accepted client connections        | []     |
| `nginxplus_connections_active`   | Gauge   | Active client connections          | []     |
| `nginxplus_connections_dropped`  | Counter | Dropped client connections dropped | []     |
| `nginxplus_connections_idle`     | Gauge   | Idle client connections            | []     |

#### [HTTP](https://nginx.org/en/docs/http/ngx_http_api_module.html#http_)

| Name                              | Type    | Description           | Labels |
| --------------------------------- | ------- | --------------------- | ------ |
| `nginxplus_http_requests_total`   | Counter | Total http requests   | []     |
| `nginxplus_http_requests_current` | Gauge   | Current http requests | []     |

#### [SSL](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_ssl_object)

| Name                              | Type    | Description                         | Labels |
| --------------------------------- | ------- | ----------------------------------- | ------ |
| `nginxplus_ssl_handshakes`        | Counter | Successful SSL handshakes           | []     |
| `nginxplus_ssl_handshakes_failed` | Counter | Failed SSL handshakes               | []     |
| `nginxplus_ssl_session_reuses`    | Counter | Session reuses during SSL handshake | []     |

#### [HTTP Server Zones](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_server_zone)

| Name                                     | Type    | Description                                        | Labels                                                                                                                                     |
| ---------------------------------------- | ------- | -------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `nginxplus_server_zone_processing`       | Gauge   | Client requests that are currently being processed | `server_zone`                                                                                                                              |
| `nginxplus_server_zone_requests`         | Counter | Total client requests                              | `server_zone`                                                                                                                              |
| `nginxplus_server_zone_responses`        | Counter | Total responses sent to clients                    | `code` (the response status code. The values are: `1xx`, `2xx`, `3xx`, `4xx` and `5xx`), `server_zone`                                     |
| `nginxplus_server_zone_responses_codes`  | Counter | Total responses sent to clients by code            | `code` (the response status code. The possible values are [here](https://www.nginx.com/resources/wiki/extending/api/http/)), `server_zone` |
| `nginxplus_server_zone_discarded`        | Counter | Requests completed without sending a response      | `server_zone`                                                                                                                              |
| `nginxplus_server_zone_received`         | Counter | Bytes received from clients                        | `server_zone`                                                                                                                              |
| `nginxplus_server_zone_sent`             | Counter | Bytes sent to clients                              | `server_zone`                                                                                                                              |
| `nginxplus_server_ssl_handshakes`        | Counter | Successful SSL handshakes                          | `server_zone`                                                                                                                              |
| `nginxplus_server_ssl_handshakes_failed` | Counter | Failed SSL handshakes                              | `server_zone`                                                                                                                              |
| `nginxplus_server_ssl_session_reuses`    | Counter | Session reuses during SSL handshake                | `server_zone`                                                                                                                              |

#### [Stream Server Zones](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_stream_server_zone)

| Name                                            | Type    | Description                                           | Labels                                                                                    |
| ----------------------------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `nginxplus_stream_server_zone_processing`       | Gauge   | Client connections that are currently being processed | `server_zone`                                                                             |
| `nginxplus_stream_server_zone_connections`      | Counter | Total connections                                     | `server_zone`                                                                             |
| `nginxplus_stream_server_zone_sessions`         | Counter | Total sessions completed                              | `code` (the response status code. The values are: `2xx`, `4xx`, and `5xx`), `server_zone` |
| `nginxplus_stream_server_zone_discarded`        | Counter | Connections completed without creating a session      | `server_zone`                                                                             |
| `nginxplus_stream_server_zone_received`         | Counter | Bytes received from clients                           | `server_zone`                                                                             |
| `nginxplus_stream_server_zone_sent`             | Counter | Bytes sent to clients                                 | `server_zone`                                                                             |
| `nginxplus_stream_server_ssl_handshakes`        | Counter | Successful SSL handshakes                             | `server_zone`                                                                             |
| `nginxplus_stream_server_ssl_handshakes_failed` | Counter | Failed SSL handshakes                                 | `server_zone`                                                                             |
| `nginxplus_stream_server_ssl_session_reuses`    | Counter | Session reuses during SSL handshake                   | `server_zone`                                                                             |

#### [HTTP Upstreams](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_upstream)

> Note: for the `state` metric, the string values are converted to float64 using the following rule: `"up"` -> `1.0`,
> `"draining"` -> `2.0`, `"down"` -> `3.0`, `"unavail"` –> `4.0`, `"checking"` –> `5.0`, `"unhealthy"` -> `6.0`.

| Name                                                | Type    | Description                                                                                                                                                    | Labels                                                                                                                                            |
| --------------------------------------------------- | ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `nginxplus_upstream_server_state`                   | Gauge   | Current state                                                                                                                                                  | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_active`                  | Gauge   | Active connections                                                                                                                                             | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_limit`                   | Gauge   | Limit for connections which corresponds to the max_conns parameter of the upstream server. Zero value means there is no limit                                  | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_requests`                | Counter | Total client requests                                                                                                                                          | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_responses`               | Counter | Total responses sent to clients                                                                                                                                | `code` (the response status code. The values are: `1xx`, `2xx`, `3xx`, `4xx` and `5xx`), `server`, `upstream`                                     |
| `nginxplus_upstream_server_responses_codes`         | Counter | Total responses sent to clients by code                                                                                                                        | `code` (the response status code. The possible values are [here](https://www.nginx.com/resources/wiki/extending/api/http/)), `server`, `upstream` |
| nginxplus_upstream_server_sent`                     | Counter | Bytes sent to this server                                                                                                                                      | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_received`                | Counter | Bytes received to this server                                                                                                                                  | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_fails`                   | Counter | Number of unsuccessful attempts to communicate with the server                                                                                                 | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_unavail`                 | Counter | How many times the server became unavailable for client requests (state 'unavail') due to the number of unsuccessful attempts reaching the max_fails threshold | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_header_time`             | Gauge   | Average time to get the response header from the server                                                                                                        | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_response_time`           | Gauge   | Average time to get the full response from the server                                                                                                          | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_health_checks_checks`    | Counter | Total health check requests                                                                                                                                    | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_health_checks_fails`     | Counter | Failed health checks                                                                                                                                           | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_health_checks_unhealthy` | Counter | How many times the server became unhealthy (state 'unhealthy')                                                                                                 | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_ssl_handshakes`          | Counter | Successful SSL handshakes                                                                                                                                      | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_ssl_handshakes_failed`   | Counter | Failed SSL handshakes                                                                                                                                          | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_server_ssl_session_reuses`      | Counter | Session reuses during SSL handshake                                                                                                                            | `server`, `upstream`                                                                                                                              |
| `nginxplus_upstream_keepalive`                      | Gauge   | Idle keepalive connections                                                                                                                                     | `upstream`                                                                                                                                        |
| `nginxplus_upstream_zombies`                        | Gauge   | Servers removed from the group but still processing active client requests                                                                                     | `upstream`                                                                                                                                        |

#### [Stream Upstreams](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_stream_upstream)

> Note: for the `state` metric, the string values are converted to float64 using the following rule: `"up"` -> `1.0`,
> `"down"` -> `3.0`, `"unavail"` –> `4.0`, `"checking"` –> `5.0`, `"unhealthy"` -> `6.0`.

| Name                                                       | Type    | Description                                                                                                                                                       | Labels                |
| ---------------------------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| `nginxplus_stream_upstream_server_state`                   | Gauge   | Current state                                                                                                                                                     | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_active`                  | Gauge   | Active connections                                                                                                                                                | `server` , `upstream` |
| `nginxplus_stream_upstream_server_limit`                   | Gauge   | Limit for connections which corresponds to the max_conns parameter of the upstream server. Zero value means there is no limit                                     | `server` , `upstream` |
| `nginxplus_stream_upstream_server_connections`             | Counter | Total number of client connections forwarded to this server                                                                                                       | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_connect_time`            | Gauge   | Average time to connect to the upstream server                                                                                                                    | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_first_byte_time`         | Gauge   | Average time to receive the first byte of data                                                                                                                    | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_response_time`           | Gauge   | Average time to receive the last byte of data                                                                                                                     | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_sent`                    | Counter | Bytes sent to this server                                                                                                                                         | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_received`                | Counter | Bytes received from this server                                                                                                                                   | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_fails`                   | Counter | Number of unsuccessful attempts to communicate with the server                                                                                                    | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_unavail`                 | Counter | How many times the server became unavailable for client connections (state 'unavail') due to the number of unsuccessful attempts reaching the max_fails threshold | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_health_checks_checks`    | Counter | Total health check requests                                                                                                                                       | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_health_checks_fails`     | Counter | Failed health checks                                                                                                                                              | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_health_checks_unhealthy` | Counter | How many times the server became unhealthy (state 'unhealthy')                                                                                                    | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_ssl_handshakes`          | Counter | Successful SSL handshakes                                                                                                                                         | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_ssl_handshakes_failed`   | Counter | Failed SSL handshakes                                                                                                                                             | `server`, `upstream`  |
| `nginxplus_stream_upstream_server_ssl_session_reuses`      | Counter | Session reuses during SSL handshake                                                                                                                               | `server`, `upstream`  |
| `nginxplus_stream_upstream_zombies`                        | Gauge   | Servers removed from the group but still processing active client connections                                                                                     | `upstream`            |

#### [Stream Zone Sync](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_stream_zone_sync)

| Name                                              | Type    | Description                                                  | Labels |
| ------------------------------------------------- | ------- | ------------------------------------------------------------ | ------ |
| `nginxplus_stream_zone_sync_zone_records_pending` | Gauge   | The number of records that need to be sent to the cluster    | `zone` |
| `nginxplus_stream_zone_sync_zone_records_total`   | Gauge   | The total number of records stored in the shared memory zone | `zone` |
| `nginxplus_stream_zone_sync_zone_bytes_in`        | Counter | Bytes received by this node                                  | []     |
| `nginxplus_stream_zone_sync_zone_bytes_out`       | Counter | Bytes sent by this node                                      | []     |
| `nginxplus_stream_zone_sync_zone_msgs_in`         | Counter | Total messages received by this node                         | []     |
| `nginxplus_stream_zone_sync_zone_msgs_out`        | Counter | Total messages sent by this node                             | []     |
| `nginxplus_stream_zone_sync_zone_nodes_online`    | Gauge   | Number of peers this node is connected to                    | []     |

#### [Location Zones](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_location_zone)

| Name                                      | Type    | Description                                   | Labels                                                                                                                                       |
| ----------------------------------------- | ------- | --------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `nginxplus_location_zone_requests`        | Counter | Total client requests                         | `location_zone`                                                                                                                              |
| `nginxplus_location_zone_responses`       | Counter | Total responses sent to clients               | `code` (the response status code. The values are: `1xx`, `2xx`, `3xx`, `4xx` and `5xx`), `location_zone`                                     |
| `nginxplus_location_zone_responses_codes` | Counter | Total responses sent to clients by code       | `code` (the response status code. The possible values are [here](https://www.nginx.com/resources/wiki/extending/api/http/)), `location_zone` |
| `nginxplus_location_zone_discarded`       | Counter | Requests completed without sending a response | `location_zone`                                                                                                                              |
| `nginxplus_location_zone_received`        | Counter | Bytes received from clients                   | `location_zone`                                                                                                                              |
| `nginxplus_location_zone_sent`            | Counter | Bytes sent to clients                         | `location_zone`                                                                                                                              |

#### [Resolver](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_resolver_zone)

| Name                          | Type    | Description                                    | Labels     |
| ----------------------------- | ------- | ---------------------------------------------- | ---------- |
| `nginxplus_resolver_name`     | Counter | Total requests to resolve names to addresses   | `resolver` |
| `nginxplus_resolver_srv`      | Counter | Total requests to resolve SRV records          | `resolver` |
| `nginxplus_resolver_addr`     | Counter | Total requests to resolve addresses to names   | `resolver` |
| `nginxplus_resolver_noerror`  | Counter | Total number of successful responses           | `resolver` |
| `nginxplus_resolver_formerr`  | Counter | Total number of FORMERR responses              | `resolver` |
| `nginxplus_resolver_servfail` | Counter | Total number of SERVFAIL responses             | `resolver` |
| `nginxplus_resolver_nxdomain` | Counter | Total number of NXDOMAIN responses             | `resolver` |
| `nginxplus_resolver_notimp`   | Counter | Total number of NOTIMP responses               | `resolver` |
| `nginxplus_resolver_refused`  | Counter | Total number of REFUSED responses              | `resolver` |
| `nginxplus_resolver_timedout` | Counter | Total number of timed out request              | `resolver` |
| `nginxplus_resolver_unknown`  | Counter | Total requests completed with an unknown error | `resolver` |

#### [HTTP Requests Rate Limiting](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_limit_req_zone)

| Name                                       | Type    | Description                                                                 | Labels |
| ------------------------------------------ | ------- | --------------------------------------------------------------------------- | ------ |
| `nginxplus_limit_request_passed`           | Counter | Total number of requests that were neither limited nor accounted as limited | `zone` |
| `nginxplus_limit_request_rejected`         | Counter | Total number of requests that were rejected                                 | `zone` |
| `nginxplus_limit_request_delayed`          | Counter | Total number of requests that were delayed                                  | `zone` |
| `nginxplus_limit_request_rejected_dry_run` | Counter | Total number of requests accounted as rejected in the dry run mode          | `zone` |
| `nginxplus_limit_request_delayed_dry_run`  | Counter | Total number of requests accounted as delayed in the dry run mode           | `zone` |

#### [HTTP Connections Limiting](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_limit_conn_zone)

| Name                                          | Type    | Description                                                                    | Labels |
| --------------------------------------------- | ------- | ------------------------------------------------------------------------------ | ------ |
| `nginxplus_limit_connection_passed`           | Counter | Total number of connections that were neither limited nor accounted as limited | `zone` |
| `nginxplus_limit_connection_rejected`         | Counter | Total number of connections that were rejected                                 | `zone` |
| `nginxplus_limit_connection_rejected_dry_run` | Counter | Total number of connections accounted as rejected in the dry run mode          | `zone` |

#### [Stream Connections Limiting](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_stream_limit_conn_zone)

| Name                                                 | Type    | Description                                                                    | Labels |
| ---------------------------------------------------- | ------- | ------------------------------------------------------------------------------ | ------ |
| `nginxplus_stream_limit_connection_passed`           | Counter | Total number of connections that were neither limited nor accounted as limited | `zone` |
| `nginxplus_stream_limit_connection_rejected`         | Counter | Total number of connections that were rejected                                 | `zone` |
| `nginxplus_stream_limit_connection_rejected_dry_run` | Counter | Total number of connections accounted as rejected in the dry run mode          | `zone` |

#### [Cache](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_cache)

| Name                                        | Type    | Description                                                             | Labels  |
| ------------------------------------------- | ------- | ----------------------------------------------------------------------- | ------- |
| `nginxplus_cache_size`                      | Gauge   | Total size of the cache                                                 | `cache` |
| `nginxplus_cache_max_size`                  | Gauge   | Maximum size of the cache                                               | `cache` |
| `nginxplus_cache_cold`                      | Gauge   | Is the cache considered cold                                            | `cache` |
| `nginxplus_cache_hit_responses`             | Counter | Total number of cache hits                                              | `cache` |
| `nginxplus_cache_hit_bytes`                 | Counter | Total number of bytes returned from cache hits                          | `cache` |
| `nginxplus_cache_stale_responses`           | Counter | Total number of stale cache hits                                        | `cache` |
| `nginxplus_cache_stale_bytes`               | Counter | Total number of bytes returned from stale cache hits                    | `cache` |
| `nginxplus_cache_updating_responses`        | Counter | Total number of cache hits while cache is updating                      | `cache` |
| `nginxplus_cache_updating_bytes`            | Counter | Total number of bytes returned from cache while cache is updating       | `cache` |
| `nginxplus_cache_revalidated_responses`     | Counter | Total number of cache revalidations                                     | `cache` |
| `nginxplus_cache_revalidated_bytes`         | Counter | Total number of bytes returned from cache revalidations                 | `cache` |
| `nginxplus_cache_miss_responses`            | Counter | Total number of cache misses                                            | `cache` |
| `nginxplus_cache_miss_bytes`                | Counter | Total number of bytes returned from cache misses                        | `cache` |
| `nginxplus_cache_expired_responses`         | Counter | Total number of cache hits with expired TTL                             | `cache` |
| `nginxplus_cache_expired_bytes`             | Counter | Total number of bytes returned from cache hits with expired TTL         | `cache` |
| `nginxplus_cache_expired_responses_written` | Counter | Total number of cache hits with expired TTL written to cache            | `cache` |
| `nginxplus_cache_expired_bytes_written`     | Counter | Total number of bytes written to cache from cache hits with expired TTL | `cache` |
| `nginxplus_cache_bypass_responses`          | Counter | Total number of cache bypasses                                          | `cache` |
| `nginxplus_cache_bypass_bytes`              | Counter | Total number of bytes returned from cache bypasses                      | `cache` |
| `nginxplus_cache_bypass_responses_written`  | Counter | Total number of cache bypasses written to cache                         | `cache` |
| `nginxplus_cache_bypass_bytes_written`      | Counter | Total number of bytes written to cache from cache bypasses              | `cache` |

#### [Worker](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_worker)

| Name                                     | Type    | Description                                                              | Labels      |
| ---------------------------------------- | ------- | ------------------------------------------------------------------------ | ----------- |
| `nginxplus_worker_connection_accepted`   | Counter | The total number of accepted client connections                          | `id`, `pid` |
| `nginxplus_worker_connection_dropped`    | Counter | The total number of dropped client connections                           | `id`, `pid` |
| `nginxplus_worker_connection_active`     | Gauge   | The current number of active client connections                          | `id`, `pid` |
| `nginxplus_worker_connection_idle`       | Gauge   | The current number of idle client connection                             | `id`, `pid` |
| `nginxplus_worker_http_requests_total`   | Counter | The total number of client requests received                             | `id`, `pid` |
| `nginxplus_worker_http_requests_current` | Gauge   | The current number of client requests that are currently being processed | `id`, `pid` |

Connect to the `/metrics` page of the running exporter to see the complete list of metrics along with their
descriptions. Note: to see server zones related metrics you must configure [status
zones](https://nginx.org/en/docs/http/ngx_http_api_module.html#status_zone) and to see upstream related metrics you
must configure upstreams with a [shared memory zone](https://nginx.org/en/docs/http/ngx_http_upstream_module.html#zone).

## Troubleshooting

The exporter logs errors to the standard output. When using Docker, if the exporter doesn’t work as expected, check its
logs using [docker logs](https://docs.docker.com/engine/reference/commandline/logs/) command.

## Releases

### Docker images

We publish the Docker image on [DockerHub](https://hub.docker.com/r/nginx/nginx-prometheus-exporter/),
[GitHub Container](https://github.com/nginx/nginx-prometheus-exporter/pkgs/container/nginx-prometheus-exporter),
[Amazon ECR Public Gallery](https://gallery.ecr.aws/nginx/nginx-prometheus-exporter) and
[Quay.io](https://quay.io/repository/nginx/nginx-prometheus-exporter).

As an alternative, you can choose the _edge_ version built from the [latest commit](https://github.com/nginx/nginx-prometheus-exporter/commits/main)
from the main branch. The edge version is useful for experimenting with new features that are not yet published in a
stable release.

### Binaries

We publish the binaries for multiple Operating Systems and architectures on the GitHub [releases page](https://github.com/nginx/nginx-prometheus-exporter/releases).

### Homebrew

You can add the NGINX homebrew tap with

```console
brew tap nginx/tap
```

and then install the formula with

```console
brew install nginx-prometheus-exporter
```

### Snap

You can install the NGINX Prometheus Exporter from the [Snap Store](https://snapcraft.io/nginx-prometheus-exporter).

```console
snap install nginx-prometheus-exporter
```

### Scoop

You can add the NGINX Scoop bucket with

```console
scoop bucket add nginx https://github.com/nginx/scoop-bucket.git
```

and then install the package with

```console
scoop install nginx-prometheus-exporter
```

### Nix

First include NUR in your packageOverrides as explained in the [NUR documentation](https://github.com/nix-community/NUR#installation).

Then you can use the exporter with the following command:

```console
nix-shell --packages nur.repos.nginx.nginx-prometheus-exporter
```

or install it with:

```console
nix-env -f '<nixpkgs>' -iA nur.repos.nginx.nginx-prometheus-exporter
```

## Building the Exporter

You can build the exporter using the provided Makefile. Before building the exporter, make sure the following software
is installed on your machine:

- make
- git
- Docker for building the container image
- Go for building the binary

### Building the Docker Image

To build the Docker image with the exporter, run:

```console
make container
```

Note: go is not required, as the exporter binary is built in a Docker container. See the [Dockerfile](build/Dockerfile).

### Building the Binary

To build the binary, run:

```console
make
```

Note: the binary is built for the OS/arch of your machine. To build binaries for other platforms, see the
[Makefile](Makefile).

The binary is built with the name `nginx-prometheus-exporter`.

## Grafana Dashboard

The official Grafana dashboard is provided with the exporter for NGINX. Check the [Grafana
Dashboard](./grafana/README.md) documentation for more information.

## SBOM (Software Bill of Materials)

We generate SBOMs for the binaries and the Docker image.

### Binaries

The SBOMs for the binaries are available in the releases page. The SBOMs are generated using
[syft](https://github.com/anchore/syft) and are available in SPDX format.

### Docker Image

The SBOM for the Docker image is available in the
[DockerHub](https://hub.docker.com/r/nginx/nginx-prometheus-exporter),
[GitHub Container registry](https://github.com/nginx/nginx-prometheus-exporter/pkgs/container/nginx-prometheus-exporter),
[Amazon ECR Public Gallery](https://gallery.ecr.aws/nginx/nginx-prometheus-exporter) and
[Quay.io](https://quay.io/repository/nginx/nginx-prometheus-exporter) repositories. The SBOMs are generated using
[syft](https://github.com/anchore/syft) and stored as an attestation in the image manifest.

For example to retrieve the SBOM for `linux/amd64` from Docker Hub and analyze it using
[grype](https://github.com/anchore/grype) you can run the following command:

```console
docker buildx imagetools inspect nginx/nginx-prometheus-exporter:edge --format '{{ json (index .SBOM "linux/amd64").SPDX }}' | grype
```

## Provenance

We generate provenance for the Docker image and it's available in the
[DockerHub](https://hub.docker.com/r/nginx/nginx-prometheus-exporter),
[GitHub Container registry](https://github.com/nginx/nginx-prometheus-exporter/pkgs/container/nginx-prometheus-exporter),
[Amazon ECR Public Gallery](https://gallery.ecr.aws/nginx/nginx-prometheus-exporter) and
[Quay.io](https://quay.io/repository/nginx/nginx-prometheus-exporter) repositories, stored as an attestation in the
image manifest.

For example to retrieve the provenance for `linux/amd64` from Docker Hub you can run the following command:

```console
docker buildx imagetools inspect nginx/nginx-prometheus-exporter:edge --format '{{ json (index .Provenance "linux/amd64").SLSA }}'
```

## Contacts

We’d like to hear your feedback! If you have any suggestions or experience issues with the NGINX Prometheus Exporter,
please create an issue or send a pull request on GitHub. You can contact us directly via <integrations@nginx.com> or on
the [NGINX Community Slack](https://nginxcommunity.slack.com/channels/nginx-prometheus-exporter) in the
`#nginx-prometheus-exporter` channel.

## Contributing

If you'd like to contribute to the project, please read our [Contributing guide](CONTRIBUTING.md).

## Support

The commercial support is available for NGINX Plus customers when the NGINX Prometheus Exporter is used with NGINX
Ingress Controller.

## License

[Apache License, Version 2.0](LICENSE).
