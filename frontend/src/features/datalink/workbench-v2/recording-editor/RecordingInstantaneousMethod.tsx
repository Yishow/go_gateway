import { useEffect, useId, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import { RECORDING_PROFILE_BOUNDS } from '../../../../utils/studioV2RecordingProfile';
import { ParameterHelp } from '../components/ParameterHelp';
import type { GroupCandidate } from '../state/writeGroup/candidates';
import type { EditorDraft } from '../state/writeGroup/draft';
import {
  applyInstantaneousInterval, applyInstantaneousReduction, draftReduction, instantaneousIntervalValid, instantaneousMethodChanged,
  instantaneousSaveRejected, parseIntervalText,
} from './instantaneousReduction';

export interface RecordingInstantaneousMethodProps {
  draft: EditorDraft;
  savedGroup?: WriteGroup;
  candidates: GroupCandidate[];
  disabled: boolean;
  /** The editor's last save failure; a configuration rejection is shown at these fields and the draft is kept. */
  saveError?: unknown;
  onReplaceDraft: (next: EditorDraft) => void;
}

const bounds = RECORDING_PROFILE_BOUNDS.instantaneous;

/**
 * Reduction and interval of a saved instantaneous profile record. Both live in the editor's single draft and are
 * saved through its revision-checked save; mean is only written after the operator confirms gauge use.
 */
export function RecordingInstantaneousMethod({ draft, savedGroup, candidates, disabled, saveError, onReplaceDraft }: RecordingInstantaneousMethodProps) {
  const { t } = useTranslation('recording-editor-constraints');
  const { t: editorText } = useTranslation('recording-editor');
  const id = useId();
  const profile = draft.recording_profile;
  const interval = profile && profile.kind === 'instantaneous' ? profile.interval : bounds.default;
  const [text, setText] = useState(String(interval));
  const [consentPending, setConsentPending] = useState(false);
  // An interval replaced from outside (save readback, reload) is shown as saved; typed text that parses to it stays.
  useEffect(() => { setText(current => parseIntervalText(current) === interval ? current : String(interval)); }, [interval]);
  const reduction = draftReduction(draft);
  const intervalValid = instantaneousIntervalValid(draft);
  const changed = instantaneousMethodChanged(draft, savedGroup);
  const rejected = changed && instantaneousSaveRejected(saveError);
  const name = `${id}-reduction`;
  const radio = 'inline-flex min-h-11 items-center gap-2 text-sm text-[var(--gw-text)]';
  return (
    <section className="space-y-3 text-[var(--gw-text)]" data-testid="recording-instant-method">
      <fieldset className="space-y-1" disabled={disabled} aria-describedby={rejected ? `${id}-rejected` : undefined}>
        <legend className="text-sm font-medium">{t('instantaneousMethod.reductionLabel')}</legend>
        <div className="flex flex-wrap gap-x-6">
          <label className={radio}>
            <input type="radio" name={name} checked={reduction === 'latest'} data-testid="recording-instant-reduction-latest"
              onChange={() => { setConsentPending(false); if (!disabled) onReplaceDraft(applyInstantaneousReduction(draft, 'latest', savedGroup, candidates)); }} />
            {t('instantaneousMethod.latest')}
          </label>
          <label className={radio}>
            <input type="radio" name={name} checked={reduction === 'mean'} data-testid="recording-instant-reduction-mean"
              onChange={() => { if (!disabled) setConsentPending(true); }} />
            {t('instantaneousMethod.mean')}
          </label>
        </div>
      </fieldset>
      {consentPending && reduction !== 'mean' && (
        <div role="group" aria-labelledby={`${id}-consent`} className="space-y-2 rounded-lg border border-[var(--gw-control)] bg-[var(--gw-surface)] p-3" data-testid="recording-instant-mean-consent">
          <p id={`${id}-consent`} className="text-sm font-medium">{t('instantaneousMethod.meanConsentTitle')}</p>
          <p className="text-sm">{editorText('compact.gaugeConsent')}</p>
          <div className="flex flex-wrap items-center gap-2">
            <button type="button" className="gw-secondary min-h-11 rounded-md border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-1.5 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-ground)] disabled:opacity-50" disabled={disabled} data-testid="recording-instant-confirm-gauge"
              onClick={() => { if (disabled) return; setConsentPending(false); onReplaceDraft(applyInstantaneousReduction(draft, 'mean', savedGroup, candidates)); }}>
              {editorText('compact.confirmGauge')}
            </button>
            <ParameterHelp compact help="member.gaugeConfirm" />
            <button type="button" className="gw-secondary min-h-11 rounded-md border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-1.5 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-ground)]" data-testid="recording-instant-cancel-mean" onClick={() => setConsentPending(false)}>
              {t('instantaneousMethod.cancelMean')}
            </button>
          </div>
        </div>
      )}
      <div className="space-y-1">
        <div className="flex items-center gap-1">
          <label htmlFor={`${id}-interval`} className="text-sm font-medium">{t('instantaneousMethod.intervalLabel')}</label>
          <ParameterHelp compact help="circuit.instantaneousInterval" />
        </div>
        <input id={`${id}-interval`} type="text" inputMode="numeric" autoComplete="off" value={text} disabled={disabled}
          aria-invalid={!intervalValid} aria-describedby={intervalValid ? undefined : `${id}-interval-error`}
          className="min-h-11 w-28 rounded-md border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-2 text-sm tabular-nums text-[var(--gw-text)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--gw-focus)] aria-[invalid=true]:border-[var(--gw-danger-fg)]"
          data-testid="recording-instant-interval"
          onChange={event => { setText(event.target.value); if (!disabled) onReplaceDraft(applyInstantaneousInterval(draft, parseIntervalText(event.target.value))); }} />
        {!intervalValid && (
          <p id={`${id}-interval-error`} className="text-sm text-[var(--gw-danger-fg)]" data-testid="recording-instant-interval-error">
            {t('instantaneousMethod.intervalInvalid', { min: bounds.min, max: bounds.max })}
          </p>
        )}
      </div>
      {rejected && <p id={`${id}-rejected`} role="alert" className="text-sm text-[var(--gw-danger-fg)]" data-testid="recording-instant-rejected">{t('instantaneousMethod.rejected')}</p>}
      {changed && <p role="status" className="text-sm text-[var(--gw-muted)]" data-testid="recording-instant-pending-note">{t('instantaneousMethod.pendingNote')}</p>}
    </section>
  );
}
