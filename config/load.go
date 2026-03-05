package config

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
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

