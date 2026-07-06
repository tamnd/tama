import type { IconProps } from './icon'

// Speaker with sound waves, single tone. Audio prompts in speech bubbles.
export function SpeakerIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="none" aria-hidden="true">
      <path
        d="M3.5 9.2v5.6c0 .7.5 1.2 1.2 1.2h2.6l4.3 3.6c.8.7 2 .1 2-1V5.4c0-1.1-1.2-1.7-2-1L7.3 8H4.7c-.7 0-1.2.5-1.2 1.2z"
        fill="currentColor"
      />
      <path
        d="M16.5 9a4.6 4.6 0 0 1 0 6"
        stroke="currentColor"
        strokeWidth="2.5"
        strokeLinecap="round"
      />
      <path
        d="M19 6.5a8.2 8.2 0 0 1 0 11"
        stroke="currentColor"
        strokeWidth="2.5"
        strokeLinecap="round"
      />
    </svg>
  )
}
