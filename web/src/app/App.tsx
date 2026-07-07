import { lazy, Suspense, useEffect, type ReactNode } from 'react'
import { BrowserRouter, Navigate, Outlet, Route, Routes, useNavigate } from 'react-router-dom'
import { setUnauthorizedHandler } from '@/api/client'
import { ErrorBoundary } from '@/app/ErrorBoundary'
import { Shell } from '@/app/Shell'
import { Card } from '@/components/Card'
import { ToastProvider } from '@/components/Toast'
import { Home } from '@/screens/Home/Home'
import { Login } from '@/screens/Login/Login'
import { AuthProvider, useAuth } from '@/state/auth'

// The gallery is compiled in for dev builds and behind the __GALLERY__ flag;
// production builds drop the chunk entirely. It renders outside the router
// on a plain pathname check, same as before the router landed. The shell
// preview is the document the gallery's resizable iframes load.
const galleryEnabled = import.meta.env.DEV || __GALLERY__
const Gallery = galleryEnabled ? lazy(() => import('@/dev/Gallery')) : null
const ShellPreview = galleryEnabled ? lazy(() => import('@/dev/ShellPreview')) : null

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
  if (ShellPreview && window.location.pathname === '/dev/gallery/shell') {
    return (
      <ErrorBoundary>
        <Suspense fallback={null}>
          <ShellPreview />
        </Suspense>
      </ErrorBoundary>
    )
  }
  return (
    <ErrorBoundary>
      <ToastProvider>
        <AuthProvider>
          <BrowserRouter>
            <UnauthorizedRedirect />
            <Routes>
              <Route
                element={
                  <RequireAuth>
                    <Shell>
                      <Outlet />
                    </Shell>
                  </RequireAuth>
                }
              >
                <Route path="/" element={<Home />} />
                <Route path="/practice" element={<ScreenStub title="Practice" />} />
                <Route path="/leaderboards" element={<ScreenStub title="Leaderboards" />} />
                <Route path="/quests" element={<ScreenStub title="Quests" />} />
                <Route path="/shop" element={<ScreenStub title="Shop" />} />
                <Route path="/profile" element={<ScreenStub title="Profile" />} />
                <Route path="/more" element={<ScreenStub title="More" />} />
              </Route>
              <Route path="/login" element={<Login />} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </BrowserRouter>
        </AuthProvider>
      </ToastProvider>
    </ErrorBoundary>
  )
}

// The rail links to screens that land in later milestones; until then each
// route renders this stub inside the shell.
function ScreenStub({ title }: { title: string }) {
  return (
    <Card>
      <h1>{title}</h1>
      <p>Nothing here yet. This screen lands in a later milestone.</p>
    </Card>
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
