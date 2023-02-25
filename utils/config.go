package utils

import (
	"encoding/json"
	"os"
)

type ShirazConfig struct {
	ProjectPath        string   `json:"projectPath"`
	CoverageFolderPath string   `json:"coverageFolderPath"`
	Ignore             []string `json:"ignore"`
}

func GetConfig() (ShirazConfig, error) {
	cb, cbErr := os.ReadFile("./shiraz.json")
	if cbErr != nil {
		return ShirazConfig{}, cbErr
	}

	var conf ShirazConfig
	uErr := json.Unmarshal(cb, &conf)
	if uErr != nil {
		return ShirazConfig{}, uErr
	}

	if conf.ProjectPath == "" {
		conf.ProjectPath = "."
	}

	if conf.CoverageFolderPath == "" {
		conf.CoverageFolderPath = "./coverage/"
	}

	return conf, nil
}

func GetConfigOrDefault() ShirazConfig {
	userDefined, uE := GetConfig()

	if uE == nil {
		return userDefined
	}

	return ShirazConfig{
		ProjectPath:        ".",
		CoverageFolderPath: "./coverage/",
		Ignore:             make([]string, 0),
	}
}
