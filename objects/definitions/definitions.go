package definitions

type Definitions struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ParentClassID string `json:"parentClassID"`
	AllowInherit  bool   `json:"allowInherit"`
}
