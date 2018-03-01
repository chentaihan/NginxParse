package util

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	StructName     []string
	SourceCodeUrl  string
}

var ConfigInfo Config

func LoadConfig() error {
	configFile := GetConfigFile(NGINX_PARSE)
	if _, err := toml.DecodeFile(configFile, &ConfigInfo); err != nil {
		return err
	}
	return nil
}
