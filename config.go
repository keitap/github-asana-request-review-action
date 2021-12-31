// Package githubasana provides GitHub and Asana pull request task integration.
package githubasana

import (
	"gopkg.in/yaml.v2"
)

type (
	GithubLogin = string
	AsanaGID    = string
)

type Config struct {
	DueDate  int
	Holidays map[string]bool
	Accounts map[GithubLogin]AsanaGID `yaml:"accounts"`
}

func LoadConfig(data []byte) (*Config, error) {
	c := &Config{
		DueDate:  1,
		Holidays: map[string]bool{},
	}

	err := yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
