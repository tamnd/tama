import type { IconProps } from './icon'

// Person bust, single tone. The PROFILE tab.
export function ProfileIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <circle cx="12" cy="7" r="4.6" />
      <path d="M12 13.6c-4.7 0-8.2 2.6-8.2 6 0 1.4 1.1 2.4 2.5 2.4h11.4c1.4 0 2.5-1 2.5-2.4 0-3.4-3.5-6-8.2-6z" />
    </svg>
  )
}
