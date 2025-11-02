package modules

import (
	"integral/config"
	"os"
)

type ViModeModule struct {
	Mode string
}

func (m *ViModeModule) Initialize(cfg *config.PromptConfig) bool {
	var set bool
	m.Mode, set = os.LookupEnv("VI_KEYMAP")
	return set
}
func (m *ViModeModule) Render(cfg *config.PromptConfig) RenderedModule {
	var final RenderedModule
	switch m.Mode {
	case "NORMAL":
		final = RenderIcon(cfg.ViMode.Normal.Icon, cfg.ViMode.Normal.Color)
	case "INSERT":
		final = RenderIcon(cfg.ViMode.Insert.Icon, cfg.ViMode.Insert.Color)
	case "VISUAL":
		final = RenderIcon(cfg.ViMode.Visual.Icon, cfg.ViMode.Visual.Color)
	case "V-LINE":
		final = RenderIcon(cfg.ViMode.ViLine.Icon, cfg.ViMode.ViLine.Color)
	}
	return final
}
