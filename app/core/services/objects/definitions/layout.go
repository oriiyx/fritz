package definitions

type LayoutType string

type Layout struct {
	Type       LayoutType      `json:"type" validate:"required"`
	Components []DataComponent `json:"components" validate:"required,dive"`
}
