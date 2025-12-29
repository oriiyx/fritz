import {apiClient} from '../lib/api'
import {getErrorDetails, getErrorMessage} from '../lib/errorHandler'

export interface CreateEntityRequest {
    parent_id?: string | null
    key: string
    path: string
    type?: string
    published: boolean
    data?: Record<string, unknown>
}

export interface CreateEntityResponse {
    id: string
    entity_class: string
    parent_id: string | null
    o_key: string
    o_path: string
    o_type: string
    published: boolean
    created_at: string
    updated_at: string
}

export const entitiesApi = {
    // Create new entity instance
    createEntity: async (
        definitionId: string,
        request: CreateEntityRequest
    ): Promise<CreateEntityResponse> => {
        try {
            const response = await apiClient.post<CreateEntityResponse>(
                `/api/v1/entities/${definitionId}/create`,
                request
            )
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Create entity error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },
}