import { lazy, Suspense, useEffect, useState, type CSSProperties } from 'react'
import { Button } from './components/Button'
import { Card, CardHeader } from './components/Card'
import './app.css'

// The gallery is compiled in for dev builds and behind the __GALLERY__ flag;
// production builds drop the chunk entirely.
const galleryEnabled = import.meta.env.DEV || __GALLERY__
const Gallery = galleryEnabled ? lazy(() => import('./dev/Gallery')) : null

export function App() {
  if (Gallery && window.location.pathname === '/dev/gallery') {
    return (
      <Suspense fallback={null}>
        <Gallery />
      </Suspense>
    )
  }
  return <Home />
}

// Static shell for the foundation milestone: the three-column layout, the nav
// rail, a demo path, and the sidebar widgets. Real data arrives with M3/M5.

type Course = {
  id: string
  base: { code: string; name: string }
  target: { code: string; name: string; native: string }
}

const NAV = [
  { icon: '\u{1F3E0}', label: 'Learn', active: true },
  { icon: '\u{1F4AA}', label: 'Practice', active: false },
  { icon: '\u{1F6E1}️', label: 'Leaderboards', active: false },
  { icon: '\u{1F3AF}', label: 'Quests', active: false },
  { icon: '\u{1F6CD}️', label: 'Shop', active: false },
  { icon: '\u{1F464}', label: 'Profile', active: false },
] as const

type NodeState = 'active' | 'done' | 'locked' | 'chest'

const DEMO_PATH: { state: NodeState; icon: string }[] = [
  { state: 'active', icon: '★' },
  { state: 'locked', icon: '★' },
  { state: 'locked', icon: '\u{1F4D6}' },
  { state: 'chest', icon: '\u{1F381}' },
  { state: 'locked', icon: '★' },
  { state: 'locked', icon: '\u{1F3C6}' },
]

function Home() {
  const [catalogSize, setCatalogSize] = useState<number | null>(null)

  useEffect(() => {
    fetch('/api/catalog?from=en')
      .then((r) => r.json())
      .then((courses: Course[]) => setCatalogSize(courses.length))
      .catch(() => setCatalogSize(null))
  }, [])

  return (
    <div className="layout">
      <nav className="nav">
        <div className="nav-logo">
          <img src="/tama.svg" alt="" />
          tama
        </div>
        {NAV.map((item) => (
          <button key={item.label} className={`nav-item${item.active ? ' active' : ''}`}>
            <span className="icon">{item.icon}</span>
            {item.label.toUpperCase()}
          </button>
        ))}
      </nav>

      <main className="path">
        <div className="unit-banner">
          <div>
            <div className="crumb">Section 1, Unit 1</div>
            <h2>Order food and drink</h2>
          </div>
          <button className="guidebook" title="Guidebook">
            {'\u{1F4D3}'}
          </button>
        </div>

        <div className="nodes">
          {DEMO_PATH.map((node, i) => (
            <div className="node-row" key={i}>
              {node.state === 'active' && <div className="start-bubble">START</div>}
              <button
                className={`node ${node.state}`}
                disabled={node.state === 'locked'}
                style={{ '--node-offset': `${nodeOffset(i)}px` } as CSSProperties}
              >
                {node.state === 'chest' ? '\u{1F381}' : node.icon}
              </button>
            </div>
          ))}
          <img className="mascot" src="/tama.svg" alt="Tama the cat" />
        </div>
      </main>

      <aside className="sidebar">
        <div className="stat-row">
          <span className="stat streak">
            <span className="icon">{'\u{1F525}'}</span> 0
          </span>
          <span className="stat gems">
            <span className="icon">{'\u{1F48E}'}</span> 500
          </span>
          <span className="stat hearts">
            <span className="icon">{'❤️'}</span> 5
          </span>
        </div>

        <Card>
          <CardHeader>Daily quests</CardHeader>
          <p>Earn 30 XP</p>
          <div className="quest-bar">
            <div className="quest-bar-fill" style={{ '--quest-progress': '0%' } as CSSProperties} />
          </div>
        </Card>

        <Card>
          <CardHeader>Pick a course</CardHeader>
          <p>
            {catalogSize === null
              ? 'Any language to any language.'
              : `${catalogSize} courses from English, or generate any other pair.`}
          </p>
          <Button variant="primary">Get started</Button>
        </Card>
      </aside>
    </div>
  )
}

// The path snakes left and right as it descends, like a winding trail.
function nodeOffset(i: number): number {
  const wave = [0, -44, -70, -44, 0, 44]
  return wave[i % wave.length] ?? 0
}
