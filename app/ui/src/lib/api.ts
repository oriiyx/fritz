import axios, {AxiosError} from 'axios'
import Cookies from 'js-cookie'
import {ApiErrorResponse, getErrorDetails} from './errorHandler'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
})

// Request interceptor
apiClient.interceptors.request.use(
    (config) => {
        // Add CSRF token if available
        const csrfToken = Cookies.get('csrf_token')
        if (csrfToken) {
            config.headers['X-CSRF-Token'] = csrfToken
        }
        return config
    },
    (error: unknown) => {
        return Promise.reject(error)
    }
)

// Response interceptor for centralized error handling
apiClient.interceptors.response.use(
    (response) => response,
    (error: unknown) => {
        // Log errors in development
        if (import.meta.env.DEV) {
            const details = getErrorDetails(error)
            console.error('[API Error]', details)
        }

        // Handle specific error codes globally if needed
        if (error instanceof Error && 'response' in error) {
            const axiosError = error as AxiosError<ApiErrorResponse>

            // Handle 401 - redirect to login (could be done here or in components)
            if (axiosError.response?.status === 401) {
                // You could dispatch a global auth event here
                console.debug('Unauthorized - authentication required')
            }

            // Handle 403 - forbidden
            if (axiosError.response?.status === 403) {
                console.warn('Access forbidden')
            }

            // Handle 500 - server error
            if (axiosError.response?.status === 500) {
                console.error('Server error occurred')
            }
        }

        return Promise.reject(error)
    }
)