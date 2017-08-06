package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/vektah/goparsify/json"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)

		defer func() {
			pprof.StopCPUProfile()
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}()
	}
	if *memprofile != "" {
		runtime.MemProfileRate = 1
	}

	for i := 0; i < 10000; i++ {
		_, err := json.Unmarshal(`{"true":true, "false":false, "null": null}`)
		if err != nil {
			panic(err)
		}
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
