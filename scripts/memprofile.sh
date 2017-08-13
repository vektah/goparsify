#!/bin/bash

set -eu

go build ./json/profile/json.go
./json.exe -memprofile mem.out
go tool pprof json.exe mem.out
