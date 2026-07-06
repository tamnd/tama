// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { ChoiceCard } from '@/components/ChoiceCard'

afterEach(cleanup)

describe('ChoiceCard', () => {
  it('wraps a native radio labeled by the card text', () => {
    render(
      <ChoiceCard name="answer" value="a">
        el agua
      </ChoiceCard>,
    )
    const radio = screen.getByRole('radio', { name: 'el agua' }) as HTMLInputElement
    expect(radio.value).toBe('a')
    fireEvent.click(radio)
    expect(radio.checked).toBe(true)
  })

  it('supports the checkbox type', () => {
    render(<ChoiceCard type="checkbox">Remind me</ChoiceCard>)
    expect(screen.getByRole('checkbox', { name: 'Remind me' })).toBeTruthy()
  })

  it('renders the keyboard badge hidden from screen readers', () => {
    const { container } = render(<ChoiceCard badge={2}>la leche</ChoiceCard>)
    const badge = container.querySelector('.tama-choice__badge') as HTMLElement
    expect(badge.textContent).toBe('2')
    expect(badge.getAttribute('aria-hidden')).toBe('true')
  })

  it('disables the input with the card', () => {
    render(<ChoiceCard disabled>el pan</ChoiceCard>)
    const radio = screen.getByRole('radio', { name: 'el pan' }) as HTMLInputElement
    expect(radio.disabled).toBe(true)
  })
})
