import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type {
  SourceRuleCandidateSetStatus,
  SourceRuleCandidateSetView,
  SourceRuleDatabaseOutputCandidateView,
  SourceRuleOutputStatus,
} from '../../../types/sourceRuleCandidates';
import { parseTableKey } from './databaseTargetBoardUtils';

type DatabaseOutputReviewPanelProps = {
  reviewRuleId: string | null;
  reviewRevisionId: string | null;
  reviewSet: SourceRuleCandidateSetView<SourceRuleDatabaseOutputCandidateView> | null;
  reviewLoading: boolean;
  selectedConnectorId: string;
  selectedConnectorName: string | null;
  selectedTableKey: string;
  selectedTagId: string;
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

function getCandidateTone(status: SourceRuleOutputStatus, isSelected: boolean) {
  if (status === 'blocked') {
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

function isVisibleInConnectorScope(
  candidate: SourceRuleDatabaseOutputCandidateView,
  selectedConnectorId: string,
) {
  if (!selectedConnectorId) {
    return !candidate.connector_id;
  }
  return !candidate.connector_id || candidate.connector_id === selectedConnectorId;
}

function buildScopePath(
  candidate: SourceRuleDatabaseOutputCandidateView,
  selectedConnectorName: string | null,
  selectedTableKey: string,
  t: ReturnType<typeof useTranslation>['t'],
) {
  if (
    candidate.connector_id &&
    candidate.table_schema &&
    candidate.table_name &&
    candidate.column_name
  ) {
    return t('workbench.output.database.reviewSurface.persistedPath', {
      connector: candidate.connector_id,
      table: `${candidate.table_schema}.${candidate.table_name}`,
      column: candidate.column_name,
    });
  }

  if (selectedConnectorName) {
    const currentTable = parseTableKey(selectedTableKey);
    return t('workbench.output.database.reviewSurface.currentScopePath', {
      connector: selectedConnectorName,
      table: currentTable ? `${currentTable.schema}.${currentTable.name}` : '—',
    });
  }

  return t('workbench.output.database.reviewSurface.unscoped');
}

export function DatabaseOutputReviewPanel({
  reviewRuleId,
  reviewRevisionId,
  reviewSet,
  reviewLoading,
  selectedConnectorId,
  selectedConnectorName,
  selectedTableKey,
  selectedTagId,
}: DatabaseOutputReviewPanelProps) {
  const { t } = useTranslation();
  const scopedCandidates = useMemo(
    () =>
      (reviewSet?.candidates ?? []).filter((candidate) =>
        isVisibleInConnectorScope(candidate, selectedConnectorId),
      ),
    [reviewSet?.candidates, selectedConnectorId],
  );
  const currentTable = parseTableKey(selectedTableKey);

  return (
    <section
      className="space-y-4 rounded-2xl border border-cyan-500/20 bg-slate-950/60 p-4"
      data-testid="database-output-review-surface"
    >
      <div className="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
        <div className="space-y-1">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
            {t('workbench.output.database.reviewSurface.eyebrow')}
          </p>
          <h4 className="text-lg font-semibold text-slate-50">
            {t('workbench.output.database.reviewSurface.title')}
          </h4>
          <p className="max-w-3xl text-sm text-slate-300">
            {t('workbench.output.database.reviewSurface.description')}
          </p>
        </div>

        <div
          className={`rounded-full border px-3 py-1 text-xs font-semibold ${getSetTone(
            reviewSet?.status ?? 'deferred',
          )}`}
          data-testid="database-review-set-status"
        >
          {t(`workbench.output.database.reviewSurface.status.${reviewSet?.status ?? 'deferred'}`)}
        </div>
      </div>

      <div className="grid gap-3 md:grid-cols-3">
        <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.output.database.reviewSurface.revision')}
          </p>
          <p className="mt-2 font-mono text-sm text-slate-50" data-testid="database-review-revision">
            {reviewRevisionId ?? '—'}
          </p>
        </article>
        <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.output.database.reviewSurface.currentConnector')}
          </p>
          <p className="mt-2 text-sm font-semibold text-slate-50">
            {selectedConnectorName ?? t('workbench.output.database.reviewSurface.none')}
          </p>
        </article>
        <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.output.database.reviewSurface.currentTable')}
          </p>
          <p className="mt-2 text-sm font-semibold text-slate-50">
            {currentTable
              ? `${currentTable.schema}.${currentTable.name}`
              : t('workbench.output.database.reviewSurface.none')}
          </p>
        </article>
      </div>

      {!reviewRuleId ? (
        <p className="rounded-xl border border-dashed border-slate-700 bg-slate-900/60 px-3 py-4 text-sm text-slate-400">
          {t('workbench.output.database.reviewSurface.noRule')}
        </p>
      ) : reviewLoading ? (
        <p className="rounded-xl border border-dashed border-slate-700 bg-slate-900/60 px-3 py-4 text-sm text-slate-400">
          {t('workbench.output.database.reviewSurface.loading')}
        </p>
      ) : (
        <>
          {reviewSet?.reason ? (
            <p className="rounded-xl border border-amber-500/30 bg-amber-500/10 px-3 py-4 text-sm text-amber-100">
              {t('workbench.output.database.reviewSurface.setReason', {
                reason: reviewSet.reason,
              })}
            </p>
          ) : null}

          {scopedCandidates.length > 0 ? (
            <div className="grid gap-3 xl:grid-cols-[minmax(0,1.15fr)_minmax(260px,0.85fr)]">
              <div className="space-y-2">
                {scopedCandidates.map((candidate) => {
                  const isSelected = Boolean(candidate.tag_id && candidate.tag_id === selectedTagId);
                  return (
                    <article
                      key={candidate.id}
                      data-testid={`database-review-candidate-${candidate.id}`}
                      className={`space-y-2 rounded-2xl border p-4 ${getCandidateTone(
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
                                {t('workbench.output.database.reviewSurface.currentSelection')}
                              </span>
                            ) : null}
                          </div>
                          <p className="text-xs text-slate-400">
                            {t('workbench.output.database.reviewSurface.candidateMeta', {
                              point: candidate.display_name,
                              address: candidate.address,
                              dataType: candidate.data_type,
                            })}
                          </p>
                        </div>
                        <span className="rounded-full bg-slate-950/70 px-2 py-1 text-[11px] font-semibold text-slate-200">
                          {t(`workbench.output.database.reviewSurface.status.${candidate.status}`)}
                        </span>
                      </div>

                      <p className="text-xs text-slate-300">
                        {buildScopePath(candidate, selectedConnectorName, selectedTableKey, t)}
                      </p>

                      {candidate.blocking_reason ? (
                        <p className="text-xs text-amber-100">{candidate.blocking_reason}</p>
                      ) : null}

                      {!candidate.tag_id ? (
                        <p className="text-xs text-slate-400">
                          {t('workbench.output.database.reviewSurface.tagPending')}
                        </p>
                      ) : null}
                    </article>
                  );
                })}
              </div>

              <aside className="rounded-2xl border border-slate-800 bg-slate-900/70 p-4">
                <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                  {t('workbench.output.database.reviewSurface.currentConnector')}
                </p>
                <p className="mt-2 text-sm text-slate-300">
                  {selectedConnectorName ?? t('workbench.output.database.reviewSurface.none')}
                </p>
                <p className="mt-4 text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                  {t('workbench.output.database.reviewSurface.currentTable')}
                </p>
                <p className="mt-2 text-sm text-slate-300">
                  {currentTable
                    ? `${currentTable.schema}.${currentTable.name}`
                    : t('workbench.output.database.reviewSurface.none')}
                </p>
              </aside>
            </div>
          ) : (
            <p className="rounded-xl border border-dashed border-slate-700 bg-slate-900/60 px-3 py-4 text-sm text-slate-400">
              {t('workbench.output.database.reviewSurface.empty')}
            </p>
          )}
        </>
      )}
    </section>
  );
}
