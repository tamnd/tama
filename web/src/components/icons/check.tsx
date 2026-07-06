import type { IconProps } from './icon'

// Bold checkmark, single tone. Completed nodes, correct feedback.
export function CheckIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="none" aria-hidden="true">
      <path
        d="M4 12.8l5.2 5.2L20 7"
        stroke="currentColor"
        strokeWidth="4"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}
