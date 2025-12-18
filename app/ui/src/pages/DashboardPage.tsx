import {useAuthStore} from '../stores/authStore'
import {ChartBarIcon, CubeIcon, DocumentTextIcon, UsersIcon,} from '@heroicons/react/24/outline'
import {Stat, Stats} from "@/components/Stats.tsx";
import {Card, CardBody, CardTitle} from "@/components/Card.tsx";
import {Button} from "@/components/Button.tsx";
import {useNavigate} from "@tanstack/react-router";

export function DashboardPage() {
    const {user} = useAuthStore()
    const navigate = useNavigate()

    const statsData = [
        {name: 'Total Entities', value: '0', icon: CubeIcon, color: 'text-primary'},
        {name: 'Definitions', value: '0', icon: DocumentTextIcon, color: 'text-secondary'},
        {name: 'Active Users', value: '1', icon: UsersIcon, color: 'text-accent'},
        {name: 'Activities', value: '0', icon: ChartBarIcon, color: 'text-info'},
    ]

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-3xl font-bold">Dashboard</h1>
                <p className="mt-2 text-base-content/70">Welcome back, {user?.display_name || 'User'}!</p>
            </div>

            {/* Stats Grid */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {statsData.map((stat) => {
                    const Icon = stat.icon
                    return (
                        <Stats key={stat.name}>
                            <Stat
                                title={stat.name}
                                value={stat.value}
                                figure={<Icon className={`h-8 w-8 ${stat.color}`}/>}
                            />
                        </Stats>
                    )
                })}
            </div>

            {/* Quick Actions */}
            <Card>
                <CardBody>
                    <CardTitle>Quick Actions</CardTitle>
                    <div className="mt-4 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                        <Button variant="primary" outline iconLeft={<CubeIcon className="h-5 w-5"/>}
                                onClick={() => navigate({to: '/entities'})}>
                            Create Entity
                        </Button>
                        <Button variant="secondary" outline iconLeft={<DocumentTextIcon className="h-5 w-5"/>}>
                            New Definition
                        </Button>
                        <Button variant="accent" outline iconLeft={<UsersIcon className="h-5 w-5"/>}>
                            Manage Users
                        </Button>
                    </div>
                </CardBody>
            </Card>

            {/* Recent Activity */}
            <Card>
                <CardBody>
                    <CardTitle>Recent Activity</CardTitle>
                    <div className="py-8 text-center text-base-content/50">No recent activity to display</div>
                </CardBody>
            </Card>
        </div>
    )
}
