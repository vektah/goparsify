#!/bin/bash

if [ $# != 1 ] ; then
    echo Run this in a directory containing benchmarks and pass it the name of a benchmark. It will dump allocations out to trace.log
    exit
fi

set -eu

go test -c

GODEBUG=allocfreetrace=1 ./$(basename $(pwd)).test.exe  -test.run=none -test.bench=$1 -test.benchmem -test.benchtime=1ns 2> >(sed -n '/benchmark.go:75/,$p' > trace.log)
