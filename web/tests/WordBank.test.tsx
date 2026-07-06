// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen, within } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { WordBank } from '../src/components/WordBank'

afterEach(cleanup)

const TOKENS = [
  { id: 'w1', text: 'I' },
  { id: 'w2', text: 'drink' },
  { id: 'w3', text: 'water' },
]

function answerLine() {
  return within(screen.getByRole('group', { name: 'Your answer' }))
}

function bank() {
  return within(screen.getByRole('group', { name: 'Word bank' }))
}

function serialized(container: HTMLElement): string {
  const input = container.querySelector('input[type="hidden"]') as HTMLInputElement
  return input.value
}

describe('WordBank', () => {
  it('renders the bank chips and empty dashed slots', () => {
    const { container } = render(<WordBank tokens={TOKENS} />)
    expect(bank().getAllByRole('button')).toHaveLength(3)
    expect(container.querySelectorAll('.tama-word-bank__slot')).toHaveLength(3)
    expect(serialized(container)).toBe('')
  })

  it('clicking a bank chip moves it to the answer line and leaves a hollow twin', () => {
    const { container } = render(<WordBank tokens={TOKENS} />)
    fireEvent.click(bank().getByRole('button', { name: 'drink' }))
    expect(answerLine().getByRole('button', { name: 'drink' })).toBeTruthy()
    const hollow = bank().getByText('drink') as HTMLButtonElement
    expect(hollow.className).toContain('tama-tap-token--hollow')
    expect(hollow.disabled).toBe(true)
    expect(serialized(container)).toBe('w2')
  })

  it('clicking an answer chip returns it to the bank', () => {
    const { container } = render(<WordBank tokens={TOKENS} />)
    fireEvent.click(bank().getByRole('button', { name: 'water' }))
    fireEvent.click(answerLine().getByRole('button', { name: 'water' }))
    expect(answerLine().queryByRole('button')).toBeNull()
    expect(bank().getByRole('button', { name: 'water' })).toBeTruthy()
    expect(serialized(container)).toBe('')
  })

  it('number keys pick bank chips and Backspace returns the last', () => {
    const { container } = render(<WordBank tokens={TOKENS} />)
    const root = container.querySelector('.tama-word-bank') as HTMLElement
    fireEvent.keyDown(root, { key: '3' })
    fireEvent.keyDown(root, { key: '1' })
    expect(serialized(container)).toBe('w3 w1')
    fireEvent.keyDown(root, { key: 'Backspace' })
    expect(serialized(container)).toBe('w3')
  })

  it('ignores number keys for chips already used', () => {
    const { container } = render(<WordBank tokens={TOKENS} />)
    const root = container.querySelector('.tama-word-bank') as HTMLElement
    fireEvent.keyDown(root, { key: '2' })
    fireEvent.keyDown(root, { key: '2' })
    expect(serialized(container)).toBe('w2')
  })

  it('reports the ordered ids through onChange', () => {
    const onChange = vi.fn()
    render(<WordBank tokens={TOKENS} onChange={onChange} />)
    fireEvent.click(bank().getByRole('button', { name: 'I' }))
    fireEvent.click(bank().getByRole('button', { name: 'drink' }))
    expect(onChange).toHaveBeenLastCalledWith(['w1', 'w2'])
  })
})
