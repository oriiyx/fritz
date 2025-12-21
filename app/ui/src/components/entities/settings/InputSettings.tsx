import {Input} from '@/components/Input'
import {ComponentSettingsProps, ComponentWithSettings} from './ComponentSettingsTypes'

export function InputSettings({
                                  component,
                                  onSettingsChange
                              }: ComponentSettingsProps<ComponentWithSettings & { type: 'input' }>) {
    const settings = component.settings

    return (
        <div className="space-y-4">
            <Input
                label="Default Value"
                value={settings.defaultValue || ''}
                onChange={(e) => onSettingsChange('defaultValue', e.target.value)}
                fullWidth
            />
            <Input
                label="Column Length"
                type="number"
                value={settings.columnLength?.toString() || ''}
                onChange={(e) =>
                    onSettingsChange(
                        'columnLength',
                        e.target.value ? parseInt(e.target.value) : undefined
                    )
                }
                fullWidth
            />
            <Input
                label="Regex Validation"
                value={settings.regexValidation || ''}
                onChange={(e) => onSettingsChange('regexValidation', e.target.value)}
                fullWidth
                className="font-mono text-sm"
            />
        </div>
    )
}