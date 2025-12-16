export type LoadingVariant = 'spinner' | 'dots' | 'ring' | 'ball' | 'bars' | 'infinity'

export type LoadingSize = 'xs' | 'sm' | 'md' | 'lg'

export interface LoadingProps {
    variant?: LoadingVariant
    size?: LoadingSize
    className?: string
}

export function Loading({variant = 'spinner', size = 'md', className = ''}: LoadingProps) {
    const classes = [
        'loading',
        `loading-${variant}`,
        size !== 'md' && `loading-${size}`,
        className,
    ]
        .filter(Boolean)
        .join(' ')

    return <span className={classes}></span>
}

export interface LoadingOverlayProps {
    message?: string
    variant?: LoadingVariant
    size?: LoadingSize
}

export function LoadingOverlay({message, variant = 'spinner', size = 'lg'}: LoadingOverlayProps) {
    return (
        <div className="flex h-screen items-center justify-center">
            <div className="flex flex-col items-center gap-4">
                <Loading variant={variant} size={size}/>
                {message && <p className="text-base-content/70">{message}</p>}
            </div>
        </div>
    )
}
