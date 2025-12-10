package kernel

import "context"

// Plugin defines the contract all plugins must implement
type Plugin interface {
	// Name returns unique plugin identifier
	Name() string

	// Priority determines initialization order (lower = earlier)
	Priority() int

	// Initialize is called when kernel starts
	Initialize(ctx context.Context, k *Kernel) error

	// Shutdown is called when kernel stops
	Shutdown(ctx context.Context) error
}

// Plugins manages plugin collection
type Plugins struct {
	registered  []Plugin
	initialized map[string]Plugin
}

func NewPlugins() *Plugins {
	return &Plugins{
		registered:  make([]Plugin, 0),
		initialized: make(map[string]Plugin),
	}
}

func (p *Plugins) Register(plugin Plugin) {
	p.registered = append(p.registered, plugin)
}

func (p *Plugins) Get(name string) (Plugin, bool) {
	plugin, exists := p.initialized[name]
	return plugin, exists
}
