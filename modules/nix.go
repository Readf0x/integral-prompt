package modules

import (
	"integral/config"
	"os"
	"strings"
)

type NixModule struct {
	InNixShell bool
}

func (m *NixModule) Initialize(cfg *config.PromptConfig) bool {
	path := strings.Split(os.Getenv("PATH"), ":")
	if strings.HasPrefix(path[0], "/nix/store/") {
		m.InNixShell = true
		return true
	}
	return false
}
func (m *NixModule) Render(cfg *config.PromptConfig) RenderedModule {
	return RenderIcon(cfg.NixShell.Icon, cfg.NixShell.Color)
}

