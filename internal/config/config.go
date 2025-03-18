package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	App AppConfig `yaml:"app"`
}

type AppConfig struct {
	APIUrl string `yaml:"api_url"`

	ProjectURLPath string `yaml:"project_url_path"`
	APIAuthToken   string `yaml:"api_auth_token"`
	WorkspaceGID   string `yaml:"workspace_gid"`

	UserConfig    UserConfig    `yaml:"user"`
	ProjectConfig ProjectConfig `yaml:"project"`
}

type ProjectConfig struct {
	Path  string `yaml:"path"`
	Limit uint8  `yaml:"limit"`
}

type UserConfig struct {
	Path  string `yaml:"path"`
	Limit uint8  `yaml:"limit"`
}

func LoadConfig(configFilePath string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(configFilePath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
