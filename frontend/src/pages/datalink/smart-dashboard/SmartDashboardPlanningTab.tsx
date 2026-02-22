import { Filter, FolderOpen, Save } from 'lucide-react';
import type { TFunction } from 'i18next';
import type { DataType } from '../../../types/datalink';

type SourceTemplate = {
  id: string;
  name: string;
  dataType: DataType;
  count: number;
  startAddress: string;
  updatedAt: string;
  lastUsedAt: string;
  version: number;
};

export interface SmartDashboardPlanningTabProps {
  /** 規劃摘要（唯讀，Source Planner 那邊已編輯） */
  planDataType: DataType;
  planCount: number;
  totalPlannedCells: number;
  planConflictCount: number;
  typedPlanValidation: { valid: boolean };
  allocationMessage: string;
  /** 模板儲存 */
  templateName: string;
  setTemplateName: (value: string) => void;
  handleSaveTemplate: () => void;
  /** 模板列表 */
  sourceTemplates: SourceTemplate[];
  handleLoadTemplate: (template: SourceTemplate) => void;
  handleDeleteTemplate: (templateId: string) => void;
  staleTemplateCount: number;
  sourceTemplateSchemaVersion: number;
  handleUpgradeTemplates: () => void;
  /** 批次命名前綴 */
  batchNamePrefix: string;
  setBatchNamePrefix: (value: string) => void;
  normalizeNamingPrefix: (value: string) => string;
  namePreview: Array<{ sequence: number; name: string; conflict: boolean }>;
  nameConflictCount: number;
  /** 衝突過濾（控制 MemoryGrid） */
  showConflictsOnly: boolean;
  setShowConflictsOnly: (updater: (prev: boolean) => boolean) => void;
  t: TFunction;
}

