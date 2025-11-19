package definitions

type EntityDefinition struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ParentEntityID string `json:"parentEntityID"`
	AllowInherit   bool   `json:"allowInherit"`
	Layout         Layout `json:"layout"`
}
