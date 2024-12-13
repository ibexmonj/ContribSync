package commands

import (
	"fmt"
	"github.com/ibexmonj/ContribSync/config"
	"regexp"
)

func ShowConfig(cfg *config.Config) {
	fmt.Println("Loaded Configuration:")
	fmt.Printf("Reminder Time: %s\n", cfg.Reminder.Time)
	fmt.Printf("Reminder Title: %s\n", cfg.Reminder.Title)
	fmt.Printf("Reminder Message: %s\n", cfg.Reminder.Message)
	fmt.Printf("Jira Enabled: %t, Base URL: %s\n", cfg.Plugins.Jira.Enabled, cfg.Plugins.Jira.BaseURL)
	fmt.Printf("GitHub Enabled: %t, API Token: %s\n", cfg.Plugins.GitHub.Enabled, cfg.Plugins.GitHub.APIToken)
}

func SetConfig(cfg *config.Config, key, value string) error {
	switch key {
	case "reminder.time":
		// Validate time format
		if err := validateTimeFormat(value); err != nil {
			return err
		}
		cfg.Reminder.Time = value
	case "reminder.title":
		// Ensure title is not empty
		if value == "" {
			return fmt.Errorf("reminder title cannot be empty")
		}
		cfg.Reminder.Title = value
	case "reminder.message":
		// Ensure message is not empty
		if value == "" {
			return fmt.Errorf("reminder message cannot be empty")
		}
		cfg.Reminder.Message = value
	case "plugins.jira.enabled":
		cfg.Plugins.Jira.Enabled = (value == "true")
	case "plugins.jira.base_url":
		cfg.Plugins.Jira.BaseURL = value
	case "plugins.github.api_token":
		cfg.Plugins.GitHub.APIToken = value
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	// Use the imported package explicitly for SaveConfig
	return config.SaveConfig(cfg)
}

// validateTimeFormat checks if the provided time string matches the HH:MM format
func validateTimeFormat(timeStr string) error {
	timeFormat := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`) // Matches 00:00 to 23:59
	if !timeFormat.MatchString(timeStr) {
		return fmt.Errorf("invalid time format: %s (expected HH:MM)", timeStr)
	}
	return nil
}
