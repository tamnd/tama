// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { createRef } from 'react'
import { afterEach, describe, expect, it } from 'vitest'
import { Popover } from '../src/components/Popover'

afterEach(cleanup)

function renderPopover(variant?: 'neutral' | 'green') {
  const anchorRef = createRef<HTMLDivElement>()
  const utils = render(
    <>
      <div ref={anchorRef}>anchor</div>
      <Popover anchorRef={anchorRef} variant={variant}>
        Lesson 1 of 4
      </Popover>
    </>,
  )
  const popover = utils.container.querySelector('.tama-popover') as HTMLElement
  return { popover, ...utils }
}

describe('Popover', () => {
  it('renders the neutral variant with a centered tail', () => {
    const { popover } = renderPopover()
    expect(popover.className).toContain('tama-popover--neutral')
    expect(popover.querySelector('.tama-popover__tail')).toBeTruthy()
    expect(screen.getByText('Lesson 1 of 4')).toBeTruthy()
  })

  it('renders the green variant with the same tail', () => {
    const { popover } = renderPopover('green')
    expect(popover.className).toContain('tama-popover--green')
    const tail = popover.querySelector('.tama-popover__tail') as HTMLElement
    expect(tail.getAttribute('aria-hidden')).toBe('true')
  })

  it('falls back to the measured anchor rect without CSS anchor support', () => {
    // jsdom has no anchor positioning, so this is always the JS path.
    const { popover } = renderPopover()
    expect(popover.className).toContain('tama-popover--measured')
    expect(popover.style.top).not.toBe('')
    expect(popover.style.left).not.toBe('')
  })

  it('never steals focus', () => {
    render(<button type="button">stay focused</button>)
    const button = screen.getByRole('button', { name: 'stay focused' })
    button.focus()
    renderPopover('green')
    expect(document.activeElement).toBe(button)
  })
})
