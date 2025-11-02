package modules

import (
	"integral/config"

	"github.com/distatus/battery"
)

type BatteryModule struct {
	Charge   uint8
	Charging bool
}

func (m *BatteryModule) Initialize(cfg *config.PromptConfig) bool {
	batteries, err := battery.GetAll()
	if err != nil {
		Logger.Println(err)
		return false
	}
	if len(batteries) == 0 {
		return false
	}
	b := batteries[cfg.Battery.Id]
	if b.Full == 0 {
		return false
	}
	m.Charge = uint8(b.Current / b.Full * 100)
	m.Charging = b.State.Raw == battery.Charging
	return true
}
func (m *BatteryModule) Render(cfg *config.PromptConfig) RenderedModule {
	var icon config.Char
	var color config.Color
	if m.Charging {
		icon = cfg.Battery.IconEntries.Charging.Icon
		color = cfg.Battery.IconEntries.Charging.Color
	} else {
		icon = cfg.Battery.IconEntries.Discharging.Icon
		color = cfg.Battery.IconEntries.Discharging.Color
	}
	return RenderCounter(m.Charge, icon, color)
}

