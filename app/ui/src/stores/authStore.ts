import {create} from 'zustand'
import {persist} from 'zustand/middleware'

export interface User {
    id: string
    email: string
    full_name: string
    avatar_url?: string
}

interface AuthState {
    user: User | null
    isAuthenticated: boolean
    isLoading: boolean
    setUser: (user: User | null) => void
    logout: () => void
    setLoading: (loading: boolean) => void
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            isAuthenticated: false,
            isLoading: true,
            setUser: (user) =>
                set({
                    user,
                    isAuthenticated: !!user,
                    isLoading: false,
                }),
            logout: () =>
                set({
                    user: null,
                    isAuthenticated: false,
                }),
            setLoading: (loading) => set({isLoading: loading}),
        }),
        {
            name: 'fritz-auth',
        }
    )
)
