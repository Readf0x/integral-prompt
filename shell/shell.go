package shell

import (
	"fmt"
	"integral/config"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Fg(str string, color config.Color) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", color, str)
}
func Bold(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}
func Underline(str string) string {
	return fmt.Sprintf("\033[4m%s\033[0m", str)
}

func PromptFmt(prompt []string) string {
	return "PROMPT=\"\n%{" + strings.Join(prompt, "%}\n%{%G") + "%}\""
}
func RPromptFmt(prompt string) string {
	return "RPROMPT=\"" + prompt + "\""
}
func Init() {
	share, err := findShare()
	if err != nil {
		log.Fatal(err)
	}
	printFile(share + "/init.zsh")
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

