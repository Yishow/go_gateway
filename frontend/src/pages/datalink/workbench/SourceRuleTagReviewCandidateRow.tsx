import { useTranslation } from 'react-i18next';
import type { Tag } from '../../../types/datalink';
import type { SourceRuleTagCandidateView } from '../../../types/sourceRuleCandidates';
import type {
  SourceRuleTagReviewDecision,
  SourceRuleTagReviewDecisionAction,
} from '../../../types/sourceRuleTagReviewDecisions';
import type { DatabaseGroupingSuggestion } from './databaseGroupingSuggestions';
import { getDataTypeBitWidth } from './sourceCanvasModel';

type DecisionPayload = {
  tag_key?: string;
  override_tag_id?: string;
};

type SourceRuleTagReviewCandidateRowProps = {
  candidate: SourceRuleTagCandidateView;
  currentDecision: SourceRuleTagReviewDecision | null;
  effectiveTagKey: string;
  databaseSuggestion: DatabaseGroupingSuggestion;
  groupingDraft: {
    groupKey: string;
    columnName: string;
  };
  renameValue: string;
  overrideSelection: string;
  tags: Tag[];
  staleReview: boolean;
  reviewDecisionsLoading: boolean;
  reviewDecisionsError: boolean;
  rowBusy: boolean;
  rowPendingAction: SourceRuleTagReviewDecisionAction | null;
  onGroupingDraftChange: (
    candidateId: string,
    key: 'groupKey' | 'columnName',
    value: string,
  ) => void;
  onSaveGroupingOverride: (candidate: SourceRuleTagCandidateView) => void;
  onResetGroupingOverride: (candidate: SourceRuleTagCandidateView) => void;
  onRenameDraftChange: (candidateId: string, value: string) => void;
  onOverrideSelectionChange: (candidateId: string, value: string) => void;
  onSaveDecision: (
    candidate: SourceRuleTagCandidateView,
    action: SourceRuleTagReviewDecisionAction,
    request: DecisionPayload,
  ) => void;
};

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

function getReviewDecisionTone(action: SourceRuleTagReviewDecisionAction) {
  switch (action) {
    case 'rename':
      return 'border-cyan-500/30 bg-cyan-500/[0.12] text-cyan-100';
    case 'skip':
      return 'border-amber-500/30 bg-amber-500/[0.12] text-amber-100';
    case 'override':
      return 'border-violet-500/30 bg-violet-500/[0.12] text-violet-100';
  }
}

