import type { IconProps } from './icon'

// Closed book with a page slot near the spine, single tone. Course and
// guidebook entries.
export function BookIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        d="M4 5.5A3.5 3.5 0 0 1 7.5 2H19a1 1 0 0 1 1 1v13a1 1 0 0 1-1 1H7.5c-.9 0-1.5.7-1.5 1.5S6.6 20 7.5 20H19a1 1 0 0 1 1 1 1 1 0 0 1-1 1H7.5A3.5 3.5 0 0 1 4 18.5z"
        fill="currentColor"
      />
    </svg>
  )
}
