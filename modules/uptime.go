package modules

import (
	"integral/config"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

type UptimeModule struct {
	Uptime time.Duration
}

func (m *UptimeModule) Initialize(cfg *config.PromptConfig) bool {
	raw, err := host.Uptime()
	if err != nil {
		Logger.Println(err)
		return false
	}
	m.Uptime = time.Duration(raw)
	return true
}

// m.Uptime = time.Duration(raw).String()
func (m *UptimeModule) Render(cfg *config.PromptConfig) RenderedModule {
	str := m.Uptime.String()
	return RenderedModule{
		Length: len(str),
		Fmt:    Sh.Fg(str+string(cfg.Uptime.Icon), cfg.Uptime.Color),
		Wrap:   false,
	}
}

