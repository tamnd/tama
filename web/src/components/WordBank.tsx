import { useLayoutEffect, useRef, useState, type KeyboardEvent } from 'react'
import { TapToken } from './TapToken'
import './WordBank.css'

export type WordBankToken = {
  id: string
  text: string
}

type WordBankProps = {
  /** Bank chips in display order; ids must be unique. */
  tokens: WordBankToken[]
  /** Name of the hidden input that serializes the answer. */
  name?: string
  /** Accessible name for the answer line. */
  label?: string
  onChange?: (ids: string[]) => void
}

// The tap-the-words exercise control: dashed answer slots above, the chip
// bank below. Number keys 1-9 pick bank chips, Backspace returns the last,
// and every move FLIP-animates from old to new position. A hidden input
// carries the ordered token ids for the grading flow.
export function WordBank({
  tokens,
  name = 'answer',
  label = 'Your answer',
  onChange,
}: WordBankProps) {
  const [selected, setSelected] = useState<string[]>([])

  // One ref per token id, always pointing at its live (interactive) chip:
  // the bank chip while unselected, the answer chip once picked. FLIP then
  // reads the previous rect from wherever the chip last sat.
  const chips = useRef(new Map<string, HTMLButtonElement>())
  const rects = useRef(new Map<string, DOMRect>())

  useLayoutEffect(() => {
    chips.current.forEach((el, id) => {
      const next = el.getBoundingClientRect()
      const prev = rects.current.get(id)
      rects.current.set(id, next)
      if (!prev || typeof el.animate !== 'function') return
      const dx = prev.left - next.left
      const dy = prev.top - next.top
      if (!dx && !dy) return
      // Reduced motion zeroes --dur-base (motion.css), so the move snaps.
      const raw = parseFloat(getComputedStyle(el).getPropertyValue('--dur-base'))
      const duration = Number.isNaN(raw) ? 300 : raw
      if (duration === 0) return
      el.animate([{ transform: `translate(${dx}px, ${dy}px)` }, { transform: 'translate(0, 0)' }], {
        duration,
        easing: 'cubic-bezier(0.22, 1, 0.36, 1)',
      })
    })
  })

  function chipRef(id: string) {
    return (el: HTMLButtonElement | null) => {
      if (el) chips.current.set(id, el)
      else chips.current.delete(id)
    }
  }

  // Mirrors the state so taps landing before a re-render still build on
  // the latest answer instead of a stale closure.
  const current = useRef(selected)

  function update(next: string[]) {
    current.current = next
    setSelected(next)
    onChange?.(next)
  }

  function pick(id: string) {
    if (!current.current.includes(id)) update([...current.current, id])
  }

  function put(id: string) {
    update(current.current.filter((s) => s !== id))
  }

  function onKeyDown(event: KeyboardEvent<HTMLDivElement>) {
    if (event.key === 'Backspace') {
      if (current.current.length === 0) return
      event.preventDefault()
      update(current.current.slice(0, -1))
      return
    }
    const n = Number(event.key)
    if (n >= 1 && n <= 9) {
      const token = tokens[n - 1]
      if (token && !current.current.includes(token.id)) {
        event.preventDefault()
        pick(token.id)
      }
    }
  }

  return (
    <div className="tama-word-bank" onKeyDown={onKeyDown}>
      <div className="tama-word-bank__answer" role="group" aria-label={label}>
        {tokens.map((_, i) => {
          const token = tokens.find((t) => t.id === selected[i])
          return (
            <span key={i} className="tama-word-bank__slot">
              {token && (
                <TapToken
                  ref={chipRef(token.id)}
                  className="tama-word-bank__chip"
                  onClick={() => put(token.id)}
                >
                  {token.text}
                </TapToken>
              )}
            </span>
          )
        })}
      </div>
      <div className="tama-word-bank__bank" role="group" aria-label="Word bank">
        {tokens.map((token) =>
          selected.includes(token.id) ? (
            <TapToken key={token.id} hollow>
              {token.text}
            </TapToken>
          ) : (
            <TapToken
              key={token.id}
              ref={chipRef(token.id)}
              className="tama-word-bank__chip"
              onClick={() => pick(token.id)}
            >
              {token.text}
            </TapToken>
          ),
        )}
      </div>
      <input type="hidden" name={name} value={selected.join(' ')} readOnly />
    </div>
  )
}
