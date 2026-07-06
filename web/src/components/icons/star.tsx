import type { IconProps } from './icon'

// Five-point star, single tone. Marks the active path node.
export function StarIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <path d="M12 1.8c.5 0 .9.3 1.1.7l2.5 5.4 5.9.7c1 .1 1.4 1.4.7 2.1l-4.4 4 1.2 5.8c.2 1-.9 1.7-1.7 1.2L12 18.8l-5.3 2.9c-.8.5-1.9-.2-1.7-1.2l1.2-5.8-4.4-4c-.7-.7-.3-2 .7-2.1l5.9-.7 2.5-5.4c.2-.4.6-.7 1.1-.7z" />
    </svg>
  )
}
