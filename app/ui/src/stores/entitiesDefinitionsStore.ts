import {create} from 'zustand'
import {persist} from 'zustand/middleware'

export type ComponentType = 'input' | 'integer' | 'date'

export interface InputSettings {
    defaultValue?: string
    columnLength?: number
    regexValidation?: string
}

export interface IntegerSettings {
    defaultValue?: number
    minValue?: number
    maxValue?: number
    unsigned?: boolean
}

export interface DateSettings {
    defaultValue?: string
}

export type ComponentSettings = InputSettings | IntegerSettings | DateSettings

export interface DefinitionComponent {
    type: ComponentType
    name: string
    title: string
    dbtype: string
    mandatory: boolean
    invisible: boolean
    notEditable: boolean
    settings: ComponentSettings
}

export interface DefinitionLayout {
    type: string
    components: DefinitionComponent[]
}

export interface CreateNewEntityDefinition {
    id: string
    name: string
}

export interface EntityDefinition {
    id: string
    name: string
    description: string
    allowInherit: boolean
    layout: DefinitionLayout
}

interface EntityDefinitionState {
    entityDefinitions: EntityDefinition[]
    selectedDefinition: EntityDefinition | null
    isLoading: boolean
    setEntityDefinitions: (definitions: EntityDefinition[]) => void
    setSelectedDefinition: (definition: EntityDefinition | null) => void
    addEntityDefinition: (definition: CreateNewEntityDefinition) => void
    updateEntityDefinition: (id: string, definition: EntityDefinition) => void
    removeEntityDefinition: (id: string) => void
    setLoading: (loading: boolean) => void
}

export const useEntityDefinitionsStore = create<EntityDefinitionState>()(
    persist(
        (set) => ({
            entityDefinitions: [],
            selectedDefinition: null,
            isLoading: false,
            setEntityDefinitions: (definitions) =>
                set({entityDefinitions: definitions, isLoading: false}),
            setSelectedDefinition: (definition) =>
                set({selectedDefinition: definition}),
            addEntityDefinition: (definition) =>
                set((state) => ({
                    entityDefinitions: [...state.entityDefinitions, definition],
                })),
            updateEntityDefinition: (id, definition) =>
                set((state) => ({
                    entityDefinitions: state.entityDefinitions.map((def) =>
                        def.id === id ? definition : def
                    ),
                    selectedDefinition:
                        state.selectedDefinition?.id === id
                            ? definition
                            : state.selectedDefinition,
                })),
            removeEntityDefinition: (id) =>
                set((state) => ({
                    entityDefinitions: state.entityDefinitions.filter(
                        (def) => def.id !== id
                    ),
                    selectedDefinition:
                        state.selectedDefinition?.id === id
                            ? null
                            : state.selectedDefinition,
                })),
            setLoading: (loading) => set({isLoading: loading}),
        }),
        {
            name: 'fritz-entities-definitions',
        }
    )
)