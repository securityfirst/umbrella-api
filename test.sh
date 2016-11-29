#!/bin/bash
. env.sh
echo "Using: $FIPS"
go test -v $(glide nv)