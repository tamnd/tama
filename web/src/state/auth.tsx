import { createContext, useCallback, useContext, useMemo, useState, type ReactNode } from 'react'
import { me } from '@/api/client'
import type { User } from '@/api/types'

// Session state, hand-rolled on React context. See README.md for why there
// is no store library behind this.

export type AuthStatus = 'unknown' | 'loading' | 'authed' | 'anon'

type AuthState = {
  user: User | null
  status: AuthStatus
  /** Asks the server who owns the cookie; drives the / route guard. */
  refresh: () => Promise<void>
  /** Records the user that login or register returned. */
  signIn: (user: User) => void
  /** Drops the client-side copy; the cookie dies with the logout call. */
  signOut: () => void
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [status, setStatus] = useState<AuthStatus>('unknown')

  const refresh = useCallback(async () => {
    setStatus('loading')
    try {
      setUser(await me())
      setStatus('authed')
    } catch {
      setUser(null)
      setStatus('anon')
    }
  }, [])

  const signIn = useCallback((next: User) => {
    setUser(next)
    setStatus('authed')
  }, [])

  const signOut = useCallback(() => {
    setUser(null)
    setStatus('anon')
  }, [])

  const value = useMemo(
    () => ({ user, status, refresh, signIn, signOut }),
    [user, status, refresh, signIn, signOut],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must run under AuthProvider')
  return ctx
}
