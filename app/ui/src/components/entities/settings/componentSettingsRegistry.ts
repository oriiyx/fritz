import {DataComponentType} from '@/generated/definitions'
import {ComponentSettingsProps, ComponentWithSettings} from './ComponentSettingsTypes'
import {InputSettings} from './InputSettings'
import {IntegerSettings} from './IntegerSettings'
import {DateSettings} from './DateSettings'
import * as React from 'react'

/**
 * Type helper for settings components.
 *
 * This creates a type for a React component that can handle settings for a specific component type.
 *
 * How it works:
 * 1. Takes a type parameter T which must be one of our component types ('input' | 'integer' | 'date')
 * 2. Uses TypeScript's Extract utility to pull out the specific component shape from our discriminated union
 * 3. Returns a React.FC type that accepts ComponentSettingsProps with that specific shape
 *
 * Example:
 *   SettingsComponentFor<'input'>
 *   resolves to:
 *   React.FC<ComponentSettingsProps<{ type: 'input', settings: InputSettings, ... }>>
 *
 * This allows TypeScript to know that InputSettings only receives 'input' type components,
 * IntegerSettings only receives 'integer' type components, etc.
 */
type SettingsComponentFor<T extends ComponentWithSettings['type']> = React.FC<ComponentSettingsProps<Extract<ComponentWithSettings, {
    type: T
}>>>

/**
 * Registry mapping component types to their settings components
 *
 * Each entry uses 'as SettingsComponentFor<type>' to tell TypeScript the specific
 * component type this settings component handles. This bridges the gap between
 * the specific types each component expects and the registry's need to store them all.
 */
export const componentSettingsRegistry = {
    input: InputSettings as SettingsComponentFor<'input'>,
    integer: IntegerSettings as SettingsComponentFor<'integer'>,
    date: DateSettings as SettingsComponentFor<'date'>,
} as const

/**
 * Get the settings component for a given component type
 * Returns null if no settings component is registered
 */
export function getSettingsComponent(
    componentType: DataComponentType
): React.FC<ComponentSettingsProps<ComponentWithSettings>> | null {
    const component = componentSettingsRegistry[componentType as keyof typeof componentSettingsRegistry]
    return (component as React.FC<ComponentSettingsProps<ComponentWithSettings>>) || null
}