#!/usr/bin/env bash
set -xe

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

version=$1

if [ -z $version ]; then
    version=v0.1
fi

go build -o conntrack_exporter conntrack_exporter.go
chmod +x conntrack_exporter

tar -zcvf conntrack_exporter-linux-amd64-${version}.tar.gz \
  conntrack_exporter config/conntrack_exporter.yaml config/conntrack_exporter.service README.md
