package modules

import (
	"integral/config"
	"os"
	"strconv"
)

type JobsModule struct {
	Jobs uint8
}

func (m *JobsModule) Initialize(cfg *config.PromptConfig) bool {
	j, err := strconv.ParseUint(os.Args[5], 10, 8)
	if err != nil {
		Logger.Println(err)
		return false
	}
	if j != 0 {
		m.Jobs = uint8(j)
		return true
	}
	return false
}
func (m *JobsModule) Render(cfg *config.PromptConfig) RenderedModule {
	return RenderCounter(m.Jobs, cfg.Jobs.Icon, cfg.Jobs.Color)
}

