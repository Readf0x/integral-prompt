package main

import (
	"integral/config"
	"integral/modules"
	"os"
	"runtime/pprof"
)

func main() {
	module := &modules.GitModule{}
	cfg := config.GetDefault()

	f, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	module.Initialize(cfg)
}

