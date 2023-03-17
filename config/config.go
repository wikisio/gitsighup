package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Services []Service
}

type Service struct {
	Name       string `yaml:"name"`
	ConfigPath string `yaml:"configPath"`
}

var GlobalConfig *Config

func LoadConfig() error {
	var configFile = flag.String("c", "", "the yaml config file")
	flag.Parse()

	var f, err = os.Open(*configFile)
	if err != nil {
		fmt.Printf("Failed to open config file, %s, error: %v", configFile, err)
		return nil
	}
	defer f.Close()

	var config Config
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		fmt.Printf("Failed to read yaml content, error: %v", err)
		return nil
	}

	fmt.Printf("config file %s loaded", *configFile)
	GlobalConfig = &config
	return nil
}
