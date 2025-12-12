package kernel

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
