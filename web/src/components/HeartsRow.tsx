import { useEffect, useRef, useState } from 'react'
import { HeartIcon } from './icons/heart'
import './HeartsRow.css'

const TOTAL = 5

type HeartsRowProps = {
  /** Hearts left, 0 through 5. */
  remaining: number
  /** Super state: one infinity-marked heart instead of the row. */
  unlimited?: boolean
}

// The five-heart row: filled cardinal for what's left, hollow swan for
// what's spent. A heart lost while mounted pops (pop-in reversed) and then
// hollows over dur-fast; unlimited swaps the row for a single heart wearing
// the infinity.
export function HeartsRow({ remaining, unlimited = false }: HeartsRowProps) {
  const clamped = Math.max(0, Math.min(TOTAL, remaining))

  // Hearts that just went from filled to spent get the losing class so the
  // pop and the hollow transition play; a refill clears the set.
  const prev = useRef(clamped)
  const [losing, setLosing] = useState<number[]>([])
  useEffect(() => {
    if (clamped < prev.current) {
      const lost: number[] = []
      for (let i = clamped; i < prev.current; i++) lost.push(i)
      setLosing(lost)
    } else if (clamped > prev.current) {
      setLosing([])
    }
    prev.current = clamped
  }, [clamped])

  if (unlimited) {
    return (
      <span className="tama-hearts tama-hearts--unlimited">
        <span className="tama-hearts__heart" aria-hidden="true">
          <HeartIcon infinity />
        </span>
        <span className="visually-hidden" aria-live="polite">
          Hearts: unlimited
        </span>
      </span>
    )
  }

  return (
    <span className="tama-hearts">
      {Array.from({ length: TOTAL }, (_, i) => {
        const spent = i >= clamped
        const className = [
          'tama-hearts__heart',
          spent ? 'tama-hearts__heart--spent' : '',
          spent && losing.includes(i) ? 'tama-hearts__heart--losing' : '',
        ]
          .filter(Boolean)
          .join(' ')
        return (
          <span key={i} className={className} aria-hidden="true">
            <HeartIcon />
          </span>
        )
      })}
      <span className="visually-hidden" aria-live="polite">
        {`Hearts: ${clamped} of ${TOTAL}`}
      </span>
    </span>
  )
}
