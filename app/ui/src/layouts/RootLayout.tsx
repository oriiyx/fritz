import {useEffect} from 'react'
import {Outlet, useNavigate, useRouterState} from '@tanstack/react-router'
import {useQuery} from '@tanstack/react-query'
import {useAuthStore} from '../stores/authStore'
import {authApi} from '../services/authService'

export function RootLayout() {
    const navigate = useNavigate()
    const {isAuthenticated, setUser, setLoading} = useAuthStore()
    const routerState = useRouterState()
    const currentPath = routerState.location.pathname

    // Check authentication status on mount
    const {data: user, isLoading} = useQuery({
        queryKey: ['auth-check'],
        queryFn: authApi.checkAuth,
        retry: false,
        staleTime: 5 * 60 * 1000, // 5 minutes
    })

    useEffect(() => {
        if (!isLoading) {
            setUser(user || null)

            // Redirect logic
            if (user && currentPath === '/login') {
                navigate({to: '/'})
            } else if (!user && currentPath !== '/login') {
                navigate({to: '/login'})
            }
        }
    }, [user, isLoading, currentPath, setUser, navigate])

    // Show loading state while checking authentication
    if (isLoading) {
        return (
            <div className="flex h-screen items-center justify-center">
                <div className="loading loading-spinner loading-lg"></div>
            </div>
        )
    }

    return <Outlet/>
}
