import { useEffect, useState } from 'react'
import './ThemeToggle.css'

export type ThemePreference = 'light' | 'dark' | 'system'

export const THEME_STORAGE_KEY = 'tama-theme'

const PREFERENCES: ThemePreference[] = ['light', 'dark', 'system']

// The stored choice, defaulting to system when nothing (or garbage) is
// saved. The inline script in index.html reads the same key before first
// paint, so the toggle only has to keep it up to date.
export function readThemePreference(): ThemePreference {
  try {
    const stored = localStorage.getItem(THEME_STORAGE_KEY)
    if (stored === 'light' || stored === 'dark' || stored === 'system') return stored
  } catch {
    // Storage can be walled off (private mode); fall through to system.
  }
  return 'system'
}

function systemPrefersDark(): boolean {
  return window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false
}

// Resolves a preference to the data-theme attribute the token sheets key
// on. Always sets an explicit value, matching the first-paint script.
export function applyTheme(pref: ThemePreference) {
  const dark = pref === 'dark' || (pref === 'system' && systemPrefersDark())
  document.documentElement.dataset.theme = dark ? 'dark' : 'light'
}

// The manual theme switch: light, dark, or follow the OS. The choice lands
// in localStorage so the next load paints right from the inline script;
// once settings exist server-side it will sync there too.
export function ThemeToggle() {
  const [pref, setPref] = useState<ThemePreference>(readThemePreference)

  useEffect(() => {
    applyTheme(pref)
    if (pref !== 'system') return
    // In system mode, follow the OS live instead of waiting for a reload.
    const media = window.matchMedia?.('(prefers-color-scheme: dark)')
    if (!media) return
    const follow = () => applyTheme('system')
    media.addEventListener('change', follow)
    return () => media.removeEventListener('change', follow)
  }, [pref])

  function choose(next: ThemePreference) {
    setPref(next)
    try {
      localStorage.setItem(THEME_STORAGE_KEY, next)
    } catch {
      // No storage, no persistence; the in-page switch still works.
    }
  }

  return (
    <div className="tama-theme-toggle" role="group" aria-label="Theme">
      {PREFERENCES.map((option) => (
        <button
          key={option}
          type="button"
          className={
            pref === option
              ? 'tama-theme-toggle__option tama-theme-toggle__option--on'
              : 'tama-theme-toggle__option'
          }
          aria-pressed={pref === option}
          onClick={() => choose(option)}
        >
          {option}
        </button>
      ))}
    </div>
  )
}
