package definitions

type EntityDefinition struct {
	ID           string `json:"id" validate:"required,max=255"`
	Name         string `json:"name" validate:"required,max=255"`
	Description  string `json:"description" validate:"max=1000"`
	AllowInherit bool   `json:"allowInherit"`
	Layout       Layout `json:"layout"`
}
