#!/bin/bash
realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}
root=$(dirname `realpath "$0"`)
export GOPATH=$root:$GOPATH
export FIPS=$root/fips.csv