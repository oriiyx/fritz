import {
    ArrowRightOnRectangleIcon,
    Bars3Icon,
    Cog6ToothIcon,
    DocumentTextIcon,
    HomeIcon,
    UserCircleIcon,
    UsersIcon
} from '@heroicons/react/24/outline'
import {Link, useNavigate} from '@tanstack/react-router'
import {useAuthStore} from '@/stores/authStore.ts'
import {ComponentType} from 'react'
import {useMutation} from '@tanstack/react-query'
import {authApi} from '@/services/authService.ts'

interface NavItem {
    name: string
    path: string
    icon: ComponentType<{ className?: string }>
}

const navItems: NavItem[] = [
    {name: 'Dashboard', path: '/', icon: HomeIcon},
    {name: 'Definitions', path: '/definitions', icon: DocumentTextIcon},
    {name: 'Users', path: '/users', icon: UsersIcon},
    {name: 'Settings', path: '/settings', icon: Cog6ToothIcon},
]

export function Header() {
    const {user, logout: logoutStore} = useAuthStore()
    const navigate = useNavigate()

    const logoutMutation = useMutation({
        mutationFn: () => authApi.logout('google'),
        onSuccess: () => {
            logoutStore()
            navigate({to: '/login'})
        },
    })

    const handleLogout = () => {
        logoutMutation.mutate()
    }

    return (
        <header className="navbar bg-base-300/20">
            <div className="flex-1">
                <a href="/" className="btn btn-ghost text-xl font-bold">
                    Fritz
                </a>
            </div>
            <div className="flex-none">
                <div className="flex items-center gap-4">
                    <div className="form-control">
                        <input
                            type="text"
                            placeholder="Search..."
                            className="input input-bordered w-auto"
                        />
                    </div>
                    {user && (
                        <div className="dropdown dropdown-end">
                            <div tabIndex={0} role="button" className="btn btn-ghost btn-circle">
                                <Bars3Icon className="h-6 w-6"/>
                            </div>
                            <ul tabIndex={0}
                                className="menu dropdown-content z-[1] mt-3 w-52 rounded-box bg-base-100 p-2 shadow">
                                <li className="menu-title">
                                    <span>{user.display_name}</span>
                                    <span className="text-xs font-normal">{user.email}</span>
                                </li>
                                <div className="divider my-0"></div>
                                {navItems.map((item) => {
                                    const Icon = item.icon
                                    return (
                                        <li key={item.path}>
                                            <Link
                                                to={item.path}
                                                activeProps={{
                                                    className: 'active',
                                                }}
                                                className="flex items-center gap-3"
                                            >
                                                <Icon className="h-5 w-5"/>
                                                <span>{item.name}</span>
                                            </Link>
                                        </li>
                                    )
                                })}
                                <div className="divider my-0"></div>
                                <li>
                                    <a href="/profile">
                                        <UserCircleIcon className="h-5 w-5"/>
                                        Profile
                                    </a>
                                </li>
                                <li>
                                    <button onClick={handleLogout} disabled={logoutMutation.isPending}>
                                        <ArrowRightOnRectangleIcon className="h-5 w-5"/>
                                        {logoutMutation.isPending ? 'Logging out...' : 'Logout'}
                                    </button>
                                </li>
                            </ul>
                        </div>
                    )}
                </div>
            </div>
        </header>
    )
}