package modules

import (
	"fmt"
	"integral/config"
	"integral/shell"
	"log"
)

var Logger *log.Logger
var Sh shell.Shell

type RenderedModule struct {
	Length int
	Fmt    string
	Wrap   bool
	Color  config.Color
}

func RenderCounter(num uint8, icon config.Char, color config.Color) RenderedModule {
	raw := fmt.Sprintf("%d%c", num, icon)
	return RenderedModule{
		Length: len(raw),
		Fmt:    Sh.Fg(raw, color),
		Wrap:   false,
		Color:  color,
	}
}
func RenderIcon(icon config.Char, color config.Color) RenderedModule {
	raw := string(icon)
	return RenderedModule{
		Length: 1,
		Fmt:    Sh.Fg(raw, color),
		Wrap:   false,
		Color:  color,
	}
}

type Module interface {
	Initialize(*config.PromptConfig) bool
	Render(*config.PromptConfig) RenderedModule
}

type MultiModule interface {
	Initialize(*config.PromptConfig) bool
	Render(*config.PromptConfig) []RenderedModule
}

