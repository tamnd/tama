import type { IconProps } from './icon'

// Report flag on a pole, single tone. The "report this sentence" affordance.
export function FlagReportIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <rect x="4.5" y="1.5" width="3" height="21" rx="1.5" />
      <path d="M9 3.5h10.2c.8 0 1.3.9.9 1.6l-2.4 3.9 2.4 3.9c.4.7-.1 1.6-.9 1.6H9z" />
    </svg>
  )
}
