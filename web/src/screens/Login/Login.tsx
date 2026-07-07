import { useState, type FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { ApiError, login, register } from '@/api/client'
import { Button } from '@/components/Button'
import { Card } from '@/components/Card'
import { Field } from '@/components/Field'
import { TextInput } from '@/components/TextInput'
import { useAuth } from '@/state/auth'
import './Login.css'

// The auth screen: log in and sign up as two tabs over one form. Nothing is
// stored client-side; the session cookie the server sets is the session.

type Mode = 'login' | 'register'

const COPY = {
  login: { tab: 'Log in', title: 'Welcome back', cta: 'Log in' },
  register: { tab: 'Sign up', title: 'Create your profile', cta: 'Create account' },
} as const

const MODES: Mode[] = ['login', 'register']

export function Login() {
  const [mode, setMode] = useState<Mode>('login')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const navigate = useNavigate()
  const { signIn } = useAuth()

  function switchMode(next: Mode) {
    setMode(next)
    setError(null)
  }

  async function submit() {
    setBusy(true)
    setError(null)
    try {
      const call = mode === 'login' ? login : register
      const user = await call({ username, password })
      signIn(user)
      await navigate('/', { replace: true })
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'something went wrong, try again')
    } finally {
      setBusy(false)
    }
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    void submit()
  }

  return (
    <main className="login">
      <img className="login__mascot" src="/tama.svg" alt="" />
      <Card className="login__card">
        <div className="login__tabs">
          {MODES.map((m) => (
            <Button
              key={m}
              size="small"
              variant={mode === m ? 'blue' : 'secondary'}
              aria-pressed={mode === m}
              onClick={() => switchMode(m)}
            >
              {COPY[m].tab}
            </Button>
          ))}
        </div>
        <h1 className="login__title">{COPY[mode].title}</h1>
        <form className="login__form" onSubmit={handleSubmit}>
          <Field
            label="Username"
            hint={mode === 'register' ? '3-24 characters: a-z, 0-9, _' : undefined}
          >
            <TextInput
              name="username"
              autoComplete="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </Field>
          <Field label="Password" hint={mode === 'register' ? 'At least 8 characters' : undefined}>
            <TextInput
              type="password"
              name="password"
              autoComplete={mode === 'register' ? 'new-password' : 'current-password'}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </Field>
          {error && (
            <p className="login__error" role="alert">
              {error}
            </p>
          )}
          <Button type="submit" variant="primary" size="large" fullWidth disabled={busy}>
            {COPY[mode].cta}
          </Button>
        </form>
      </Card>
    </main>
  )
}
