package config

import "time"

const Version = "v0.1"

// [TODO] add jsonschema info

type IconEntry struct {
	Name  string `json:"name,omitempty"`
	Color Color  `json:"color,omitempty"`
	Icon  rune   `json:"icon,omitempty"`
}

type SingleIconEntry struct {
	Color Color `json:"color,omitempty"`
	Icon  rune  `json:"icon,omitempty"`
}

type IconConfig struct {
	DefaultIcon *SingleIconEntry `json:"default_icon,omitempty"`
	IconEntries *[]IconEntry      `json:"icons,omitempty"`
}

type BatteryConfig struct {
	Id          int                `json:"battery,omitempty"`
	DefaultIcon *SingleIconEntry `json:"default_icon,omitempty"`
	IconEntries *BatteryIconConfig `json:"icons,omitempty"`
}
type BatteryIconConfig struct {
	Charging    SingleIconEntry `json:"charging"`
	Discharging SingleIconEntry `json:"discharging"`
}

type ErrorConfig struct {
	DefaultIcon *SingleIconEntry `json:"default_icon,omitempty"`
	IconEntries *[]ErrorEntry     `json:"icons,omitempty"`
}
type ErrorEntry struct {
	Code  uint64 `json:"code,omitempty"`
	Color Color  `json:"color,omitempty"`
	Icon  rune   `json:"icon,omitempty"`
}

type ViModeConfig struct {
	Insert *IconEntry `json:"insert,omitempty"`
	Normal *IconEntry `json:"normal,omitempty"`
	Visual *IconEntry `json:"visual,omitempty"`
	ViLine *IconEntry `json:"visual_line,omitempty"`
}

type CounterConfig struct {
	Color Color `json:"color,omitempty"`
	Icon  rune  `json:"icon,omitempty"`
}

type LineConfig struct {
	Color   Color   `json:"color,omitempty"`
	Symbols [3]rune `json:"symbols,omitempty"`
}

type GitConfig struct {
	Branch   *CounterConfig `json:"branch,omitempty"`
	Unstaged *CounterConfig `json:"unstaged,omitempty"`
	Staged   *CounterConfig `json:"staged,omitempty"`
	Push     *CounterConfig `json:"push,omitempty"`
	Pull     *CounterConfig `json:"pull,omitempty"`
}

type CpuConfig struct {
	Time  time.Duration `json:"time,omitempty"`
	Color Color         `json:"color,omitempty"`
	Icon  rune          `json:"icon,omitempty"`
}

type DirConfig struct {
	Color       Color         `json:"color,omitempty"`
	Replace     *[]*[2]string `json:"replace,omitempty"`
	ReplaceHome bool          `json:"replace_home"`
	HomeIcon    rune          `json:"home_icon"`
}

type DistroboxConfig struct {
	TextColor   Color             `json:"color,omitempty"`
	DefaultIcon *SingleIconEntry `json:"default_icon,omitempty"`
	IconEntries *[]IconEntry      `json:"icons,omitempty"`
}

type SshConfig struct {
	User *DisplayConfig `json:"user,omitempty"`
	At   *DisplayConfig `json:"at,omitempty"`
	Host *DisplayConfig `json:"host,omitempty"`
}

type DisplayConfig struct {
	Color   Color `json:"color,omitempty"`
	Visible bool  `json:"visible,omitempty"`
}

type TimeConfig struct {
	Format string `json:"format,omitempty"`
	Color  Color  `json:"color,omitempty"`
	Icon   rune   `json:"icon,omitempty"`
}

type PromptConfig struct {
	Version      string      `json:"version,omitempty"`
	Modules      *[]string   `json:"modules,omitempty"`
	ModulesRight *[]string   `json:"modules_right,omitempty"`
	RightSize    uint8       `json:"right_width"`
	Length       uint8       `json:"length,omitempty"`
	Line         *LineConfig `json:"line,omitempty"`
	// Per Module Config
	Battery   *BatteryConfig    `json:"battery,omitempty"`
	Cpu       *CpuConfig        `json:"cpu,omitempty"`
	Dir       *DirConfig        `json:"dir,omitempty"`
	Direnv    *IconConfig       `json:"direnv,omitempty"`
	Distrobox *DistroboxConfig  `json:"distrobox,omitempty"`
	Error     *ErrorConfig      `json:"error,omitempty"`
	Git       *GitConfig        `json:"git,omitempty"`
	Jobs      *CounterConfig    `json:"jobs,omitempty"`
	NixShell  *SingleIconEntry  `json:"nix_shell,omitempty"`
	Ssh       *SshConfig        `json:"ssh,omitempty"`
	Time      *TimeConfig       `json:"time,omitempty"`
	Uptime    *CounterConfig    `json:"uptime,omitempty"`
	ViMode    *ViModeConfig     `json:"vi_mode,omitempty"`
}

