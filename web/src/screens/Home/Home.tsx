import { useEffect, useState, type CSSProperties } from 'react'
import { catalog } from '@/api/client'
import { Button } from '@/components/Button'
import { Card, CardHeader } from '@/components/Card'
import { BookIcon } from '@/components/icons/book'
import { ChestNode, PathNode, type PathNodeState } from '@/components/PathNode'
import './Home.css'

// The learn screen's main column: unit banner, a demo path, and the course
// picker card. The Shell owns the rail, stat bar, and sidebar around it.
// Real path data arrives with M3/M5.

type DemoEntry =
  { kind: 'node'; state: PathNodeState; label: string } | { kind: 'chest'; label: string }

const DEMO_PATH: DemoEntry[] = [
  { kind: 'node', state: 'active', label: 'Lesson 1, current' },
  { kind: 'node', state: 'locked', label: 'Lesson 2, locked' },
  { kind: 'node', state: 'locked', label: 'Lesson 3, locked' },
  { kind: 'chest', label: 'Reward chest' },
  { kind: 'node', state: 'locked', label: 'Lesson 4, locked' },
  { kind: 'node', state: 'legendary', label: 'Legendary challenge' },
]

export function Home() {
  const [catalogSize, setCatalogSize] = useState<number | null>(null)

  useEffect(() => {
    catalog('en')
      .then((courses) => setCatalogSize(courses.length))
      .catch(() => setCatalogSize(null))
  }, [])

  return (
    <div className="tama-home">
      <div className="tama-home__banner">
        <div>
          <div className="tama-home__crumb">Section 1, Unit 1</div>
          <h2>Order food and drink</h2>
        </div>
        <button type="button" className="tama-home__guidebook" aria-label="Guidebook">
          <BookIcon size={24} />
        </button>
      </div>

      <div className="tama-home__path">
        {DEMO_PATH.map((entry, i) => (
          <div
            key={i}
            className="tama-home__path-row"
            style={{ '--node-offset': `${nodeOffset(i)}px` } as CSSProperties}
          >
            {entry.kind === 'chest' ? (
              <ChestNode label={entry.label} />
            ) : (
              <PathNode state={entry.state} label={entry.label} />
            )}
          </div>
        ))}
        <img className="tama-home__mascot" src="/tama.svg" alt="Tama the cat" />
      </div>

      <Card className="tama-home__course">
        <CardHeader>Pick a course</CardHeader>
        <p>
          {catalogSize === null
            ? 'Any language to any language.'
            : `${catalogSize} courses from English, or generate any other pair.`}
        </p>
        <Button variant="primary">Get started</Button>
      </Card>
    </div>
  )
}

// The path snakes left and right as it descends, like a winding trail.
function nodeOffset(i: number): number {
  const wave = [0, -44, -70, -44, 0, 44]
  return wave[i % wave.length] ?? 0
}
