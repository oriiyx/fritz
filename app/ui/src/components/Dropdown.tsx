import {ReactNode, useState} from 'react'

export type DropdownPosition = 'top' | 'bottom' | 'left' | 'right' | 'end'

export interface DropdownProps {
    trigger: ReactNode
    children: ReactNode
    position?: DropdownPosition
    hover?: boolean
    className?: string
}

export function Dropdown({
                             trigger,
                             children,
                             position = 'bottom',
                             hover = false,
                             className = '',
                         }: DropdownProps) {
    const [isOpen, setIsOpen] = useState(false)

    const dropdownClasses = [
        'dropdown',
        position && `dropdown-${position}`,
        hover && 'dropdown-hover',
        isOpen && 'dropdown-open',
        className,
    ]
        .filter(Boolean)
        .join(' ')

    return (
        <div className={dropdownClasses}>
            <div
                tabIndex={0}
                role="button"
                onClick={() => !hover && setIsOpen(!isOpen)}
                onBlur={() => !hover && setIsOpen(false)}
            >
                {trigger}
            </div>
            {(isOpen || hover) && (
                <div tabIndex={0} className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
                    {children}
                </div>
            )}
        </div>
    )
}
