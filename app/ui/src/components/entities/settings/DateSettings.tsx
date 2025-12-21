import {Input} from '@/components/Input'
import {ComponentSettingsProps, ComponentWithSettings} from './ComponentSettingsTypes'

export function DateSettings({
                                 component,
                                 onSettingsChange
                             }: ComponentSettingsProps<ComponentWithSettings & { type: 'date' }>) {
    const settings = component.settings

    return (
        <div className="space-y-4">
            <Input
                label="Default Value"
                type="date"
                value={settings.defaultValue || ''}
                onChange={(e) => onSettingsChange('defaultValue', e.target.value)}
                fullWidth
            />
        </div>
    )
}