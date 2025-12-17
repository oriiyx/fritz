import {Outlet} from '@tanstack/react-router'
import {Header} from '../components/layout/Header'
import {Footer} from '../components/layout/Footer'
import {Sidebar} from '../components/layout/Sidebar'

export function DashboardLayout() {
    return (
        <div className="drawer lg:drawer-open">
            <input id="sidebar-drawer" type="checkbox" className="drawer-toggle"/>
            <div className="drawer-content flex flex-col">
                <Header/>
                <main className="flex-1 overflow-y-auto bg-base-200 p-6">
                    <Outlet/>
                </main>
                <Footer/>
            </div>
            <div className="drawer-side">
                <label htmlFor="sidebar-drawer" aria-label="close sidebar" className="drawer-overlay"></label>
                <Sidebar/>
            </div>
        </div>
    )
}