import {apiClient} from '../lib/api'
import {getErrorDetails, getErrorMessage} from '../lib/errorHandler'

// Entity from database
export interface Entity {
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
    created_by: string | null
    updated_by: string | null
}

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
    id: string
    parent_id?: string | null
    key: string
    path: string
    type?: string
    published: boolean
    data: Record<string, unknown>
}

// Request for transitioning entity data in the backend
export interface TransitionEntityRequest {
    data: Record<string, unknown>
}

// Response from entity creation
export interface CreateEntityResponse {
    entity: Entity
}

// Response from entity transition and save
export interface ModifyEntityResponse {
    entity: Entity
    data: Record<string, unknown>
}

// Response from reading entity data
export interface ReadEntityResponse {
    entity: Entity
    data: Record<string, unknown>
}

export const entitiesApi = {
    /**
     * Get entity metadata by ID
     */
    getEntity: async (id: string): Promise<Entity> => {
        try {
            const response = await apiClient.post<Entity>(
                '/api/v1/entities',
                {id}
            )
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Get entity error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },

    /**
     * Read entity data (requires has_data === true)
     * This fetches the actual data stored in entity_{class_id} table
     */
    readEntity: async (
        definitionId: string,
        entityId: string
    ): Promise<ReadEntityResponse> => {
        try {
            const response = await apiClient.post<ReadEntityResponse>(
                `/api/v1/entities/${definitionId}/read`,
                {id: entityId}
            )
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Read entity error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },

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
        request: SaveEntityRequest
    ): Promise<ModifyEntityResponse> => {
        try {
            const response = await apiClient.post<ModifyEntityResponse>(
                `/api/v1/entities/${definitionId}/save`,
                request
            )
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Save entity error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },

    /**
     * Transition entity data (first-time save when has_data === false)
     * This is specifically for creating the initial data entry
     */
    transitionEntity: async (
        definitionId: string,
        entityId: string,
        request: TransitionEntityRequest
    ): Promise<ModifyEntityResponse> => {
        try {
            const response = await apiClient.post<ModifyEntityResponse>(
                `/api/v1/entities/${definitionId}/${entityId}/transition`,
                request
            )
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Transition entity error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },

    /**
     * Delete entity
     */
    deleteEntity: async (
        definitionId: string,
        entityId: string,
    ): Promise<void> => {
        try {
            await apiClient.post<void>(
                `/api/v1/entities/${definitionId}/delete`,
                {
                    id: entityId
                }
            )
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Deleting entity error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },
}