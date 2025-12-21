import {useState} from 'react'
import {PlusIcon, ChevronDownIcon} from '@heroicons/react/24/outline'
import {DataComponentDefinition, DataComponentCategory} from '@/generated/definitions'

interface Props {
    dataComponentTypes: DataComponentDefinition[]
    onAddComponent: (componentType: DataComponentDefinition) => void
}

export function AddComponentDropdown({dataComponentTypes, onAddComponent}: Props) {
    const [isOpen, setIsOpen] = useState(false)

    // Group components by category
    const groupedComponents = dataComponentTypes.reduce(
        (acc, component) => {
            const category = component.category
            if (!acc[category]) {
                acc[category] = []
            }
            acc[category].push(component)
            return acc
        },
        {} as Record<DataComponentCategory, DataComponentDefinition[]>
    )

    const handleSelect = (component: DataComponentDefinition) => {
        onAddComponent(component)
        setIsOpen(false)
    }

    const getCategoryLabel = (category: DataComponentCategory): string => {
        const labels: Record<DataComponentCategory, string> = {
            text: 'Text',
            numeric: 'Numeric',
            date: 'Date',
        }
        return labels[category]
    }

    const getCategoryIcon = (category: DataComponentCategory): string => {
        const icons: Record<DataComponentCategory, string> = {
            text: '📝',
            numeric: '#️⃣',
            date: '📅',
        }
        return icons[category]
    }

    return (
        <div className="dropdown dropdown-bottom w-full">
            <button
                tabIndex={0}
                onClick={() => setIsOpen(!isOpen)}
                className="btn btn-primary btn-block gap-2"
            >
                <PlusIcon className="h-5 w-5"/>
                Add Component
                <ChevronDownIcon className={`h-4 w-4 transition-transform ${isOpen ? 'rotate-180' : ''}`}/>
            </button>

            {isOpen && (
                <ul
                    tabIndex={0}
                    className="menu dropdown-content z-[1] mt-2 w-full rounded-box bg-base-100 p-2 shadow-lg"
                >
                    {Object.entries(groupedComponents).map(([category, components]) => (
                        <li key={category}>
                            {/* Category Header */}
                            <div className="menu-title flex items-center gap-2">
                                <span>{getCategoryIcon(category as DataComponentCategory)}</span>
                                <span>{getCategoryLabel(category as DataComponentCategory)}</span>
                            </div>

                            {/* Components in Category */}
                            <ul>
                                {components.map((component) => (
                                    <li key={component.id}>
                                        <button
                                            onClick={() => handleSelect(component)}
                                            className="flex items-start gap-3"
                                        >
                                            <div className="flex-1">
                                                <div className="font-medium">{component.label}</div>
                                                <div className="text-xs text-base-content/60">
                                                    {component.tooltip}
                                                </div>
                                            </div>
                                        </button>
                                    </li>
                                ))}
                            </ul>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    )
}