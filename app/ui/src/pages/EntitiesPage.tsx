import {useEffect, useState} from 'react'
import {useQuery} from '@tanstack/react-query'
import {ChevronDownIcon, ChevronRightIcon} from '@heroicons/react/24/outline'
import {Card, CardBody, CardTitle} from '@/components/Card'
import {Alert} from '@/components/Alert'
import {Loading} from '@/components/Loading'
import {entitiesApi} from '@/services/entitiesService'
import type {DefinitionComponent, EntityDefinition,} from '@/stores/entitiesDefinitionsStore'
import {EntityList} from "@/components/entities/EntityList.tsx";

export function EntitiesPage() {
    // Local state management
    const [selectedDefinition, setSelectedDefinition] = useState<EntityDefinition | null>(null)

    // Fetch entity definitions using React Query
    const {data: entityDefinitions = [], isLoading, error} = useQuery({
        queryKey: ['entity-definitions'],
        queryFn: entitiesApi.getEntityDefinitions,
    })

    // Fetch entity definitions using React Query
    const {
        data: dataComponentTypes = [],
        isLoading: isDataComponentLoading,
        error: dataComponentError
    } = useQuery({
        queryKey: ['data-component-types'],
        queryFn: entitiesApi.getDataComponentTypes,
    })

    const handleSelectDefinition = (definition: EntityDefinition) => {
        setSelectedDefinition(definition)
    }

    useEffect(() => {
        console.log('Data Component Types', dataComponentTypes)
    }, [dataComponentTypes]);

    if (isLoading || isDataComponentLoading) {
        return (
            <div className="flex h-full items-center justify-center">
                <Loading variant="spinner" size="lg"/>
            </div>
        )
    }

    if (error || dataComponentError) {
        return (
            <div className="space-y-6">
                <Alert variant="error">
                    Failed to load entity definitions. Please try again.
                </Alert>
            </div>
        )
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold">Entity Definitions</h1>
                    <p className="mt-2 text-base-content/70">
                        Manage your entity schemas and structures
                    </p>
                </div>
            </div>

            {/* Split View - 1/4 for list, 3/4 for details */}
            <div className="grid grid-cols-1 gap-6 lg:grid-cols-4">
                {/* Left Panel - Entity List (1/4) */}
                <div className="lg:col-span-1">
                    <EntityList
                        entityDefinitions={entityDefinitions}
                        selectedDefinition={selectedDefinition}
                        handleSelectDefinition={handleSelectDefinition}
                    />
                </div>

                {/* Right Panel - Entity Details (3/4) */}
                <div className="lg:col-span-3">
                    <Card>
                        <CardBody>
                            <CardTitle>Definition Details</CardTitle>
                            {selectedDefinition ? (
                                <div className="space-y-6">
                                    {/* Basic Info */}
                                    <div className="space-y-3">
                                        <div>
                                            <label className="label">
                                                <span className="label-text font-semibold">Name</span>
                                            </label>
                                            <div className="text-lg">{selectedDefinition.name}</div>
                                        </div>
                                        <div>
                                            <label className="label">
                                                <span className="label-text font-semibold">ID</span>
                                            </label>
                                            <div className="font-mono text-sm text-base-content/70">
                                                {selectedDefinition.id}
                                            </div>
                                        </div>
                                        {selectedDefinition.description && (
                                            <div>
                                                <label className="label">
                                                    <span className="label-text font-semibold">Description</span>
                                                </label>
                                                <div className="text-base-content/70">
                                                    {selectedDefinition.description}
                                                </div>
                                            </div>
                                        )}
                                    </div>

                                    <div className="divider"></div>

                                    {/* Layout Tree */}
                                    {selectedDefinition.layout?.components?.length > 0 && (
                                        <div>
                                            <label className="label">
                                                <span className="label-text font-semibold">Layout Structure</span>
                                            </label>
                                            <div className="rounded-lg border border-base-300 bg-base-200 p-4">
                                                <LayoutTree definition={selectedDefinition}/>
                                            </div>
                                        </div>
                                    )}
                                </div>
                            ) : (
                                <div className="py-12 text-center text-base-content/50">
                                    <p>Select a definition from the list to view details</p>
                                </div>
                            )}
                        </CardBody>
                    </Card>
                </div>
            </div>
        </div>
    )
}

