package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type BldrConfig struct {
	Version int          `yaml:"version"`
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

	return &c, c.validate()
}

func (c *BldrConfig) validate() error {
	if c.Version != 0 {
		return fmt.Errorf("bldr version %d is not valid", c.Version)
	}

	if len(c.Structs) == 0 {
		return fmt.Errorf("no structs found in YAML file")
	}

	// Ensure we are not writing two structs to the same target output location.
	outLocations := make([]string, 0, len(c.Structs))
	for _, st := range c.Structs {
		if contains(outLocations, st.Go.Output) {
			return fmt.Errorf("struct %s has duplicate output: %s", st.Name, st.Go.Output)
		}
		outLocations = append(outLocations, st.Go.Output)
	}

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
