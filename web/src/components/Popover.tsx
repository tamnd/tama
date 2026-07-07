import { useEffect, useId, useRef, useState, type ReactNode, type RefObject } from 'react'
import './Popover.css'

export type PopoverVariant = 'neutral' | 'green'

/** Gap between the anchor's bottom edge and the popover, tail included. */
const ANCHOR_GAP = 12

type PopoverProps = {
  /** The path node (or any element) the centered tail points at. */
  anchorRef: RefObject<HTMLElement | null>
  /** green is the active-node look: feather fill, snow text. */
  variant?: PopoverVariant
  className?: string
  children: ReactNode
}

const supportsAnchor =
  typeof CSS !== 'undefined' &&
  typeof CSS.supports === 'function' &&
  CSS.supports('anchor-name: --tama')

// The lesson popover shell that hangs under a path node. It is plain
// content, never focused and never focusable, so opening one steals
// nothing from the keyboard. Positioning prefers CSS anchor positioning;
// where that is missing we measure the anchor rect and pin the popover
// with fixed coordinates, re-measuring on scroll and resize.
export function Popover({ anchorRef, variant = 'neutral', className, children }: PopoverProps) {
  const reactId = useId()
  const ref = useRef<HTMLDivElement>(null)
  const [measured, setMeasured] = useState<{ top: number; left: number } | null>(null)

  useEffect(() => {
    const anchor = anchorRef.current
    if (!anchor) return

    if (supportsAnchor) {
      const name = `--tama-popover-${reactId.replace(/[^a-zA-Z0-9_-]/g, '')}`
      anchor.style.setProperty('anchor-name', name)
      ref.current?.style.setProperty('position-anchor', name)
      return () => anchor.style.removeProperty('anchor-name')
    }

    function measure() {
      const rect = anchorRef.current?.getBoundingClientRect()
      if (rect) setMeasured({ top: rect.bottom + ANCHOR_GAP, left: rect.left + rect.width / 2 })
    }
    measure()
    window.addEventListener('resize', measure)
    window.addEventListener('scroll', measure, true)
    return () => {
      window.removeEventListener('resize', measure)
      window.removeEventListener('scroll', measure, true)
    }
  }, [anchorRef, reactId])

  const classes = ['tama-popover', `tama-popover--${variant}`]
  classes.push(supportsAnchor ? 'tama-popover--anchored' : 'tama-popover--measured')
  if (className) classes.push(className)

  return (
    <div
      ref={ref}
      className={classes.join(' ')}
      style={measured ? { top: `${measured.top}px`, left: `${measured.left}px` } : undefined}
    >
      <span className="tama-popover__tail" aria-hidden="true" />
      {children}
    </div>
  )
}
