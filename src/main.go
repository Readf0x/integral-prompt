package main

import (
  "fmt"
  config "integral/conf"
  "log"
  "os"
  "strconv"

  "golang.org/x/term"
)

func main() {
  width, _, err := term.GetSize(int(os.Stdin.Fd()))
  if err != nil {
    log.Fatal(err)
  }

  test_prompt := fmt.Sprint(
    fg("⌠", config.Yellow), fg("○ ", config.Green), fg("~ ", config.Blue), "\n",
    fg("⌡", config.Yellow))

  fmt.Println(fg("Width:", config.Yellow), bold(strconv.Itoa(width)))
  fmt.Println(fg("=========", config.BrightBlack))
  fmt.Println(test_prompt)
  fmt.Println(fg("=========", config.BrightBlack))
}

// Formatting

func fg(str string, color config.Color) string {
  return fmt.Sprintf("%%F{%d}%s%%f", color, str)
}

func bold(str string) string {
  return fmt.Sprintf("%%B%s%%b", str)
}

func underline(str string) string {
  return fmt.Sprintf("%%U%s%%u", str)
}

func standout(str string) string {
  return fmt.Sprintf("%%S%s%%s", str)
}

