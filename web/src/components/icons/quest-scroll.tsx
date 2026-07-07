import type { IconProps } from './icon'

// Quest scroll: a rolled top edge over a sheet with two cut-out text lines,
// single tone.
export function QuestScrollIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="currentColor" aria-hidden="true">
      <rect x="3" y="1.5" width="18" height="4.5" rx="2.25" />
      <path
        fillRule="evenodd"
        d="M5.5 5h13v14.5A3 3 0 0 1 15.5 22.5h-7a3 3 0 0 1-3-3zM8.5 9.5h7a1 1 0 0 1 0 2h-7a1 1 0 0 1 0-2zm0 4.5h7a1 1 0 0 1 0 2h-7a1 1 0 0 1 0-2z"
      />
    </svg>
  )
}
