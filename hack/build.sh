#!/bin/bash

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -o ./cmd/poc-kafka/poc-kafka ./cmd/poc-kafka/main.go
