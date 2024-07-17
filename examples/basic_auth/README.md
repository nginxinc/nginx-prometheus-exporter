# NGINX Prometheus Exporter with Web Configuration for Basic Authentication

This example shows how to run NGINX Prometheus Exporter with web configuration. In this folder you will find an example
configuration `web-config.yml` that enables basic authentication. It is configured to have a single user `alice` with
password `password`.

The full documentation for the web configuration can be found
[here](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [Prerequisites](#prerequisites)
- [Running NGINX Prometheus Exporter with Web Configuration in Basic Authentication mode](#running-nginx-prometheus-exporter-with-web-configuration-in-basic-authentication-mode)
- [Verification](#verification)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Prerequisites

- NGINX Prometheus Exporter binary. See the [main README](../../README.md) for installation instructions.
- NGINX or NGINX Plus running on the same machine.

## Running NGINX Prometheus Exporter with Web Configuration in Basic Authentication mode

You can run NGINX Prometheus Exporter with web configuration in Basic Authentication mode using the following command:

```console
nginx-prometheus-exporter --web.config.file=web-config.yml --nginx.scrape-uri="http://127.0.0.1:8080/stub_status"
```

Depending on your environment, you may need to specify the full path to the binary or change the path to the web
configuration file.

## Verification

Run `curl -u alice:password http://localhost:9113/metrics` to see the metrics exposed by the exporter. Without the `-u`
flag, the request will fail with `401 Unauthorized`.
