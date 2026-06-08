import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { DbConnector } from '../../state/types';
import { studioV2WorkspaceDatabaseAPI } from '../../../../../services/studioV2WorkspaceDatabase';
import type {
  StudioV2WorkspaceReadinessIssue,
  StudioV2WorkspaceReadinessSummary,
} from '../../../../../types/studioV2WorkspaceReadiness';
import {
  readinessIssueSolution,
  readinessStepNumber,
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
  schemaActionsDisabled?: boolean;
  schemaPreviewSignature?: string;
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
  schemaActionsDisabled = false,
  schemaPreviewSignature,
  readinessSummary,
  onActivate,
  onNavigateStep,
}: CommitSummaryProps) {
  const { t } = useTranslation('workbench-v2');

  const { kind, schema, table } = connector;
  const previewStateKey = schemaPreviewSignature ?? `${kind}|${schema}|${table}`;

  // 建立資料表（dry-run 預覽 + 實際執行）的本地狀態
  const [schemaPhase, setSchemaPhase] = useState<'idle' | 'previewing' | 'previewed' | 'creating' | 'done' | 'error'>('idle');
  const [schemaStatements, setSchemaStatements] = useState<string[]>([]);
  const [schemaError, setSchemaError] = useState<string | null>(null);
  const [schemaExecuted, setSchemaExecuted] = useState(0);
  const schemaBusy = schemaPhase === 'previewing' || schemaPhase === 'creating';
  const schemaControlsDisabled = schemaBusy || schemaActionsDisabled;

  useEffect(() => {
    setSchemaPhase('idle');
    setSchemaStatements([]);
    setSchemaError(null);
    setSchemaExecuted(0);
  }, [previewStateKey]);

  const handlePreviewSchema = async () => {
    setSchemaPhase('previewing');
    setSchemaError(null);
    try {
      const result = await studioV2WorkspaceDatabaseAPI.generateSchema(true);
      setSchemaStatements(result.statements);
      setSchemaPhase('previewed');
    } catch (error) {
      setSchemaError(error instanceof Error ? error.message : String(error));
      setSchemaPhase('error');
    }
  };

  const handleCreateSchema = async () => {
    setSchemaPhase('creating');
    setSchemaError(null);
    try {
      const result = await studioV2WorkspaceDatabaseAPI.generateSchema(false);
      setSchemaExecuted(result.executed);
      setSchemaPhase('done');
    } catch (error) {
      setSchemaError(error instanceof Error ? error.message : String(error));
      setSchemaPhase('error');
    }
  };

  // 格式化資料庫寫入目標路徑
  const dbTargetText = kind === 'sqlite' || !schema
    ? `${kind} → ${table}`
    : `${kind} → ${schema}.${table}`;

  // 提交按鈕不可用條件：有衝突，或者啟用的寫入欄位為 0
  const hasReadinessBlocker = (readinessSummary?.blocking_count ?? 0) > 0;
  const blockingIssues = readinessSummary?.issues
    .filter((issue) => issue.severity === 'blocking')
    .slice(0, 3) ?? [];
  const isSubmitDisabled = hasConflict || enabledTargetCount === 0 || hasReadinessBlocker;

  const summaryItems = [
    { label: t('step4.summary_devices', '採集裝置'), value: `${deviceCount} 台` },
    { label: t('step4.summary_rules', '點位規則'), value: `${ruleCount} 組` },
    { label: t('step4.summary_points', '採集點位'), value: `${pointCount} 個` },
    { label: t('step4.summary_mappings', '點位映射'), value: `${mappingCount} 筆` },
    { label: t('step4.summary_database', '資料庫寫入'), value: dbTargetText, isDb: true }
  ];

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

        {/* 摘要列 */}
        <div className="divide-y divide-gray-800/40 border border-gray-800/60 rounded-xl overflow-hidden bg-gray-950/40">
          {summaryItems.map((item, idx) => (
            <div key={idx} className="flex justify-between items-center p-3 text-xs">
              <span className="text-gray-400 select-none">{item.label}</span>
              <span className={`font-semibold ${item.isDb ? 'font-mono text-blue-400 max-w-[200px] truncate' : 'text-gray-200'}`}>
                {item.value}
              </span>
            </div>
          ))}
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
              {blockingIssues.map((issue) => (
                <BlockerResolutionItem
                  key={`${issue.code}-${issue.scope}`}
                  issue={issue}
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

        {/* 建立資料表（可選，dry-run 預覽後執行）。啟動時後端也會自動確保。 */}
        {kind !== 'sqlite' && (
          <div className="space-y-2 border border-gray-800/60 rounded-xl p-3 bg-gray-950/40">
            <div className="flex items-center justify-between gap-2">
              <span className="text-xs text-gray-400 select-none">
                {t('step4.schema_section_title', '目標資料表')}
              </span>
              <div className="flex gap-2">
                <button
                  type="button"
                  disabled={schemaControlsDisabled}
                  onClick={handlePreviewSchema}
                  className="text-[11px] px-2 py-1 rounded-md border border-gray-700 text-gray-300 hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {schemaPhase === 'previewing'
                    ? t('step4.schema_previewing', '預覽中...')
                    : t('step4.schema_preview_btn', '預覽 DDL')}
                </button>
                <button
                  type="button"
                  disabled={schemaControlsDisabled}
                  onClick={handleCreateSchema}
                  className="text-[11px] px-2 py-1 rounded-md border border-blue-700 text-blue-300 hover:bg-blue-900/40 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {schemaPhase === 'creating'
                    ? t('step4.schema_creating', '建立中...')
                    : t('step4.schema_create_btn', '建立資料表')}
                </button>
              </div>
            </div>

            {schemaPhase === 'previewed' && (
              schemaStatements.length === 0 ? (
                <p className="text-[11px] text-emerald-400 select-none">
                  {t('step4.schema_no_changes', '資料表已存在，無需建立')}
                </p>
              ) : (
                <pre className="text-[10px] text-gray-300 font-mono bg-black/40 rounded-md p-2 max-h-32 overflow-auto whitespace-pre-wrap">
                  {schemaStatements.join('\n')}
                </pre>
              )
            )}

            {schemaPhase === 'done' && (
              <p className="text-[11px] text-emerald-400 select-none">
                ✓ {t('step4.schema_done', '已建立 {{count}} 項資料表結構', { count: schemaExecuted })}
              </p>
            )}

            {schemaPhase === 'error' && schemaError && (
              <p className="text-[11px] text-rose-400 break-words">
                {t('step4.schema_error', '建表失敗')}: {schemaError}
              </p>
            )}
          </div>
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
  issue: StudioV2WorkspaceReadinessIssue;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
}

function BlockerResolutionItem({ issue, onNavigateStep }: BlockerResolutionItemProps) {
  const { t } = useTranslation('workbench-v2');
  const targetStep = readinessStepNumber(issue.step);

  return (
    <div className="rounded-lg border border-slate-800/70 bg-slate-950/50 px-3 py-2">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="font-mono text-[11px] font-semibold text-amber-100">
            {issue.code}
          </div>
          <div className="mt-1 text-[11px] leading-5 text-slate-300">
            {issue.message}
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
