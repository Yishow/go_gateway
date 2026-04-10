import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type Dispatch,
  type ReactNode,
  type SetStateAction,
} from 'react';
import {
  type DevicePanelState,
  type DeviceTestHistoryEntry,
  INSPECTOR_SELECTION_NONE,
  WORKBENCH_CROSS_STEP_CONTEXT_INITIAL,
  WORKBENCH_OUTPUT_SELECTION_INITIAL,
  WORKBENCH_SOURCE_PLANNING_INITIAL,
  WORKBENCH_STEPS,
  type InspectorSelection,
  type OutputTarget,
  type WorkbenchCrossStepContext,
    type WorkbenchOutputSelectionState,
    type WorkbenchSourcePlanningState,
    type WorkbenchStep,
    type WorkbenchTagGroupingOverride,
    type WorkbenchTagGroupingOverrideMap,
    WORKBENCH_TAG_GROUPING_OVERRIDES_INITIAL,
  } from './workbenchTypes';

type WorkbenchContextValue = {
  // Step navigation
  steps: readonly WorkbenchStep[];
  activeStep: WorkbenchStep;
  setActiveStep: (step: WorkbenchStep) => void;

  // Device context (cross-step)
  selectedDeviceId: string | null;
  setSelectedDeviceId: (deviceId: string | null) => void;

  // Inspector selection — cleared on every step navigation
  inspectorSelection: InspectorSelection;
  setInspectorSelection: (selection: InspectorSelection) => void;
  clearInspectorSelection: () => void;

  // Output target (Step 4)
  activeOutputTarget: OutputTarget;
  setActiveOutputTarget: (target: OutputTarget) => void;
  outputSelectionState: WorkbenchOutputSelectionState;
  setOutputSelectionState: Dispatch<SetStateAction<WorkbenchOutputSelectionState>>;
  setSelectedOutputTagId: (target: OutputTarget, tagId: string) => void;
  clearOutputSelectionState: () => void;

  // Cross-step traceability context — persists across step navigation (spec §8.1)
  crossStepContext: WorkbenchCrossStepContext;
  setFocusedRuleId: (ruleId: string | null) => void;
  setFocusedTagIds: (tagIds: ReadonlyArray<string>) => void;
  clearCrossStepContext: () => void;

  // Step 2 shared source planning state — persists across step navigation
  sourcePlanningState: WorkbenchSourcePlanningState;
  setSourcePlanningState: Dispatch<SetStateAction<WorkbenchSourcePlanningState>>;
  setSourcePlannerStartAddress: (deviceId: string, startAddress: string) => void;
  clearSourcePlanningState: () => void;

  // Step 3 grouped database suggestion overrides — persist into Output planning
  tagGroupingOverrides: WorkbenchTagGroupingOverrideMap;
  setTagGroupingOverride: (
    ruleId: string,
    pointId: string,
    override: WorkbenchTagGroupingOverride,
  ) => void;
  clearTagGroupingOverride: (ruleId: string, pointId: string) => void;
  clearTagGroupingOverrides: () => void;

  // Step 1 shared UI state
  devicePanelState: DevicePanelState;
  openCreateDevicePanel: () => void;
  openEditDevicePanel: (deviceId: string) => void;
  openCloneDevicePanel: (deviceId: string) => void;
  closeDevicePanel: () => void;
  recentDeviceTests: Readonly<Record<string, ReadonlyArray<DeviceTestHistoryEntry>>>;
  recordDeviceTest: (
    deviceId: string,
    entry: Omit<DeviceTestHistoryEntry, 'id'>,
  ) => void;

  /** 來源步驟表面訊息（規則 API、收集開關等），供 Context Bar 與畫布區共用。 */
  sourceStepNotice: string | null;
  setSourceStepNotice: Dispatch<SetStateAction<string | null>>;

  /**
   * 來源步驟且已選設備時，由來源畫布區註冊至右欄檢查面板上方的 Runtime／規則摘要區；
   * 離開來源步驟或卸載時應清空。
   */
  sourceStepInspectorBanner: ReactNode | null;
  setSourceStepInspectorBanner: Dispatch<SetStateAction<ReactNode | null>>;
};

const WorkbenchContext = createContext<WorkbenchContextValue | null>(null);

