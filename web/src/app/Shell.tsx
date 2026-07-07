import type { ReactNode } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import { logout } from '@/api/client'
import { Button } from '@/components/Button'
import { Card, CardHeader } from '@/components/Card'
import { HeartsRow } from '@/components/HeartsRow'
import { DumbbellIcon } from '@/components/icons/dumbbell'
import { HomePathIcon } from '@/components/icons/home-path'
import { MoreDotsIcon } from '@/components/icons/more-dots'
import { ProfileIcon } from '@/components/icons/profile'
import { QuestScrollIcon } from '@/components/icons/quest-scroll'
import { ShieldIcon } from '@/components/icons/shield'
import { ShopBagIcon } from '@/components/icons/shop-bag'
import { ProgressBar } from '@/components/ProgressBar'
import { StatChip } from '@/components/StatChip'
import { ThemeToggle } from '@/components/ThemeToggle'
import { useAuth } from '@/state/auth'
import './Shell.css'

// The three-region app frame: nav rail on the left, the centered content
// unit (592px main, 48px gutter, 368px sidebar) in the middle, and the
// sidebar widgets on the right. Under 700px (see breakpoints.ts) the rail
// turns into the bottom tab bar; under 1264px the sidebar folds into the
// top stat bar. One nav element plays both rail and tab bar, so there is
// exactly one navigation landmark.

type NavItem = {
  to: string
  label: string
  icon: ReactNode
  /** Rail-only entries drop out of the mobile tab bar. */
  desktopOnly?: boolean
}

const NAV_ITEMS: NavItem[] = [
  { to: '/', label: 'Learn', icon: <HomePathIcon size={32} /> },
  { to: '/practice', label: 'Practice', icon: <DumbbellIcon size={32} />, desktopOnly: true },
  { to: '/leaderboards', label: 'Leaderboards', icon: <ShieldIcon size={32} /> },
  { to: '/quests', label: 'Quests', icon: <QuestScrollIcon size={32} /> },
  { to: '/shop', label: 'Shop', icon: <ShopBagIcon size={32} /> },
  { to: '/profile', label: 'Profile', icon: <ProfileIcon size={32} /> },
  { to: '/more', label: 'More', icon: <MoreDotsIcon size={32} />, desktopOnly: true },
]

// Fixture numbers for the M2 shell; M3 onward wires the real profile in.
const STATS = { streak: 0, gems: 500, hearts: 5 }

export function Shell({ children }: { children: ReactNode }) {
  return (
    <div className="tama-shell">
      <a className="tama-shell__skip" href="#main">
        Skip to main content
      </a>

      <nav className="tama-shell__rail" aria-label="Main">
        <NavLink to="/" className="tama-shell__logo">
          <img src="/tama.svg" alt="" />
          tama
        </NavLink>
        <ul className="tama-shell__nav-list">
          {NAV_ITEMS.map((item) => (
            <li key={item.to} className={item.desktopOnly ? 'tama-shell__nav-row--desktop' : ''}>
              <NavLink
                to={item.to}
                end={item.to === '/'}
                className={({ isActive }) =>
                  isActive
                    ? 'tama-shell__nav-item tama-shell__nav-item--active'
                    : 'tama-shell__nav-item'
                }
              >
                <span className="tama-shell__nav-icon" aria-hidden="true">
                  {item.icon}
                </span>
                <span className="tama-shell__nav-label">{item.label}</span>
              </NavLink>
            </li>
          ))}
        </ul>
        <MoreArea />
      </nav>

      <div className="tama-shell__center">
        <header className="tama-shell__statbar">
          <span className="tama-shell__flag">
            <span aria-hidden="true">EN</span>
            <span className="visually-hidden">Course: English</span>
          </span>
          <StatChip kind="streak" value={STATS.streak} />
          <StatChip kind="gems" value={STATS.gems} />
          <HeartsRow remaining={STATS.hearts} />
        </header>

        <div className="tama-shell__content">
          <main id="main" className="tama-shell__main">
            {children}
          </main>
          <aside className="tama-shell__sidebar" aria-label="Progress">
            <SidebarWidget title="Streak">
              <div className="tama-shell__widget-row">
                <StatChip kind="streak" value={STATS.streak} />
                <p>Do a lesson today to start a streak.</p>
              </div>
            </SidebarWidget>
            <SidebarWidget title="Gems">
              <div className="tama-shell__widget-row">
                <StatChip kind="gems" value={STATS.gems} />
                <p>Earn gems in lessons, spend them in the shop.</p>
              </div>
            </SidebarWidget>
            <SidebarWidget title="Daily quests">
              <p>Earn 30 XP</p>
              <ProgressBar value={0} max={30} label="Earn 30 XP progress" />
            </SidebarWidget>
            <SidebarWidget title="Leaderboards">
              <div className="tama-shell__widget-row">
                <span className="tama-shell__widget-shield" aria-hidden="true">
                  <ShieldIcon size={32} />
                </span>
                <p>Complete 10 more lessons to start competing.</p>
              </div>
            </SidebarWidget>
          </aside>
        </div>
      </div>
    </div>
  )
}

// One sidebar card: caps header over whatever the widget carries. M3+
// swaps the fixture children for live data without touching the frame.
export function SidebarWidget({ title, children }: { title: string; children: ReactNode }) {
  return (
    <Card className="tama-shell__widget">
      <CardHeader>{title}</CardHeader>
      {children}
    </Card>
  )
}

// The rail's MORE area: theme switch and the session exit. Hidden on
// mobile with the rest of the rail chrome.
function MoreArea() {
  const { signOut } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    void (async () => {
      try {
        await logout()
      } catch {
        // A dead session is already logged out.
      }
      signOut()
      await navigate('/login', { replace: true })
    })()
  }

  return (
    <div className="tama-shell__more">
      <ThemeToggle />
      <Button variant="secondary" size="small" onClick={handleLogout}>
        Log out
      </Button>
    </div>
  )
}
