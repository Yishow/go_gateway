import { useTranslation } from 'react-i18next';
import type { ModbusShareDesiredMapping, ModbusShareStatus } from '../../../../../types/modbusShare';
import type { WorkbenchV2State } from '../../state/types';
import type { ModbusShareCandidateReviewState } from '../../../../../hooks/datalink/useModbusShareCandidateReview';
import type { SourceRuleLocalModbusOutputCandidateView } from '../../../../../types/sourceRuleCandidates';

export interface ShareAllocationRow {
  id: string;
  sourceLabel: string;
  mapping: ModbusShareDesiredMapping;
}

type ShareState = 'valid' | 'blocked' | 'invalidated' | 'failed';

export function allocationRowsFromCanonicalMappings(
  mappings: ModbusShareDesiredMapping[],
  state: Pick<WorkbenchV2State, 'rules' | 'points' | 'mappings'>,
): ShareAllocationRow[] {
  return mappings
    .filter((mapping) => Number.isInteger(mapping.share_start_register) && Number.isInteger(mapping.zero_based_register) &&
      Number.isInteger(mapping.span_registers) && mapping.span_registers > 0 && Number.isInteger(mapping.stride_registers) &&
      mapping.stride_registers > 0)
    .map((mapping) => {
      const browserMapping = Object.values(state.mappings).find((candidate) => candidate.tag_id === mapping.tag_id);
      const point = browserMapping
        ? state.points.find((candidate) => candidate.id === browserMapping.point_id)
        : undefined;
      const rule = state.rules.find((candidate) => candidate.id === mapping.source_rule_id);
      return {
        id: mapping.tag_id,
        sourceLabel: point?.address ?? mapping.tag_key ?? mapping.display_name ?? rule?.name ?? mapping.tag_id,
        mapping,
      };
    });
}

export function allocationState(status: ModbusShareStatus): ShareState {
  if (status.outcome === 'invalidated_unknown' || (status.invalidated_count ?? 0) > 0) {
    return 'invalidated';
  }
  if (status.candidate_set_status === 'blocked' || status.candidate_statuses?.some((item) => item.startsWith('blocked'))) {
    return 'blocked';
  }
  if (status.error || status.recovery || status.dirty_state || status.failed || status.bind_state === 'fail' || status.lifecycle_state === 'failed') {
    return 'failed';
  }
  return 'valid';
}

interface LocalModbusReviewSurfaceProps {
  rows: ShareAllocationRow[];
  status: ModbusShareStatus;
  showDiagnostic?: boolean;
  candidatesByTagID?: Map<string, SourceRuleLocalModbusOutputCandidateView>;
  candidateReviewState?: ModbusShareCandidateReviewState;
}

