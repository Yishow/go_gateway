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
        <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
          <h3 className="text-sm font-semibold text-slate-100">Source Planner</h3>
          <span className="text-xs text-slate-400">
            {planDataType} x {planCount} = {totalPlannedCells} cells
          </span>
        </div>
        <div className="grid grid-cols-1 gap-3 xl:grid-cols-[1fr_1fr_auto_auto]">
          <label className="order-1 text-xs text-slate-300 xl:order-1">
            Start Address
            <input
              value={planStartAddress}
              onChange={(e) => setPlanStartAddress(e.target.value.toUpperCase())}
              placeholder="例如: 40001"
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </label>
          <label className="order-2 text-xs text-slate-300 xl:order-2">
            Source Count
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
          <label className="order-3 text-xs text-slate-300 xl:order-3">
            Data Type
            <select
              value={planDataType}
              onChange={(e) => setPlanDataType(e.target.value as DataType)}
              className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="int16">int16 (1 cell)</option>
              <option value="int32">int32 (2 cells)</option>
              <option value="float32">float32 (2 cells)</option>
              <option value="int64">int64 (4 cells)</option>
              <option value="float64">float64 (4 cells)</option>
            </select>
          </label>
          <div className="order-4 flex items-end gap-2 xl:order-4">
            <button
              type="button"
              onClick={handleAutoAllocate}
              className="inline-flex min-h-11 items-center gap-2 rounded-lg border border-sky-500/40 bg-sky-500/15 px-3 py-2 text-xs font-medium text-sky-100 hover:bg-sky-500/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500"
            >
              <WandSparkles className="h-4 w-4" />
              Auto
            </button>
            <button
              type="button"
              onClick={handleApplyPlan}
              className="inline-flex min-h-11 items-center gap-2 rounded-lg border border-emerald-500/40 bg-emerald-500/20 px-3 py-2 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
            >
              套用到 Grid
            </button>
          </div>
        </div>
        <div className="mt-3 grid grid-cols-1 gap-3 xl:grid-cols-[1fr_auto]">
          <label className="text-xs text-slate-300">
            儲存模板
            <div className="mt-1 flex gap-2">
              <input
                value={templateName}
                onChange={(e) => setTemplateName(e.target.value)}
                placeholder="例如: line-a-float32-10"
                className="w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="button"
                onClick={handleSaveTemplate}
                disabled={!templateName.trim()}
                className="inline-flex min-h-11 items-center gap-2 rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 disabled:cursor-not-allowed disabled:opacity-40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
              >
                <Save className="h-4 w-4" />
                儲存
              </button>
            </div>
          </label>
          <div className="flex items-end gap-2">
            <button
              type="button"
              onClick={() => setShowConflictsOnly((prev) => !prev)}
              className={`inline-flex min-h-11 items-center gap-2 rounded-lg border px-3 py-2 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500 ${
                showConflictsOnly
                  ? 'border-amber-500/40 bg-amber-500/20 text-amber-100'
                  : 'border-slate-700 bg-slate-800/70 text-slate-200'
              }`}
            >
              <Filter className="h-4 w-4" />
              只看衝突
            </button>
          </div>
        </div>
        <div className="mt-3 rounded-lg border border-white/10 bg-slate-800/40 p-3">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <label className="text-xs text-slate-300">
              批次命名前綴
              <input
                value={batchNamePrefix}
                onChange={(e) => setBatchNamePrefix(normalizeNamingPrefix(e.target.value))}
                className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="例如: LINEA"
              />
            </label>
            <div className="text-xs text-slate-400">
              預覽 {namePreview.length} 筆，名稱衝突 {nameConflictCount} 筆
            </div>
          </div>
          <div className="mt-2 grid grid-cols-2 gap-2 md:grid-cols-3 xl:grid-cols-4">
            {namePreview.slice(0, 12).map((item) => (
              <div
                key={item.sequence}
                className={`rounded-md border px-2 py-1 text-[11px] font-mono ${
                  item.conflict
                    ? 'border-rose-500/40 bg-rose-500/10 text-rose-200'
                    : 'border-slate-700 bg-slate-900/70 text-slate-200'
                }`}
              >
                {item.name}
              </div>
            ))}
          </div>
          <p className="mt-2 text-[11px] text-slate-500">
            命名規則: 僅允許英數、底線、連字號；系統會自動轉大寫並附加三位流水號。
          </p>
        </div>
        {staleTemplateCount > 0 && (
          <div className="mt-3 rounded-lg border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-100">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <p>
                偵測到 {staleTemplateCount} 個舊版模板，建議升級到 v{sourceTemplateSchemaVersion} 以確保流程一致性。
              </p>
              <button
                type="button"
                onClick={handleUpgradeTemplates}
                className="rounded-md border border-amber-400/40 bg-amber-500/20 px-2 py-1 text-[11px] font-semibold hover:bg-amber-500/30"
              >
                升級模板
              </button>
            </div>
          </div>
        )}
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <span className="text-xs text-slate-400">衝突格數: {planConflictCount}</span>
          {!typedPlanValidation.valid && (
            <span className="text-xs text-rose-300">來源數量需介於 1 到 200，且型別占格規則必須有效。</span>
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
                aria-label={`Delete template ${template.name}`}
              >
                x
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
        <div className="border-b border-white/10 px-4 py-3">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
              <h3 className="text-sm font-semibold text-slate-100">Dashboard Memory Grid</h3>
              <p className="mt-1 text-[11px] text-slate-400">
                bind_state / mapping_count / conflict_count（預檢一致）
              </p>
            </div>
            <button
              type="button"
              onClick={goToLocalModbusWorkbench}
              className="min-h-9 rounded-lg border border-blue-500/40 bg-blue-500/15 px-3 py-2 text-xs font-semibold text-blue-100 hover:bg-blue-500/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              完整工作台
            </button>
          </div>
          <div className="mt-2 grid grid-cols-1 gap-2 text-[11px] sm:grid-cols-3">
            <div className="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-slate-300">
              bind_state:
              <span className={`ml-1 font-semibold ${modbusStatus?.bind_state === 'pass' ? 'text-emerald-300' : 'text-rose-300'}`}>
                {(modbusStatus?.bind_state ?? 'fail').toUpperCase()}
              </span>
            </div>
            <div className="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-slate-300">
              mapping_count:
              <span className="ml-1 font-semibold text-slate-100">{modbusStatus?.mapping_count ?? 0}</span>
            </div>
            <div className="rounded-lg border border-white/10 bg-slate-800/60 px-3 py-2 text-slate-300">
              conflict_count:
              <span className={`ml-1 font-semibold ${planConflictCount > 0 ? 'text-amber-300' : 'text-emerald-300'}`}>
                {planConflictCount}
              </span>
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
