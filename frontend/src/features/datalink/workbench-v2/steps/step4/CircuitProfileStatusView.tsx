import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type {
  WriteGroupDelivery,
  WriteGroupDeliveryStages,
  WriteGroupSourceEvidenceItem,
  WriteGroupSourceRecord,
} from '../../../../../types/studioV2WriteGroupDelivery';
import { formatPreciseUtc } from '../../../../../utils/preciseInstant';
import { SourceReview } from './SourceReview';

export interface CircuitProfileStatusViewProps {
  group: WriteGroup;
  delivery: WriteGroupDelivery | undefined;
  /** Display names by tag ID; a source without one is shown by its tag ID. */
  sourceLabels?: Record<string, string>;
  /** Opens this saved group in Advanced, where a source that cannot be named or verified is repaired. */
  onOpenAdvanced?: (groupId: string) => void;
}

/** Reasons the backend can attach to a record. Anything else is shown as unavailable, never as raw text. */
const KNOWN_REASONS = new Set([
  'no-samples', 'missing-value', 'max-age-exceeded', 'type-mismatch', 'mapping-failed', 'counter-unavailable',
  'quality-bad', 'quality-missing', 'quality-stale', 'quality-invalid', 'quality-uncertain', 'first-reading', 'binding-changed',
  'no-eligible-reading', 'decrease-unproven', 'contradictory-same-time', 'non-monotonic-time', 'mixed-numeric-basis',
  'rollover-evidence-insufficient', 'rollover-range-invalid', 'rollover-multiple-wraps-possible', 'rollover-rate-bound-violated',
  'non-numeric-reading', 'delta-not-representable', 'coverage-unverified', 'bad', 'invalid', 'stale', 'non-numeric',
]);
const KNOWN_STATUSES = new Set(['ok', 'partial', 'baseline', 'missing', 'bad', 'stale', 'invalid', 'reset_or_uncertain', 'gap']);
/** Newest records that need a look: one that cannot be verified, or one that stored a failure. */
const ABNORMAL_STATUSES = new Set(['bad', 'stale', 'invalid', 'reset_or_uncertain']);
const needsCheck = (item: WriteGroupSourceEvidenceItem) => Boolean(item.latest && (!item.latest.available || ABNORMAL_STATUSES.has(item.latest.status)));

const base = 'step4.basic.status';
const locallyDurable = (s: WriteGroupDeliveryStages) => s.queued + s.retrying + s.blocked + s.quarantined + s.unknown;

