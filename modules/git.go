package modules

import (
	"bufio"
	"bytes"
	"integral/config"
	"os"
	"os/exec"
	"strconv"
	"strings"
	git "github.com/go-git/go-git/v6"
)

type GitModule struct {
	Branch   string
	Unstaged uint8
	Staged   uint8
	Push     uint8
	Pull     uint8
}

func (m *GitModule) Initialize(cfg *config.PromptConfig) bool {
	cwd, err := os.Getwd()
	if err != nil {
		Logger.Println(err)
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
		cmd := exec.Command("git", "status", "--porcelain=v1")
		cmd.Dir = cwd
		out, err := cmd.Output()
		if err == nil {
			scanner := bufio.NewScanner(bytes.NewReader(out))
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) < 2 {
					continue
				}
				x, y := line[0], line[1]
				if x != ' ' && x != '?' {
					m.Staged++
				}
				if y != ' ' && y != '?' {
					m.Unstaged++
				}
			}
		}
	}

	if cfg.Git.ShowPP {
		cmd := exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
		cmd.Dir = cwd
		out, err := cmd.Output()
		if err == nil {
			// Output format: "A\tB" where:
			// A = commits ahead (need to push)
			// B = commits behind (need to pull)
			fields := strings.Fields(string(out))
			if len(fields) == 2 {
				ahead, _ := strconv.Atoi(fields[0])
				behind, _ := strconv.Atoi(fields[1])
				m.Push, m.Pull = uint8(ahead), uint8(behind)
			}
		}
	}

	return true
}
func (m *GitModule) Render(cfg *config.PromptConfig) []RenderedModule {
	list := []RenderedModule{
		{
			Length: len(m.Branch),
			Fmt:    Sh.Fg(m.Branch, cfg.Git.Branch.Color),
			Wrap:   true,
			Color:  cfg.Git.Branch.Color,
		},
	}

	if m.Unstaged != 0 {
		list = append(list, RenderCounter(m.Unstaged, cfg.Git.Unstaged.Icon, cfg.Git.Unstaged.Color))
	}
	if m.Staged != 0 {
		list = append(list, RenderCounter(m.Staged, cfg.Git.Staged.Icon, cfg.Git.Staged.Color))
	}
	if m.Push != 0 {
		list = append(list, RenderCounter(m.Push, cfg.Git.Push.Icon, cfg.Git.Push.Color))
	}
	if m.Pull != 0 {
		list = append(list, RenderCounter(m.Pull, cfg.Git.Pull.Icon, cfg.Git.Pull.Color))
	}

	return list
}

