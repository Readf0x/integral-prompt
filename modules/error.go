package modules

import (
	"integral/config"
	"os"
	"strconv"
)

type ErrorModule struct {
	ExitCode uint64
}

func (m *ErrorModule) Initialize(cfg *config.PromptConfig) bool {
	c, err := strconv.ParseUint(os.Args[4], 10, 8)
	if err != nil {
		Logger.Panicln(err)
		return false
	}
	if c != 0 {
		m.ExitCode = c
		return true
	}
	return false
}
func (m *ErrorModule) Render(cfg *config.PromptConfig) RenderedModule {
	color, icon := cfg.Error.DefaultIcon.Color, cfg.Error.DefaultIcon.Icon
	for _, c := range *cfg.Error.IconEntries {
		if m.ExitCode == c.Code {
			color, icon = c.Color, c.Icon
		}
	}
	return RenderIcon(icon, color)
}

