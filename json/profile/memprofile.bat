go build
profile.exe -memprofile mem.out
go tool pprof --inuse_objects profile.exe mem.out
