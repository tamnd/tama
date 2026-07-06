import type { IconProps } from './icon'

type ChestProps = IconProps & {
  /** Lid up with the reward glow showing. */
  open?: boolean
}

// Reward chest, two tones. Canonical fills are bee over fox; the path node
// repaints both custom properties grey once the chest has been opened.
export function ChestIcon({ size = 24, open = false }: ChestProps) {
  if (open) {
    return (
      <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
        {/* lid tipped back */}
        <path
          d="M4.5 6.5 6 2.8A1.5 1.5 0 0 1 7.4 1.9h9.2a1.5 1.5 0 0 1 1.4.9l1.5 3.7z"
          fill="var(--chest-band, var(--fox))"
        />
        {/* body */}
        <path
          d="M3 9.5h18V20a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"
          fill="var(--chest-body, var(--bee))"
        />
        {/* opening */}
        <rect x="3" y="9.5" width="18" height="3" fill="var(--chest-band, var(--fox))" />
        {/* clasp */}
        <rect x="10.5" y="14" width="3" height="4.5" rx="1" fill="var(--chest-band, var(--fox))" />
      </svg>
    )
  }
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      {/* lid */}
      <path d="M3 11V8a4 4 0 0 1 4-4h10a4 4 0 0 1 4 4v3z" fill="var(--chest-band, var(--fox))" />
      {/* body */}
      <path d="M3 11h18v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" fill="var(--chest-body, var(--bee))" />
      {/* clasp */}
      <rect x="10.5" y="9" width="3" height="5.5" rx="1" fill="var(--chest-band, var(--fox))" />
    </svg>
  )
}
