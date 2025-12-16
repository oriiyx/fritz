import {ImgHTMLAttributes, ReactNode} from 'react'

export type AvatarSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl'

export interface AvatarProps extends ImgHTMLAttributes<HTMLImageElement> {
    size?: AvatarSize
    rounded?: boolean
    ring?: boolean
    ringColor?: string
    online?: boolean
    offline?: boolean
    placeholder?: ReactNode
    fallback?: string
}

const sizeMap: Record<AvatarSize, string> = {
    xs: 'w-8',
    sm: 'w-12',
    md: 'w-16',
    lg: 'w-24',
    xl: 'w-32',
}

export function Avatar({
                           size = 'md',
                           rounded = true,
                           ring = false,
                           ringColor = '',
                           online = false,
                           offline = false,
                           placeholder,
                           fallback,
                           className = '',
                           src,
                           alt = '',
                           ...props
                       }: AvatarProps) {
    const avatarClasses = ['avatar', online && 'online', offline && 'offline', ring && 'ring', ringColor]
        .filter(Boolean)
        .join(' ')

    const imgClasses = [sizeMap[size], rounded && 'rounded-full', className].filter(Boolean).join(' ')

    // If we have a placeholder (like an icon), use that
    if (placeholder) {
        return (
            <div className={avatarClasses}>
                <div className={imgClasses}>{placeholder}</div>
            </div>
        )
    }

    // If we have a source image, use that
    if (src) {
        return (
            <div className={avatarClasses}>
                <div className={imgClasses}>
                    <img src={src} alt={alt} {...props} />
                </div>
            </div>
        )
    }

    // Otherwise show fallback text (initials)
    if (fallback) {
        return (
            <div className={avatarClasses}>
                <div className={`${imgClasses} flex items-center justify-center bg-neutral text-neutral-content`}>
                    <span className="text-xl">{fallback}</span>
                </div>
            </div>
        )
    }

    return null
}
