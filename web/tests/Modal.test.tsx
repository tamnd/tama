// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, beforeAll, describe, expect, it, vi } from 'vitest'
import { Button } from '../src/components/Button'
import { Modal } from '../src/components/Modal'

// jsdom ships HTMLDialogElement but not showModal/close yet; stub the two
// methods on the prototype so the component's native path runs unchanged.
beforeAll(() => {
  HTMLDialogElement.prototype.showModal ??= function (this: HTMLDialogElement) {
    this.open = true
  }
  HTMLDialogElement.prototype.close ??= function (this: HTMLDialogElement) {
    this.open = false
  }
})

afterEach(cleanup)

function renderModal(props: Partial<Parameters<typeof Modal>[0]> = {}) {
  const onClose = vi.fn()
  const utils = render(
    <Modal
      open
      onClose={onClose}
      title="Ready?"
      footer={
        <>
          <Button>Start</Button>
          <Button variant="secondary">Not now</Button>
        </>
      }
      {...props}
    >
      <p>Body copy</p>
    </Modal>,
  )
  const dialog = utils.container.querySelector('dialog') as HTMLDialogElement
  return { onClose, dialog, ...utils }
}

describe('Modal', () => {
  it('shows through the native dialog when open', () => {
    const { dialog, rerender, onClose } = renderModal()
    expect(dialog.open).toBe(true)
    expect(screen.getByText('Ready?')).toBeTruthy()
    expect(screen.getByText('Body copy')).toBeTruthy()
    rerender(
      <Modal open={false} onClose={onClose}>
        <p>Body copy</p>
      </Modal>,
    )
    expect(dialog.open).toBe(false)
  })

  it('asks to close on Esc (the cancel event) and stays open itself', () => {
    const { dialog, onClose } = renderModal()
    const cancel = new Event('cancel', { cancelable: true })
    fireEvent(dialog, cancel)
    expect(onClose).toHaveBeenCalledTimes(1)
    // The native close is always vetoed; the open prop drives the dialog.
    expect(cancel.defaultPrevented).toBe(true)
    expect(dialog.open).toBe(true)
  })

  it('closes on a backdrop click but not on a content click', () => {
    const { dialog, onClose } = renderModal()
    fireEvent.click(screen.getByText('Body copy'))
    expect(onClose).not.toHaveBeenCalled()
    // Only backdrop clicks land on the dialog element itself.
    fireEvent.click(dialog)
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('blocking ignores Esc and the backdrop', () => {
    const { dialog, onClose } = renderModal({ blocking: true })
    fireEvent(dialog, new Event('cancel', { cancelable: true }))
    fireEvent.click(dialog)
    expect(onClose).not.toHaveBeenCalled()
    expect(dialog.open).toBe(true)
  })

  it('keeps the footer primary first in tab order, last on screen', () => {
    const { container } = renderModal()
    const footer = container.querySelector('.tama-modal__footer') as HTMLElement
    const buttons = Array.from(footer.querySelectorAll('button')).map((b) => b.textContent)
    // DOM (= tab) order is primary first; column-reverse flips the paint order.
    expect(buttons).toEqual(['Start', 'Not now'])
  })
})
