package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/light-speak/lighthouse/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Schema SchemaConfig `yaml:"schema"`
}

type SchemaConfig struct {
	Ext  []string `yaml:"ext"`
	Path []string `yaml:"path"`
}

func ReadConfig(path string) (*Config, error) {
	var configFile string
	for _, ext := range []string{"yaml", "yml"} {
		filePath := filepath.Join(path, fmt.Sprintf("lighthouse.%s", ext))
		if _, err := os.Stat(filePath); err == nil {
			configFile = filePath
			break
		}
	}

	if configFile == "" {
		return nil, &errors.ConfigError{Message: "config file not found, please create a lighthouse.yaml(yml) file"}
	}

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, &errors.ConfigError{Message: err.Error()}
	}

	return &config, nil
}
