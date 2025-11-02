package modules

import (
	"integral/config"
	"os"
	"os/exec"
)

type DirenvModule struct {
}

func (m *DirenvModule) Initialize(cfg *config.PromptConfig) bool {
	if _, set := os.LookupEnv("DIRENV_DIR"); set {
		return true
	}
	return false
}
func (m *DirenvModule) Render(cfg *config.PromptConfig) RenderedModule {
	var color config.Color = cfg.Direnv.DefaultIcon.Color
	var icon config.Char = cfg.Direnv.DefaultIcon.Icon
	if cfg.Direnv.IconEntries != nil {
		for _, entry := range *cfg.Direnv.IconEntries {
			_, err := exec.LookPath(entry.Name)
			if err == nil {
				color = entry.Color
				icon = entry.Icon
			}
		}
	}
	return RenderIcon(icon, color)
}

