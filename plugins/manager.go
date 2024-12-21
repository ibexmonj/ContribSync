package plugins

import (
	"errors"
	"fmt"
	"github.com/ibexmonj/ContribSync/plugins/core"
	"plugin"
)

type PluginManager struct {
	plugins map[string]Plugin // Loaded plugins by name
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

// LoadCorePlugins loads built-in plugins
func (pm *PluginManager) LoadCorePlugins() {
	jira := &core.JiraPlugin{}
	pm.plugins["jira"] = jira
	err := jira.Init()
	if err != nil {
		return
	}

	// Add other built-in plugins (e.g., GitHub)
}

// LoadExternalPlugin loads a shared object (.so) plugin
func (pm *PluginManager) LoadExternalPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	symbol, err := p.Lookup("PluginInstance")
	if err != nil {
		return errors.New("plugin must export a 'PluginInstance'")
	}

	instance, ok := symbol.(Plugin)
	if !ok {
		return errors.New("invalid plugin type")
	}

	err = instance.Init()
	if err != nil {
		return err
	}
	name, _ := instance.Info()
	pm.plugins[name] = instance
	fmt.Printf("Loaded external plugin: %s\n", name)
	return nil
}

// ExecutePlugin executes a plugin by name
func (pm *PluginManager) ExecutePlugin(name string, args []string) error {
	p, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}
	return p.Execute(args)
}

// ListPlugins lists all loaded plugins
func (pm *PluginManager) ListPlugins() {
	for name, p := range pm.plugins {
		_, desc := p.Info()
		fmt.Printf("Plugin: %s - %s\n", name, desc)
	}
}
