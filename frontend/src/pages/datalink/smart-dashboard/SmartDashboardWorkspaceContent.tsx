import { Filter, FolderOpen, Save, WandSparkles } from 'lucide-react';
import type { RefObject } from 'react';
import type { TFunction } from 'i18next';
import { MemoryGrid } from '../../../components/datalink/MemoryGrid';
import type { Point, DataType, ProtocolType, ModbusShareStatus } from '../../../types/datalink';
import type { PlannedAllocation } from '../../../components/datalink/MemoryGrid';
import type { FlowState, FlowSegment } from '../../../features/flow/stateMachine';
import SmartDashboardFlowStatusSection from './SmartDashboardFlowStatusSection';

type IntentStage = 'idle' | 'grid' | 'commit';
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

interface SmartDashboardWorkspaceContentProps {
  planDataType: DataType;
  setPlanDataType: (value: DataType) => void;
  planCount: number;
  setPlanCount: (value: number) => void;
  totalPlannedCells: number;
  planStartAddress: string;
  setPlanStartAddress: (value: string) => void;
  handleAutoAllocate: () => void;
  handleApplyPlan: () => void;
  templateName: string;
  setTemplateName: (value: string) => void;
  handleSaveTemplate: () => void;
  showConflictsOnly: boolean;
  setShowConflictsOnly: (updater: (prev: boolean) => boolean) => void;
  batchNamePrefix: string;
  setBatchNamePrefix: (value: string) => void;
  normalizeNamingPrefix: (value: string) => string;
  namePreview: Array<{ sequence: number; name: string; conflict: boolean }>;
  nameConflictCount: number;
  staleTemplateCount: number;
  sourceTemplateSchemaVersion: number;
  handleUpgradeTemplates: () => void;
  planConflictCount: number;
  typedPlanValidation: { valid: boolean };
  allocationMessage: string;
  sourceTemplates: SourceTemplate[];
  handleLoadTemplate: (template: SourceTemplate) => void;
  handleDeleteTemplate: (templateId: string) => void;
  flowSegments: FlowSegment[];
  flowState: FlowState;
  statusStyle: Record<string, string>;
  hasError: boolean;
  t: TFunction;
  gridSectionRef: RefObject<HTMLElement | null>;
  resolveIntentMotionClass: (stage: IntentStage, reducedMotion: boolean) => string;
  guideStage: IntentStage;
  reducedMotion: boolean;
  motionTokens: { stageHandoffMs: number };
  goToLocalModbusWorkbench: () => void;
  modbusStatus: ModbusShareStatus | null;
  selectedDevice: { id: string; protocol: ProtocolType };
  allPoints: Point[];
  linkedAddresses: string[];
  selectedAddresses: string[];
  plannedAllocations: PlannedAllocation[];
  setSelectedAddresses: (addresses: string[]) => void;
  handleCellClick: (address: string) => void;
  getGridCenterAddress: (protocol: ProtocolType) => string;
}

