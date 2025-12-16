import {apiClient} from '../lib/api'
import {User} from '../stores/authStore'

export interface LoginResponse {
    user: User
    message?: string
}

export const authApi = {
    // Check current authentication status
    checkAuth: async (): Promise<User | null> => {
        try {
            const response = await apiClient.get<User>('/api/v1/auth/me')
            return response.data
        } catch (error) {
            console.log(error)
            return null
        }
    },

    // Initiate OAuth login (redirect-based)
    login: (provider: 'google' | 'github') => {
        window.location.href = `${apiClient.defaults.baseURL}/auth/${provider}`
    },

    // Logout
    logout: async (provider: 'google' | 'github') => {
        await apiClient.post(`/auth/${provider}/logout`)
    },
}
