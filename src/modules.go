package main

import (
	"fmt"
	config "integral/conf"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/distatus/battery"
	// "github.com/go-git/go-git/v5"
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
	Raw   string
	Wrap  bool
	Color config.Color
}

func renderCounter(num uint8, icon rune, color config.Color) RenderedModule {
	return RenderedModule{
		Raw:   fg(fmt.Sprintf("%d%c", num, icon), color),
		Wrap:  false,
		Color: color,
	}
}
func renderIcon(icon rune, color config.Color) RenderedModule {
	return RenderedModule{
		Raw:   fg(string(icon), color),
		Wrap:  false,
		Color: color,
	}
}

type Module interface {
	initialize(*config.PromptConfig) bool
	render(*config.PromptConfig) RenderedModule
}

type MultiModule interface {
	initialize(*config.PromptConfig) bool
	render(*config.PromptConfig) []RenderedModule
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
	if m.Charging {
		icon = cfg.Battery.IconEntries.Charging.Icon
		color = cfg.Battery.IconEntries.Charging.Color
	} else {
		icon = cfg.Battery.IconEntries.Discharging.Icon
		color = cfg.Battery.IconEntries.Discharging.Color
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
func (m *DirModule) initialize(cfg *config.PromptConfig) bool {
	var err error
	m.CWD, err = os.Getwd()
	if err != nil {
		logger.Println(err)
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
func (m *DirModule) render(cfg *config.PromptConfig) RenderedModule {
	return RenderedModule{
		Raw:   fg(m.CWD, cfg.Dir.Color),
		Wrap:  true,
		Color: cfg.Dir.Color,
	}
}

type DirenvModule struct {
}
func (m *DirenvModule) initialize(cfg *config.PromptConfig) bool {
	if _, set := os.LookupEnv("DIRENV_DIR"); set {
		return true
	}
	return false
}
func (m *DirenvModule) render(cfg *config.PromptConfig) RenderedModule {
	var color config.Color = cfg.Direnv.DefaultIcon.Color
	var icon rune = cfg.Direnv.DefaultIcon.Icon
	if cfg.Direnv.IconEntries != nil {
		for _, entry := range *cfg.Direnv.IconEntries {
			_, err := exec.LookPath(entry.Name)
			if err == nil {
				color = entry.Color
				icon = entry.Icon
			}
		}
	}
	return renderIcon(icon, color)
}

type DistroboxModule struct {
	Distro string
}

func (m *DistroboxModule) initialize(cfg *config.PromptConfig) bool {
	var set bool
	m.Distro, set = os.LookupEnv("CONTAINER_ID")
	return set
}
func (m *DistroboxModule) render(cfg *config.PromptConfig) RenderedModule {
	color, icon := cfg.Distrobox.DefaultIcon.Color, cfg.Distrobox.DefaultIcon.Icon
	for _, distro := range *cfg.Distrobox.IconEntries {
		if m.Distro == distro.Name {
			color, icon = distro.Color, distro.Icon
		}
	}
	return RenderedModule{
		Raw:   fmt.Sprint(fg(m.Distro, cfg.Distrobox.TextColor), fg(string(icon), color)),
		Wrap:  true,
		Color: color,
	}
}

type ErrorModule struct {
	ExitCode uint64
}
func (m *ErrorModule) initialize(cfg *config.PromptConfig) bool {
	c, err := strconv.ParseUint(os.Args[1], 10, 8)
	if err != nil {
		logger.Panicln(err)
		return false
	}
	if c != 0 {
		m.ExitCode = c
		return true
	}
	return false
}
func (m *ErrorModule) render(cfg *config.PromptConfig) RenderedModule {
	color, icon := cfg.Error.DefaultIcon.Color, cfg.Error.DefaultIcon.Icon
	for _, c := range *cfg.Error.IconEntries {
		if m.ExitCode == c.Code {
			color, icon = c.Color, c.Icon
		}
	}
	return renderIcon(icon, color)
}

type GitModule struct {
	Branch   string
	Unstaged uint16
	Staged   uint16
	Push     uint16
	Pull     uint16
}
// func (m *GitModule) initialize(cfg *config.PromptConfig) bool {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		logger.Println(err)
// 		return false
// 	}
// 	repo, err := git.PlainOpen(cwd)
// 	if err != nil {
// 		logger.Println(err)
// 		return false
// 	}
// 	// repo.
// 	return true
// }

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
