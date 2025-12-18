import {Outlet} from '@tanstack/react-router'

export function RootLayout() {
    // Remove all authentication checking logic from here
    // Let the router handle authentication via beforeLoad guards
    return <Outlet/>
}