// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { TextInput } from '@/components/TextInput'

afterEach(cleanup)

describe('TextInput', () => {
  it('renders a textbox with the placeholder', () => {
    render(<TextInput placeholder="Type your answer" />)
    const input = screen.getByPlaceholderText('Type your answer')
    expect(input.className).toContain('tama-text-input')
  })

  it('passes native props through', () => {
    render(<TextInput aria-label="Name" defaultValue="Tama" disabled />)
    const input = screen.getByRole('textbox', { name: 'Name' }) as HTMLInputElement
    expect(input.value).toBe('Tama')
    expect(input.disabled).toBe(true)
  })

  it('carries aria-invalid for the error border', () => {
    render(<TextInput aria-label="Name" aria-invalid />)
    const input = screen.getByRole('textbox', { name: 'Name' })
    expect(input.getAttribute('aria-invalid')).toBe('true')
  })
})
