package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AccessToken string `env:"ACCESS_TOKEN,required"`
	Repos       []Repo `yaml:"repos"`
	Port        int    `yaml:"port" env:"PORT" env-default:"80"`
}

type Repo struct {
	Repo       string `yaml:"repo"`
	Deployment string `yaml:"deployment"`
}

func GetConfig(filename string) (*Config, error) {
	config := &Config{}
	err := cleanenv.ReadConfig(filename, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
