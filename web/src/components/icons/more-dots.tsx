import type { IconProps } from './icon'

// Three horizontal dots, single tone. The MORE tab and overflow menus.
export function MoreDotsIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <circle cx="4.5" cy="12" r="2.4" />
      <circle cx="12" cy="12" r="2.4" />
      <circle cx="19.5" cy="12" r="2.4" />
    </svg>
  )
}
