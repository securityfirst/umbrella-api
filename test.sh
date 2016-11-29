#!/bin/bash
. env.sh
echo "Using: $FIPS"
cd $root/src/umbrella
go test -v $(glide nv)