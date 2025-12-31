package definitions

type DataComponentType string
type DataComponentCategory string

const (
	CategoryText    DataComponentCategory = "text"
	CategoryNumeric DataComponentCategory = "numeric"
	CategoryDate    DataComponentCategory = "date"
)

const (
	ComponentInput    DataComponentType = "input"
	ComponentTextarea DataComponentType = "textarea"
	ComponentInteger  DataComponentType = "integer"
	ComponentFloat4   DataComponentType = "float4"
	ComponentFloat8   DataComponentType = "float8"
	ComponentDate     DataComponentType = "date"
)

// DataComponentDefinition Common metadata for all data components
type DataComponentDefinition struct {
	ID            DataComponentType     `json:"id"`
	Label         string                `json:"label"`
	Category      DataComponentCategory `json:"category"`
	Tooltip       string                `json:"tooltip"`
	Icon          string                `json:"icon,omitempty"`
	DefaultDBType DBType                `json:"defaultDBType"`
}

var DataComponentRegistry = map[DataComponentType]DataComponentDefinition{
	ComponentInput: {
		ID:            ComponentInput,
		Label:         "Input",
		Category:      CategoryText,
		Tooltip:       "Single line text input field",
		Icon:          "text-cursor",
		DefaultDBType: DataTypeVarchar,
	},
	ComponentTextarea: {
		ID:            ComponentTextarea,
		Label:         "Textarea",
		Category:      CategoryText,
		Tooltip:       "Multi line text input field",
		Icon:          "text-cursor",
		DefaultDBType: DataTypeText,
	},
	ComponentInteger: {
		ID:            ComponentInteger,
		Label:         "Integer",
		Category:      CategoryNumeric,
		Tooltip:       "Whole number field",
		Icon:          "hash",
		DefaultDBType: DataTypeInteger,
	},
	ComponentFloat4: {
		ID:            ComponentFloat4,
		Label:         "Float 4-byte",
		Category:      CategoryNumeric,
		Tooltip:       "Float 4-byte number field",
		Icon:          "hash",
		DefaultDBType: DataTypeFloat4,
	},
	ComponentFloat8: {
		ID:            ComponentFloat8,
		Label:         "Float 8-byte",
		Category:      CategoryNumeric,
		Tooltip:       "Float 8-byte number field",
		Icon:          "hash",
		DefaultDBType: DataTypeFloat8,
	},
	ComponentDate: {
		ID:            ComponentDate,
		Label:         "Date",
		Category:      CategoryDate,
		Tooltip:       "Date picker field",
		Icon:          "calendar",
		DefaultDBType: DataTypeDate,
	},
}

func GetDataComponentDefinition(ct DataComponentType) (DataComponentDefinition, bool) {
	def, ok := DataComponentRegistry[ct]

	return def, ok
}

func GetDataComponentsByCategory(cat DataComponentCategory) []DataComponentDefinition {
	var result []DataComponentDefinition
	for _, def := range DataComponentRegistry {
		if def.Category == cat {
			result = append(result, def)
		}
	}

	return result
}

func GetAllDataComponents() []DataComponentDefinition {
	result := make([]DataComponentDefinition, 0, len(DataComponentRegistry))
	for _, def := range DataComponentRegistry {
		result = append(result, def)
	}

	return result
}
