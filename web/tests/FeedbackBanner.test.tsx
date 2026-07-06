// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { FeedbackBanner } from '../src/components/FeedbackBanner'

afterEach(cleanup)

describe('FeedbackBanner', () => {
  it('correct announces itself and defaults to CONTINUE', () => {
    render(<FeedbackBanner kind="correct" title="Nicely done!" meaning="I drink water." />)
    const banner = screen.getByRole('status')
    expect(banner.className).toContain('tama-feedback--correct')
    expect(banner.getAttribute('aria-live')).toBe('polite')
    expect(screen.getByText('Correct.')).toBeTruthy()
    expect(screen.getByText('I drink water.')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Continue' })).toBeTruthy()
  })

  it('incorrect defaults to GOT IT with the cardinal treatment', () => {
    render(<FeedbackBanner kind="incorrect" title="Correct solution:" />)
    expect(screen.getByRole('status').className).toContain('tama-feedback--incorrect')
    expect(screen.getByText('Incorrect.')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Got it' })).toBeTruthy()
  })

  it('runs the action on click and on Enter', () => {
    const onAction = vi.fn()
    render(<FeedbackBanner kind="correct" title="Nice!" onAction={onAction} />)
    fireEvent.click(screen.getByRole('button', { name: 'Continue' }))
    fireEvent.keyDown(window, { key: 'Enter' })
    expect(onAction).toHaveBeenCalledTimes(2)
  })

  it('renders the REPORT placeholder', () => {
    render(<FeedbackBanner kind="correct" title="Nice!" />)
    expect(screen.getByRole('button', { name: 'Report' })).toBeTruthy()
  })
})