/** 側欄「規劃」Tab：模板管理、批次命名、狀態摘要 */
export default function SmartDashboardPlanningTab({
  planDataType,
  planCount,
  totalPlannedCells,
  planConflictCount,
  typedPlanValidation,
  allocationMessage,
  templateName,
  setTemplateName,
  handleSaveTemplate,
  sourceTemplates,
  handleLoadTemplate,
  handleDeleteTemplate,
  staleTemplateCount,
  sourceTemplateSchemaVersion,
  handleUpgradeTemplates,
  batchNamePrefix,
  setBatchNamePrefix,
  normalizeNamingPrefix,
  namePreview,
  nameConflictCount,
  showConflictsOnly,
  setShowConflictsOnly,
  t,
}: SmartDashboardPlanningTabProps) {
  return (
    <div className="space-y-4 p-3">
      {/* 規劃摘要 */}
      <div className="flex flex-wrap items-center gap-1.5">
        <span className="rounded-full bg-slate-800/80 px-2 py-0.5 text-[11px] text-slate-400 ring-1 ring-white/10">
          {planDataType} × {planCount} = {totalPlannedCells} 格
        </span>
        {planConflictCount > 0 && (
          <span className="rounded-full bg-amber-500/15 px-2 py-0.5 text-[11px] text-amber-300 ring-1 ring-amber-500/30">
            衝突 {planConflictCount}
          </span>
        )}
        {!typedPlanValidation.valid && (
          <span className="rounded-full bg-rose-500/15 px-2 py-0.5 text-[11px] text-rose-300 ring-1 ring-rose-500/30">
            驗證失敗
          </span>
        )}
        {allocationMessage && (
          <span className="text-[11px] text-sky-300">{allocationMessage}</span>
        )}
      </div>

      {/* 儲存模板 */}
      <div>
        <p className="mb-1.5 text-[11px] font-semibold uppercase tracking-wider text-slate-500">
          {t('smartDashboard.sourcePlanner.saveTemplate')}
        </p>
        <div className="flex gap-2">
          <input
            value={templateName}
            onChange={(e) => setTemplateName(e.target.value)}
            placeholder={t('smartDashboard.sourcePlanner.templatePlaceholder')}
            className="min-w-0 flex-1 rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="button"
            onClick={handleSaveTemplate}
            disabled={!templateName.trim()}
            className="inline-flex shrink-0 items-center gap-1.5 rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 disabled:cursor-not-allowed disabled:opacity-40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
          >
            <Save className="h-3.5 w-3.5" />
            {t('smartDashboard.sourcePlanner.save')}
          </button>
        </div>
      </div>

      {/* 舊版模板升級 */}
      {staleTemplateCount > 0 && (
        <div className="rounded-lg border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-100">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <p>
              {t('smartDashboard.sourcePlanner.staleTemplates', {
                count: staleTemplateCount,
                version: sourceTemplateSchemaVersion,
              })}
            </p>
            <button
              type="button"
              onClick={handleUpgradeTemplates}
              className="rounded-md border border-amber-400/40 bg-amber-500/20 px-2 py-1 text-[11px] font-semibold hover:bg-amber-500/30"
            >
              {t('smartDashboard.sourcePlanner.upgradeTemplates')}
            </button>
          </div>
        </div>
      )}

      {/* 已儲存模板快捷列 */}
      {sourceTemplates.length > 0 && (
        <div>
          <p className="mb-1.5 text-[11px] font-semibold uppercase tracking-wider text-slate-500">
            模板列表
          </p>
          <div className="flex flex-wrap gap-1.5">
            {sourceTemplates.slice(0, 8).map((template) => (
              <div
                key={template.id}
                className="inline-flex items-center gap-1 rounded-full border border-white/10 bg-slate-800/60 py-1 pl-2 pr-1 text-[11px] text-slate-200"
              >
                <button
                  type="button"
                  onClick={() => handleLoadTemplate(template)}
                  className="inline-flex cursor-pointer items-center gap-1 hover:text-white"
                >
                  <FolderOpen className="h-3.5 w-3.5" />
                  {template.name}
                </button>
                <span className="rounded bg-slate-700/60 px-1.5 py-0.5 text-[10px] text-slate-300">
                  v{template.version}
                </span>
                <button
                  type="button"
                  onClick={() => handleDeleteTemplate(template.id)}
                  className="rounded-full px-1 text-slate-400 hover:bg-slate-700 hover:text-white"
                  aria-label={`Delete template ${template.name}`}
                >
                  ×
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* 批次命名前綴 */}
      <div className="rounded-lg border border-white/10 bg-slate-800/40 p-3">
        <div className="mb-1.5 flex items-center justify-between gap-2">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-slate-500">
            {t('smartDashboard.sourcePlanner.batchNamePrefix')}
          </p>
          <span className="shrink-0 text-[11px] text-slate-400">
            {t('smartDashboard.sourcePlanner.previewCount', { count: namePreview.length })}
            {nameConflictCount > 0 && (
              <span className="ml-1 text-rose-400">
                {t('smartDashboard.sourcePlanner.conflictCount', { count: nameConflictCount })}
              </span>
            )}
          </span>
        </div>
        <div className="flex items-center gap-2">
          <input
            value={batchNamePrefix}
            onChange={(e) => setBatchNamePrefix(normalizeNamingPrefix(e.target.value))}
            className="h-9 min-w-0 flex-1 rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder={t('smartDashboard.sourcePlanner.batchNamePrefixPlaceholder')}
            aria-label={t('smartDashboard.sourcePlanner.batchNamePrefix')}
          />
          <button
            type="button"
            onClick={() => setShowConflictsOnly((prev) => !prev)}
            className={`inline-flex h-9 shrink-0 items-center gap-1.5 rounded-lg border px-2.5 py-1.5 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500 ${
              showConflictsOnly
                ? 'border-amber-500/40 bg-amber-500/20 text-amber-100'
                : 'border-slate-700 bg-slate-800/70 text-slate-200'
            }`}
          >
            <Filter className="h-3.5 w-3.5" />
            {t('smartDashboard.sourcePlanner.showConflictsOnly')}
          </button>
        </div>
        {namePreview.length > 0 && (
          <div className="mt-2 flex flex-wrap gap-1.5">
            {namePreview.slice(0, 12).map((item) => (
              <span
                key={item.sequence}
                className={`rounded-md border px-2 py-0.5 text-[11px] font-mono ${
                  item.conflict
                    ? 'border-rose-500/40 bg-rose-500/10 text-rose-200'
                    : 'border-slate-700 bg-slate-900/70 text-slate-200'
                }`}
              >
                {item.name}
              </span>
            ))}
          </div>
        )}
        <p className="mt-1.5 text-[11px] text-slate-500">
          {t('smartDashboard.sourcePlanner.batchNamingHint')}
        </p>
      </div>
    </div>
  );
}
