import type { IconProps } from './icon'

// Streak flame, two tones: the canonical bee core over a fox body.
export function StreakFlameIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        d="M12 1.5c.6 3.4-.6 5.1-2.3 7C7.8 10.7 6 13 6 15.9a6 6 0 0 0 12 0c0-2.5-1.2-4.4-2.3-6.1-1.4-2.2-2.9-4.5-3.7-8.3z"
        fill="var(--fox)"
      />
      <path
        d="M12 10.2c.3 1.8 1 2.9 1.7 4 .6.9 1.3 2 1.3 3.1a3 3 0 0 1-6 0c0-1.5.9-2.7 1.8-3.9.5-.8 1-1.6 1.2-3.2z"
        fill="var(--bee)"
      />
    </svg>
  )
}
