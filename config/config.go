package config

type Config struct {
	Services []Service
}

type Service struct {
	Name       string `yaml:"name"`
	ConfigPath string `yaml:"configPath"`
}
