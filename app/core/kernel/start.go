package kernel

import (
	"context"
	"fmt"
	"sort"
)

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
