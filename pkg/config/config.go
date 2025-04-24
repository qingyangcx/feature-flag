package config

import (
	"encoding/json"
	"feature-flag/pkg/logger"
	"log"
	"os"
	"path/filepath"
)

type Database struct {
	Url      string `json:"url"`
	Db       string `json:"db"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type FeatureFlag struct {
	Port int `json:"port"`
}

type Config struct {
	LoggerConfig logger.LoggerConfig `json:"loggerConfig"`
	FeatureFlag  FeatureFlag         `json:"featureFlag"`
	Database     *Database           `json:"database"`
}

var GlobalConfig = &Config{}
var GlobalWorkDir string = ""

func LoadConfig(workDir, configPath string) error {
	if workDir == "" {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		GlobalWorkDir = filepath.Dir(filepath.Dir(exePath))
	} else {
		GlobalWorkDir = workDir
	}
	log.Println("work dir=", GlobalWorkDir)
	configPathAbs := GlobalWorkDir + "/" + configPath
	content, err := os.ReadFile(configPathAbs)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(content, GlobalConfig); err != nil {
		return err
	}
	return logger.Init(GlobalWorkDir, &GlobalConfig.LoggerConfig)
}
