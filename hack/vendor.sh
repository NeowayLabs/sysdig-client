#!/usr/bin/env bash

set -o nounset
set -o errexit

rm -Rf vendor Godeps Gopkg.*

dep init
