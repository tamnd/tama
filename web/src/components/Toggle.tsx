import type { ComponentPropsWithRef, ReactNode } from 'react'
import './Toggle.css'

type ToggleProps = {
  /** Label text after the switch; pass aria-label instead when omitted. */
  children?: ReactNode
} & Omit<ComponentPropsWithRef<'input'>, 'type' | 'children'>

// The settings switch: a native checkbox under a 32x20 pill track. The
// knob travels over dur-fast and the track goes feather green when on.
export function Toggle({ className, disabled, children, ...rest }: ToggleProps) {
  const classes = ['tama-toggle']
  if (disabled) classes.push('tama-toggle--disabled')
  if (className) classes.push(className)

  return (
    <label className={classes.join(' ')}>
      <input
        type="checkbox"
        className="tama-toggle__input visually-hidden"
        disabled={disabled}
        {...rest}
      />
      <span className="tama-toggle__track" aria-hidden="true">
        <span className="tama-toggle__knob" />
      </span>
      {children != null && <span className="tama-toggle__text">{children}</span>}
    </label>
  )
}
