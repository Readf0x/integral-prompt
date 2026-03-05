package main

import (
	"fmt"
	"integral/config"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/distatus/battery"
	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
)

type RenderedModule struct {
	Length int
	Fmt    string
	Wrap   bool
	Color  config.Color
}

func renderCounter(num uint8, icon config.Char, color config.Color) RenderedModule {
	raw := fmt.Sprintf("%d%c", num, icon)
	return RenderedModule{
		Length: len(raw),
		Fmt:    sh.Fg(raw, color),
		Wrap:   false,
		Color:  color,
	}
}
func renderIcon(icon config.Char, color config.Color) RenderedModule {
	raw := string(icon)
	return RenderedModule{
		Length: 1,
		Fmt:    sh.Fg(raw, color),
		Wrap:   false,
		Color:  color,
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
	var icon config.Char
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
		Length: len(m.CWD),
		Fmt:    sh.Fg(m.CWD, cfg.Dir.Color),
		Wrap:   true,
		Color:  cfg.Dir.Color,
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
		Length: len(m.Distro) + 1,
		Fmt:    fmt.Sprint(sh.Fg(m.Distro, cfg.Distrobox.TextColor), sh.Fg(string(icon), color)),
		Wrap:   true,
		Color:  color,
	}
}

type ErrorModule struct {
	ExitCode uint64
}

