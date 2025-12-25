import {useState} from 'react'
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {Card, CardBody, CardTitle} from '@/components/Card'
import {Alert} from '@/components/Alert'
import {Loading} from '@/components/Loading'
import {Button} from '@/components/Button'
import {entitiesApi} from '@/services/entitiesService'
import {DataComponent, DataComponentDefinition, DataComponentType, EntityDefinition,} from '@/generated/definitions'
import {EntityList} from '@/components/entities/EntityList.tsx'
import {ComponentLayoutTree} from '@/components/entities/ComponentLayoutTree'
import {AddComponentDropdown} from '@/components/entities/AddComponentDropdown'
import {ComponentSettingsPanel} from '@/components/entities/ComponentSettingsPanel'
import {CheckCircleIcon} from '@heroicons/react/24/outline'

export function EntitiesPage() {
    const queryClient = useQueryClient()

    // Local state management
    const [selectedDefinition, setSelectedDefinition] = useState<EntityDefinition | null>(null)
    const [selectedComponent, setSelectedComponent] = useState<DataComponent | null>(null)
    const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false)
    const [validationError, setValidationError] = useState<string | null>(null)

    // Fetch entity definitions using React Query
    const {data: entityDefinitions = [], isLoading, error} = useQuery({
        queryKey: ['entity-definitions'],
        queryFn: entitiesApi.getEntityDefinitions,
    })

    // Fetch data component types using React Query
    const {
        data: dataComponentTypes = [],
        isLoading: isDataComponentLoading,
        error: dataComponentError,
    } = useQuery({
        queryKey: ['data-component-types'],
        queryFn: entitiesApi.getDataComponentTypes,
    })

    // Save mutation
    const saveMutation = useMutation({
        mutationFn: async (definition: EntityDefinition) => {
            await entitiesApi.updateEntityDefinition(definition.id, definition)
        },
        onSuccess: () => {
            setHasUnsavedChanges(false)
        },
    })

    // Delete mutation
    const deleteMutation = useMutation({
        mutationFn: async (id: string) => {
            await entitiesApi.deleteEntityDefinition(id)
        },
        onSuccess: () => {
            // Invalidate and refetch entity definitions
            queryClient.invalidateQueries({queryKey: ['entity-definitions']})

            // Clear selection
            setSelectedDefinition(null)
            setSelectedComponent(null)
            setHasUnsavedChanges(false)
        },
    })

    const handleSelectDefinition = (definition: EntityDefinition) => {
        if (hasUnsavedChanges) {
            if (!confirm('You have unsaved changes. Do you want to discard them?')) {
                return
            }
        }
        setSelectedDefinition(definition)
        setSelectedComponent(null)
        setHasUnsavedChanges(false)
        setValidationError(null)
    }

    const handleComponentsReorder = (newComponents: DataComponent[]) => {
        if (!selectedDefinition) return

        const updatedDefinition = {
            ...selectedDefinition,
            layout: {
                ...selectedDefinition.layout,
                components: newComponents,
            },
        }

        setSelectedDefinition(updatedDefinition)
        setHasUnsavedChanges(true)
        setValidationError(null)
    }

    const handleComponentSelect = (component: DataComponent) => {
        setSelectedComponent(component)
    }

    const handleComponentUpdate = (updatedComponent: DataComponent) => {
        if (!selectedDefinition) return

        const updatedComponents = selectedDefinition.layout.components.map((c) =>
            c.name === updatedComponent.name ? updatedComponent : c
        )

        const updatedDefinition = {
            ...selectedDefinition,
            layout: {
                ...selectedDefinition.layout,
                components: updatedComponents,
            },
        }

        setSelectedDefinition(updatedDefinition)
        setSelectedComponent(updatedComponent)
        setHasUnsavedChanges(true)
    }

    const handleComponentDelete = (component: DataComponent) => {
        if (!selectedDefinition) return

        if (!confirm(`Delete component "${component.title}"?`)) return

        const updatedComponents = selectedDefinition.layout.components.filter(
            (c) => c.name !== component.name
        )

        const updatedDefinition = {
            ...selectedDefinition,
            layout: {
                ...selectedDefinition.layout,
                components: updatedComponents,
            },
        }

        setSelectedDefinition(updatedDefinition)

        if (selectedComponent?.name === component.name) {
            setSelectedComponent(null)
        }

        setHasUnsavedChanges(true)
    }

    const handleAddComponent = (componentType: DataComponentDefinition) => {
        if (!selectedDefinition) return

        const baseName = componentType.id.toLowerCase()
        const existingNames = selectedDefinition.layout.components.map((c) => c.name)
        let counter = 1
        let newName = baseName

        while (existingNames.includes(newName)) {
            newName = `${baseName}${counter}`
            counter++
        }

        const newComponent: DataComponent = {
            type: componentType.id as DataComponentType,
            name: newName,
            title: componentType.label,
            dbtype: componentType.defaultDBType,
            mandatory: false,
            invisible: false,
            notEditable: false,
            settings: {},
        }

        const updatedComponents = [...selectedDefinition.layout.components, newComponent]

        const updatedDefinition = {
            ...selectedDefinition,
            layout: {
                ...selectedDefinition.layout,
                components: updatedComponents,
            },
        }

        setSelectedDefinition(updatedDefinition)
        setSelectedComponent(newComponent)
        setHasUnsavedChanges(true)
        setValidationError(null)
    }

    const handleSave = () => {
        if (!selectedDefinition) return

        // Validation: check if at least one component exists
        if (selectedDefinition.layout.components.length === 0) {
            setValidationError('Definition must have at least one component before saving')
            return
        }

        // Clear validation error
        setValidationError(null)

        saveMutation.mutate(selectedDefinition)
    }

    const handleDelete = () => {
        if (!selectedDefinition) return

        const confirmMessage = `Are you sure you want to delete "${selectedDefinition.name}"?\n\nThis action cannot be undone and will permanently remove this entity definition and all its components.`

        if (!confirm(confirmMessage)) return

        deleteMutation.mutate(selectedDefinition.id)
    }

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
                <Alert variant="error">Failed to load entity definitions. Please try again.</Alert>
            </div>
        )
    }

    return (
        <div className="space-y-6 pb-24">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold">Entity Definitions</h1>
                    <p className="mt-2 text-base-content/70">Manage your entity schemas and structures</p>
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
                    {selectedDefinition ? (
                        <div className="space-y-6">
                            {/* Basic Info Card */}
                            <Card>
                                <CardBody>
                                    <div className="flex items-center justify-between">
                                        <CardTitle>Definition Details</CardTitle>
                                        <Button
                                            variant="error"
                                            outline
                                            size="sm"
                                            onClick={handleDelete}
                                            disabled={deleteMutation.isPending}
                                            loading={deleteMutation.isPending}
                                        >
                                            {deleteMutation.isPending ? 'Deleting...' : 'Delete Definition'}
                                        </Button>
                                    </div>
                                    <div className="mt-4 space-y-3">
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
                                </CardBody>
                            </Card>

                            {/* Component Management Split View */}
                            <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
                                {/* Left: Component Tree (1/3) */}
                                <div className="space-y-4">
                                    <Card>
                                        <CardBody>
                                            <CardTitle>Components</CardTitle>
                                            <div className="mt-4 space-y-4">
                                                <AddComponentDropdown
                                                    dataComponentTypes={dataComponentTypes}
                                                    onAddComponent={handleAddComponent}
                                                />
                                                <ComponentLayoutTree
                                                    components={selectedDefinition.layout.components}
                                                    selectedComponent={selectedComponent}
                                                    onComponentsReorder={handleComponentsReorder}
                                                    onComponentSelect={handleComponentSelect}
                                                    onComponentDelete={handleComponentDelete}
                                                />
                                            </div>
                                        </CardBody>
                                    </Card>
                                </div>

                                {/* Right: Settings Panel (2/3) */}
                                <div className="lg:col-span-2">
                                    <ComponentSettingsPanel
                                        component={selectedComponent}
                                        onComponentUpdate={handleComponentUpdate}
                                    />
                                </div>
                            </div>
                        </div>
                    ) : (
                        <Card>
                            <CardBody>
                                <div className="flex h-96 items-center justify-center text-base-content/50">
                                    <div className="text-center">
                                        <p className="text-lg">No definition selected</p>
                                        <p className="mt-2 text-sm">
                                            Select a definition from the list to view and edit its components
                                        </p>
                                    </div>
                                </div>
                            </CardBody>
                        </Card>
                    )}
                </div>
            </div>

            {/* Fixed Bottom Bar - Only show when a definition is selected */}
            {selectedDefinition && (
                <div className="fixed bottom-0 left-0 right-0 z-50 border-t border-base-300 bg-base-100 shadow-lg">
                    <div className="mx-auto flex items-center justify-between px-6 py-4">
                        <div className="flex items-center gap-3">
                            {validationError ? (
                                <>
                                    <span className="badge badge-error">Validation Error</span>
                                    <span className="text-sm text-error">{validationError}</span>
                                </>
                            ) : hasUnsavedChanges ? (
                                <>
                                    <span className="badge badge-warning">Unsaved Changes</span>
                                    <span className="text-sm text-base-content/70">
                                        You have unsaved changes to {selectedDefinition.name}
                                    </span>
                                </>
                            ) : (
                                <>
                                    <CheckCircleIcon className="h-5 w-5 text-success"/>
                                    <span className="text-sm text-base-content/70">All changes saved</span>
                                </>
                            )}
                        </div>
                        <Button
                            variant="primary"
                            onClick={handleSave}
                            disabled={!hasUnsavedChanges || saveMutation.isPending}
                            loading={saveMutation.isPending}
                        >
                            {saveMutation.isPending ? 'Saving...' : 'Save Changes'}
                        </Button>
                    </div>
                </div>
            )}
        </div>
    )
}