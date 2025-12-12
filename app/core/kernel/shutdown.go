package kernel

import (
	"context"
	"fmt"
)

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
