import {DataComponent, InputSettings} from '@/generated/definitions'
import {Textarea} from "@/components/Textarea.tsx";

interface TextareaFieldProps {
    component: DataComponent
    value: unknown
    onChange: (value: unknown) => void
    error?: string
    disabled?: boolean
}

export function TextareaField({component, value, onChange, error, disabled}: TextareaFieldProps) {
    const settings = (component.settings || {}) as InputSettings

    return (
        <Textarea
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