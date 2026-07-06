// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { TapToken } from '@/components/TapToken'

afterEach(cleanup)

describe('TapToken', () => {
  it('renders a pressable chip', () => {
    const onClick = vi.fn()
    render(<TapToken onClick={onClick}>water</TapToken>)
    const chip = screen.getByRole('button', { name: 'water' })
    expect(chip.className).toContain('tama-tap-token')
    fireEvent.click(chip)
    expect(onClick).toHaveBeenCalledTimes(1)
  })

  it('hollow is a non-interactive placeholder', () => {
    const onClick = vi.fn()
    render(
      <TapToken hollow onClick={onClick}>
        water
      </TapToken>,
    )
    const chip = screen.getByText('water') as HTMLButtonElement
    expect(chip.className).toContain('tama-tap-token--hollow')
    expect(chip.disabled).toBe(true)
    fireEvent.click(chip)
    expect(onClick).not.toHaveBeenCalled()
  })
})
