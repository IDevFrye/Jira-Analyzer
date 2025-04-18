package configreader

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

func LoadConfig(r io.Reader) (Config, error) {
	var cfg Config
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse config error: %w", err)
	}

	return cfg, nil
}
