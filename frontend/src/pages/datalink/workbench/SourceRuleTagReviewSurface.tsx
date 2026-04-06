import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import {
  useSourceRuleTagReviewDecisionsQuery,
  useUpsertSourceRuleTagReviewDecisionMutation,
} from '../../../hooks/datalink/useSourceRuleTagReviewDecisions';
import { useSourceRulesQuery } from '../../../hooks/datalink/useSourceRules';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import type { SourceRuleRecord } from '../../../types/datalink';
import type {
  SourceRuleCandidateSetStatus,
  SourceRuleTagCandidateView,
} from '../../../types/sourceRuleCandidates';
import type {
  SourceRuleTagReviewDecisionAction,
} from '../../../types/sourceRuleTagReviewDecisions';
import { SourceRuleTagReviewCandidateRow } from './SourceRuleTagReviewCandidateRow';
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

function getErrorMessage(error: unknown, fallback: string) {
  return error instanceof Error ? error.message : fallback;
}

type ReviewFeedback = {
  tone: 'success' | 'error';
  message: string;
};

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
  const { data: tags = [] } = useTagsQuery();
  const candidateQuery = useSourceRuleCandidatesQuery(activeRule?.id ?? null);
  const reviewDecisionsQuery = useSourceRuleTagReviewDecisionsQuery(activeRule?.id ?? null);
  const upsertDecisionMutation = useUpsertSourceRuleTagReviewDecisionMutation(activeRule?.id ?? null);
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
  const activeDecisions = useMemo(
    () => (reviewDecisionsQuery.data ?? []).filter((decision) => !decision.stale),
    [reviewDecisionsQuery.data],
  );
  const decisionByCandidateId = useMemo(
    () =>
      new Map(
        activeDecisions.map((decision) => [decision.candidate_id, decision] as const),
      ),
    [activeDecisions],
  );
  const [renameDrafts, setRenameDrafts] = useState<Record<string, string>>({});
  const [overrideSelections, setOverrideSelections] = useState<Record<string, string>>({});
  const [reviewFeedback, setReviewFeedback] = useState<ReviewFeedback | null>(null);
  const [pendingDecision, setPendingDecision] = useState<{
    candidateId: string;
    action: SourceRuleTagReviewDecisionAction;
  } | null>(null);

  useEffect(() => {
    if (!activeRule?.id || crossStepContext.focusedRuleId === activeRule.id) {
      return;
    }

    setFocusedRuleId(activeRule.id);
  }, [activeRule?.id, crossStepContext.focusedRuleId, setFocusedRuleId]);

  useEffect(() => {
    setRenameDrafts({});
    setOverrideSelections({});
    setReviewFeedback(null);
    setPendingDecision(null);
  }, [activeRule?.id, openRevisionId]);

  if (!selectedDeviceId || persistedRules.length === 0 || !activeRule) {
    return null;
  }

  const handleRefreshReview = async () => {
    await Promise.all([rulesQuery.refetch(), candidateQuery.refetch()]);
  };

  const handleRenameDraftChange = (candidateId: string, value: string) => {
    setRenameDrafts((current) => ({
      ...current,
      [candidateId]: value,
    }));
  };

  const handleOverrideSelectionChange = (candidateId: string, value: string) => {
    setOverrideSelections((current) => ({
      ...current,
      [candidateId]: value,
    }));
  };

  const handleSaveDecision = async (
    candidate: SourceRuleTagCandidateView,
    action: SourceRuleTagReviewDecisionAction,
    request: {
      tag_key?: string;
      override_tag_id?: string;
    },
  ) => {
    setPendingDecision({
      candidateId: candidate.id,
      action,
    });

    try {
      const decision = await upsertDecisionMutation.mutateAsync({
        candidate_id: candidate.id,
        action,
        ...request,
      });
      if (action === 'rename') {
        setRenameDrafts((current) => ({
          ...current,
          [candidate.id]: decision.tag_key ?? request.tag_key ?? '',
        }));
        setReviewFeedback({
          tone: 'success',
          message: t('workbench.tag.reviewSurface.feedback.renameSaved', {
            tagKey: decision.tag_key ?? request.tag_key ?? '',
          }),
        });
      } else if (action === 'skip') {
        setReviewFeedback({
          tone: 'success',
          message: t('workbench.tag.reviewSurface.feedback.skipSaved', {
            tagKey: candidate.tag_key,
          }),
        });
      } else {
        setOverrideSelections((current) => ({
          ...current,
          [candidate.id]: decision.override_tag_id ?? request.override_tag_id ?? '',
        }));
        setRenameDrafts((current) => ({
          ...current,
          [candidate.id]: decision.tag_key ?? '',
        }));
        setReviewFeedback({
          tone: 'success',
          message: t('workbench.tag.reviewSurface.feedback.overrideSaved', {
            tagKey: decision.tag_key ?? '',
          }),
        });
      }
    } catch (error) {
      setReviewFeedback({
        tone: 'error',
        message: getErrorMessage(error, t('workbench.tag.reviewSurface.feedback.saveFailed')),
      });
    } finally {
      setPendingDecision(null);
    }
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

      {reviewDecisionsQuery.isError ? (
        <div
          className="rounded-2xl border border-rose-500/30 bg-rose-500/10 px-4 py-4 text-sm text-rose-100"
          data-testid="source-rule-tag-review-decision-error"
        >
          <p className="font-medium">{t('workbench.tag.reviewSurface.feedback.loadSavedFailed')}</p>
          <p className="mt-1 text-rose-100/80">
            {reviewDecisionsQuery.error instanceof Error ? reviewDecisionsQuery.error.message : ''}
          </p>
        </div>
      ) : null}

      {reviewFeedback ? (
        <div
          className={`rounded-2xl border px-4 py-3 text-sm ${
            reviewFeedback.tone === 'success'
              ? 'border-emerald-500/30 bg-emerald-500/[0.12] text-emerald-100'
              : 'border-rose-500/30 bg-rose-500/10 text-rose-100'
          }`}
          data-testid="source-rule-tag-review-feedback"
        >
          {reviewFeedback.message}
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
              {tagCandidates.map((candidate) => {
                const currentDecision = decisionByCandidateId.get(candidate.id) ?? null;
                const renameValue =
                  renameDrafts[candidate.id] ?? currentDecision?.tag_key?.trim() ?? candidate.tag_key;
                const overrideSelection =
                  overrideSelections[candidate.id] ?? currentDecision?.override_tag_id ?? '';
                const rowPendingAction =
                  pendingDecision?.candidateId === candidate.id ? pendingDecision.action : null;
                const rowBusy =
                  pendingDecision?.candidateId === candidate.id && upsertDecisionMutation.isPending;

                return (
                  <SourceRuleTagReviewCandidateRow
                    key={candidate.id}
                    candidate={candidate}
                    currentDecision={currentDecision}
                    renameValue={renameValue}
                    overrideSelection={overrideSelection}
                    tags={tags}
                    staleReview={staleReview}
                    reviewDecisionsLoading={reviewDecisionsQuery.isLoading}
                    reviewDecisionsError={reviewDecisionsQuery.isError}
                    rowBusy={rowBusy}
                    rowPendingAction={rowPendingAction}
                    onRenameDraftChange={handleRenameDraftChange}
                    onOverrideSelectionChange={handleOverrideSelectionChange}
                    onSaveDecision={(rowCandidate, action, request) => {
                      void handleSaveDecision(rowCandidate, action, request);
                    }}
                  />
                );
              })}
            </ul>
          )}
        </>
      )}
    </section>
  );
}
