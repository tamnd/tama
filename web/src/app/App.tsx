import { lazy, Suspense, useEffect, type ReactNode } from 'react'
import { BrowserRouter, Navigate, Route, Routes, useNavigate } from 'react-router-dom'
import { setUnauthorizedHandler } from '@/api/client'
import { ErrorBoundary } from '@/app/ErrorBoundary'
import { Home } from '@/screens/Home/Home'
import { Login } from '@/screens/Login/Login'
import { AuthProvider, useAuth } from '@/state/auth'

// The gallery is compiled in for dev builds and behind the __GALLERY__ flag;
// production builds drop the chunk entirely. It renders outside the router
// on a plain pathname check, same as before the router landed.
const galleryEnabled = import.meta.env.DEV || __GALLERY__
const Gallery = galleryEnabled ? lazy(() => import('@/dev/Gallery')) : null

export function App() {
  if (Gallery && window.location.pathname === '/dev/gallery') {
    return (
      <ErrorBoundary>
        <Suspense fallback={null}>
          <Gallery />
        </Suspense>
      </ErrorBoundary>
    )
  }
  return (
    <ErrorBoundary>
      <AuthProvider>
        <BrowserRouter>
          <UnauthorizedRedirect />
          <Routes>
            <Route
              path="/"
              element={
                <RequireAuth>
                  <Home />
                </RequireAuth>
              }
            />
            <Route path="/login" element={<Login />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ErrorBoundary>
  )
}

// Points the API client's single 401 interceptor at the router, once for the
// whole app; an expired session lands on /login from any call site.
function UnauthorizedRedirect() {
  const navigate = useNavigate()
  const { signOut } = useAuth()

  useEffect(() => {
    setUnauthorizedHandler(() => {
      signOut()
      if (window.location.pathname !== '/login') {
        void navigate('/login', { replace: true })
      }
    })
    return () => setUnauthorizedHandler(undefined)
  }, [navigate, signOut])

  return null
}

// Gates / on the session: ask the server via me() once, hold the frame blank
// while it answers, and bounce anonymous visitors to /login.
function RequireAuth({ children }: { children: ReactNode }) {
  const { status, refresh } = useAuth()

  useEffect(() => {
    if (status === 'unknown') void refresh()
  }, [status, refresh])

  if (status === 'unknown' || status === 'loading') return null
  if (status === 'anon') return <Navigate to="/login" replace />
  return children
}
