import type { IconProps } from './icon'

// Faceted gem, two tones: a macaw core over a humpback body. Currency
// counters and the shop.
export function GemIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        d="M8 2.5h8a1.5 1.5 0 0 1 1.2.6l3.9 5.1c.4.5.4 1.3 0 1.8l-8 11.3a1.4 1.4 0 0 1-2.3 0L2.9 10a1.5 1.5 0 0 1 0-1.8l3.9-5.1A1.5 1.5 0 0 1 8 2.5z"
        fill="var(--gem-body, var(--humpback))"
      />
      <path d="M12 5.6l4.8 3.6-4.8 8.2-4.8-8.2z" fill="var(--gem-core, var(--macaw))" />
    </svg>
  )
}
