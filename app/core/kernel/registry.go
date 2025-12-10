package kernel

import (
	"fmt"
	"sync"
)

// Registry provides service container for dependency injection
type Registry struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]interface{}),
	}
}

// Register adds a service to the registry
func (r *Registry) Register(name string, service interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.services[name]; exists {
		return fmt.Errorf("service '%s' already registered", name)
	}

	r.services[name] = service
	return nil
}

// Get retrieves a service from registry
func (r *Registry) Get(name string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.services[name]
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	return service, nil
}

// MustGet retrieves service or panics
func (r *Registry) MustGet(name string) interface{} {
	service, err := r.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}
