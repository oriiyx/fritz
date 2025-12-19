import {apiClient} from '../lib/api'
import {CreateNewEntityDefinition, EntityDefinition} from '../stores/entitiesDefinitionsStore'
import {getErrorDetails, getErrorMessage} from '../lib/errorHandler'
import {DataComponent} from "@/stores/dataComponentTypesStore.ts";

export interface CreateEntityDefinitionRequest {
    id: string
    name: string
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
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Entity definition error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },

    // Get single entity definition by ID
    getEntityDefinition: async (id: string): Promise<EntityDefinition> => {
        try {
            const response = await apiClient.get<EntityDefinition>(`/api/v1/definitions/${id}`)
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Entity definition error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },

    // Create new entity definition
    createEntityDefinition: async (
        definition: EntityDefinition
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
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Create entity definition error:', errorDetails)
            throw new Error(getErrorMessage(error))
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
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Update entity definition error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },

    // Delete entity definition
    deleteEntityDefinition: async (id: string): Promise<void> => {
        try {
            await apiClient.delete(`/api/v1/definitions/${id}`)
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Delete entity definition error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },

    // Get available data component types
    getDataComponentTypes: async (): Promise<DataComponent[]> => {
        try {
            const response = await apiClient.get<DataComponent[]>('/api/v1/definitions/data-component-types')
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Data component types error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },
}