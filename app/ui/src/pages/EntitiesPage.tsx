import {useState} from 'react'
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {ChevronDownIcon, ChevronRightIcon, CubeIcon, PlusIcon} from '@heroicons/react/24/outline'
import {Card, CardBody, CardTitle} from '@/components/Card'
import {Button} from '@/components/Button'
import {Input} from '@/components/Input'
import {Alert} from '@/components/Alert'
import {Loading} from '@/components/Loading'
import {entitiesApi} from '@/services/entitiesService'
import {DefinitionComponent, EntityDefinition, useEntityDefinitionsStore,} from '@/stores/entitiesDefinitionsStore'

export function EntitiesPage() {
    const queryClient = useQueryClient()
    const {
        entityDefinitions,
        selectedDefinition,
        setEntityDefinitions,
        setSelectedDefinition,
        addEntityDefinition,
    } = useEntityDefinitionsStore()

    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
    const [newEntityId, setNewEntityId] = useState('')
    const [newEntityName, setNewEntityName] = useState('')
    const [createError, setCreateError] = useState<string | null>(null)

    // Fetch entity definitions
    const {isLoading, error} = useQuery({
        queryKey: ['entity-definitions'],
        queryFn: async () => {
            const data = await entitiesApi.getEntityDefinitions()
            setEntityDefinitions(data)
            return data
        },
    })

    // Create entity definition mutation
    const createMutation = useMutation({
        mutationFn: entitiesApi.createEntityDefinition,
        onSuccess: () => {
            const data = entitiesApi.getEntityDefinitions()
            setEntityDefinitions(data)
            queryClient.invalidateQueries({queryKey: ['entity-definitions']})
            setIsCreateModalOpen(false)
            setNewEntityId('')
            setNewEntityName('')
            setCreateError(null)
        },
        onError: (error: any) => {
            const errorMessage =
                error?.response?.data?.error || 'Failed to create entity definition'
            setCreateError(errorMessage)
        },
    })

    const handleCreateEntity = () => {
        setCreateError(null)

        if (!newEntityId.trim() || !newEntityName.trim()) {
            setCreateError('Both ID and Name are required')
            return
        }

        // Validate ID format (lowercase, no spaces, alphanumeric + underscore/dash)
        if (!/^[a-z0-9_-]+$/.test(newEntityId)) {
            setCreateError('ID must be lowercase letters, numbers, underscores, or dashes only')
            return
        }

        createMutation.mutate({
            id: newEntityId,
            name: newEntityName,
        })
    }

    const handleSelectDefinition = (definition: EntityDefinition) => {
        setSelectedDefinition(definition)
    }

    if (isLoading) {
        return (
            <div className="flex h-full items-center justify-center">
                <Loading variant="spinner" size="lg"/>
            </div>
        )
    }

    if (error) {
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
                <Button
                    variant="primary"
                    iconLeft={<PlusIcon className="h-5 w-5"/>}
                    onClick={() => setIsCreateModalOpen(true)}
                >
                    Create New
                </Button>
            </div>

            {/* Split View */}
            <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
                {/* Left Panel - Entity List */}
                <Card>
                    <CardBody>
                        <CardTitle>Available Definitions</CardTitle>
                        {entityDefinitions.length === 0 ? (
                            <div className="py-12 text-center text-base-content/50">
                                <CubeIcon className="mx-auto h-12 w-12 opacity-50"/>
                                <p className="mt-4">No entity definitions yet</p>
                                <p className="mt-1 text-sm">Create your first one to get started</p>
                            </div>
                        ) : (
                            <div className="space-y-2">
                                {entityDefinitions.map((definition) => (
                                    <button
                                        key={definition.id}
                                        onClick={() => handleSelectDefinition(definition)}
                                        className={`w-full rounded-lg border-2 p-4 text-left transition-all hover:border-primary ${
                                            selectedDefinition?.id === definition.id
                                                ? 'border-primary bg-primary/10'
                                                : 'border-base-300'
                                        }`}
                                    >
                                        <div className="flex items-start justify-between">
                                            <div className="flex-1">
                                                <h3 className="font-semibold">{definition.name}</h3>
                                                <p className="text-sm text-base-content/60">
                                                    ID: {definition.id}
                                                </p>
                                                {definition.description && (
                                                    <p className="mt-1 text-sm text-base-content/70">
                                                        {definition.description}
                                                    </p>
                                                )}
                                                <div className="mt-2 flex gap-2">
                                                    <span
                                                        className="badge badge-sm badge-outline">
                                                        {definition.layout?.components?.length} components
                                                    </span>
                                                    {definition.allowInherit && (
                                                        <span className="badge badge-sm badge-primary">
                                                            Inheritable
                                                        </span>
                                                    )}
                                                </div>
                                            </div>
                                            <ChevronRightIcon className="h-5 w-5 text-base-content/40"/>
                                        </div>
                                    </button>
                                ))}
                            </div>
                        )}
                    </CardBody>
                </Card>

                {/* Right Panel - Entity Details */}
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

            {/* Create Modal */}
            {isCreateModalOpen && (
                <div className="modal modal-open">
                    <div className="modal-box">
                        <h3 className="text-lg font-bold">Create New Entity Definition</h3>
                        <p className="py-2 text-sm text-base-content/70">
                            Define a new entity type. You can add components later.
                        </p>

                        {createError && (
                            <Alert variant="error" className="mt-4">
                                {createError}
                            </Alert>
                        )}

                        <div className="space-y-4 py-4">
                            <Input
                                label="Entity ID"
                                placeholder="e.g., product, customer, order"
                                value={newEntityId}
                                onChange={(e) => setNewEntityId(e.target.value.toLowerCase())}
                                fullWidth
                                error={
                                    newEntityId && !/^[a-z0-9_-]+$/.test(newEntityId)
                                        ? 'Only lowercase letters, numbers, underscores, and dashes'
                                        : undefined
                                }
                            />
                            <Input
                                label="Entity Name"
                                placeholder="e.g., Product, Customer, Order"
                                value={newEntityName}
                                onChange={(e) => setNewEntityName(e.target.value)}
                                fullWidth
                            />
                        </div>

                        <div className="modal-action">
                            <Button
                                variant="ghost"
                                onClick={() => {
                                    setIsCreateModalOpen(false)
                                    setNewEntityId('')
                                    setNewEntityName('')
                                    setCreateError(null)
                                }}
                                disabled={createMutation.isPending}
                            >
                                Cancel
                            </Button>
                            <Button
                                variant="primary"
                                onClick={handleCreateEntity}
                                loading={createMutation.isPending}
                                disabled={!newEntityId.trim() || !newEntityName.trim()}
                            >
                                Create
                            </Button>
                        </div>
                    </div>
                    <div className="modal-backdrop" onClick={() => setIsCreateModalOpen(false)}/>
                </div>
            )}
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