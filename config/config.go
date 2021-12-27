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
	Intro       Intro       `yaml:"intro"`
	Source      Source      `yaml:"source"`
	Destination Destination `yaml:"destination"`
	Filters     Filters     `yaml:"filters"`
	Options     Options     `yaml:"options"`
}

// Contains intro information for prepended content
type Intro struct {
	Path     string  `yaml:"path"`
	Duration float64 `yaml:"duration"`
}

// Contains required fields to pull raw content from platform
type Source struct {
	Platform    content.ContentType `yaml:"platform"`
	Credentials string              `yaml:"credentials"`
	Query       Query               `yaml:"query"`
}

// Contains query arguments to be used to gather content
type Query struct {
	GameID string `yaml:"gameID"`
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
	Platform    content.ContentType `yaml:"platform"`
	Credentials string              `yaml:"credentials"`
	Title       string              `yaml:"title"`
	Description string              `yaml:"description"`
	Category    string              `yaml:"category"`
	Keywords    string              `yaml:"keywords"`
	Privacy     content.Privacy     `yaml:"privacy"`
}

// Contains options for content editing
type Options struct {
	Overlay content.Overlay `yaml:"overlay"`
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
