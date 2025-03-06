package config

import (
	"fmt"
	"github.com/spf13/viper"
	"regexp"
)

type Config struct {
	Reminder struct {
		Time    string `mapstructure:"time"`
		Title   string `mapstructure:"title"`
		Message string `mapstructure:"message"`
	} `mapstructure:"reminder"`
	Plugins struct {
		Jira struct {
			Enabled bool   `mapstructure:"enabled"`
			BaseURL string `mapstructure:"base_url"`
		} `mapstructure:"jira"`
		GitHub struct {
			Enabled  bool   `mapstructure:"enabled"`
			APIToken string `mapstructure:"api_token"`
		} `mapstructure:"github"`
	} `mapstructure:"plugins"`
}

var ConfigData Config

func LoadConfig() error {
	viper.SetConfigName("config") // Name of the file (without extension)
	viper.SetConfigType("yaml")   // File type
	viper.AddConfigPath(".")      // Look in the current directory
	viper.AutomaticEnv()          // Support environment variables

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("⚠️  No config file found. Creating default config.yaml...")
		if err := SaveConfig(); err != nil {
			return fmt.Errorf("failed to save default config: %w", err)
		}
	}

	if err := viper.Unmarshal(&ConfigData); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	if err := ValidateConfig(&ConfigData); err != nil {
		return err
	}

	return nil
}

func SaveConfig() error {
	if err := viper.WriteConfigAs("config.yaml"); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func ValidateConfig(cfg *Config) error {
	timeFormat := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`) // Matches 00:00 to 23:59
	if !timeFormat.MatchString(cfg.Reminder.Time) {
		return fmt.Errorf("invalid reminder time format: %s (expected HH:MM)", cfg.Reminder.Time)
	}
	return nil
}

func setDefaults() {
	viper.SetDefault("reminder.time", "17:00")
	viper.SetDefault("reminder.title", "Contribution Reminder")
	viper.SetDefault("reminder.message", "Don't forget to log your contributions!")

	viper.SetDefault("plugins.jira.enabled", false)
	viper.SetDefault("plugins.jira.base_url", "")

	viper.SetDefault("plugins.github.enabled", false)
	viper.SetDefault("plugins.github.api_token", "")
}
