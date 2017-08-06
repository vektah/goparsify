go build
profile.exe -cpuprofile cpu.out
go tool pprof profile.exe cpu.out
