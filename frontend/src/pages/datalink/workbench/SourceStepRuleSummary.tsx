import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useSourceRuleCandidatesQuery } from '../../../hooks/datalink/useSourceRuleCandidates';
import type { ProtocolType } from '../../../types/datalink';
import { buildTagDatabaseGroupingSuggestions } from './databaseGroupingSuggestions';
import { buildDatabaseRowPlans } from './databaseRowPlannerModel';
import type { OutputTarget } from './workbenchTypes';
import type { AddressCanvasItem, SourceRule } from './sourceCanvasModel';
import {
  buildSourcePlanningDatabaseAdvisory,
  buildSourceRuleCoverageSummary,
} from './sourceStepRuleSummaryModel';

type SourceStepRuleSummaryProps = {
  activeOutputTarget: OutputTarget;
  activeRuleId: string | null;
  applyDisabled: boolean;
  batchCreateSummary: { successCount: number; failureCount: number } | null;
  items: ReadonlyArray<AddressCanvasItem>;
  previewRuleId: string | null;
  protocol: ProtocolType;
  rules: ReadonlyArray<SourceRule>;
  onApplyPlanning: () => void;
  onHoverRule: (ruleId: string | null) => void;
  onSelectRule: (ruleId: string) => void;
};

function buildRuleLabel(rule: SourceRule) {
  return `${rule.startAddress} · ${rule.count} · ${rule.dataType}`;
}

