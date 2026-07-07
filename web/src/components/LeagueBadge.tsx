import type { CSSProperties } from 'react'
import type { IconSize } from './icons/icon'
import { ShieldIcon } from './icons/shield'
import './LeagueBadge.css'

// The M2 rungs of the league ladder; M7 adds the leagues above gold and
// extends this map alongside the tokens in tokens.css.
export const LEAGUE_TONES = {
  bronze: 'var(--league-bronze)',
  silver: 'var(--league-silver)',
  gold: 'var(--league-gold)',
} as const

export type LeagueTone = keyof typeof LEAGUE_TONES

type LeagueBadgeProps = {
  tone: LeagueTone
  size?: IconSize
}

// The league shield in its metal tone. The shield inherits currentColor,
// so the badge only has to set the color.
export function LeagueBadge({ tone, size = 32 }: LeagueBadgeProps) {
  return (
    <span className="tama-league" style={{ '--league-tone': LEAGUE_TONES[tone] } as CSSProperties}>
      <ShieldIcon size={size} />
      <span className="visually-hidden">{`${tone} league`}</span>
    </span>
  )
}
