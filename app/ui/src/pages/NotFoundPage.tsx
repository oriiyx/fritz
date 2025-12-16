import {Link} from '@tanstack/react-router'
import {HomeIcon} from '@heroicons/react/24/outline'

export function NotFoundPage() {
    return (
        <div className="hero min-h-screen bg-base-200">
            <div className="hero-content text-center">
                <div className="max-w-md">
                    <h1 className="text-9xl font-bold text-primary">404</h1>
                    <h2 className="text-3xl font-bold mt-4">Page Not Found</h2>
                    <p className="py-6">
                        Oops! The page you're looking for doesn't exist. It might have been moved or deleted.
                    </p>
                    <Link to="/" className="btn btn-primary gap-2">
                        <HomeIcon className="h-5 w-5"/>
                        Back to Dashboard
                    </Link>
                </div>
            </div>
        </div>
    )
}
