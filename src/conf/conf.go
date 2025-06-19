package config

const Version = "v0.1"

// [TODO] add jsonschema info

type IconEntry struct {
	Name  string `json:"name,omitempty"`
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
	Id			 		int   						`json:"battery,omitempty"`
	Color       *Color             `json:"color,omitempty"`
	Icons       *[2]rune           `json:"icons,omitempty"`
	IconEntries *BatteryIconConfig `json:"icons,omitzero"`
}
type BatteryIconConfig struct {
	Charging    IconEntry `json:"charging"`
	Discharging IconEntry `json:"discharging"`
}

type ErrorConfig struct {
	Color Color  `json:"color,omitempty"`
	Icon  rune   `json:"icon,omitempty"`
	Icons []rune `json:"icons,omitempty"`

	DefaultIcon SingleIconConfig `json:"default_icon,omitzero"`
	IconEntries []ErrorEntry     `json:"icons,omitempty"`
}
type ErrorEntry struct {
	Code  uint8 `json:"code,omitempty"`
	Color Color `json:"color,omitempty"`
	Icon  rune  `json:"icon"`
}

type ViModeConfig struct {
	Color       Color            `json:"color,omitempty"`
	Icons       [2]rune          `json:"icons,omitempty"`
	IconEntries ViModeIconConfig `json:"icons,omitzero"`
}
type ViModeIconConfig struct {
	Insert     IconEntry `json:"insert"`
	Normal     IconEntry `json:"normal"`
	Visual     IconEntry `json:"visual"`
	VisualLine IconEntry `json:"visual_line"`
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
	Branch   CounterConfig `json:"branch"`
	Unstaged CounterConfig `json:"unstaged"`
	Staged   CounterConfig `json:"staged"`
	Push     CounterConfig `json:"push"`
	Pull     CounterConfig `json:"pull"`
}

type PromptConfig struct {
	Version      string     `json:"version,omitempty"`
	ModulesLeft  []string   `json:"modules_left,omitempty"`
	ModulesRight []string   `json:"modules_right,omitempty"`
	Length       uint8      `json:"length"`
	Line         LineConfig `json:"line,omitzero"`
	// Per Module Config
	Battery  BatteryConfig    `json:"battery,omitzero"`
	Direnv   IconConfig       `json:"direnv,omitzero"`
	Error    ErrorConfig      `json:"error,omitzero"`
	Git      GitConfig        `json:"git,omitzero"`
	NixShell SingleIconConfig `json:"nix_shell,omitzero"`
	ViMode   ViModeConfig     `json:"vi_mode,omitzero"`
}

var defaultConfig = PromptConfig{
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
		Color:   Yellow,
		Symbols: [3]rune{'⌠', '⎮', '⌡'},
	},
	Battery: BatteryConfig{
		Id: 0,
		IconEntries: &BatteryIconConfig{
			Charging: IconEntry{
				Color: BrightGreen,
				Icon:  '🗲',
			},
			Discharging: IconEntry{
				Color: Green,
				Icon:  '󰁹',
			},
		},
	},
	Direnv: IconConfig{
		Color: Magenta,
		Icon:  '⌁',
	},
	Error: ErrorConfig{
		IconEntries: []ErrorEntry{
			{Code: 1, Color: Red, Icon: '✘'},
			{Code: 2, Color: Blue, Icon: '?'},
			{Code: 127, Color: Blue, Icon: '?'},
			{Code: 126, Color: Yellow, Icon: '⚠'},
			{Code: 130, Color: BrightBlack, Icon: '☠'},
			{Code: 137, Color: BrightRed, Icon: '☠'},
			{Code: 148, Color: Cyan, Icon: '*'},
		},
		DefaultIcon: SingleIconConfig{Color: Red, Icon: '✘'},
	},
	Git: GitConfig{
		Branch: CounterConfig{
			Color: Yellow,
			Icon:  '⎇',
		},
		Unstaged: CounterConfig{
			Color: Red,
			Icon:  '✘',
		},
		Staged: CounterConfig{
			Color: Green,
			Icon:  '+',
		},
		Push: CounterConfig{
			Color: Cyan,
			Icon:  '↑',
		},
		Pull: CounterConfig{
			Color: Cyan,
			Icon:  '↓',
		},
	},
	NixShell: SingleIconConfig{
		Color: Cyan,
		Icon:  '❄',
	},
	ViMode: ViModeConfig{
		IconEntries: ViModeIconConfig{
			Insert: IconEntry{
				Color: Green,
				Icon:  '○',
			},
			Normal: IconEntry{
				Color: Red,
				Icon:  '●',
			},
			Visual: IconEntry{
				Color: Magenta,
				Icon:  '◒',
			},
			VisualLine: IconEntry{
				Color: Magenta,
				Icon:  '◐',
			},
		},
	},
}

// Returns a pointer to the default configuration.
func GetDefault() *PromptConfig { return &defaultConfig }
