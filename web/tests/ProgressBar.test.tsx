// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { ProgressBar } from '../src/components/ProgressBar'

afterEach(cleanup)

describe('ProgressBar', () => {
  it('exposes the progressbar role with the aria values', () => {
    render(<ProgressBar value={3} max={10} label="Lesson progress" />)
    const bar = screen.getByRole('progressbar', { name: 'Lesson progress' })
    expect(bar.getAttribute('aria-valuenow')).toBe('3')
    expect(bar.getAttribute('aria-valuemin')).toBe('0')
    expect(bar.getAttribute('aria-valuemax')).toBe('10')
  })

  it('drives the fill width through the custom property', () => {
    const { container } = render(<ProgressBar value={5} max={10} />)
    const fill = container.querySelector('.tama-progress__fill') as HTMLElement
    expect(fill.style.getPropertyValue('--progress')).toBe('50%')
  })

  it('clamps the fill to the track', () => {
    const { container } = render(<ProgressBar value={15} max={10} />)
    const fill = container.querySelector('.tama-progress__fill') as HTMLElement
    expect(fill.style.getPropertyValue('--progress')).toBe('100%')
  })

  it('runs one shimmer sweep per increment', () => {
    const { container, rerender } = render(<ProgressBar value={2} max={10} />)
    expect(container.querySelector('.tama-progress__shimmer')).toBeNull()
    rerender(<ProgressBar value={3} max={10} />)
    expect(container.querySelector('.tama-progress__shimmer')).toBeTruthy()
  })

  it('shows the combo flame only when the caller says so', () => {
    const { container, rerender } = render(<ProgressBar value={4} max={10} />)
    expect(container.querySelector('.tama-progress__flame')).toBeNull()
    rerender(<ProgressBar value={5} max={10} showFlame />)
    expect(container.querySelector('.tama-progress__flame')).toBeTruthy()
  })
})
