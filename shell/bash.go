package shell

import (
	"log"
	"strings"
)

var Bash = Shell{
	Fg:         gFg,
	Bold:       gBold,
	Underline:  gUnderline,
	PromptFmt:  bashPromptFmt,
	RPromptFmt: gRPromptFmt,
	Init:       bashInit,
	SupportsRP: false,
}

func bashPromptFmt(prompt []string) string {
	return "PS1='\n" + strings.Join(prompt, "\n") + "'"
}

func bashInit() {
	share, err := findShare()
	if err != nil {
		log.Fatal(err)
	}
	printFile(share + "/init.bash")
}
