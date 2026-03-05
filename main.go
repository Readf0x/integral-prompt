package main

import (
	"fmt"
	"integral/config"
	"integral/shell"
	"log"
	"os"
)

var logger = log.New(os.Stderr, fmt.Sprintf("\033[%dmError:\033[%dm ", 31, 39), 1)

const version = "v0.3.0"

func main() {
	cfg := getConfig()

	if len(os.Args) < 2 {
		logger.Fatalln("Not enough arguments!")
	}
	switch os.Args[1] {
	case "transient":
		fmt.Print(shell.Fg(string(cfg.Line.Symbols[3]), cfg.Line.Color))
	case "render":
		render(cfg)
	case "init":
		shell.Init()
	case "version":
		fmt.Println(version)
	default:
		logger.Fatalln("Unknown command")
	}
}

func getConfig() *config.PromptConfig {
	return config.GetDefault()
}
