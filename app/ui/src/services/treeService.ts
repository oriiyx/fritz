import {apiClient} from '@/lib/api.ts'
import {getErrorDetails, getErrorMessage} from '@/lib/errorHandler.ts'

export interface TreeNode {
    id: string
    entity_class: string
    parent_id: string | null
    o_key: string
    o_path: string
    o_type: 'folder' | 'object' | 'variant'
    published: boolean
    created_at: string
    updated_at: string
    has_children: boolean
    children_count: number
}

export interface GetChildrenParams {
    parent_id: string
    limit?: number
    offset?: number
    definition_id?: string
}

export interface GetChildrenResponse {
    items: TreeNode[]
    total: number
    limit: number
    offset: number
    has_more: boolean
}

export const treeApi = {
    // Get children of a parent node
    getChildren: async (params: GetChildrenParams): Promise<GetChildrenResponse> => {
        try {
            const response = await apiClient.post<GetChildrenResponse>(
                '/api/v1/entities/tree/children',
                {
                    parent_id: params.parent_id,
                    limit: params.limit || 25,
                    offset: params.offset || 0,
                    definition_id: params.definition_id,
                }
            )
            return response.data
        } catch (error: unknown) {
            const errorDetails = getErrorDetails(error)
            console.error('Tree children fetch error:', errorDetails)
            throw new Error(getErrorMessage(error))
        }
    },
}