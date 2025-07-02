package main

import (
	"fmt"
	config "integral/conf"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

var logger = log.New(os.Stderr, fmt.Sprintf("\033[%dmError:\033[%dm ", 31, 39), 1)

func main() {
	if len(os.Args) < 2 {
		logger.Fatalln("Not enough arguments!")
	}
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	prompt := finalize(config.GetDefault(), width)

	for _, line := range prompt {
		fmt.Println(line)
	}
}

func finalize(cfg *config.PromptConfig, size int) []string {
	// lines := make([]string, 0, cfg.Length+2)

	// render right prompt
	// right := make(chan []RenderedModule)
	// go render(cfg, cfg.ModulesRight, right)

	// render main prompt
	main := make(chan []RenderedModule)
	go render(cfg, cfg.Modules, main)

	//assembly
	lines := assemble(size, <-main, int(cfg.Length+2), cfg.Line)

	lines = append(lines, fg(string(cfg.Line.Symbols[2]), cfg.Line.Color))
	return lines
}

func trueLength(str string) int {
	return utf8.RuneCountInString(regexp.MustCompile(`%[FfBbUuSs]({\d*})?`).ReplaceAllString(strings.ReplaceAll(str, "%%", "%"), ""))
}
func digit(num int) int {
	if num < 10 {
		return 1
	}
	return 2
}
func assemble(width int, mods []RenderedModule, maxLines int, cfg *config.LineConfig) []string {
	lines := make([]string, 0, maxLines)
	lines = append(lines, fg(string(cfg.Symbols[0]), cfg.Color))
	var currentLine = 0
	maxm := len(mods) - 1
	for i, mod := range mods {
		curLength := trueLength(lines[currentLine])
		modLength := trueLength(mod.Raw)
		if curLength+modLength > width {
			if mod.Wrap {
				wrapPos := width - curLength + 4 + digit(int(cfg.Color))
				rag := modLength / width
				lines[currentLine] += fmt.Sprintf("%%F{%d}", mod.Color) + mod.Raw[:wrapPos]
				for i := range rag {
					lines = append(lines, fg(string(cfg.Symbols[1]), cfg.Color))
					currentLine += 1
					if i != rag-1 {
						lines[currentLine] += fmt.Sprintf("%%F{%d}", mod.Color) + mod.Raw[wrapPos+((width-1)*i):wrapPos+((width-1)*(i+1))]
					} else {
						lines[currentLine] += fmt.Sprintf("%%F{%d}", mod.Color) + mod.Raw[wrapPos+((width-1)*i):]
					}
				}
			} else {
				lines = append(lines, fg(string(cfg.Symbols[1]), cfg.Color))
				lines[currentLine+1] += mod.Raw
				currentLine += 1
			}
		} else {
			lines[currentLine] += mod.Raw
		}
		if i != maxm {
			lines[currentLine] += " "
		}
	}

	return lines
}

func render(cfg *config.PromptConfig, modules *[]string, c chan []RenderedModule) {
	rendered := make([]RenderedModule, 0, len(*modules))
	for _, module := range *modules {
		var m Module
		var M MultiModule
		switch module {
		case "battery":
			m = &BatteryModule{}
		case "cpu":
			m = &CpuModule{}
		case "dir":
			m = &DirModule{}
		case "direnv":
			m = &DirenvModule{}
		case "distrobox":
			m = &DistroboxModule{}
		case "error":
			m = &ErrorModule{}
		}
		if m != nil {
			if m.initialize(cfg) {
				rendered = append(rendered, m.render(cfg))
			}
		} else if M != nil {
			if M.initialize(cfg) {
				for _, r := range M.render(cfg) {
					rendered = append(rendered, r)
				}
			}
		}
	}
	c <- rendered
}

// Formatting

// Foreground color
func fg(str string, color config.Color) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return fmt.Sprintf("%%F{%d}%s%%f", color, str)
}

// Bold text
func bold(str string) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return fmt.Sprintf("%%B%s%%b", str)
}

// Underline text
func underline(str string) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return fmt.Sprintf("%%U%s%%u", str)
}

// Standout text
func standout(str string) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return fmt.Sprintf("%%S%s%%s", str)
}
