import * as React from 'react'
import {useState} from 'react'
import {Bars3Icon, TrashIcon} from '@heroicons/react/24/outline'
import {DataComponent} from '@/generated/definitions'

interface Props {
    components: DataComponent[]
    selectedComponent: DataComponent | null
    onComponentsReorder: (newComponents: DataComponent[]) => void
    onComponentSelect: (component: DataComponent) => void
    onComponentDelete: (component: DataComponent) => void
}

export function ComponentLayoutTree({
                                        components,
                                        selectedComponent,
                                        onComponentsReorder,
                                        onComponentSelect,
                                        onComponentDelete,
                                    }: Props) {
    const [draggedIndex, setDraggedIndex] = useState<number | null>(null)
    const [dropIndicator, setDropIndicator] = useState<{
        index: number
        position: 'above' | 'below'
    } | null>(null)

    const handleDragStart = (e: React.DragEvent, index: number) => {
        setDraggedIndex(index)
        e.dataTransfer.effectAllowed = 'move'
    }

    const handleDragOver = (e: React.DragEvent, index: number) => {
        e.preventDefault()

        if (draggedIndex === null || draggedIndex === index) return

        const rect = e.currentTarget.getBoundingClientRect()
        const midpoint = rect.top + rect.height / 2
        const position = e.clientY < midpoint ? 'above' : 'below'

        setDropIndicator({index, position})
    }

    const handleDrop = (e: React.DragEvent, dropIndex: number) => {
        e.preventDefault()

        if (draggedIndex === null || !dropIndicator) return

        const newComponents = [...components]
        const draggedComponent = newComponents[draggedIndex]

        // Remove dragged component
        newComponents.splice(draggedIndex, 1)

        // Calculate new position
        let insertIndex = dropIndex
        if (dropIndicator.position === 'below') {
            insertIndex++
        }
        if (draggedIndex < dropIndex) {
            insertIndex--
        }

        // Insert at new position
        newComponents.splice(insertIndex, 0, draggedComponent)

        onComponentsReorder(newComponents)
        setDraggedIndex(null)
        setDropIndicator(null)
    }

    const handleDragEnd = () => {
        setDraggedIndex(null)
        setDropIndicator(null)
    }

    const handleClick = (component: DataComponent) => {
        onComponentSelect(component)
    }

    const handleDelete = (e: React.MouseEvent, component: DataComponent) => {
        e.stopPropagation() // Prevent selecting when deleting
        onComponentDelete(component)
    }

    if (components.length === 0) {
        return (
            <div
                className="flex h-48 items-center justify-center rounded-lg border-2 border-dashed border-base-300 bg-base-100">
                <div className="text-center text-base-content/50">
                    <p className="text-sm">No components yet</p>
                    <p className="mt-1 text-xs">Add your first component to get started</p>
                </div>
            </div>
        )
    }

    return (
        <div className="space-y-2">
            {components.map((component, index) => (
                <div key={component.id}>
                    {/* Drop indicator above */}
                    {dropIndicator?.index === index && dropIndicator.position === 'above' && (
                        <div className="my-1 h-0.5 bg-primary"/>
                    )}

                    <div
                        onDragOver={(e) => handleDragOver(e, index)}
                        onDrop={(e) => handleDrop(e, index)}
                        onClick={() => handleClick(component)}
                        className={`
                            group flex cursor-pointer items-center gap-3 rounded-lg border-2 p-3 transition-all
                            ${
                            selectedComponent?.id === component.id
                                ? 'border-primary bg-primary/10'
                                : 'border-base-300 hover:border-base-400'
                        }
                            ${draggedIndex === index ? 'opacity-50' : ''}
                        `}
                    >
                        {/* Drag Handle - ONLY this is draggable */}
                        <div
                            draggable
                            onDragStart={(e) => handleDragStart(e, index)}
                            onDragEnd={handleDragEnd}
                            className="cursor-grab text-base-content/40 transition-colors hover:text-primary active:cursor-grabbing"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <Bars3Icon className="h-5 w-5"/>
                        </div>

                        {/* Component Info */}
                        <div className="min-w-0 flex-1">
                            <div className="flex items-center gap-1">
                                <span className="badge badge-primary badge-xs">{component.type}</span>
                                <span title={component.title}
                                      className="truncate text-xs font-medium">{component.title}</span>
                            </div>
                            <span className="text-xs font-mono text-base-content/60">
                                {component.name}
                            </span>
                        </div>

                        {/* Delete Button - Shows on hover */}
                        <button
                            onClick={(e) => handleDelete(e, component)}
                            className="opacity-0 transition-opacity group-hover:opacity-100 btn btn-ghost btn-sm btn-circle"
                            title="Delete component"
                        >
                            <TrashIcon className="h-4 w-4 text-error"/>
                        </button>
                    </div>

                    {/* Drop indicator below */}
                    {dropIndicator?.index === index && dropIndicator.position === 'below' && (
                        <div className="my-1 h-0.5 bg-primary"/>
                    )}
                </div>
            ))}
        </div>
    )
}