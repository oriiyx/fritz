import {ButtonHTMLAttributes, ReactNode} from 'react'

export type ButtonVariant =
    | 'primary'
    | 'secondary'
    | 'accent'
    | 'ghost'
    | 'link'
    | 'info'
    | 'success'
    | 'warning'
    | 'error'

export type ButtonSize = 'xs' | 'sm' | 'md' | 'lg'

export type ButtonShape = 'circle' | 'square'

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: ButtonVariant
    outline?: boolean
    size?: ButtonSize
    shape?: ButtonShape
    wide?: boolean
    block?: boolean
    glass?: boolean
    loading?: boolean
    active?: boolean
    iconLeft?: ReactNode
    iconRight?: ReactNode
    children?: ReactNode
}

export function Button({
                           variant,
                           outline = false,
                           size = 'md',
                           shape,
                           wide = false,
                           block = false,
                           glass = false,
                           loading = false,
                           active = false,
                           iconLeft,
                           iconRight,
                           disabled,
                           className = '',
                           children,
                           ...props
                       }: ButtonProps) {
    const classes = [
        'btn',
        variant && `btn-${variant}`,
        outline && 'btn-outline',
        size !== 'md' && `btn-${size}`,
        shape && `btn-${shape}`,
        wide && 'btn-wide',
        block && 'btn-block',
        glass && 'glass',
        active && 'btn-active',
        loading && 'loading',
        className,
    ]
        .filter(Boolean)
        .join(' ')

    return (
        <button className={classes} disabled={disabled || loading} {...props}>
            {loading ? (
                <span className="loading loading-spinner"></span>
            ) : (
                <>
                    {iconLeft && <span className="flex items-center">{iconLeft}</span>}
                    {children}
                    {iconRight && <span className="flex items-center">{iconRight}</span>}
                </>
            )}
        </button>
    )
}
