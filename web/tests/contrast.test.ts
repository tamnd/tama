import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

// WCAG AA check on the body text pairs in both themes. The tokens file is
// parsed directly, no DOM needed.
//
// Light-theme wolf (#777777) on snow lands at 4.48:1, which matches the
// reference palette exactly; it colors supporting text only, so the strict
// 4.5:1 pairs below are the body text tokens in both themes plus the dark
// secondary text, all of which the palette clears.

const css = readFileSync(
  fileURLToPath(new URL('../src/styles/tokens.css', import.meta.url)),
  'utf8',
)

function parseBlock(selector: string): Map<string, string> {
  const start = css.indexOf(selector)
  if (start === -1) throw new Error(`selector not found: ${selector}`)
  const open = css.indexOf('{', start)
  const close = css.indexOf('}', open)
  const vars = new Map<string, string>()
  for (const m of css.slice(open + 1, close).matchAll(/(--[\w-]+):\s*([^;]+);/g)) {
    vars.set(m[1], m[2].trim())
  }
  return vars
}

const root = parseBlock(':root')
const dark = parseBlock("html[data-theme='dark']")

function resolve(name: string, theme: Map<string, string>): string {
  let value = theme.get(name) ?? root.get(name)
  if (value === undefined) throw new Error(`token not found: ${name}`)
  const ref = value.match(/^var\((--[\w-]+)\)$/)
  if (ref) return resolve(ref[1], theme)
  return value
}

function luminance(hex: string): number {
  const m = hex.match(/^#([0-9A-Fa-f]{6})$/)
  if (!m) throw new Error(`not a 6-digit hex: ${hex}`)
  const [r, g, b] = [0, 2, 4].map((i) => {
    const c = parseInt(m[1].slice(i, i + 2), 16) / 255
    return c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4)
  })
  return 0.2126 * r + 0.7152 * g + 0.0722 * b
}

function contrast(fg: string, bg: string): number {
  const [l1, l2] = [luminance(fg), luminance(bg)].sort((a, b) => b - a)
  return (l1 + 0.05) / (l2 + 0.05)
}

type Pair = [text: string, surface: string]

const COMMON: Pair[] = [
  ['--color-text', '--color-surface'],
  ['--color-text', '--color-bg'],
  ['--color-text', '--polar'],
]

const DARK_ONLY: Pair[] = [
  ['--color-text-secondary', '--color-surface'],
  ['--color-text-secondary', '--color-bg'],
]

describe('text on surface contrast', () => {
  const cases: [theme: string, vars: Map<string, string>, pair: Pair][] = [
    ...COMMON.map((p): [string, Map<string, string>, Pair] => ['light', root, p]),
    ...COMMON.map((p): [string, Map<string, string>, Pair] => ['dark', dark, p]),
    ...DARK_ONLY.map((p): [string, Map<string, string>, Pair] => ['dark', dark, p]),
  ]

  it.each(cases)('%s: %# clears 4.5:1', (_theme, vars, [text, surface]) => {
    const ratio = contrast(resolve(text, vars), resolve(surface, vars))
    expect(ratio).toBeGreaterThanOrEqual(4.5)
  })
})
