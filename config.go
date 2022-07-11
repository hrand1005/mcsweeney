package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CategoryID  string `yaml:"category-id"`
	Days        int    `yaml:"days"`
	DB          string `yaml:"database"`
	Description string `yaml:"description"`
	GameID      string `yaml:"game-id"`
	First       int    `yaml:"first"`
	Tags        string `yaml:"tags"`
	Title       string `yaml:"title"`
}

func LoadConfig(path string) (Config, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return Config{}, fmt.Errorf("LoadConfig: failed to load config from file: %v", err)
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return Config{}, fmt.Errorf("LoadConfig: failed to read file: %v", err)
	}

	conf := &Config{}
	err = yaml.Unmarshal(raw, conf)
	if err != nil {
		return Config{}, fmt.Errorf("LoadConfig: failed to unmarshal to yaml: %v", err)
	}

	return *conf, nil
}
