package modules

import (
	"fmt"
	"integral/config"
	"os"
)

type SshPlus struct {
	User   string
	Host   string
	Distro string
}

func (m *SshPlus) Initialize(cfg *config.PromptConfig) bool {
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
func (m *SshPlus) Render(cfg *config.PromptConfig) RenderedModule {
	var ssh, db, final string
	var ln int
	if m.User != "" {
		var user, at, host string
		if cfg.Ssh.User.Visible {
			ln += len(m.User)
			user = Sh.Fg(m.User, cfg.Ssh.User.Color)
		}
		if cfg.Ssh.At.Visible {
			ln += 1
			at = Sh.Fg("@", cfg.Ssh.At.Color)
		}
		if cfg.Ssh.Host.Visible {
			ln += len(m.Host)
			host = Sh.Fg(m.Host, cfg.Ssh.Host.Color)
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
		db = fmt.Sprint(Sh.Fg(m.Distro, cfg.Distrobox.TextColor), Sh.Fg(string(icon), color))
	}
	if ssh != "" && db != "" {
		final = ssh + Sh.Fg("[", cfg.Ssh.At.Color) + db + Sh.Fg("]", cfg.Ssh.At.Color)
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

