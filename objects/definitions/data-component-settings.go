package definitions

import (
	"fmt"
	"regexp"
	"time"
)

// SettingsValidator interface
type SettingsValidator interface {
	Validate() error
}

type InputSettings struct {
	DefaultValue    string `json:"defaultValue,omitempty"`
	ColumnLength    *int   `json:"columnLength,omitempty"`
	RegexValidation string `json:"regexValidation,omitempty"`
}

func (s InputSettings) Validate() error {
	if s.ColumnLength != nil && *s.ColumnLength <= 0 {
		return fmt.Errorf("columnLength must be greater than 0")
	}
	if s.RegexValidation != "" {
		if _, err := regexp.Compile(s.RegexValidation); err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	}
	return nil
}

type IntegerSettings struct {
	DefaultValue *int `json:"defaultValue,omitempty"`
	MinValue     *int `json:"minValue,omitempty"`
	MaxValue     *int `json:"maxValue,omitempty"`
	Unsigned     bool `json:"unsigned"`
}

func (s IntegerSettings) Validate() error {
	if s.MinValue != nil && s.MaxValue != nil && *s.MinValue > *s.MaxValue {
		return fmt.Errorf("minValue cannot be greater than maxValue")
	}
	if s.Unsigned && s.DefaultValue != nil && *s.DefaultValue < 0 {
		return fmt.Errorf("defaultValue cannot be negative when unsigned is true")
	}
	if s.Unsigned && s.MinValue != nil && *s.MinValue < 0 {
		return fmt.Errorf("minValue cannot be negative when unsigned is true")
	}
	return nil
}

type DateSettings struct {
	DefaultValue *time.Time `json:"defaultValue,omitempty"`
}

func (s DateSettings) Validate() error {
	// No specific validation needed for now
	return nil
}
