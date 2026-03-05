package modules

import (
	"integral/config"
	"os"
	"strings"
)

type DirModule struct {
	CWD string
}

func (m *DirModule) Initialize(cfg *config.PromptConfig) bool {
	var err error
	m.CWD, err = os.Getwd()
	if err != nil {
		Logger.Println(err)
		return false
	}
	for _, replaceList := range *cfg.Dir.Replace {
		m.CWD = strings.ReplaceAll(m.CWD, (*replaceList)[0], (*replaceList)[1])
	}
	if cfg.Dir.ReplaceHome {
		m.CWD = strings.Replace(m.CWD, os.Getenv("HOME"), string(cfg.Dir.HomeIcon), 1)
	}
	return true
}
func (m *DirModule) Render(cfg *config.PromptConfig) RenderedModule {
	return RenderedModule{
		Length: len(m.CWD),
		Fmt:    Sh.Fg(m.CWD, cfg.Dir.Color),
		Wrap:   true,
		Color:  cfg.Dir.Color,
	}
}

