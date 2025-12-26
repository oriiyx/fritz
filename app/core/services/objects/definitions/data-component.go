package definitions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DataComponent represents an actual configured data component instance
type DataComponent struct {
	Type   DataComponentType `json:"type" validate:"required"`
	Name   string            `json:"name" validate:"required,max=255"`
	Title  string            `json:"title" validate:"required,max=255"`
	DBType DBType            `json:"dbtype" validate:"required"`

	// Common properties
	Mandatory   bool `json:"mandatory"`
	Invisible   bool `json:"invisible"`
	NotEditable bool `json:"notEditable"`

	// Type-specific settings stored as raw JSON
	Settings json.RawMessage `json:"settings"`
}

// GetDefinition retrieves the type definition for this component
func (dc *DataComponent) GetDefinition() (DataComponentDefinition, bool) {
	return GetDataComponentDefinition(dc.Type)
}

// GetSettings unmarshalls settings into the correct type based on component type
func (dc *DataComponent) GetSettings() (interface{}, error) {
	if dc.Settings == nil {
		return nil, nil
	}

	switch dc.Type {
	case ComponentInput:
		var settings InputSettings
		if err := json.Unmarshal(dc.Settings, &settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal input settings: %w", err)
		}
		return settings, nil

	case ComponentInteger:
		var settings IntegerSettings
		if err := json.Unmarshal(dc.Settings, &settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal integer settings: %w", err)
		}
		return settings, nil

	case ComponentDate:
		var settings DateSettings
		if err := json.Unmarshal(dc.Settings, &settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal date settings: %w", err)
		}
		return settings, nil

	default:
		return nil, fmt.Errorf("unknown component type: %s", dc.Type)
	}
}

// ValidateSettings validates the settings for this component
func (dc *DataComponent) ValidateSettings() error {
	settings, err := dc.GetSettings()
	if err != nil {
		return err
	}

	if settings == nil {
		return nil
	}

	if validator, ok := settings.(SettingsValidator); ok {
		return validator.Validate()
	}

	return nil
}

// ToColumnDefinition generates SQL add column definition to table
func (dc *DataComponent) ToColumnDefinition() string {
	parts := []string{dc.Name, string(dc.DBType)}

	// Get settings to check for default value
	settings, _ := dc.GetSettings()

	// Handle default values based on type
	switch dc.Type {
	case ComponentInput:
		if s, ok := settings.(InputSettings); ok && s.DefaultValue != "" {
			parts = append(parts, fmt.Sprintf("DEFAULT '%s'", s.DefaultValue))
		}
	case ComponentInteger:
		if s, ok := settings.(IntegerSettings); ok && s.DefaultValue != nil {
			parts = append(parts, fmt.Sprintf("DEFAULT %d", *s.DefaultValue))
		}
	case ComponentDate:
		if s, ok := settings.(DateSettings); ok && s.DefaultValue != nil {
			parts = append(parts, fmt.Sprintf("DEFAULT '%s'", s.DefaultValue.Format("2006-01-02")))
		}
	}

	if dc.Mandatory {
		parts = append(parts, "NOT NULL")
	}

	return strings.Join(parts, " ")
}

// GetGoType returns the Go type that SQLC generates based on DB type and nullability
// when using pgx/v5 driver (which is what Fritz uses)
func (dc *DataComponent) GetGoType() string {
	switch dc.DBType {
	case DataTypeVarchar, DataTypeText, DataTypeChar:
		if dc.Mandatory {
			return "string"
		}
		return "pgtype.Text"

	case DataTypeInteger:
		if dc.Mandatory {
			return "int32"
		}
		return "pgtype.Int4"

	case DataTypeBigInt:
		if dc.Mandatory {
			return "int64"
		}
		return "pgtype.Int8"

	case DataTypeSmallInt:
		if dc.Mandatory {
			return "int16"
		}
		return "pgtype.Int2"

	case DataTypeBoolean:
		if dc.Mandatory {
			return "bool"
		}
		return "pgtype.Bool"

	case DataTypeDate:
		// Date is always pgtype.Date regardless of nullability
		return "pgtype.Date"

	case DataTypeTimestamp, DataTypeTimestampTZ:
		// Timestamp types are always pgtype.Timestamptz
		return "pgtype.Timestamptz"

	default:
		return "string"
	}
}

// GetGoTypeImport returns the import needed for this type (if any)
func (dc *DataComponent) GetGoTypeImport() string {
	goType := dc.GetGoType()
	if strings.HasPrefix(goType, "pgtype.") {
		return "github.com/jackc/pgx/v5/pgtype"
	}
	return ""
}
