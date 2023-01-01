package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Hooks struct {
	Activate   []string `yaml:"activate"`
	Deactivate []string `yaml:"deactivcate"`
}

type Display struct {
	Resolution string `yaml:"resolution"`
	Primary    bool   `yaml:"primary"`
	Rotation   string `yaml:"rotation"`
	Order      []struct {
		Display  string `yaml:"display"`
		Position string `yaml:"position"`
	} `yaml:"order"`
}

type Workspace struct {
	Hooks    Hooks              `yaml:"hooks"`
	Displays map[string]Display `yaml:"displays"`
}

type WorkspaceSwitcherConfiguration struct {
	Hooks      Hooks                `yaml:"hooks"`
	Aliases    map[string]string    `yaml:"aliases"`
	Workspaces map[string]Workspace `yaml:"workspaces"`
}

func loadConfig(path string) (*WorkspaceSwitcherConfiguration, error) {
	var config WorkspaceSwitcherConfiguration
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	return &config, err
}

func (c *WorkspaceSwitcherConfiguration) GetWorkspaceNames() []string {
	var result []string
	for key := range c.Workspaces {
		result = append(result, key)
	}

	return result
}
