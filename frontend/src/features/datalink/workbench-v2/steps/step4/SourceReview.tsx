import * as React from 'react';
import { useTranslation } from 'react-i18next';

export interface SourceReviewProps {
  /** Test and element prefix, e.g. `circuit-profile-sources`. */
  prefix: string;
  /** Every source of the review, including the unavailable ones. */
  count: number;
  /** Sources that need a look: saved members no current source matches, or records that cannot be verified. */
  unavailable: number;
  /** The summary text for `unavailable`; saved members that no longer match by default. */
  exceptionKey?: string;
  /** Opens the saved group that holds them; the summary never edits anything itself. */
  onRepair?: () => void;
  /** A new explicit request reveals and focuses the picker without changing its selection. */
  expandRequest?: number;
  /** The rows. They are rendered only while expanded. */
  children: React.ReactNode;
}

/**
 * A source list that starts as its count and its exceptions. The unavailable sources and their repair stay visible
 * while collapsed; expanding only reveals the rows and changes nothing. The expanded list scrolls inside itself.
 */
export const SourceReview: React.FC<SourceReviewProps> = ({ prefix, count, unavailable, exceptionKey = 'step4.basic.sources.unavailable', onRepair, expandRequest = 0, children }) => {
  const { t } = useTranslation('workbench-v2');
  const [expanded, setExpanded] = React.useState(false);
  const toggleRef = React.useRef<HTMLButtonElement>(null);
  React.useEffect(() => {
    if (expandRequest > 0) {
      setExpanded(true);
      toggleRef.current?.focus();
    }
  }, [expandRequest]);
  const listId = `${prefix}-list`;
  return (
    <div className="space-y-1.5">
      <div className="flex flex-wrap items-center justify-between gap-2 text-xs text-slate-300">
        <p data-testid={`${prefix}-summary`}>{t('step4.basic.sources.count', { count })}</p>
        <button ref={toggleRef} type="button" data-testid={`${prefix}-toggle`} aria-expanded={expanded} aria-controls={listId} onClick={() => setExpanded((open) => !open)}
          className="rounded border border-slate-700 px-2 py-0.5 text-[11px] text-slate-300 hover:bg-slate-800">
          {t(`step4.basic.sources.${expanded ? 'hide' : 'show'}`)}
        </button>
      </div>
      {unavailable > 0 && (
        <div role="status" data-testid={`${prefix}-exceptions`} className="flex flex-wrap items-center gap-2 text-xs text-amber-300">
          <span>{t(exceptionKey, { count: unavailable })}</span>
          {onRepair && (
            <button type="button" data-testid={`${prefix}-repair`} onClick={onRepair} className="rounded border border-amber-600/60 px-2 py-0.5 text-[11px]">
              {t('step4.basic.sources.repair')}
            </button>
          )}
        </div>
      )}
      <div id={listId} hidden={!expanded}>{expanded ? children : null}</div>
    </div>
  );
};
