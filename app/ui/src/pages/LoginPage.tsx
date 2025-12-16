import {authApi} from '../services/authService'
import {Card, CardBody} from "@/components/Card.tsx";
import {Button} from "@/components/Button.tsx";
import {GitHubIcon, GoogleIcon} from "@/components/CustomIcons.tsx";

export function LoginPage() {
    const handleGoogleLogin = () => {
        authApi.login('google')
    }

    const handleGithubLogin = () => {
        authApi.login('github')
    }

    return (
        <div className="hero min-h-screen bg-base-200">
            <div className="hero-content flex-col gap-20 lg:flex-row">
                <div className="max-w-xs text-center lg:text-left">
                    <h1 className="text-5xl font-bold">Welcome to Fritz</h1>
                    <p className="py-6">
                        Sign in to access your Product Information Management system. Use your Google or GitHub
                        account to get started.
                    </p>
                </div>

                <Card className="w-full max-w-sm flex-shrink-0">
                    <CardBody>
                        <div className="text-center">
                            <h2 className="text-2xl font-bold">Sign In</h2>
                            <p className="mt-2 text-sm text-base-content/70">
                                Choose your preferred authentication method
                            </p>
                        </div>

                        <div className="form-control mt-6 gap-4">
                            <Button variant="primary" outline iconLeft={<GoogleIcon/>} onClick={handleGoogleLogin}>
                                Continue with Google
                            </Button>

                            <Button outline iconLeft={<GitHubIcon/>} onClick={handleGithubLogin}>
                                Continue with GitHub
                            </Button>
                        </div>

                        <div className="divider text-xs">OAuth Authentication</div>

                        <p className="text-center text-xs text-base-content/50">
                            By signing in, you agree to our Terms of Service
                        </p>
                    </CardBody>
                </Card>
            </div>
        </div>
    )
}
