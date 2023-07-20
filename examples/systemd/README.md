# NGINX Prometheus Exporter with systemd-activated socket

This example shows how to run NGINX Prometheus Exporter with systemd-activated socket.

## Prerequisites

- Linux machine with [systemd](https://www.freedesktop.org/wiki/Software/systemd/).
- NGINX Prometheus Exporter binary in `/usr/local/bin/nginx-prometheus-exporter` or a location of your choice. See the
  [main README](../../README.md) for installation instructions.
- NGINX or NGINX Plus running on the same machine.

## Customization

Modify `nginx_exporter.service` and `nginx_exporter.socket` to match your environment.

The default configuration assumes that NGINX Prometheus Exporter binary is located in
`/usr/local/bin/nginx-prometheus-exporter`.

The `ExecStart` directive has the flag `--web.systemd-socket` which tells the exporter to listen on the socket specified
in the `nginx_exporter.socket` file.

The `ListenStream` directive in `nginx_exporter.socket` specifies the socket to listen on. The default configuration
uses `9113` port, but the address can be written in various formats, for example `/run/nginx_exporter.sock`. To see the
full list of supported formats, run `man systemd.socket`.

## Installation

1. Copy `nginx_exporter.service` and `nginx_exporter.socket` to `/etc/systemd/system/`
2. Run `systemctl daemon-reload`
3. Run `systemctl start nginx_exporter`
4. Run `systemctl status nginx_exporter` to check the status of the service

## Verification

1. Run `curl http://localhost:9113/metrics` to see the metrics exposed by the exporter.
