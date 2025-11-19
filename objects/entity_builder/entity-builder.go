package entity_builder

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oriiyx/fritz/objects/definitions"
	"github.com/oriiyx/fritz/utils/helpers/slug"
)

const entitiesDefinitionsFilePathTemplate = "var/entities/definitions/%s"

type EntityBuilder struct {
}

func (e *EntityBuilder) StoreDefinitionIntoEntityFile(definition *definitions.EntityDefinition) error {
	slugEntityName := slug.CreateSlug(definition.Name)
	filename := fmt.Sprintf("entity_%s.json", slugEntityName)
	filePath := fmt.Sprintf(entitiesDefinitionsFilePathTemplate, filename)

	file, _ := os.OpenFile(filePath, os.O_CREATE, os.ModePerm)
	defer file.Close()

	encoder := json.NewEncoder(file)
	err := encoder.Encode(definition)
	if err != nil {
		return err
	}

	return nil
}

func (e *EntityBuilder) LoadDefinitionsFromEntityFiles() {

}
