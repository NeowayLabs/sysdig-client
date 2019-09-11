#!/usr/bin/env bash

set -o errexit
set -o nounset

timeout=50s
if [ $# -eq 0 ]
  then
        for d in $(go list ./...); do
            go test -tags integration -timeout $timeout -race $d
        done
        exit
fi

pkg=$1
testname=$2
echo "Running test pkg:" $pkg " name: " $testname
go test -tags integration -timeout $timeout -race $pkg --run $testname
