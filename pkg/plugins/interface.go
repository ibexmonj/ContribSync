package plugins

type Plugin interface {
	Init() error                             // Initialize the plugin
	Execute(args []string) error             // Execute plugin-specific commands
	Info() (name string, description string) // Return plugin name and description
}
