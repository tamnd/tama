import type { CSSProperties } from 'react'
import './Avatar.css'

export type AvatarSize = 32 | 48 | 96

type AvatarProps = {
  /** Display name; drives the fallback letter and its color. */
  name: string
  /** Picture URL; the initial fallback renders without one. */
  src?: string
  size?: AvatarSize
}

// The fallback swatches: brand colors that keep a snow initial readable.
const AVATAR_COLORS = [
  '--feather-green',
  '--macaw',
  '--cardinal',
  '--bee',
  '--fox',
  '--beetle',
  '--humpback',
] as const

// Same name, same color, on every machine and every render: a small
// deterministic hash over the code points picks from the brand set.
export function avatarColor(name: string): string {
  let hash = 0
  for (const ch of name) {
    hash = (hash * 31 + (ch.codePointAt(0) ?? 0)) >>> 0
  }
  return `var(${AVATAR_COLORS[hash % AVATAR_COLORS.length]})`
}

// Circular avatar with a 2px swan border at 32, 48, or 96. No picture
// means the name's first letter on its deterministic brand color.
export function Avatar({ name, src, size = 48 }: AvatarProps) {
  const initial = [...name.trim()][0]?.toLocaleUpperCase() ?? '?'

  return (
    <span
      className={`tama-avatar tama-avatar--${size}`}
      role="img"
      aria-label={name}
      style={{ '--avatar-bg': avatarColor(name) } as CSSProperties}
    >
      {src ? (
        <img className="tama-avatar__image" src={src} alt="" />
      ) : (
        <span className="tama-avatar__initial" aria-hidden="true">
          {initial}
        </span>
      )}
    </span>
  )
}
