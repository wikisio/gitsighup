package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Services  []Service `yaml:"services"`
	ConfigSrv ConfigSrv `yaml:"configsrv"`
}

type Service struct {
	Name       string       `yaml:"name"`
	NameSpace  string       `yaml:"namespace"`
	ConfigPath []ConfigPath `yaml:"configPath"`
}

type ConfigPath struct {
	Src string `yaml:"src"`
	Dst string `yaml:"dst"`
}

type ConfigSrv struct {
	EndPoint string `yaml:"endpoint"`
}

type Request struct {
	Username string
	Password string
	Url      string
	Type     string
}

var ApiUrl = string("http://127.0.0.1:3000/api/v1/configsrv/")

var RequestMsg = Request{
	Username: "testuser",
	Password: "123456",
	Url:      "http://127.0.0.1:3000/api/v1/session/login",
	Type:     "POST",
}

var GlobalConfigFile string
var GlobalConfig *Config

func LoadConfig() error {

	var f, err = os.Open(GlobalConfigFile)
	if err != nil {
		fmt.Printf("Failed to open config file, %s, error: %v\n", GlobalConfigFile, err)
		return nil
	}
	defer f.Close()

	var config Config
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		fmt.Printf("Failed to read yaml content, error: %v\n", err)
		return nil
	}

	fmt.Printf("config file %s loaded\n", GlobalConfigFile)
	GlobalConfig = &config
	return nil
}
