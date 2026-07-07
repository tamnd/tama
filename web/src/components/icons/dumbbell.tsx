import type { IconProps } from './icon'

// Dumbbell, single tone. The practice hub.
export function DumbbellIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <rect x="1.5" y="9.5" width="3" height="5" rx="1.2" />
      <rect x="4.5" y="6.5" width="4" height="11" rx="1.6" />
      <rect x="8.5" y="10.25" width="7" height="3.5" />
      <rect x="15.5" y="6.5" width="4" height="11" rx="1.6" />
      <rect x="19.5" y="9.5" width="3" height="5" rx="1.2" />
    </svg>
  )
}