export function SourceRuleTagReviewCandidateRow({
  candidate,
  currentDecision,
  effectiveTagKey,
  databaseSuggestion,
  groupingDraft,
  renameValue,
  overrideSelection,
  tags,
  staleReview,
  reviewDecisionsLoading,
  reviewDecisionsError,
  rowBusy,
  rowPendingAction,
  onGroupingDraftChange,
  onSaveGroupingOverride,
  onResetGroupingOverride,
  onRenameDraftChange,
  onOverrideSelectionChange,
  onSaveDecision,
}: SourceRuleTagReviewCandidateRowProps) {
  const { t } = useTranslation();
  const appliedCandidate = getCandidateStatusKey(candidate) === 'applied';
  const actionDisabled = staleReview || reviewDecisionsLoading || reviewDecisionsError;
  const disableRename = actionDisabled || appliedCandidate || rowBusy;
  const disableSkip = actionDisabled || appliedCandidate || rowBusy;
  const disableOverride = actionDisabled || appliedCandidate || rowBusy || tags.length === 0;
  const disableGrouping = staleReview || rowBusy;
  const invalidGroupingDraft = groupingDraft.groupKey.trim() !== '' && groupingDraft.columnName.trim() === '';

  return (
    <li
      className="rounded-2xl border border-slate-800 bg-slate-900/80 px-4 py-4"
      data-testid={`source-rule-tag-review-row-${candidate.id}`}
    >
      <div className="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
        <div className="space-y-1">
          <p className="text-sm font-semibold text-slate-50">
            {candidate.display_name || effectiveTagKey}
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
          {currentDecision ? (
            <span
              className={`inline-flex items-center rounded-full border px-3 py-1 text-xs font-semibold ${getReviewDecisionTone(currentDecision.action)}`}
              data-testid={`source-rule-tag-review-decision-${candidate.id}`}
            >
              {t(`workbench.tag.reviewSurface.reviewDecision.${currentDecision.action}`)}
            </span>
          ) : null}
        </div>
      </div>

      <dl className="mt-4 grid gap-3 text-sm text-slate-300 md:grid-cols-3">
        <div>
          <dt className="text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500">
            {t('workbench.tag.reviewSurface.fields.tagKey')}
          </dt>
          <dd className="mt-1 font-mono text-slate-100">{effectiveTagKey}</dd>
          {effectiveTagKey !== candidate.tag_key ? (
            <p className="mt-1 text-xs text-slate-400">
              {t('workbench.tag.reviewSurface.fields.generatedTagKey')}: {candidate.tag_key}
            </p>
          ) : null}
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

      <div className="mt-4 rounded-2xl border border-slate-800/80 bg-slate-950/55 p-3">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div className="space-y-1">
            <p className="text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500">
              {t('workbench.tag.reviewSurface.databaseGrouping.title')}
            </p>
            <div className="flex flex-wrap items-center gap-2">
              <span
                className="inline-flex items-center rounded-full border border-cyan-500/25 bg-cyan-500/[0.12] px-3 py-1 text-xs font-semibold text-cyan-100"
                data-testid={`source-rule-tag-review-database-group-${candidate.id}`}
              >
                {databaseSuggestion.groupKey
                  ? databaseSuggestion.groupKey
                  : t('workbench.tag.reviewSurface.databaseGrouping.singleRow')}
              </span>
              <span
                className="inline-flex items-center rounded-full border border-slate-700 bg-slate-900/75 px-3 py-1 text-xs font-semibold text-slate-200"
                data-testid={`source-rule-tag-review-database-column-${candidate.id}`}
              >
                {databaseSuggestion.columnName || '—'}
              </span>
              {databaseSuggestion.writeIntervalSeconds ? (
                <span className="text-xs text-slate-400">
                  {t('workbench.tag.reviewSurface.databaseGrouping.interval', {
                    seconds: databaseSuggestion.writeIntervalSeconds,
                  })}
                </span>
              ) : null}
            </div>
          </div>
          {databaseSuggestion.overridden ? (
            <span className="rounded-full border border-violet-500/30 bg-violet-500/[0.12] px-3 py-1 text-xs font-semibold text-violet-100">
              {t('workbench.tag.reviewSurface.databaseGrouping.overrideBadge')}
            </span>
          ) : null}
        </div>
        <p
          className="mt-2 text-xs text-slate-400"
          data-testid={`source-rule-tag-review-database-members-${candidate.id}`}
        >
          {databaseSuggestion.groupKey
            ? `${t('workbench.tag.reviewSurface.databaseGrouping.members')}: ${databaseSuggestion.memberTagKeys.join(', ')}`
            : t('workbench.tag.reviewSurface.databaseGrouping.singleMember')}
        </p>

        <div className="mt-3 grid gap-3 xl:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_auto_auto]">
          <label className="flex min-w-0 flex-col gap-1 text-xs font-medium text-slate-300">
            <span>{t('workbench.tag.reviewSurface.databaseGrouping.groupKey')}</span>
            <input
              value={groupingDraft.groupKey}
              disabled={disableGrouping}
              data-testid={`source-rule-tag-review-database-group-input-${candidate.id}`}
              onChange={(event) =>
                onGroupingDraftChange(candidate.id, 'groupKey', event.target.value)
              }
              className="h-10 rounded-xl border-0 bg-slate-900/70 px-3 text-sm text-slate-100 ring-1 ring-slate-700/50 focus:outline-none focus:ring-2 focus:ring-cyan-400/35 disabled:cursor-not-allowed disabled:opacity-60"
              placeholder={t('workbench.tag.reviewSurface.databaseGrouping.groupPlaceholder')}
            />
          </label>
          <label className="flex min-w-0 flex-col gap-1 text-xs font-medium text-slate-300">
            <span>{t('workbench.tag.reviewSurface.databaseGrouping.columnName')}</span>
            <input
              value={groupingDraft.columnName}
              disabled={disableGrouping}
              data-testid={`source-rule-tag-review-database-column-input-${candidate.id}`}
              onChange={(event) =>
                onGroupingDraftChange(candidate.id, 'columnName', event.target.value)
              }
              className="h-10 rounded-xl border-0 bg-slate-900/70 px-3 text-sm text-slate-100 ring-1 ring-slate-700/50 focus:outline-none focus:ring-2 focus:ring-cyan-400/35 disabled:cursor-not-allowed disabled:opacity-60"
              placeholder={t('workbench.tag.reviewSurface.databaseGrouping.columnPlaceholder')}
            />
          </label>
          <button
            type="button"
            disabled={disableGrouping || invalidGroupingDraft}
            data-testid={`source-rule-tag-review-database-save-${candidate.id}`}
            onClick={() => onSaveGroupingOverride(candidate)}
            className="inline-flex h-10 items-center justify-center rounded-xl border border-violet-400/30 bg-violet-500/[0.1] px-4 text-sm font-semibold text-violet-100 transition hover:border-violet-300/50 hover:bg-violet-500/[0.14] disabled:cursor-not-allowed disabled:opacity-60"
          >
            {t('workbench.tag.reviewSurface.databaseGrouping.save')}
          </button>
          <button
            type="button"
            disabled={disableGrouping}
            data-testid={`source-rule-tag-review-database-reset-${candidate.id}`}
            onClick={() => onResetGroupingOverride(candidate)}
            className="inline-flex h-10 items-center justify-center rounded-xl border border-slate-700 bg-slate-900/75 px-4 text-sm font-semibold text-slate-100 transition hover:border-slate-500 hover:bg-slate-900 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {t('workbench.tag.reviewSurface.databaseGrouping.reset')}
          </button>
        </div>
      </div>

      <div className="mt-4 rounded-2xl border border-slate-800/80 bg-slate-950/55 p-3">
        <div className="grid gap-3 xl:grid-cols-[minmax(0,1fr)_auto_auto]">
          <label className="flex min-w-0 flex-col gap-1 text-xs font-medium text-slate-300">
            <span>{t('workbench.tag.reviewSurface.actions.renameLabel')}</span>
            <input
              value={renameValue}
              disabled={disableRename}
              data-testid={`source-rule-tag-review-rename-input-${candidate.id}`}
              onChange={(event) => onRenameDraftChange(candidate.id, event.target.value)}
              className="h-10 rounded-xl border-0 bg-slate-900/70 px-3 text-sm text-slate-100 ring-1 ring-slate-700/50 focus:outline-none focus:ring-2 focus:ring-cyan-400/35 disabled:cursor-not-allowed disabled:opacity-60"
              placeholder={t('workbench.tag.reviewSurface.actions.renamePlaceholder')}
            />
          </label>
          <button
            type="button"
            disabled={disableRename || renameValue.trim() === ''}
            data-testid={`source-rule-tag-review-rename-save-${candidate.id}`}
            onClick={() =>
              onSaveDecision(candidate, 'rename', {
                tag_key: renameValue.trim(),
              })
            }
            className="inline-flex h-10 items-center justify-center rounded-xl border border-cyan-400/30 bg-cyan-500/[0.1] px-4 text-sm font-semibold text-cyan-100 transition hover:border-cyan-300/50 hover:bg-cyan-500/[0.14] disabled:cursor-not-allowed disabled:opacity-60"
          >
            {rowBusy && rowPendingAction === 'rename'
              ? t('workbench.tag.reviewSurface.actions.saving')
              : t('workbench.tag.reviewSurface.actions.renameSave')}
          </button>
          <button
            type="button"
            disabled={disableSkip}
            data-testid={`source-rule-tag-review-skip-${candidate.id}`}
            onClick={() => onSaveDecision(candidate, 'skip', {})}
            className="inline-flex h-10 items-center justify-center rounded-xl border border-amber-400/30 bg-amber-500/[0.1] px-4 text-sm font-semibold text-amber-100 transition hover:border-amber-300/50 hover:bg-amber-500/[0.14] disabled:cursor-not-allowed disabled:opacity-60"
          >
            {rowBusy && rowPendingAction === 'skip'
              ? t('workbench.tag.reviewSurface.actions.saving')
              : t('workbench.tag.reviewSurface.actions.skip')}
          </button>
        </div>

        <div className="mt-3 grid gap-3 xl:grid-cols-[minmax(0,1fr)_auto]">
          <label className="flex min-w-0 flex-col gap-1 text-xs font-medium text-slate-300">
            <span>{t('workbench.tag.reviewSurface.actions.overrideLabel')}</span>
            <select
              value={overrideSelection}
              disabled={disableOverride}
              data-testid={`source-rule-tag-review-override-select-${candidate.id}`}
              onChange={(event) => onOverrideSelectionChange(candidate.id, event.target.value)}
              className="h-10 rounded-xl border-0 bg-slate-900/70 px-3 text-sm text-slate-100 ring-1 ring-slate-700/50 focus:outline-none focus:ring-2 focus:ring-cyan-400/35 disabled:cursor-not-allowed disabled:opacity-60"
            >
              <option value="">
                {t('workbench.tag.reviewSurface.actions.overridePlaceholder')}
              </option>
              {tags.map((tag) => (
                <option key={tag.id} value={tag.id}>
                  {tag.key}
                </option>
              ))}
            </select>
          </label>
          <button
            type="button"
            disabled={disableOverride || overrideSelection.trim() === ''}
            data-testid={`source-rule-tag-review-override-save-${candidate.id}`}
            onClick={() =>
              onSaveDecision(candidate, 'override', {
                override_tag_id: overrideSelection,
              })
            }
            className="inline-flex h-10 items-center justify-center rounded-xl border border-violet-400/30 bg-violet-500/[0.1] px-4 text-sm font-semibold text-violet-100 transition hover:border-violet-300/50 hover:bg-violet-500/[0.14] disabled:cursor-not-allowed disabled:opacity-60"
          >
            {rowBusy && rowPendingAction === 'override'
              ? t('workbench.tag.reviewSurface.actions.saving')
              : t('workbench.tag.reviewSurface.actions.overrideSave')}
          </button>
        </div>

        {appliedCandidate ? (
          <p className="mt-3 text-xs text-slate-400">
            {t('workbench.tag.reviewSurface.actions.appliedLocked')}
          </p>
        ) : tags.length === 0 ? (
          <p className="mt-3 text-xs text-slate-400">
            {t('workbench.tag.reviewSurface.actions.noOverrideTags')}
          </p>
        ) : null}
      </div>

      {candidate.blocking_reason ? (
        <p className="mt-3 text-xs text-amber-200">{candidate.blocking_reason}</p>
      ) : null}
    </li>
  );
}
