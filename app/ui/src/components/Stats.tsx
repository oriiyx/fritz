import {ReactNode} from 'react'

export interface StatsProps {
    children: ReactNode
    vertical?: boolean
    horizontal?: boolean
    className?: string
}

export interface StatProps {
    title: string
    value: string | number
    desc?: string
    figure?: ReactNode
    valueClassName?: string
    className?: string
}

export function Stats({children, vertical = false, horizontal = false, className = ''}: StatsProps) {
    const classes = [
        'stats',
        'shadow',
        vertical && 'stats-vertical',
        horizontal && 'stats-horizontal',
        className,
    ]
        .filter(Boolean)
        .join(' ')

    return <div className={classes}>{children}</div>
}

export function Stat({title, value, desc, figure, valueClassName = '', className = ''}: StatProps) {
    return (
        <div className={`stat ${className}`}>
            {figure && <div className="stat-figure">{figure}</div>}
            <div className="stat-title">{title}</div>
            <div className={`stat-value ${valueClassName}`}>{value}</div>
            {desc && <div className="stat-desc">{desc}</div>}
        </div>
    )
}
