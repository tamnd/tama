import type { ComponentPropsWithRef, ReactNode } from 'react'
import './ChoiceCard.css'

type ChoiceCardProps = {
  /** Radio for pick-one exercises, checkbox for multi-select. */
  type?: 'radio' | 'checkbox'
  /** Keyboard shortcut digit shown bottom-right. */
  badge?: number
  children: ReactNode
} & Omit<ComponentPropsWithRef<'input'>, 'type' | 'children'>

// A selectable card over a native radio or checkbox. The input stays in
// the tree visually hidden; checked and focus states style the card
// through :has(), so the semantics are entirely the browser's.
export function ChoiceCard({
  type = 'radio',
  badge,
  disabled,
  className,
  children,
  ...rest
}: ChoiceCardProps) {
  const classes = ['tama-choice']
  if (badge != null) classes.push('tama-choice--badged')
  if (disabled) classes.push('tama-choice--disabled')
  if (className) classes.push(className)

  return (
    <label className={classes.join(' ')}>
      <input
        type={type}
        className="tama-choice__input visually-hidden"
        disabled={disabled}
        {...rest}
      />
      <span className="tama-choice__text">{children}</span>
      {badge != null && (
        <span className="tama-choice__badge" aria-hidden="true">
          {badge}
        </span>
      )}
    </label>
  )
}
