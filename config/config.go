package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// TODO: More descriptive name
// Contains the full configuration of a content strategy
type Config struct {
	Source      Source      `yaml:"source"`
	Destination Destination `yaml:"destination"`
	Options     Options     `yaml:"options"`
}

// Contains required fields to pull raw content from platform
type Source struct {
	Platform    string `yaml:"platform"`
	Credentials string `yaml:"credentials"`
	Query       Query  `yaml:"query"`
}

// Contains query arguments to be used to gather content
type Query struct {
	GameID string `yaml:"gameID"`
	First  int    `yaml:"first"`
	Days   int    `yaml:"days"`
}

// Contains required fields to push content to platform
type Destination struct {
	Platform    string `yaml:"platform"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Category    string `yaml:"category"`
	Keywords    string `yaml:"keywords"`
	Privacy     string `yaml:"privacy"`
}

// Contains options for content editing
type Options struct {
	Overlay Overlay `yaml:"overlay"`
}

// Contains configurable overlay fields
type Overlay struct {
	Font     string `yaml:"font"`
	Size     string `yaml:"size"`
	Color    string `yaml:"color"`
	Duration int    `yaml:"duration"`
	Fade     int    `yaml:"fade"`
}

// Loads config from given yaml file, returns Config pointer
func LoadConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read %s: %w", path, err)
	}

	c := Config{}
	err = yaml.Unmarshal(raw, &c)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal: %v", err)
	}

	return &c, nil
}
