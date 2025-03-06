package main

import (
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"github.com/ibexmonj/ContribSync/pkg/plugins"
	"github.com/spf13/cobra"
	"os"

	"github.com/ibexmonj/ContribSync/commands"
)

func main() {
	if err := logger.InitLogger("info"); err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to initialize logger")
		os.Exit(1)
	}

	logger.Logger.Info().Msg("Starting csync... ")

	var rootCmd = &cobra.Command{
		Use:   "csync",
		Short: "Csync - Contribution Sync CLI",
		Long:  "Csync helps manage and log your contributions across various platforms.",
	}

	rootCmd.AddCommand(commands.NewConfigCommand())
	rootCmd.AddCommand(commands.NewReminderCommand())
	rootCmd.AddCommand(commands.NewCompletionCommand())

	pluginManager := plugins.NewPluginManager()
	pluginManager.LoadCorePlugins()

	rootCmd.AddCommand(commands.NewPluginCommand(pluginManager))

	if err := rootCmd.Execute(); err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to execute command")
		os.Exit(1)
	}
}
