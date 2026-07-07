import type { IconProps } from './icon'

// Settings gear: eight square teeth around a ring, single tone.
export function GearIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <g>
        <rect x="10.4" y="1.5" width="3.2" height="5" rx="1.3" />
        <rect x="10.4" y="17.5" width="3.2" height="5" rx="1.3" />
      </g>
      <g transform="rotate(45 12 12)">
        <rect x="10.4" y="1.5" width="3.2" height="5" rx="1.3" />
        <rect x="10.4" y="17.5" width="3.2" height="5" rx="1.3" />
      </g>
      <g transform="rotate(90 12 12)">
        <rect x="10.4" y="1.5" width="3.2" height="5" rx="1.3" />
        <rect x="10.4" y="17.5" width="3.2" height="5" rx="1.3" />
      </g>
      <g transform="rotate(135 12 12)">
        <rect x="10.4" y="1.5" width="3.2" height="5" rx="1.3" />
        <rect x="10.4" y="17.5" width="3.2" height="5" rx="1.3" />
      </g>
      <path
        fillRule="evenodd"
        d="M12 4.5a7.5 7.5 0 1 1 0 15 7.5 7.5 0 0 1 0-15zm0 4.2a3.3 3.3 0 1 0 0 6.6 3.3 3.3 0 0 0 0-6.6z"
      />
    </svg>
  )
}
