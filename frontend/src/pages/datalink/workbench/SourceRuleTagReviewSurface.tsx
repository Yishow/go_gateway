import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import { useSourceRulesQuery } from '../../../hooks/datalink/useSourceRules';
import type { SourceRuleRecord } from '../../../types/datalink';
import type {
  SourceRuleCandidateSetStatus,
  SourceRuleTagCandidateView,
} from '../../../types/sourceRuleCandidates';
import { getDataTypeBitWidth } from './sourceCanvasModel';
import { useWorkbench } from './WorkbenchProvider';

function resolveActiveRuleId(
  rules: SourceRuleRecord[],
  preferredIds: Array<string | null | undefined>,
) {
  for (const candidateId of preferredIds) {
    if (candidateId && rules.some((rule) => rule.id === candidateId)) {
      return candidateId;
    }
  }

  return rules[0]?.id ?? null;
}

function buildRuleLabel(rule: SourceRuleRecord) {
  return `${rule.start_address} · ${rule.naming_prefix} · ${rule.data_type}`;
}

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

function getCandidateStatusKey(candidate: SourceRuleTagCandidateView) {
  return candidate.tag_id || candidate.mapping_id ? 'applied' : 'generated';
}

function getCandidateTone(candidate: SourceRuleTagCandidateView) {
  return getCandidateStatusKey(candidate) === 'applied'
    ? 'border-emerald-500/30 bg-emerald-500/[0.12] text-emerald-100'
    : 'border-sky-500/30 bg-sky-500/[0.12] text-sky-100';
}

function getMappingIntentTone(candidate: SourceRuleTagCandidateView) {
  return candidate.mapping_id
    ? 'border-emerald-500/20 bg-emerald-500/[0.1] text-emerald-100'
    : 'border-amber-500/20 bg-amber-500/[0.12] text-amber-100';
}

