import type { IconProps } from './icon'

// House with a door cut-out, single tone. The LEARN tab pointing home to
// the path.
export function HomePathIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <path
        fillRule="evenodd"
        d="M11 2.4a1.6 1.6 0 0 1 2 0l8.3 6.6c.9.7.4 2.1-.7 2.1H19v8.4a2.5 2.5 0 0 1-2.5 2.5h-9A2.5 2.5 0 0 1 5 19.5v-8.4H3.4c-1.1 0-1.6-1.4-.7-2.1zM12 13.4a2.1 2.1 0 0 0-2.1 2.1V22h4.2v-6.5a2.1 2.1 0 0 0-2.1-2.1z"
      />
    </svg>
  )
}
