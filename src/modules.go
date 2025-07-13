package main

import (
	"fmt"
	config "integral/conf"
	"integral/shell"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/distatus/battery"
	"github.com/go-git/go-git/v6"
	"github.com/shirou/gopsutil/v4/cpu"
)

type RenderedModule struct {
	Raw   string
	Wrap  bool
	Color config.Color
}

func renderCounter(num uint8, icon rune, color config.Color) RenderedModule {
	return RenderedModule{
		Raw:   shell.Fg(fmt.Sprintf("%d%c", num, icon), color),
		Wrap:  false,
		Color: color,
	}
}
func renderIcon(icon rune, color config.Color) RenderedModule {
	return RenderedModule{
		Raw:   shell.Fg(string(icon), color),
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
		Raw:   shell.Fg(m.CWD, cfg.Dir.Color),
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
		Raw:   fmt.Sprint(shell.Fg(m.Distro, cfg.Distrobox.TextColor), shell.Fg(string(icon), color)),
		Wrap:  true,
		Color: color,
	}
}

type ErrorModule struct {
	ExitCode uint64
}

func (m *ErrorModule) initialize(cfg *config.PromptConfig) bool {
	c, err := strconv.ParseUint(os.Args[3], 10, 8)
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
	Unstaged uint8
	Staged   uint8
	Push     uint8
	Pull     uint8
}

func (m *GitModule) initialize(cfg *config.PromptConfig) bool {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Println(err)
		return false
	}
	repo, err := git.PlainOpen(cwd)
	if err != nil {
		logger.Println(err)
		return false
	}

	head, err := repo.Head()
	if err != nil {
		logger.Println(err)
		return false
	}
	m.Branch = head.Name().Short()

	wt, err := repo.Worktree()
	if err != nil {
		logger.Println(err)
		return false
	}
	status, err := wt.Status()
	if err != nil {
		logger.Println(err)
		return false
	}

	for _, entry := range status {
		switch {
		case entry.Worktree != git.Unmodified:
			m.Unstaged++
		case entry.Staging != git.Unmodified:
			m.Staged++
		}
	}

	// [TODO] implement push and pull detection

	return true
}
func (m *GitModule) render(cfg *config.PromptConfig) []RenderedModule {
	// [TODO] implement git render
	return []RenderedModule{}
}

type JobsModule struct {
	Jobs uint8
}
func (m *JobsModule) initialize(cfg *config.PromptConfig) bool {
	j, err := strconv.ParseUint(os.Args[4], 10, 8)
	if err != nil {
		logger.Println(err)
		return false
	}
	if j != 0 {
		m.Jobs = uint8(j)
		return true
	}
	return false
}
func (m *JobsModule) render(cfg *config.PromptConfig) RenderedModule {
	return renderCounter(m.Jobs, cfg.Jobs.Icon, cfg.Jobs.Color)
}

type NixModule struct {
	InNixShell bool
}
func (m *NixModule) initialize(cfg *config.PromptConfig) bool {
	path := strings.Split(os.Getenv("PATH"), ":")
	if strings.HasPrefix(path[0], "/nix/store/") {
		m.InNixShell = true
		return true
	}
	return false
}
func (m *NixModule) render(cfg *config.PromptConfig) RenderedModule {
	return renderIcon(cfg.NixShell.Icon, cfg.NixShell.Color)
}

type SshModule struct {
	User string
	Host string
}
func (m *SshModule) initialize(cfg *config.PromptConfig) bool {
	if _, set := os.LookupEnv("SSH_CONNECTION"); set {
		m.User = os.Getenv("USER")
		m.Host = os.Getenv("HOSTNAME")
		return true
	}
	return false
}
// [TODO] make really fancy wrapping logic to make this wrappable
func (m *SshModule) render(cfg *config.PromptConfig) RenderedModule {
	var user, at, host string
	if cfg.Ssh.User.Visible {
		user = shell.Fg(m.User, cfg.Ssh.User.Color)
	}
	if cfg.Ssh.At.Visible {
		at = shell.Fg("@", cfg.Ssh.At.Color)
	}
	if cfg.Ssh.Host.Visible {
		host = shell.Fg(m.Host, cfg.Ssh.Host.Color)
	}
	return RenderedModule{
		Raw: fmt.Sprint(user, at, host),
		Wrap: false,
	}
}

type SshPlus struct {
	User   string
	Host   string
	Distro string
}
func (m *SshPlus) initialize(cfg *config.PromptConfig) bool {
	if os.Getenv("SSH_CONNECTION") != "" {
		m.User = os.Getenv("USER")
		m.Host = os.Getenv("HOSTNAME")
	}
	m.Distro, _ = os.LookupEnv("CONTAINER_ID")
	if m.User != "" || m.Distro != "" {
		return true
	}
	return false
}
func (m *SshPlus) render(cfg *config.PromptConfig) RenderedModule {
	var ssh, db, final string
	if m.User != "" {
		var user, at, host string
		if cfg.Ssh.User.Visible {
			user = shell.Fg(m.User, cfg.Ssh.User.Color)
		}
		if cfg.Ssh.At.Visible {
			at = shell.Fg("@", cfg.Ssh.At.Color)
		}
		if cfg.Ssh.Host.Visible {
			host = shell.Fg(m.Host, cfg.Ssh.Host.Color)
		}
		ssh = fmt.Sprint(user, at, host)
	}
	if m.Distro != "" {
		color, icon := cfg.Distrobox.DefaultIcon.Color, cfg.Distrobox.DefaultIcon.Icon
		for _, distro := range *cfg.Distrobox.IconEntries {
			if m.Distro == distro.Name {
				color, icon = distro.Color, distro.Icon
			}
		}
		db = fmt.Sprint(shell.Fg(m.Distro, cfg.Distrobox.TextColor), shell.Fg(string(icon), color))
	}
	if ssh != "" && db != "" {
		final = ssh + shell.Fg("[", cfg.Ssh.At.Color) + db + shell.Fg("]", cfg.Ssh.At.Color)
	} else if ssh != "" {
		final = ssh
	} else if db != "" {
		final = db
	}
	return RenderedModule{
		Raw: final,
		Wrap: false,
	}
}

type TimeModule struct {
	Time string
}
func (m *TimeModule) initialize(cfg *config.PromptConfig) bool {
	m.Time = time.Now().Format(cfg.Time.Format)
	return true
}
func (m *TimeModule) render(cfg *config.PromptConfig) RenderedModule {
	return RenderedModule{
		Raw: shell.Fg(m.Time, cfg.Time.Color),
		Wrap: true,
		Color: cfg.Time.Color,
	}
}

// [TODO] uptime
// frankly I have no clue how to read /proc/uptime
type UptimeModule struct {
	Uptime string
}

type ViModeModule struct {
	Mode string
}
func (m *ViModeModule) initialize(cfg *config.PromptConfig) bool {
	var set bool
	m.Mode, set = os.LookupEnv("VI_KEYMAP")
	return set
}
func (m *ViModeModule) render(cfg *config.PromptConfig) RenderedModule {
	var final RenderedModule
	switch m.Mode {
	case "NORMAL":
		final = renderIcon(cfg.ViMode.Normal.Icon, cfg.ViMode.Normal.Color)
	case "INSERT":
		final = renderIcon(cfg.ViMode.Insert.Icon, cfg.ViMode.Insert.Color)
	case "VISUAL":
		final = renderIcon(cfg.ViMode.Visual.Icon, cfg.ViMode.Visual.Color)
	case "V-LINE":
		final = renderIcon(cfg.ViMode.ViLine.Icon, cfg.ViMode.ViLine.Color)
	}
	return final
}
