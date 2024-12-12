package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config structure to represent the YAML config file
type Config struct {
	Reminder struct {
		Time string `yaml:"time"`
	} `yaml:"reminder"`
	Plugins struct {
		Jira struct {
			Enabled bool   `yaml:"enabled"`
			BaseURL string `yaml:"base_url"`
		} `yaml:"jira"`
		GitHub struct {
			Enabled  bool   `yaml:"enabled"`
			APIToken string `yaml:"api_token"`
		} `yaml:"github"`
	} `yaml:"plugins"`
}

// LoadConfig reads the config.yaml file and parses it into the Config struct
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SaveConfig writes the updated Config struct back to the YAML file
func SaveConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile("config.yaml", data, 0644)
}
