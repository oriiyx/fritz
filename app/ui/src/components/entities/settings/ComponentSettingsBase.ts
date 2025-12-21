import {DataComponent} from '@/generated/definitions'

/**
 * Base props interface that all component-specific settings components should implement
 */
export interface ComponentSettingsProps {
    component: DataComponent
    onSettingsChange: (settingKey: string, value: any) => void
}

/**
 * Base interface for component settings components
 * Used for type safety in the registry
 */
export type ComponentSettingsComponent = React.FC<ComponentSettingsProps>
