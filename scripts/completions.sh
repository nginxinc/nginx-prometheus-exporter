#!/bin/sh
set -e
rm -rf completions
mkdir completions
for shell in bash zsh; do
    go run . --completion-script-$shell >completions/nginx-prometheus-exporter.$shell
done
