package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
)

// Config structure to represent the YAML config file
type Config struct {
	Reminder struct {
		Time    string `yaml:"time"`
		Title   string `yaml:"title"`
		Message string `yaml:"message"`
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
		// If the file doesn't exist, create a new config with defaults
		if os.IsNotExist(err) {
			cfg := &Config{}
			SetDefaults(cfg) // Populate defaults
			if saveErr := SaveConfig(cfg); saveErr != nil {
				return nil, saveErr
			}
			return cfg, nil
		}
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Apply default values to missing fields in an existing config
	SetDefaults(&config)

	// Optionally save updated config with defaults filled in
	if saveErr := SaveConfig(&config); saveErr != nil {
		return nil, saveErr
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

// ValidateConfig checks if the configuration values are valid
func ValidateConfig(cfg *Config) error {
	// Validate reminder time format (HH:MM)
	timeFormat := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`) // Matches 00:00 to 23:59
	if !timeFormat.MatchString(cfg.Reminder.Time) {
		return fmt.Errorf("invalid reminder time format: %s (expected HH:MM)", cfg.Reminder.Time)
	}

	// Add more validations here as 	needed
	return nil
}

func SetDefaults(cfg *Config) {
	// Set default reminder time
	if cfg.Reminder.Time == "" {
		cfg.Reminder.Time = "17:00"
	}
	// Set default reminder title
	if cfg.Reminder.Title == "" {
		cfg.Reminder.Title = "Contribution Reminder"
	}
	// Set default reminder message
	if cfg.Reminder.Message == "" {
		cfg.Reminder.Message = "Don't forget to log your contributions!"
	}

	// Add defaults for other fields as needed
	if cfg.Plugins.Jira.BaseURL == "" {
		cfg.Plugins.Jira.BaseURL = ""
	}
	if cfg.Plugins.GitHub.APIToken == "" {
		cfg.Plugins.GitHub.APIToken = ""
	}
}
