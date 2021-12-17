package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// TODO: More descriptive name
type Config struct {
	Source      string `yaml:"source"`
	ClientID    string `yaml:"clientID"`
	Token       string `yaml:"token"`
    Query   Query `yaml:"query"`
	Destination string `yaml:"destination"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Category    string `yaml:"category"`
	Keywords    string `yaml:"keywords"`
	Privacy     string `yaml:"privacy"`
}

type Query struct {
    GameID string `yaml:"gameID"`
	First       int    `yaml:"first"`
	StartTime   string `yaml:"started_at"`
}

// It may be appropriate to get more information than just a token
func NewConfig(path string) (*Config, error) {
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
