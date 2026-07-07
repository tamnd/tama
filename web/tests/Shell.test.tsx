// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { Shell } from '../src/app/Shell'
import { AuthProvider } from '../src/state/auth'

beforeEach(() => {
  localStorage.clear()
  document.documentElement.removeAttribute('data-theme')
})

afterEach(cleanup)

function renderShell(path = '/') {
  return render(
    <AuthProvider>
      <MemoryRouter initialEntries={[path]}>
        <Shell>
          <p>screen content</p>
        </Shell>
      </MemoryRouter>
    </AuthProvider>,
  )
}

describe('Shell', () => {
  it('ships the landmark structure: banner, labeled nav, main', () => {
    renderShell()
    expect(screen.getByRole('banner')).toBeTruthy()
    expect(screen.getByRole('navigation', { name: 'Main' })).toBeTruthy()
    expect(screen.getByRole('main')).toBeTruthy()
  })

  it('puts the skip link first, pointing at the main region', () => {
    const { container } = renderShell()
    const shell = container.querySelector('.tama-shell')
    const first = shell?.firstElementChild
    expect(first?.tagName).toBe('A')
    expect(first?.getAttribute('href')).toBe('#main')
    expect(first?.textContent).toMatch(/skip to main content/i)
    expect(screen.getByRole('main').id).toBe('main')
  })

  it('renders the seven destinations in the documented order', () => {
    renderShell()
    const nav = screen.getByRole('navigation', { name: 'Main' })
    const labels = Array.from(nav.querySelectorAll('.tama-shell__nav-label')).map(
      (el) => el.textContent,
    )
    expect(labels).toEqual([
      'Learn',
      'Practice',
      'Leaderboards',
      'Quests',
      'Shop',
      'Profile',
      'More',
    ])
  })

  it('marks only the current route active', () => {
    renderShell('/quests')
    const quests = screen.getByRole('link', { name: 'Quests' })
    expect(quests.getAttribute('aria-current')).toBe('page')
    expect(quests.className).toContain('tama-shell__nav-item--active')
    const learn = screen.getByRole('link', { name: 'Learn' })
    expect(learn.getAttribute('aria-current')).toBeNull()
    expect(learn.className).not.toContain('tama-shell__nav-item--active')
  })

  it('keeps the tab bar entries out of the desktop-only rows', () => {
    renderShell()
    const desktopOnly = Array.from(
      document.querySelectorAll('.tama-shell__nav-row--desktop .tama-shell__nav-label'),
    ).map((el) => el.textContent)
    expect(desktopOnly).toEqual(['Practice', 'More'])
  })

  it('carries the stat bar fixtures: course slot, counters, hearts', () => {
    renderShell()
    const banner = screen.getByRole('banner')
    expect(banner.textContent).toContain('Course: English')
    expect(banner.textContent).toContain('Streak: 0 days')
    expect(banner.textContent).toContain('Gems: 500')
    expect(banner.textContent).toContain('Hearts: 5 of 5')
  })

  it('renders the four sidebar widgets as cards', () => {
    renderShell()
    const sidebar = document.querySelector('.tama-shell__sidebar')
    const headers = Array.from(sidebar?.querySelectorAll('.tama-card__header') ?? []).map(
      (el) => el.textContent,
    )
    expect(headers).toEqual(['Streak', 'Gems', 'Daily quests', 'Leaderboards'])
  })
})
