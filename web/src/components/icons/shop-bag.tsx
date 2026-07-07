import type { IconProps } from './icon'

// Shopping bag with a rounded handle, single tone. The SHOP tab.
export function ShopBagIcon({ size = 24 }: IconProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} fill="none" aria-hidden="true">
      <path
        d="M8.5 11V6.75a3.5 3.5 0 0 1 7 0V11"
        stroke="currentColor"
        strokeWidth="2.5"
        strokeLinecap="round"
      />
      <path
        d="M5.4 8h13.2a1.5 1.5 0 0 1 1.5 1.4l.8 10.4a2.5 2.5 0 0 1-2.5 2.7H5.6a2.5 2.5 0 0 1-2.5-2.7l.8-10.4A1.5 1.5 0 0 1 5.4 8z"
        fill="currentColor"
      />
    </svg>
  )
}
