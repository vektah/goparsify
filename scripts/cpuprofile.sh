#!/bin/bash

set -eu

go build ./json/profile/json.go
./json.exe -cpuprofile cpu.out
go tool pprof json.exe cpu.out