export function LocalModbusReviewSurface({
  rows,
  status,
  showDiagnostic = false,
  candidatesByTagID,
  candidateReviewState = 'not-requested',
}: LocalModbusReviewSurfaceProps) {
  const { t } = useTranslation('workbench-v2');
  const state = allocationState(status);
  const candidateStatuses = status.candidate_statuses?.join(', ');
  const hasDiagnostic = Boolean(status.error || status.recovery || status.dirty_state || candidateStatuses || state === 'invalidated');

  return (
    <>
      {showDiagnostic && hasDiagnostic && (
        <div
          data-testid="step4-share-output-diagnostic"
          className="mt-3 rounded-xl border border-red-500/25 bg-red-950/20 p-3 text-xs text-red-100"
          role="status"
        >
          <div className="font-semibold">{t('step4.share_output_diagnostic_title')}</div>
          <div className="mt-1 flex flex-wrap gap-x-3 gap-y-1 font-mono text-red-200/90">
            <span>{state}</span>
            {candidateStatuses && <span>{candidateStatuses}</span>}
            {(status.error?.code ?? status.recovery?.code ?? status.dirty_state) && <span>{status.error?.code ?? status.recovery?.code ?? status.dirty_state}</span>}
          </div>
          {(status.error?.action ?? status.recovery?.action) && (
            <p className="mt-2 text-red-100/90">
              <span className="font-semibold">{t('step4.share_output_diagnostic_action')}</span>{' '}
              {status.error?.action ?? status.recovery?.action}
            </p>
          )}
          {(status.error?.retryable ?? status.recovery?.retryable) && (
            <p className="mt-1 text-red-200/75">{t('step4.share_output_diagnostic_retryable')}</p>
          )}
        </div>
      )}
      {rows.length > 0 && (
        <div className="mt-4 border-t border-slate-700/60 pt-4" data-testid="step4-share-allocation-review">
          <div className="flex flex-wrap items-baseline justify-between gap-2">
            <h3 className="text-sm font-semibold text-slate-100">
              {t('step4.share_output_review_title')}
            </h3>
            <span className="text-xs text-slate-400">
              {t('step4.share_output_review_subtitle')}
            </span>
          </div>
          <div className="mt-3 space-y-2">
            {rows.map(({ id, sourceLabel, mapping }) => {
              const humanEnd = mapping.share_start_register + mapping.span_registers - 1;
              const zeroBasedEnd = mapping.zero_based_register + mapping.span_registers - 1;
              const capacity = mapping.capacity_registers ?? status.capacity_registers;
              const candidate = candidatesByTagID?.get(mapping.tag_id);
              const candidateStatus = candidate?.status ?? (
                candidateReviewState === 'loading'
                  ? 'candidate review loading'
                  : candidateReviewState === 'error'
                    ? 'candidate review unavailable'
                    : undefined
              );
              const displayState = candidateStatus ?? state;
              const isBlocked = displayState !== 'valid' && displayState !== 'ready';
              const statusLabel = ['valid', 'ready', 'blocked', 'invalidated', 'failed'].includes(displayState)
                ? t(`step4.share_output_status_${displayState}`, { defaultValue: displayState })
                : displayState;
              return (
                <article
                  key={id}
                  data-testid={`step4-share-output-row-${id}`}
                  className="rounded-xl border border-slate-700/70 bg-slate-950/35 p-3"
                >
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <div className="min-w-0">
                      <div className="truncate font-mono text-sm font-semibold text-slate-100">{sourceLabel}</div>
                      <div className="mt-0.5 font-mono text-[11px] text-slate-400">{mapping.tag_id}</div>
                    </div>
                    <span className={`rounded-full border px-2 py-0.5 text-[11px] font-semibold ${
                      !isBlocked
                        ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
                        : 'border-red-500/30 bg-red-500/10 text-red-200'
                    }`}>
                      {statusLabel}
                    </span>
                  </div>
                  <dl className="mt-3 grid grid-cols-2 gap-x-3 gap-y-2 text-xs sm:grid-cols-3 lg:grid-cols-6">
                    <div><dt className="text-slate-500">{t('step4.share_output_human_register')}</dt><dd className="font-mono text-slate-200">{mapping.share_start_register}–{humanEnd}</dd></div>
                    <div><dt className="text-slate-500">{t('step4.share_output_zero_based_register')}</dt><dd className="font-mono text-slate-200">{mapping.zero_based_register}–{zeroBasedEnd}</dd></div>
                    <div><dt className="text-slate-500">{t('step4.share_output_span')}</dt><dd className="font-mono text-slate-200">{mapping.span_registers}</dd></div>
                    <div><dt className="text-slate-500">{t('step4.share_output_stride')}</dt><dd className="font-mono text-slate-200">{mapping.stride_registers}</dd></div>
                    <div><dt className="text-slate-500">{t('step4.share_output_capacity')}</dt><dd className="font-mono text-slate-200">{capacity ?? '—'}</dd></div>
                    <div><dt className="text-slate-500">{t('step4.share_output_datatype')}</dt><dd className="font-mono text-slate-200">{mapping.data_type}</dd></div>
                  </dl>
                  {candidate?.blocking_reason && (
                    <p className="mt-2 rounded-lg border border-red-500/20 bg-red-500/5 px-2.5 py-2 text-xs text-red-100" role="status">
                      <span className="font-semibold">{t('step4.share_output_blocking_reason')}</span>{' '}
                      {candidate.blocking_reason}
                    </p>
                  )}
                  {mapping.ownership_proof && (
                    <p className="mt-2 text-xs text-slate-400">
                      <span className="font-semibold text-slate-300">{t('step4.share_output_ownership')}</span>{' '}
                      {mapping.ownership_proof.verified
                        ? t('step4.share_output_ownership_verified')
                        : t('step4.share_output_ownership_unverified')}
                      {' · '}{mapping.ownership_proof.basis}
                    </p>
                  )}
                </article>
              );
            })}
          </div>
        </div>
      )}
    </>
  );
}
