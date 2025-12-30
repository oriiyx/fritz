import {DataComponent, DateSettings} from '@/generated/definitions'
import {Input} from '@/components/Input'

interface DateFieldProps {
    component: DataComponent
    value: unknown
    onChange: (value: unknown) => void
    error?: string
    disabled?: boolean
}

export function DateField({component, value, onChange, error, disabled}: DateFieldProps) {
    const settings = (component.settings || {}) as DateSettings

    // Convert value to date string format (YYYY-MM-DD) for input
    const getDateString = (): string => {
        if (!value) return ''

        try {
            const date = new Date(String(value))
            if (isNaN(date.getTime())) return ''
            return date.toISOString().split('T')[0]
        } catch {
            return ''
        }
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const inputValue = e.target.value
        if (!inputValue) {
            onChange(null)
            return
        }

        // Convert to ISO string for storage
        const date = new Date(inputValue)
        if (!isNaN(date.getTime())) {
            onChange(date.toISOString())
        }
    }

    return (
        <Input
            type="date"
            label={component.title + (component.mandatory ? ' *' : '')}
            value={getDateString()}
            onChange={handleChange}
            error={error}
            disabled={disabled}
            fullWidth
            placeholder={settings.defaultValue ? `Default: ${settings.defaultValue}` : undefined}
        />
    )
}