package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads and unmarshals a YAML config file into GeneratorConfig.
func Load(path string) (*GeneratorConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: cannot read %s: %w", path, err)
	}
	var cfg GeneratorConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: invalid YAML in %s: %w", path, err)
	}
	return &cfg, nil
}
