import type { IconProps } from './icon'

type HeartProps = IconProps & {
  /** Spent heart: swan outline, no fill. */
  hollow?: boolean
  /** Unlimited-hearts marker drawn over the fill. */
  infinity?: boolean
}

// Heart, canonical cardinal fill. The hearts row hollows spent ones by
// swapping the fill and stroke custom properties, so CSS can transition the
// change; the unlimited state stamps an infinity over the fill.
export function HeartIcon({ size = 24, hollow = false, infinity = false }: HeartProps) {
  return (
    <svg viewBox="0 0 24 24" width={size} height={size} aria-hidden="true">
      <path
        className="tama-icon-heart__shape"
        d="M12 20.7C7.2 16.9 2.4 13.2 2.4 8.7 2.4 5.6 4.9 3.2 8 3.2c1.6 0 3.1.7 4 1.9.9-1.2 2.4-1.9 4-1.9 3.1 0 5.6 2.4 5.6 5.5 0 4.5-4.8 8.2-9.6 12z"
        fill={hollow ? 'transparent' : 'var(--heart-fill, var(--cardinal))'}
        stroke={hollow ? 'var(--heart-stroke, var(--swan))' : 'transparent'}
        strokeWidth="2.5"
        strokeLinejoin="round"
      />
      {infinity && (
        <path
          d="M12 10.5c-.8 1.2-1.5 1.9-2.6 1.9a1.9 1.9 0 1 1 0-3.8c1.1 0 1.8.7 2.6 1.9.8 1.2 1.5 1.9 2.6 1.9a1.9 1.9 0 1 0 0-3.8c-1.1 0-1.8.7-2.6 1.9z"
          fill="none"
          stroke="var(--heart-mark, var(--snow))"
          strokeWidth="2"
          strokeLinecap="round"
        />
      )}
    </svg>
  )
}
