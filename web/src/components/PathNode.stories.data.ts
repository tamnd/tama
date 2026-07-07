// Fixture data for the gallery's serpentine demo. Nodes snake left and
// right as the path descends; the horizontal offset cycles with the index.

import type { PathNodeProgress, PathNodeState } from './PathNode'

export const SERPENTINE_OFFSETS = [0, 44, 0, -44] as const

export function serpentineOffset(index: number): number {
  return SERPENTINE_OFFSETS[index % SERPENTINE_OFFSETS.length] ?? 0
}

export type DemoEntry =
  | { kind: 'node'; state: PathNodeState; label: string; progress?: PathNodeProgress }
  | { kind: 'chest'; opened: boolean; label: string }
  | { kind: 'gate'; passed: boolean; label: string }

export const DEMO_PATH: DemoEntry[] = [
  { kind: 'node', state: 'completed', label: 'Lesson 1, completed' },
  { kind: 'node', state: 'completed', label: 'Lesson 2, completed' },
  { kind: 'node', state: 'active', label: 'Lesson 3, current', progress: { passed: 2, total: 5 } },
  { kind: 'chest', opened: false, label: 'Reward chest' },
  { kind: 'node', state: 'locked', label: 'Lesson 4, locked' },
  { kind: 'node', state: 'locked', label: 'Lesson 5, locked' },
  { kind: 'gate', passed: false, label: 'Unit checkpoint' },
  { kind: 'node', state: 'legendary', label: 'Legendary challenge' },
]
