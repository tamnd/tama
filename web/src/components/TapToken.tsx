import type { ComponentPropsWithRef } from 'react'
import './TapToken.css'

type TapTokenProps = {
  /** The spent placeholder left in the bank: polar slab, non-interactive. */
  hollow?: boolean
} & ComponentPropsWithRef<'button'>

// The word-bank chip: a small secondary-button face that presses down.
// Hollow tokens keep the bank layout stable after their word moves up.
export function TapToken({ hollow = false, disabled, className, children, ...rest }: TapTokenProps) {
  const classes = ['tama-tap-token']
  if (hollow) classes.push('tama-tap-token--hollow')
  if (className) classes.push(className)

  return (
    <button type="button" className={classes.join(' ')} disabled={hollow || disabled} {...rest}>
      {children}
    </button>
  )
}
