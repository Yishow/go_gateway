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
  enabledTargetCount: number;
  hasConflict: boolean;
  readinessSummary?: StudioV2WorkspaceReadinessSummary | null;
  onActivate: () => void;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

/**
 * 提交前設定摘要與啟動按鈕元件
 * 落地需求：「Commit sequence and animation」之觸發與摘要區塊
 */
export function CommitSummary({
  deviceCount,
  ruleCount,
  pointCount,
  mappingCount,
  connector,
  enabledTargetCount,
  hasConflict,
  readinessSummary,
  onActivate,
  onNavigateStep,
}: CommitSummaryProps) {
  const { t } = useTranslation('workbench-v2');

  const { kind, schema, table } = connector;

  // 格式化資料庫寫入目標路徑
  const dbTargetText = kind === 'sqlite' || !schema
    ? `${kind} → ${table}`
    : `${kind} → ${schema}.${table}`;

  // 提交按鈕不可用條件：有衝突，或者啟用的寫入欄位為 0
  const readinessView = summarizeWorkspaceReadiness(readinessSummary);
  const hasReadinessBlocker = readinessView.hasBlockers;
  const blockingIssues = groupWorkspaceReadinessIssues(readinessSummary?.issues ?? [])
    .filter((group) => group.issue.severity === 'blocking')
    .slice(0, 3);
  const isSubmitDisabled = hasConflict || enabledTargetCount === 0 || hasReadinessBlocker;

  return (
    <div className="bg-gray-900/10 border border-gray-800 rounded-2xl p-6 flex flex-col justify-between h-full backdrop-blur-sm">
      <div className="space-y-6">
        {/* 標題 */}
        <div>
          <h3 className="text-lg font-semibold text-white">
            {t('step4.summary_title', '確認設定並第一次啟動')}
          </h3>
          <p className="text-xs text-gray-500 mt-1">
            {t('step4.summary_subtitle', 'Step 4 只會啟動目前合法、可用、尚未啟動的設備。')}
          </p>
        </div>

        <div className="rounded-xl border border-gray-800/60 bg-gray-950/40 p-3">
          <div className="text-xs text-slate-300">
            {t('step4.summary_compact', {
              defaultValue: '本次會套用 {{devices}} 台裝置、{{rules}} 組規則、{{points}} 個點位與 {{targets}} 個寫入欄位。',
              devices: deviceCount,
              rules: ruleCount,
              points: pointCount,
              targets: enabledTargetCount,
            })}
          </div>
          <div className="mt-2 text-[11px] font-mono text-blue-300">
            {dbTargetText}
          </div>
          <div className="mt-1 text-[11px] text-slate-500">
            {t('step4.summary_mapping_count', {
              defaultValue: '啟用中的映射列：{{count}}',
              count: mappingCount,
            })}
          </div>
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
            className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-3"
            data-testid="step4-blocker-resolution"
          >
            <div className="text-xs font-semibold text-amber-100">
              {t('step4.blocker_resolution_title', '要讓它可套用，先完成：')}
            </div>
            <div className="mt-1 text-[11px] leading-5 text-amber-100/80">
              {t('step4.blocker_resolution_subtitle', '這些 blocker 解除後，啟動按鈕會恢復可用。')}
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

        {/* 提示訊息 */}
        {enabledTargetCount === 0 && !hasConflict && (
          <p className="text-xs text-amber-500 text-center select-none">
            ⚠️ {t('step4.no_enabled_targets_warning', '尚未啟用任何資料表欄位寫入')}
          </p>
        )}

        {/* 寬版啟動按鈕 */}
        <button
          type="button"
          disabled={isSubmitDisabled}
          onClick={onActivate}
          className={`
            w-full py-3 px-4 rounded-xl font-medium text-sm text-center transition-all duration-200 flex items-center justify-center gap-2
            ${isSubmitDisabled
              ? 'bg-gray-800 text-gray-500 cursor-not-allowed border border-gray-800/20'
              : 'bg-blue-600 hover:bg-blue-500 text-white cursor-pointer shadow-lg shadow-blue-600/10 hover:shadow-blue-500/20 active:scale-[0.98]'
            }
          `}
        >
          <span>🚀</span>
          {t('step4.activate_btn', '第一次啟動設備')}
        </button>

        <p className="text-[10px] text-gray-500 text-center leading-normal select-none">
          {t('step4.activate_info', '系統會逐台啟動符合條件的設備，並保留每台成功或失敗結果。')}
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
    <div className="rounded-lg border border-slate-800/70 bg-slate-950/50 px-3 py-2">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="flex items-center gap-2 font-mono text-[11px] font-semibold text-amber-100">
            <span>{readinessIssueLabel(issue, t)}</span>
            {count > 1 ? (
              <span className="rounded-full border border-current/30 px-1.5 py-0.5 text-[10px]">
                x{count}
              </span>
            ) : null}
          </div>
          <div className="mt-1 text-[11px] leading-5 text-slate-300">
            {readinessIssueMessage(issue, t)}
          </div>
          <div className="mt-1 text-[11px] leading-5 text-slate-500">
            {issue.step} · {groupedIssueScopeLabel(issue, count, t)}
          </div>
          <div className="mt-1 text-[11px] leading-5 text-slate-400">
            {readinessIssueSolution(issue, t)}
          </div>
        </div>
        <button
          type="button"
          data-testid={`step4-blocker-action-${issue.code}`}
          disabled={!onNavigateStep}
          onClick={() => onNavigateStep?.(targetStep)}
          className="shrink-0 rounded-md border border-cyan-500/40 bg-cyan-500/10 px-2 py-1 text-[11px] font-semibold text-cyan-200 transition-colors hover:bg-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {t('step4.fix_step_action', { step: issue.step })}
        </button>
      </div>
    </div>
  );
}
