// @vitest-environment jsdom

import { cleanup, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { Field } from '../src/components/Field'
import { TextInput } from '../src/components/TextInput'

afterEach(cleanup)

describe('Field', () => {
  it('associates the label with the control', () => {
    render(
      <Field label="Name">
        <TextInput />
      </Field>,
    )
    expect(screen.getByLabelText('Name')).toBeTruthy()
  })

  it('wires the hint through aria-describedby', () => {
    render(
      <Field label="Name" hint="As it appears on your profile">
        <TextInput />
      </Field>,
    )
    const input = screen.getByLabelText('Name')
    const hint = screen.getByText('As it appears on your profile')
    expect(input.getAttribute('aria-describedby')).toBe(hint.id)
    expect(input.getAttribute('aria-invalid')).toBeNull()
  })

  it('wires the error and flips aria-invalid', () => {
    render(
      <Field label="Name" hint="Pick something short" error="Name is required">
        <TextInput />
      </Field>,
    )
    const input = screen.getByLabelText('Name')
    const hint = screen.getByText('Pick something short')
    const error = screen.getByText('Name is required')
    expect(error.className).toContain('tama-field__error')
    expect(input.getAttribute('aria-describedby')).toBe(`${hint.id} ${error.id}`)
    expect(input.getAttribute('aria-invalid')).toBe('true')
  })
})
