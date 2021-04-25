package pkg

import (
	"gopkg.in/yaml.v2"
)

type githubLogin = string
type asanaGID = string

type config struct {
	Accounts map[githubLogin]asanaGID `yaml:"accounts"`
}

func loadConfig(data []byte) (*config, error) {
	c := &config{}

	err := yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