func (m *ErrorModule) initialize(cfg *config.PromptConfig) bool {
	c, err := strconv.ParseUint(os.Args[4], 10, 8)
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

	var repo *git.Repository
	for i := 0; i < cfg.Git.RecurseCount; i++ {
		repo, err = git.PlainOpen(cwd)
		if err == nil {
			break
		}
		cwd += "/.."
	}
	if err != nil {
		return false
	}

	head, err := repo.Head()
	if err != nil {
		return false
	}
	m.Branch = head.Name().Short()

	if cfg.Git.ShowWT {
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
	}

	if cfg.Git.ShowPP {
		config, err := repo.Config()
		if err != nil {
			logger.Println(err)
			return false
		}
		branchCfg, ok := config.Branches[m.Branch]
		if ok {
			remoteRef, err := repo.Reference(plumbing.NewRemoteReferenceName(branchCfg.Remote, m.Branch), true)
			if err == nil {
				local, err := getCommits(repo, head.Hash())
				if err != nil {
					logger.Println(err)
					return false
				}
				remote, err := getCommits(repo, remoteRef.Hash())
				if err != nil {
					logger.Println(err)
					return false
				}

				m.Push, m.Pull = diffCommits(local, remote)
			}
		}
	}

	return true
}
func getCommits(repo *git.Repository, hash plumbing.Hash) ([]plumbing.Hash, error) {
	var hashes []plumbing.Hash

	commit, err := repo.CommitObject(hash)
	if err != nil {
		return nil, err
	}

	err = object.NewCommitPreorderIter(commit, nil, nil).ForEach(func(c *object.Commit) error {
		hashes = append(hashes, c.Hash)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return hashes, nil
}
func diffCommits(a, b []plumbing.Hash) (onlyInA, onlyInB uint8) {
	set := make(map[plumbing.Hash]struct{}, len(a))
	for _, h := range a {
		set[h] = struct{}{}
	}
	for _, h := range b {
		if _, ok := set[h]; ok {
			delete(set, h)
		} else {
			onlyInB++
		}
	}
	onlyInA = uint8(len(set))
	return
}

func (m *GitModule) render(cfg *config.PromptConfig) []RenderedModule {
	list := []RenderedModule{
		{
			Length: len(m.Branch),
			Fmt:    sh.Fg(m.Branch, cfg.Git.Branch.Color),
			Wrap:   true,
			Color:  cfg.Git.Branch.Color,
		},
	}

	if m.Unstaged != 0 {
		list = append(list, renderCounter(m.Unstaged, cfg.Git.Unstaged.Icon, cfg.Git.Unstaged.Color))
	}
	if m.Staged != 0 {
		list = append(list, renderCounter(m.Staged, cfg.Git.Staged.Icon, cfg.Git.Staged.Color))
	}
	if m.Push != 0 {
		list = append(list, renderCounter(m.Push, cfg.Git.Push.Icon, cfg.Git.Push.Color))
	}
	if m.Pull != 0 {
		list = append(list, renderCounter(m.Pull, cfg.Git.Pull.Icon, cfg.Git.Pull.Color))
	}

	return list
}

type JobsModule struct {
	Jobs uint8
}

func (m *JobsModule) initialize(cfg *config.PromptConfig) bool {
	j, err := strconv.ParseUint(os.Args[5], 10, 8)
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
	var ln int
	if cfg.Ssh.User.Visible {
		ln += len(m.User)
		user = sh.Fg(m.User, cfg.Ssh.User.Color)
	}
	if cfg.Ssh.At.Visible {
		ln += 1
		at = sh.Fg("@", cfg.Ssh.At.Color)
	}
	if cfg.Ssh.Host.Visible {
		ln += len(m.Host)
		host = sh.Fg(m.Host, cfg.Ssh.Host.Color)
	}
	return RenderedModule{
		Length: ln,
		Fmt:    fmt.Sprint(user, at, host),
		Wrap:   false,
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
	var ln int
	if m.User != "" {
		var user, at, host string
		if cfg.Ssh.User.Visible {
			ln += len(m.User)
			user = sh.Fg(m.User, cfg.Ssh.User.Color)
		}
		if cfg.Ssh.At.Visible {
			ln += 1
			at = sh.Fg("@", cfg.Ssh.At.Color)
		}
		if cfg.Ssh.Host.Visible {
			ln += len(m.Host)
			host = sh.Fg(m.Host, cfg.Ssh.Host.Color)
		}
		ssh = fmt.Sprint(user, at, host)
	}
	if m.Distro != "" {
		ln += len(m.Distro)
		color, icon := cfg.Distrobox.DefaultIcon.Color, cfg.Distrobox.DefaultIcon.Icon
		for _, distro := range *cfg.Distrobox.IconEntries {
			if m.Distro == distro.Name {
				color, icon = distro.Color, distro.Icon
			}
		}
		db = fmt.Sprint(sh.Fg(m.Distro, cfg.Distrobox.TextColor), sh.Fg(string(icon), color))
	}
	if ssh != "" && db != "" {
		final = ssh + sh.Fg("[", cfg.Ssh.At.Color) + db + sh.Fg("]", cfg.Ssh.At.Color)
		ln += 2
	} else if ssh != "" {
		final = ssh
	} else if db != "" {
		final = db
	}
	return RenderedModule{
		Length: ln,
		Fmt:    final,
		Wrap:   false,
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
		Length: len(m.Time),
		Fmt:    sh.Fg(m.Time, cfg.Time.Color),
		Wrap:   true,
		Color:  cfg.Time.Color,
	}
}

type UptimeModule struct {
	Uptime time.Duration
}

func (m *UptimeModule) initialize(cfg *config.PromptConfig) bool {
	raw, err := host.Uptime()
	if err != nil {
		logger.Println(err)
		return false
	}
	m.Uptime = time.Duration(raw)
	return true
}

// m.Uptime = time.Duration(raw).String()
func (m *UptimeModule) render(cfg *config.PromptConfig) RenderedModule {
	str := m.Uptime.String()
	return RenderedModule{
		Length: len(str),
		Fmt:    sh.Fg(str+string(cfg.Uptime.Icon), cfg.Uptime.Color),
		Wrap:   false,
	}
}

type ViModeModule struct {
	Mode string
}

func (m *ViModeModule) initialize(cfg *config.PromptConfig) bool {
	var set bool
	m.Mode, set = os.LookupEnv("VI_KEYMAP")
	return set && os.Args[2] == "zsh"
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
