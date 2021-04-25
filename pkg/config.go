package pkg

import (
	"gopkg.in/yaml.v2"
)

type GithubLogin = string
type AsanaGID = string

type Config struct {
	Accounts map[GithubLogin]AsanaGID `yaml:"accounts"`
}

func loadConfig(data []byte) (*Config, error) {
	c := &Config{}

	err := yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
