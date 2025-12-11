package hooks

// Entity Hook Constants
const (
	HookBeforeEntityCreate = "entity.before_create"
	HookAfterEntityCreate  = "entity.after_create"
	HookBeforeEntityUpdate = "entity.before_update"
	HookAfterEntityUpdate  = "entity.after_update"
	HookBeforeEntityDelete = "entity.before_delete"
	HookAfterEntityDelete  = "entity.after_delete"
)

// Controller Hook Constants
const (
	HookBeforeRouterMiddleware = "controller.before_middleware"
	HookAfterRouterMiddleware  = "controller.after_middleware"
)
