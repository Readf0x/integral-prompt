package main

import (
	"fmt"
	"integral/config"
	"integral/shell"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

func render(cfg *config.PromptConfig) {
	if len(os.Args) > 5 {
		logger.Fatalln("Not enough arguments!")
	}
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	prompt := finalize(cfg, width)

	fmt.Println()
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
	go generate(cfg, cfg.Modules, main)

	//assembly
	lines := assemble(size, <-main, int(cfg.Length+2), cfg.Line)

	lines = append(lines, shell.Fg(string(cfg.Line.Symbols[2]), cfg.Line.Color))
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
	lines = append(lines, shell.Fg(string(cfg.Symbols[0]), cfg.Color))
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
					lines = append(lines, shell.Fg(string(cfg.Symbols[1]), cfg.Color))
					currentLine += 1
					if i != rag-1 {
						lines[currentLine] += fmt.Sprintf("%%F{%d}", mod.Color) + mod.Raw[wrapPos+((width-1)*i):wrapPos+((width-1)*(i+1))]
					} else {
						lines[currentLine] += fmt.Sprintf("%%F{%d}", mod.Color) + mod.Raw[wrapPos+((width-1)*i):]
					}
				}
			} else {
				lines = append(lines, shell.Fg(string(cfg.Symbols[1]), cfg.Color))
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

func generate(cfg *config.PromptConfig, modules *[]string, c chan []RenderedModule) {
	rendered := make([]RenderedModule, 0, len(*modules))
	for _, module := range *modules {
		var m Module
		var M MultiModule
		switch module {
		case "visym":
			m = &ViModeModule{}
		case "dir":
			m = &DirModule{}
		case "error":
			m = &ErrorModule{}
		// case "git":
		// 	M = &GitModule{}
		case "battery":
			m = &BatteryModule{}
		case "cpu":
			m = &CpuModule{}
		case "direnv":
			m = &DirenvModule{}
		case "distrobox":
			m = &DistroboxModule{}
		case "jobs":
			m = &JobsModule{}
		case "nix":
			m = &NixModule{}
		case "ssh":
			m = &SshModule{}
		case "ssh+":
			m = &SshPlus{}
		case "time":
			m = &TimeModule{}
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
