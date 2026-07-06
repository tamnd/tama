import { readdirSync, readFileSync } from 'node:fs'
import { join, relative } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

// Repo-wide CSS guards: depth comes from borders, colors from tokens, and
// the typeface from base.css. Anything else is a regression.

const srcDir = fileURLToPath(new URL('../src', import.meta.url))

function walk(dir: string): string[] {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const full = join(dir, entry.name)
    return entry.isDirectory() ? walk(full) : [full]
  })
}

const files = walk(srcDir)
const cssFiles = files.filter((f) => f.endsWith('.css'))
const codeFiles = files.filter((f) => /\.(css|tsx|ts)$/.test(f))

function offenders(list: string[], pattern: RegExp, exempt: string[] = []): string[] {
  return list
    .filter((f) => !exempt.includes(relative(srcDir, f)))
    .filter((f) => pattern.test(readFileSync(f, 'utf8')))
    .map((f) => relative(srcDir, f))
}

describe('css guards', () => {
  it('no box-shadow anywhere, depth is borders', () => {
    expect(offenders(codeFiles, /box-shadow/)).toEqual([])
  })

  it('no hardcoded hex colors outside tokens.css', () => {
    const hex = /#(?:[0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})\b/
    expect(offenders(cssFiles, hex, ['styles/tokens.css'])).toEqual([])
  })

  it('no font-family outside base.css', () => {
    expect(offenders(cssFiles, /font-family/, ['styles/base.css'])).toEqual([])
  })
})
