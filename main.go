package main

import (
	"fmt"
	"integral/config"
	"integral/shell"
	"log"
	"os"
)

var logger = log.New(os.Stderr, fmt.Sprintf("\033[%dmError:\033[%dm ", 31, 39), 1)
var sh = shell.Shell{}

const version = "v0.3.0"

func main() {
	cfg := getConfig()

	if len(os.Args) < 2 {
		logger.Fatalln("Not enough arguments!")
	}
	switch os.Args[1] {
	case "transient":
		fmt.Print(shell.Raw.Fg(string(cfg.Line.Symbols[3]), cfg.Line.Color))
	case "render":
		var err error
		sh, err = shell.GetShell(os.Args[2])
		if err != nil {
			logger.Fatal(err)
		}
		render(cfg)
	case "init":
		var err error
		sh, err = shell.GetShell(os.Args[2])
		if err != nil {
			logger.Fatal(err)
		}
		sh.Init()
	case "version":
		fmt.Println(version)
	default:
		logger.Fatalln("Unknown command")
	}
}

func getConfig() *config.PromptConfig {
	return config.GetDefault()
}
