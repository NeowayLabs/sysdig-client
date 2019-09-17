#!/usr/bin/env bash

set -o errexit
set -o nounset

timeout=15s
if [ $# -eq 0 ]
  then
        for d in $(go list ./...); do
            go test -tags unit -timeout $timeout -race $d
        done
        exit
fi

pkg=$1
testname=$2
echo "Running test pkg:" $pkg " name: " $testname
go test -tags unit -timeout $timeout -race $pkg --run $testname
