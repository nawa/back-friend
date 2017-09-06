package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Port  int
	SQLDb *SQLDb `yaml:"sqlDb"`
}

type SQLDb struct {
	Type, Host, Port, Name, User, Password string
}

// FromFile : creates config from file
func FromFile(cfgLocation string) (*AppConfig, error) {
	yamlFile, err := ioutil.ReadFile(cfgLocation)
	if err != nil {
		return nil, fmt.Errorf("Fail to read config file: %v", err)
	}
	cfg, err := NewConfig(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("Invalid config: %v", err)
	}
	return cfg, nil
}

// NewConfig : creates config from bytes
func NewConfig(yml []byte) (cfg *AppConfig, err error) {
	cfg = new(AppConfig)
	err = yaml.Unmarshal(yml, cfg)
	return
}
