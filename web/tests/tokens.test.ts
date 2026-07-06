import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

// Every brand hex must appear verbatim in tokens.css, so a refactor cannot
// silently drift a color. Values come from the M2 design system spec.

const css = readFileSync(
  fileURLToPath(new URL('../src/styles/tokens.css', import.meta.url)),
  'utf8',
)

const LIGHT: [string, string][] = [
  ['--feather-green', '#58CC02'],
  ['--mask-green', '#89E219'],
  ['--macaw', '#1CB0F6'],
  ['--cardinal', '#FF4B4B'],
  ['--bee', '#FFC800'],
  ['--fox', '#FF9600'],
  ['--beetle', '#CE82FF'],
  ['--humpback', '#2B70C9'],
  ['--whale', '#084C77'],
  ['--eel', '#4B4B4B'],
  ['--wolf', '#777777'],
  ['--hare', '#AFAFAF'],
  ['--swan', '#E5E5E5'],
  ['--polar', '#F7F7F7'],
  ['--snow', '#FFFFFF'],
  ['--feather-green-deep', '#58A700'],
  ['--macaw-deep', '#1899D6'],
  ['--cardinal-deep', '#EA2B2B'],
  ['--bee-deep', '#E5B800'],
  ['--beetle-deep', '#A568CC'],
  ['--swan-deep', '#CFCFCF'],
  ['--correct-fill', '#D7FFB8'],
  ['--incorrect-fill', '#FFDFE0'],
]

// The dark base palette from the spec. Token assignments live in the
// html[data-theme="dark"] block; here we only pin the hex values.
const DARK_HEX = ['#131F24', '#202F36', '#37464F', '#F1F7FB', '#DCE6EC', '#52656D']

describe('tokens.css', () => {
  it.each(LIGHT)('%s is %s', (name, hex) => {
    expect(css).toContain(`${name}: ${hex}`)
  })

  it('has a dark theme block', () => {
    expect(css).toContain("html[data-theme='dark']")
  })

  it.each(DARK_HEX)('dark palette contains %s', (hex) => {
    const darkBlock = css.slice(css.indexOf('html[data-theme='))
    expect(darkBlock).toContain(hex)
  })
})
