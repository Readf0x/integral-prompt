package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

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

