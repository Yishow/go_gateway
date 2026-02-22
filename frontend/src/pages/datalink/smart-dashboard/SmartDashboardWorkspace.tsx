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
  handleCellClick: (address: string) => void;
  getGridCenterAddress: (protocol: ProtocolType) => string;
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
  getGridCenterAddress,
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
            getGridCenterAddress={getGridCenterAddress}
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
