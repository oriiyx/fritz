package kernel

import (
	"context"
	"fmt"
	"sort"
)

type Kernel struct {
	registry *Registry
	plugins  *Plugins
	hooks    *Hooks
}

func New() *Kernel {
	return &Kernel{
		registry: NewRegistry(),
		plugins:  NewPlugins(),
		hooks:    NewHooks(),
	}
}

// Registry returns the service registry
func (k *Kernel) Registry() *Registry {
	return k.registry
}

// Hooks returns the hook manager
func (k *Kernel) Hooks() *Hooks {
	return k.hooks
}

// RegisterPlugin adds a plugin to the kernel
func (k *Kernel) RegisterPlugin(plugin Plugin) {
	k.plugins.Register(plugin)
}

// Start initializes kernel and all plugins
func (k *Kernel) Start(ctx context.Context) error {
	// Sort plugins by priority
	sort.Slice(k.plugins.registered, func(i, j int) bool {
		return k.plugins.registered[i].Priority() < k.plugins.registered[j].Priority()
	})

	// Initialize each plugin
	for _, plugin := range k.plugins.registered {
		if err := plugin.Initialize(ctx, k); err != nil {
			return fmt.Errorf("failed to initialize plugin '%s': %w", plugin.Name(), err)
		}

		k.plugins.initialized[plugin.Name()] = plugin
	}

	return nil
}

// Shutdown gracefully stops all plugins
func (k *Kernel) Shutdown(ctx context.Context) error {
	// Shutdown in reverse order
	for i := len(k.plugins.registered) - 1; i >= 0; i-- {
		plugin := k.plugins.registered[i]
		if err := plugin.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown plugin '%s': %w", plugin.Name(), err)
		}
	}

	return nil
}
