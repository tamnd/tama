// @vitest-environment jsdom

import { cleanup, render } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { BellIcon } from '../src/components/icons/bell'
import { CrownIcon } from '../src/components/icons/crown'
import { GemIcon } from '../src/components/icons/gem'
import { HeartIcon } from '../src/components/icons/heart'
import { HomePathIcon } from '../src/components/icons/home-path'
import { LightningIcon } from '../src/components/icons/lightning'
import { MicIcon } from '../src/components/icons/mic'
import { ShieldIcon } from '../src/components/icons/shield'
import { StreakFlameIcon } from '../src/components/icons/streak-flame'

afterEach(cleanup)

function svgOf(ui: React.ReactElement): SVGSVGElement {
  const { container } = render(ui)
  return container.querySelector('svg') as SVGSVGElement
}

describe('icons', () => {
  it('draw on the 24 grid and default to 24px', () => {
    const svg = svgOf(<GemIcon />)
    expect(svg.getAttribute('viewBox')).toBe('0 0 24 24')
    expect(svg.getAttribute('width')).toBe('24')
    expect(svg.getAttribute('height')).toBe('24')
  })

  it('render at the size the prop asks for', () => {
    expect(svgOf(<ShieldIcon size={32} />).getAttribute('width')).toBe('32')
    expect(svgOf(<MicIcon size={16} />).getAttribute('height')).toBe('16')
    expect(svgOf(<HomePathIcon size={20} />).getAttribute('width')).toBe('20')
  })

  it('are hidden from screen readers', () => {
    for (const svg of [svgOf(<BellIcon />), svgOf(<CrownIcon />), svgOf(<HeartIcon />)]) {
      expect(svg.getAttribute('aria-hidden')).toBe('true')
    }
  })

  it('gamification icons carry their canonical token fills', () => {
    expect(svgOf(<GemIcon />).innerHTML).toContain('var(--macaw)')
    expect(svgOf(<GemIcon />).innerHTML).toContain('var(--humpback)')
    expect(svgOf(<HeartIcon />).innerHTML).toContain('var(--cardinal)')
    expect(svgOf(<LightningIcon />).innerHTML).toContain('var(--bee)')
    expect(svgOf(<CrownIcon />).innerHTML).toContain('var(--bee)')
    const flame = svgOf(<StreakFlameIcon />).innerHTML
    expect(flame).toContain('var(--bee)')
    expect(flame).toContain('var(--fox)')
  })

  it('chrome icons are single tone on currentColor', () => {
    for (const svg of [svgOf(<BellIcon />), svgOf(<ShieldIcon />), svgOf(<HomePathIcon />)]) {
      expect(svg.outerHTML).toContain('currentColor')
      expect(svg.outerHTML).not.toContain('var(--')
    }
  })

  it('heart hollows and stamps the infinity on request', () => {
    const filled = svgOf(<HeartIcon />)
    expect(filled.querySelectorAll('path').length).toBe(1)
    const hollow = svgOf(<HeartIcon hollow />)
    const shape = hollow.querySelector('path') as SVGPathElement
    expect(shape.getAttribute('fill')).toBe('transparent')
    expect(shape.getAttribute('stroke')).toContain('var(--swan)')
    const unlimited = svgOf(<HeartIcon infinity />)
    expect(unlimited.querySelectorAll('path').length).toBe(2)
  })
})
