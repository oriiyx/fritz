import {useEffect, useState} from 'react'
import {useMutation, useQuery} from '@tanstack/react-query'
import {DataComponent, EntityDefinition} from '@/generated/definitions'
import {entitiesApi, Entity, ModifyEntityResponse} from '@/services/entitiesService'
import {Card, CardBody, CardTitle} from '@/components/Card'
import {Button} from '@/components/Button'
import {Alert} from '@/components/Alert'
import {InputField} from './fields/InputField'
import {IntegerField} from './fields/IntegerField'
import {DateField} from './fields/DateField'
import {CheckCircleIcon} from '@heroicons/react/24/outline'
import {Loading} from '@/components/Loading'

interface EntityFormProps {
    entity: Entity
    definition: EntityDefinition
    onSaveSuccess?: (response: ModifyEntityResponse) => void
}

export function EntityForm({entity, definition, onSaveSuccess}: EntityFormProps) {
    const [formData, setFormData] = useState<Record<string, unknown>>({})
    const [errors, setErrors] = useState<Record<string, string>>({})
    const [hasChanges, setHasChanges] = useState(false)
    const [isInitialized, setIsInitialized] = useState(false)

    const isTransitionMode = !entity.has_data
    const isEditMode = entity.has_data

    // Fetch existing data if entity has data
    const {
        data: entityData,
        isLoading: isLoadingData,
        error: dataError,
    } = useQuery({
        queryKey: ['entity-data', entity.entity_class, entity.id],
        queryFn: () => entitiesApi.readEntity(entity.entity_class, entity.id),
        enabled: isEditMode,
        staleTime: 0, // Always fetch fresh data
    })

    // Initialize form data
    useEffect(() => {
        if (isTransitionMode && !isInitialized) {
            // Transition mode: Initialize with defaults
            const initialData: Record<string, unknown> = {}
            definition.layout.components.forEach(component => {
                if (component.invisible) return

                const settings = component.settings || {}
                switch (component.type) {
                    case 'input':
                        initialData[component.name] = settings.defaultValue || ''
                        break
                    case 'integer':
                        initialData[component.name] = settings.defaultValue ?? null
                        break
                    case 'date':
                        initialData[component.name] = settings.defaultValue || null
                        break
                    default:
                        initialData[component.name] = null
                }
            })
            setFormData(initialData)
            setIsInitialized(true)
        } else if (isEditMode && entityData && !isInitialized) {
            // Edit mode: Initialize with existing data
            setFormData(entityData.data)
            setIsInitialized(true)
        }
    }, [isTransitionMode, isEditMode, entityData, definition.layout.components, isInitialized])

    const saveMutation = useMutation({
        mutationFn: async () => {
            if (isTransitionMode) {
                return entitiesApi.transitionEntity(entity.entity_class, entity.id, {data: formData})
            } else {
                return entitiesApi.saveEntity(entity.entity_class, {
                    id: entity.id,
                    parent_id: entity.parent_id,
                    key: entity.o_key,
                    path: entity.o_path,
                    published: entity.published,
                    type: entity.o_type,
                    data: formData,
                })
            }
        },
        onSuccess: (response) => {
            setHasChanges(false)
            onSaveSuccess?.(response)
        },
    })

    const handleFieldChange = (name: string, value: unknown) => {
        setFormData(prev => ({
            ...prev,
            [name]: value
        }))
        setHasChanges(true)

        // Clear error for this field
        if (errors[name]) {
            setErrors(prev => {
                const newErrors = {...prev}
                delete newErrors[name]
                return newErrors
            })
        }
    }

    const validateForm = (): boolean => {
        const newErrors: Record<string, string> = {}

        definition.layout.components.forEach(component => {
            if (component.invisible) return

            const value = formData[component.name]

            // Check mandatory fields
            if (component.mandatory) {
                if (value === null || value === undefined || value === '') {
                    newErrors[component.name] = `${component.title} is required`
                }
            }

            // Type-specific validation
            if (value !== null && value !== undefined && value !== '') {
                const settings = component.settings || {}

                switch (component.type) {
                    case 'input':
                        if (settings.regexValidation) {
                            try {
                                const regex = new RegExp(settings.regexValidation)
                                if (!regex.test(String(value))) {
                                    newErrors[component.name] = `${component.title} does not match the required pattern`
                                }
                            } catch {
                                // Invalid regex - skip validation
                            }
                        }
                        if (settings.columnLength && String(value).length > settings.columnLength) {
                            newErrors[component.name] = `${component.title} must be at most ${settings.columnLength} characters`
                        }
                        break

                    case 'integer':
                        const numValue = Number(value)
                        if (settings.minValue !== undefined && numValue < settings.minValue) {
                            newErrors[component.name] = `${component.title} must be at least ${settings.minValue}`
                        }
                        if (settings.maxValue !== undefined && numValue > settings.maxValue) {
                            newErrors[component.name] = `${component.title} must be at most ${settings.maxValue}`
                        }
                        if (settings.unsigned && numValue < 0) {
                            newErrors[component.name] = `${component.title} must be a positive number`
                        }
                        break
                }
            }
        })

        setErrors(newErrors)
        return Object.keys(newErrors).length === 0
    }

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        if (!validateForm()) {
            return
        }

        saveMutation.mutate()
    }

    const renderField = (component: DataComponent) => {
        if (component.invisible) return null

        const commonProps = {
            component,
            value: formData[component.name],
            onChange: (value: unknown) => handleFieldChange(component.name, value),
            error: errors[component.name],
            disabled: component.notEditable,
        }

        switch (component.type) {
            case 'input':
                return <InputField key={component.id} {...commonProps} />
            case 'integer':
                return <IntegerField key={component.id} {...commonProps} />
            case 'date':
                return <DateField key={component.id} {...commonProps} />
            default:
                return (
                    <div key={component.id} className="form-control fieldset">
                        <label className="label">
                            <span className="label-text">{component.title}</span>
                        </label>
                        <div className="text-sm text-base-content/50">
                            Unsupported field type: {component.type}
                        </div>
                    </div>
                )
        }
    }

    // Loading state for edit mode
    if (isEditMode && (isLoadingData || !isInitialized)) {
        return (
            <Card>
                <CardBody>
                    <div className="flex h-96 items-center justify-center">
                        <div className="flex flex-col items-center gap-4">
                            <Loading size="lg"/>
                            <p className="text-base-content/70">Loading entity data...</p>
                        </div>
                    </div>
                </CardBody>
            </Card>
        )
    }

    // Error state for edit mode
    if (isEditMode && dataError) {
        return (
            <Card>
                <CardBody>
                    <Alert variant="error">
                        Failed to load entity data: {dataError.message}
                    </Alert>
                </CardBody>
            </Card>
        )
    }

    const visibleComponents = definition.layout.components.filter(c => !c.invisible)

    return (
        <form onSubmit={handleSubmit}>
            <Card>
                <CardBody>
                    <div className="flex items-center justify-between">
                        <div>
                            <CardTitle>{definition.name}</CardTitle>
                            {definition.description && (
                                <p className="text-sm text-base-content/70 mt-1">{definition.description}</p>
                            )}
                        </div>
                        {isEditMode && (
                            <span className="badge badge-info badge-sm">Edit Mode</span>
                        )}
                        {isTransitionMode && (
                            <span className="badge badge-warning badge-sm">First Save</span>
                        )}
                    </div>

                    {saveMutation.error && (
                        <Alert variant="error" className="mt-4">
                            {saveMutation.error.message}
                        </Alert>
                    )}

                    <div className="mt-6 space-y-4">
                        {visibleComponents.map(renderField)}
                    </div>
                </CardBody>
            </Card>

            {/* Fixed Bottom Bar */}
            <div className="fixed bottom-0 left-0 right-0 z-50 border-t border-base-300 bg-base-100 shadow-lg">
                <div className="mx-auto flex items-center justify-between px-6 py-4">
                    <div className="flex items-center gap-3">
                        {Object.keys(errors).length > 0 ? (
                            <>
                                <span className="badge badge-error">Validation Error</span>
                                <span className="text-sm text-error">
                                    Please fix the errors above
                                </span>
                            </>
                        ) : hasChanges ? (
                            <>
                                <span className="badge badge-warning">Unsaved Changes</span>
                                <span className="text-sm text-base-content/70">
                                    You have unsaved changes
                                </span>
                            </>
                        ) : saveMutation.isSuccess ? (
                            <>
                                <CheckCircleIcon className="h-5 w-5 text-success"/>
                                <span className="text-sm text-base-content/70">All changes saved</span>
                            </>
                        ) : isTransitionMode ? (
                            <>
                                <span className="badge badge-info">New Entity</span>
                                <span className="text-sm text-base-content/70">
                                    Fill in the form and save
                                </span>
                            </>
                        ) : (
                            <>
                                <CheckCircleIcon className="h-5 w-5 text-success"/>
                                <span className="text-sm text-base-content/70">No changes</span>
                            </>
                        )}
                    </div>
                    <Button
                        type="submit"
                        variant="primary"
                        disabled={!hasChanges || saveMutation.isPending}
                        loading={saveMutation.isPending}
                    >
                        {saveMutation.isPending ? 'Saving...' : 'Save Entity'}
                    </Button>
                </div>
            </div>
        </form>
    )
}