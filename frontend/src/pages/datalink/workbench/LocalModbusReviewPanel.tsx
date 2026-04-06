import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type { ModbusShareStatus } from '../../../types/datalink';
import type {
  SourceRuleCandidateSetStatus,
  SourceRuleCandidateSetView,
  SourceRuleLocalModbusOutputCandidateView,
  SourceRuleLocalModbusOutputStatus,
} from '../../../types/sourceRuleCandidates';
import { toModbusDisplayRegister } from './workbenchOutputTypes';

type LocalModbusReviewPanelProps = {
  reviewRuleId: string | null;
  reviewRevisionId: string | null;
  reviewSet: SourceRuleCandidateSetView<SourceRuleLocalModbusOutputCandidateView> | null;
  reviewLoading: boolean;
  selectedTagId: string;
  status: ModbusShareStatus | null;
  conflictCount: number;
};

type VerificationState = {
  tone: 'ready' | 'pending' | 'blocked';
  labelKey: string;
  detailKey: string;
};

function getSetTone(status: SourceRuleCandidateSetStatus) {
  switch (status) {
    case 'blocked':
      return 'border-rose-500/30 bg-rose-500/10 text-rose-100';
    case 'deferred':
      return 'border-amber-500/30 bg-amber-500/10 text-amber-100';
    case 'ready':
    default:
      return 'border-cyan-500/30 bg-cyan-500/10 text-cyan-100';
  }
}

function getCandidateTone(status: SourceRuleLocalModbusOutputStatus, isSelected: boolean) {
  if (status === 'blocked_conflict') {
    return isSelected
      ? 'border-rose-400/60 bg-rose-500/15 ring-1 ring-rose-400/40'
      : 'border-rose-500/25 bg-rose-500/10';
  }
  if (status === 'out_of_sync') {
    return isSelected
      ? 'border-amber-400/60 bg-amber-500/15 ring-1 ring-amber-400/40'
      : 'border-amber-500/25 bg-amber-500/10';
  }
  return isSelected
    ? 'border-cyan-400/60 bg-cyan-500/15 ring-1 ring-cyan-400/40'
    : 'border-slate-700 bg-slate-900/60';
}

function getVerificationTone(tone: VerificationState['tone']) {
  switch (tone) {
    case 'blocked':
      return 'border-rose-500/25 bg-rose-500/10 text-rose-100';
    case 'ready':
      return 'border-emerald-500/25 bg-emerald-500/10 text-emerald-100';
    case 'pending':
    default:
      return 'border-amber-500/25 bg-amber-500/10 text-amber-100';
  }
}

function getCandidateVerificationState(
  candidate: SourceRuleLocalModbusOutputCandidateView,
  serverEnabled: boolean,
): VerificationState {
  if (candidate.status === 'blocked_conflict') {
    return {
      tone: 'blocked',
      labelKey: 'workbench.output.modbusReviewSurface.verification.blocked',
      detailKey: 'workbench.output.modbusReviewSurface.verification.conflictDetail',
    };
  }
  if (candidate.status === 'out_of_sync') {
    return {
      tone: 'blocked',
      labelKey: 'workbench.output.modbusReviewSurface.verification.blocked',
      detailKey: 'workbench.output.modbusReviewSurface.verification.outOfSyncDetail',
    };
  }
  if (!candidate.tag_id) {
    return {
      tone: 'pending',
      labelKey: 'workbench.output.modbusReviewSurface.verification.pending',
      detailKey: 'workbench.output.modbusReviewSurface.verification.tagPendingDetail',
    };
  }
  if (candidate.register === undefined) {
    return {
      tone: 'pending',
      labelKey: 'workbench.output.modbusReviewSurface.verification.pending',
      detailKey: 'workbench.output.modbusReviewSurface.verification.registerPendingDetail',
    };
  }
  if (!serverEnabled) {
    return {
      tone: 'pending',
      labelKey: 'workbench.output.modbusReviewSurface.verification.pending',
      detailKey: 'workbench.output.modbusReviewSurface.verification.serverPendingDetail',
    };
  }
  return {
    tone: 'ready',
    labelKey: 'workbench.output.modbusReviewSurface.verification.ready',
    detailKey: 'workbench.output.modbusReviewSurface.verification.readyDetail',
  };
}

