package modules

import (
	"fmt"
	"integral/config"
	"os"
)

type SshModule struct {
	User string
	Host string
}

func (m *SshModule) Initialize(cfg *config.PromptConfig) bool {
	if _, set := os.LookupEnv("SSH_CONNECTION"); set {
		var err error
		m.User = os.Getenv("USER")
		m.Host, err = os.Hostname()
		if err != nil {
			m.Host = os.Getenv("HOSTNAME")
		}
		return true
	}
	return false
}

// [TODO] make really fancy wrapping logic to make this wrappable
func (m *SshModule) Render(cfg *config.PromptConfig) RenderedModule {
	var user, at, host string
	var ln int
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
	return RenderedModule{
		Length: ln,
		Fmt:    fmt.Sprint(user, at, host),
		Wrap:   false,
	}
}