var defaultConfig = PromptConfig{
	Version: Version,
	Modules: &[]string{
		"direnv",
		"nix",
		"visym",
		"error",
		"dir",
		"ssh+",
		"git",
		"jobs",
	},
	ModulesRight: &[]string{
		"time",
	},
	Line: &LineConfig{
		Color:   Yellow,
		Symbols: [3]rune{'⌠', '⎮', '⌡'},
	},
	Battery: &BatteryConfig{
		Id: 0,
		IconEntries: &BatteryIconConfig{
			Charging: SingleIconEntry{
				Color: BrightGreen,
				Icon:  '🗲',
			},
			Discharging: SingleIconEntry{
				Color: Green,
				Icon:  '󰁹',
			},
		},
	},
	Direnv: &IconConfig{
		DefaultIcon: &SingleIconEntry{
			Color: Magenta,
			Icon:  '⌁',
		},
	},
	Error: &ErrorConfig{
		IconEntries: &[]ErrorEntry{
			{Code: 1,   Color: Red,         Icon: '✘'},
			{Code: 2,   Color: Blue,        Icon: '?'},
			{Code: 127, Color: Blue,        Icon: '?'},
			{Code: 126, Color: Yellow,      Icon: '⚠'},
			{Code: 130, Color: BrightBlack, Icon: '☠'},
			{Code: 137, Color: BrightRed,   Icon: '☠'},
			{Code: 148, Color: Cyan,        Icon: '*'},
		},
		DefaultIcon: &SingleIconEntry{Color: Red, Icon: '✘'},
	},
	Git: &GitConfig{
		Branch: &CounterConfig{
			Color: Yellow,
			Icon:  '⎇',
		},
		Unstaged: &CounterConfig{
			Color: Red,
			Icon:  '✘',
		},
		Staged: &CounterConfig{
			Color: Green,
			Icon:  '+',
		},
		Push: &CounterConfig{
			Color: Cyan,
			Icon:  '↑',
		},
		Pull: &CounterConfig{
			Color: Cyan,
			Icon:  '↓',
		},
	},
	NixShell: &SingleIconEntry{
		Color: Cyan,
		Icon:  '❄',
	},
	ViMode: &ViModeConfig{
		Insert: &IconEntry{
			Color: Green,
			Icon:  '○',
		},
		Normal: &IconEntry{
			Color: Red,
			Icon:  '●',
		},
		Visual: &IconEntry{
			Color: Magenta,
			Icon:  '◒',
		},
		ViLine: &IconEntry{
			Color: Magenta,
			Icon:  '◐',
		},
	},
	Cpu: &CpuConfig{
		Time:  time.Millisecond * 100,
		Color: White,
		Icon:  '%',
	},
	Dir: &DirConfig{
		Color: Blue,
		Replace: &[]*[2]string{},
		ReplaceHome: true,
		HomeIcon: '~',
	},
	Distrobox: &DistroboxConfig{
		TextColor: Green,
		DefaultIcon: &SingleIconEntry{
			Color: White,
			Icon: '',
		},
		IconEntries: &[]IconEntry{
			{ Name: "fedora", Color: Blue,      Icon: '' },
			{ Name: "ubuntu", Color: BrightRed, Icon: '' },
			{ Name: "debian", Color: Red,       Icon: '' },
			{ Name: "arch",   Color: Cyan,      Icon: '' },
		},
	},
	Jobs: &CounterConfig{
		Color: Magenta,
		Icon: '⚙',
	},
	Time: &TimeConfig{
		Format: time.TimeOnly,
		Color: White,
		// Icon: idfk man does anyone even really use this? I know I don't.
		// Look if you *actually* want me to finish the time module and make all the icons work
		// make an issue, because otherwise I'm not touching this thing.
	},
	Ssh: &SshConfig{
		User: &DisplayConfig{
			Color: Magenta,
			Visible: true,
		},
		At: &DisplayConfig{
			Color: BrightBlack,
			Visible: true,
		},
		Host: &DisplayConfig{
			Color: Cyan,
			Visible: true,
		},
	},
}

// Returns a pointer to the default configuration.
func GetDefault() *PromptConfig { return &defaultConfig }
