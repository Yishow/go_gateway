import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import type {
  StudioV2WorkspaceReadinessIssue,
  StudioV2WorkspaceReadinessStep,
  StudioV2WorkspaceReadinessSummary,
} from '../../../../types/studioV2WorkspaceReadiness';

export type WorkspaceReadinessStepNumber = 1 | 2 | 3 | 4;

interface WorkspaceReadinessPanelProps {
  summary?: StudioV2WorkspaceReadinessSummary | null;
  dataTestId: string;
  compact?: boolean;
  maxIssues?: number;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

export const WorkspaceReadinessPanel: React.FC<WorkspaceReadinessPanelProps> = ({
  summary,
  dataTestId,
  compact = false,
  maxIssues = 4,
  onNavigateStep,
}) => {
  const { t } = useTranslation('workbench-v2');

  if (!summary) {
    return null;
  }

  const hasBlockers = summary.blocking_count > 0;
  const issues = compact ? [] : summary.issues.slice(0, maxIssues);
  const containerClassName = compact
    ? 'rounded-xl border border-slate-700/60 bg-slate-950/50 px-4 py-2'
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

      {!compact && issues.length > 0 ? (
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
              <div className="mt-2 flex items-center justify-between gap-2">
                <p className="text-[11px] leading-5 text-slate-300">
                  {readinessIssueSolution(issue, t)}
                </p>
                {onNavigateStep ? (
                  <button
                    type="button"
                    data-testid={`workspace-readiness-action-${issue.code}`}
                    onClick={() => onNavigateStep(readinessStepNumber(issue.step))}
                    className="shrink-0 rounded-md border border-cyan-500/40 bg-cyan-500/10 px-2 py-1 text-[11px] font-semibold text-cyan-200 transition-colors hover:bg-cyan-500/20"
                  >
                    {t('workspace_readiness.go_to_step', { step: issue.step })}
                  </button>
                ) : null}
              </div>
            </div>
          ))}
        </div>
      ) : !compact ? (
        <p className="mt-3 text-[11px] text-slate-500">{t('workspace_readiness.no_issues')}</p>
      ) : null}
    </section>
  );
};

export function readinessStepNumber(step: StudioV2WorkspaceReadinessStep): WorkspaceReadinessStepNumber {
  switch (step) {
    case 'Step 2':
      return 2;
    case 'Step 3':
      return 3;
    case 'Step 4':
      return 4;
    case 'Step 1':
    default:
      return 1;
  }
}

export function readinessIssueSolution(
  issue: StudioV2WorkspaceReadinessIssue,
  t: TFunction<'workbench-v2'>,
): string {
  switch (issue.code) {
    case 'point-missing':
      return t('workspace_readiness.solution_point_missing');
    default:
      break;
  }

  switch (issue.step) {
    case 'Step 2':
      return t('workspace_readiness.solution_step2');
    case 'Step 3':
      return t('workspace_readiness.solution_step3');
    case 'Step 4':
      return t('workspace_readiness.solution_step4');
    case 'Step 1':
    default:
      return t('workspace_readiness.solution_step1');
  }
}
