import {Outlet} from '@tanstack/react-router'
import {Header} from '../components/layout/Header'
import {Footer} from '../components/layout/Footer'
import {Sidebar} from '../components/layout/Sidebar'

export function DashboardLayout() {
    return (
        <div className="flex h-screen flex-col">
            <Header/>
            <div className="flex flex-1 overflow-hidden">
                <Sidebar/>
                <main className="flex-1 overflow-y-auto bg-base-200 p-6">
                    <Outlet/>
                </main>
            </div>
            <Footer/>
        </div>
    )
}
