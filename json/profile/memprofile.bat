go build
profile.exe -memprofile mem.out
go tool pprof profile.exe mem.out
