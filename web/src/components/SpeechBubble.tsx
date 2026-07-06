import type { ReactNode } from 'react'
import { SpeakerIcon } from './icons/speaker'
import './SpeechBubble.css'

type SpeechBubbleProps = {
  /** The exercise prompt. */
  children: ReactNode
  /** Renders the audio button before the text when set. */
  onPlayAudio?: () => void
  /** Accessible name for the audio button. */
  audioLabel?: string
  className?: string
}

// The character's line in an exercise: a snow bubble with a border-built
// tail pointing left at the character slot, and an optional audio button
// ahead of the prompt.
export function SpeechBubble({
  children,
  onPlayAudio,
  audioLabel = 'Play audio',
  className,
}: SpeechBubbleProps) {
  const classes = ['tama-speech-bubble']
  if (className) classes.push(className)

  return (
    <div className={classes.join(' ')}>
      {onPlayAudio && (
        <button
          type="button"
          className="tama-speech-bubble__audio"
          aria-label={audioLabel}
          onClick={onPlayAudio}
        >
          <SpeakerIcon size={24} />
        </button>
      )}
      <p className="tama-speech-bubble__text">{children}</p>
    </div>
  )
}
