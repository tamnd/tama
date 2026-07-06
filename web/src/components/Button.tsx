import type { ButtonHTMLAttributes, ReactNode } from 'react'
import './Button.css'

export type ButtonVariant = 'primary' | 'blue' | 'secondary' | 'danger' | 'super'
export type ButtonSize = 'small' | 'default' | 'large'

type ButtonProps = {
  variant?: ButtonVariant
  size?: ButtonSize
  fullWidth?: boolean
  /** 20px icon rendered left of the text. */
  icon?: ReactNode
} & ButtonHTMLAttributes<HTMLButtonElement>

// The 3D button: flat face, darker bottom border for depth, press to
// collapse. Text is uppercased by the component, callers pass normal case.
export function Button({
  variant = 'primary',
  size = 'default',
  fullWidth = false,
  icon,
  className,
  children,
  ...rest
}: ButtonProps) {
  const classes = ['tama-button', `tama-button--${variant}`]
  if (size !== 'default') classes.push(`tama-button--${size}`)
  if (fullWidth) classes.push('tama-button--full')
  if (className) classes.push(className)

  return (
    <button type="button" className={classes.join(' ')} {...rest}>
      {icon != null && (
        <span className="tama-button__icon" aria-hidden="true">
          {icon}
        </span>
      )}
      {children}
    </button>
  )
}
