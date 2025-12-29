import {useState} from 'react'
import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {XMarkIcon} from '@heroicons/react/24/outline'
import {Card, CardBody, CardTitle} from '@/components/Card'
import {Input} from '@/components/Input'
import {Button} from '@/components/Button'
import {Alert} from '@/components/Alert'
import {entitiesApi} from '@/services/entitiesService'
import {definitionApi} from "@/services/definitionService.ts";

interface AddEntityModalProps {
    parentId: string | null
    parentPath: string
    onClose: () => void
    onSuccess?: () => void
}

export function AddEntityModal({parentId, parentPath, onClose, onSuccess}: AddEntityModalProps) {
    const queryClient = useQueryClient()
    const [selectedDefinitionId, setSelectedDefinitionId] = useState<string>('')
    const [entityKey, setEntityKey] = useState('')
    const [error, setError] = useState<string | null>(null)

    // Fetch available entity definitions
    const {data: definitions = [], isLoading: isLoadingDefinitions} = useQuery({
        queryKey: ['entity-definitions'],
        queryFn: definitionApi.getEntityDefinitions,
    })

    // Create entity mutation
    const createMutation = useMutation({
        mutationFn: async ({definitionId, key, path}: { definitionId: string; key: string; path: string }) => {
            return entitiesApi.createEntity(definitionId, {
                parent_id: parentId,
                key: key,
                path: path,
                published: false,
                // TODO: Data field will be implemented when backend is ready
            })
        },
        onSuccess: () => {
            // Invalidate tree queries to refresh the tree
            queryClient.invalidateQueries({queryKey: ['tree-children']})
            onSuccess?.()
            onClose()
        },
        onError: (err: any) => {
            setError(err.message || 'Failed to create entity')
        },
    })

    const handleSubmit = () => {
        // Validation
        if (!selectedDefinitionId) {
            setError('Please select an entity definition')
            return
        }

        if (!entityKey.trim()) {
            setError('Entity name is required')
            return
        }

        // Construct the entity path
        const entityPath = parentPath === '/' ? `/${entityKey}` : `${parentPath}/${entityKey}`

        setError(null)
        createMutation.mutate({
            definitionId: selectedDefinitionId,
            key: entityKey,
            path: entityPath,
        })
    }

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
            <Card className="w-full max-w-md">
                <CardBody>
                    <div className="flex items-center justify-between mb-4">
                        <CardTitle>Add New Entity</CardTitle>
                        <button
                            onClick={onClose}
                            className="btn btn-ghost btn-sm btn-circle"
                            disabled={createMutation.isPending}
                        >
                            <XMarkIcon className="h-5 w-5"/>
                        </button>
                    </div>

                    <div className="space-y-4">
                        {/* Parent Path Info */}
                        <div className="alert alert-info">
                            <div className="text-sm">
                                <strong>Parent:</strong> {parentPath || 'System Root'}
                            </div>
                        </div>

                        {error && (
                            <Alert variant="error">
                                {error}
                            </Alert>
                        )}

                        {/* Entity Definition Selector */}
                        <div className="form-control">
                            <label className="label">
                                <span className="label-text">Entity Definition *</span>
                            </label>
                            {isLoadingDefinitions ? (
                                <div className="skeleton h-12 w-full"></div>
                            ) : definitions.length === 0 ? (
                                <Alert variant="warning">
                                    No entity definitions available. Please create one first in the Definitions page.
                                </Alert>
                            ) : (
                                <select
                                    className="select select-bordered w-full"
                                    value={selectedDefinitionId}
                                    onChange={(e) => setSelectedDefinitionId(e.target.value)}
                                    disabled={createMutation.isPending}
                                >
                                    <option value="">Select a definition...</option>
                                    {definitions.map((def) => (
                                        <option key={def.id} value={def.id}>
                                            {def.name} ({def.id})
                                        </option>
                                    ))}
                                </select>
                            )}
                        </div>

                        {/* Entity Name Input */}
                        <Input
                            label="Entity Name *"
                            placeholder="e.g., my-product, homepage, etc."
                            value={entityKey}
                            onChange={(e) => setEntityKey(e.target.value)}
                            fullWidth
                            disabled={createMutation.isPending}
                        />

                        {/* Path Preview */}
                        {entityKey && (
                            <div className="alert">
                                <div className="text-sm">
                                    <strong>Path:</strong>{' '}
                                    <code className="font-mono">
                                        {parentPath === '/' ? `/${entityKey}` : `${parentPath}/${entityKey}`}
                                    </code>
                                </div>
                            </div>
                        )}

                        {/* TODO Notice */}
                        <div className="alert alert-warning">
                            <span className="text-xs">
                                Note: Entity data fields will be editable once the backend is ready.
                            </span>
                        </div>

                        {/* Action Buttons */}
                        <div className="flex gap-2 justify-end">
                            <Button
                                variant="ghost"
                                onClick={onClose}
                                disabled={createMutation.isPending}
                            >
                                Cancel
                            </Button>
                            <Button
                                variant="primary"
                                onClick={handleSubmit}
                                disabled={createMutation.isPending || !selectedDefinitionId || !entityKey.trim()}
                                loading={createMutation.isPending}
                            >
                                {createMutation.isPending ? 'Creating...' : 'Create Entity'}
                            </Button>
                        </div>
                    </div>
                </CardBody>
            </Card>
        </div>
    )
}