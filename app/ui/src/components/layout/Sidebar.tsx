import {Link} from '@tanstack/react-router'
import {Cog6ToothIcon, CubeIcon, DocumentTextIcon, HomeIcon, UsersIcon,} from '@heroicons/react/24/outline'
import {ComponentType} from "react";

interface NavItem {
    name: string
    path: string
    icon: ComponentType<{ className?: string }>
}

const navItems: NavItem[] = [
    {name: 'Dashboard', path: '/', icon: HomeIcon},
    {name: 'Entities', path: '/entities', icon: CubeIcon},
    {name: 'Definitions', path: '/definitions', icon: DocumentTextIcon},
    {name: 'Users', path: '/users', icon: UsersIcon},
    {name: 'Settings', path: '/settings', icon: Cog6ToothIcon},
]

export function Sidebar() {
    return (
        <>
            {/* Mobile drawer */}
            <div className="drawer lg:drawer-open">
                <input id="sidebar-drawer" type="checkbox" className="drawer-toggle"/>
                <div className="drawer-side z-10">
                    <label htmlFor="sidebar-drawer" className="drawer-overlay"></label>
                    <aside className="menu min-h-full w-64 bg-base-100 p-4 text-base-content shadow-xl">
                        <div className="mb-4">
                            <h2 className="px-4 text-lg font-semibold">Navigation</h2>
                        </div>
                        <ul>
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
                        </ul>
                    </aside>
                </div>
            </div>
        </>
    )
}
