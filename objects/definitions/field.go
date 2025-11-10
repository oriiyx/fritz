package definitions

import "strings"

type Field struct {
	Type        FieldType `json:"type"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	DBType      DBType    `json:"dbtype"`
	Mandatory   bool      `json:"mandatory"`
	NotEditable bool      `json:"noteditable"`
	Index       bool      `json:"index"`
	Unique      bool      `json:"unique "`
}

func (f *Field) ToColumnDefinition() string {
	parts := []string{f.Name, string(f.DBType)}

	if f.Mandatory {
		parts = append(parts, "NOT NULL")
	}
	if f.Unique {
		parts = append(parts, "UNIQUE")
	}

	return strings.Join(parts, " ")
}
