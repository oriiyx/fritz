import {createRootRoute, createRoute, createRouter} from '@tanstack/react-router'
import {RootLayout} from './layouts/RootLayout'
import {DashboardLayout} from './layouts/DashboardLayout'
import {LoginPage} from './pages/LoginPage'
import {DashboardPage} from './pages/DashboardPage'
import {NotFoundPage} from './pages/NotFoundPage'
import {EntitiesPage} from "@/pages/EntitiesPage.tsx";

// Root route
const rootRoute = createRootRoute({
    component: RootLayout,
})

// Guest routes (no authentication required)
const loginRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/login',
    component: LoginPage,
})

// Dashboard layout route (no path specified - it's just a layout wrapper)
const dashboardRoute = createRoute({
    getParentRoute: () => rootRoute,
    id: 'dashboard',  // Use id instead of path for layout-only routes
    component: DashboardLayout,
})

const entitiesRoute = createRoute({
    getParentRoute: () => dashboardRoute,
    path: '/entities',
    component: EntitiesPage,
})

// Dashboard index route - this is what actually matches '/'
const dashboardIndexRoute = createRoute({
    getParentRoute: () => dashboardRoute,
    path: '/',
    component: DashboardPage,
})

// 404 route
const notFoundRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '*',
    component: NotFoundPage,
})

// Create route tree
const routeTree = rootRoute.addChildren([
    loginRoute,
    dashboardRoute.addChildren([dashboardIndexRoute, entitiesRoute]),
    notFoundRoute,
])

// Create router instance
export const router = createRouter({
    routeTree,
    defaultPreload: 'intent',
})

// Register router for type safety
declare module '@tanstack/react-router' {
    interface Register {
        router: typeof router
    }
}