import type { ComponentPropsWithRef } from 'react'
import './TextInput.css'

type TextInputProps = ComponentPropsWithRef<'input'>

// The single-line text field. Error wiring (aria-invalid, the message id)
// comes from the Field wrapper; the CSS reacts to aria-invalid.
export function TextInput({ className, ...rest }: TextInputProps) {
  const classes = ['tama-text-input']
  if (className) classes.push(className)

  return <input className={classes.join(' ')} {...rest} />
}
