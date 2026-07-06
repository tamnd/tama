import { useState, type CSSProperties, type ReactNode } from 'react'
import { Button, type ButtonVariant } from '../components/Button'
import { Card, CardHeader, Divider } from '../components/Card'
import { ChoiceCard } from '../components/ChoiceCard'
import { FeedbackBanner, type FeedbackKind } from '../components/FeedbackBanner'
import { Field } from '../components/Field'
import { CharacterGate, ChestNode, PathNode } from '../components/PathNode'
import { DEMO_PATH, serpentineOffset } from '../components/PathNode.stories.data'
import { ProgressBar } from '../components/ProgressBar'
import { SpeechBubble } from '../components/SpeechBubble'
import { TapToken } from '../components/TapToken'
import { TextArea } from '../components/TextArea'
import { TextInput } from '../components/TextInput'
import { Toggle } from '../components/Toggle'
import { WordBank } from '../components/WordBank'
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

function PathNodeSection() {
  return (
    <Section
      title="Path nodes"
      path="web/src/components/PathNode.tsx"
      note="Active bobs its START bubble inside the gold ring; the ringed node counts levels passed. Locked ignores the press. The chest shimmers until opened, and the gate carries a silhouette until the Tama poses land."
    >
      <div className="tama-gallery__path-row">
        <PathNode state="active" label="Active node" />
        <PathNode state="completed" label="Completed node" />
        <PathNode state="locked" label="Locked node" />
        <PathNode state="legendary" label="Legendary node" />
        <PathNode
          state="active"
          label="Node with progress ring"
          progress={{ passed: 3, total: 5 }}
        />
        <PathNode state="locked" label="Jump target" jumpHere />
      </div>
      <div className="tama-gallery__path-row">
        <ChestNode label="Openable chest" />
        <ChestNode label="Opened chest" opened />
        <CharacterGate label="Locked gate" />
        <CharacterGate label="Passed gate" passed />
      </div>
      <h3 className="label-caps">Serpentine rhythm</h3>
      <div className="tama-gallery__serpentine">
        {DEMO_PATH.map((entry, i) => (
          <div
            key={i}
            className="tama-gallery__serpentine-row"
            style={{ '--node-offset': `${serpentineOffset(i)}px` } as CSSProperties}
          >
            {entry.kind === 'node' && (
              <PathNode state={entry.state} label={entry.label} progress={entry.progress} />
            )}
            {entry.kind === 'chest' && <ChestNode label={entry.label} opened={entry.opened} />}
            {entry.kind === 'gate' && <CharacterGate label={entry.label} passed={entry.passed} />}
          </div>
        ))}
      </div>
    </Section>
  )
}

const LONG_MEANING =
  'Excuse me, could you bring us two glasses of water and the check, please? ' +
  'We are in a bit of a hurry because our train leaves in twenty minutes.'

function FeedbackBannerSection() {
  const [live, setLive] = useState<FeedbackKind | null>(null)

  return (
    <Section
      title="Feedback banners"
      path="web/src/components/FeedbackBanner.tsx"
      note="The two banners below are pinned in place for review. The buttons mount the real fixed sheet: it slides up from the bottom and Enter dismisses it."
    >
      <div className="tama-gallery__grid">
        <FeedbackBanner
          kind="correct"
          title="Nicely done!"
          meaning="Watashi wa mizu o nomimasu."
          className="tama-gallery__banner-inline"
        />
        <FeedbackBanner
          kind="incorrect"
          title="Correct solution:"
          meaning={LONG_MEANING}
          className="tama-gallery__banner-inline"
        />
      </div>
      <div className="tama-gallery__row">
        <Button variant="primary" size="small" onClick={() => setLive('correct')}>
          Show correct
        </Button>
        <Button variant="danger" size="small" onClick={() => setLive('incorrect')}>
          Show incorrect
        </Button>
      </div>
      {live && (
        <FeedbackBanner
          kind={live}
          title={live === 'correct' ? 'Nicely done!' : 'Correct solution:'}
          meaning="Watashi wa mizu o nomimasu."
          onAction={() => setLive(null)}
        />
      )}
    </Section>
  )
}

