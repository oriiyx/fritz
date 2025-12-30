import {DataComponent, InputSettings} from '@/generated/definitions'
import {Input} from '@/components/Input'

interface InputFieldProps {
    component: DataComponent
    value: unknown
    onChange: (value: unknown) => void
    error?: string
    disabled?: boolean
}

export function InputField({component, value, onChange, error, disabled}: InputFieldProps) {
    const settings = (component.settings || {}) as InputSettings

    return (
        <Input
            label={component.title + (component.mandatory ? ' *' : '')}
            value={String(value ?? '')}
            onChange={(e) => onChange(e.target.value)}
            error={error}
            disabled={disabled}
            fullWidth
            maxLength={settings.columnLength}
            placeholder={settings.defaultValue ? `Default: ${settings.defaultValue}` : undefined}
        />
    )
}