export function LocalModbusReviewPanel({
  reviewRuleId,
  reviewRevisionId,
  reviewSet,
  reviewLoading,
  selectedTagId,
  status,
  conflictCount,
}: LocalModbusReviewPanelProps) {
  const { t } = useTranslation();
  const serverEnabled = Boolean(status?.enabled);
  const blockedCount = useMemo(
    () =>
      (reviewSet?.candidates ?? []).filter((candidate) =>
        candidate.status === 'blocked_conflict' || candidate.status === 'out_of_sync',
      ).length,
    [reviewSet?.candidates],
  );
  const verificationSummary = blockedCount > 0
    ? t('workbench.output.modbusReviewSurface.summary.blocked', {
      count: blockedCount,
    })
    : !serverEnabled
      ? t('workbench.output.modbusReviewSurface.summary.pending')
      : t('workbench.output.modbusReviewSurface.summary.ready');

  return (
    <section
      className="space-y-4 rounded-2xl border border-cyan-500/20 bg-slate-950/60 p-4"
      data-testid="local-modbus-review-surface"
    >
      <div className="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
        <div className="space-y-1">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
            {t('workbench.output.modbusReviewSurface.eyebrow')}
          </p>
          <h4 className="text-lg font-semibold text-slate-50">
            {t('workbench.output.modbusReviewSurface.title')}
          </h4>
          <p className="max-w-3xl text-sm text-slate-300">
            {t('workbench.output.modbusReviewSurface.description')}
          </p>
        </div>

        <div
          className={`rounded-full border px-3 py-1 text-xs font-semibold ${getSetTone(
            reviewSet?.status ?? 'deferred',
          )}`}
          data-testid="local-modbus-review-set-status"
        >
          {t(`workbench.output.modbusReviewSurface.status.${reviewSet?.status ?? 'deferred'}`)}
        </div>
      </div>

      <div className="grid gap-3 md:grid-cols-3">
        <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.output.modbusReviewSurface.revision')}
          </p>
          <p className="mt-2 font-mono text-sm text-slate-50" data-testid="local-modbus-review-revision">
            {reviewRevisionId ?? '—'}
          </p>
        </article>
        <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.output.modbusReviewSurface.serverHealth')}
          </p>
          <p className="mt-2 text-sm font-semibold text-slate-50">
            {serverEnabled
              ? t('workbench.output.modbusReviewSurface.serverRunning', { address: status?.address ?? '—' })
              : t('workbench.output.modbusReviewSurface.serverStopped')}
          </p>
          <p className="mt-2 text-xs text-slate-400">
            {t('workbench.output.modbusReviewSurface.serverMeta', {
              mappings: status?.mapping_count ?? 0,
              conflicts: conflictCount,
            })}
          </p>
        </article>
        <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.output.modbusReviewSurface.verification.title')}
          </p>
          <p className="mt-2 text-sm font-semibold text-slate-50">
            {verificationSummary}
          </p>
          <p className="mt-2 text-xs text-slate-400">
            {t('workbench.output.modbusReviewSurface.verification.summaryDetail', {
              count: reviewSet?.candidates.length ?? 0,
            })}
          </p>
        </article>
      </div>

      {!reviewRuleId ? (
        <p className="rounded-xl border border-dashed border-slate-700 bg-slate-900/60 px-3 py-4 text-sm text-slate-400">
          {t('workbench.output.modbusReviewSurface.noRule')}
        </p>
      ) : reviewLoading ? (
        <p className="rounded-xl border border-dashed border-slate-700 bg-slate-900/60 px-3 py-4 text-sm text-slate-400">
          {t('workbench.output.modbusReviewSurface.loading')}
        </p>
      ) : (
        <>
          {reviewSet?.reason ? (
            <p className="rounded-xl border border-amber-500/30 bg-amber-500/10 px-3 py-4 text-sm text-amber-100">
              {t('workbench.output.modbusReviewSurface.setReason', {
                reason: reviewSet.reason,
              })}
            </p>
          ) : null}

          {reviewSet?.candidates.length ? (
            <div className="grid gap-3 xl:grid-cols-2">
              {reviewSet.candidates.map((candidate) => {
                const isSelected = Boolean(candidate.tag_id && candidate.tag_id === selectedTagId);
                const verification = getCandidateVerificationState(candidate, serverEnabled);
                const persistedRegister = candidate.register !== undefined
                  ? t('workbench.output.modbusReviewSurface.persistedRegister', {
                    register: toModbusDisplayRegister(candidate.register),
                    count: candidate.register_count,
                  })
                  : t('workbench.output.modbusReviewSurface.registerPending');

                return (
                  <article
                    key={candidate.id}
                    data-testid={`local-modbus-review-candidate-${candidate.id}`}
                    className={`space-y-3 rounded-2xl border p-4 ${getCandidateTone(
                      candidate.status,
                      isSelected,
                    )}`}
                  >
                    <div className="flex flex-wrap items-start justify-between gap-3">
                      <div className="space-y-1">
                        <div className="flex flex-wrap items-center gap-2">
                          <p className="text-sm font-semibold text-slate-50">{candidate.tag_key}</p>
                          {isSelected ? (
                            <span className="rounded-full bg-cyan-500/15 px-2 py-0.5 text-[11px] font-medium text-cyan-200">
                              {t('workbench.output.modbusReviewSurface.currentSelection')}
                            </span>
                          ) : null}
                        </div>
                        <p className="text-xs text-slate-400">
                          {t('workbench.output.modbusReviewSurface.candidateMeta', {
                            point: candidate.display_name,
                            address: candidate.address,
                            dataType: candidate.data_type,
                          })}
                        </p>
                      </div>
                      <div className="flex flex-wrap gap-2">
                        <span className="rounded-full bg-slate-950/70 px-2 py-1 text-[11px] font-semibold text-slate-200">
                          {t(`workbench.output.modbusReviewSurface.status.${candidate.status}`)}
                        </span>
                        <span className={`rounded-full border px-2 py-1 text-[11px] font-semibold ${getVerificationTone(verification.tone)}`}>
                          {t(verification.labelKey)}
                        </span>
                      </div>
                    </div>

                    <p className="text-xs text-slate-300">{persistedRegister}</p>

                    {candidate.blocking_reason ? (
                      <p className="text-xs text-amber-100">{candidate.blocking_reason}</p>
                    ) : null}

                    <p className="text-xs text-slate-400">
                      {t(verification.detailKey)}
                    </p>
                  </article>
                );
              })}
            </div>
          ) : (
            <p className="rounded-xl border border-dashed border-slate-700 bg-slate-900/60 px-3 py-4 text-sm text-slate-400">
              {t('workbench.output.modbusReviewSurface.empty')}
            </p>
          )}
        </>
      )}
    </section>
  );
}
