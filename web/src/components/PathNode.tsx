import { CheckIcon } from './icons/check'
import { ChestIcon } from './icons/chest'
import { LockIcon } from './icons/lock'
import { StarIcon } from './icons/star'
import { TrophyIcon } from './icons/trophy'
import './PathNode.css'

export type PathNodeState = 'active' | 'completed' | 'locked' | 'legendary'

export type PathNodeProgress = {
  /** Levels passed so far, filled in feather green. */
  passed: number
  /** One ring segment per level. */
  total: number
}

type PathNodeProps = {
  state: PathNodeState
  /** Accessible name, e.g. "Order food and drink, lesson 3". */
  label: string
  /** Renders the segmented progress ring instead of the gold ring. */
  progress?: PathNodeProgress
  /** Frontier of an untouched unit: swaps the START bubble for JUMP HERE?. */
  jumpHere?: boolean
  onClick?: () => void
}

// A 70px circular 3D button, the same face-plus-darker-bottom-edge
// construction as Button but with the chunkier 8px depth edge. Locked nodes
// stay focusable for discovery but expose aria-disabled and ignore clicks.
export function PathNode({ state, label, progress, jumpHere = false, onClick }: PathNodeProps) {
  const locked = state === 'locked'
  const classes = ['tama-path-node', `tama-path-node--${state}`]
  if (progress) classes.push('tama-path-node--ringed')

  return (
    <button
      type="button"
      className={classes.join(' ')}
      aria-label={label}
      aria-disabled={locked || undefined}
      onClick={locked ? undefined : onClick}
    >
      {jumpHere ? (
        <span className="tama-path-node__bubble tama-path-node__bubble--jump" aria-hidden="true">
          Jump here?
        </span>
      ) : (
        state === 'active' && (
          <span
            className="tama-path-node__bubble tama-path-node__bubble--start"
            aria-hidden="true"
          >
            Start
          </span>
        )
      )}
      {progress && <ProgressRing passed={progress.passed} total={progress.total} />}
      <span className="tama-path-node__icon">{STATE_ICONS[state]}</span>
    </button>
  )
}

const STATE_ICONS = {
  active: <StarIcon size={32} />,
  completed: <CheckIcon size={32} />,
  locked: <LockIcon size={32} />,
  legendary: <TrophyIcon size={32} />,
}

// One segment per level, drawn as dash slices of a circle that starts at
// twelve o'clock. Feather green fills over the swan track.
function ProgressRing({ passed, total }: PathNodeProgress) {
  const r = 41
  const c = 2 * Math.PI * r
  const seg = c / total
  const gap = total > 1 ? 6 : 0

  return (
    <svg className="tama-path-node__ring" viewBox="0 0 94 94" aria-hidden="true">
      {Array.from({ length: total }, (_, i) => (
        <circle
          key={i}
          className={
            i < passed
              ? 'tama-path-node__ring-seg tama-path-node__ring-seg--passed'
              : 'tama-path-node__ring-seg'
          }
          cx="47"
          cy="47"
          r={r}
          fill="none"
          strokeWidth="4"
          strokeLinecap="round"
          strokeDasharray={`${seg - gap} ${c - seg + gap}`}
          strokeDashoffset={c / 4 - i * seg}
        />
      ))}
    </svg>
  )
}

type ChestNodeProps = {
  /** Already collected: lid up, flat grey, not interactive. */
  opened?: boolean
  label: string
  onClick?: () => void
}

// The reward chest path entry. Openable chests shimmer gently at idle;
// opened ones go flat grey through the chest icon's tone properties.
export function ChestNode({ opened = false, label, onClick }: ChestNodeProps) {
  const classes = ['tama-path-chest']
  if (opened) classes.push('tama-path-chest--opened')

  return (
    <button
      type="button"
      className={classes.join(' ')}
      aria-label={label}
      aria-disabled={opened || undefined}
      onClick={opened ? undefined : onClick}
    >
      <span className="tama-path-chest__art">
        <ChestIcon open={opened} />
      </span>
    </button>
  )
}

type CharacterGateProps = {
  /** The unit boundary has been cleared. */
  passed?: boolean
  label: string
  onClick?: () => void
}

// The oval 90x70 node at unit boundaries. The mascot slot renders a plain
// silhouette until the Tama poses land in a later slice.
export function CharacterGate({ passed = false, label, onClick }: CharacterGateProps) {
  return (
    <button
      type="button"
      className={`tama-path-gate tama-path-gate--${passed ? 'passed' : 'locked'}`}
      aria-label={label}
      aria-disabled={!passed || undefined}
      onClick={passed ? onClick : undefined}
    >
      <span className="tama-path-gate__mascot">{silhouette}</span>
    </button>
  )
}

// Placeholder cat silhouette: ears, head, haunches. Single tone so the gate
// states color it like any icon.
const silhouette = (
  <svg viewBox="0 0 48 48" fill="currentColor" aria-hidden="true">
    <path d="M14 14 11 3l10 7h6l10-7-3 11z" />
    <circle cx="24" cy="21" r="11" />
    <ellipse cx="24" cy="39" rx="14" ry="7" />
  </svg>
)
