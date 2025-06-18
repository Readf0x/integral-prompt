package main

import "fmt"

func fg(str string, color Color) string {
	return fmt.Sprintf("%%F{%d}%s%%f", color, str)
}

func main() {
	test_prompt :=
		fmt.Sprint(
		fg("⌠", Yellow), fg("○ ", Green), fg("~ ", Blue), "\n",
		fg("⌡", Yellow))

	fmt.Print(test_prompt)
}

