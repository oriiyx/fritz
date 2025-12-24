package entity_builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
)

const entitiesTableSchemaFilePathTemplate = "database/schema"

func (e *EntityBuilder) CreateEntityTable(ctx context.Context, definition *definitions.EntityDefinition) (string, error) {
	tableName := e.CreateEntityTableName(definition)

	// Build CREATE TABLE statement from definition.Layout.Components
	columns := []string{
		"id UUID PRIMARY KEY DEFAULT uuid_generate_v4()",
		"entity_id UUID NOT NULL REFERENCES entities(id) ON DELETE CASCADE",
		"created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()",
		"updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()",
	}

	for _, component := range definition.Layout.Components {
		columns = append(columns, component.ToColumnDefinition())
	}

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)",
		tableName,
		strings.Join(columns, ", "))

	// Execute SQL
	_, err := e.db.Exec(ctx, sql)

	if err != nil {
		return "", err
	}

	err = e.cw.WriteNewFile(sql, entitiesTableSchemaFilePathTemplate, fmt.Sprintf("%s.sql", tableName))
	if err != nil {
		return "", err
	}

	return tableName, nil
}

func (e *EntityBuilder) CreateEntityTableName(definition *definitions.EntityDefinition) string {
	tableName := fmt.Sprintf("entity_%s", definition.ID)
	return tableName
}