export const CircuitProfileStatusView: React.FC<CircuitProfileStatusViewProps> = ({ group, delivery, sourceLabels = {}, onOpenAdvanced }) => {
  const { t } = useTranslation('workbench-v2');
  const applied = group.applied_revision;
  const answer = delivery && delivery.group_id === group.id ? delivery : undefined;
  const stages = answer?.revision_stages?.group_revision === applied ? answer.revision_stages.stages : undefined;
  const evidence = answer?.source_evidence ?? undefined;
  const current = evidence?.state === 'available' && evidence.group_revision === applied;
  const reasonText = (reason: string | undefined) => (reason ? t(`${base}.${KNOWN_REASONS.has(reason) ? `reason.${reason}` : 'reason_unavailable'}`) : '');

  const record = (item: WriteGroupSourceEvidenceItem, kind: 'latest' | 'confirmed', data: WriteGroupSourceRecord | WriteGroupSourceEvidenceItem['latest']) => {
    if (!data) return null;
    const tag = item.tag_id;
    const testId = `circuit-status-record-${kind}-${tag}`;
    const heading = <span className="font-medium text-slate-300">{t(`${base}.${kind === 'latest' ? 'latest_window' : 'confirmed_window'}`)}</span>;
    // This legacy view only describes raw-profile slots. Mixed calculation
    // evidence stays read-only/unavailable here until its dedicated presentation.
    if (!data.available || data.profile_kind === 'mixed') {
      return <div data-testid={testId} className="text-slate-400">{heading} {t(`${base}.record_unavailable`)}</div>;
    }
    const foreignRevision = data.group_revision !== applied;
    // Proof is judged against the destination the applied revision froze, never the draft's.
    const foreignDestination = data.connector_revision !== evidence?.connector_revision;
    const cumulative = data.profile_kind === 'cumulative';
    const reason = reasonText(data.reason);
    const status = KNOWN_STATUSES.has(data.status) ? data.status : 'ok';
    const confirmed = data.delivery === 'sql_committed' && !foreignDestination && !foreignRevision;
    return (
      <div data-testid={testId} className="space-y-0.5">
        {heading}
        <div className="text-slate-400">{formatPreciseUtc(data.window_start)} → {formatPreciseUtc(data.window_end)} UTC</div>
        {foreignRevision || foreignDestination ? <div className="text-amber-300">{t(`${base}.destination_changed`)}</div> : (
          <>
            <div>{t(`${base}.status_value.${status}`)} · {t(`${base}.${data.raw_present ? 'raw_recorded' : 'null_stored'}`)}</div>
            {data.aggregation && (
              <div data-testid={`circuit-status-aggregation-${tag}`} className="space-y-0.5 text-slate-400">
                <div>{t(`${base}.method.${data.aggregation.method}`)} · {t(`${base}.coverage.${data.aggregation.coverage}`)}
                  {data.aggregation.reason && ` · ${reasonText(data.aggregation.reason)}`}</div>
                <div>{t(`${base}.valid_label`)}: <span data-testid={`circuit-status-valid-${tag}`}>{data.aggregation.valid_count}</span>
                  {' · '}{t(`${base}.rejected_label`)}: <span data-testid={`circuit-status-rejected-${tag}`}>{data.aggregation.rejected_count}</span></div>
                <div>{data.aggregation.source_type} → {data.aggregation.result_type}</div>
              </div>
            )}
            {cumulative
              ? <div data-testid={`circuit-status-delta-${tag}`}>{t(`${base}.${data.delta_known ? 'difference_recorded' : 'difference_unknown'}`)}{reason && ` · ${reason}`}</div>
              : reason && <div>{reason}</div>}
            <div>{t(`${base}.delivery.${data.delivery}`)}{confirmed && <> · {t(`${base}.receipt_confirmed`)}</>}</div>
          </>
        )}
      </div>
    );
  };

  const row = (item: WriteGroupSourceEvidenceItem) => {
    const sameRecord = item.latest?.available && item.last_sql_committed?.available && item.latest.effect_key === item.last_sql_committed.effect_key;
    return (
      <li key={item.tag_id} data-testid={`circuit-status-source-${item.tag_id}`} className="space-y-1 rounded border border-slate-900 bg-slate-950/40 px-2.5 py-1.5">
        {sourceLabels[item.tag_id]
          ? <div className="font-medium text-slate-200">{sourceLabels[item.tag_id]}</div>
          : (
            // No saved name or address: the technical key stays in diagnostics; the saved member is repaired in its group.
            <div className="flex flex-wrap items-center gap-2 font-medium text-amber-300">
              <span>{t('step4.basic.saved_member_unavailable')}</span>
              {onOpenAdvanced && (
                <button type="button" data-testid={`circuit-status-source-repair-${item.tag_id}`} onClick={() => onOpenAdvanced(group.id)}
                  className="rounded border border-amber-600/60 px-2 py-0.5 text-[11px] font-normal">{t('step4.basic.sources.repair')}</button>
              )}
            </div>
          )}
        {!item.latest && !item.last_sql_committed && <div className="text-slate-400">{t(`${base}.source_none`)}</div>}
        {record(item, 'latest', item.latest)}
        {!sameRecord && record(item, 'confirmed', item.last_sql_committed)}
      </li>
    );
  };

  const chip = (id: string, label: string, value: string) => (
    <li data-testid={`circuit-status-stage-${id}`} className="rounded border border-slate-800 px-2 py-0.5">{t(`${base}.${label}`)}: {value}</li>
  );
  return (
    <div className="space-y-2 text-xs text-slate-300">
      <ul className="flex flex-wrap gap-2 text-[11px]">
        {chip('saved', 'stage_saved', t(`${base}.state_done`))}
        {chip('applied', 'stage_applied', t(`${base}.${applied ? 'state_done' : 'state_pending'}`))}
        {chip('collecting', 'stage_collecting', stages ? String(stages.collecting) : '—')}
        {chip('local', 'stage_local', stages ? String(locallyDurable(stages)) : '—')}
        {chip('sql', 'stage_sql', stages ? String(stages.sql_committed) : '—')}
      </ul>
      {answer?.storage?.history_may_be_lost && (
        <div role="alert" data-testid="circuit-status-storage-history" className="rounded border border-amber-700/60 bg-amber-950/30 p-2 text-amber-200">
          <p data-testid="circuit-status-storage-generation">{t(`${base}.storage_generation`)}: {answer.storage.generation}</p>
          <p>{t(`${base}.storage_history_uncertain`)}</p>
        </div>
      )}
      {stages && stages.sql_committed > 0 && answer?.last_sql_committed_effect?.group_revision === applied && (
        <p data-testid="circuit-status-sql-confirmed">{t(`${base}.receipt_confirmed`)}</p>
      )}
      {applied && <p data-testid="circuit-status-readback" className="text-slate-400">{t(`${base}.readback_not_checked`)}</p>}
      {evidence && !current && (
        <p data-testid="circuit-status-evidence-note" className="text-amber-300">
          {t(`${base}.${evidence.state === 'unavailable' ? 'evidence_unavailable' : 'evidence_not_current'}`)}
        </p>
      )}
      {current && evidence && (evidence.connector_id !== group.destination.connector_id || evidence.connector_revision !== group.destination.connector_revision) && (
        <p data-testid="circuit-status-applied-destination-note" className="text-slate-400">{t(`${base}.applied_destination_note`)}</p>
      )}
      {current && evidence && (
        <SourceReview prefix={`circuit-status-sources-${group.id}`} count={evidence.sources.length} unavailable={evidence.sources.filter(needsCheck).length}
          exceptionKey="step4.basic.sources.needs_check" onRepair={onOpenAdvanced ? () => onOpenAdvanced(group.id) : undefined}>
          <ul data-testid={`circuit-status-source-list-${group.id}`} className="max-h-72 space-y-1.5 overflow-y-auto">{evidence.sources.map(row)}</ul>
        </SourceReview>
      )}
    </div>
  );
};
