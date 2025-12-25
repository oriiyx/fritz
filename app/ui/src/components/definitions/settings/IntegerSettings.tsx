import {Input} from '@/components/Input'
import {ComponentSettingsProps, ComponentWithSettings} from './ComponentSettingsTypes'

export function IntegerSettings({
                                    component,
                                    onSettingsChange
                                }: ComponentSettingsProps<ComponentWithSettings & { type: 'integer' }>) {
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
            <div className="form-control fieldset">
                <label className="label cursor-pointer justify-start gap-3">
                    <input
                        type="checkbox"
                        checked={settings.unsigned}
                        onChange={(e) => onSettingsChange('unsigned', e.target.checked)}
                        className="checkbox checkbox-primary"
                    />
                    <span className="label-text">Unsigned (positive numbers only)</span>
                </label>
            </div>
        </div>
    )
}