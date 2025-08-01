package main

import (
	"fmt"
	"integral/config"
	"integral/shell"
	"log"
	"os"
	"strconv"
	"strings"
)

func render(cfg *config.PromptConfig) {
	if len(os.Args) < 6 {
		logger.Fatalln("Not enough arguments!")
	}
	width, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	prompt, rprompt := finalize(cfg, width)

	fmt.Print(sh.PromptFmt(prompt))
	if rprompt != "" {
		fmt.Println(sh.RPromptFmt(rprompt))
	}
}

func finalize(cfg *config.PromptConfig, size int) ([]string, string) {
	// lines := make([]string, 0, cfg.Length+2)

	// render right prompt
	right := make(chan []RenderedModule)
	go generate(cfg, cfg.ModulesRight, right)

	// render main prompt
	main := make(chan []RenderedModule)
	go generate(cfg, cfg.Modules, main)

	//assembly
	lines, rprompt := assemble(size, <-main, assembleRight(<-right, cfg.RightSize), cfg.WrapMinimum, int(cfg.Length+2), cfg.Line)

	lines = append(lines, sh.Fg(string(cfg.Line.Symbols[2]), cfg.Line.Color))
	return lines, rprompt
}

func digit(num int) int {
	if num < 10 {
		return 1
	}
	return 2
}

func assemble(width int, modules []RenderedModule, rightPrompt []string, wrapMinimum int, maxLines int, cfg *config.LineConfig) ([]string, string) {
	lines := make([]string, 0, maxLines)
	lines = append(lines, sh.Fg(string(cfg.Symbols[0]), cfg.Color))

	currentLine := 0
	lastIndex := len(modules) - 1
	maxWidth := width
	var rprompt string
	if len(rightPrompt) > 0 {
		for _, s := range rightPrompt {
			maxWidth = min(maxWidth, width-shell.TrueLength(s)) - 1
		}
	}

	for i, mod := range modules {
		lineLen := shell.TrueLength(lines[currentLine])

		if lineLen+mod.Length > maxWidth {
			if mod.Wrap && mod.Length > wrapMinimum {
				wrapped := wrapModule(mod, maxWidth, lineLen, cfg.Color)
				for j, segment := range wrapped {
					if j > 0 {
						lines = append(lines, sh.Fg(string(cfg.Symbols[1]), cfg.Color))
						currentLine++
					}
					lines[currentLine] += segment
					if len(rightPrompt)-1 >= currentLine {
						lines[currentLine] += " " + rightPrompt[currentLine]
					}
				}
			} else {
				if len(rightPrompt)-1 >= currentLine {
					lines[currentLine] += strings.Repeat(" ", maxWidth-shell.TrueLength(lines[currentLine])+1) + rightPrompt[currentLine]
				}
				lines = append(lines, sh.Fg(string(cfg.Symbols[1]), cfg.Color))
				currentLine++
				lines[currentLine] += mod.Fmt
			}
		} else {
			lines[currentLine] += mod.Fmt
		}

		if i != lastIndex {
			lines[currentLine] += " "
		} else {
			if len(rightPrompt)-1 >= currentLine {
				lines[currentLine] += strings.Repeat(" ", width-shell.TrueLength(lines[currentLine])-shell.TrueLength(rightPrompt[currentLine])) + rightPrompt[currentLine]
			}
		}
	}
	if len(rightPrompt) > currentLine+1 {
		if len(rightPrompt)-currentLine-2 > 0 {
			for i := 0; i < len(rightPrompt)-1; i++ {
				lines = append(lines, sh.Fg(string(cfg.Symbols[1]), cfg.Color))
				currentLine++
				lines[currentLine] += strings.Repeat(" ", width-1-shell.TrueLength(rightPrompt[currentLine])) + rightPrompt[currentLine]
			}
		}
		if sh.SupportsRP {
			rprompt = rightPrompt[currentLine+1]
		} else {
			lines = append(lines, sh.Fg(string(cfg.Symbols[1]), cfg.Color))
			currentLine++
			lines[currentLine] += strings.Repeat(" ", width-1-shell.TrueLength(rightPrompt[currentLine])) + rightPrompt[currentLine]
		}
	}

	return lines, rprompt
}
func assembleRight(modules []RenderedModule, width int) []string {
	lines := make([]string, 0, 10)
	lines = append(lines, "")

	currentLine := 0
	lastIndex := len(modules) - 1

	for i, mod := range modules {
		if mod.Length > width {
			if mod.Wrap {
				wrapped := wrapModule(mod, width, 0, 0)
				for j, segment := range wrapped {
					if j > 0 {
						lines = append(lines, "")
						currentLine++
					}
					lines[currentLine] += segment
				}
			} else {
				lines = append(lines, "")
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
func wrapModule(mod RenderedModule, width, currentLineLen int, color config.Color) []string {
	var segments []string

	// ANSI overhead for your coloring
	colorOverhead := 3 + digit(int(color))
	if color == 0 {
		colorOverhead = 0
	}
	firstWrapLimit := width - currentLineLen + colorOverhead

	segments = append(segments, mod.Fmt[:firstWrapLimit])

	remaining := mod.Fmt[firstWrapLimit:]
	for shell.TrueLength(remaining) > 0 {
		chunkSize := width - 1
		chunkSize = min(chunkSize, shell.TrueLength(remaining))
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