function ProgressBarSection() {
  const max = 10
  const flameAt = 5
  const [value, setValue] = useState(3)

  return (
    <Section
      title="Progress bar"
      path="web/src/components/ProgressBar.tsx"
      note="Play increments the bar: each step animates the fill over 300ms and runs one shimmer sweep; the combo flame pops in at 5."
    >
      <div className="tama-gallery__grid">
        <ProgressBar value={value} max={max} showFlame={value >= flameAt} label="Demo progress" />
      </div>
      <div className="tama-gallery__row">
        <Button variant="blue" size="small" onClick={() => setValue((v) => Math.min(max, v + 1))}>
          Play
        </Button>
        <Button variant="secondary" size="small" onClick={() => setValue(0)}>
          Reset
        </Button>
        <code>
          {value} / {max}
        </code>
      </div>
    </Section>
  )
}

function SpeechBubbleSection() {
  return (
    <Section
      title="Speech bubble"
      path="web/src/components/SpeechBubble.tsx"
      note="The tail points left at the character slot; the grey circle stands in until the Tama poses land. The audio button is optional and presses down like a chip."
    >
      <div className="tama-gallery__grid">
        <div className="tama-gallery__bubble-row">
          <span className="tama-gallery__character" aria-hidden="true" />
          <SpeechBubble onPlayAudio={() => {}}>Watashi wa mizu o nomimasu.</SpeechBubble>
        </div>
        <div className="tama-gallery__bubble-row">
          <span className="tama-gallery__character" aria-hidden="true" />
          <SpeechBubble>No audio button, just the prompt in eel at text-l.</SpeechBubble>
        </div>
      </div>
    </Section>
  )
}

const BANK_TOKENS = [
  { id: 'w1', text: 'I' },
  { id: 'w2', text: 'drink' },
  { id: 'w3', text: 'water' },
  { id: 'w4', text: 'the' },
  { id: 'w5', text: 'milk' },
  { id: 'w6', text: 'eat' },
]

function WordBankSection() {
  const [answer, setAnswer] = useState<string[]>([])

  return (
    <Section
      title="Tap tokens and word bank"
      path="web/src/components/WordBank.tsx"
      note="Tap chips or press 1-9 to build the answer, Backspace returns the last chip; moves animate with FLIP and snap under reduced motion. The frozen chips below show the rest and hollow states."
    >
      <div className="tama-gallery__grid">
        <WordBank tokens={BANK_TOKENS} onChange={setAnswer} />
        <code>serialized: {answer.join(' ') || '(empty)'}</code>
      </div>
      <h3 className="label-caps">Frozen states</h3>
      <div className="tama-gallery__row">
        <TapToken>rest</TapToken>
        <TapToken hollow>hollow</TapToken>
      </div>
    </Section>
  )
}

function FormsSection() {
  return (
    <Section
      title="Inputs and forms"
      path="web/src/components/Field.tsx"
      note="Tab through for the focus treatment: the border swaps to macaw with no ring doubling. Field wires label, hint, and error to the control once."
    >
      <div className="tama-gallery__forms">
        <Field label="Name" hint="As it appears on your profile">
          <TextInput placeholder="Type your name" />
        </Field>
        <Field label="Email" error="That does not look like an email">
          <TextInput defaultValue="tama@" />
        </Field>
        <Field label="Disabled">
          <TextInput defaultValue="Cannot touch this" disabled />
        </Field>
        <Field label="Translation" hint="Grows from 2 to 6 rows as you type">
          <TextArea placeholder="Type the translation" />
        </Field>
      </div>
      <h3 className="label-caps">Choice cards</h3>
      <div className="tama-gallery__row">
        <ChoiceCard name="gallery-choice" value="agua" badge={1} defaultChecked>
          el agua
        </ChoiceCard>
        <ChoiceCard name="gallery-choice" value="leche" badge={2}>
          la leche
        </ChoiceCard>
        <ChoiceCard name="gallery-choice" value="pan" badge={3}>
          el pan
        </ChoiceCard>
        <ChoiceCard type="checkbox">checkbox flavor</ChoiceCard>
        <ChoiceCard disabled>disabled</ChoiceCard>
      </div>
      <h3 className="label-caps">Toggles</h3>
      <div className="tama-gallery__row">
        <Toggle>Off</Toggle>
        <Toggle defaultChecked>On</Toggle>
        <Toggle disabled>Disabled</Toggle>
        <Toggle disabled defaultChecked>
          Disabled on
        </Toggle>
      </div>
    </Section>
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

      <PathNodeSection />
      <FeedbackBannerSection />
      <ProgressBarSection />
      <SpeechBubbleSection />
      <WordBankSection />
      <FormsSection />
    </div>
  )
}
