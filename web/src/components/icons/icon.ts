// Shared contract for the icon set. Icons are drawn on a 24x24 viewBox and
// render only at the four grid sizes; every icon is aria-hidden, meaning is
// carried by adjacent text or an aria-label on the parent control.
// UI chrome icons are single tone and inherit currentColor; gamification
// icons carry their canonical token fills.

export type IconSize = 16 | 20 | 24 | 32

export type IconProps = {
  size?: IconSize
}
