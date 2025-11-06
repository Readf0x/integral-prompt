package config

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/invopop/jsonschema"
)

func (r *Char) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if len(s) != 1 {
		return errors.New("expected string of length 1")
	}
	*r = Char(s[0])
	return nil
}

func (r *Char) MarshalJSON() ([]byte, error) {
	return []byte(`"`+string(*r)+`"`), nil
}

var char_length uint64 = 1
func (Char) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		MaxLength: &char_length,
	}
}

func (r *Color) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "black":
		*r = Black
		return nil
	case "red":
		*r = Red
		return nil
	case "green":
		*r = Green
		return nil
	case "yellow":
		*r = Yellow
		return nil
	case "blue":
		*r = Blue
		return nil
	case "magenta":
		*r = Magenta
		return nil
	case "cyan":
		*r = Cyan
		return nil
	case "white":
		*r = White
		return nil
	case "bright_black":
		*r = BrightBlack
		return nil
	case "bright_red":
		*r = BrightRed
		return nil
	case "bright_green":
		*r = BrightGreen
		return nil
	case "bright_yellow":
		*r = BrightYellow
		return nil
	case "bright_blue":
		*r = BrightBlue
		return nil
	case "bright_magenta":
		*r = BrightMagenta
		return nil
	case "bright_cyan":
		*r = BrightCyan
		return nil
	case "bright_white":
		*r = BrightWhite
		return nil
	}
	return errors.New("Not a valid color string")
}

func (r *Color) MarshalJSON() ([]byte, error) {
	switch *r {
	case Black:
		return []byte(`"black"`), nil
	case Red:
		return []byte(`"red"`), nil
	case Green:
		return []byte(`"green"`), nil
	case Yellow:
		return []byte(`"yellow"`), nil
	case Blue:
		return []byte(`"blue"`), nil
	case Magenta:
		return []byte(`"magenta"`), nil
	case Cyan:
		return []byte(`"cyan"`), nil
	case White:
		return []byte(`"white"`), nil
	case BrightBlack:
		return []byte(`"bright_black"`), nil
	case BrightRed:
		return []byte(`"bright_red"`), nil
	case BrightGreen:
		return []byte(`"bright_green"`), nil
	case BrightYellow:
		return []byte(`"bright_yellow"`), nil
	case BrightBlue:
		return []byte(`"bright_blue"`), nil
	case BrightMagenta:
		return []byte(`"bright_magenta"`), nil
	case BrightCyan:
		return []byte(`"bright_cyan"`), nil
	case BrightWhite:
		return []byte(`"bright_white"`), nil
	}
	return []byte{}, errors.New("Not a valid color string")
}

func (Color) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			"black",
			"red",
			"green",
			"yellow",
			"blue",
			"magenta",
			"cyan",
			"white",
			"bright_black",
			"bright_red",
			"bright_green",
			"bright_yellow",
			"bright_blue",
			"bright_magenta",
			"bright_cyan",
			"bright_white",
		},
	}
}

func (h HostMap) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]json.RawMessage, len(h))
	for k, v := range h {
		b, err := v.MarshalJSON()
		if err != nil {
			return nil, err
		}
		tmp[k] = b
	}
	return json.Marshal(tmp)
}

func (h *HostMap) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	res := make(map[string]Color, len(raw))
	for k, v := range raw {
		var c Color
		if err := c.UnmarshalJSON(v); err != nil {
			return err
		}
		res[k] = c
	}
	*h = res
	return nil
}

// Returns a pointer to the default configuration.
func GetDefault() *PromptConfig { return &defaultConfig }

func LoadConfig(paths []string) *PromptConfig {
	conf := GetDefault()
	for _, p := range paths {
		file, err := os.Open(p)
		if err == nil {
			b, err := io.ReadAll(file)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			json.Unmarshal(b, conf)
			return conf
		}
	}
	return conf
}

