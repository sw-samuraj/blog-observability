#!/bin/sh -eu

LOKI_HOME="${HOME}/dev/loki"
PROMETHEUS_HOME="${HOME}/dev/prometheus"
PWD=$(pwd)

./kill-backends.sh

echo "Starting Loki in the background..."
"${LOKI_HOME}/loki-linux-amd64" -config.file "${PWD}/_config/loki/loki-local-config.yaml" &

echo "Starting Promtail in the background..."
"${LOKI_HOME}/promtail-linux-amd64" -config.file "${PWD}/_config/loki/promtail-local-config.yaml" &

echo "Starting Prometheus in the background..."
"${PROMETHEUS_HOME}/prometheus" --config.file "${PWD}/_config/prometheus/prometheus.yml" --storage.tsdb.path="_tmp/prometheus/data/" &
