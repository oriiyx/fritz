import {useState} from 'react'
import type {AnyFieldApi} from '@tanstack/react-form'
import {useForm} from '@tanstack/react-form'
import {useMutation} from '@tanstack/react-query'
import {useNavigate} from '@tanstack/react-router'
import {authApi} from '../services/authService'
import {useAuthStore} from '../stores/authStore'
import {Card, CardBody} from "@/components/Card.tsx"
import {Button} from "@/components/Button.tsx"
import {Input} from "@/components/Input.tsx"
import {Alert} from "@/components/Alert.tsx"
import {GitHubIcon, GoogleIcon} from "@/components/CustomIcons.tsx"
import {ExclamationCircleIcon} from '@heroicons/react/24/outline'

function FieldError({field}: { field: AnyFieldApi }) {
    return (
        <>
            {field.state.meta.isTouched && field.state.meta.errors.length > 0 && (
                <span className="text-error text-sm mt-1">{field.state.meta.errors.join(', ')}</span>
            )}
        </>
    )
}

export function LoginPage() {
    const navigate = useNavigate()
    const {setUser} = useAuthStore()
    const [loginError, setLoginError] = useState<string | null>(null)

    const loginMutation = useMutation({
        mutationFn: ({email, password}: { email: string; password: string }) =>
            authApi.loginWithPassword(email, password),
        onSuccess: (user) => {
            setUser(user)
            navigate({to: '/'})
        },
        onError: (error: any) => {
            console.error('Login error:', error)
            if (error?.response?.status === 401) {
                setLoginError('Invalid email or password. Please try again.')
            } else {
                setLoginError('An error occurred during login. Please try again.')
            }
        },
    })

    const form = useForm({
        defaultValues: {
            email: '',
            password: '',
        },
        onSubmit: async ({value}) => {
            setLoginError(null)
            loginMutation.mutate(value)
        },
    })

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
                        Sign in to access your Product Information Management system. Use your email and password or
                        your Google or GitHub account to get started.
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

                        {loginError && (
                            <Alert variant="error" icon={<ExclamationCircleIcon className="h-5 w-5"/>}>
                                {loginError}
                            </Alert>
                        )}

                        {/* Email/Password Form */}
                        <form
                            onSubmit={(e) => {
                                e.preventDefault()
                                e.stopPropagation()
                                form.handleSubmit().then()
                            }}
                            className="space-y-4"
                        >
                            <form.Field
                                name="email"
                                validators={{
                                    onChange: ({value}) => {
                                        if (!value) return 'Email is required'
                                        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
                                            return 'Please enter a valid email address'
                                        }
                                        return undefined
                                    },
                                }}
                                children={(field) => (
                                    <div>
                                        <Input
                                            type="email"
                                            label="Email"
                                            id={field.name}
                                            name={field.name}
                                            value={field.state.value}
                                            onBlur={field.handleBlur}
                                            onChange={(e) => field.handleChange(e.target.value)}
                                            fullWidth
                                            variant="bordered"
                                        />
                                        <FieldError field={field}/>
                                    </div>
                                )}
                            />

                            <form.Field
                                name="password"
                                validators={{
                                    onChange: ({value}) => {
                                        if (!value) return 'Password is required'
                                        if (value.length < 8) {
                                            return 'Password must be at least 8 characters'
                                        }
                                        return undefined
                                    },
                                }}
                                children={(field) => (
                                    <div>
                                        <Input
                                            type="password"
                                            label="Password"
                                            id={field.name}
                                            name={field.name}
                                            value={field.state.value}
                                            onBlur={field.handleBlur}
                                            onChange={(e) => field.handleChange(e.target.value)}
                                            fullWidth
                                            variant="bordered"
                                        />
                                        <FieldError field={field}/>
                                    </div>
                                )}
                            />

                            <form.Subscribe
                                selector={(state) => [state.canSubmit, state.isSubmitting]}
                                children={([canSubmit, isSubmitting]) => (
                                    <Button
                                        type="submit"
                                        variant="primary"
                                        block
                                        disabled={!canSubmit || loginMutation.isPending}
                                        loading={isSubmitting || loginMutation.isPending}
                                    >
                                        Sign In
                                    </Button>
                                )}
                            />
                        </form>

                        <div className="divider text-xs">OR</div>

                        {/* OAuth Options */}
                        <div className="form-control gap-4 flex flex-col">
                            <Button outline iconLeft={<GoogleIcon/>} onClick={handleGoogleLogin}>
                                Continue with Google
                            </Button>

                            <Button outline iconLeft={<GitHubIcon/>} onClick={handleGithubLogin}>
                                Continue with GitHub
                            </Button>
                        </div>

                        <p className="text-center text-xs text-base-content/50">
                            By signing in, you agree to our Terms of Service
                        </p>
                    </CardBody>
                </Card>
            </div>
        </div>
    )
}