package main

import (
	"encoding/json"
	"fmt"
	"integral/config"
	"integral/shell"
	"log"
	"os"
)

var logger = log.New(os.Stderr, fmt.Sprintf("\033[%dmError:\033[%dm ", 31, 39), 1)
var sh = shell.Shell{}

//go:generate go run gen.go
var VersionString = "%s, built from commit %s"

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
		fmt.Printf(VersionString+"\n", Version, Commit)
	case "char":
		fmt.Println([]rune(os.Args[2]))
	case "config":
		if b, err := json.Marshal(cfg); err == nil {
			fmt.Printf("%s", b)
		}
	default:
		logger.Fatalln("Unknown command")
	}
}

func getConfig() *config.PromptConfig {
	c := os.Getenv("XDG_CONFIG_HOME")
	if c == "" {
		c = os.Getenv("HOME") + "/.config"
	}
	cfg := config.LoadConfig([]string{
		os.Getenv("HOME") + "/.integralrc",
		c + "/integralrc",
		c + "/integralrc.json",
		c + "/integral/rc",
		c + "/integral/rc.json",
	})
	return cfg
}
