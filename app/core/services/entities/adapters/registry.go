package adapters

import (
	"context"
	"fmt"
	"sync"
)

// EntityAdapter is the interface all generated adapters implement
type EntityAdapter interface {
	Create(ctx context.Context, entityID any, data map[string]interface{}) (interface{}, error)
	Read(ctx context.Context, id any) (interface{}, error)
	Update(ctx context.Context, id any, data map[string]interface{}) (interface{}, error)
	Delete(ctx context.Context, id any) error
}

var (
	registry = make(map[string]EntityAdapter)
	mu       sync.RWMutex
)

// Register registers an adapter for an entity class
func Register(classID string, adapter EntityAdapter) {
	mu.Lock()
	defer mu.Unlock()
	registry[classID] = adapter
}

// Get retrieves an adapter for an entity class
func Get(classID string) (EntityAdapter, error) {
	mu.RLock()
	defer mu.RUnlock()

	adapter, exists := registry[classID]
	if !exists {
		return nil, fmt.Errorf("no adapter registered for entity class: %s", classID)
	}
	return adapter, nil
}
