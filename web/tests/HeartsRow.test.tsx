// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { HeartsRow } from '../src/components/HeartsRow'

afterEach(cleanup)

describe('HeartsRow', () => {
  it('renders five hearts, all filled at full health', () => {
    const { container } = render(<HeartsRow remaining={5} />)
    expect(container.querySelectorAll('.tama-hearts__heart').length).toBe(5)
    expect(container.querySelectorAll('.tama-hearts__heart--spent').length).toBe(0)
    expect(screen.getByText('Hearts: 5 of 5')).toBeTruthy()
  })

  it('hollows the spent hearts', () => {
    const { container } = render(<HeartsRow remaining={2} />)
    expect(container.querySelectorAll('.tama-hearts__heart--spent').length).toBe(3)
    expect(screen.getByText('Hearts: 2 of 5')).toBeTruthy()
  })

  it('clamps out-of-range values', () => {
    const { container } = render(<HeartsRow remaining={-2} />)
    expect(container.querySelectorAll('.tama-hearts__heart--spent').length).toBe(5)
    expect(screen.getByText('Hearts: 0 of 5')).toBeTruthy()
  })

  it('marks a live loss so the pop and hollow play', () => {
    const { container, rerender } = render(<HeartsRow remaining={3} />)
    expect(container.querySelector('.tama-hearts__heart--losing')).toBeNull()
    rerender(<HeartsRow remaining={2} />)
    const losing = container.querySelectorAll('.tama-hearts__heart--losing')
    expect(losing.length).toBe(1)
    expect(losing[0].className).toContain('tama-hearts__heart--spent')
  })

  it('a refill clears the losing marks', () => {
    const { container, rerender } = render(<HeartsRow remaining={3} />)
    rerender(<HeartsRow remaining={2} />)
    rerender(<HeartsRow remaining={5} />)
    expect(container.querySelector('.tama-hearts__heart--losing')).toBeNull()
    expect(container.querySelectorAll('.tama-hearts__heart--spent').length).toBe(0)
  })

  it('unlimited swaps the row for one infinity heart', () => {
    const { container } = render(<HeartsRow remaining={5} unlimited />)
    expect(container.querySelector('.tama-hearts--unlimited')).toBeTruthy()
    expect(container.querySelectorAll('.tama-hearts__heart').length).toBe(1)
    expect(container.querySelectorAll('svg path').length).toBe(2)
    expect(screen.getByText('Hearts: unlimited')).toBeTruthy()
  })
})
