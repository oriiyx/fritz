import {
    DataComponent,
    DateSettings as DateSettingsType, FloatSettings,
    InputSettings as InputSettingsType,
    IntegerSettings as IntegerSettingsType, TextareaSettings
} from '@/generated/definitions'

/**
 * Discriminated union for all component settings
 * This makes TypeScript know exactly which settings shape to expect based on the component type
 */
export type ComponentWithSettings =
    | (Omit<DataComponent, 'type' | 'settings'> & {
    type: 'input'
    settings: InputSettingsType
})
    | (Omit<DataComponent, 'type' | 'settings'> & {
    type: 'textarea'
    settings: TextareaSettings
})
     | (Omit<DataComponent, 'type' | 'settings'> & {
    type: 'float4'
    settings: FloatSettings
})
    | (Omit<DataComponent, 'type' | 'settings'> & {
    type: 'float8'
    settings: FloatSettings
})
    | (Omit<DataComponent, 'type' | 'settings'> & {
    type: 'integer'
    settings: IntegerSettingsType
})
    | (Omit<DataComponent, 'type' | 'settings'> & {
    type: 'date'
    settings: DateSettingsType
})

/**
 * Base props interface that all component-specific settings components should implement
 */
export interface ComponentSettingsProps<T extends ComponentWithSettings> {
    component: T
    onSettingsChange: (settingKey: string, value: unknown) => void
}

/**
 * Helper to narrow component type
 */
export function isInputComponent(component: DataComponent): component is ComponentWithSettings & { type: 'input' } {
    return component.type === 'input'
}

export function isIntegerComponent(component: DataComponent): component is ComponentWithSettings & { type: 'integer' } {
    return component.type === 'integer'
}

export function isFloat4Component(component: DataComponent): component is ComponentWithSettings & { type: 'float4' } {
    return component.type === 'float4'
}
export function isFloat8Component(component: DataComponent): component is ComponentWithSettings & { type: 'float8' } {
    return component.type === 'float8'
}

export function isDateComponent(component: DataComponent): component is ComponentWithSettings & { type: 'date' } {
    return component.type === 'date'
}