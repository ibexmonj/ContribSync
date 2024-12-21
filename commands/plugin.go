package commands

import (
	"fmt"
	"github.com/ibexmonj/ContribSync/plugins"

	"github.com/spf13/cobra"
)

func NewPluginCommand(pm *plugins.PluginManager) *cobra.Command {
	pluginCmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
	}

	pluginCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all loaded plugins",
		Run: func(cmd *cobra.Command, args []string) {
			pm.ListPlugins()
		},
	})

	pluginCmd.AddCommand(&cobra.Command{
		Use:   "load [path]",
		Short: "Load an external plugin",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			if err := pm.LoadExternalPlugin(path); err != nil {
				fmt.Printf("Failed to load plugin: %v\n", err)
			}
		},
	})

	pluginCmd.AddCommand(&cobra.Command{
		Use:   "exec [name] [args...]",
		Short: "Execute a plugin by name",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if err := pm.ExecutePlugin(name, args[1:]); err != nil {
				fmt.Printf("Failed to execute plugin: %v\n", err)
			}
		},
	})

	return pluginCmd
}
