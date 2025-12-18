import {apiClient} from '../lib/api'
import {CreateNewEntityDefinition, EntityDefinition} from '../stores/entitiesDefinitionsStore'

export interface CreateEntityDefinitionRequest {
    id: string
    name: string
    description?: string
    allowInherit?: boolean
    layout?: {
        type: string
        components: Array<any>
    }
}

export const entitiesApi = {
    // Get all entity definitions
    getEntityDefinitions: async (): Promise<EntityDefinition[]> => {
        try {
            const response = await apiClient.get<EntityDefinition[]>('/api/v1/definitions')
            console.log(response)
            if (response.data === null) {
                console.log("empty")
                return []
            }
            return response.data
        } catch (error: any) {
            if (error?.response?.status !== 401) {
                console.error('Entity definition error:', error)
            }
            throw error
        }
    },

    // Get single entity definition by ID
    getEntityDefinition: async (id: string): Promise<EntityDefinition> => {
        try {
            const response = await apiClient.get<EntityDefinition>(`/api/v1/definitions/${id}`)
            return response.data
        } catch (error: any) {
            console.error('Entity definition error:', error)
            throw error
        }
    },

    // Create new entity definition
    createEntityDefinition: async (
        definition: CreateEntityDefinitionRequest
    ): Promise<void> => {
        try {
            const payload: CreateNewEntityDefinition = {
                id: definition.id,
                name: definition.name,
            }

            await apiClient.post<EntityDefinition>(
                '/api/v1/definitions/create',
                payload
            )
        } catch (error: any) {
            console.error('Create entity definition error:', error)
            throw error
        }
    },

    // Update entity definition
    updateEntityDefinition: async (
        id: string,
        definition: Partial<EntityDefinition>
    ): Promise<EntityDefinition> => {
        try {
            const response = await apiClient.put<EntityDefinition>(
                `/api/v1/definitions/${id}`,
                definition
            )
            return response.data
        } catch (error: any) {
            console.error('Update entity definition error:', error)
            throw error
        }
    },

    // Delete entity definition
    deleteEntityDefinition: async (id: string): Promise<void> => {
        try {
            await apiClient.delete(`/api/v1/definitions/${id}`)
        } catch (error: any) {
            console.error('Delete entity definition error:', error)
            throw error
        }
    },

    // Get available data component types
    getDataComponentTypes: async (): Promise<any[]> => {
        try {
            const response = await apiClient.get('/api/v1/definitions/data-component-types')
            return response.data
        } catch (error: any) {
            console.error('Data component types error:', error)
            throw error
        }
    },
}