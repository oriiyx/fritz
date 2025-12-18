// Type definitions for entity components and definitions
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