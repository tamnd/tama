// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { Toggle } from '@/components/Toggle'

afterEach(cleanup)

describe('Toggle', () => {
  it('is a native checkbox labeled by its text', () => {
    render(<Toggle>Sound effects</Toggle>)
    const input = screen.getByRole('checkbox', { name: 'Sound effects' }) as HTMLInputElement
    expect(input.checked).toBe(false)
    fireEvent.click(input)
    expect(input.checked).toBe(true)
  })

  it('respects defaultChecked and disabled', () => {
    render(<Toggle defaultChecked disabled aria-label="Dark mode" />)
    const input = screen.getByRole('checkbox', { name: 'Dark mode' }) as HTMLInputElement
    expect(input.checked).toBe(true)
    expect(input.disabled).toBe(true)
  })

  it('hides the track art from screen readers', () => {
    const { container } = render(<Toggle aria-label="Sound" />)
    const track = container.querySelector('.tama-toggle__track') as HTMLElement
    expect(track.getAttribute('aria-hidden')).toBe('true')
    expect(track.querySelector('.tama-toggle__knob')).toBeTruthy()
  })
})
