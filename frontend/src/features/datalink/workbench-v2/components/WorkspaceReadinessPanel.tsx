import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import type {
  StudioV2WorkspaceReadinessIssue,
  StudioV2WorkspaceReadinessStep,
  StudioV2WorkspaceReadinessSummary,
} from '../../../../types/studioV2WorkspaceReadiness';

export type WorkspaceReadinessStepNumber = 1 | 2 | 3 | 4;

export interface GroupedWorkspaceReadinessIssue {
  issue: StudioV2WorkspaceReadinessIssue;
  count: number;
}

interface WorkspaceReadinessPanelProps {
  summary?: StudioV2WorkspaceReadinessSummary | null;
  dataTestId: string;
  compact?: boolean;
  maxIssues?: number;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

export function groupWorkspaceReadinessIssues(
  issues: StudioV2WorkspaceReadinessIssue[],
): GroupedWorkspaceReadinessIssue[] {
  const groups = new Map<string, GroupedWorkspaceReadinessIssue>();

  issues.forEach((issue) => {
    const key = [
      issue.severity,
      issue.step,
      issue.code,
      issue.message,
    ].join('|');
    const existing = groups.get(key);

    if (existing) {
      existing.count += 1;
      return;
    }

    groups.set(key, { issue, count: 1 });
  });

  return Array.from(groups.values());
}

export function summarizeWorkspaceReadiness(
  summary?: StudioV2WorkspaceReadinessSummary | null,
): {
  hasBlockers: boolean;
  blockingCount: number;
  warningCount: number;
  issues: GroupedWorkspaceReadinessIssue[];
} {
  const groupedIssues = groupWorkspaceReadinessIssues(summary?.issues ?? []);
  const fallbackBlockingCount = groupedIssues
    .filter((group) => group.issue.severity === 'blocking')
    .reduce((total, group) => total + group.count, 0);
  const fallbackWarningCount = groupedIssues
    .filter((group) => group.issue.severity === 'warning')
    .reduce((total, group) => total + group.count, 0);
  const blockingCount = summary?.blocking_count ?? fallbackBlockingCount;
  const warningCount = summary?.warning_count ?? fallbackWarningCount;

  return {
    hasBlockers: blockingCount > 0,
    blockingCount,
    warningCount,
    issues: groupedIssues,
  };
}

export function readinessIssueLabel(
  issue: StudioV2WorkspaceReadinessIssue,
  t: TFunction<'workbench-v2'>,
): string {
  return t(`workspace_readiness.issue_labels.${issue.code}`, issue.code);
}

export function readinessIssueMessage(
  issue: StudioV2WorkspaceReadinessIssue,
  t: TFunction<'workbench-v2'>,
): string {
  return t(`workspace_readiness.issue_messages.${issue.code}`, issue.message);
}

export function groupedIssueScopeLabel(
  issue: StudioV2WorkspaceReadinessIssue,
  count: number,
  t?: TFunction<'workbench-v2'>,
): string {
  if (count <= 1) {
    return issue.scope;
  }
  if (!t) {
    return `${issue.scope} +${count - 1}`;
  }
  return t('workspace_readiness.scope_grouped', {
    scope: issue.scope,
    count: count - 1,
  });
}

export function compactReadinessGuidance(
  groupedIssue: GroupedWorkspaceReadinessIssue | undefined,
  t: TFunction<'workbench-v2'>,
): string | null {
  if (!groupedIssue) {
    return null;
  }

  const label = readinessIssueLabel(groupedIssue.issue, t);
  const action = readinessIssueSolution(groupedIssue.issue, t);
  return t('workspace_readiness.compact_fix', {
    label,
    action,
  });
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

  const readinessView = summarizeWorkspaceReadiness(summary);
  const issues = compact ? [] : readinessView.issues.slice(0, maxIssues);
  const primaryIssue = readinessView.issues[0];
  const compactGuidance = compactReadinessGuidance(primaryIssue, t);
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
          <div className={`mt-1 text-sm font-semibold ${readinessView.hasBlockers ? 'text-amber-200' : 'text-emerald-200'}`}>
            {readinessView.hasBlockers ? t('workspace_readiness.blocked') : t('workspace_readiness.ready')}
          </div>
        </div>
        <div className="flex items-center gap-2 text-[11px] font-mono">
          <span className="rounded-full border border-rose-500/30 bg-rose-500/10 px-2 py-1 text-rose-200">
            {t('workspace_readiness.blocking_count', { count: readinessView.blockingCount })}
          </span>
          <span className="rounded-full border border-amber-500/30 bg-amber-500/10 px-2 py-1 text-amber-200">
            {t('workspace_readiness.warning_count', { count: readinessView.warningCount })}
          </span>
        </div>
      </div>

      {compact && compactGuidance ? (
        <p className="mt-2 text-[11px] leading-5 text-slate-400">
          {compactGuidance}
        </p>
      ) : null}

      {!compact && issues.length > 0 ? (
        <div className="mt-3 space-y-2">
          {issues.map(({ issue, count }) => (
            <div key={`${issue.code}-${issue.step}-${issue.message}`} className="rounded-lg border border-slate-800/70 bg-black/20 px-3 py-2">
              <div className="flex items-center justify-between gap-2 text-[11px]">
                <span className={`flex items-center gap-2 font-semibold ${issue.severity === 'blocking' ? 'text-rose-200' : 'text-amber-200'}`}>
                  <span>{readinessIssueLabel(issue, t)}</span>
                  {count > 1 ? (
                    <span className="rounded-full border border-current/30 px-1.5 py-0.5 text-[10px]">
                      x{count}
                    </span>
                  ) : null}
                </span>
                <span className="text-slate-500">
                  {issue.step} · {groupedIssueScopeLabel(issue, count, t)}
                </span>
              </div>
              <p className="mt-1 text-[11px] leading-5 text-slate-400">{readinessIssueMessage(issue, t)}</p>
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
