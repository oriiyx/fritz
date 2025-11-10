package definitions

type Definition struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ParentClassID string `json:"parentClassID"`
	AllowInherit  bool   `json:"allowInherit"`
	Layout        Layout `json:"layout"`
}
