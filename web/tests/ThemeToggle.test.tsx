// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ThemeToggle, THEME_STORAGE_KEY } from '../src/components/ThemeToggle'

// jsdom has no matchMedia; this stub decides what "system" resolves to.
function stubMatchMedia(matches: boolean) {
  vi.stubGlobal(
    'matchMedia',
    vi.fn().mockReturnValue({
      matches,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    }),
  )
}

beforeEach(() => {
  localStorage.clear()
  document.documentElement.removeAttribute('data-theme')
})

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('ThemeToggle', () => {
  it('defaults to system and applies the OS preference', () => {
    stubMatchMedia(true)
    render(<ThemeToggle />)
    const system = screen.getByRole('button', { name: 'system' })
    expect(system.getAttribute('aria-pressed')).toBe('true')
    expect(document.documentElement.dataset.theme).toBe('dark')
  })

  it('applies and persists an explicit dark choice', () => {
    stubMatchMedia(false)
    render(<ThemeToggle />)
    fireEvent.click(screen.getByRole('button', { name: 'dark' }))
    expect(document.documentElement.dataset.theme).toBe('dark')
    expect(localStorage.getItem(THEME_STORAGE_KEY)).toBe('dark')
  })

  it('light beats a dark OS preference', () => {
    stubMatchMedia(true)
    render(<ThemeToggle />)
    fireEvent.click(screen.getByRole('button', { name: 'light' }))
    expect(document.documentElement.dataset.theme).toBe('light')
    expect(localStorage.getItem(THEME_STORAGE_KEY)).toBe('light')
  })

  it('restores the stored preference on mount', () => {
    stubMatchMedia(false)
    localStorage.setItem(THEME_STORAGE_KEY, 'dark')
    render(<ThemeToggle />)
    expect(screen.getByRole('button', { name: 'dark' }).getAttribute('aria-pressed')).toBe('true')
    expect(document.documentElement.dataset.theme).toBe('dark')
  })

  it('going back to system follows the OS again', () => {
    stubMatchMedia(false)
    localStorage.setItem(THEME_STORAGE_KEY, 'dark')
    render(<ThemeToggle />)
    fireEvent.click(screen.getByRole('button', { name: 'system' }))
    expect(document.documentElement.dataset.theme).toBe('light')
    expect(localStorage.getItem(THEME_STORAGE_KEY)).toBe('system')
  })
})
