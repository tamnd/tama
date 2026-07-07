import type { IconProps } from './icon'

// League shield, single tone. Inherits currentColor so the league badge can
// paint it in a metal tone; the nav rail uses it plain.
export function ShieldIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        d="M11.5 1.9a1.4 1.4 0 0 1 1 0l7.6 2.9c.5.2.9.7.9 1.3v5.7c0 5.1-3.5 8.7-8.6 10.4a1.4 1.4 0 0 1-.8 0C6.5 20.5 3 16.9 3 11.8V6.1c0-.6.4-1.1.9-1.3z"
        fill="currentColor"
      />
    </svg>
  )
}
