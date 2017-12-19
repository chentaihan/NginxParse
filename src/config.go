package main

import "github.com/BurntSushi/toml"

type Config struct {
	ParseList []StructPraseInfo
}

type StructPraseInfo struct {
	StructName  string
	Include     []string
	IncludeNo   []string
	Fields      []int    //为空表示所有字段
	FieldNames  []string //存Fields中指定的字段名称
	MacroExpand bool     //宏展开
}

var ConfigInfo Config

func LoadConfig() error {
	configFile := getConfigFile(NGINX_PARSE)
	if _, err := toml.DecodeFile(configFile, &ConfigInfo); err != nil {
		return err
	}

	for i, _ := range ConfigInfo.ParseList {
		tmpInclude := make([]string, 0, len(ConfigInfo.ParseList[i].Include))
		name := ConfigInfo.ParseList[i].StructName
		name = name[0:len(name)-1]
		tmpInclude = append(tmpInclude, name+"t")
		tmpInclude = append(tmpInclude, ConfigInfo.ParseList[i].Include...)
		ConfigInfo.ParseList[i].Include = tmpInclude
	}
	return nil
}