export function SourceStepRuleSummary({
  activeOutputTarget,
  activeRuleId,
  applyDisabled,
  batchCreateSummary,
  items,
  previewRuleId,
  protocol,
  rules,
  onApplyPlanning,
  onHoverRule,
  onSelectRule,
}: SourceStepRuleSummaryProps) {
  const { t } = useTranslation();
  const activeRule = rules.find((rule) => rule.id === activeRuleId) ?? rules[0] ?? null;
  const previewRule = rules.find((rule) => rule.id === previewRuleId) ?? activeRule;
  const activeCoverage = useMemo(
    () => buildSourceRuleCoverageSummary(items, activeRule?.id ?? null),
    [activeRule?.id, items],
  );
  const databaseAdvisory = useMemo(
    () => buildSourcePlanningDatabaseAdvisory(protocol, activeRule?.dataType ?? 'int16'),
    [activeRule?.dataType, protocol],
  );
  const previewQuery = useSourceRuleCandidatesQuery(
    previewRule?.persisted ? previewRule.id : null,
  );
  const previewTagCandidates = useMemo(
    () => previewQuery.data?.tags.candidates ?? [],
    [previewQuery.data],
  );
  const previewDatabaseCandidates = useMemo(
    () => previewQuery.data?.database_outputs.candidates ?? [],
    [previewQuery.data],
  );
  const databaseGroupingSuggestions = useMemo(
    () =>
      buildTagDatabaseGroupingSuggestions(
        previewTagCandidates,
        new Map(),
        previewDatabaseCandidates,
      ),
    [previewDatabaseCandidates, previewTagCandidates],
  );
  const previewRowPlans = useMemo(
    () =>
      buildDatabaseRowPlans({
        candidates: previewDatabaseCandidates,
        selectedConnectorId: '',
        selectedTableKey: '',
        selectedWriteMode: 'insert',
        selectedTimestampColumn: '',
        connectorDefaultWriteIntervalSeconds: null,
      }),
    [previewDatabaseCandidates],
  );

  if (!activeRule) {
    return null;
  }

  return (
    <section
      className="space-y-4 rounded-2xl border border-cyan-500/20 bg-cyan-500/[0.04] p-4"
      data-testid="source-active-rule-summary"
    >
      <div className="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between">
        <div className="space-y-3">
          <div className="space-y-1">
            <p className="text-[11px] font-semibold uppercase tracking-[0.24em] text-cyan-300">
              {t('workbench.source.activeRule.eyebrow')}
            </p>
            <h3 className="text-base font-semibold text-slate-100">
              {t('workbench.source.activeRule.title')}
            </h3>
            <p className="max-w-3xl text-sm text-slate-300">
              {t('workbench.source.activeRule.description')}
            </p>
          </div>
          <div className="flex flex-wrap gap-2" data-testid="source-active-rule-selector">
            {rules.map((rule) => {
              const selected = rule.id === activeRule.id;
              const previewing = rule.id === previewRule?.id && rule.id !== activeRule.id;
              return (
                <button
                  key={rule.id}
                  type="button"
                  data-testid={`source-preview-rule-${rule.id}`}
                  onClick={() => onSelectRule(rule.id)}
                  onMouseEnter={() => onHoverRule(rule.id)}
                  onMouseLeave={() => onHoverRule(null)}
                  className={[
                    'rounded-full border px-3 py-1.5 text-xs transition',
                    selected
                      ? 'border-cyan-400/60 bg-cyan-500/15 text-cyan-100'
                      : previewing
                        ? 'border-sky-400/40 bg-sky-500/10 text-sky-100'
                        : 'border-slate-700 bg-slate-900/60 text-slate-300 hover:border-slate-500',
                  ].join(' ')}
                >
                  {buildRuleLabel(rule)}
                </button>
              );
            })}
          </div>
        </div>
        <button
          type="button"
          onClick={onApplyPlanning}
          disabled={applyDisabled}
          className="rounded-xl bg-cyan-500 px-4 py-2.5 text-sm font-semibold text-slate-950 shadow-sm shadow-cyan-900/30 transition hover:bg-cyan-400 disabled:cursor-not-allowed disabled:bg-slate-800 disabled:text-slate-500 disabled:shadow-none"
        >
          {t('workbench.source.actions.createRulePoints')}
        </button>
      </div>

      <div className="grid gap-3 xl:grid-cols-[minmax(0,1.3fr)_minmax(18rem,1fr)]">
        <section className="space-y-3 rounded-xl border border-slate-800/70 bg-slate-950/50 p-3">
          <div className="grid gap-2 sm:grid-cols-3">
            <div className="rounded-lg border border-slate-800 bg-slate-900/70 px-3 py-2">
              <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
                {t('workbench.source.activeRule.startAddress')}
              </p>
              <p
                className="mt-1 font-mono text-sm text-slate-100"
                data-testid="source-active-rule-address"
              >
                {activeRule.startAddress}
              </p>
            </div>
            <div className="rounded-lg border border-slate-800 bg-slate-900/70 px-3 py-2">
              <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
                {t('workbench.source.activeRule.count')}
              </p>
              <p
                className="mt-1 font-mono text-sm text-slate-100"
                data-testid="source-active-rule-count"
              >
                {activeRule.count}
              </p>
            </div>
            <div className="rounded-lg border border-slate-800 bg-slate-900/70 px-3 py-2">
              <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
                {t('workbench.source.activeRule.dataType')}
              </p>
              <p className="mt-1 font-mono text-sm text-slate-100">{activeRule.dataType}</p>
            </div>
          </div>

          <div className="grid gap-2 sm:grid-cols-3">
            <div className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2">
              <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
                {t('workbench.source.activeRule.coveragePlanned')}
              </p>
              <p className="mt-1 font-mono text-sm text-sky-200">{activeCoverage.planned}</p>
            </div>
            <div className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2">
              <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
                {t('workbench.source.activeRule.coverageUsed')}
              </p>
              <p className="mt-1 font-mono text-sm text-emerald-200">{activeCoverage.used}</p>
            </div>
            <div className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2">
              <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
                {t('workbench.source.activeRule.coverageConflict')}
              </p>
              <p className="mt-1 font-mono text-sm text-rose-200">{activeCoverage.conflict}</p>
            </div>
          </div>

          {batchCreateSummary ? (
            <p className="text-sm text-slate-300">
              {t('workbench.source.planner.batchSummary', batchCreateSummary)}
            </p>
          ) : null}

          {activeOutputTarget === 'database' ? (
            <div
              className="rounded-xl border border-violet-500/20 bg-violet-500/[0.08] p-3"
              data-testid="source-database-planning-hints"
            >
              <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-violet-200">
                {t('workbench.source.databaseHints.title')}
              </p>
              <p className="mt-1 text-sm text-violet-100">
                {t('workbench.source.databaseHints.description')}
              </p>
              <div className="mt-2 grid gap-2 sm:grid-cols-3">
                <div className="rounded-lg border border-violet-400/15 bg-slate-950/35 px-3 py-2">
                  <p className="text-[10px] uppercase tracking-[0.18em] text-violet-200/75">
                    {t('workbench.source.databaseHints.targetDataType')}
                  </p>
                  <p className="mt-1 font-mono text-sm text-violet-50">
                    {databaseAdvisory.targetDataType}
                  </p>
                </div>
                <div className="rounded-lg border border-violet-400/15 bg-slate-950/35 px-3 py-2">
                  <p className="text-[10px] uppercase tracking-[0.18em] text-violet-200/75">
                    {t('workbench.source.databaseHints.scale')}
                  </p>
                  <p className="mt-1 font-mono text-sm text-violet-50">
                    {databaseAdvisory.scaleMultiplier} / {databaseAdvisory.scaleOffset}
                  </p>
                </div>
                <div className="rounded-lg border border-violet-400/15 bg-slate-950/35 px-3 py-2">
                  <p className="text-[10px] uppercase tracking-[0.18em] text-violet-200/75">
                    {t('workbench.source.databaseHints.naming')}
                  </p>
                  <p className="mt-1 font-mono text-sm text-violet-50">
                    {databaseAdvisory.namingPrefix}
                  </p>
                </div>
              </div>
            </div>
          ) : null}
        </section>

        <section
          className="space-y-3 rounded-xl border border-slate-800/70 bg-slate-950/50 p-3"
          data-testid="source-tag-preview"
        >
          <div>
            <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
              {t('workbench.source.preview.eyebrow')}
            </p>
            <h4 className="mt-1 text-sm font-semibold text-slate-100">
              {previewRule?.id === activeRule.id
                ? t('workbench.source.preview.title')
                : t('workbench.source.preview.hoverTitle')}
            </h4>
            <p className="mt-1 text-xs text-slate-400">
              {previewRule ? buildRuleLabel(previewRule) : t('workbench.source.preview.empty')}
            </p>
          </div>

          {!previewRule?.persisted ? (
            <p className="text-sm text-slate-300">{t('workbench.source.preview.pending')}</p>
          ) : previewQuery.isLoading ? (
            <p className="text-sm text-slate-300">{t('workbench.source.preview.loading')}</p>
          ) : (
            <>
              <div className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2">
                <p className="text-[10px] uppercase tracking-[0.18em] text-slate-500">
                  {t('workbench.source.preview.tagCandidates')}
                </p>
                <p className="mt-1 text-sm text-slate-100">
                  {t('workbench.source.preview.tagCandidateCount', {
                    count: previewTagCandidates.length,
                  })}
                </p>
              </div>

              {previewTagCandidates.length > 0 ? (
                <div className="space-y-2">
                  {previewTagCandidates.slice(0, 3).map((candidate) => {
                    const suggestion = databaseGroupingSuggestions[candidate.id];
                    return (
                      <article
                        key={candidate.id}
                        className="rounded-lg border border-slate-800 bg-slate-900/40 px-3 py-2"
                      >
                        <p className="text-sm font-medium text-slate-100">{candidate.tag_key}</p>
                        <p className="mt-1 text-xs text-slate-400">
                          {candidate.address} · {candidate.data_type}
                        </p>
                        {activeOutputTarget === 'database' && suggestion ? (
                          <p className="mt-1 text-xs text-violet-200">
                            {suggestion.groupKey ?? 'single'} → {suggestion.columnName || 'value'}
                          </p>
                        ) : null}
                      </article>
                    );
                  })}
                </div>
              ) : (
                <p className="text-sm text-slate-300">{t('workbench.source.preview.noCandidates')}</p>
              )}

              {activeOutputTarget === 'database' ? (
                <div className="space-y-2 rounded-lg border border-violet-500/15 bg-violet-500/[0.05] p-3">
                  <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-violet-200">
                    {t('workbench.source.databaseHints.grouping')}
                  </p>
                  {previewRowPlans.length > 0 ? (
                    previewRowPlans.slice(0, 2).map((plan) => (
                      <div
                        key={plan.id}
                        className="rounded-lg border border-violet-400/15 bg-slate-950/35 px-3 py-2"
                      >
                        <p className="text-sm text-violet-50">{plan.displayKey}</p>
                        <p className="mt-1 text-xs text-violet-200/80">
                          {plan.members.map((member) => member.columnName || member.tagKey).join(', ')}
                        </p>
                      </div>
                    ))
                  ) : (
                    <p className="text-sm text-violet-100">
                      {t('workbench.source.databaseHints.pendingGrouping')}
                    </p>
                  )}
                </div>
              ) : null}
            </>
          )}
        </section>
      </div>
    </section>
  );
}
