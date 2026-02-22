import { WandSparkles } from 'lucide-react';
import type { RefObject } from 'react';
import type { TFunction } from 'i18next';
import { MemoryGrid } from '../../../components/datalink/MemoryGrid';
import type { Point, DataType, ProtocolType, ModbusShareStatus } from '../../../types/datalink';
import type { PlannedAllocation } from '../../../components/datalink/MemoryGrid';
import type { FlowState, FlowSegment } from '../../../features/flow/stateMachine';
import SmartDashboardFlowStatusSection from './SmartDashboardFlowStatusSection';

type IntentStage = 'idle' | 'grid' | 'commit';

interface SmartDashboardWorkspaceContentProps {
  /** Source Planner 緊湊列 */
  planDataType: DataType;
  setPlanDataType: (value: DataType) => void;
  planCount: number;
  setPlanCount: (value: number) => void;
  totalPlannedCells: number;
  planStartAddress: string;
  setPlanStartAddress: (value: string) => void;
  handleAutoAllocate: () => void;
  handleApplyPlan: () => void;
  planConflictCount: number;
  typedPlanValidation: { valid: boolean };
  /** 衝突過濾（由規劃 Tab 控制） */
  showConflictsOnly: boolean;
  /** Flow Status */
  flowSegments: readonly FlowSegment[];
  flowState: FlowState;
  statusStyle: Record<string, string>;
  hasError: boolean;
  t: TFunction;
  /** Memory Grid */
  gridSectionRef: RefObject<HTMLElement | null>;
  resolveIntentMotionClass: (stage: IntentStage, reducedMotion: boolean) => string;
  guideStage: IntentStage;
  reducedMotion: boolean;
  motionTokens: { stageHandoffMs: number };
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
  planConflictCount,
  typedPlanValidation,
  showConflictsOnly,
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
    <div className="flex flex-1 flex-col overflow-hidden">
      {/* ── Flow Status（置頂） ── */}
      <div className="shrink-0 border-b border-white/5">
        <SmartDashboardFlowStatusSection
          flowSegments={flowSegments}
          flowState={flowState}
          statusStyle={statusStyle}
          hasError={hasError}
          t={t}
        />
      </div>

      {/* ── Source Planner 緊湊列 ── */}
      <div className="shrink-0 border-b border-white/5 px-4 py-2.5">
        <div className="flex flex-wrap items-center gap-2">
          {/* 位址 */}
          <label className="flex items-center gap-1.5 text-xs text-slate-400">
            {t('smartDashboard.sourcePlanner.startAddress')}
            <input
              value={planStartAddress}
              onChange={(e) => setPlanStartAddress(e.target.value.toUpperCase())}
              placeholder={t('smartDashboard.sourcePlanner.startAddressPlaceholder')}
              className="w-24 rounded-md border border-slate-700 bg-slate-800/80 px-2 py-1 font-mono text-xs text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </label>
          {/* 數量 */}
          <label className="flex items-center gap-1.5 text-xs text-slate-400">
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
              className="w-16 rounded-md border border-slate-700 bg-slate-800/80 px-2 py-1 text-xs text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </label>
          {/* 型別 */}
          <select
            value={planDataType}
            onChange={(e) => setPlanDataType(e.target.value as DataType)}
            className="rounded-md border border-slate-700 bg-slate-800/80 px-2 py-1 text-xs text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="int16">{t('smartDashboard.sourcePlanner.dataTypeOption_int16')}</option>
            <option value="int32">{t('smartDashboard.sourcePlanner.dataTypeOption_int32')}</option>
            <option value="float32">{t('smartDashboard.sourcePlanner.dataTypeOption_float32')}</option>
            <option value="int64">{t('smartDashboard.sourcePlanner.dataTypeOption_int64')}</option>
            <option value="float64">{t('smartDashboard.sourcePlanner.dataTypeOption_float64')}</option>
          </select>
          {/* 摘要 */}
          <span className="rounded-full bg-slate-800/80 px-2 py-0.5 text-[11px] text-slate-400 ring-1 ring-white/10">
            = {totalPlannedCells} 格
          </span>
          {planConflictCount > 0 && (
            <span className="rounded-full bg-amber-500/15 px-2 py-0.5 text-[11px] text-amber-300 ring-1 ring-amber-500/30">
              衝突 {planConflictCount}
            </span>
          )}
          {!typedPlanValidation.valid && (
            <span className="text-[11px] text-rose-400">
              {t('smartDashboard.sourcePlanner.validationError')}
            </span>
          )}
          {/* 動作按鈕 */}
          <div className="ml-auto flex items-center gap-2">
            <button
              type="button"
              onClick={handleAutoAllocate}
              className="inline-flex items-center gap-1.5 rounded-md border border-sky-500/40 bg-sky-500/15 px-2.5 py-1 text-xs font-medium text-sky-100 hover:bg-sky-500/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500"
            >
              <WandSparkles className="h-3.5 w-3.5" />
              {t('smartDashboard.sourcePlanner.auto')}
            </button>
            <button
              type="button"
              onClick={handleApplyPlan}
              className="inline-flex items-center gap-1.5 rounded-md border border-emerald-500/40 bg-emerald-500/20 px-2.5 py-1 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
            >
              {t('smartDashboard.sourcePlanner.applyToGrid')}
            </button>
          </div>
        </div>
      </div>

      {/* ── Memory Grid ── */}
      <section
        ref={gridSectionRef}
        tabIndex={-1}
        className={`flex-1 overflow-hidden transition-all ${resolveIntentMotionClass(guideStage, reducedMotion)}`}
        style={{ transitionDuration: `${motionTokens.stageHandoffMs}ms` }}
      >
        {/* Grid 標頭：狀態 badges + 簡化 */}
        <div className="flex flex-wrap items-center gap-1.5 border-b border-white/5 px-4 py-2">
          <span className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] font-medium ${
            modbusStatus?.bind_state === 'pass'
              ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
              : 'border-rose-500/30 bg-rose-500/10 text-rose-300'
          }`}>
            bind {(modbusStatus?.bind_state ?? 'fail').toUpperCase()}
          </span>
          <span className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] ${
            planConflictCount > 0
              ? 'border-amber-500/30 bg-amber-500/10 text-amber-300'
              : 'border-white/10 bg-slate-800/60 text-slate-400'
          }`}>
            衝突 {planConflictCount}
          </span>
          <span className="ml-auto text-xs font-semibold text-slate-100">
            {t('smartDashboard.memoryGrid.title')}
          </span>
        </div>
        <div className="h-[calc(100%-40px)] overflow-auto">
          <MemoryGrid
            deviceId={selectedDevice.id}
            protocol={selectedDevice.protocol}
            centerAddress={planStartAddress || getGridCenterAddress(selectedDevice.protocol)}
            range={100}
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