export function WorkbenchProvider({ children }: { children: ReactNode }) {
  const [activeStep, setActiveStepRaw] = useState<WorkbenchStep>('device');
  const [selectedDeviceId, setSelectedDeviceIdRaw] = useState<string | null>(null);
  const [inspectorSelection, setInspectorSelection] = useState<InspectorSelection>(
    INSPECTOR_SELECTION_NONE,
  );
  const [activeOutputTarget, setActiveOutputTarget] = useState<OutputTarget>('modbus');
  const [outputSelectionState, setOutputSelectionState] = useState<WorkbenchOutputSelectionState>(
    WORKBENCH_OUTPUT_SELECTION_INITIAL,
  );
  const [crossStepContext, setCrossStepContext] = useState<WorkbenchCrossStepContext>(
    WORKBENCH_CROSS_STEP_CONTEXT_INITIAL,
  );
  const [sourcePlanningState, setSourcePlanningState] = useState<WorkbenchSourcePlanningState>(
    WORKBENCH_SOURCE_PLANNING_INITIAL,
  );
  const [tagGroupingOverrides, setTagGroupingOverrides] = useState<WorkbenchTagGroupingOverrideMap>(
    WORKBENCH_TAG_GROUPING_OVERRIDES_INITIAL,
  );
  const [devicePanelState, setDevicePanelState] = useState<DevicePanelState>(null);
  const [recentDeviceTests, setRecentDeviceTests] = useState<
    Record<string, ReadonlyArray<DeviceTestHistoryEntry>>
  >({});
  const [sourceStepNotice, setSourceStepNotice] = useState<string | null>(null);
  const [sourceStepInspectorBanner, setSourceStepInspectorBanner] = useState<ReactNode | null>(
    null,
  );

  useEffect(() => {
    setSourceStepNotice(null);
  }, [selectedDeviceId]);

  const clearInspectorSelection = useCallback(() => {
    setInspectorSelection(INSPECTOR_SELECTION_NONE);
  }, []);

  const closeDevicePanel = useCallback(() => {
    setDevicePanelState(null);
  }, []);

  const clearCrossStepContext = useCallback(() => {
    setCrossStepContext(WORKBENCH_CROSS_STEP_CONTEXT_INITIAL);
  }, []);

  const clearOutputSelectionState = useCallback(() => {
    setOutputSelectionState(WORKBENCH_OUTPUT_SELECTION_INITIAL);
  }, []);

  const setSelectedOutputTagId = useCallback((target: OutputTarget, tagId: string) => {
    setOutputSelectionState((currentState) => ({
      ...currentState,
      [target]: tagId,
    }));
  }, []);

  const clearSourcePlanningState = useCallback(() => {
    setSourcePlanningState((currentState) => ({
      ...WORKBENCH_SOURCE_PLANNING_INITIAL,
      rules: currentState.rules,
      plannerStartAddressByDeviceId: currentState.plannerStartAddressByDeviceId,
    }));
  }, []);

  const clearTagGroupingOverrides = useCallback(() => {
    setTagGroupingOverrides(WORKBENCH_TAG_GROUPING_OVERRIDES_INITIAL);
  }, []);

  const setTagGroupingOverride = useCallback(
    (ruleId: string, pointId: string, override: WorkbenchTagGroupingOverride) => {
      setTagGroupingOverrides((currentState) => ({
        ...currentState,
        [ruleId]: {
          ...(currentState[ruleId] ?? {}),
          [pointId]: override,
        },
      }));
    },
    [],
  );

  const clearTagGroupingOverride = useCallback((ruleId: string, pointId: string) => {
    setTagGroupingOverrides((currentState) => {
      const currentRuleOverrides = currentState[ruleId];
      if (!currentRuleOverrides || !currentRuleOverrides[pointId]) {
        return currentState;
      }

      const { [pointId]: _removed, ...remainingRuleOverrides } = currentRuleOverrides;
      if (Object.keys(remainingRuleOverrides).length === 0) {
        const { [ruleId]: _removedRule, ...remainingOverrides } = currentState;
        return remainingOverrides;
      }

      return {
        ...currentState,
        [ruleId]: remainingRuleOverrides,
      };
    });
  }, []);

  const setSourcePlannerStartAddress = useCallback((deviceId: string, startAddress: string) => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      plannerStartAddressByDeviceId: {
        ...currentState.plannerStartAddressByDeviceId,
        [deviceId]: startAddress,
      },
    }));
  }, []);

  // Clear inspector selection when switching steps so stale context
  // from a previous step doesn't leak into the new one.
  // Cross-step context is intentionally NOT cleared here — it persists.
  const setActiveStep = useCallback(
    (step: WorkbenchStep) => {
      setActiveStepRaw(step);
      if (step !== 'source') {
        setSourceStepInspectorBanner(null);
      }
      clearInspectorSelection();
      closeDevicePanel();
    },
    [clearInspectorSelection, closeDevicePanel],
  );

  // When a new device is selected, reset the cross-step context since
  // the source planning context is now stale.
  const setSelectedDeviceId = useCallback(
    (deviceId: string | null) => {
      setSelectedDeviceIdRaw(deviceId);
      clearCrossStepContext();
      clearOutputSelectionState();
      clearSourcePlanningState();
      clearTagGroupingOverrides();
      closeDevicePanel();
    },
    [
      clearCrossStepContext,
      clearOutputSelectionState,
      clearSourcePlanningState,
      clearTagGroupingOverrides,
      closeDevicePanel,
    ],
  );

  const setFocusedRuleId = useCallback((ruleId: string | null) => {
    setCrossStepContext((prev) => ({ ...prev, focusedRuleId: ruleId }));
  }, []);

  const setFocusedTagIds = useCallback((tagIds: ReadonlyArray<string>) => {
    setCrossStepContext((prev) => ({ ...prev, focusedTagIds: tagIds }));
  }, []);

  const openCreateDevicePanel = useCallback(() => {
    setDevicePanelState({ mode: 'create' });
  }, []);

  const openEditDevicePanel = useCallback((deviceId: string) => {
    setDevicePanelState({ mode: 'edit', deviceId });
  }, []);

  const openCloneDevicePanel = useCallback((deviceId: string) => {
    setDevicePanelState({ mode: 'clone', sourceDeviceId: deviceId });
  }, []);

  const recordDeviceTest = useCallback(
    (deviceId: string, entry: Omit<DeviceTestHistoryEntry, 'id'>) => {
      setRecentDeviceTests((currentState) => {
        const nextEntry: DeviceTestHistoryEntry = {
          ...entry,
          id: `${deviceId}-${entry.testedAt}-${currentState[deviceId]?.length ?? 0}`,
        };
        const nextEntries = [nextEntry, ...(currentState[deviceId] ?? [])].slice(0, 3);
        return {
          ...currentState,
          [deviceId]: nextEntries,
        };
      });
    },
    [],
  );

  const value = useMemo<WorkbenchContextValue>(
    () => ({
      steps: WORKBENCH_STEPS,
      activeStep,
      setActiveStep,
      selectedDeviceId,
      setSelectedDeviceId,
      inspectorSelection,
      setInspectorSelection,
      clearInspectorSelection,
      activeOutputTarget,
      setActiveOutputTarget,
      outputSelectionState,
      setOutputSelectionState,
      setSelectedOutputTagId,
      clearOutputSelectionState,
      crossStepContext,
      setFocusedRuleId,
      setFocusedTagIds,
      clearCrossStepContext,
      sourcePlanningState,
      setSourcePlanningState,
      setSourcePlannerStartAddress,
      clearSourcePlanningState,
      tagGroupingOverrides,
      setTagGroupingOverride,
      clearTagGroupingOverride,
      clearTagGroupingOverrides,
      devicePanelState,
      openCreateDevicePanel,
      openEditDevicePanel,
      openCloneDevicePanel,
      closeDevicePanel,
      recentDeviceTests,
      recordDeviceTest,
      sourceStepNotice,
      setSourceStepNotice,
      sourceStepInspectorBanner,
      setSourceStepInspectorBanner,
    }),
    [
      activeStep,
      setActiveStep,
      selectedDeviceId,
      setSelectedDeviceId,
      inspectorSelection,
      clearInspectorSelection,
      activeOutputTarget,
      outputSelectionState,
      setOutputSelectionState,
      setSelectedOutputTagId,
      clearOutputSelectionState,
      crossStepContext,
      setFocusedRuleId,
      setFocusedTagIds,
      clearCrossStepContext,
      sourcePlanningState,
      setSourcePlanningState,
      setSourcePlannerStartAddress,
      clearSourcePlanningState,
      tagGroupingOverrides,
      setTagGroupingOverride,
      clearTagGroupingOverride,
      clearTagGroupingOverrides,
      devicePanelState,
      openCreateDevicePanel,
      openEditDevicePanel,
      openCloneDevicePanel,
      closeDevicePanel,
      recentDeviceTests,
      recordDeviceTest,
      sourceStepNotice,
      sourceStepInspectorBanner,
    ],
  );

  return <WorkbenchContext.Provider value={value}>{children}</WorkbenchContext.Provider>;
}

export function useWorkbench() {
  const context = useContext(WorkbenchContext);

  if (!context) {
    throw new Error('useWorkbench must be used within WorkbenchProvider');
  }

  return context;
}
