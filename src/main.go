package main

import (
	"fmt"
	config "integral/conf"
	"log"
	"os"
	"strconv"

	"golang.org/x/term"
)

var logger = log.New(os.Stderr, fg("Error:", config.Yellow), 1)

func main() {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	prompt := render(config.GetDefault(), &defmodules)

	fmt.Println(fg("Width:", config.Yellow), bold(strconv.Itoa(width)))
	fmt.Println(fg("=========", config.BrightBlack))
	for _, line := range prompt {
		fmt.Println(line)
	}
	fmt.Println(fg("=========", config.BrightBlack))
}

func assemble(cfg *config.PromptConfig, size int, debug bool) []string {
	lines := make([]string, 0, *cfg.Length+2)

	// render right prompt
	// rc := make(chan string)

	// render main prompt
	render(cfg, cfg.Modules)

	return lines
}

func render(cfg *config.PromptConfig, modules *[]string) []RenderedModule {
	rendered := make([]RenderedModule, 0, len(*modules))
	for _, module := range *modules {
		var m Module
		switch module {
		case "battery":
			m = &BatteryModule{}
		case "cpu":
			m = &CpuModule{}
		}
		if m != nil {
			if m.initialize(cfg) {
				rendered = append(rendered, m.render(cfg))
			}
		}
	}
	return rendered
}

// Formatting

// Foreground color
func fg(str string, color config.Color) string {
	return fmt.Sprintf("%%F{%d}%s%%f", color, str)
}

// Bold text
func bold(str string) string {
	return fmt.Sprintf("%%B%s%%b", str)
}

// Underline text
func underline(str string) string {
	return fmt.Sprintf("%%U%s%%u", str)
}

// Standout text
func standout(str string) string {
	return fmt.Sprintf("%%S%s%%s", str)
}
