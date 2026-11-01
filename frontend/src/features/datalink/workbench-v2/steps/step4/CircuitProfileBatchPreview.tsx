import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupProfileKind } from '../../../../../types/studioV2WriteGroup';
import type { BatchPreviewItem, BatchProfileChange } from '../../state/circuitProfileBatch';
import {
  BATCH_PROFILE_ORDER,
  type BatchProfileStatus,
  type BatchTargetExecutionResult,
  type BatchTargetStatus,
} from '../../state/circuitProfileBatchOutcome';

const STATUS_KEY: Record<BatchTargetStatus, string> = {
  saved: 'batch_status_saved', failed: 'batch_status_failed', unknown: 'batch_status_unknown',
  'not-attempted': 'batch_status_not_attempted', unchanged: 'batch_status_unchanged', partial: 'batch_status_partial',
};
const STATUS_TONE: Record<BatchTargetStatus, string> = {
  saved: 'border-[var(--gw-focus)] bg-[var(--gw-success-bg)] text-[var(--gw-success-fg)]',
  failed: 'border-[var(--gw-danger-border)] bg-[var(--gw-danger-bg)] text-[var(--gw-danger-fg)]',
  unknown: 'border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] text-[var(--gw-warning-fg)]',
  partial: 'border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] text-[var(--gw-warning-fg)]',
  unchanged: 'border-[var(--gw-control)] bg-[var(--gw-surface)] text-[var(--gw-text)]',
  'not-attempted': 'border-[var(--gw-control)] bg-[var(--gw-surface)] text-[var(--gw-muted)]',
};

export interface CircuitProfileBatchPreviewProps {
  items: BatchPreviewItem[];
  results: Record<string, BatchTargetExecutionResult>;
  stale: boolean;
  errorMessage?: string;
  hasUnknown: boolean;
  refreshFailed: boolean;
  verifying: boolean;
  onVerify?: () => void;
}

export const CircuitProfileBatchPreview: React.FC<CircuitProfileBatchPreviewProps> = ({
  items, results, stale, errorMessage, hasUnknown, refreshFailed, verifying, onVerify,
}) => {
  const { t } = useTranslation('workbench-v2');
  const badge = (testId: string, status: BatchTargetStatus) => (
    <span data-testid={testId} className={`rounded border px-1.5 py-0.5 text-xs ${STATUS_TONE[status]}`}>
      {t(`step4.basic.circuit.${STATUS_KEY[status]}`)}
    </span>
  );
  const changeText = (kind: WriteGroupProfileKind, change: BatchProfileChange) => {
    const unit = kind === 'instantaneous' ? 's' : 'm';
    if (change.isNew) return `${t('step4.basic.circuit.batch_new_profile')} (${change.newInterval}${unit})`;
    if (change.unchanged) return `${change.oldInterval}${unit}`;
    return `${change.oldInterval}${unit} → ${change.newInterval}${unit}`;
  };
  return (
    <>
      <p className="font-medium text-[var(--gw-text)]">{t('step4.basic.circuit.preview_title')}</p>
      {stale && <p role="alert" className="text-xs text-[var(--gw-warning-fg)]">{t('step4.basic.reason.stale_intent')}</p>}
      <ul className="space-y-1.5 rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] p-2.5">
        {items.map((item) => {
          const result = results[item.circuit_id];
          return (
            <li key={item.circuit_id} className="flex items-center justify-between border-b border-[var(--gw-control)] py-1 last:border-0">
              <div className="space-y-0.5">
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-[var(--gw-text)]">{item.circuit_name || item.circuit_id}</span>
                  {result && badge(`circuit-profile-batch-status-${item.circuit_id}`, result.status)}
                </div>
                {item.selected ? (
                  <div className="space-y-0.5 text-xs text-[var(--gw-accent)]">
                    {BATCH_PROFILE_ORDER.map((kind) => {
                      const change = item[kind];
                      const profile = result?.profiles[kind];
                      if (!change) return null;
                      return (
                        <div key={kind} className="flex items-center gap-2">
                          <span>{t(`step4.basic.circuit.${kind}_label`)}: {changeText(kind, change)}</span>
                          {profile && badge(`circuit-profile-batch-status-${item.circuit_id}-${kind}`, profile.status as BatchProfileStatus)}
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <div data-testid={`circuit-profile-batch-unselected-${item.circuit_id}`} className="text-xs text-[var(--gw-muted)]">
                    {t('step4.basic.circuit.batch_unselected')}
                  </div>
                )}
              </div>
              {item.hasOverride && (
                <span data-testid={`circuit-profile-batch-override-${item.circuit_id}`} className="rounded bg-[var(--gw-warning-bg)] px-1.5 py-0.5 text-xs text-[var(--gw-warning-fg)]">
                  {t('step4.basic.circuit.batch_override')}
                </span>
              )}
            </li>
          );
        })}
      </ul>
      {errorMessage && <p role="alert" className="text-xs text-[var(--gw-danger-fg)]">{errorMessage}</p>}
      {hasUnknown && (
        <div role="alert" data-testid="circuit-profile-batch-unknown-note" className="space-y-1.5 text-xs text-[var(--gw-warning-fg)]">
          <p>{t('step4.basic.circuit.batch_unknown_note')}</p>
          {onVerify && (
            <button type="button" data-testid="circuit-profile-batch-verify" disabled={verifying} onClick={onVerify}
              className="rounded border border-[var(--gw-warning-border)] px-2 py-1 text-xs text-[var(--gw-warning-fg)] disabled:opacity-40">
              {t('step4.basic.circuit.batch_verify')}
            </button>
          )}
        </div>
      )}
      {refreshFailed && (
        <p role="alert" data-testid="circuit-profile-batch-refresh-error" className="text-xs text-[var(--gw-warning-fg)]">
          {t('step4.basic.circuit.batch_refresh_failed')}
        </p>
      )}
    </>
  );
};
