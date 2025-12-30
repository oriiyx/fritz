import {apiClient} from '../lib/api'
import {getErrorDetails, getErrorMessage} from '../lib/errorHandler'

// Request for creating entity metadata only (no data)
export interface CreateEntityRequest {
    parent_id?: string | null
    key: string
    path: string
    type?: string
    published: boolean
}

// Request for saving entity data
export interface SaveEntityRequest {
    data: Record<string, unknown>
}

// Response from entity creation
export interface CreateEntityResponse {
    entity: {
        id: string
        entity_class: string
        parent_id: string | null
        o_key: string
        o_path: string
        o_type: string
        published: boolean
        has_data: boolean // NEW: Indicates if entity has data saved
        created_at: string
        updated_at: string
    }
}

// Response from entity save
export interface SaveEntityResponse {
    entity: {
        id: string
        entity_class: string
        parent_id: string | null
        o_key: string
        o_path: string
        o_type: string
        published: boolean
        has_data: boolean
        created_at: string
        updated_at: string
    }
    data: Record<string, unknown> // The actual entity data
}

export const entitiesApi = {
    /**
     * Create new entity instance (metadata only)
     * This creates the entity in the entities table but doesn't populate
     * the entity_{class_id} table yet. Call saveEntity() to save data.
     */
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

    /**
     * Save entity data (works for both first save and updates)
     * - If entity.has_data === false: Creates entry in entity_{class_id} table
     * - If entity.has_data === true: Updates existing entry in entity_{class_id} table
     */
    saveEntity: async (
        definitionId: string,
        entityId: string,
        request: SaveEntityRequest
    ): Promise<SaveEntityResponse> => {
        try {
            const response = await apiClient.post<SaveEntityResponse>(
                `/api/v1/entities/${definitionId}/${entityId}/save`,
                request
            )
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Save entity error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },
}