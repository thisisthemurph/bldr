package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type BldrConfig struct {
	Version string       `yaml:"version"`
	Structs []StructInfo `yaml:"structs"`
}

type StructInfo struct {
	Name string     `yaml:"name"`
	Path string     `yaml:"path"`
	Go   GoResource `yaml:"go"`
}

type GoResource struct {
	Package string `yaml:"package"`
	Output  string `yaml:"out"`
}

func Read() (*BldrConfig, error) {
	data, err := os.ReadFile("test.yml")
	if err != nil {
		return nil, err
	}

	c := BldrConfig{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
