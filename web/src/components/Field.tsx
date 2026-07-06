import { cloneElement, useId, type ReactElement } from 'react'
import './Field.css'

type FieldProps = {
  label: string
  hint?: string
  error?: string
  /** A single form control; Field injects id, aria-describedby, aria-invalid. */
  children: ReactElement<Record<string, unknown>>
}

// Label, control, hint, and error with the screen reader wiring written
// once: the control gets the label's target id, the hint and error ids
// land in aria-describedby, and an error also sets aria-invalid.
export function Field({ label, hint, error, children }: FieldProps) {
  const id = useId()
  const hintId = hint ? `${id}-hint` : undefined
  const errorId = error ? `${id}-error` : undefined
  const describedBy = [hintId, errorId].filter(Boolean).join(' ')

  const control = cloneElement(children, {
    id,
    'aria-describedby': describedBy || undefined,
    'aria-invalid': error ? true : undefined,
  })

  return (
    <div className="tama-field">
      <label htmlFor={id} className="label-caps tama-field__label">
        {label}
      </label>
      {control}
      {hint && (
        <p id={hintId} className="tama-field__hint">
          {hint}
        </p>
      )}
      {error && (
        <p id={errorId} className="tama-field__error">
          {error}
        </p>
      )}
    </div>
  )
}
