package main

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type twitchConfig struct {
	GameID string `yaml:"game-id"`
	First  int    `yaml:"first"`
	Days   int    `yaml:"days"`
	DB     string `yaml:"database"`
}

func LoadTwitchConfig(path string) (twitchConfig, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return twitchConfig{}, fmt.Errorf("LoadTwitchConfig: failed to load config from file: %v", err)
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return twitchConfig{}, fmt.Errorf("LoadTwitchConfig: failed to read file: %v", err)
	}

	tConf := &twitchConfig{}
	err = yaml.Unmarshal(raw, tConf)
	if err != nil {
		return twitchConfig{}, fmt.Errorf("LoadTwitchConfig: failed to unmarshal to yaml: %v", err)
	}

	return *tConf, nil
}
