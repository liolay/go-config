package util

import (
	"gopkg.in/yaml.v2"
)

type ClientConfig struct {
	Server string    `yaml:"server"`
	Tick   int       `yaml:"tick"`
	App    []AppNode `yaml:"app"`
}

type AppNode struct {
	Name     string   `yaml:"name"`
	Profile  string   `yaml:"profile"`
	Label    string   `yaml:"label"`
	HomePath []string `yaml:"homePath"`
}

type ServerConfig struct {
	HomePath     string `yaml:"homePath"`
	Port         string    `yaml:"port"`
	DefaultRepo  string `yaml:"defaultRepo"`
	SshKey       string `yaml:"sshKey"`
	SearchSubDir bool   `yaml:"searchSubDir"`

	Route []RRoute `yaml:"route"`
}

type RRoute struct {
	Pattern []string `yaml:"pattern"`
	Repo    string   `yaml:"repo"`
	Model   RepoModel
}

const (
	OnlyOne       RepoModel = 0
	AppOne        RepoModel = 1
	AppProfileOne RepoModel = 2
)

type RepoModel int

func parseConfig(content []byte, t interface{}) (error) {
	err := yaml.Unmarshal(content, t)
	return err
}

func ParseClientConfig(content []byte) (*ClientConfig, error) {
	config := new(ClientConfig)
	if err := parseConfig(content, config); err != nil {
		return nil, err
	}
	return config, nil
}

func ParseServerConfig(content []byte) (*ServerConfig, error) {
	config := new(ServerConfig)
	if err := parseConfig(content, config); err != nil {
		return nil, err
	}
	return config, nil
}
