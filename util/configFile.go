package util

import (
	"gopkg.in/yaml.v2"
)

type ClientConfig struct {
	HomePath []string `yaml:"homePath"`
	Server   string   `yaml:"server"`
	Tick     int      `yaml:"tick"`
	Apps []struct {
		Name     string   `yaml:"name"`
		Profile  string   `yaml:"profile"`
		Label    string   `yaml:"lable"`
		RootPath []string `yaml:"rootPath"`
	} `yaml:"apps"`
}

type ServerConfig struct {
	DefaultRepo  string `yaml:"defaultRepo"`
	SshKey       string `yaml:"sshKey"`
	SearchSubDir bool   `yaml:"searchSubDir"`

	Route []struct {
		Pattern []string `yaml:"pattern"`
		Repo    string   `yaml:"repo"`
	} `yaml:"route"`
}

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
