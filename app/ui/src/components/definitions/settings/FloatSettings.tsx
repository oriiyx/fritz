import {Input} from '@/components/Input'
import {ComponentSettingsProps, ComponentWithSettings} from './ComponentSettingsTypes'

export function FloatSettings({
                                  component,
                                  onSettingsChange
                              }: ComponentSettingsProps<ComponentWithSettings & { type: 'float4' | 'float8' }>) {
    const settings = component.settings

    return (
        <div className="space-y-4">
            <Input
                label="Default Value"
                type="number"
                value={settings.defaultValue?.toString() || ''}
                onChange={(e) =>
                    onSettingsChange(
                        'defaultValue',
                        e.target.value ? parseInt(e.target.value) : undefined
                    )
                }
                fullWidth
            />
            <Input
                label="Minimum Value"
                type="number"
                value={settings.minValue?.toString() || ''}
                onChange={(e) =>
                    onSettingsChange(
                        'minValue',
                        e.target.value ? parseInt(e.target.value) : undefined
                    )
                }
                fullWidth
            />
            <Input
                label="Maximum Value"
                type="number"
                value={settings.maxValue?.toString() || ''}
                onChange={(e) =>
                    onSettingsChange(
                        'maxValue',
                        e.target.value ? parseInt(e.target.value) : undefined
                    )
                }
                fullWidth
            />
        </div>
    )
}