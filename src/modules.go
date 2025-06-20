package main

import (
	"fmt"
	config "integral/conf"

	"github.com/distatus/battery"
	"github.com/shirou/gopsutil/v4/cpu"
)

var defmodules = []string{
	"battery",
	"cpu",
	"dir",
	"direnv",
	"distrobox",
	"error",
	"git",
	"jobs",
	"nix",
	"ssh",
	"sshplus",
	"time",
	"uptime",
	"visym",
}

type RenderedModule struct {
	Raw  string
	Wrap bool
}

func renderCounter(num uint8, icon rune, color config.Color) RenderedModule {
	i := string(icon)
	if icon == '%' {
		i = "%%"
	}
	return RenderedModule{
		Raw:  fg(fmt.Sprintf("%d%s", num, i), color),
		Wrap: false,
	}
}

type Module interface {
	initialize(*config.PromptConfig) bool
	render(*config.PromptConfig) RenderedModule
}

type BatteryModule struct {
	Charge   uint8
	Charging bool
}

func (m *BatteryModule) initialize(cfg *config.PromptConfig) bool {
	batteries, err := battery.GetAll()
	if err != nil {
		logger.Println(err)
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
func (m *BatteryModule) render(cfg *config.PromptConfig) RenderedModule {
	var icon rune
	var color config.Color
	if cfg.Battery.Icons != nil {
		color = cfg.Battery.Color
		if m.Charging {
			icon = cfg.Battery.Icons[0]
		} else {
			icon = cfg.Battery.Icons[1]
		}
	} else {
		if m.Charging {
			icon = cfg.Battery.IconEntries.Charging.Icon
			color = cfg.Battery.IconEntries.Charging.Color
		} else {
			icon = cfg.Battery.IconEntries.Discharging.Icon
			color = cfg.Battery.IconEntries.Discharging.Color
		}
	}
	return renderCounter(m.Charge, icon, color)
}

type CpuModule struct {
	Usage uint8
}
func (m *CpuModule) initialize(cfg *config.PromptConfig) bool {
	percent, err := cpu.Percent(cfg.Cpu.Time, false)
	if err != nil {
		logger.Println(err)
		return false
	}
	m.Usage = uint8(percent[0])
	return true
}
func (m *CpuModule) render(cfg *config.PromptConfig) RenderedModule {
	// [TODO] add multi icon support
	return renderCounter(m.Usage, cfg.Cpu.Icon, cfg.Cpu.Color)
}

type DirModule struct {
	CWD string
}

type DirenvModule struct {
	InDirenvShell bool
}

type DistroboxModule struct {
	Distro string
}

type ErrorModule struct {
	ExitCode uint8
}

type GitModule struct {
	Branch   string
	Unstaged uint16
	Staged   uint16
	Push     uint16
	Pull     uint16
}

type JobsModule struct {
	Jobs uint8
}

type NixModule struct {
	InNixShell bool
}

type SshModule struct {
	User string
	Host string
}

type SshPlus struct {
	User   string
	Host   string
	Distro string
}

type TimeModule struct {
	Time string
}

type UptimeModule struct {
	Uptime string
}

type Mode uint8

const (
	ViNormal Mode = iota
	ViInsert
	ViVisual
	ViVisualLine
)

type ViModeModule struct {
	Mode Mode
}
