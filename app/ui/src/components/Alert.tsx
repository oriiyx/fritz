import {ReactNode} from 'react'

export type AlertVariant = 'info' | 'success' | 'warning' | 'error'

export interface AlertProps {
    children: ReactNode
    variant?: AlertVariant
    icon?: ReactNode
    className?: string
}

export function Alert({children, variant = 'info', icon, className = ''}: AlertProps) {
    const classes = ['alert', variant && `alert-${variant}`, className].filter(Boolean).join(' ')

    return (
        <div className={classes} role="alert">
            {icon && <span>{icon}</span>}
            <span>{children}</span>
        </div>
    )
}
