import { useEffect, useRef, type MouseEvent, type ReactNode } from 'react'
import './Modal.css'

type ModalProps = {
  open: boolean
  /** Fired when Esc or a backdrop click asks to close; the caller drops `open`. */
  onClose: () => void
  /** Esc and backdrop clicks stop closing; only the footer buttons get out. */
  blocking?: boolean
  title?: string
  /**
   * Button children, stacked full width. Pass the primary FIRST: it lands
   * first in tab order after the content, and the footer's column-reverse
   * shows it last, at the bottom, where the product puts it.
   */
  footer?: ReactNode
  className?: string
  children: ReactNode
}

// The dialog is native, so top-layer stacking, focus trapping, and inert
// background come for free from showModal(). We only sync the open prop,
// veto cancel so React state stays the one source of truth, and route
// backdrop clicks. Depth is the 2px swan border, never a shadow.
export function Modal({
  open,
  onClose,
  blocking = false,
  title,
  footer,
  className,
  children,
}: ModalProps) {
  const ref = useRef<HTMLDialogElement>(null)

  useEffect(() => {
    const dialog = ref.current
    if (!dialog) return
    if (open && !dialog.open) dialog.showModal()
    else if (!open && dialog.open) dialog.close()
  }, [open])

  // Esc fires cancel. Always veto the native close, then ask the owner,
  // so a blocking modal cannot be escaped and a normal one closes through
  // the same setState path as everything else.
  useEffect(() => {
    const dialog = ref.current
    if (!dialog) return
    function onCancel(event: Event) {
      event.preventDefault()
      if (!blocking) onClose()
    }
    dialog.addEventListener('cancel', onCancel)
    return () => dialog.removeEventListener('cancel', onCancel)
  }, [blocking, onClose])

  // Clicks land on the dialog element itself only when they hit the
  // backdrop; the inner wrapper swallows everything over the surface.
  function onBackdropClick(event: MouseEvent<HTMLDialogElement>) {
    if (event.target === ref.current && !blocking) onClose()
  }

  const classes = ['tama-modal']
  if (className) classes.push(className)

  return (
    <dialog ref={ref} className={classes.join(' ')} onClick={onBackdropClick}>
      <div className="tama-modal__inner">
        {title && <h2 className="tama-modal__title">{title}</h2>}
        <div className="tama-modal__content">{children}</div>
        {footer && <div className="tama-modal__footer">{footer}</div>}
      </div>
    </dialog>
  )
}
