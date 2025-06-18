package main

type Color uint8
const (
  Black Color = iota
  Red
  Green
  Yellow
  Blue
  Magenta
  Cyan
  White
  BrightBlack
  BrightRed
  BrightGreen
  BrightYellow
  BrightBlue
  BrightMagenta
  BrightCyan
  BrightWhite
)

type IconEntry struct {
  Name  string `json:"name"`
  Color string `json:"color,omitempty"`
  Icon  rune   `json:"icon"`
}

type SingleIconConfig struct {
  Color string `json:"color,omitempty"`
  Icon  rune   `json:"icon,omitempty"`
}

type IconConfig struct {
  Color string `json:"color,omitempty"`
  Icon  rune   `json:"icon,omitempty"`
  Icons []rune `json:"icons,omitempty"`

  DefaultIcon *SingleIconConfig `json:"default_icon,omitempty"`
  IconEntries []IconEntry       `json:"icons,omitempty"`
}

type PromptConfig struct {
  Version string   `json:"version"`
  Modules []string `json:"modules"`
}

