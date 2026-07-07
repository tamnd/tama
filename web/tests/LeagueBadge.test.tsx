// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { LEAGUE_TONES, LeagueBadge } from '../src/components/LeagueBadge'

afterEach(cleanup)

describe('LeagueBadge', () => {
  it('paints the shield through the tone custom property', () => {
    const { container } = render(<LeagueBadge tone="bronze" />)
    const badge = container.firstElementChild as HTMLElement
    expect(badge.style.getPropertyValue('--league-tone')).toBe('var(--league-bronze)')
    expect(badge.querySelector('svg')).toBeTruthy()
  })

  it('names the league for screen readers', () => {
    render(<LeagueBadge tone="gold" />)
    expect(screen.getByText('gold league')).toBeTruthy()
  })

  it('exports the M2 tone list for M7 to extend', () => {
    expect(Object.keys(LEAGUE_TONES)).toEqual(['bronze', 'silver', 'gold'])
    for (const tone of Object.values(LEAGUE_TONES)) {
      expect(tone).toMatch(/^var\(--league-[a-z]+\)$/)
    }
  })
})
