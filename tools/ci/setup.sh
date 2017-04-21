#!/usr/bin/env bash

set -e

GOOGLE_SRC="${GOPATH}/src/google.golang.org"
mkdir -p $GOOGLE_SRC
git clone https://github.com/grpc/grpc-go "${GOOGLE_SRC}/grpc"
cd "${GOOGLE_SRC}/grpc"
git checkout v1.0.3
