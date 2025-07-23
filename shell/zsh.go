package shell

import (
	"log"
	"strings"
)

var Zsh = Shell{
	Fg:         gFg,
	Bold:       gBold,
	Underline:  gUnderline,
	PromptFmt:  zshPromptFmt,
	RPromptFmt: zshRPromptFmt,
	Init:       zshInit,
}

func zshPromptFmt(prompt []string) string {
	return "PROMPT=\"\n%{" + strings.Join(prompt, "%}\n%{%G") + "%}\""
}
func zshRPromptFmt(prompt string) string {
	return "RPROMPT=\"" + prompt + "\""
}
func zshInit() {
	share, err := findShare()
	if err != nil {
		log.Fatal(err)
	}
	printFile(share + "/init.zsh")
}

