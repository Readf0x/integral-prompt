package shell

import (
	"fmt"
	"integral/config"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Shell struct {
	Fg           func(string, config.Color) string
	Bold         func(string) string
	Underline    func(string) string
	PromptFmt    func(prompt []string) string
	RPromptFmt   func(prompt string) string
	Init         func()
	SupportsRP   bool
}

var Raw = Shell{
	Fg:         gFg,
	Bold:       gBold,
	Underline:  gUnderline,
	PromptFmt:  gPromptFmt,
	RPromptFmt: gRPromptFmt,
	SupportsRP: true,
}

func gFg(str string, color config.Color) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", color, str)
}
func gBold(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}
func gUnderline(str string) string {
	return fmt.Sprintf("\033[4m%s\033[0m", str)
}
func gPromptFmt(prompt []string) string {
	return strings.Join(prompt, "\n")
}
func gRPromptFmt(prompt string) string {
	i, _ := strconv.Atoi(os.Args[3])
	return strings.Repeat(" ", i - TrueLength(prompt) - 1) + prompt
}

func GetShell(str string) (sh Shell, err error) {
	switch str {
	case "zsh":
		return Zsh, nil
	case "raw":
		return Raw, nil
	}
	return Shell{}, fmt.Errorf("Invalid shell name")
}

func printFile(path string) error {
  file, err := os.Open(path)
  if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
  }
  defer file.Close()

  _, err = io.Copy(os.Stdout, file)
  if err != nil {
		return fmt.Errorf("failed to print file contents: %w", err)
  }
  return nil
}
func findShare() (string, error) {
	xdgDataDirs := os.Getenv("XDG_DATA_DIRS")
	if xdgDataDirs == "" {
		xdgDataDirs = "/usr/local/share:/usr/share"
	}
	shareDirs := strings.Split(xdgDataDirs, ":")

	for _, base := range shareDirs {
		fullPath := filepath.Join(base, "integral")
		info, err := os.Stat(fullPath)
		if err == nil && info.IsDir() {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("share path not found in $XDG_DATA_DIRS")
}

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)
func TrueLength(str string) int {
	clean := ansiRegexp.ReplaceAllString(str, "")
	return utf8.RuneCountInString(clean)
}

