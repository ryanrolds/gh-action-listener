package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Schema      string                `yaml:"schema" env-required:"true"`
	AccessToken string                `env:"ACCESS_TOKEN" env-required:"true"`
	Deployments map[string]Deployment `yaml:"deployments"`
	Resources   map[string]Resource   `yaml:"resources"`
	Port        int                   `yaml:"port" env:"PORT" env-default:"80"`
}

type Deployment struct {
	DeploymentName string `yaml:"deployment_name"`
	ContainerName  string `yaml:"container_name"`
	Namespace      string `yaml:"namespace"`
	Image          string `yaml:"image"`
}

type Resource struct {
	ID        string `yaml:"id"`
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

func GetConfig(filename string) (*Config, error) {
	config := &Config{}
	err := cleanenv.ReadConfig(filename, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
