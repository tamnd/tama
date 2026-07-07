import type { IconProps } from './icon'

// Notification bell with clapper, single tone.
export function BellIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <path d="M12 1.8a6.7 6.7 0 0 0-6.7 6.7v3.6l-1.5 4c-.3.9.3 1.8 1.2 1.8h14c.9 0 1.5-.9 1.2-1.8l-1.5-4V8.5A6.7 6.7 0 0 0 12 1.8z" />
      <path d="M9.4 19.5a2.6 2.6 0 0 0 5.2 0z" />
    </svg>
  )
}
