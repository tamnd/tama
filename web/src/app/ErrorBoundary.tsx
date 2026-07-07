import { Component, type ErrorInfo, type ReactNode } from 'react'
import { Button } from '@/components/Button'
import './ErrorBoundary.css'

type Props = { children: ReactNode }
type State = { failed: boolean }

// The one render-error net for the whole app: log the error, offer a reload.
export class ErrorBoundary extends Component<Props, State> {
  state: State = { failed: false }

  static getDerivedStateFromError(): State {
    return { failed: true }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('render error', error, info.componentStack)
  }

  render() {
    if (!this.state.failed) return this.props.children
    return (
      <main className="error-screen">
        <h1>Something broke</h1>
        <p>The screen hit an error it could not recover from.</p>
        <Button variant="primary" onClick={() => window.location.reload()}>
          Reload
        </Button>
      </main>
    )
  }
}
