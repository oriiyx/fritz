import {ReactNode} from 'react'

export interface CardProps {
    children: ReactNode
    bordered?: boolean
    imageFull?: boolean
    compact?: boolean
    side?: boolean
    className?: string
}

export interface CardTitleProps {
    children: ReactNode
    className?: string
}

export interface CardBodyProps {
    children: ReactNode
    className?: string
}

export interface CardActionsProps {
    children: ReactNode
    justify?: 'start' | 'center' | 'end'
    className?: string
}

export function Card({
                         children,
                         bordered = false,
                         imageFull = false,
                         compact = false,
                         side = false,
                         className = '',
                     }: CardProps) {
    const classes = [
        'card',
        'bg-base-100',
        '',
        bordered && 'card-bordered',
        imageFull && 'image-full',
        compact && 'card-compact',
        side && 'card-side',
        className,
    ]
        .filter(Boolean)
        .join(' ')

    return <div className={classes}>{children}</div>
}

export function CardTitle({children, className = ''}: CardTitleProps) {
    return <h2 className={`card-title ${className}`}>{children}</h2>
}

export function CardBody({children, className = ''}: CardBodyProps) {
    return <div className={`card-body ${className}`}>{children}</div>
}

export function CardActions({children, justify = 'end', className = ''}: CardActionsProps) {
    const classes = ['card-actions', justify !== 'end' && `justify-${justify}`, className]
        .filter(Boolean)
        .join(' ')

    return <div className={classes}>{children}</div>
}
