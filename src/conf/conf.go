package config

const Version = "v0.1"

// [TODO] add jsonschema info

type IconEntry struct {
  Color Color  `json:"color,omitempty"`
  Icon  rune   `json:"icon"`
}

type SingleIconConfig struct {
  Color Color `json:"color,omitempty"`
  Icon  rune  `json:"icon,omitempty"`
}

type IconConfig struct {
  Color Color  `json:"color,omitempty"`
  Icon  rune   `json:"icon,omitempty"`
  Icons []rune `json:"icons,omitempty"`

  DefaultIcon *SingleIconConfig `json:"default_icon,omitempty"`
  IconEntries []IconEntry       `json:"icons,omitempty"`
}

type BatteryConfig struct {
  Color       Color                `json:"color,omitempty"`
  Icons       [2]rune              `json:"icons,omitempty"`
  IconEntries [2]BatteryIconConfig `json:"icons,omitempty"`
}
type BatteryIconConfig struct {
  Charging    IconEntry `json:"charging"`
  Discharging IconEntry `json:"discharging"`
}

type CounterConfig struct {
  Color Color `json:"color"`
  Icon  rune  `json:"icon"`
}

type LineConfig struct {
  Color   Color   `json:"color"`
  Symbols [3]rune `json:"symbols"`
}

type GitConfig struct {
  Color   Color         `json:"color"`
  Icon    rune          `json:"icon"`
  Changes CounterConfig `json:"changes"`
  Push    CounterConfig `json:"push"`
  Pull    CounterConfig `json:"pull"`
}

type PromptConfig struct {
  Version      string     `json:"version,omitempty"`
  ModulesLeft  []string   `json:"modules_left,omitempty"`
  ModulesRight []string   `json:"modules_right,omitempty"`
  Line         LineConfig `json:"line,omitzero"`
  // Per Module Config
  Battery      IconConfig `json:"battery,omitzero"`
  Direnv       IconConfig `json:"direnv,omitzero"`
  Error        IconConfig `json:"error,omitzero"`
  NixShell     IconConfig `json:"nix_shell,omitzero"`
  ViMode       IconConfig `json:"vi_mode,omitzero"`
  Git          GitConfig  `json:"git,omitzero"`
}

func GetDefault() PromptConfig {
  return PromptConfig{
    Version: Version,
    ModulesLeft: []string{
      "direnv",
      "nix",
      "visym",
      "error",
      "dir",
      "sshplus",
      "git",
      "jobs",
    },
    ModulesRight: []string{
      "time",
    },
    Line: LineConfig{
      Color: Yellow,
      Symbols: [3]rune{ '⌠', '⎮', '⌡' },
    },
    Battery: IconConfig{
      Color: BrightGreen,
      IconEntries: []IconEntry{
      },
    },
  }
}

