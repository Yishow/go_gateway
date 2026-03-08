import { ChevronLeft, ChevronRight, WandSparkles, ChevronDown, ChevronUp } from "lucide-react";
import { useState } from "react";
import type { RefObject } from "react";
import type { TFunction } from "i18next";
import { MemoryGrid } from "../../../components/datalink/MemoryGrid";
import type { Point, DataType, ProtocolType, ModbusShareStatus } from "../../../types/datalink";
import type { PlannedAllocation } from "../../../components/datalink/MemoryGrid";
import type { FlowState, FlowSegment } from "../../../features/flow/stateMachine";
import SmartDashboardFlowStatusSection from "./SmartDashboardFlowStatusSection";
import { useSmartDashboardWorkspaceContentState } from "./useSmartDashboardWorkspaceContentState";
import { designSystem } from "../../../styles/designSystem";

type IntentStage = "idle" | "grid" | "commit";

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
  handleCellClick: (address: string, point?: Point, e?: React.MouseEvent) => void;
  onCellContextMenu?: (address: string, point: Point, e: React.MouseEvent) => void;
  getGridCenterAddress: (protocol: ProtocolType) => string;
  /** Grid 視窗起始位址（滾輪/按鈕切換 100 格用） */
  gridViewStartAddress: string;
  onGridViewShift: (delta: number) => void;
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
  onCellContextMenu,
  getGridCenterAddress,
  gridViewStartAddress,
  onGridViewShift,
}: SmartDashboardWorkspaceContentProps) {
  const {
    gridScrollRef,
    gridRange,
    modbusArea,
    centerAddress,
    handleModbusAreaChange,
    handlePlanStartAddressChange,
    handlePlanCountChange,
  } = useSmartDashboardWorkspaceContentState({
    planStartAddress,
    setPlanStartAddress,
    setPlanCount,
    selectedDeviceProtocol: selectedDevice.protocol,
    gridViewStartAddress,
    getGridCenterAddress,
    onGridViewShift,
  });

  const [isFlowStatusExpanded, setIsFlowStatusExpanded] = useState(false);

  return (
    <div className="flex flex-1 flex-col overflow-hidden">
      {/* ── Flow Status（可收折，預設收合以聚焦主流程） ── */}
      <div className={`shrink-0 border-b border-white/5 transition-all duration-200 ${isFlowStatusExpanded ? '' : 'overflow-hidden'}`}>
        <button
          type="button"
          onClick={() => setIsFlowStatusExpanded(!isFlowStatusExpanded)}
          className="w-full flex items-center justify-between px-4 py-2 text-xs text-slate-400 hover:text-slate-200 hover:bg-slate-900/30 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          aria-expanded={isFlowStatusExpanded}
          aria-label={t('smartDashboard.flowStatusToggle', { defaultValue: '切換流程狀態顯示' })}
        >
          <span className="font-medium">{t('smartDashboard.flowTitle')}</span>
          {isFlowStatusExpanded ? (
            <ChevronUp className="h-4 w-4" />
          ) : (
            <ChevronDown className="h-4 w-4" />
          )}
        </button>
        {isFlowStatusExpanded && (
          <div className="px-4 pb-3">
            <SmartDashboardFlowStatusSection
              flowSegments={flowSegments}
              flowState={flowState}
              statusStyle={statusStyle}
              hasError={hasError}
              t={t}
            />
          </div>
        )}
      </div>

      {/* ── Source Planner 緊湊列 ── */}
      <div className="shrink-0 border-b border-white/5 px-2 sm:px-4 py-2.5">
        <div className="flex flex-wrap items-center gap-2">
          {/* Modbus 區域（僅 Modbus 設備顯示，切換 FC 會更新預設起始位址） */}
          {selectedDevice.protocol.startsWith("modbus") && (
            <label className="flex items-center gap-1.5 text-xs text-slate-400">
              {t("smartDashboard.sourcePlanner.modbusArea")}
              <select
                value={modbusArea}
                onChange={(e) => {
                  handleModbusAreaChange(e.target.value as "0" | "1" | "3" | "4");
                }}
                name="modbus-area"
                className="rounded-md border border-slate-700 bg-slate-800/80 px-2 py-1 text-xs text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                aria-label={t("smartDashboard.sourcePlanner.modbusArea")}
              >
                <option value="0">{t("smartDashboard.sourcePlanner.modbusAreaCoil")}</option>
                <option value="1">{t("smartDashboard.sourcePlanner.modbusAreaDiscrete")}</option>
                <option value="3">{t("smartDashboard.sourcePlanner.modbusAreaInputReg")}</option>
                <option value="4">{t("smartDashboard.sourcePlanner.modbusAreaHoldingReg")}</option>
              </select>
            </label>
          )}
          {/* 位址 */}
          <label className={`${designSystem.forms.label.base} flex items-center gap-1.5 text-xs`}>
            {t("smartDashboard.sourcePlanner.startAddress")}
            <input
              value={planStartAddress}
              onChange={(e) => handlePlanStartAddressChange(e.target.value)}
              placeholder={t("smartDashboard.sourcePlanner.startAddressPlaceholder")}
              name={designSystem.forms.name.pointAddress}
              autoComplete={designSystem.forms.autocomplete.text}
              inputMode={designSystem.forms.inputmode.text}
              className="w-24 rounded-md border border-slate-700 bg-slate-800/80 px-2 py-1 font-mono text-xs text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </label>
          {/* 數量 */}
          <label className={`${designSystem.forms.label.base} flex items-center gap-1.5 text-xs`}>
            {t("smartDashboard.sourcePlanner.sourceCount")}
            <input
              type="number"
              min={1}
              max={200}
              value={planCount}
              onChange={(e) => handlePlanCountChange(e.target.value)}
              name="plan-count"
              autoComplete={designSystem.forms.autocomplete.number}
              inputMode={designSystem.forms.inputmode.numeric}
              className="w-16 rounded-md border border-slate-700 bg-slate-800/80 px-2 py-1 text-xs text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </label>
          {/* 型別 */}
          <select
            value={planDataType}
            onChange={(e) => setPlanDataType(e.target.value as DataType)}
            name={designSystem.forms.name.pointDataType}
            className="rounded-md border border-slate-700 bg-slate-800/80 px-2 py-1 text-xs text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="int16">{t("smartDashboard.sourcePlanner.dataTypeOption_int16")}</option>
            <option value="int32">{t("smartDashboard.sourcePlanner.dataTypeOption_int32")}</option>
            <option value="float32">{t("smartDashboard.sourcePlanner.dataTypeOption_float32")}</option>
            <option value="int64">{t("smartDashboard.sourcePlanner.dataTypeOption_int64")}</option>
            <option value="float64">{t("smartDashboard.sourcePlanner.dataTypeOption_float64")}</option>
          </select>
          {/* 摘要 */}
          <span className="rounded-full bg-slate-800/80 px-2 py-0.5 text-[11px] text-slate-400 ring-1 ring-white/10">
            {t("smartDashboard.sourcePlanner.plannedCells", { count: totalPlannedCells })}
          </span>
          {planConflictCount > 0 && (
            <span className="rounded-full bg-amber-500/15 px-2 py-0.5 text-[11px] text-amber-300 ring-1 ring-amber-500/30">
              {t("smartDashboard.sourcePlanner.conflictCount", { count: planConflictCount })}
            </span>
          )}
          {!typedPlanValidation.valid && (
            <span className="text-[11px] text-rose-400">{t("smartDashboard.sourcePlanner.validationError")}</span>
          )}
          {/* 動作按鈕 */}
          <div className="ml-auto flex items-center gap-2">
            <button
              type="button"
              onClick={handleAutoAllocate}
              className="inline-flex items-center gap-1.5 rounded-md border border-sky-500/40 bg-sky-500/15 px-2.5 py-1 text-xs font-medium text-sky-100 hover:bg-sky-500/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500"
            >
              <WandSparkles className="h-3.5 w-3.5" />
              {t("smartDashboard.sourcePlanner.auto")}
            </button>
            <button
              type="button"
              onClick={handleApplyPlan}
              className="inline-flex items-center gap-1.5 rounded-md border border-emerald-500/40 bg-emerald-500/20 px-2.5 py-1 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
            >
              {t("smartDashboard.sourcePlanner.applyToGrid")}
            </button>
          </div>
        </div>
      </div>

      {/* ── Memory Grid ── */}
      <section
        ref={gridSectionRef}
        tabIndex={-1}
        className={`flex min-h-0 flex-1 flex-col overflow-hidden transition-all ${resolveIntentMotionClass(guideStage, reducedMotion)}`}
        style={{ transitionDuration: `${motionTokens.stageHandoffMs}ms` }}
      >
        {/* Grid 標頭：切換 100 格按鈕（左） + 狀態 badges + 標題（右） */}
        <div className="shrink-0 flex flex-wrap items-center gap-1.5 border-b border-white/5 px-4 py-2">
          <div className="flex items-center gap-1">
            <button
              type="button"
              onClick={() => onGridViewShift(-100)}
              className={`${designSystem.components.button.icon} rounded border border-slate-600 bg-slate-800/80 p-1 text-slate-300 hover:bg-slate-700 hover:text-slate-100`}
              title={t("smartDashboard.memoryGrid.prevRange")}
              aria-label={t("smartDashboard.memoryGrid.prevRange")}
            >
              <ChevronLeft className="h-4 w-4" />
            </button>
            <button
              type="button"
              onClick={() => onGridViewShift(100)}
              className={`${designSystem.components.button.icon} rounded border border-slate-600 bg-slate-800/80 p-1 text-slate-300 hover:bg-slate-700 hover:text-slate-100`}
              title={t("smartDashboard.memoryGrid.nextRange")}
              aria-label={t("smartDashboard.memoryGrid.nextRange")}
            >
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
          <span
            className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] font-medium ${
              modbusStatus?.bind_state === "pass"
                ? "border-emerald-500/30 bg-emerald-500/10 text-emerald-300"
                : "border-rose-500/30 bg-rose-500/10 text-rose-300"
            }`}
          >
            {t("smartDashboard.memoryGrid.bindBadge", { state: (modbusStatus?.bind_state ?? "fail").toUpperCase() })}
          </span>
          <span
            className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] ${
              planConflictCount > 0
                ? "border-amber-500/30 bg-amber-500/10 text-amber-300"
                : "border-white/10 bg-slate-800/60 text-slate-400"
            }`}
          >
            {t("smartDashboard.sourcePlanner.conflictCount", { count: planConflictCount })}
          </span>
          <span className="ml-auto text-xs font-semibold text-slate-100">{t("smartDashboard.memoryGrid.title")}</span>
        </div>
        <div
          ref={gridScrollRef}
          className="min-h-0 flex-1 flex flex-col overflow-auto"
          role="region"
          aria-label={t("smartDashboard.memoryGrid.title")}
        >
          <MemoryGrid
            deviceId={selectedDevice.id}
            protocol={selectedDevice.protocol}
            centerAddress={centerAddress}
            range={gridRange}
            existingPoints={allPoints}
            linkedAddresses={linkedAddresses}
            selectedAddresses={selectedAddresses}
            plannedAllocations={plannedAllocations}
            showConflictsOnly={showConflictsOnly}
            onSelect={setSelectedAddresses}
            onCellClick={handleCellClick}
            onCellContextMenu={onCellContextMenu}
          />
        </div>
      </section>
    </div>
  );
}
