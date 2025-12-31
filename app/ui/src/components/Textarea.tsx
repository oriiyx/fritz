import {forwardRef, TextareaHTMLAttributes} from 'react'

export type InputVariant = 'bordered' | 'ghost'

export type InputSize = 'xs' | 'sm' | 'md' | 'lg'

export interface TextAreaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
    variant?: InputVariant
    inputSize?: InputSize
    label?: string
    error?: string
    success?: boolean
    fullWidth?: boolean
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextAreaProps>(
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
        const textareaClasses = [
            'textarea',
            variant && `textarea-${variant}`,
            inputSize !== 'md' && `input-${inputSize}`,
            error && 'textarea-error',
            success && 'textarea-success',
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
                <textarea ref={ref} className={textareaClasses} disabled={disabled} {...props} />
                {error && (
                    <label className="label">
                        <span className="label-text-alt text-error">{error}</span>
                    </label>
                )}
            </div>
        )
    }
)

Textarea.displayName = 'Textarea'
