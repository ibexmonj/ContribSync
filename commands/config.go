package commands

import (
	"fmt"
	"github.com/ibexmonj/ContribSync/config"
)

func ShowConfig(config *config.Config) {
	fmt.Println("Loaded Configuration:")
	fmt.Printf("Reminder Time: %s\n", config.Reminder.Time)
	fmt.Printf("Jira Enabled: %t, Base URL: %s\n", config.Plugins.Jira.Enabled, config.Plugins.Jira.BaseURL)
	fmt.Printf("GitHub Enabled: %t, API Token: %s\n", config.Plugins.GitHub.Enabled, config.Plugins.GitHub.APIToken)
}

func SetConfig(cfg *config.Config, key, value string) error {
	switch key {
	case "reminder.time":
		cfg.Reminder.Time = value
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
