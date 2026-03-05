package modules

import (
	"fmt"
	"integral/config"
)

type SshPlus struct {
	Ssh    SshModule
	HasSsh bool
	Db     DistroboxModule
	HasDb  bool
}

func (m *SshPlus) Initialize(cfg *config.PromptConfig) bool {
	m.HasSsh = m.Ssh.Initialize(cfg)
	m.HasDb = m.Db.Initialize(cfg)
	return m.HasSsh || m.HasDb
}
func (m *SshPlus) Render(cfg *config.PromptConfig) RenderedModule {
	var state uint8 = 0
	if m.HasSsh {
		state |= 1
	}
	if m.HasDb {
		state |= 1 << 1
	}
	switch state {
	case 0b01:
		return m.Ssh.Render(cfg)
	case 0b10:
		return m.Db.Render(cfg)
	case 0b11:
		ssh := m.Ssh.Render(cfg)
		db := m.Db.Render(cfg)
		return RenderedModule{
			Length: ssh.Length + db.Length + 2,
			Fmt: fmt.Sprint(
				ssh.Fmt,
				Sh.Fg("[", cfg.Ssh.At.Color),
				db.Fmt,
				Sh.Fg("]", cfg.Ssh.At.Color),
			),
			Wrap: false,
		}
	}
	return RenderedModule{}
}

