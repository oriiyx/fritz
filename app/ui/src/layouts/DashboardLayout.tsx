import {Outlet} from '@tanstack/react-router'
import {useEffect} from 'react'
import {useQuery} from '@tanstack/react-query'
import {Header} from '../components/layout/Header'
import {Footer} from '../components/layout/Footer'
import {Sidebar} from '../components/layout/Sidebar'
import {useAuthStore} from '../stores/authStore'
import {authApi} from '../services/authService'

export function DashboardLayout() {
    const {setUser} = useAuthStore()

    // We only fetch user data for UI updates
    // The beforeLoad guard already verified we're authenticated
    // So this query is just for keeping the zustand store fresh
    const {data: user} = useQuery({
        queryKey: ['current-user'],
        queryFn: authApi.checkAuth,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: false,
        refetchOnMount: false, // Don't refetch on mount since beforeLoad just checked
        retry: false,
    })

    // Update zustand store when user data changes
    useEffect(() => {
        if (user) {
            setUser(user)
        }
    }, [user, setUser])

    return (
        <div className="flex h-screen">
            <Sidebar />
            <div className="flex flex-1 flex-col">
                <Header />
                <main className="flex-1 overflow-y-auto bg-base-100 p-4">
                    <Outlet />
                </main>
                <Footer />
            </div>
        </div>
    )
}