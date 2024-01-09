package utils

import (
	"encoding/json"
	"os"
	"strings"
)

type testConifg struct {
	Command string `json:"command"`
	Output  string `json:"output"`
}

type ShirazConfig struct {
	Test               testConifg        `json:"test"`
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

func GetDefaultConfig() ShirazConfig {
	return ShirazConfig{
		Test: testConifg{
			Command: "go test -v ./...",
			Output:  "pkgname",
		},
		ProjectPath:        ".",
		CoverageFolderPath: "./coverage/",
		Ignore:             make([]string, 0),
	}
}

func GetConfigOrDefault() ShirazConfig {
	userDefined, uE := GetConfig()
	defaultConf := GetDefaultConfig()

	if uE != nil {
		return GetDefaultConfig()
	}

	if userDefined.CoverageFolderPath == "" {
		userDefined.CoverageFolderPath = defaultConf.CoverageFolderPath
	}

	if userDefined.ProjectPath == "" {
		userDefined.ProjectPath = defaultConf.ProjectPath
	}

	if userDefined.Test.Command == "" {
		userDefined.Test.Command = defaultConf.Test.Command
	}

	if userDefined.Test.Output == "" {
		userDefined.Test.Output = defaultConf.Test.Output
	}

	return userDefined
}
