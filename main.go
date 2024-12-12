package main

import (
	"fmt"
	"os"

	"github.com/ibexmonj/ContribSync/commands"
	"github.com/ibexmonj/ContribSync/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
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
		commands.Help()
	case "config":
		if len(os.Args) > 2 && os.Args[2] == "set" {
			if len(os.Args) < 5 {
				fmt.Println("Usage: csync config set <key> <value>")
				return
			}
			key := os.Args[3]
			value := os.Args[4]
			if err := commands.SetConfig(cfg, key, value); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Configuration updated successfully.")
			}
		} else {
			commands.ShowConfig(cfg)
		}
	case "reminder":
		commands.ReminderCommand(cfg)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		commands.Help()
	}
}
