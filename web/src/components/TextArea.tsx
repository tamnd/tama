import { useLayoutEffect, useRef, type ComponentPropsWithRef, type InputEvent } from 'react'
import './TextArea.css'

type TextAreaProps = ComponentPropsWithRef<'textarea'>

// The multi-line field for type-what-you-hear and translations. Grows with
// its content between 2 and 6 rows; the measured height rides a custom
// property so the inline style only ever carries the dynamic value.
export function TextArea({ className, rows = 2, ref, onInput, ...rest }: TextAreaProps) {
  const classes = ['tama-textarea']
  if (className) classes.push(className)

  const inner = useRef<HTMLTextAreaElement | null>(null)

  useLayoutEffect(() => {
    if (inner.current) fit(inner.current)
  }, [])

  function setRef(el: HTMLTextAreaElement | null) {
    inner.current = el
    if (typeof ref === 'function') ref(el)
    else if (ref) ref.current = el
  }

  function handleInput(event: InputEvent<HTMLTextAreaElement>) {
    fit(event.currentTarget)
    onInput?.(event)
  }

  return (
    <textarea
      ref={setRef}
      rows={rows}
      className={classes.join(' ')}
      onInput={handleInput}
      {...rest}
    />
  )
}

// Measure at natural height, then pin; min/max-height in the CSS clamp the
// result to the 2 to 6 row band.
function fit(el: HTMLTextAreaElement) {
  el.style.setProperty('--textarea-height', 'auto')
  const borders = el.offsetHeight - el.clientHeight
  el.style.setProperty('--textarea-height', `${el.scrollHeight + borders}px`)
}
