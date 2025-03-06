package plugins

import (
	"fmt"
	"github.com/ibexmonj/ContribSync/pkg/logger"
	"plugin"
)

type PluginManager struct {
	Plugins map[string]Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		Plugins: make(map[string]Plugin),
	}
}

func (pm *PluginManager) RegisterPlugin(plugin Plugin) {
	name, desc := plugin.Info()
	pm.Plugins[name] = plugin
	logger.Logger.Info().
		Str("plugin", name).
		Msgf("✅ Loaded plugin: %s - %s", name, desc)
}

// LoadCorePlugins initializes built-in plugins
func (pm *PluginManager) LoadCorePlugins() {
	logger.Logger.Info().Msg("🔍 Loading core plugins...")

	pm.RegisterPlugin(&GitHubPlugin{})
	pm.RegisterPlugin(&JiraPlugin{})
	pm.RegisterPlugin(&SlackPlugin{})

	logger.Logger.Info().
		Int("plugin_count", len(pm.Plugins)).
		Msg("✅ Core plugins loaded successfully")
}

func (pm *PluginManager) LoadExternalPlugin(path string) error {
	logger.Logger.Info().Str("plugin_path", path).Msg("🔗 Loading external plugin")

	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	sym, err := p.Lookup("PluginInstance")
	if err != nil {
		return fmt.Errorf("failed to find PluginInstance in %s: %w", path, err)
	}

	pluginInstance, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("invalid plugin format in %s", path)
	}

	name, desc := pluginInstance.Info()
	pm.Plugins[name] = pluginInstance

	logger.Logger.Info().Str("plugin", name).Msgf("✅ Loaded external plugin: %s - %s", name, desc)
	return nil
}

func (pm *PluginManager) ExecutePlugin(name string, args []string) error {
	plugin, exists := pm.Plugins[name]
	if !exists {
		logger.Logger.Error().Str("plugin", name).Msg("❌ Plugin not found")
		return fmt.Errorf("plugin not found: %s", name)
	}

	logger.Logger.Info().Str("plugin", name).Msg("🚀 Executing plugin")
	return plugin.Execute(args)
}

func (pm *PluginManager) ListPlugins() {
	fmt.Println("\n🔌 Loaded Plugins:")
	for name, p := range pm.Plugins {
		_, desc := p.Info()
		fmt.Printf("   ✅ %s - %s\n", name, desc)
	}
}
