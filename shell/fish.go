package shell

import (
	"fmt"
	"log"
	"strings"
)

var Fish = Shell{
	Fg:         gFg,
	Bold:       gBold,
	Underline:  gUnderline,
	PromptFmt:  fishPromptFmt,
	RPromptFmt: fishRPromptFmt,
	Init:       fishInit,
	SupportsRP: true,
}

func fishPromptFmt(prompt []string) string {
	return fmt.Sprintf("echo \"\n%s\"", strings.Join(prompt, "\n"))
}
func fishRPromptFmt(prompt string) string {
	return fmt.Sprintf("set -g integral_right_prompt_string \"%s\"", prompt)
}
func fishInit() {
	share, err := findShare()
	if err != nil {
		log.Fatal(err)
	}
	printFile(share + "/init.fish")
}

