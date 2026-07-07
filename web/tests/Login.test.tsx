// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, describe, expect, it } from 'vitest'
import { Login } from '@/screens/Login/Login'
import { AuthProvider } from '@/state/auth'

afterEach(cleanup)

function renderLogin() {
  return render(
    <AuthProvider>
      <MemoryRouter initialEntries={['/login']}>
        <Login />
      </MemoryRouter>
    </AuthProvider>,
  )
}

describe('Login', () => {
  it('renders the credential fields and a submit button', () => {
    const { container } = renderLogin()
    expect(screen.getByLabelText('Username')).toBeTruthy()
    expect(screen.getByLabelText('Password')).toBeTruthy()
    const submit = container.querySelector('button[type="submit"]')
    expect(submit?.textContent).toBe('Log in')
  })

  it('switches to the register tab on the same screen', () => {
    const { container } = renderLogin()
    fireEvent.click(screen.getByRole('button', { name: 'Sign up' }))
    expect(screen.getByText('Create your profile')).toBeTruthy()
    expect(screen.getByText('At least 8 characters')).toBeTruthy()
    const submit = container.querySelector('button[type="submit"]')
    expect(submit?.textContent).toBe('Create account')
  })
})
