import {createRootRoute, createRoute, createRouter, redirect} from '@tanstack/react-router'
import {RootLayout} from './layouts/RootLayout'
import {DashboardLayout} from './layouts/DashboardLayout'
import {LoginPage} from './pages/LoginPage'
import {DashboardPage} from './pages/DashboardPage'
import {NotFoundPage} from './pages/NotFoundPage'
import {DefinitionsPage} from "@/pages/DefinitionsPage.tsx"
import {EntityPage} from "@/pages/EntityPage.tsx"
import {authApi} from './services/authService'

// Root route - no authentication check here
const rootRoute = createRootRoute({
    component: RootLayout,
})

// Login route - redirect to dashboard if already authenticated
const loginRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/login',
    component: LoginPage,
    beforeLoad: async () => {
        try {
            const user = await authApi.checkAuth()
            if (user) {
                // Already authenticated, go to dashboard
                throw redirect({to: '/', replace: true})
            }
        } catch (error) {
            // Expected error when not authenticated - allow login page to show
            // Don't redirect, just let the login page render
        }
    },
})

// Protected routes - check authentication before loading
const dashboardRoute = createRoute({
    getParentRoute: () => rootRoute,
    id: 'dashboard',
    component: DashboardLayout,
    beforeLoad: async ({location}) => {
        try {
            const user = await authApi.checkAuth()

            // If checkAuth returns null/undefined, redirect to login
            if (!user) {
                throw redirect({
                    to: '/login',
                    search: {
                        redirect: location.href,
                    },
                    replace: true,
                })
            }

            // User authenticated, allow access
            return {user}
        } catch (error: any) {
            // Any error (including 401) means not authenticated
            throw redirect({
                to: '/login',
                search: {
                    redirect: location.href,
                },
                replace: true,
            })
        }
    },
})

const definitionsRoute = createRoute({
    getParentRoute: () => dashboardRoute,
    path: '/definitions',
    component: DefinitionsPage,
})

const entityRoute = createRoute({
    getParentRoute: () => dashboardRoute,
    path: '/entity/$id',
    component: EntityPage,
})

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
    dashboardRoute.addChildren([dashboardIndexRoute, definitionsRoute, entityRoute]),
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