export default function SmartDashboardWorkspaceContent({
  planDataType,
  setPlanDataType,
  planCount,
  setPlanCount,
  totalPlannedCells,
  planStartAddress,
  setPlanStartAddress,
  handleAutoAllocate,
  handleApplyPlan,
  templateName,
  setTemplateName,
  handleSaveTemplate,
  showConflictsOnly,
  setShowConflictsOnly,
  batchNamePrefix,
  setBatchNamePrefix,
  normalizeNamingPrefix,
  namePreview,
  nameConflictCount,
  staleTemplateCount,
  sourceTemplateSchemaVersion,
  handleUpgradeTemplates,
  planConflictCount,
  typedPlanValidation,
  allocationMessage,
  sourceTemplates,
  handleLoadTemplate,
  handleDeleteTemplate,
  flowSegments,
  flowState,
  statusStyle,
  hasError,
  t,
  gridSectionRef,
  resolveIntentMotionClass,
  guideStage,
  reducedMotion,
  motionTokens,
  goToLocalModbusWorkbench,
  modbusStatus,
  selectedDevice,
  allPoints,
  linkedAddresses,
  selectedAddresses,
  plannedAllocations,
  setSelectedAddresses,
  handleCellClick,
  getGridCenterAddress,
}: SmartDashboardWorkspaceContentProps) {
  return (
    <div className="scrollbar-thin scrollbar-thumb-slate-700/50 scrollbar-track-transparent flex-1 overflow-auto p-4 sm:p-8">
      <section className="mb-6 rounded-2xl border border-white/10 bg-slate-900/60 p-4">
        {/* 標題列：名稱、cell 統計、主要動作按鈕 */}
        <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
          <div className="flex items-center gap-2">
            <h3 className="text-sm font-semibold text-slate-100">{t('smartDashboard.sourcePlanner.title')}</h3>
            <span className="rounded-full bg-slate-800/80 px-2 py-0.5 text-[11px] text-slate-400 ring-1 ring-white/10">
              {t('smartDashboard.sourcePlanner.cells', {
                dataType: planDataType,
                count: planCount,
                total: totalPlannedCells,
              })}
            </span>
            {planConflictCount > 0 && (
              <span className="rounded-full bg-amber-500/15 px-2 py-0.5 text-[11px] text-amber-300 ring-1 ring-amber-500/30">
                {t('smartDashboard.sourcePlanner.conflict')} {planConflictCount}
              </span>
            )}
          </div>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handleAutoAllocate}
              className="inline-flex items-center gap-1.5 rounded-lg border border-sky-500/40 bg-sky-500/15 px-3 py-1.5 text-xs font-medium text-sky-100 hover:bg-sky-500/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500"
            >
              <WandSparkles className="h-3.5 w-3.5" />
              {t('smartDashboard.sourcePlanner.auto')}
            </button>
            <button
              type="button"
              onClick={handleApplyPlan}
              className="inline-flex items-center gap-1.5 rounded-lg border border-emerald-500/40 bg-emerald-500/20 px-3 py-1.5 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
            >
              {t('smartDashboard.sourcePlanner.applyToGrid')}
            </button>
          </div>
        </div>

        {/* 控制欄：四個輸入欄合一行 */}
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-[1fr_1fr_1fr_1.6fr]">
          <label className="text-xs text-slate-300">
            {t('smartDashboard.sourcePlanner.startAddress')}
            <input
              value={planStartAddress}
              onChange={(e) => setPlanStartAddress(e.target.value.toUpperCase())}
              placeholder={t('smartDashboard.sourcePlanner.startAddressPlaceholder')}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </label>
          <label className="text-xs text-slate-300">
            {t('smartDashboard.sourcePlanner.sourceCount')}
            <input
              type="number"
              min={1}
              max={200}
              value={planCount}
              onChange={(e) => {
                const raw = Number(e.target.value);
                const bounded = Number.isFinite(raw) ? Math.min(200, Math.max(1, Math.floor(raw))) : 1;
                setPlanCount(bounded);
              }}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </label>
          <label className="text-xs text-slate-300">
            {t('smartDashboard.sourcePlanner.dataType')}
            <select
              value={planDataType}
              onChange={(e) => setPlanDataType(e.target.value as DataType)}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="int16">{t('smartDashboard.sourcePlanner.dataTypeOption_int16')}</option>
              <option value="int32">{t('smartDashboard.sourcePlanner.dataTypeOption_int32')}</option>
              <option value="float32">{t('smartDashboard.sourcePlanner.dataTypeOption_float32')}</option>
              <option value="int64">{t('smartDashboard.sourcePlanner.dataTypeOption_int64')}</option>
              <option value="float64">{t('smartDashboard.sourcePlanner.dataTypeOption_float64')}</option>
            </select>
          </label>
          <label className="text-xs text-slate-300">
            {t('smartDashboard.sourcePlanner.saveTemplate')}
            <div className="mt-1 flex gap-2">
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
          </label>
        </div>

        {/* 批次命名前綴（標題列含預覽數、輸入與按鈕同一列等高對齊） */}
        <div className="mt-3 rounded-lg border border-white/10 bg-slate-800/40 p-3">
          <div className="mb-1.5 flex flex-wrap items-center justify-between gap-2">
            <label className="text-xs text-slate-300">{t('smartDashboard.sourcePlanner.batchNamePrefix')}</label>
            <span className="shrink-0 text-xs text-slate-400">
              {t('smartDashboard.sourcePlanner.previewCount', { count: namePreview.length })}
              {nameConflictCount > 0 && (
                <span className="ml-1 text-rose-400">{t('smartDashboard.sourcePlanner.conflictCount', { count: nameConflictCount })}</span>
              )}
            </span>
          </div>
          <div className="flex flex-wrap items-center gap-3">
            <input
              value={batchNamePrefix}
              onChange={(e) => setBatchNamePrefix(normalizeNamingPrefix(e.target.value))}
              className="h-10 min-w-0 flex-1 rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder={t('smartDashboard.sourcePlanner.batchNamePrefixPlaceholder')}
              aria-label={t('smartDashboard.sourcePlanner.batchNamePrefix')}
            />
            <button
              type="button"
              onClick={() => setShowConflictsOnly((prev) => !prev)}
              className={`h-10 shrink-0 inline-flex items-center gap-1.5 rounded-lg border px-3 py-2 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500 ${
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

        {/* 舊版模板升級提示（條件顯示） */}
        {staleTemplateCount > 0 && (
          <div className="mt-3 rounded-lg border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-100">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <p>
                {t('smartDashboard.sourcePlanner.staleTemplates', { count: staleTemplateCount, version: sourceTemplateSchemaVersion })}
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

        {/* 狀態列：驗證錯誤、配置訊息、已存模板快捷列 */}
        <div className="mt-3 flex flex-wrap items-center gap-2">
          {!typedPlanValidation.valid && (
            <span className="text-xs text-rose-300">{t('smartDashboard.sourcePlanner.validationError')}</span>
          )}
          {allocationMessage && (
            <span className="text-xs text-sky-200">{allocationMessage}</span>
          )}
          {sourceTemplates.slice(0, 6).map((template) => (
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
              <span className="text-[10px] text-slate-400">
                {new Date(template.lastUsedAt).toLocaleDateString()}
              </span>
              <button
                type="button"
                onClick={() => handleDeleteTemplate(template.id)}
                className="rounded-full px-1 text-slate-400 hover:bg-slate-700 hover:text-white"
                aria-label={t('smartDashboard.sourcePlanner.deleteTemplateAria', { name: template.name })}
              >
                ×
              </button>
            </div>
          ))}
        </div>
      </section>
      <SmartDashboardFlowStatusSection
        flowSegments={flowSegments}
        flowState={flowState}
        statusStyle={statusStyle}
        hasError={hasError}
        t={t}
      />
      <section
        ref={gridSectionRef}
        tabIndex={-1}
        className={`rounded-2xl border border-white/10 bg-slate-900/50 transition-all ${resolveIntentMotionClass(guideStage, reducedMotion)}`}
        style={{ transitionDuration: `${motionTokens.stageHandoffMs}ms` }}
      >
        <div className="border-b border-white/10 px-4 py-2.5">
          <div className="flex flex-wrap items-center gap-2">
            <h3 className="text-sm font-semibold text-slate-100">{t('smartDashboard.memoryGrid.title')}</h3>
            <div className="flex flex-wrap items-center gap-1.5">
              <span className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] font-medium ${
                modbusStatus?.bind_state === 'pass'
                  ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
                  : 'border-rose-500/30 bg-rose-500/10 text-rose-300'
              }`}>
                {t('smartDashboard.memoryGrid.bindState')}
                <span className="font-bold">{(modbusStatus?.bind_state ?? 'fail').toUpperCase()}</span>
              </span>
              <span className="inline-flex items-center gap-1 rounded-full border border-white/10 bg-slate-800/60 px-2 py-0.5 text-[11px] text-slate-300">
                {t('smartDashboard.memoryGrid.mappingCount')}
                <span className="font-semibold text-slate-100">{modbusStatus?.mapping_count ?? 0}</span>
              </span>
              <span className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] ${
                planConflictCount > 0
                  ? 'border-amber-500/30 bg-amber-500/10 text-amber-300'
                  : 'border-white/10 bg-slate-800/60 text-slate-300'
              }`}>
                {t('smartDashboard.memoryGrid.conflictCount')}
                <span className="font-semibold">{planConflictCount}</span>
              </span>
            </div>
            <div className="ml-auto">
              <button
                type="button"
                onClick={goToLocalModbusWorkbench}
                className="rounded-lg border border-blue-500/40 bg-blue-500/15 px-3 py-1.5 text-xs font-semibold text-blue-100 hover:bg-blue-500/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              >
                {t('smartDashboard.memoryGrid.fullWorkbench')}
              </button>
            </div>
          </div>
        </div>
        <div className="max-h-[340px] overflow-auto md:max-h-[360px] lg:max-h-[320px]">
          <MemoryGrid
            deviceId={selectedDevice.id}
            protocol={selectedDevice.protocol}
            centerAddress={planStartAddress || getGridCenterAddress(selectedDevice.protocol)}
            range={220}
            existingPoints={allPoints}
            linkedAddresses={linkedAddresses}
            selectedAddresses={selectedAddresses}
            plannedAllocations={plannedAllocations}
            showConflictsOnly={showConflictsOnly}
            onSelect={setSelectedAddresses}
            onCellClick={handleCellClick}
          />
        </div>
      </section>
    </div>
  );
}
