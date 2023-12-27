#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run . --help-man | gzip -c -9 >manpages/nginx-prometheus-exporter.1.gz
