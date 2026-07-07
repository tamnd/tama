import { CountUp } from './CountUp'
import { CrownIcon } from './icons/crown'
import { GemIcon } from './icons/gem'
import { LightningIcon } from './icons/lightning'
import { StreakFlameIcon } from './icons/streak-flame'
import './StatChip.css'

export type StatKind = 'streak' | 'gems' | 'xp'

type StatChipProps = {
  kind: StatKind
  value: number
}

const ICONS: Record<StatKind, (size: 16 | 20 | 24 | 32) => React.ReactNode> = {
  streak: (size) => <StreakFlameIcon size={size} />,
  gems: (size) => <GemIcon size={size} />,
  xp: (size) => <LightningIcon size={size} />,
}

function labelFor(kind: StatKind, value: number): string {
  switch (kind) {
    case 'streak':
      return `Streak: ${value} ${value === 1 ? 'day' : 'days'}`
    case 'gems':
      return `Gems: ${value}`
    case 'xp':
      return `XP: ${value}`
  }
}

// Top-bar counter: an icon plus a 700 number, no border, no background.
// Streak counts in fox, gems in macaw, XP in bee; a zero-day streak goes
// hare all over, the unlit flame. The count sweeps up through CountUp and
// screen readers get the full phrase through the live label.
export function StatChip({ kind, value }: StatChipProps) {
  const zero = kind === 'streak' && value === 0
  const className = ['tama-statchip', `tama-statchip--${kind}`, zero ? 'tama-statchip--zero' : '']
    .filter(Boolean)
    .join(' ')

  return (
    <span className={className}>
      <span className="tama-statchip__icon" aria-hidden="true">
        {ICONS[kind](24)}
      </span>
      <span className="tama-statchip__value" aria-hidden="true">
        <CountUp value={value} />
      </span>
      <span className="visually-hidden" aria-live="polite">
        {labelFor(kind, value)}
      </span>
    </span>
  )
}

type CrownChipProps = {
  level: number
}

// Crown level chip: the crown over the level number in a bee-tinted pill,
// worn by path headers and the profile.
export function CrownChip({ level }: CrownChipProps) {
  return (
    <span className="tama-crownchip">
      <span className="tama-crownchip__icon" aria-hidden="true">
        <CrownIcon size={20} />
      </span>
      <span className="tama-crownchip__value" aria-hidden="true">
        <CountUp value={level} />
      </span>
      <span className="visually-hidden" aria-live="polite">{`Crown level: ${level}`}</span>
    </span>
  )
}
