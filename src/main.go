package main

import (
	"fmt"
	config "integral/conf"
	"integral/shell"
	"log"
	"os"
)

var logger = log.New(os.Stderr, fmt.Sprintf("\033[%dmError:\033[%dm ", 31, 39), 1)

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
	}
}

func getConfig() *config.PromptConfig {
	return config.GetDefault()
}
