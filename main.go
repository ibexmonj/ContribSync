package main

import (
	"fmt"
	"os"

	"github.com/ibexmonj/ContribSync/commands"
	"github.com/ibexmonj/ContribSync/config"
	"github.com/ibexmonj/ContribSync/utils"
)

func main() {
	// Initialize logger with configurable log level
	logLevel := os.Getenv("LOG_LEVEL") // Use an environment variable for log level
	if logLevel == "" {
		logLevel = "info" // Default to info level
	}

	if err := utils.InitLogger(logLevel); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	utils.Logger.Info().Msg("Starting csync application")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to load configuration")
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Validate configuration
	if err := config.ValidateConfig(cfg); err != nil {
		utils.Logger.Error().Err(err).Msg("Invalid configuration")
		fmt.Printf("Invalid configuration: %v\n", err)
		return
	}

	// Check if a command is provided
	if len(os.Args) < 2 {
		commands.Help()
		return
	}

	// Parse the command
	command := os.Args[1]

	// Handle commands
	switch command {
	case "help":
		utils.Logger.Info().Msg("Executing 'help' command")
		commands.Help()
	case "config":
		utils.Logger.Info().Msg("Executing 'config' command")
		if len(os.Args) > 2 && os.Args[2] == "set" {
			if len(os.Args) < 5 {
				fmt.Println("Usage: csync config set <key> <value>")
				return
			}
			key := os.Args[3]
			value := os.Args[4]
			if err := commands.SetConfig(cfg, key, value); err != nil {
				utils.Logger.Error().Err(err).Msg("Error setting config")
				fmt.Printf("Error: %v\n", err)
			} else {
				utils.Logger.Info().Str("key", key).Str("value", value).Msg("Config updated")
				fmt.Println("Configuration updated successfully.")
			}
		} else {
			commands.ShowConfig(cfg)
		}
	case "reminder":
		utils.Logger.Info().Msg("Executing 'reminder' command")
		args := os.Args[2:] // Pass subcommands like "test"
		commands.ReminderCommand(cfg, args)
	default:
		utils.Logger.Warn().Str("command", command).Msg("Unknown command")
		fmt.Printf("Unknown command: %s\n", command)
		commands.Help()
	}
}
