import {DataComponentType} from '@/generated/definitions'
import {ComponentSettingsComponent} from './ComponentSettingsBase'
import {InputSettings} from './InputSettings'
import {IntegerSettings} from './IntegerSettings'
import {DateSettings} from './DateSettings'

/**
 * Registry mapping component types to their settings components
 * Add new component types here as they are created
 */
export const componentSettingsRegistry: Record<DataComponentType, ComponentSettingsComponent> = {
    input: InputSettings,
    integer: IntegerSettings,
    date: DateSettings,
}

/**
 * Get the settings component for a given component type
 * Returns null if no settings component is registered
 */
export function getSettingsComponent(
    componentType: DataComponentType
): ComponentSettingsComponent | null {
    return componentSettingsRegistry[componentType] || null
}
