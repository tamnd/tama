import { useEffect, useState } from 'react'

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

export function App() {
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
                style={{ marginLeft: nodeOffset(i) }}
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

        <div className="card">
          <h3>Daily quests</h3>
          <p>Earn 30 XP</p>
          <div className="quest-bar">
            <div className="quest-bar-fill" style={{ width: '0%' }} />
          </div>
        </div>

        <div className="card">
          <h3>Pick a course</h3>
          <p>
            {catalogSize === null
              ? 'Any language to any language.'
              : `${catalogSize} courses from English, or generate any other pair.`}
          </p>
          <button className="btn btn-primary">Get started</button>
        </div>
      </aside>
    </div>
  )
}

// The path snakes left and right as it descends, like a winding trail.
function nodeOffset(i: number): number {
  const wave = [0, -44, -70, -44, 0, 44]
  return wave[i % wave.length]
}