export function SourceRuleTagReviewSurface() {
  const { t } = useTranslation();
  const { selectedDeviceId, crossStepContext, sourcePlanningState, setFocusedRuleId } =
    useWorkbench();
  const rulesQuery = useSourceRulesQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
    selectedDeviceId ? { refetchInterval: 5000, refetchOnWindowFocus: true } : undefined,
  );
  const { data: rules = [] } = rulesQuery;

  const persistedRules = useMemo(
    () => rules.filter((rule) => rule.device_id === selectedDeviceId),
    [rules, selectedDeviceId],
  );
  const activeRuleId = resolveActiveRuleId(persistedRules, [
    crossStepContext.focusedRuleId,
    sourcePlanningState.selectedRuleId,
  ]);
  const activeRule = persistedRules.find((rule) => rule.id === activeRuleId) ?? null;
  const candidateQuery = useSourceRuleCandidatesQuery(activeRule?.id ?? null);
  const candidateView = candidateQuery.data;
  const tagSet = candidateView?.tags;
  const tagCandidates = tagSet?.candidates ?? [];
  const pendingMappingCount = tagCandidates.filter((candidate) => !candidate.mapping_id).length;
  const latestRevisionId = activeRule?.revision_id ?? null;
  const openRevisionId = candidateView?.revision_id ?? null;
  const staleReview = Boolean(
    latestRevisionId && openRevisionId && latestRevisionId !== openRevisionId,
  );
  const refreshingReview = rulesQuery.isRefetching || candidateQuery.isRefetching;

  if (!selectedDeviceId || persistedRules.length === 0 || !activeRule) {
    return null;
  }

  const handleRefreshReview = async () => {
    await Promise.all([rulesQuery.refetch(), candidateQuery.refetch()]);
  };

  return (
    <section
      className="shrink-0 space-y-4 border-b border-slate-800/70 bg-slate-950/40 px-4 py-4"
      data-testid="source-rule-tag-review-surface"
    >
      <div className="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
        <div className="space-y-1">
          <p className="text-[10px] font-semibold uppercase tracking-[0.2em] text-cyan-300/90">
            {t('workbench.tag.reviewSurface.eyebrow')}
          </p>
          <h3 className="text-base font-semibold text-slate-50">
            {t('workbench.tag.reviewSurface.title')}
          </h3>
          <p className="max-w-3xl text-sm leading-6 text-slate-400">
            {t('workbench.tag.reviewSurface.description')}
          </p>
        </div>

        <label className="flex min-w-0 flex-col gap-1.5 text-xs font-medium text-slate-300">
          <span>{t('workbench.tag.reviewSurface.activeRule')}</span>
          <select
            value={activeRule.id}
            data-testid="source-rule-tag-review-rule-select"
            onChange={(event) => setFocusedRuleId(event.target.value)}
            className="h-[42px] min-w-0 rounded-xl border-0 bg-slate-900/70 px-3 text-sm text-slate-100 ring-1 ring-slate-700/50 focus:outline-none focus:ring-2 focus:ring-cyan-400/35 xl:min-w-[280px]"
          >
            {persistedRules.map((rule) => (
              <option key={rule.id} value={rule.id}>
                {buildRuleLabel(rule)}
              </option>
            ))}
          </select>
        </label>
      </div>

      {staleReview ? (
        <div
          className="rounded-2xl border border-amber-500/40 bg-amber-500/10 px-4 py-4 text-amber-50"
          data-testid="source-rule-tag-review-stale"
        >
          <div className="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
            <div className="space-y-2">
              <p className="text-sm font-semibold">
                {t('workbench.tag.reviewSurface.stale.title')}
              </p>
              <p className="text-sm leading-6 text-amber-50/85">
                {t('workbench.tag.reviewSurface.stale.description')}
              </p>
              <dl className="grid gap-3 text-xs text-amber-100/85 md:grid-cols-2">
                <div>
                  <dt className="font-semibold uppercase tracking-[0.14em] text-amber-200/90">
                    {t('workbench.tag.reviewSurface.stale.openRevision')}
                  </dt>
                  <dd
                    className="mt-1 font-mono text-sm text-amber-50"
                    data-testid="source-rule-tag-review-open-revision"
                  >
                    {openRevisionId}
                  </dd>
                </div>
                <div>
                  <dt className="font-semibold uppercase tracking-[0.14em] text-amber-200/90">
                    {t('workbench.tag.reviewSurface.stale.latestRevision')}
                  </dt>
                  <dd
                    className="mt-1 font-mono text-sm text-amber-50"
                    data-testid="source-rule-tag-review-latest-revision"
                  >
                    {latestRevisionId}
                  </dd>
                </div>
              </dl>
            </div>

            <button
              type="button"
              onClick={() => void handleRefreshReview()}
              disabled={refreshingReview}
              data-testid="source-rule-tag-review-refresh"
              className="inline-flex h-10 items-center justify-center rounded-xl border border-amber-300/40 bg-slate-950/40 px-4 text-sm font-semibold text-amber-50 transition hover:border-amber-200/60 hover:bg-slate-950/60 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {t(
                `workbench.tag.reviewSurface.stale.${refreshingReview ? 'refreshing' : 'refresh'}`,
              )}
            </button>
          </div>
        </div>
      ) : null}

      {candidateQuery.isLoading ? (
        <div
          className="rounded-2xl border border-slate-800 bg-slate-900/75 px-4 py-6 text-sm text-slate-300"
          data-testid="source-rule-tag-review-loading"
        >
          {t('workbench.tag.reviewSurface.loading')}
        </div>
      ) : candidateQuery.isError ? (
        <div
          className="rounded-2xl border border-rose-500/30 bg-rose-500/10 px-4 py-4 text-sm text-rose-100"
          data-testid="source-rule-tag-review-error"
        >
          <p className="font-medium">{t('workbench.tag.reviewSurface.loadFailed')}</p>
          <p className="mt-1 text-rose-100/80">
            {candidateQuery.error instanceof Error ? candidateQuery.error.message : ''}
          </p>
        </div>
      ) : (
        <>
          <div className="grid gap-3 md:grid-cols-3">
            <div className="rounded-2xl border border-slate-800 bg-slate-900/75 px-4 py-3">
              <dt className="text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500">
                {t('workbench.tag.reviewSurface.revision')}
              </dt>
              <dd
                className="mt-2 font-mono text-sm text-slate-100"
                data-testid="source-rule-tag-review-revision"
              >
                {candidateView?.revision_id ?? '—'}
              </dd>
            </div>
            <div className="rounded-2xl border border-slate-800 bg-slate-900/75 px-4 py-3">
              <dt className="text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500">
                {t('workbench.tag.reviewSurface.metrics.generated')}
              </dt>
              <dd className="mt-2 font-mono text-2xl font-semibold text-slate-50">
                {tagCandidates.length}
              </dd>
            </div>
            <div className="rounded-2xl border border-slate-800 bg-slate-900/75 px-4 py-3">
              <dt className="text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500">
                {t('workbench.tag.reviewSurface.metrics.pendingMappings')}
              </dt>
              <dd className="mt-2 font-mono text-2xl font-semibold text-amber-100">
                {pendingMappingCount}
              </dd>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <span
              className={`inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold ${getSetTone(tagSet?.status ?? 'ready')}`}
              data-testid="source-rule-tag-review-set-status"
            >
              {t(`workbench.tag.reviewSurface.setStatus.${tagSet?.status ?? 'ready'}`)}
            </span>
            {tagSet?.reason ? (
              <span className="text-xs text-slate-400">
                {t('workbench.tag.reviewSurface.reason')}: {tagSet.reason}
              </span>
            ) : null}
          </div>

          {tagCandidates.length === 0 ? (
            <div
              className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/60 px-4 py-6"
              data-testid="source-rule-tag-review-empty"
            >
              <p className="text-sm font-medium text-slate-100">
                {t('workbench.tag.reviewSurface.empty.title')}
              </p>
              <p className="mt-2 text-sm leading-6 text-slate-400">
                {t('workbench.tag.reviewSurface.empty.description')}
              </p>
            </div>
          ) : (
            <ul className="grid gap-3" data-testid="source-rule-tag-review-candidate-list">
              {tagCandidates.map((candidate) => (
                <li
                  key={candidate.id}
                  className="rounded-2xl border border-slate-800 bg-slate-900/80 px-4 py-4"
                  data-testid={`source-rule-tag-review-row-${candidate.id}`}
                >
                  <div className="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
                    <div className="space-y-1">
                      <p className="text-sm font-semibold text-slate-50">
                        {candidate.display_name || candidate.tag_key}
                      </p>
                      <p className="text-xs text-slate-400">
                        {t('workbench.tag.reviewSurface.candidateMeta', {
                          address: candidate.address,
                          dataType: candidate.data_type,
                          bitWidth: getDataTypeBitWidth(candidate.data_type),
                        })}
                      </p>
                    </div>
                    <div className="flex flex-wrap gap-2">
                      <span
                        className={`inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold ${getCandidateTone(candidate)}`}
                        data-testid={`source-rule-tag-review-status-${candidate.id}`}
                      >
                        {t(`workbench.tag.reviewSurface.status.${getCandidateStatusKey(candidate)}`)}
                      </span>
                      <span
                        className={`inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold ${getMappingIntentTone(candidate)}`}
                        data-testid={`source-rule-tag-review-mapping-intent-${candidate.id}`}
                      >
                        {t(
                          `workbench.tag.reviewSurface.mappingIntent.${candidate.mapping_id ? 'applied' : 'pending'}`,
                        )}
                      </span>
                    </div>
                  </div>

                  <dl className="mt-4 grid gap-3 text-sm text-slate-300 md:grid-cols-3">
                    <div>
                      <dt className="text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500">
                        {t('workbench.tag.reviewSurface.fields.tagKey')}
                      </dt>
                      <dd className="mt-1 font-mono text-slate-100">{candidate.tag_key}</dd>
                    </div>
                    <div>
                      <dt className="text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500">
                        {t('workbench.tag.reviewSurface.fields.pointId')}
                      </dt>
                      <dd className="mt-1 font-mono text-slate-100">{candidate.point_id}</dd>
                    </div>
                    <div>
                      <dt className="text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500">
                        {t('workbench.tag.reviewSurface.fields.mappingStatus')}
                      </dt>
                      <dd className="mt-1 text-slate-100">
                        {t(`workbench.tag.reviewSurface.mappingStatus.${candidate.status}`)}
                      </dd>
                    </div>
                  </dl>

                  {candidate.blocking_reason ? (
                    <p className="mt-3 text-xs text-amber-200">{candidate.blocking_reason}</p>
                  ) : null}
                </li>
              ))}
            </ul>
          )}
        </>
      )}
    </section>
  );
}
