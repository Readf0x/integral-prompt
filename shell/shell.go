package shell

import (
	"fmt"
	"integral/config"
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

