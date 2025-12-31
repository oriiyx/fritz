package definition_builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
)

type TableChangeset struct {
	Added    strings.Builder
	Removed  strings.Builder
	Modified strings.Builder
}

func (e *Builder) UpdateTableFromChangeset(changeset *ComponentChangeset, tablename string, ctx context.Context) error {
	tc := e.CreateTableChangesetBasis(tablename)

	// handle adding new columns to the table
	var add []string
	for _, component := range changeset.Added {
		add = append(add, fmt.Sprintf("ADD COLUMN %s", component.ToColumnDefinition()))
	}
	if add != nil {
		tc.Added.WriteString(strings.Join(add, ", "))
	}

	// handle removing columns from the table
	var remove []string
	for _, component := range changeset.Removed {
		remove = append(remove, fmt.Sprintf("DROP COLUMN %s", component.Name))
	}
	if remove != nil {
		tc.Removed.WriteString(strings.Join(remove, ", "))
	}

	// handle modified columns from the table
	var modify []string
	for _, component := range changeset.Modified {
		// Change the type (with USING for safe conversion)
		modify = append(modify, fmt.Sprintf(
			"ALTER COLUMN %s TYPE %s USING %s::%s",
			component.Name,
			string(component.DBType),
			component.Name,
			string(component.DBType),
		))

		// Handle default value changes before nullability
		settings, _ := component.GetSettings()
		hasDefault := false
		var defaultValue string

		switch component.Type {
		case definitions.ComponentInput:
			if s, ok := settings.(definitions.InputSettings); ok && s.DefaultValue != "" {
				hasDefault = true
				defaultValue = fmt.Sprintf("'%s'", s.DefaultValue)
			}
		case definitions.ComponentTextarea:
			if s, ok := settings.(definitions.TextareaSettings); ok && s.DefaultValue != "" {
				hasDefault = true
				defaultValue = fmt.Sprintf("'%s'", s.DefaultValue)
			}
		case definitions.ComponentInteger:
			if s, ok := settings.(definitions.IntegerSettings); ok && s.DefaultValue != nil {
				hasDefault = true
				defaultValue = fmt.Sprintf("%d", *s.DefaultValue)
			}
		case definitions.ComponentFloat4:
			if s, ok := settings.(definitions.FloatSettings); ok && s.DefaultValue != nil {
				hasDefault = true
				defaultValue = fmt.Sprintf("%d", *s.DefaultValue)
			}
		case definitions.ComponentFloat8:
			if s, ok := settings.(definitions.FloatSettings); ok && s.DefaultValue != nil {
				hasDefault = true
				defaultValue = fmt.Sprintf("%d", *s.DefaultValue)
			}
		case definitions.ComponentDate:
			if s, ok := settings.(definitions.DateSettings); ok && s.DefaultValue != nil {
				hasDefault = true
				defaultValue = fmt.Sprintf("'%s'", s.DefaultValue.Format("2006-01-02"))
			}
		}

		if hasDefault {
			modify = append(modify, fmt.Sprintf("ALTER COLUMN %s SET DEFAULT %s", component.Name, defaultValue))
		} else {
			modify = append(modify, fmt.Sprintf("ALTER COLUMN %s DROP DEFAULT", component.Name))
		}

		// Handle nullability changes
		if component.Mandatory {
			modify = append(modify, fmt.Sprintf("ALTER COLUMN %s SET NOT NULL", component.Name))
		} else {
			modify = append(modify, fmt.Sprintf("ALTER COLUMN %s DROP NOT NULL", component.Name))
		}

	}
	if modify != nil {
		s := strings.Join(modify, ", ")
		tc.Modified.WriteString(s)
	}

	baseLen := len(e.generatePrefix(tablename))

	if tc.Added.Len() > baseLen {
		_, err := e.db.Exec(ctx, tc.Added.String())
		if err != nil {
			return fmt.Errorf("failed to add columns: %w", err)
		}
	}

	if tc.Removed.Len() > baseLen {
		_, err := e.db.Exec(ctx, tc.Removed.String())
		if err != nil {
			return fmt.Errorf("failed to remove columns: %w", err)
		}
	}

	if tc.Modified.Len() > baseLen {
		_, err := e.db.Exec(ctx, tc.Modified.String())
		if err != nil {
			return fmt.Errorf("failed to modify columns: %w", err)
		}
	}

	return nil
}

func (e *Builder) CreateTableChangesetBasis(tablename string) *TableChangeset {
	tc := &TableChangeset{}
	prefix := e.generatePrefix(tablename)

	tc.Added.WriteString(prefix)
	tc.Removed.WriteString(prefix)
	tc.Modified.WriteString(prefix)

	return tc
}

func (e *Builder) generatePrefix(tablename string) string {
	return fmt.Sprintf("ALTER TABLE IF EXISTS %s ", tablename)
}
