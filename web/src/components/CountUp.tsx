import { useEffect, useRef, useState } from 'react'
import './CountUp.css'

type CountUpProps = {
  value: number
  /** Sweep length in ms; the spec's 600ms unless a caller needs otherwise. */
  duration?: number
}

// People who asked for less motion get the final number straight away,
// whether the ask came from the OS or the in-app class.
function prefersReducedMotion(): boolean {
  if (typeof window.matchMedia === 'function' && window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return true
  }
  return document.documentElement.classList.contains('reduce-motion')
}

// The celebration counter: sweeps from 0 to value over 600ms with an
// ease-out, 900 weight. Later increases sweep from the shown number so a
// live counter never restarts at zero.
export function CountUp({ value, duration = 600 }: CountUpProps) {
  const shownRef = useRef(prefersReducedMotion() ? value : 0)
  const [shown, setShown] = useState(shownRef.current)

  useEffect(() => {
    if (prefersReducedMotion()) {
      shownRef.current = value
      setShown(value)
      return
    }
    const from = shownRef.current
    if (from === value) return
    const start = performance.now()
    let frame = requestAnimationFrame(function tick(now: number) {
      const t = Math.min(1, (now - start) / duration)
      const eased = 1 - Math.pow(1 - t, 3)
      const next = Math.round(from + (value - from) * eased)
      shownRef.current = next
      setShown(next)
      if (t < 1) frame = requestAnimationFrame(tick)
    })
    return () => cancelAnimationFrame(frame)
  }, [value, duration])

  return <span className="tama-countup">{shown}</span>
}
