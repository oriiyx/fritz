package definitions

import (
	"encoding/json"
	"os"
	"testing"
)

func TestUnmarshalDefinition(t *testing.T) {
	// Read the example.json file
	data, err := os.ReadFile("example.json")
	if err != nil {
		t.Fatalf("Failed to read example.json: %v", err)
	}

	// Unmarshal into Definition struct
	var def Definition
	err = json.Unmarshal(data, &def)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Validate the unmarshalled data
	if def.ID != "Product" {
		t.Errorf("Expected ID to be 'Product', got '%s'", def.ID)
	}

	if def.Name != "Product" {
		t.Errorf("Expected Name to be 'Product', got '%s'", def.Name)
	}

	if def.Description != "Product entity" {
		t.Errorf("Expected Description to be 'Product entity', got '%s'", def.Description)
	}

	if !def.AllowInherit {
		t.Error("Expected AllowInherit to be true")
	}

	if def.Layout.Type != "panel" {
		t.Errorf("Expected Layout.Type to be 'panel', got '%s'", def.Layout.Type)
	}

	if len(def.Layout.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(def.Layout.Fields))
	}

	// Validate first field (name)
	if len(def.Layout.Fields) > 0 {
		field := def.Layout.Fields[0]
		if field.Name != "name" {
			t.Errorf("Expected first field name to be 'name', got '%s'", field.Name)
		}
		if field.Title != "Product Name" {
			t.Errorf("Expected first field title to be 'Product Name', got '%s'", field.Title)
		}
		if field.Type != TypeInput {
			t.Errorf("Expected first field type to be TypeInput, got '%s'", field.Type)
		}
		if !field.Mandatory {
			t.Error("Expected first field to be mandatory")
		}
	}

	t.Logf("Successfully unmarshaled Definition: %+v", def)
}

func TestGenerateDefinition(t *testing.T) {
	var fields []Field
	nameField := Field{
		Type:        TypeInput,
		Name:        "name",
		Title:       "Product Name",
		DBType:      DataTypeVarchar.WithSize(255),
		Mandatory:   true,
		NotEditable: false,
		Index:       true,
		Unique:      false,
	}

	priceField := Field{
		Type:      TypeNumeric,
		Name:      "price",
		Title:     "Price",
		DBType:    DataTypeDecimal.WithPrecision(10, 2),
		Mandatory: true,
	}
	fields = append(fields, priceField, nameField)

	def := Definition{
		ID:           "Product",
		Name:         "Product",
		Description:  "Product entity",
		AllowInherit: true,
		Layout: Layout{
			Type:   "panel",
			Fields: fields,
		},
	}

	jsonData, _ := json.MarshalIndent(def, "", "  ")
	println(string(jsonData))

	t.Logf("Successfully unmarshaled Definition: %+v", def)
}
