import { useState, type CSSProperties, type ReactNode } from 'react'
import { Button, type ButtonVariant } from '../components/Button'
import { Card, CardHeader, Divider } from '../components/Card'
import './Gallery.css'

// Eyeball QA for the design system: one section per component, every state
// rendered from fixture data. Dev-only, see the __GALLERY__ flag in App.

const BRAND_TOKENS = [
  '--feather-green',
  '--mask-green',
  '--macaw',
  '--cardinal',
  '--bee',
  '--fox',
  '--beetle',
  '--humpback',
  '--whale',
]

const GREY_TOKENS = ['--eel', '--wolf', '--hare', '--swan', '--polar', '--snow']

const DEEP_TOKENS = [
  '--feather-green-deep',
  '--macaw-deep',
  '--cardinal-deep',
  '--bee-deep',
  '--beetle-deep',
  '--swan-deep',
]

const FILL_TOKENS = ['--correct-fill', '--incorrect-fill']

const TYPE_SCALE = ['--text-xs', '--text-s', '--text-m', '--text-l', '--text-xl', '--text-2xl']

const BUTTON_VARIANTS: ButtonVariant[] = ['primary', 'blue', 'secondary', 'danger', 'super']

const RADII = ['--radius-s', '--radius-m', '--radius-l', '--radius-pill']

// Values come off the live stylesheet, so the swatches always show what the
// active theme actually resolves, never a copy that can drift.
function tokenValue(name: string): string {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

function Swatch({ name }: { name: string }) {
  return (
    <li className="tama-gallery__swatch">
      <span
        className="tama-gallery__swatch-chip"
        style={{ '--swatch': `var(${name})` } as CSSProperties}
      />
      <code>{name}</code>
      <code className="tama-gallery__swatch-value">{tokenValue(name)}</code>
    </li>
  )
}

function SwatchGroup({ title, tokens }: { title: string; tokens: string[] }) {
  return (
    <div className="tama-gallery__group">
      <h3 className="label-caps">{title}</h3>
      <ul className="tama-gallery__swatches">
        {tokens.map((name) => (
          <Swatch key={name} name={name} />
        ))}
      </ul>
    </div>
  )
}

function Section({
  title,
  path,
  note,
  children,
}: {
  title: string
  path: string
  note?: string
  children: ReactNode
}) {
  return (
    <section className="tama-gallery__section">
      <header className="tama-gallery__section-header">
        <h2>{title}</h2>
        <code>{path}</code>
      </header>
      {note && <p className="tama-gallery__note">{note}</p>}
      {children}
    </section>
  )
}

const starIcon = (
  <svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
    <path d="M12 2l2.9 6.3 6.9.8-5.1 4.7 1.4 6.8L12 17.2 5.9 20.6l1.4-6.8L2.2 9.1l6.9-.8z" />
  </svg>
)

export default function Gallery() {
  const [theme, setTheme] = useState<'light' | 'dark'>(() =>
    document.documentElement.dataset.theme === 'dark' ? 'dark' : 'light',
  )

  function toggleTheme() {
    const next = theme === 'light' ? 'dark' : 'light'
    if (next === 'dark') {
      document.documentElement.dataset.theme = 'dark'
    } else {
      delete document.documentElement.dataset.theme
    }
    setTheme(next)
  }

  return (
    <div className="tama-gallery">
      <header className="tama-gallery__header">
        <h1>Component gallery</h1>
        <Button variant="secondary" size="small" onClick={toggleTheme}>
          {theme === 'light' ? 'Dark theme' : 'Light theme'}
        </Button>
      </header>

      <Section title="Tokens" path="web/src/styles/tokens.css">
        <SwatchGroup title="Brand" tokens={BRAND_TOKENS} />
        <SwatchGroup title="Greys" tokens={GREY_TOKENS} />
        <SwatchGroup title="Depth edges" tokens={DEEP_TOKENS} />
        <SwatchGroup title="Feedback fills" tokens={FILL_TOKENS} />
      </Section>

      <Section title="Typography" path="web/src/styles/base.css">
        {TYPE_SCALE.map((name) => (
          <p
            key={name}
            className="tama-gallery__specimen"
            style={{ '--specimen-size': `var(${name})` } as CSSProperties}
          >
            <code>
              {name}: {tokenValue(name)}
            </code>
            The quick brown fox jumps over the lazy dog
          </p>
        ))}
        <p className="tama-gallery__specimen">
          <code>weights</code>
          <span className="tama-gallery__w-body">Body 700</span>{' '}
          <span className="tama-gallery__w-heading">Heading 800</span>{' '}
          <span className="tama-gallery__w-black">Celebration 900</span>
        </p>
        <p className="tama-gallery__specimen">
          <code>.label-caps</code>
          <span className="label-caps">Section header label</span>
        </p>
      </Section>

      <Section
        title="Buttons"
        path="web/src/components/Button.tsx"
        note="Hover, active, and focus are live states: point, press, and tab through the grid to check them. Press one and watch the row below it, nothing should shift."
      >
        <div className="tama-gallery__grid">
          {BUTTON_VARIANTS.map((variant) => (
            <div key={variant} className="tama-gallery__row">
              <Button variant={variant}>{variant}</Button>
              <Button variant={variant} disabled>
                Disabled
              </Button>
            </div>
          ))}
          <div className="tama-gallery__row">
            <Button size="small">Small 40</Button>
            <Button>Default 50</Button>
            <Button size="large">Large 58</Button>
          </div>
          <div className="tama-gallery__row">
            <Button variant="blue" icon={starIcon}>
              With icon
            </Button>
          </div>
          <Button variant="primary" size="large" fullWidth>
            Full width
          </Button>
        </div>
      </Section>

      <Section title="Cards" path="web/src/components/Card.tsx">
        <div className="tama-gallery__grid">
          <Card>
            <CardHeader>Default card</CardHeader>
            <p>Snow surface, 2px swan border, 16px radius, 24px padding.</p>
            <Divider />
            <p>The divider above bleeds across the card padding.</p>
          </Card>
          <Card variant="well">
            <p>Inset well: polar fill, no border, 12px radius.</p>
          </Card>
          <Card variant="clickable" onClick={() => {}}>
            <CardHeader>Clickable card</CardHeader>
            <p>Presses down like a secondary button.</p>
          </Card>
        </div>
        <h3 className="label-caps">Radius reference</h3>
        <div className="tama-gallery__row">
          {RADII.map((name) => (
            <span
              key={name}
              className="tama-gallery__radius-chip"
              style={{ '--chip-radius': `var(${name})` } as CSSProperties}
            >
              {name.replace('--radius-', '')} {tokenValue(name)}
            </span>
          ))}
        </div>
      </Section>
    </div>
  )
}
