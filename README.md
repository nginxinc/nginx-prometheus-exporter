[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nginxinc/nginx-prometheus-exporter/badge)](https://api.securityscorecards.dev/projects/github.com/nginxinc/nginx-prometheus-exporter) [![CI](https://github.com/nginxinc/nginx-prometheus-exporter/workflows/Continuous%20Integration/badge.svg)](https://github.com/nginxinc/nginx-prometheus-exporter/actions?query=workflow%3A%22Continuous+Integration%22) [![FOSSA Status](https://app.fossa.com/api/projects/custom%2B5618%2Fgithub.com%2Fnginxinc%2Fnginx-prometheus-exporter.svg?type=shield)](https://app.fossa.com/projects/custom%2B5618%2Fgithub.com%2Fnginxinc%2Fnginx-prometheus-exporter?ref=badge_shield) [![Go Report Card](https://goreportcard.com/badge/github.com/nginxinc/nginx-prometheus-exporter)](https://goreportcard.com/report/github.com/nginxinc/nginx-prometheus-exporter) ![GitHub all releases](https://img.shields.io/github/downloads/nginxinc/nginx-prometheus-exporter/total?logo=github) ![GitHub release (latest by SemVer)](https://img.shields.io/github/downloads/nginxinc/nginx-prometheus-exporter/latest/total?sort=semver&logo=github) [![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nginxinc/nginx-prometheus-exporter?logo=github&sort=semver)](https://github.com/nginxinc/nginx-prometheus-exporter/releases/latest) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nginxinc/nginx-prometheus-exporter?logo=go) [![Docker Pulls](https://img.shields.io/docker/pulls/nginx/nginx-prometheus-exporter?logo=docker&logoColor=white)](https://hub.docker.com/r/nginx/nginx-prometheus-exporter) ![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/nginx/nginx-prometheus-exporter?logo=docker&logoColor=white&sort=semver) [![Slack](https://img.shields.io/badge/slack-%23nginx--prometheus--exporter-green?logo=slack)](https://nginxcommunity.slack.com/channels/nginx-prometheus-exporter)

# NGINX Prometheus Exporter

NGINX Prometheus exporter makes it possible to monitor NGINX or NGINX Plus using Prometheus.

## Overview

[NGINX](https://nginx.org) exposes a handful of metrics via the [stub_status page](https://nginx.org/en/docs/http/ngx_http_stub_status_module.html#stub_status). [NGINX Plus](https://www.nginx.com/products/nginx/) provides a richer set of metrics via the [API](https://nginx.org/en/docs/http/ngx_http_api_module.html) and the [monitoring dashboard](https://www.nginx.com/products/nginx/live-activity-monitoring/). NGINX Prometheus exporter fetches the metrics from a single NGINX or NGINX Plus, converts the metrics into appropriate Prometheus metrics types and finally exposes them via an HTTP server to be collected by [Prometheus](https://prometheus.io/).

## Getting Started

In this section, we show how to quickly run NGINX Prometheus Exporter for NGINX or NGINX Plus.

### A Note about NGINX Ingress Controller

If you’d like to use the NGINX Prometheus Exporter with [NGINX Ingress Controller](https://github.com/nginxinc/kubernetes-ingress/) for Kubernetes, see [this doc](https://docs.nginx.com/nginx-ingress-controller/logging-and-monitoring/prometheus/) for the installation instructions.

### Prerequisites

We assume that you have already installed Prometheus and NGINX or NGINX Plus. Additionally, you need to:
* Expose the built-in metrics in NGINX/NGINX Plus:
    * For NGINX, expose the [stub_status page](https://nginx.org/en/docs/http/ngx_http_stub_status_module.html#stub_status) at `/stub_status` on port `8080`.
    * For NGINX Plus, expose the [API](https://nginx.org/en/docs/http/ngx_http_api_module.html#api) at `/api` on port `8080`.
* Configure Prometheus to scrape metrics from the server with the exporter. Note that the default scrape port of the exporter is `9113` and the default metrics path -- `/metrics`.

### Running the Exporter in a Docker Container

To start the exporter we use the [docker run](https://docs.docker.com/engine/reference/run/) command.

* To export NGINX metrics, run:
    ```
    $ docker run -p 9113:9113 nginx/nginx-prometheus-exporter:0.10.0 -nginx.scrape-uri=http://<nginx>:8080/stub_status
    ```
    where `<nginx>` is the IP address/DNS name, through which NGINX is available.

* To export NGINX Plus metrics, run:
    ```
    $ docker run -p 9113:9113 nginx/nginx-prometheus-exporter:0.10.0 -nginx.plus -nginx.scrape-uri=http://<nginx-plus>:8080/api
    ```
    where `<nginx-plus>` is the IP address/DNS name, through which NGINX Plus is available.

### Running the Exporter Binary

* To export NGINX metrics, run:
    ```
    $ nginx-prometheus-exporter -nginx.scrape-uri=http://<nginx>:8080/stub_status
    ```
    where `<nginx>` is the IP address/DNS name, through which NGINX is available.

* To export NGINX Plus metrics:
    ```
    $ nginx-prometheus-exporter -nginx.plus -nginx.scrape-uri=http://<nginx-plus>:8080/api
    ```
    where `<nginx-plus>` is the IP address/DNS name, through which NGINX Plus is available.

* To export and scrape NGINX metrics with unix domain sockets, run:
    ```
    $ nginx-prometheus-exporter -nginx.scrape-uri=unix:<nginx>:/stub_status -web.listen-address=unix:/path/to/socket.sock
    ```
    where `<nginx>` is the path to unix domain socket, through which NGINX stub status is available.

**Note**. The `nginx-prometheus-exporter` is not a daemon. To run the exporter as a system service (daemon), configure the init system of your Linux server (such as systemd or Upstart) accordingly. Alternatively, you can run the exporter in a Docker container.

## Usage

### Command-line Arguments

```
Usage of ./nginx-prometheus-exporter:
  -nginx.plus
        Start the exporter for NGINX Plus. By default, the exporter is started for NGINX. The default value can be overwritten by NGINX_PLUS environment variable.
  -nginx.retries int
        A number of retries the exporter will make on start to connect to the NGINX stub_status page/NGINX Plus API before exiting with an error. The default value can be overwritten by NGINX_RETRIES environment variable.
  -nginx.retry-interval duration
        An interval between retries to connect to the NGINX stub_status page/NGINX Plus API on start. The default value can be overwritten by NGINX_RETRY_INTERVAL environment variable. (default 5s)
  -nginx.scrape-uri string
        A URI or unix domain socket path for scraping NGINX or NGINX Plus metrics.
        For NGINX, the stub_status page must be available through the URI. For NGINX Plus -- the API. The default value can be overwritten by SCRAPE_URI environment variable. (default "http://127.0.0.1:8080/stub_status")
  -nginx.ssl-ca-cert string
        Path to the PEM encoded CA certificate file used to validate the servers SSL certificate. The default value can be overwritten by SSL_CA_CERT environment variable.
  -nginx.ssl-client-cert string
        Path to the PEM encoded client certificate file to use when connecting to the server. The default value can be overwritten by SSL_CLIENT_CERT environment variable.
  -nginx.ssl-client-key string
        Path to the PEM encoded client certificate key file to use when connecting to the server. The default value can be overwritten by SSL_CLIENT_KEY environment variable.
  -nginx.ssl-verify
        Perform SSL certificate verification. The default value can be overwritten by SSL_VERIFY environment variable. (default true)
  -nginx.timeout duration
        A timeout for scraping metrics from NGINX or NGINX Plus. The default value can be overwritten by TIMEOUT environment variable. (default 5s)
  -prometheus.const-labels value
        A comma separated list of constant labels that will be used in every metric. Format is label1=value1,label2=value2... The default value can be overwritten by CONST_LABELS environment variable.
  -web.listen-address string
        An address or unix domain socket path to listen on for web interface and telemetry. The default value can be overwritten by LISTEN_ADDRESS environment variable. (default ":9113")
  -web.telemetry-path string
        A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable. (default "/metrics")
  -web.secured-metrics
        Expose metrics using https. The default value can be overwritten by SECURED_METRICS variable.  (default false)
  -web.ssl-server-cert string
        Path to the PEM encoded certificate for the nginx-exporter metrics server(when web.secured-metrics=true). The default value can be overwritten by SSL_SERVER_CERT variable.
  -web.ssl-server-key string
        Path to the PEM encoded key for the nginx-exporter metrics server (when web.secured-metrics=true). The default value can be overwritten by SSL_SERVER_KEY variable.
  -version
        Display the NGINX exporter version. (default false)
```

## Exported Metrics

### Common metrics:
Name | Type | Description | Labels
----|----|----|----|
`nginxexporter_build_info` | Gauge | Shows the exporter build information. | `gitCommit`, `version` |
`nginx_up` | Gauge | Shows the status of the last metric scrape: `1` for a successful scrape and `0` for a failed one | [] |

### Metrics for NGINX OSS:
#### [Stub status metrics](https://nginx.org/en/docs/http/ngx_http_stub_status_module.html)
Name | Type | Description | Labels
----|----|----|----|
`nginx_connections_accepted` | Counter | Accepted client connections. | [] |
`nginx_connections_active` | Gauge | Active client connections. | [] |
`nginx_connections_handled` | Counter | Handled client connections. | [] |
`nginx_connections_reading` | Gauge | Connections where NGINX is reading the request header. | [] |
`nginx_connections_waiting` | Gauge | Idle client connections. | [] |
`nginx_connections_writing` | Gauge | Connections where NGINX is writing the response back to the client. | [] |
`nginx_http_requests_total` | Counter | Total http requests. | [] |

### Metrics for NGINX Plus:
#### [Connections](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_connections)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_connections_accepted` | Counter | Accepted client connections | [] |
`nginxplus_connections_active` | Gauge | Active client connections | [] |
`nginxplus_connections_dropped` | Counter | Dropped client connections dropped | [] |
`nginxplus_connections_idle` | Gauge | Idle client connections | [] |

#### [HTTP](https://nginx.org/en/docs/http/ngx_http_api_module.html#http_)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_http_requests_total` | Counter | Total http requests | [] |
`nginxplus_http_requests_current` | Gauge | Current http requests | [] |

#### [SSL](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_ssl_object)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_ssl_handshakes` | Counter | Successful SSL handshakes | [] |
`nginxplus_ssl_handshakes_failed` | Counter | Failed SSL handshakes | [] |
`nginxplus_ssl_session_reuses` | Counter | Session reuses during SSL handshake | [] |

#### [HTTP Server Zones](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_server_zone)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_server_zone_processing` | Gauge | Client requests that are currently being processed | `server_zone` |
`nginxplus_server_zone_requests` | Counter | Total client requests | `server_zone` |
`nginxplus_server_zone_responses` | Counter | Total responses sent to clients | `code` (the response status code. The values are: `1xx`, `2xx`, `3xx`, `4xx` and `5xx`), `server_zone` |
`nginxplus_server_zone_responses_codes` | Counter | Total responses sent to clients by code | `code` (the response status code. The possible values are [here](https://www.nginx.com/resources/wiki/extending/api/http/)), `server_zone` |
`nginxplus_server_zone_discarded` | Counter | Requests completed without sending a response | `server_zone` |
`nginxplus_server_zone_received` | Counter | Bytes received from clients | `server_zone` |
`nginxplus_server_zone_sent` | Counter | Bytes sent to clients | `server_zone` |
`nginxplus_server_ssl_handshakes` | Counter | Successful SSL handshakes | `server_zone` |
`nginxplus_server_ssl_handshakes_failed` | Counter | Failed SSL handshakes | `server_zone` |
`nginxplus_server_ssl_session_reuses` | Counter | Session reuses during SSL handshake | `server_zone` |

#### [Stream Server Zones](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_stream_server_zone)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_stream_server_zone_processing` | Gauge | Client connections that are currently being processed | `server_zone` |
`nginxplus_stream_server_zone_connections` | Counter | Total connections | `server_zone` |
`nginxplus_stream_server_zone_sessions` | Counter | Total sessions completed | `code` (the response status code. The values are: `2xx`, `4xx`, and `5xx`), `server_zone` |
`nginxplus_stream_server_zone_discarded` | Counter | Connections completed without creating a session | `server_zone` |
`nginxplus_stream_server_zone_received` | Counter | Bytes received from clients | `server_zone` |
`nginxplus_stream_server_zone_sent` | Counter | Bytes sent to clients | `server_zone` |
`nginxplus_stream_server_ssl_handshakes` | Counter | Successful SSL handshakes | `server_zone` |
`nginxplus_stream_server_ssl_handshakes_failed` | Counter | Failed SSL handshakes | `server_zone` |
`nginxplus_stream_server_ssl_session_reuses` | Counter | Session reuses during SSL handshake | `server_zone` |

#### [HTTP Upstreams](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_upstream)

> Note: for the `state` metric, the string values are converted to float64 using the following rule: `"up"` -> `1.0`, `"draining"` -> `2.0`, `"down"` -> `3.0`, `"unavail"` –> `4.0`, `"checking"` –> `5.0`, `"unhealthy"` -> `6.0`.

Name | Type | Description | Labels
----|----|----|----|
`nginxplus_upstream_server_state` | Gauge | Current state | `server`, `upstream` |
`nginxplus_upstream_server_active` | Gauge | Active connections | `server`, `upstream` |
`nginxplus_upstream_server_limit` | Gauge | Limit for connections which corresponds to the max_conns parameter of the upstream server. Zero value means there is no limit | `server`, `upstream` |
`nginxplus_upstream_server_requests` | Counter | Total client requests | `server`, `upstream` |
`nginxplus_upstream_server_responses` | Counter | Total responses sent to clients | `code` (the response status code. The values are: `1xx`, `2xx`, `3xx`, `4xx` and `5xx`), `server`, `upstream` |
`nginxplus_upstream_server_responses_codes` | Counter | Total responses sent to clients by code | `code` (the response status code. The possible values are [here](https://www.nginx.com/resources/wiki/extending/api/http/)), `server`, `upstream` |
`nginxplus_upstream_server_sent` | Counter | Bytes sent to this server | `server`, `upstream` |
`nginxplus_upstream_server_received` | Counter | Bytes received to this server | `server`, `upstream` |
`nginxplus_upstream_server_fails` | Counter | Number of unsuccessful attempts to communicate with the server | `server`, `upstream` |
`nginxplus_upstream_server_unavail` | Counter | How many times the server became unavailable for client requests (state 'unavail') due to the number of unsuccessful attempts reaching the max_fails threshold | `server`, `upstream` |
`nginxplus_upstream_server_header_time` | Gauge | Average time to get the response header from the server | `server`, `upstream` |
`nginxplus_upstream_server_response_time` | Gauge | Average time to get the full response from the server | `server`, `upstream` |
`nginxplus_upstream_server_health_checks_checks` | Counter | Total health check requests | `server`, `upstream` |
`nginxplus_upstream_server_health_checks_fails` | Counter | Failed health checks | `server`, `upstream` |
`nginxplus_upstream_server_health_checks_unhealthy` | Counter | How many times the server became unhealthy (state 'unhealthy') | `server`, `upstream` |
`nginxplus_upstream_server_ssl_handshakes` | Counter | Successful SSL handshakes | `server`, `upstream` |
`nginxplus_upstream_server_ssl_handshakes_failed` | Counter | Failed SSL handshakes | `server`, `upstream` |
`nginxplus_upstream_server_ssl_session_reuses` | Counter | Session reuses during SSL handshake | `server`, `upstream` |
`nginxplus_upstream_keepalives` | Gauge | Idle keepalive connections | `upstream` |
`nginxplus_upstream_zombies` | Gauge | Servers removed from the group but still processing active client requests | `upstream` |

#### [Stream Upstreams](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_stream_upstream)

> Note: for the `state` metric, the string values are converted to float64 using the following rule: `"up"` -> `1.0`, `"down"` -> `3.0`, `"unavail"` –> `4.0`, `"checking"` –> `5.0`, `"unhealthy"` -> `6.0`.

Name | Type | Description | Labels
----|----|----|----|
`nginxplus_stream_upstream_server_state` | Gauge | Current state | `server`, `upstream` |
`nginxplus_stream_upstream_server_active` | Gauge | Active connections | `server` , `upstream` |
`nginxplus_stream_upstream_server_limit` | Gauge | Limit for connections which corresponds to the max_conns parameter of the upstream server. Zero value means there is no limit | `server` , `upstream` |
`nginxplus_stream_upstream_server_connections` | Counter | Total number of client connections forwarded to this server | `server`, `upstream` |
`nginxplus_stream_upstream_server_connect_time` | Gauge | Average time to connect to the upstream server | `server`, `upstream`
`nginxplus_stream_upstream_server_first_byte_time` | Gauge | Average time to receive the first byte of data | `server`, `upstream` |
`nginxplus_stream_upstream_server_response_time` | Gauge | Average time to receive the last byte of data | `server`, `upstream` |
`nginxplus_stream_upstream_server_sent` | Counter | Bytes sent to this server | `server`, `upstream` |
`nginxplus_stream_upstream_server_received` | Counter | Bytes received from this server | `server`, `upstream` |
`nginxplus_stream_upstream_server_fails` | Counter | Number of unsuccessful attempts to communicate with the server | `server`, `upstream` |
`nginxplus_stream_upstream_server_unavail` | Counter | How many times the server became unavailable for client connections (state 'unavail') due to the number of unsuccessful attempts reaching the max_fails threshold | `server`, `upstream` |
`nginxplus_stream_upstream_server_health_checks_checks` | Counter | Total health check requests | `server`, `upstream` |
`nginxplus_stream_upstream_server_health_checks_fails` | Counter | Failed health checks | `server`, `upstream` |
`nginxplus_stream_upstream_server_health_checks_unhealthy` | Counter | How many times the server became unhealthy (state 'unhealthy') | `server`, `upstream` |
`nginxplus_stream_upstream_server_ssl_handshakes` | Counter | Successful SSL handshakes | `server`, `upstream` |
`nginxplus_stream_upstream_server_ssl_handshakes_failed` | Counter | Failed SSL handshakes | `server`, `upstream` |
`nginxplus_stream_upstream_server_ssl_session_reuses` | Counter | Session reuses during SSL handshake | `server`, `upstream` |
`nginxplus_stream_upstream_zombies` | Gauge | Servers removed from the group but still processing active client connections | `upstream`|

#### [Stream Zone Sync](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_stream_zone_sync)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_stream_zone_sync_zone_records_pending` | Gauge | The number of records that need to be sent to the cluster | `zone` |
`nginxplus_stream_zone_sync_zone_records_total` | Gauge | The total number of records stored in the shared memory zone | `zone` |
`nginxplus_stream_zone_sync_zone_bytes_in` | Counter | Bytes received by this node | [] |
`nginxplus_stream_zone_sync_zone_bytes_out` | Counter | Bytes sent by this node | [] |
`nginxplus_stream_zone_sync_zone_msgs_in` | Counter | Total messages received by this node | [] |
`nginxplus_stream_zone_sync_zone_msgs_out` | Counter | Total messages sent by this node | [] |
`nginxplus_stream_zone_sync_zone_nodes_online` | Gauge | Number of peers this node is connected to | [] |

#### [Location Zones](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_location_zone)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_location_zone_requests` | Counter | Total client requests | `location_zone` |
`nginxplus_location_zone_responses` | Counter | Total responses sent to clients | `code` (the response status code. The values are: `1xx`, `2xx`, `3xx`, `4xx` and `5xx`), `location_zone` |
`nginxplus_location_zone_responses_codes` | Counter | Total responses sent to clients by code | `code` (the response status code. The possible values are [here](https://www.nginx.com/resources/wiki/extending/api/http/)), `location_zone` |
`nginxplus_location_zone_discarded` | Counter | Requests completed without sending a response | `location_zone` |
`nginxplus_location_zone_received` | Counter | Bytes received from clients | `location_zone` |
`nginxplus_location_zone_sent` | Counter | Bytes sent to clients | `location_zone` |

#### [Resolver](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_resolver_zone)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_resolver_name` | Counter | Total requests to resolve names to addresses | `resolver` |
`nginxplus_resolver_srv` | Counter | Total requests to resolve SRV records | `resolver` |
`nginxplus_resolver_addr` | Counter | Total requests to resolve addresses to names | `resolver` |
`nginxplus_resolver_noerror` | Counter | Total number of successful responses | `resolver` |
`nginxplus_resolver_formerr` | Counter | Total number of FORMERR responses | `resolver` |
`nginxplus_resolver_servfail` | Counter | Total number of SERVFAIL responses | `resolver` |
`nginxplus_resolver_nxdomain` | Counter | Total number of NXDOMAIN responses | `resolver` |
`nginxplus_resolver_notimp` | Counter | Total number of NOTIMP responses | `resolver` |
`nginxplus_resolver_refused` | Counter | Total number of REFUSED responses | `resolver` |
`nginxplus_resolver_timedout` | Counter | Total number of timed out request | `resolver` |
`nginxplus_resolver_unknown` | Counter | Total requests completed with an unknown error | `resolver`|

#### [HTTP Requests Rate Limiting](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_limit_req_zone)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_limit_request_passed` | Counter | Total number of requests that were neither limited nor accounted as limited | `zone` |
`nginxplus_limit_request_rejected` | Counter | Total number of requests that were rejected | `zone` |
`nginxplus_limit_request_delayed` | Counter | Total number of requests that were delayed | `zone` |
`nginxplus_limit_request_rejected_dry_run` | Counter | Total number of requests accounted as rejected in the dry run mode | `zone` |
`nginxplus_limit_request_delayed_dry_run` | Counter | Total number of requests accounted as delayed in the dry run mode | `zone` |

#### [HTTP Connections Limiting](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_http_limit_conn_zone)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_limit_connection_passed` | Counter | Total number of connections that were neither limited nor accounted as limited | `zone` |
`nginxplus_limit_connection_rejected` | Counter | Total number of connections that were rejected | `zone` |
`nginxplus_limit_connection_rejected_dry_run` | Counter | Total number of connections accounted as rejected in the dry run mode | `zone` |


#### [Stream Connections Limiting](https://nginx.org/en/docs/http/ngx_http_api_module.html#def_nginx_stream_limit_conn_zone)
Name | Type | Description | Labels
----|----|----|----|
`nginxplus_stream_limit_connection_passed` | Counter | Total number of connections that were neither limited nor accounted as limited | `zone` |
`nginxplus_stream_limit_connection_rejected` | Counter | Total number of connections that were rejected | `zone` |
`nginxplus_stream_limit_connection_rejected_dry_run` | Counter | Total number of connections accounted as rejected in the dry run mode | `zone` |

Connect to the `/metrics` page of the running exporter to see the complete list of metrics along with their descriptions. Note: to see server zones related metrics you must configure [status zones](https://nginx.org/en/docs/http/ngx_http_status_module.html#status_zone) and to see upstream related metrics you must configure upstreams with a [shared memory zone](https://nginx.org/en/docs/http/ngx_http_upstream_module.html#zone).

## Troubleshooting

The exporter logs errors to the standard output. When using Docker, if the exporter doesn’t work as expected, check its logs using [docker logs](https://docs.docker.com/engine/reference/commandline/logs/) command.

## Releases

### Docker images
We publish the Docker image on [DockerHub](https://hub.docker.com/r/nginx/nginx-prometheus-exporter/), [GitHub Container](https://github.com/nginxinc/nginx-prometheus-exporter/pkgs/container/nginx-prometheus-exporter), [Amazon ECR Public Gallery](https://gallery.ecr.aws/nginx/nginx-prometheus-exporter) and [Quay.io](https://quay.io/repository/nginx/nginx-prometheus-exporter).

As an alternative, you can choose the *edge* version built from the [latest commit](https://github.com/nginxinc/nginx-prometheus-exporter/commits/main) from the main branch. The edge version is useful for experimenting with new features that are not yet published in a stable release.

### Binaries
We publish the binaries for multiple Operating Systems and architectures on the GitHub [releases page](https://github.com/nginxinc/nginx-prometheus-exporter/releases).

### Homebrew
You can add the NGINX homebrew tap with
```
$ brew tap nginxinc/tap
```
and then install the formula with
```
$ brew install nginx-prometheus-exporter
```

## Building the Exporter

You can build the exporter using the provided Makefile. Before building the exporter, make sure the following software is installed on your machine:
* make
* git
* Docker for building the container image
* Go for building the binary

### Building the Docker Image

To build the Docker image with the exporter, run:
```
$ make container
```

Note: go is not required, as the exporter binary is built in a Docker container. See the [Dockerfile](build/Dockerfile).

### Building the Binary

To build the binary, run:
```
$ make
```

Note: the binary is built for the OS/arch of your machine. To build binaries for other platforms, see the [Makefile](Makefile).

The binary is built with the name `nginx-prometheus-exporter`.

## Grafana Dashboard
The official Grafana dashboard is provided with the exporter for NGINX. Check the [Grafana Dashboard](./grafana/README.md) documentation for more information.

## SBOM (Software Bill of Materials)

We generate SBOMs for the binaries and the Docker image.

### Binaries

The SBOMs for the binaries are available in the releases page. The SBOMs are generated using [syft](https://github.com/anchore/syft) and are available in SPDX format.

### Docker Image

The SBOM for the Docker image is available in the [DockerHub](https://hub.docker.com/r/nginx/nginx-prometheus-exporter), [GitHub Container registry](https://github.com/nginxinc/nginx-prometheus-exporter/pkgs/container/nginx-prometheus-exporter), [Amazon ECR Public Gallery](https://gallery.ecr.aws/nginx/nginx-prometheus-exporter) and [Quay.io](https://quay.io/repository/nginx/nginx-prometheus-exporter) repositories. The SBOMs are generated using [syft](https://github.com/anchore/syft) and stored as an attestation in the image manifest.

For example to retrieve the SBOM for `linux/amd64` from Docker Hub and analyze it using [grype](https://github.com/anchore/grype) you can run the following command:
```
$ docker buildx imagetools inspect nginx/nginx-prometheus-exporter:edge --format '{{ json (index .SBOM "linux/amd64").SPDX }}' | grype
```

## Contacts

We’d like to hear your feedback! If you have any suggestions or experience issues with the NGINX Prometheus Exporter, please create an issue or send a pull request on GitHub.
You can contact us directly via integrations@nginx.com or on the [NGINX Community Slack](https://nginxcommunity.slack.com/channels/nginx-prometheus-exporter) in the `#nginx-prometheus-exporter` channel.

## Contributing

If you'd like to contribute to the project, please read our [Contributing guide](CONTRIBUTING.md).

## Support

The commercial support is available for NGINX Plus customers when the NGINX Prometheus Exporter is used with NGINX Ingress Controller.

## License

[Apache License, Version 2.0](LICENSE).
