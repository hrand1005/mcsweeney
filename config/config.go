package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mcsweeney/content"
)

// TODO: More descriptive name
// Contains the full configuration of a content strategy
type Config struct {
	Name        string      `yaml:"name"`
	Intro       Intro       `yaml:"intro"`
	Outro       Outro       `yaml:"outro"`
	Source      Source      `yaml:"source"`
	Destination Destination `yaml:"destination"`
	Filters     Filters     `yaml:"filters"`
	Options     Options     `yaml:"options"`
}

// Contains intro information for prepended content
type Intro struct {
	Path         string  `yaml:"path"`
	Duration     float64 `yaml:"duration"`
	OverlayStart float64 `yaml:"overlay-start"`
	Font         string  `yaml:"font"`
}

// Contains outro information for prepended content
type Outro struct {
	Path     string  `yaml:"path"`
	Duration float64 `yaml:"duration"`
}

// Contains required fields to pull raw content from platform
type Source struct {
	Platform    content.Platform `yaml:"platform"`
	Credentials string           `yaml:"credentials"`
	Query       Query            `yaml:"query"`
}

// Contains query arguments to be used to gather content
type Query struct {
	GameID string `yaml:"game-id"`
	First  int    `yaml:"first"`
	Days   int    `yaml:"days"`
}

// Contains filters for source material
type Filters struct {
	Language  string   `yaml:"language"`
	Blacklist []string `yaml:"blacklist"`
}

// Contains required fields to push content to platform
type Destination struct {
	Platform    content.Platform `yaml:"platform"`
	Credentials string           `yaml:"credentials"`
	Title       string           `yaml:"title"`
	Description string           `yaml:"description"`
	CategoryID  string           `yaml:"category-id"`
	Keywords    string           `yaml:"keywords"`
	Privacy     content.Privacy  `yaml:"privacy"`
	TokenCache  string           `yaml:"token-cache"`
}

// Contains options for content editing
type Options struct {
	Overlay Overlay `yaml:"overlay"`
}

// Contains fields for applying overlays to content
type Overlay struct {
	Font       string `yaml:"font"`
	Background string `yaml:"background"`
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
