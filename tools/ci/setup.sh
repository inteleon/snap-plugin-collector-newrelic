#!/usr/bin/env bash

set -e

GOOGLE_SRC="${GOPATH}/src/google.golang.org"
mkdir -p $GOOGLE_SRC
git clone https://github.com/grpc/grpc-go "${GOOGLE_SRC}/grpc"
cd "${GOOGLE_SRC}/grpc"
git checkout v1.0.3

YFRONTO_SRC="${GOPATH}/src/github.com/yfronto"
mkdir -p $YFRONTO_SRC
git clone https://github.com/inteleon/newrelic.git "${YFRONTO_SRC}/newrelic"
cd "${YFRONTO_SRC}/newrelic"
git checkout feature/support-for-component-metrics
