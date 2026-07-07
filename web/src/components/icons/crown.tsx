import type { IconProps } from './icon'

// Crown, canonical bee fill. Level chips on path headers and profile.
export function CrownIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        d="M11 3.6a1.2 1.2 0 0 1 2 0l3.1 4.6 4-3c.9-.6 2 .1 1.9 1.2l-1.5 9.7c-.1.7-.7 1.3-1.5 1.3H5c-.8 0-1.4-.6-1.5-1.3L2 6.4c-.1-1.1 1-1.8 1.9-1.2l4 3z"
        fill="var(--crown, var(--bee))"
      />
      <path
        d="M4.6 19.4h14.8v1.4a1.4 1.4 0 0 1-1.4 1.4H6a1.4 1.4 0 0 1-1.4-1.4z"
        fill="var(--crown, var(--bee))"
      />
    </svg>
  )
}
