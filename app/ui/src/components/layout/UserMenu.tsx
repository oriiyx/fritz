import {useState} from 'react'
import {ArrowRightOnRectangleIcon, UserCircleIcon} from '@heroicons/react/24/outline'
import {useNavigate} from '@tanstack/react-router'
import {useMutation} from '@tanstack/react-query'
import {useAuthStore, User} from '@/stores/authStore.ts'
import {authApi} from '@/services/authService.ts'

interface UserMenuProps {
    user: User
}

export function UserMenu({user}: UserMenuProps) {
    const navigate = useNavigate()
    const {logout: logoutStore} = useAuthStore()
    const [isOpen, setIsOpen] = useState(false)

    const logoutMutation = useMutation({
        mutationFn: () => authApi.logout('google'), // Default to google, adjust as needed
        onSuccess: () => {
            logoutStore()
            navigate({to: '/login'})
        },
    })

    const handleLogout = () => {
        logoutMutation.mutate()
    }

    return (
        <div className="dropdown dropdown-end">
            <div
                tabIndex={0}
                role="button"
                className="avatar btn btn-circle btn-ghost"
                onClick={() => setIsOpen(!isOpen)}
            >
                {user.avatar_url ? (
                    <div className="w-10 rounded-full">
                        <img alt={user.full_name} src={user.avatar_url}/>
                    </div>
                ) : (
                    <UserCircleIcon className="h-8 w-8"/>
                )}
            </div>
            {isOpen && (
                <ul
                    tabIndex={0}
                    className="menu dropdown-content menu-sm z-[1] mt-3 w-52 rounded-box bg-base-100 p-2 shadow"
                >
                    <li className="menu-title">
                        <span>{user.full_name}</span>
                        <span className="text-xs font-normal">{user.email}</span>
                    </li>
                    <li>
                        <a href="/profile">
                            <UserCircleIcon className="h-4 w-4"/>
                            Profile
                        </a>
                    </li>
                    <li>
                        <button onClick={handleLogout} disabled={logoutMutation.isPending}>
                            <ArrowRightOnRectangleIcon className="h-4 w-4"/>
                            {logoutMutation.isPending ? 'Logging out...' : 'Logout'}
                        </button>
                    </li>
                </ul>
            )}
        </div>
    )
}
