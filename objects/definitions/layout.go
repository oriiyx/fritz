package definitions

type LayoutType string

type Layout struct {
	Type   LayoutType `json:"type"`
	Fields []Field    `json:"fields"`
}
