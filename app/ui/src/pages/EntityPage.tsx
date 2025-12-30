import {useQuery, useQueryClient} from '@tanstack/react-query'
import {useNavigate, useParams} from '@tanstack/react-router'
import {entitiesApi} from '@/services/entitiesService'
import {definitionApi} from '@/services/definitionService'
import {EntityForm} from '@/components/entity/EntityForm'
import {Alert} from '@/components/Alert'
import {Loading} from '@/components/Loading'
import {Button} from '@/components/Button'
import {ArrowLeftIcon} from '@heroicons/react/24/outline'
import {EntityHeader} from "@/components/entity/EntityHeader.tsx"

export function EntityPage() {
    const {id} = useParams({from: '/dashboard/entity/$id'})
    const navigate = useNavigate()
    const queryClient = useQueryClient()

    // Fetch entity metadata
    const {
        data: entity,
        isLoading: isEntityLoading,
        error: entityError,
    } = useQuery({
        queryKey: ['entity', id],
        queryFn: () => entitiesApi.getEntity(id),
        staleTime: 0, // Always fetch fresh to detect has_data changes
    })

    // Fetch definition based on entity's class
    const {
        data: definition,
        isLoading: isDefinitionLoading,
        error: definitionError,
    } = useQuery({
        queryKey: ['definition', entity?.entity_class],
        queryFn: () => definitionApi.getEntityDefinition(entity!.entity_class),
        enabled: !!entity?.entity_class,
    })

    const isLoading = isEntityLoading || isDefinitionLoading
    const error = entityError || definitionError

    const handleSaveSuccess = () => {
        // Invalidate entity query to refetch metadata and detect has_data change
        queryClient.invalidateQueries({queryKey: ['entity', id]})

        // Also invalidate entity-data query for edit mode
        if (entity) {
            queryClient.invalidateQueries({queryKey: ['entity-data', entity.entity_class, entity.id]})
        }
    }

    if (isLoading) {
        return (
            <div className="flex h-96 items-center justify-center">
                <div className="flex flex-col items-center gap-4">
                    <Loading size="lg"/>
                    <p className="text-base-content/70">Loading entity...</p>
                </div>
            </div>
        )
    }

    if (error) {
        return (
            <div className="p-6">
                <Alert variant="error">
                    Could not load entity: {error.message}
                </Alert>
                <div className="mt-4">
                    <Button
                        variant="ghost"
                        iconLeft={<ArrowLeftIcon className="h-4 w-4"/>}
                        onClick={() => navigate({to: '/'})}
                    >
                        Back to Dashboard
                    </Button>
                </div>
            </div>
        )
    }

    if (!entity || !definition) {
        return (
            <div className="p-6">
                <Alert variant="error">
                    Entity or definition not found
                </Alert>
                <div className="mt-4">
                    <Button
                        variant="ghost"
                        iconLeft={<ArrowLeftIcon className="h-4 w-4"/>}
                        onClick={() => navigate({to: '/'})}
                    >
                        Back to Dashboard
                    </Button>
                </div>
            </div>
        )
    }

    return (
        <div className="p-6 pb-24">
            <EntityHeader entity={entity} definition={definition}/>
            <EntityForm
                entity={entity}
                definition={definition}
                onSaveSuccess={handleSaveSuccess}
            />
        </div>
    )
}