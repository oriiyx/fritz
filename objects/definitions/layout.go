package definitions

type LayoutType string

type Layout struct {
	Type     LayoutType `json:"type"`
	Children Field      `json:"children"`
}
