package definition_builder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
	"github.com/oriiyx/fritz/app/core/utils/helpers/slug"
	"github.com/oriiyx/fritz/app/core/utils/rw"
	"github.com/rs/zerolog"
)

const EntitiesDefinitionsFilePathTemplate = "var/entities/definitions"

type Builder struct {
	db     *pgxpool.Pool
	cw     *rw.CustomWriter
	logger *zerolog.Logger
}

func NewDefinitionsBuilder(logger *zerolog.Logger, db *pgxpool.Pool, cw *rw.CustomWriter) *Builder {
	return &Builder{
		db:     db,
		cw:     cw,
		logger: logger,
	}
}

// ValidateNewDefinition validates everything that new definition has to adhere to
//
// [ ] - Duplicate ID
//
// [ ] - Duplicate Name
//
// [ ] - Duplicate Component Names
func (e *Builder) ValidateNewDefinition(definition *definitions.EntityDefinition) ([]byte, error) {
	existingDefinitions, err := e.LoadDefinitionsFromEntityFiles()
	if err != nil {
		return nil, err
	}

	// Check for duplicate ID
	for _, existing := range existingDefinitions {
		if existing.ID == definition.ID {
			return []byte(`{"error": "entity id already exists", "conflictingId": "` + definition.ID + `"}`), nil
		}
	}

	// Check for duplicate Name
	for _, existing := range existingDefinitions {
		if existing.Name == definition.Name {
			return []byte(`{"error": "entity name already exists", "conflictingName": "` + definition.Name + `"}`), nil
		}
	}

	validationOfExistingDefinitionsResult, err := e.ValidateExistingDefinition(definition)
	if err != nil {
		return nil, err
	}

	return validationOfExistingDefinitionsResult, nil
}

// ValidateExistingDefinition validates definition points that touch only core of the definition
//
// [ ] - Duplicate Component Names
func (e *Builder) ValidateExistingDefinition(definition *definitions.EntityDefinition) ([]byte, error) {
	// Check for duplicated component names
	componentNames := make(map[string]bool, 1)
	for _, component := range definition.Layout.Components {
		if componentNames[component.Name] {
			return []byte(`{"error": "entity duplicated component name", "conflictingComponentName": "` + component.Name + `"}`), nil
		}

		componentNames[component.Name] = true
	}

	return nil, nil
}

func (e *Builder) StoreDefinitionIntoEntityFile(definition *definitions.EntityDefinition) error {
	filename := e.CreateEntityDefinitionFileName(definition)

	// First create directory
	if err := os.MkdirAll(EntitiesDefinitionsFilePathTemplate, 0755); err != nil {
		return fmt.Errorf("failed to create var/entities/definitions directory: %w", err)
	}

	// Second create the .json file
	filePath := filepath.Join(EntitiesDefinitionsFilePathTemplate, filename)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create definition entity file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(definition)
	if err != nil {
		return err
	}

	return nil
}

func (e *Builder) CreateEntityDefinitionFileName(definition *definitions.EntityDefinition) string {
	return fmt.Sprintf("entity_%s.json", slug.CreateSlug(definition.ID))
}

func (e *Builder) LoadDefinitionsFromEntityFiles() ([]*definitions.EntityDefinition, error) {
	// load all the entity files that are stored
	entries, err := os.ReadDir(EntitiesDefinitionsFilePathTemplate)
	if err != nil {
		return nil, err
	}

	var entities []*definitions.EntityDefinition

	// loop over all the entries and prepare the
	for _, e := range entries {
		filePath := filepath.Join(EntitiesDefinitionsFilePathTemplate, e.Name())
		file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to create definition entity file: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)

		var entity definitions.EntityDefinition
		err = decoder.Decode(&entity)
		if err != nil {
			return nil, err
		}

		entities = append(entities, &entity)
	}

	return entities, nil
}

// LoadDefinitionByID loads a specific entity definition by its ID
func (e *Builder) LoadDefinitionByID(id string) (*definitions.EntityDefinition, error) {
	slugName := slug.CreateSlug(id)
	filename := fmt.Sprintf("entity_%s.json", slugName)
	filePath := filepath.Join(EntitiesDefinitionsFilePathTemplate, filename)

	var definition definitions.EntityDefinition
	if err := e.cw.ReadJSONFromFile(filePath, &definition); err != nil {
		return nil, fmt.Errorf("failed to load definition for '%s': %w", id, err)
	}

	return &definition, nil
}
