package config

import (
    "fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)


// TODO: More descriptive name
type Config struct {
	Source   string `yaml:"source"`
	ClientID string `yaml:"clientID"`
	Token    string `yaml:"token"`
	GameID   string `yaml:"gameID"`
	First    int    `yaml:"first"`
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

