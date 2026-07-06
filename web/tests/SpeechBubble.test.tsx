// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { SpeechBubble } from '@/components/SpeechBubble'

afterEach(cleanup)

describe('SpeechBubble', () => {
  it('renders the prompt text', () => {
    render(<SpeechBubble>Watashi wa mizu o nomimasu.</SpeechBubble>)
    expect(screen.getByText('Watashi wa mizu o nomimasu.')).toBeTruthy()
  })

  it('has no audio button without a handler', () => {
    render(<SpeechBubble>Hello</SpeechBubble>)
    expect(screen.queryByRole('button')).toBeNull()
  })

  it('plays audio from the labeled button before the text', () => {
    const onPlayAudio = vi.fn()
    render(<SpeechBubble onPlayAudio={onPlayAudio}>Hello</SpeechBubble>)
    const button = screen.getByRole('button', { name: 'Play audio' })
    fireEvent.click(button)
    expect(onPlayAudio).toHaveBeenCalledTimes(1)
  })

  it('hides the speaker icon from screen readers', () => {
    render(<SpeechBubble onPlayAudio={() => {}}>Hello</SpeechBubble>)
    const svg = screen.getByRole('button').querySelector('svg')
    expect(svg?.getAttribute('aria-hidden')).toBe('true')
  })
})
