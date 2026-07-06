import {
  createContext,
  useCallback,
  useContext,
  useRef,
  useState,
  type ReactNode,
} from 'react'
import { XIcon } from './icons/x'
import './Toast.css'

const AUTO_DISMISS_MS = 4000
const MAX_TOASTS = 3

export type ShowToast = (message: string, options?: { icon?: ReactNode }) => void

type ToastItem = {
  id: number
  message: string
  icon?: ReactNode
}

const ToastContext = createContext<ShowToast | null>(null)

export function useToast(): ShowToast {
  const show = useContext(ToastContext)
  if (!show) throw new Error('useToast needs a ToastProvider above it')
  return show
}

// Owns the queue and renders the top-center stack. Toasts announce through
// role="status" and never take focus; each one leaves on its own after 4s
// or sooner through the dismiss button. A fourth arrival evicts the oldest
// so the stack never grows past three.
export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([])
  const nextId = useRef(1)

  const dismiss = useCallback((id: number) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id))
  }, [])

  const show = useCallback<ShowToast>(
    (message, options) => {
      const id = nextId.current++
      setToasts((prev) => [...prev, { id, message, icon: options?.icon }].slice(-MAX_TOASTS))
      window.setTimeout(() => dismiss(id), AUTO_DISMISS_MS)
    },
    [dismiss],
  )

  return (
    <ToastContext.Provider value={show}>
      {children}
      <div className="tama-toast-viewport">
        {toasts.map((toast) => (
          <Toast key={toast.id} icon={toast.icon} onDismiss={() => dismiss(toast.id)}>
            {toast.message}
          </Toast>
        ))}
      </div>
    </ToastContext.Provider>
  )
}

type ToastProps = {
  icon?: ReactNode
  onDismiss: () => void
  children: ReactNode
}

// One toast chip: icon slot, message, and a visible dismiss button so
// nobody has to wait out the timer.
export function Toast({ icon, onDismiss, children }: ToastProps) {
  return (
    <div className="tama-toast" role="status">
      {icon != null && (
        <span className="tama-toast__icon" aria-hidden="true">
          {icon}
        </span>
      )}
      <span className="tama-toast__message">{children}</span>
      <button type="button" className="tama-toast__dismiss" aria-label="Dismiss" onClick={onDismiss}>
        <XIcon size={16} />
      </button>
    </div>
  )
}
