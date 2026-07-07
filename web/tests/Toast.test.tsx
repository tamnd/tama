// @vitest-environment jsdom

import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { useRef } from 'react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ToastProvider, useToast } from '@/components/Toast'

beforeEach(() => vi.useFakeTimers())
afterEach(() => {
  cleanup()
  vi.useRealTimers()
})

// A minimal consumer: each press fires one numbered toast.
function Demo() {
  const toast = useToast()
  const n = useRef(0)
  return (
    <button type="button" onClick={() => toast(`toast ${++n.current}`)}>
      fire
    </button>
  )
}

function renderDemo() {
  render(
    <ToastProvider>
      <Demo />
    </ToastProvider>,
  )
  return screen.getByRole('button', { name: 'fire' })
}

describe('Toast', () => {
  it('announces through role=status and auto-dismisses after 4s', () => {
    const fire = renderDemo()
    fireEvent.click(fire)
    expect(screen.getByRole('status').textContent).toContain('toast 1')
    act(() => vi.advanceTimersByTime(3999))
    expect(screen.queryByRole('status')).toBeTruthy()
    act(() => vi.advanceTimersByTime(1))
    expect(screen.queryByRole('status')).toBeNull()
  })

  it('dismisses early through the visible button', () => {
    const fire = renderDemo()
    fireEvent.click(fire)
    fireEvent.click(screen.getByRole('button', { name: 'Dismiss' }))
    expect(screen.queryByRole('status')).toBeNull()
  })

  it('stacks at most three, evicting the oldest', () => {
    const fire = renderDemo()
    for (let i = 0; i < 4; i++) fireEvent.click(fire)
    const messages = screen.getAllByRole('status').map((t) => t.textContent)
    expect(messages).toHaveLength(3)
    expect(messages.join(' ')).not.toContain('toast 1')
    expect(messages.join(' ')).toContain('toast 4')
  })

  it('renders the icon slot hidden from screen readers', () => {
    function IconDemo() {
      const toast = useToast()
      return (
        <button type="button" onClick={() => toast('with icon', { icon: <svg /> })}>
          fire
        </button>
      )
    }
    render(
      <ToastProvider>
        <IconDemo />
      </ToastProvider>,
    )
    fireEvent.click(screen.getByRole('button', { name: 'fire' }))
    const icon = screen.getByRole('status').querySelector('.tama-toast__icon') as HTMLElement
    expect(icon.getAttribute('aria-hidden')).toBe('true')
  })
})
