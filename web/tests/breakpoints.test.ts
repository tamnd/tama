import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'
import { BP_MOBILE } from '../src/breakpoints'

// Media queries cannot read custom properties, so the mobile breakpoint
// lives twice: documented in the tokens.css header comment and exported
// from breakpoints.ts. This test is the sync guard between the two.

function read(rel: string): string {
  return readFileSync(fileURLToPath(new URL(rel, import.meta.url)), 'utf8')
}

describe('breakpoints', () => {
  it('matches the value documented in tokens.css', () => {
    const match = /The mobile breakpoint is (\d+)px/.exec(read('../src/styles/tokens.css'))
    expect(match).not.toBeNull()
    expect(Number(match?.[1])).toBe(BP_MOBILE)
  })

  it('is the value the shell media query uses', () => {
    expect(read('../src/app/Shell.css')).toContain(`(width < ${BP_MOBILE}px)`)
  })
})
