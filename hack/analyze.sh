#!/bin/bash

set -o errexit
set -o nounset

echo "performing static analysis on the code"

golangci-lint run

echo "done"
