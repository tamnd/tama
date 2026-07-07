import { useRef, useState, type CSSProperties, type ReactNode } from 'react'
import { Avatar } from '@/components/Avatar'
import { Button, type ButtonVariant } from '@/components/Button'
import { Card, CardHeader, Divider } from '@/components/Card'
import { ChoiceCard } from '@/components/ChoiceCard'
import { CountUp } from '@/components/CountUp'
import { FeedbackBanner, type FeedbackKind } from '@/components/FeedbackBanner'
import { Field } from '@/components/Field'
import { HeartsRow } from '@/components/HeartsRow'
import { BellIcon } from '@/components/icons/bell'
import { BookIcon } from '@/components/icons/book'
import { CheckIcon } from '@/components/icons/check'
import { ChestIcon } from '@/components/icons/chest'
import { CrownIcon } from '@/components/icons/crown'
import { DumbbellIcon } from '@/components/icons/dumbbell'
import { FlagReportIcon } from '@/components/icons/flag-report'
import { GearIcon } from '@/components/icons/gear'
import { GemIcon } from '@/components/icons/gem'
import { HeartIcon } from '@/components/icons/heart'
import { HomePathIcon } from '@/components/icons/home-path'
import type { IconSize } from '@/components/icons/icon'
import { LightningIcon } from '@/components/icons/lightning'
import { LockIcon } from '@/components/icons/lock'
import { MicIcon } from '@/components/icons/mic'
import { MoreDotsIcon } from '@/components/icons/more-dots'
import { ProfileIcon } from '@/components/icons/profile'
import { QuestScrollIcon } from '@/components/icons/quest-scroll'
import { ShieldIcon } from '@/components/icons/shield'
import { ShopBagIcon } from '@/components/icons/shop-bag'
import { SpeakerIcon } from '@/components/icons/speaker'
import { SpeakerSlowIcon } from '@/components/icons/speaker-slow'
import { StarIcon } from '@/components/icons/star'
import { StreakFlameIcon } from '@/components/icons/streak-flame'
import { TrophyIcon } from '@/components/icons/trophy'
import { XIcon } from '@/components/icons/x'
import { LeagueBadge, LEAGUE_TONES, type LeagueTone } from '@/components/LeagueBadge'
import { Modal } from '@/components/Modal'
import { CharacterGate, ChestNode, PathNode } from '@/components/PathNode'
import { DEMO_PATH, serpentineOffset } from '@/components/PathNode.stories.data'
import { Popover } from '@/components/Popover'
import { ProgressBar } from '@/components/ProgressBar'
import { SpeechBubble } from '@/components/SpeechBubble'
import { CrownChip, StatChip } from '@/components/StatChip'
import { TapToken } from '@/components/TapToken'
import { TextArea } from '@/components/TextArea'
import { TextInput } from '@/components/TextInput'
import { ToastProvider, useToast } from '@/components/Toast'
import { Toggle } from '@/components/Toggle'
import { WordBank } from '@/components/WordBank'
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

const ICON_SIZES: IconSize[] = [16, 20, 24, 32]

// The full set, gamification tones first, then the single-tone chrome.
const ICON_SET: [string, (size: IconSize) => ReactNode][] = [
  ['streak-flame', (s) => <StreakFlameIcon size={s} />],
  ['gem', (s) => <GemIcon size={s} />],
  ['heart', (s) => <HeartIcon size={s} />],
  ['lightning', (s) => <LightningIcon size={s} />],
  ['crown', (s) => <CrownIcon size={s} />],
  ['chest', (s) => <ChestIcon size={s} />],
  ['chest open', (s) => <ChestIcon size={s} open />],
  ['trophy', (s) => <TrophyIcon size={s} />],
  ['star', (s) => <StarIcon size={s} />],
  ['check', (s) => <CheckIcon size={s} />],
  ['x', (s) => <XIcon size={s} />],
  ['lock', (s) => <LockIcon size={s} />],
  ['book', (s) => <BookIcon size={s} />],
  ['dumbbell', (s) => <DumbbellIcon size={s} />],
  ['shield', (s) => <ShieldIcon size={s} />],
  ['quest-scroll', (s) => <QuestScrollIcon size={s} />],
  ['bell', (s) => <BellIcon size={s} />],
  ['gear', (s) => <GearIcon size={s} />],
  ['flag-report', (s) => <FlagReportIcon size={s} />],
  ['speaker', (s) => <SpeakerIcon size={s} />],
  ['speaker-slow', (s) => <SpeakerSlowIcon size={s} />],
  ['mic', (s) => <MicIcon size={s} />],
  ['home-path', (s) => <HomePathIcon size={s} />],
  ['profile', (s) => <ProfileIcon size={s} />],
  ['shop-bag', (s) => <ShopBagIcon size={s} />],
  ['more-dots', (s) => <MoreDotsIcon size={s} />],
]

