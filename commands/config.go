package commands

import (
	"fmt"
	"github.com/ibexmonj/ContribSync/config"
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"github.com/spf13/cobra"
	"regexp"
)

func ShowConfig(cfg *config.Config) {
	logger.Logger.Info().Msg("Loaded Configuration")

	fmt.Printf("\n📌 Reminder Settings:\n")
	fmt.Printf("   ⏰ Time: %s\n", cfg.Reminder.Time)
	fmt.Printf("   📝 Title: %s\n", cfg.Reminder.Title)
	fmt.Printf("   💬 Message: %s\n", cfg.Reminder.Message)

	fmt.Printf("\n🔧 Plugin Settings:\n")
	fmt.Printf("   🏷️ Jira: Enabled: %t, Base URL: %s\n", cfg.Plugins.Jira.Enabled, cfg.Plugins.Jira.BaseURL)
	fmt.Printf("   🏷️ GitHub: Enabled: %t, API Token: %s\n", cfg.Plugins.GitHub.Enabled, maskToken(cfg.Plugins.GitHub.APIToken))
}

func SetConfig(cfg *config.Config, key, value string) error {
	switch key {
	case "reminder.time":
		if err := validateTimeFormat(value); err != nil {
			logger.Logger.Warn().Str("key", key).Str("value", value).Msg("Invalid time format")
			return err
		}
		cfg.Reminder.Time = value
	case "reminder.title":
		if value == "" {
			logger.Logger.Warn().Str("key", key).Msg("Empty title value")
			return fmt.Errorf("reminder title cannot be empty")
		}
		cfg.Reminder.Title = value
	case "reminder.message":
		if value == "" {
			logger.Logger.Warn().Str("key", key).Msg("Empty message value")
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
		logger.Logger.Warn().Str("key", key).Msg("Unknown configuration key")
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	if err := config.SaveConfig(); err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to save configuration")
		return err
	}

	logger.Logger.Info().Str("key", key).Str("value", value).Msg("Configuration updated successfully")
	return nil
}

func validateTimeFormat(timeStr string) error {
	timeFormat := regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`) // Matches 00:00 to 23:59
	if !timeFormat.MatchString(timeStr) {
		return fmt.Errorf("invalid time format: %s (expected HH:MM)", timeStr)
	}
	return nil
}

func maskToken(token string) string {
	if len(token) > 6 {
		return token[:3] + "..." + token[len(token)-3:] // Masking in format "abc...xyz"
	}
	return "******"
}

func NewConfigCommand() *cobra.Command {
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "View or update configuration for csync.",
	}

	configCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.LoadConfig(); err != nil {
				logger.Logger.Error().Err(err).Msg("Failed to load configuration")
				fmt.Printf("❌ Error loading config: %v\n", err)
				return
			}
			ShowConfig(&config.ConfigData)
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.LoadConfig(); err != nil {
				logger.Logger.Error().Err(err).Msg("Failed to load configuration")
				fmt.Printf("❌ Error loading config: %v\n", err)
				return
			}

			key := args[0]
			value := args[1]

			if err := SetConfig(&config.ConfigData, key, value); err != nil {
				logger.Logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to set configuration")
				fmt.Printf("❌ Error: %v\n", err)
				return
			}

			fmt.Println("✅ Configuration updated successfully.")
		},
	})

	return configCmd
}
