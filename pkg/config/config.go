package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AccessToken string                `env:"ACCESS_TOKEN,required"`
	Repos       map[string]Deployment `yaml:"repos"`
	Port        int                   `yaml:"port" env:"PORT" env-default:"80"`
}

type Deployment struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

func GetConfig(filename string) (*Config, error) {
	config := &Config{}
	err := cleanenv.ReadConfig(filename, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
