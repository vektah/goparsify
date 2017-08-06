go build
profile.exe -cpuprofile cpu.out
go tool pprof --inuse_objects profile.exe cpu.out
