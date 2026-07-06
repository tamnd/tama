import type { IconProps } from './icon'

// Padlock, single tone. Locked path nodes and gated features.
export function LockIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        d="M7.5 10.5V7a4.5 4.5 0 0 1 9 0v3.5"
        fill="none"
        stroke="currentColor"
        strokeWidth="3"
        strokeLinecap="round"
      />
      <rect x="4" y="10" width="16" height="11.5" rx="3" fill="currentColor" />
    </svg>
  )
}
