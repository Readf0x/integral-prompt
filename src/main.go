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

	// prompt := assemble(config.GetDefault(), width, true)

	fmt.Println(fg("Width:", config.Yellow), bold(strconv.Itoa(width)))
	fmt.Println(fg("=========", config.BrightBlack))
	// for _, line := range prompt {
	// 	fmt.Println(line)
	// }
	b := BatteryModule{}
	b.initialize(config.GetDefault())
	fmt.Println(b.render(config.GetDefault()).Raw)
	fmt.Println(fg("=========", config.BrightBlack))
}

func assemble(cfg *config.PromptConfig, size int, debug bool) []string {
	lines := make([]string, 0, cfg.Length+2)

	// rc := make(chan string)
	// render right prompt

	return lines
}

func render() {

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
