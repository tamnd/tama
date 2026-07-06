import type { IconProps } from './icon'

// Trophy cup, single tone. Legendary nodes and league rewards.
export function TrophyIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        d="M5.5 5H3v1.5A4.5 4.5 0 0 0 7.5 11M18.5 5H21v1.5A4.5 4.5 0 0 1 16.5 11"
        fill="none"
        stroke="currentColor"
        strokeWidth="2.5"
        strokeLinecap="round"
      />
      <path
        fill="currentColor"
        d="M6 2h12v7a6 6 0 0 1-12 0zM10.5 14.5h3V18h-3z"
      />
      <rect x="6.5" y="17.5" width="11" height="4.5" rx="1.5" fill="currentColor" />
    </svg>
  )
}
