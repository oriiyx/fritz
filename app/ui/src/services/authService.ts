import {apiClient} from '../lib/api'
import {User} from '../stores/authStore'

export const authApi = {
    // Check current authentication status
    checkAuth: async (): Promise<User | null> => {
        try {
            const response = await apiClient.get<User>('/api/v1/auth/me')
            return response.data
        } catch (error: any) {
            // Don't log 401 errors - they're expected when not logged in
            if (error?.response?.status !== 401) {
                console.error('Auth check error:', error)
            }
            return null
        }
    },

    // Email/password login
    loginWithPassword: async (email: string, password: string): Promise<User> => {
        const response = await apiClient.post<User>('/api/v1/auth/login', {
            email: email,
            password: password
        })
        return response.data
    },

    // Initiate OAuth login (redirect-based)
    login: (provider: 'google' | 'github') => {
        window.location.href = `${apiClient.defaults.baseURL}/api/v1/auth/${provider}`
    },

    // Logout
    logout: async (provider: 'google' | 'github') => {
        await apiClient.post(`/api/v1/auth/${provider}/logout`)
    },
}