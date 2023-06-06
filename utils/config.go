package utils

import (
	"encoding/json"
	"os"
	"strings"
)

type ShirazConfig struct {
	ProjectPath        string            `json:"projectPath"`
	CoverageFolderPath string            `json:"coverageFolderPath"`
	Env                map[string]string `json:"env"`
	Ignore             []string          `json:"ignore"`
	IgnoreFiles        []string
	IgnoreFolders      []string
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

	if len(conf.Ignore) > 0 {
		for _, i := range conf.Ignore {
			if strings.Contains(i, ".go") {
				conf.IgnoreFiles = append(conf.IgnoreFiles, i)
			} else {
				conf.IgnoreFolders = append(conf.IgnoreFolders, i)
			}
		}
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