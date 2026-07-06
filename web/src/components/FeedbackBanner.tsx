import { useEffect } from 'react'
import { Button } from './Button'
import { CheckIcon } from './icons/check'
import { XIcon } from './icons/x'
import './FeedbackBanner.css'

export type FeedbackKind = 'correct' | 'incorrect'

type FeedbackBannerProps = {
  kind: FeedbackKind
  /** Headline string from the caller, e.g. "Nicely done!". */
  title: string
  /** The correct answer or translation, rendered selectable. */
  meaning?: string
  /** Defaults to CONTINUE on correct, GOT IT on incorrect. */
  actionLabel?: string
  /** Enter triggers this too while the banner is up. */
  onAction?: () => void
  className?: string
}

// The bottom sheet that grades an answer. Slides up over 300ms, announces
// itself politely, and hands the room to one big button.
export function FeedbackBanner({
  kind,
  title,
  meaning,
  actionLabel,
  onAction,
  className,
}: FeedbackBannerProps) {
  useEffect(() => {
    if (!onAction) return
    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Enter') onAction?.()
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [onAction])

  const correct = kind === 'correct'
  const classes = ['tama-feedback', `tama-feedback--${kind}`]
  if (className) classes.push(className)

  return (
    <div className={classes.join(' ')} role="status" aria-live="polite">
      <div className="tama-feedback__inner">
        <span className="tama-feedback__roundel" aria-hidden="true">
          {correct ? <CheckIcon size={24} /> : <XIcon size={24} />}
        </span>
        <div className="tama-feedback__text">
          <h2 className="tama-feedback__title">
            <span className="visually-hidden">{correct ? 'Correct.' : 'Incorrect.'}</span>
            {title}
          </h2>
          {meaning && <p className="tama-feedback__meaning">{meaning}</p>}
        </div>
        <button type="button" className="tama-feedback__report">
          Report
        </button>
        <div className="tama-feedback__action">
          <Button
            variant={correct ? 'primary' : 'danger'}
            size="large"
            fullWidth
            onClick={onAction}
          >
            {actionLabel ?? (correct ? 'Continue' : 'Got it')}
          </Button>
        </div>
      </div>
    </div>
  )
}
