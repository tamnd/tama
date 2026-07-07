// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { Avatar, avatarColor } from '../src/components/Avatar'

afterEach(cleanup)

describe('Avatar', () => {
  it('falls back to the uppercased first letter', () => {
    render(<Avatar name="tama" />)
    const avatar = screen.getByRole('img', { name: 'tama' })
    expect(avatar.textContent).toBe('T')
  })

  it('picks the same brand color for the same name every time', () => {
    const first = render(<Avatar name="tama" />)
    const a = (first.container.firstElementChild as HTMLElement).style.getPropertyValue(
      '--avatar-bg',
    )
    cleanup()
    const second = render(<Avatar name="tama" />)
    const b = (second.container.firstElementChild as HTMLElement).style.getPropertyValue(
      '--avatar-bg',
    )
    expect(a).toBe(b)
    expect(a).toMatch(/^var\(--[\w-]+\)$/)
  })

  it('the hash spreads names across the set', () => {
    const colors = new Set(['tama', 'duo', 'lily', 'oscar', 'zari'].map(avatarColor))
    expect(colors.size).toBeGreaterThan(1)
  })

  it('renders the picture when given one, with the name on the wrapper', () => {
    const { container } = render(<Avatar name="tama" src="/tama.png" />)
    const img = container.querySelector('img') as HTMLImageElement
    expect(img.getAttribute('alt')).toBe('')
    expect(container.querySelector('.tama-avatar__initial')).toBeNull()
    expect(screen.getByRole('img', { name: 'tama' })).toBeTruthy()
  })

  it('sizes through the modifier class', () => {
    const { container } = render(<Avatar name="tama" size={96} />)
    expect(container.querySelector('.tama-avatar--96')).toBeTruthy()
  })
})