// Tree component for displaying layout structure
function LayoutTree({definition}: { definition: EntityDefinition }) {
    const [expandedComponents, setExpandedComponents] = useState<Set<string>>(new Set())

    const toggleComponent = (componentName: string) => {
        setExpandedComponents((prev) => {
            const next = new Set(prev)
            if (next.has(componentName)) {
                next.delete(componentName)
            } else {
                next.add(componentName)
            }
            return next
        })
    }

    return (
        <div className="space-y-2">
            {/* Root Layout Node */}
            <div className="flex items-center gap-2 font-mono text-sm">
                <span className="badge badge-sm badge-neutral">{definition.layout.type}</span>
                <span className="text-base-content/60">Layout Root</span>
            </div>

            {/* Components */}
            {definition.layout.components.length === 0 ? (
                <div className="ml-8 py-4 text-sm text-base-content/50">No components defined</div>
            ) : (
                <div className="ml-4 space-y-2">
                    {definition.layout.components.map((component) => (
                        <ComponentNode
                            key={component.name}
                            component={component}
                            isExpanded={expandedComponents.has(component.name)}
                            onToggle={() => toggleComponent(component.name)}
                        />
                    ))}
                </div>
            )}
        </div>
    )
}

// Individual component node in tree
function ComponentNode({
                           component,
                           isExpanded,
                           onToggle,
                       }: {
    component: DefinitionComponent
    isExpanded: boolean
    onToggle: () => void
}) {
    return (
        <div className="rounded border border-base-300 bg-base-100">
            <button
                onClick={onToggle}
                className="flex w-full items-center gap-2 p-3 text-left hover:bg-base-200"
            >
                {isExpanded ? (
                    <ChevronDownIcon className="h-4 w-4 flex-shrink-0"/>
                ) : (
                    <ChevronRightIcon className="h-4 w-4 flex-shrink-0"/>
                )}
                <div className="flex flex-1 items-center gap-2">
                    <span className="badge badge-sm badge-primary">{component.type}</span>
                    <span className="font-medium">{component.title}</span>
                    <span className="font-mono text-xs text-base-content/50">
                        ({component.name})
                    </span>
                </div>
            </button>

            {isExpanded && (
                <div className="border-t border-base-300 bg-base-200/50 p-3 text-sm">
                    <div className="grid gap-2">
                        <div className="flex justify-between">
                            <span className="text-base-content/60">Database Type:</span>
                            <span className="font-mono">{component.dbtype}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-base-content/60">Mandatory:</span>
                            <span>{component.mandatory ? 'Yes' : 'No'}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-base-content/60">Visible:</span>
                            <span>{component.invisible ? 'No' : 'Yes'}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-base-content/60">Editable:</span>
                            <span>{component.notEditable ? 'No' : 'Yes'}</span>
                        </div>

                        {/* Settings */}
                        {Object.keys(component.settings).length > 0 && (
                            <>
                                <div className="divider my-1"></div>
                                <div className="text-xs font-semibold text-base-content/60">Settings:</div>
                                {Object.entries(component.settings).map(([key, value]) => (
                                    <div key={key} className="flex justify-between">
                                        <span className="text-base-content/60">{key}:</span>
                                        <span className="font-mono text-xs">
                                            {typeof value === 'boolean'
                                                ? value
                                                    ? 'true'
                                                    : 'false'
                                                : value?.toString() || 'null'}
                                        </span>
                                    </div>
                                ))}
                            </>
                        )}
                    </div>
                </div>
            )}
        </div>
    )
}