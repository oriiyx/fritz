import {forwardRef, InputHTMLAttributes} from 'react'

export type InputVariant = 'bordered' | 'ghost'

export type InputSize = 'xs' | 'sm' | 'md' | 'lg'

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    variant?: InputVariant
    inputSize?: InputSize
    label?: string
    error?: string
    success?: boolean
    fullWidth?: boolean
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
    (
        {
            variant = 'bordered',
            inputSize = 'md',
            label,
            error,
            success = false,
            fullWidth = false,
            className = '',
            disabled,
            ...props
        },
        ref
    ) => {
        const inputClasses = [
            'input',
            variant && `input-${variant}`,
            inputSize !== 'md' && `input-${inputSize}`,
            error && 'input-error',
            success && 'input-success',
            fullWidth && 'w-full',
            className,
        ]
            .filter(Boolean)
            .join(' ')

        return (
            <div className={`form-control fieldset ${fullWidth ? 'w-full' : ''}`}>
                {label && (
                    <label className="label">
                        <span className="label-text">{label}</span>
                    </label>
                )}
                <input ref={ref} className={inputClasses} disabled={disabled} {...props} />
                {error && (
                    <label className="label">
                        <span className="label-text-alt text-error">{error}</span>
                    </label>
                )}
            </div>
        )
    }
)

Input.displayName = 'Input'
