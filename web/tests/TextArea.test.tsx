// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { TextArea } from '@/components/TextArea'

afterEach(cleanup)

describe('TextArea', () => {
  it('renders two rows by default', () => {
    render(<TextArea aria-label="Answer" />)
    const area = screen.getByRole('textbox', { name: 'Answer' }) as HTMLTextAreaElement
    expect(area.className).toContain('tama-textarea')
    expect(area.rows).toBe(2)
  })

  it('re-measures its height on input', () => {
    render(<TextArea aria-label="Answer" />)
    const area = screen.getByRole('textbox', { name: 'Answer' }) as HTMLTextAreaElement
    fireEvent.input(area, { target: { value: 'line one\nline two\nline three' } })
    expect(area.style.getPropertyValue('--textarea-height')).toMatch(/px$/)
  })

  it('carries aria-invalid for the error border', () => {
    render(<TextArea aria-label="Answer" aria-invalid />)
    expect(screen.getByRole('textbox').getAttribute('aria-invalid')).toBe('true')
  })
})
