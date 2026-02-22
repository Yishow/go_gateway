import type { Dispatch, RefObject, SetStateAction } from 'react';
import type { TFunction } from 'i18next';
import type { Device, Point, DataType, ProtocolType, ModbusShareStatus } from '../../../types/datalink';
import type { PlannedAllocation } from '../../../components/datalink/MemoryGrid';
import type { FlowState, FlowSegment } from '../../../features/flow/stateMachine';
import SmartDashboardWorkspaceContent from './SmartDashboardWorkspaceContent';
import SmartDashboardWorkspaceEmptyState from './SmartDashboardWorkspaceEmptyState';

type IntentStage = 'idle' | 'grid' | 'commit';

interface SmartDashboardWorkspaceProps {
  selectedDevice: Device | null;
  /** Source Planner */
  planDataType: DataType;
  setPlanDataType: Dispatch<SetStateAction<DataType>>;
  planCount: number;
  setPlanCount: (value: number) => void;
  totalPlannedCells: number;
  planStartAddress: string;
  setPlanStartAddress: (value: string) => void;
  handleAutoAllocate: () => void;
  handleApplyPlan: () => void;
  planConflictCount: number;
  typedPlanValidation: { valid: boolean };
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
  allPoints: Point[];
  linkedAddresses: string[];
  selectedAddresses: string[];
  plannedAllocations: PlannedAllocation[];
  setSelectedAddresses: (addresses: string[]) => void;
  handleCellClick: (address: string, point?: Point, e?: React.MouseEvent) => void;
  onCellContextMenu?: (address: string, point: Point, e: React.MouseEvent) => void;
  getGridCenterAddress: (protocol: ProtocolType) => string;
  /** Grid 視窗起始位址（可與 planStartAddress 不同，用於滾輪/按鈕切換 100 格） */
  gridViewStartAddress: string;
  onGridViewShift: (delta: number) => void;
  handleChooseDevice: () => void;
  handleCreateDevice: () => void;
}

export default function SmartDashboardWorkspace({
  selectedDevice,
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
  handleChooseDevice,
  handleCreateDevice,
}: SmartDashboardWorkspaceProps) {
  return (
    <div className="min-w-0 flex flex-col">
      <div className="relative flex flex-1 flex-col overflow-hidden rounded-3xl border border-white/5 bg-slate-900/40 shadow-2xl backdrop-blur-md">
        {selectedDevice ? (
          <SmartDashboardWorkspaceContent
            planDataType={planDataType}
            setPlanDataType={setPlanDataType}
            planCount={planCount}
            setPlanCount={setPlanCount}
            totalPlannedCells={totalPlannedCells}
            planStartAddress={planStartAddress}
            setPlanStartAddress={setPlanStartAddress}
            handleAutoAllocate={handleAutoAllocate}
            handleApplyPlan={handleApplyPlan}
            planConflictCount={planConflictCount}
            typedPlanValidation={typedPlanValidation}
            showConflictsOnly={showConflictsOnly}
            flowSegments={flowSegments}
            flowState={flowState}
            statusStyle={statusStyle}
            hasError={hasError}
            t={t}
            gridSectionRef={gridSectionRef}
            resolveIntentMotionClass={resolveIntentMotionClass}
            guideStage={guideStage}
            reducedMotion={reducedMotion}
            motionTokens={motionTokens}
            modbusStatus={modbusStatus}
            selectedDevice={selectedDevice}
            allPoints={allPoints}
            linkedAddresses={linkedAddresses}
            selectedAddresses={selectedAddresses}
            plannedAllocations={plannedAllocations}
            setSelectedAddresses={setSelectedAddresses}
            handleCellClick={handleCellClick}
            onCellContextMenu={onCellContextMenu}
            getGridCenterAddress={getGridCenterAddress}
            gridViewStartAddress={gridViewStartAddress}
            onGridViewShift={onGridViewShift}
          />
        ) : (
          <SmartDashboardWorkspaceEmptyState
            t={t}
            onChooseDevice={handleChooseDevice}
            onCreateDevice={handleCreateDevice}
          />
        )}
      </div>
    </div>
  );
}
