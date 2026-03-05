package modules

import (
	"integral/config"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CpuModule struct {
	Usage uint8
}

func (m *CpuModule) Initialize(cfg *config.PromptConfig) bool {
	percent, err := cpu.Percent(cfg.Cpu.Time, false)
	if err != nil {
		Logger.Println(err)
		return false
	}
	m.Usage = uint8(percent[0])
	return true
}
func (m *CpuModule) Render(cfg *config.PromptConfig) RenderedModule {
	// [TODO] add multi icon support
	return RenderCounter(m.Usage, cfg.Cpu.Icon, cfg.Cpu.Color)
}

