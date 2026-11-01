import { useTranslation } from 'react-i18next';
import type { DbConnector } from '../../state/types';
import type {
  StudioV2WorkspaceReadinessSummary,
} from '../../../../../types/studioV2WorkspaceReadiness';
import {
  groupWorkspaceReadinessIssues,
  groupedIssueScopeLabel,
  readinessIssueLabel,
  readinessIssueMessage,
  readinessIssueSolution,
  readinessStepNumber,
  summarizeWorkspaceReadiness,
  WorkspaceReadinessPanel,
  type WorkspaceReadinessStepNumber,
} from '../../components/WorkspaceReadinessPanel';

/**
 * CommitSummary 元件屬性
 */
interface CommitSummaryProps {
  deviceCount: number;
  ruleCount: number;
  pointCount: number;
  mappingCount: number;
  connector: DbConnector;
  /** Saved write groups of the workspace and how many are applied; never a gate by itself. */
  groupSummary: { state: 'loading' | 'error' | 'ready'; total: number; applied: number };
  readinessSummary?: StudioV2WorkspaceReadinessSummary | null;
  configurationSaved?: boolean;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

/**
 * 設定與 readiness 的狀態總覽。它不提供任何啟動入口：錄製只從已審核的設備範圍開始，
 * 全域啟動 API 只留給明確呼叫它的舊呼叫端。
 */
export function CommitSummary({
  deviceCount,
  ruleCount,
  pointCount,
  mappingCount,
  connector,
  groupSummary,
  readinessSummary,
  configurationSaved = false,
  onNavigateStep,
}: CommitSummaryProps) {
  const { t } = useTranslation('workbench-v2');

  const { kind, schema, table } = connector;

  // 格式化資料庫寫入目標路徑
  const dbTargetText = kind === 'sqlite' || !schema
    ? `${kind} → ${table}`
    : `${kind} → ${schema}.${table}`;

  // 啟用只依後端持久化的 readiness；是否需要資料庫輸出由後端判定（只用 Share 時不需要）。
  const readinessView = summarizeWorkspaceReadiness(readinessSummary);
  const hasReadinessBlocker = readinessView.hasBlockers;
  const blockingIssues = groupWorkspaceReadinessIssues(readinessSummary?.issues ?? [])
    .filter((group) => group.issue.severity === 'blocking')
    .slice(0, 3);

  return (
    <div className="bg-[var(--gw-surface)] border border-[var(--gw-border)] rounded-2xl p-6 flex flex-col justify-between h-full">
      <div className="space-y-6">
        {/* 標題 */}
        <div>
          <h3 className="text-lg font-semibold text-[var(--gw-text)]">
            {t('step4.summary_title')}
          </h3>
          <p className="text-xs text-[var(--gw-muted)] mt-1">
            {t('step4.summary_subtitle')}
          </p>
        </div>

        <div className="rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3">
          <div className="text-xs text-[var(--gw-text)]">
            {t('step4.summary_compact', {
              devices: deviceCount,
              rules: ruleCount,
              points: pointCount,
              groups: groupSummary.total,
            })}
          </div>
          <div className="mt-2 text-xs font-mono text-[var(--gw-accent)]">
            {dbTargetText}
          </div>
          <div className="mt-1 text-xs text-[var(--gw-muted)]">
            {t('step4.summary_mapping_count', {
              count: mappingCount,
            })}
          </div>
          {configurationSaved && (
            <div
              data-testid="step4-configuration-saved"
              className="mt-2 text-xs text-[var(--gw-success-fg)]"
            >
              {t('step4.configuration_saved')}
            </div>
          )}
        </div>
      </div>

      <div className="mt-8 space-y-3">
        <WorkspaceReadinessPanel
          summary={readinessSummary}
          dataTestId="step4-readiness-panel"
          maxIssues={4}
          onNavigateStep={onNavigateStep}
        />

        {hasReadinessBlocker && blockingIssues.length > 0 && (
          <div
            className="rounded-xl border border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] p-3"
            data-testid="step4-blocker-resolution"
          >
            <div className="text-xs font-semibold text-[var(--gw-warning-fg)]">
              {t('step4.blocker_resolution_title')}
            </div>
            <div className="mt-1 text-xs leading-5 text-[var(--gw-warning-fg)]">
              {t('step4.blocker_resolution_subtitle')}
            </div>
            <div className="mt-3 space-y-2">
              {blockingIssues.map((group) => (
                <BlockerResolutionItem
                  key={`${group.issue.code}-${group.issue.step}-${group.issue.message}`}
                  groupedIssue={group}
                  onNavigateStep={onNavigateStep}
                />
              ))}
            </div>
          </div>
        )}

        {/* 資料庫輸出狀態：只陳述事實，不當成啟用的門檻 */}
        <p className="text-xs text-[var(--gw-warning-fg)] text-center" role="status" data-testid="step4-group-summary-note">
          {groupSummary.state === 'error'
            ? t('step4.group.summary.unknown')
            : groupSummary.state === 'loading'
              ? t('step4.group.summary.loading')
              : groupSummary.total === 0
                ? t('step4.group.summary.none')
                : groupSummary.applied === 0
                  ? t('step4.group.summary.saved_not_applied', { count: groupSummary.total })
                  : t('step4.group.summary.applied', { applied: groupSummary.applied, total: groupSummary.total })}
        </p>
      </div>
    </div>
  );
}

interface BlockerResolutionItemProps {
  groupedIssue: ReturnType<typeof groupWorkspaceReadinessIssues>[number];
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

function BlockerResolutionItem({ groupedIssue, onNavigateStep }: BlockerResolutionItemProps) {
  const { t } = useTranslation('workbench-v2');
  const { issue, count } = groupedIssue;
  const targetStep = readinessStepNumber(issue.step);

  return (
    <div className="rounded-lg border border-[var(--gw-border)] bg-[var(--gw-surface)] px-3 py-2">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="flex items-center gap-2 font-mono text-xs font-semibold text-[var(--gw-warning-fg)]">
            <span>{readinessIssueLabel(issue, t)}</span>
            {count > 1 ? (
              <span className="rounded-full border border-current/30 px-1.5 py-0.5 text-xs">
                x{count}
              </span>
            ) : null}
          </div>
          <div className="mt-1 text-xs leading-5 text-[var(--gw-text)]">
            {readinessIssueMessage(issue, t)}
          </div>
          <div className="mt-1 text-xs leading-5 text-[var(--gw-muted)]">
            {issue.step} · {groupedIssueScopeLabel(issue, count, t)}
          </div>
          <div className="mt-1 text-xs leading-5 text-[var(--gw-muted)]">
            {readinessIssueSolution(issue, t)}
          </div>
        </div>
        <button
          type="button"
          data-testid={`step4-blocker-action-${issue.code}`}
          disabled={!onNavigateStep}
          onClick={() => onNavigateStep?.(targetStep)}
          className="shrink-0 rounded-md border border-[var(--gw-focus)] bg-[var(--gw-active)] px-2 py-1 text-xs font-semibold text-[var(--gw-accent)] transition-colors hover:bg-[var(--gw-active)] disabled:cursor-not-allowed disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]"
        >
          {t('step4.fix_step_action', { step: issue.step })}
        </button>
      </div>
    </div>
  );
}
