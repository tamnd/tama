// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { CountUp } from '../src/components/CountUp'

afterEach(() => {
  cleanup()
  document.documentElement.classList.remove('reduce-motion')
  delete (window as { matchMedia?: unknown }).matchMedia
})

function stubMatchMedia(matches: boolean) {
  window.matchMedia = ((query: string) => ({
    matches,
    media: query,
    addEventListener() {},
    removeEventListener() {},
  })) as typeof window.matchMedia
}

describe('CountUp', () => {
  it('starts the sweep from zero', () => {
    stubMatchMedia(false)
    render(<CountUp value={120} />)
    expect(screen.getByText('0')).toBeTruthy()
  })

  it('lands on the value and stays there', async () => {
    stubMatchMedia(false)
    render(<CountUp value={12} duration={30} />)
    await screen.findByText('12')
  })

  it('renders the final value immediately under OS reduced motion', () => {
    stubMatchMedia(true)
    render(<CountUp value={340} />)
    expect(screen.getByText('340')).toBeTruthy()
  })

  it('honors the in-app reduce-motion class too', () => {
    document.documentElement.classList.add('reduce-motion')
    render(<CountUp value={99} />)
    expect(screen.getByText('99')).toBeTruthy()
  })
})
