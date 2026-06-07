import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { StudioV2WorkspaceReadinessSummary } from '../../../../types/studioV2WorkspaceReadiness';

interface WorkspaceReadinessPanelProps {
  summary?: StudioV2WorkspaceReadinessSummary | null;
  dataTestId: string;
  compact?: boolean;
  maxIssues?: number;
}

export const WorkspaceReadinessPanel: React.FC<WorkspaceReadinessPanelProps> = ({
  summary,
  dataTestId,
  compact = false,
  maxIssues = 4,
}) => {
  const { t } = useTranslation('workbench-v2');

  if (!summary) {
    return null;
  }

  const hasBlockers = summary.blocking_count > 0;
  const issues = summary.issues.slice(0, maxIssues);
  const containerClassName = compact
    ? 'rounded-xl border border-slate-700/60 bg-slate-950/50 px-4 py-3'
    : 'rounded-xl border border-slate-700/60 bg-slate-900/40 p-3';

  return (
    <section className={containerClassName} data-testid={dataTestId}>
      <div className="flex items-center justify-between gap-3">
        <div>
          <div className="text-[10px] uppercase tracking-wider text-slate-500">
            {t('workspace_readiness.title')}
          </div>
          <div className={`mt-1 text-sm font-semibold ${hasBlockers ? 'text-amber-200' : 'text-emerald-200'}`}>
            {hasBlockers ? t('workspace_readiness.blocked') : t('workspace_readiness.ready')}
          </div>
        </div>
        <div className="flex items-center gap-2 text-[11px] font-mono">
          <span className="rounded-full border border-rose-500/30 bg-rose-500/10 px-2 py-1 text-rose-200">
            {t('workspace_readiness.blocking_count', { count: summary.blocking_count })}
          </span>
          <span className="rounded-full border border-amber-500/30 bg-amber-500/10 px-2 py-1 text-amber-200">
            {t('workspace_readiness.warning_count', { count: summary.warning_count })}
          </span>
        </div>
      </div>

      {issues.length > 0 ? (
        <div className="mt-3 space-y-2">
          {issues.map((issue) => (
            <div key={`${issue.code}-${issue.scope}`} className="rounded-lg border border-slate-800/70 bg-black/20 px-3 py-2">
              <div className="flex items-center justify-between gap-2 text-[11px]">
                <span className={`font-semibold ${issue.severity === 'blocking' ? 'text-rose-200' : 'text-amber-200'}`}>
                  {issue.code}
                </span>
                <span className="text-slate-500">
                  {issue.step} · {issue.scope}
                </span>
              </div>
              <p className="mt-1 text-[11px] leading-5 text-slate-400">{issue.message}</p>
            </div>
          ))}
        </div>
      ) : (
        <p className="mt-3 text-[11px] text-slate-500">{t('workspace_readiness.no_issues')}</p>
      )}
    </section>
  );
};
