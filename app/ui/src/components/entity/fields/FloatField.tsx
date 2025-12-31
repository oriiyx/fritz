import {DataComponent, FloatSettings} from '@/generated/definitions'
import {Input} from '@/components/Input'
import * as React from "react";

interface FloatFieldProps {
    component: DataComponent
    value: unknown
    onChange: (value: unknown) => void
    error?: string
    disabled?: boolean
}

export function FloatField({component, value, onChange, error, disabled}: FloatFieldProps) {
    const settings = (component.settings || {}) as FloatSettings

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const inputValue = e.target.value

        // Allow empty string for clearing
        if (inputValue === '') {
            onChange(null)
            return
        }

        // Parse as float
        const floatValue = parseFloat(inputValue)
        if (!isNaN(floatValue)) {
            onChange(floatValue)
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
            min={settings.minValue}
            max={settings.maxValue}
            placeholder={settings.defaultValue !== undefined ? `Default: ${settings.defaultValue}` : undefined}
        />
    )
}