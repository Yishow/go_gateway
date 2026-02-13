import type { Dispatch, RefObject, SetStateAction } from 'react';
import type { TFunction } from 'i18next';
import type { Device, Point, DataType, ProtocolType, ModbusShareStatus } from '../../../types/datalink';
import type { PlannedAllocation } from '../../../components/datalink/MemoryGrid';
import type { FlowState, FlowSegment } from '../../../features/flow/stateMachine';
import SmartDashboardWorkspaceContent from './SmartDashboardWorkspaceContent';
import SmartDashboardWorkspaceEmptyState from './SmartDashboardWorkspaceEmptyState';

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

interface SmartDashboardWorkspaceProps {
  selectedDevice: Device | null;
  planDataType: DataType;
  setPlanDataType: Dispatch<SetStateAction<DataType>>;
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
      <div className="relative flex min-h-[420px] flex-1 flex-col overflow-hidden rounded-3xl border border-white/5 bg-slate-900/40 shadow-2xl backdrop-blur-md">
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
            templateName={templateName}
            setTemplateName={setTemplateName}
            handleSaveTemplate={handleSaveTemplate}
            showConflictsOnly={showConflictsOnly}
            setShowConflictsOnly={setShowConflictsOnly}
            batchNamePrefix={batchNamePrefix}
            setBatchNamePrefix={setBatchNamePrefix}
            normalizeNamingPrefix={normalizeNamingPrefix}
            namePreview={namePreview}
            nameConflictCount={nameConflictCount}
            staleTemplateCount={staleTemplateCount}
            sourceTemplateSchemaVersion={sourceTemplateSchemaVersion}
            handleUpgradeTemplates={handleUpgradeTemplates}
            planConflictCount={planConflictCount}
            typedPlanValidation={typedPlanValidation}
            allocationMessage={allocationMessage}
            sourceTemplates={sourceTemplates}
            handleLoadTemplate={handleLoadTemplate}
            handleDeleteTemplate={handleDeleteTemplate}
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
            goToLocalModbusWorkbench={goToLocalModbusWorkbench}
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
