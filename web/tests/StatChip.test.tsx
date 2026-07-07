// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { CrownChip, StatChip } from '../src/components/StatChip'

afterEach(cleanup)

describe('StatChip', () => {
  it('announces the streak through a polite live label', () => {
    render(<StatChip kind="streak" value={12} />)
    const label = screen.getByText('Streak: 12 days')
    expect(label.getAttribute('aria-live')).toBe('polite')
    expect(label.className).toContain('visually-hidden')
  })

  it('uses the singular for a one-day streak', () => {
    render(<StatChip kind="streak" value={1} />)
    expect(screen.getByText('Streak: 1 day')).toBeTruthy()
  })

  it('greys the whole chip at zero days', () => {
    const { container } = render(<StatChip kind="streak" value={0} />)
    expect(container.querySelector('.tama-statchip--zero')).toBeTruthy()
  })

  it('a live streak is not the zero state', () => {
    const { container } = render(<StatChip kind="streak" value={5} />)
    expect(container.querySelector('.tama-statchip--zero')).toBeNull()
  })

  it('labels gems and XP with their own phrases', () => {
    render(<StatChip kind="gems" value={505} />)
    expect(screen.getByText('Gems: 505')).toBeTruthy()
    render(<StatChip kind="xp" value={40} />)
    expect(screen.getByText('XP: 40')).toBeTruthy()
  })

  it('hides the icon and the animated number from screen readers', () => {
    const { container } = render(<StatChip kind="gems" value={9} />)
    const icon = container.querySelector('.tama-statchip__icon') as HTMLElement
    const value = container.querySelector('.tama-statchip__value') as HTMLElement
    expect(icon.getAttribute('aria-hidden')).toBe('true')
    expect(value.getAttribute('aria-hidden')).toBe('true')
  })
})

describe('CrownChip', () => {
  it('wears the crown and announces the level', () => {
    const { container } = render(<CrownChip level={7} />)
    expect(container.querySelector('.tama-crownchip svg')).toBeTruthy()
    const label = screen.getByText('Crown level: 7')
    expect(label.getAttribute('aria-live')).toBe('polite')
  })
})