function IconPanel({ variant }: { variant: 'light' | 'dark' }) {
  return (
    <div className={`tama-gallery__icon-panel tama-gallery__icon-panel--${variant}`}>
      {ICON_SET.map(([name, draw]) => (
        <div key={name} className="tama-gallery__icon-cell">
          <span className="tama-gallery__icon-sizes">
            {ICON_SIZES.map((size) => (
              <span key={size}>{draw(size)}</span>
            ))}
          </span>
          <code>{name}</code>
        </div>
      ))}
    </div>
  )
}

function IconsSection() {
  return (
    <Section
      title="Icons"
      path="web/src/components/icons/"
      note="Every icon at 16, 20, 24, and 32 on a light and a dark swatch. Gamification icons keep their canonical fills on both; chrome icons inherit the swatch's currentColor."
    >
      <IconPanel variant="light" />
      <IconPanel variant="dark" />
    </Section>
  )
}

function StatSection() {
  const [streak, setStreak] = useState(12)
  const [gems, setGems] = useState(505)
  const [xp, setXp] = useState(120)
  const [hearts, setHearts] = useState(5)
  const [replay, setReplay] = useState(0)

  return (
    <Section
      title="Stat displays"
      path="web/src/components/StatChip.tsx"
      note="Counters sweep up through CountUp and announce themselves on polite live labels. Lose a heart to watch it pop and hollow; the unlimited row is the Super state."
    >
      <h3 className="label-caps">Stat chips</h3>
      <div className="tama-gallery__row">
        <StatChip kind="streak" value={streak} />
        <StatChip kind="streak" value={0} />
        <StatChip kind="gems" value={gems} />
        <StatChip kind="xp" value={xp} />
        <Button
          variant="secondary"
          size="small"
          onClick={() => {
            setStreak((v) => v + 1)
            setGems((v) => v + 30)
            setXp((v) => v + 15)
          }}
        >
          Earn a day
        </Button>
      </div>
      <h3 className="label-caps">Hearts</h3>
      <div className="tama-gallery__row">
        <HeartsRow remaining={hearts} />
        <Button variant="danger" size="small" onClick={() => setHearts((v) => Math.max(0, v - 1))}>
          Lose a heart
        </Button>
        <Button variant="secondary" size="small" onClick={() => setHearts(5)}>
          Refill
        </Button>
      </div>
      <div className="tama-gallery__row">
        <HeartsRow remaining={3} />
        <HeartsRow remaining={0} />
        <HeartsRow remaining={5} unlimited />
      </div>
      <h3 className="label-caps">Avatars</h3>
      <div className="tama-gallery__row">
        <Avatar name="Tama" size={32} />
        <Avatar name="Tama" size={48} />
        <Avatar name="Tama" size={96} />
        <Avatar name="Duo" size={48} />
        <Avatar name="Lily" size={48} />
        <Avatar name="Oscar" size={48} />
        <Avatar name="Zari" size={48} />
      </div>
      <h3 className="label-caps">League badges and crown levels</h3>
      <div className="tama-gallery__row">
        {(Object.keys(LEAGUE_TONES) as LeagueTone[]).map((tone) => (
          <LeagueBadge key={tone} tone={tone} />
        ))}
        <CrownChip level={1} />
        <CrownChip level={12} />
      </div>
      <h3 className="label-caps">Count-up</h3>
      <div className="tama-gallery__row">
        <span className="tama-gallery__countup">
          <CountUp key={replay} value={340} />
        </span>
        <Button variant="blue" size="small" onClick={() => setReplay((n) => n + 1)}>
          Replay
        </Button>
      </div>
    </Section>
  )
}

