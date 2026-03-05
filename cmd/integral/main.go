package main

import (
	"encoding/json"
	"fmt"
	"integral/config"
	"integral/modules"
	"integral/shell"
	"log"
	"os"
)

var logger = log.New(os.Stderr, fmt.Sprintf("\033[%dmError:\033[%dm ", 31, 39), 1)
var sh = shell.Shell{}

//go:generate go run gen.go
var VersionString = "%s, built from commit %s"

func main() {
	modules.Logger = logger
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
		modules.Sh = sh
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
	case "config":
		if b, err := json.Marshal(cfg); err == nil {
			fmt.Printf("%s", b)
		}
	case "--help":
		fallthrough
	case "-h":
		fallthrough
	case "help":
		printHelp()
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

func printHelp() {
	fmt.Print(
		"usage: integral {transient,render,init,version,config,help} ...\n\n" +
		"CLI for integral prompt\n\n" +
		"positional arguments:\n" +
		"  {transient,render,init,version,config,help}\n" +
		"    transient\trender transient prompt\n" +
		"    render\t\trender full prompt\n" +
		"    init\t\tprint shell init script\n" +
		"    version\t\tprint version info\n" +
		"    config\t\tprint config as json\n" +
		"    help\t\tshow this menu\n",
	)
	os.Exit(0)
}
