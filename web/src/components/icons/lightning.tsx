import type { IconProps } from './icon'

// Lightning bolt, canonical bee fill. XP counters and boosts.
export function LightningIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        d="M13.6 1.7c.7-.9 2.1-.2 1.8.9L13.7 9h4.9c.9 0 1.4 1 .8 1.7L10.4 22.3c-.7.9-2.1.2-1.8-.9L10.3 15H5.4c-.9 0-1.4-1-.8-1.7z"
        fill="var(--lightning, var(--bee))"
      />
    </svg>
  )
}
