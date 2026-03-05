package config

import (
	"encoding/json"
	"os"
)

// Returns a pointer to the default configuration.
func GetDefault() *PromptConfig { return &defaultConfig }

func LoadConfig(paths []string) *PromptConfig {
	conf := GetDefault()
	for _, p := range paths {
		file, err := os.Open(p)
		if err == nil {
			var b []byte
			file.Read(b)
			defer file.Close()

			json.Unmarshal(b, conf)
			return conf
		}
	}
	return conf
}

