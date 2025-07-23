package main

import (
	"fmt"
	"integral/config"
	"log"
	"os"
	"regexp"
	"strconv"
	"unicode/utf8"
)

func render(cfg *config.PromptConfig) {
	if len(os.Args) < 6 {
		logger.Fatalln("Not enough arguments!")
	}
	width, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	prompt := finalize(cfg, width)

	fmt.Println(sh.PromptFmt(prompt))
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

	lines = append(lines, sh.Fg(string(cfg.Line.Symbols[2]), cfg.Line.Color))
	return lines
}

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func trueLength(str string) int {
	clean := ansiRegexp.ReplaceAllString(str, "")
	return utf8.RuneCountInString(clean)
}
func digit(num int) int {
	if num < 10 {
		return 1
	}
	return 2
}

func assemble(width int, modules []RenderedModule, maxLines int, cfg *config.LineConfig) []string {
	lines := make([]string, 0, maxLines)
	lines = append(lines, sh.Fg(string(cfg.Symbols[0]), cfg.Color))

	currentLine := 0
	lastIndex := len(modules) - 1

	for i, mod := range modules {
		lineLen := trueLength(lines[currentLine])

		if lineLen+mod.Length > width {
			if mod.Wrap {
				wrapped := wrapModule(mod, width, lineLen, cfg)
				for j, segment := range wrapped {
					if j > 0 {
						lines = append(lines, sh.Fg(string(cfg.Symbols[1]), cfg.Color))
						currentLine++
					}
					lines[currentLine] += segment
				}
			} else {
				lines = append(lines, sh.Fg(string(cfg.Symbols[1]), cfg.Color))
				currentLine++
				lines[currentLine] += mod.Fmt
			}
		} else {
			lines[currentLine] += mod.Fmt
		}

		if i != lastIndex {
			lines[currentLine] += " "
		}
	}

	return lines
}
func wrapModule(mod RenderedModule, width, currentLineLen int, cfg *config.LineConfig) []string {
	var segments []string

	// ANSI overhead for your coloring
	colorOverhead := 3 + digit(int(cfg.Color))
	firstWrapLimit := width - currentLineLen + colorOverhead

	segments = append(segments, mod.Fmt[:firstWrapLimit])

	remaining := mod.Fmt[firstWrapLimit:]
	for trueLength(remaining) > 0 {
		chunkSize := width - 1
		chunkSize = min(chunkSize, trueLength(remaining))
		segments = append(segments, fmt.Sprintf("\033[%dm%s", mod.Color, remaining[:chunkSize]))
		remaining = remaining[chunkSize:]
	}

	return segments
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
		case "git":
			M = &GitModule{}
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
