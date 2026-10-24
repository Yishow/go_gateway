import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupCutover } from '../../../../../../types/studioV2WriteGroupCutover';
import { boundaryIsScheduled } from '../../../state/cutoverFlow';
import type { CutoverOutcome } from './useCutoverFlow';

export interface CutoverOutcomeViewProps {
  outcome: CutoverOutcome;
  groupName: (groupId: string) => string;
  busy: boolean;
  readonly: boolean;
  onCheck: () => void;
  onResend: () => void;
}

/** What happened to the one confirmed request: recorded, refused with the old group untouched, or not known. */
export const CutoverOutcomeView: React.FC<CutoverOutcomeViewProps> = ({ outcome, groupName, busy, readonly, onCheck, onResend }) => {
  const { t } = useTranslation('workbench-v2');
  if (outcome.state === 'applied') {
    const { cutover } = outcome;
    return (
      <div data-testid="cutover-record" className="space-y-1 rounded border border-emerald-800 bg-emerald-950/20 p-3 text-xs text-emerald-100">
        <p className="font-semibold">{t('step4.cutover.record_title')}</p>
        <p data-testid={boundaryIsScheduled(cutover.boundary) ? 'cutover-scheduled' : 'cutover-boundary-reached'}>
          {t(boundaryIsScheduled(cutover.boundary) ? 'step4.cutover.scheduled' : 'step4.cutover.boundary_reached', { boundary: cutover.boundary })}
        </p>
        <p>{t('step4.cutover.record_id')} <code data-testid="cutover-record-id">{cutover.id}</code></p>
        <p>{t('step4.cutover.record_boundary')} <code data-testid="cutover-record-boundary">{cutover.boundary}</code></p>
        <p data-testid="cutover-record-started">{t('step4.cutover.record_started')} {cutover.started.map((item) => groupName(item.group_id)).join(', ') || '—'}</p>
        <p data-testid="cutover-record-ended">{t('step4.cutover.record_ended')} {cutover.ended.map((item) => groupName(item.group_id)).join(', ') || '—'}</p>
        <p>{t('step4.cutover.record_moved', { count: cutover.moved.length })}</p>
      </div>
    );
  }
  if (outcome.state === 'failed') {
    return (
      <div role="alert" data-testid="cutover-failed" className="space-y-1 rounded border border-rose-800 bg-rose-950/30 p-3 text-xs text-rose-200">
        <p>{t('step4.cutover.failed')}</p>
        <p>{outcome.message}</p>
      </div>
    );
  }
  if (outcome.state === 'unknown') {
    return (
      <div role="alert" data-testid="cutover-unknown" className="space-y-2 rounded border border-amber-700 bg-amber-950/30 p-3 text-xs text-amber-200">
        <p>{t('step4.cutover.unknown')}</p>
        {outcome.checked === 'none' && <p data-testid="cutover-check-none">{t('step4.cutover.check_none')}</p>}
        {outcome.checked === 'failed' && <p data-testid="cutover-check-failed">{t('step4.cutover.check_failed')}</p>}
        <div className="flex flex-wrap gap-2">
          <button type="button" data-testid="cutover-check" disabled={busy} onClick={onCheck} className="rounded border border-amber-600 px-3 py-1.5 disabled:opacity-40">{t('step4.cutover.check')}</button>
          {outcome.checked === 'none' && <button type="button" data-testid="cutover-resend" disabled={busy || readonly} onClick={onResend} className="rounded border border-amber-600 px-3 py-1.5 disabled:opacity-40">{t('step4.cutover.resend')}</button>}
        </div>
      </div>
    );
  }
  return null;
};

export const CutoverHistory: React.FC<{ cutovers: readonly WriteGroupCutover[]; groupName: (groupId: string) => string }> = ({ cutovers, groupName }) => {
  const { t } = useTranslation('workbench-v2');
  if (cutovers.length === 0) return null;
  return (
    <div data-testid="cutover-history" className="space-y-1 text-xs text-slate-400">
      <p className="font-semibold text-slate-300">{t('step4.cutover.history_title')}</p>
      <ul className="space-y-1">
        {cutovers.map((item) => (
          <li key={item.id} data-testid={`cutover-history-${item.id}`}>
            {t(`step4.cutover.history_${item.kind}`)} · <code>{item.boundary}</code> · {item.started.map((ref) => groupName(ref.group_id)).join(', ') || '—'}
          </li>
        ))}
      </ul>
    </div>
  );
};