// The app will mount one ToastProvider at its root; the gallery scopes one
// to this section so the demo stays self-contained.
function OverlaySection() {
  return (
    <ToastProvider>
      <OverlayDemos />
    </ToastProvider>
  )
}

function OverlayDemos() {
  const [modal, setModal] = useState(false)
  const [blocking, setBlocking] = useState(false)
  const [popovers, setPopovers] = useState(true)
  const greenAnchor = useRef<HTMLDivElement>(null)
  const neutralAnchor = useRef<HTMLDivElement>(null)
  const toast = useToast()
  const count = useRef(0)

  function fireToast() {
    count.current += 1
    toast(`Streak extended! Day ${count.current}`, { icon: <StreakFlameIcon size={20} /> })
  }

  return (
    <Section
      title="Overlays"
      path="web/src/components/Modal.tsx"
      note="The modal closes on Esc and backdrop click unless blocking; its footer stacks buttons with the primary last on screen but first in tab order. Popovers hang under their anchor nodes and never take focus. Toasts stack top-center, at most three, and leave after 4s."
    >
      <div className="tama-gallery__row">
        <Button variant="blue" size="small" onClick={() => setModal(true)}>
          Open modal
        </Button>
        <Button variant="secondary" size="small" onClick={() => setBlocking(true)}>
          Open blocking modal
        </Button>
        <Button variant="primary" size="small" onClick={fireToast}>
          Fire toast
        </Button>
        <Button
          variant="secondary"
          size="small"
          onClick={() => {
            fireToast()
            fireToast()
            fireToast()
            fireToast()
          }}
        >
          Fire four toasts
        </Button>
        <Button variant="secondary" size="small" onClick={() => setPopovers((v) => !v)}>
          {popovers ? 'Hide popovers' : 'Show popovers'}
        </Button>
      </div>

      <Modal
        open={modal}
        onClose={() => setModal(false)}
        title="Ready to practice?"
        footer={
          <>
            <Button onClick={() => setModal(false)}>Start lesson</Button>
            <Button variant="secondary" onClick={() => setModal(false)}>
              Not now
            </Button>
          </>
        }
      >
        <p>
          A quick five-minute session keeps the streak alive. Esc or the backdrop closes this one.
        </p>
      </Modal>

      <Modal
        open={blocking}
        onClose={() => setBlocking(false)}
        blocking
        title="Hearts are out"
        footer={
          <Button variant="danger" onClick={() => setBlocking(false)}>
            Got it
          </Button>
        }
      >
        <p>Blocking: Esc and backdrop clicks are ignored, only the button gets out.</p>
      </Modal>

      <h3 className="label-caps">Popover variants</h3>
      <div className="tama-gallery__popover-row">
        <div ref={greenAnchor} className="tama-gallery__popover-anchor">
          <PathNode state="active" label="Active node with popover" />
        </div>
        <div ref={neutralAnchor} className="tama-gallery__popover-anchor">
          <PathNode state="locked" label="Locked node with popover" />
        </div>
      </div>
      {popovers && (
        <>
          <Popover anchorRef={greenAnchor} variant="green">
            <strong>Order food and drink</strong>
            <br />
            Lesson 1 of 4
          </Popover>
          <Popover anchorRef={neutralAnchor}>Complete the level above to unlock this!</Popover>
        </>
      )}
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
      <OverlaySection />
      <IconsSection />
      <StatSection />
    </div>
  )
}
