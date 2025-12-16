import {Bars3Icon} from '@heroicons/react/24/outline'
import {UserMenu} from './UserMenu.tsx'
import {useAuthStore} from '@/stores/authStore.ts'

export function Header() {
    const {user} = useAuthStore()

    return (
        <header className="navbar bg-base-100 shadow-lg">
            <div className="flex-1">
                <label htmlFor="sidebar-drawer" className="btn btn-ghost btn-circle lg:hidden">
                    <Bars3Icon className="h-6 w-6"/>
                </label>
                <a href="/app/ui/public" className="btn btn-ghost text-xl font-bold">
                    Fritz
                </a>
            </div>
            <div className="flex-none">
                <div className="flex items-center gap-4">
                    <div className="form-control">
                        <input
                            type="text"
                            placeholder="Search..."
                            className="input input-bordered w-24 md:w-auto"
                        />
                    </div>
                    {user && <UserMenu user={user}/>}
                </div>
            </div>
        </header>
    )
}
