package definitions

type LayoutType string

type Layout struct {
	Type       LayoutType      `json:"type"`
	Components []DataComponent `json:"components"`
}
