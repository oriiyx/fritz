package entity_builder

import (
	"reflect"

	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
)

type ComponentChangeset struct {
	Added     []definitions.DataComponent
	Removed   []definitions.DataComponent
	Modified  []definitions.DataComponent
	Unchanged []definitions.DataComponent
}

func (e *EntityBuilder) CompareDefinitions(existing, new *definitions.EntityDefinition) (*ComponentChangeset, error) {
	changeset := ComponentChangeset{
		Added:     make([]definitions.DataComponent, 0),
		Removed:   make([]definitions.DataComponent, 0),
		Modified:  make([]definitions.DataComponent, 0),
		Unchanged: make([]definitions.DataComponent, 0),
	}

	listOfNewComponents := make(map[string]definitions.DataComponent)
	listOfExistingComponents := make(map[string]definitions.DataComponent)

	for _, component := range new.Layout.Components {
		listOfNewComponents[component.Name] = component
	}

	for _, component := range existing.Layout.Components {
		listOfExistingComponents[component.Name] = component
	}

	for id, newComponent := range listOfNewComponents {
		existingComponent, exists := listOfExistingComponents[id]
		if !exists {
			// it should be added
			changeset.Added = append(changeset.Added, newComponent)
		} else {
			// it could be modified or unchanged
			e.CompareComponents(&existingComponent, &newComponent, &changeset)
		}
	}

	for id, existingComponent := range listOfExistingComponents {
		_, exists := listOfNewComponents[id]
		if !exists {
			changeset.Removed = append(changeset.Removed, existingComponent)
		}
	}

	return &changeset, nil
}

func (e *EntityBuilder) CompareComponents(existing, new *definitions.DataComponent, changeset *ComponentChangeset) {
	// handle unchanged
	areEqual := reflect.DeepEqual(new, existing)
	if areEqual {
		changeset.Unchanged = append(changeset.Unchanged, *new)
		return
	}

	changeset.Modified = append(changeset.Modified, *new)
}
