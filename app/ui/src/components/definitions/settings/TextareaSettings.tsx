import {Input} from '@/components/Input'
import {ComponentSettingsProps, ComponentWithSettings} from './ComponentSettingsTypes'

export function TextareaSettings({
                                  component,
                                  onSettingsChange
                              }: ComponentSettingsProps<ComponentWithSettings & { type: 'textarea' }>) {
    const settings = component.settings

    return (
        <div className="space-y-4">
            <Input
                label="Default Value"
                value={settings.defaultValue || ''}
                onChange={(e) => onSettingsChange('defaultValue', e.target.value)}
                fullWidth
            />
        </div>
    )
}