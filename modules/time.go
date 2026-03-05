package modules

import (
	"integral/config"
	"time"
)

type TimeModule struct {
	Time string
}

func (m *TimeModule) Initialize(cfg *config.PromptConfig) bool {
	m.Time = time.Now().Format(cfg.Time.Format)
	return true
}
func (m *TimeModule) Render(cfg *config.PromptConfig) RenderedModule {
	return RenderedModule{
		Length: len(m.Time),
		Fmt:    Sh.Fg(m.Time, cfg.Time.Color),
		Wrap:   true,
		Color:  cfg.Time.Color,
	}
}

