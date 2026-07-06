import type { ReactNode } from 'react'
import './Card.css'

export type CardVariant = 'default' | 'well' | 'clickable'

type CardProps = {
  variant?: CardVariant
  onClick?: () => void
  className?: string
  children: ReactNode
}

// Surface container. The well variant is the inset polar fill for stat
// wells and list rows; clickable gets the 3D press treatment and renders
// as a real button.
export function Card({ variant = 'default', onClick, className, children }: CardProps) {
  const classes = ['tama-card']
  if (variant !== 'default') classes.push(`tama-card--${variant}`)
  if (className) classes.push(className)

  if (variant === 'clickable') {
    return (
      <button type="button" className={classes.join(' ')} onClick={onClick}>
        {children}
      </button>
    )
  }
  return <div className={classes.join(' ')}>{children}</div>
}

// Section header inside cards: the caps label with a 16px bottom margin.
export function CardHeader({ children }: { children: ReactNode }) {
  return <h3 className="label-caps tama-card__header">{children}</h3>
}

// A 2px rule, full bleed across the card padding.
export function Divider() {
  return <hr className="tama-card__divider" />
}
