import {ReactNode} from 'react'

export type BadgeVariant =
    | 'neutral'
    | 'primary'
    | 'secondary'
    | 'accent'
    | 'ghost'
    | 'info'
    | 'success'
    | 'warning'
    | 'error'

export type BadgeSize = 'xs' | 'sm' | 'md' | 'lg'

export interface BadgeProps {
    children: ReactNode
    variant?: BadgeVariant
    size?: BadgeSize
    outline?: boolean
    className?: string
}

export function Badge({
                          children,
                          variant = 'neutral',
                          size = 'md',
                          outline = false,
                          className = '',
                      }: BadgeProps) {
    const classes = [
        'badge',
        variant && `badge-${variant}`,
        size !== 'md' && `badge-${size}`,
        outline && 'badge-outline',
        className,
    ]
        .filter(Boolean)
        .join(' ')

    return <span className={classes}>{children}</span>
}
