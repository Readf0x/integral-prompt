package modules

import (
	"fmt"
	"integral/config"
	"os"
)

type DistroboxModule struct {
	Distro string
}

func (m *DistroboxModule) Initialize(cfg *config.PromptConfig) bool {
	var set bool
	m.Distro, set = os.LookupEnv("CONTAINER_ID")
	return set
}
func (m *DistroboxModule) Render(cfg *config.PromptConfig) RenderedModule {
	color, icon := cfg.Distrobox.DefaultIcon.Color, cfg.Distrobox.DefaultIcon.Icon
	for _, distro := range *cfg.Distrobox.IconEntries {
		if m.Distro == distro.Name {
			color, icon = distro.Color, distro.Icon
		}
	}
	return RenderedModule{
		Length: len(m.Distro) + 1,
		Fmt:    fmt.Sprint(Sh.Fg(m.Distro, cfg.Distrobox.TextColor), Sh.Fg(string(icon), color)),
		Wrap:   true,
		Color:  color,
	}
}

