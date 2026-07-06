// @vitest-environment jsdom

import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { CharacterGate, ChestNode, PathNode } from '../src/components/PathNode'

afterEach(cleanup)

describe('PathNode', () => {
  it('active renders the star, the START bubble, and presses', () => {
    const onClick = vi.fn()
    render(<PathNode state="active" label="Lesson 1" onClick={onClick} />)
    const node = screen.getByRole('button', { name: 'Lesson 1' })
    expect(node.className).toContain('tama-path-node--active')
    expect(node).not.toHaveProperty('ariaDisabled', 'true')
    expect(node.querySelector('.tama-path-node__bubble--start')).toBeTruthy()
    fireEvent.click(node)
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('completed renders without a bubble and presses', () => {
    const onClick = vi.fn()
    render(<PathNode state="completed" label="Lesson 1" onClick={onClick} />)
    const node = screen.getByRole('button', { name: 'Lesson 1' })
    expect(node.className).toContain('tama-path-node--completed')
    expect(node.querySelector('.tama-path-node__bubble')).toBeNull()
    fireEvent.click(node)
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('locked is aria-disabled and ignores clicks', () => {
    const onClick = vi.fn()
    render(<PathNode state="locked" label="Lesson 9" onClick={onClick} />)
    const node = screen.getByRole('button', { name: 'Lesson 9' })
    expect(node.getAttribute('aria-disabled')).toBe('true')
    fireEvent.click(node)
    expect(onClick).not.toHaveBeenCalled()
  })

  it('legendary keeps its modifier class', () => {
    render(<PathNode state="legendary" label="Challenge" />)
    const node = screen.getByRole('button', { name: 'Challenge' })
    expect(node.className).toContain('tama-path-node--legendary')
  })

  it('draws one ring segment per level with the passed ones filled', () => {
    render(<PathNode state="active" label="Lesson 2" progress={{ passed: 2, total: 5 }} />)
    const node = screen.getByRole('button', { name: 'Lesson 2' })
    expect(node.querySelectorAll('.tama-path-node__ring-seg')).toHaveLength(5)
    expect(node.querySelectorAll('.tama-path-node__ring-seg--passed')).toHaveLength(2)
  })

  it('jumpHere swaps the bubble', () => {
    render(<PathNode state="locked" label="Unit 3 start" jumpHere />)
    const node = screen.getByRole('button', { name: 'Unit 3 start' })
    expect(node.querySelector('.tama-path-node__bubble--jump')).toBeTruthy()
    expect(node.querySelector('.tama-path-node__bubble--start')).toBeNull()
  })
})

describe('ChestNode', () => {
  it('openable chest takes the click', () => {
    const onClick = vi.fn()
    render(<ChestNode label="Chest" onClick={onClick} />)
    fireEvent.click(screen.getByRole('button', { name: 'Chest' }))
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('opened chest is inert', () => {
    const onClick = vi.fn()
    render(<ChestNode label="Chest" opened onClick={onClick} />)
    const chest = screen.getByRole('button', { name: 'Chest' })
    expect(chest.getAttribute('aria-disabled')).toBe('true')
    expect(chest.className).toContain('tama-path-chest--opened')
    fireEvent.click(chest)
    expect(onClick).not.toHaveBeenCalled()
  })
})

describe('CharacterGate', () => {
  it('locked gate is aria-disabled', () => {
    const onClick = vi.fn()
    render(<CharacterGate label="Checkpoint" onClick={onClick} />)
    const gate = screen.getByRole('button', { name: 'Checkpoint' })
    expect(gate.getAttribute('aria-disabled')).toBe('true')
    fireEvent.click(gate)
    expect(onClick).not.toHaveBeenCalled()
  })

  it('passed gate presses', () => {
    const onClick = vi.fn()
    render(<CharacterGate label="Checkpoint" passed onClick={onClick} />)
    const gate = screen.getByRole('button', { name: 'Checkpoint' })
    expect(gate.className).toContain('tama-path-gate--passed')
    fireEvent.click(gate)
    expect(onClick).toHaveBeenCalledOnce()
  })
})
