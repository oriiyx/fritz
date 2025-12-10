package kernel

import "context"

// HookFunc is a callback function for hooks
type HookFunc func(ctx context.Context, data interface{}) error

// HookPriority defines execution order
type HookPriority int

const (
	PriorityHigh   HookPriority = 10
	PriorityNormal HookPriority = 50
	PriorityLow    HookPriority = 100
)

type hookEntry struct {
	priority HookPriority
	fn       HookFunc
}

// Hooks manages event-driven callbacks
type Hooks struct {
	hooks map[string][]hookEntry
}

func NewHooks() *Hooks {
	return &Hooks{
		hooks: make(map[string][]hookEntry),
	}
}

// Register adds a hook callback
func (h *Hooks) Register(name string, priority HookPriority, fn HookFunc) {
	entry := hookEntry{priority: priority, fn: fn}
	h.hooks[name] = append(h.hooks[name], entry)

	// Sort by priority (implement if needed)
}

// Trigger executes all callbacks for a hook
func (h *Hooks) Trigger(ctx context.Context, name string, data interface{}) error {
	entries, exists := h.hooks[name]
	if !exists {
		return nil
	}

	for _, entry := range entries {
		if err := entry.fn(ctx, data); err != nil {
			return err
		}
	}

	return nil
}
