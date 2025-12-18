import {AxiosError, isAxiosError} from 'axios'

/**
 * Standard error response structure from the Fritz API
 */
export interface ApiErrorResponse {
    error: string
    errors?: string[]
    conflictingId?: string
    conflictingName?: string
    conflictingComponentName?: string
}

/**
 * Extracts a user-friendly error message from various error types
 * @param error - The caught error (unknown type)
 * @returns A user-friendly error message
 */
export function getErrorMessage(error: unknown): string {
    // Check if it's an Axios error using the type guard
    if (isAxiosError<ApiErrorResponse>(error)) {
        // API returned an error response
        if (error.response?.data?.error) {
            return error.response.data.error
        }

        // Network error or timeout
        if (error.code === 'ECONNABORTED') {
            return 'Request timeout. Please try again.'
        }

        if (error.code === 'ERR_NETWORK') {
            return 'Network error. Please check your connection.'
        }

        // Generic Axios error
        return error.message || 'An unexpected error occurred'
    }

    // Check if it's a standard Error object
    if (error instanceof Error) {
        return error.message
    }

    // Fallback for unknown error types
    if (typeof error === 'string') {
        return error
    }

    return 'An unexpected error occurred'
}

/**
 * Extracts detailed error information for debugging
 * @param error - The caught error (unknown type)
 * @returns Detailed error information
 */
export function getErrorDetails(error: unknown): {
    message: string
    statusCode?: number
    data?: ApiErrorResponse
} {
    if (isAxiosError<ApiErrorResponse>(error)) {
        return {
            message: getErrorMessage(error),
            statusCode: error.response?.status,
            data: error.response?.data,
        }
    }

    return {
        message: getErrorMessage(error),
    }
}

/**
 * Type guard to check if error is an AxiosError
 * @param error - The error to check
 */
export function isApiError(error: unknown): error is AxiosError<ApiErrorResponse> {
    return isAxiosError<ApiErrorResponse>(error)
}

/**
 * Custom error class for application-specific errors
 */
export class AppError extends Error {
    constructor(
        message: string,
        public statusCode?: number,
        public originalError?: unknown
    ) {
        super(message)
        this.name = 'AppError'
    }
}