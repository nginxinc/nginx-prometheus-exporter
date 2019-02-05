# Changelog

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