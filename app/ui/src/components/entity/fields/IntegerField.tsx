import {DataComponent, IntegerSettings} from '@/generated/definitions'
import {Input} from '@/components/Input'

interface IntegerFieldProps {
    component: DataComponent
    value: unknown
    onChange: (value: unknown) => void
    error?: string
    disabled?: boolean
}

export function IntegerField({component, value, onChange, error, disabled}: IntegerFieldProps) {
    const settings = (component.settings || {}) as IntegerSettings

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const inputValue = e.target.value

        // Allow empty string for clearing
        if (inputValue === '') {
            onChange(null)
            return
        }

        // Parse as integer
        const numValue = parseInt(inputValue, 10)
        if (!isNaN(numValue)) {
            onChange(numValue)
        }
    }

    return (
        <Input
            type="number"
            label={component.title + (component.mandatory ? ' *' : '')}
            value={value !== null && value !== undefined ? String(value) : ''}
            onChange={handleChange}
            error={error}
            disabled={disabled}
            fullWidth
            min={settings.unsigned ? 0 : settings.minValue}
            max={settings.maxValue}
            placeholder={settings.defaultValue !== undefined ? `Default: ${settings.defaultValue}` : undefined}
        />
    )
}