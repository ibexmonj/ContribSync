package main

import (
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"github.com/ibexmonj/ContribSync/pkg/plugins"
	"github.com/spf13/cobra"
	"os"

	"github.com/ibexmonj/ContribSync/commands"
)

func main() {
	// Initialize logger
	if err := logger.InitLogger("info"); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to initialize logger")
		os.Exit(1)
	}

	logger.Logger.Info().Msg("Starting csync application")

	// Root command
	var rootCmd = &cobra.Command{
		Use:   "csync",
		Short: "Csync - Contribution Sync CLI",
		Long:  "Csync helps manage and log your contributions across various platforms.",
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewConfigCommand())
	rootCmd.AddCommand(commands.NewReminderCommand())
	//	rootCmd.AddCommand(commands.NewHelpCommand())
	rootCmd.AddCommand(commands.NewCompletionCommand())

	// Initialize the PluginManager
	pluginManager := plugins.NewPluginManager()
	pluginManager.LoadCorePlugins()

	// Add plugin commands
	rootCmd.AddCommand(commands.NewPluginCommand(pluginManager))

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to execute command")
		os.Exit(1)
	}
}
