package main

import (
	"github.com/ibexmonj/ContribSync/plugins"
	"github.com/spf13/cobra"
	"os"

	"github.com/ibexmonj/ContribSync/commands"
	"github.com/ibexmonj/ContribSync/utils"
)

func main() {
	// Initialize logger
	if err := utils.InitLogger("info"); err != nil {
		utils.Logger.Fatal().Err(err).Msg("Failed to initialize logger")
		os.Exit(1)
	}

	utils.Logger.Info().Msg("Starting csync application")

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
		utils.Logger.Fatal().Err(err).Msg("Failed to execute command")
		os.Exit(1)
	}
}
