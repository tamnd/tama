import { useEffect, useRef, useState, type CSSProperties } from 'react'
import { StreakFlameIcon } from './icons/streak-flame'
import './ProgressBar.css'

type ProgressBarProps = {
  value: number
  max: number
  /** The caller decides the combo threshold; the bar only renders it. */
  showFlame?: boolean
  /** Accessible name, e.g. "Lesson progress". */
  label?: string
}

// The 16px lesson progress pill: swan track, feather fill with the
// mask-green top highlight, a shimmer sweep per increment, and the combo
// flame slot at the right end.
export function ProgressBar({ value, max, showFlame = false, label }: ProgressBarProps) {
  const pct = max > 0 ? Math.min(100, Math.max(0, (value / max) * 100)) : 0

  // Each increment remounts the shimmer band (keyed by count) so its
  // one-shot sweep replays.
  const prev = useRef(value)
  const [sweep, setSweep] = useState(0)
  useEffect(() => {
    if (value > prev.current) setSweep((n) => n + 1)
    prev.current = value
  }, [value])

  return (
    <div
      className="tama-progress"
      role="progressbar"
      aria-valuenow={value}
      aria-valuemin={0}
      aria-valuemax={max}
      aria-label={label}
    >
      <div className="tama-progress__fill" style={{ '--progress': `${pct}%` } as CSSProperties}>
        {sweep > 0 && <span key={sweep} className="tama-progress__shimmer" aria-hidden="true" />}
      </div>
      {showFlame && (
        <span className="tama-progress__flame" aria-hidden="true">
          <StreakFlameIcon size={24} />
        </span>
      )}
    </div>
  )
}
