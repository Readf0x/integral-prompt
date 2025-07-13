package shell

import (
	"fmt"
	"integral/config"
	"strings"
)

type Shell uint8
const (
	ZSH Shell = iota
	NUSH
)

// Formatting

// Foreground color
func Fg(str string, color config.Color) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return fmt.Sprintf("%%F{%d}%s%%f", color, str)
}

// Bold text
func Bold(str string) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return fmt.Sprintf("%%B%s%%b", str)
}

// Underline text
func Underline(str string) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return fmt.Sprintf("%%U%s%%u", str)
}

// Standout text
func Standout(str string) string {
	str = strings.ReplaceAll(str, "%", "%%")
	return fmt.Sprintf("%%S%s%%s", str)
}
