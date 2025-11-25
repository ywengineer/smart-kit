#!/bin/bash
RUN_NAME=payment
mkdir -p output/bin
cp script/* output 2>/dev/null
chmod +x output/bootstrap.sh
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o output/bin/${RUN_NAME}