import { deviceAPI } from '../../../services/datalink';
import { useRecordingRangeSubmission } from './useRecordingRangeSubmission';
import { readSourceRangeDevice } from './readSourceRangeDevice';
import type { DeviceSourceRequest } from './deviceSourceSubmission';
import { useLocatorDeviceSourceLease } from './useLocatorDeviceSourceLease';
import { useDeviceSourceSubmission } from './useDeviceSourceSubmission';
import type { StudioLocatorEntryRead } from '../../../types/studioV2LocatorEntry';
import { useStudioV2LocatorHydration } from './studioV2LocatorHydration';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { useCreateStudioV2WorkspaceDeviceMutation, useDeleteStudioV2WorkspaceDeviceMutation, useStudioV2WorkspaceDevicesQuery, useUpdateStudioV2WorkspaceDeviceAvailabilityMutation, useUpdateStudioV2WorkspaceDeviceMutation } from '../../../hooks/datalink/useStudioV2WorkspaceDevices';
import { hydrateStudioV2Device, isStudioV2DeviceValid, toStudioV2DeviceCreateRequest, toStudioV2DeviceUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { canHydrateStudioV2Rule, hydrateStudioV2Rule } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import type { Device, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import { useWorkbenchV2State, workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { useStudioV2RulesQuery } from '../../../hooks/datalink/useStudioV2Rules';
import { useStudioV2RuleAutosave } from './useStudioV2RuleAutosave';
import { useStudioV2MappingAutosave } from './useStudioV2MappingAutosave';
import { useStudioV2DatabaseAutosave } from './useStudioV2DatabaseAutosave';
import { hydrateStudioV2DatabaseConnector, hydrateStudioV2DatabaseRowGroups } from '../../../features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import { deriveAllPoints } from '../../../features/datalink/workbench-v2/state/sourceRule';
import { isPersistedMCRule, mergePersistedMCState, withPersistedMCEvidence } from '../../../features/datalink/workbench-v2/state/mcPersistedPoints';
import { useStudioV2PersistedMCPoints, type PersistedMCWorkspaceScope } from './useStudioV2PersistedMCPoints';
import { inferHydratedProgress } from './hydratedProgress';
import type { ProtocolType } from '../../../types/datalink';
import { getStudioV2AutosaveBarrier, hasPendingSave, isUnpersistedNonIdle } from './studioV2AutosaveBarrier';
import { getSafeErrorStateMessage } from '../../../utils/typedErrors';
type SaveMeta = {
  inFlight: boolean;
  pending: boolean;
};
const draftLossStorageKey = 'wbv2_unrecovered_draft';
function saveMetaFor(store: Record<string, SaveMeta>, deviceId: string): SaveMeta {
  if (!store[deviceId]) {
    store[deviceId] = { inFlight: false, pending: false };
  }
  return store[deviceId];
}
function errorMessageOf(error: unknown, fallback: string): string {
  return getSafeErrorStateMessage(error, fallback);
}
const invalidDeviceAvailabilityReason = 'device form is invalid';
function isDraftPending(saveState?: string): boolean {
  return saveState === 'saving' || saveState === 'save-error' || saveState === 'draft-invalid';
}
function hasLocalUnpersistedSetupDrafts(state: WorkbenchV2State): boolean {
  if (state.devices.some((device) => hasPendingSave(device.save_state, device.persisted) && (device.persisted !== true || device.save_state !== 'save-error' || device.runtime_apply_status !== 'apply_failed'))) {
    return true;
  }
  if (state.rules.some((rule) => isUnpersistedNonIdle(rule.save_state, rule.persisted) || isDraftPending(rule.save_state))) {
    return true;
  }
  if (Object.values(state.mappings).some((mapping) => isUnpersistedNonIdle(mapping.save_state, mapping.persisted) || isDraftPending(mapping.save_state))) {
    return true;
  }
  if (isUnpersistedNonIdle(state.db.connector.save_state, state.db.connector.persisted) || isDraftPending(state.db.connector.save_state)) {
    return true;
  }
  return Object.values(state.db.targets).some((target) => isUnpersistedNonIdle(target.save_state, target.persisted) || isDraftPending(target.save_state));
}
function isSetupMutationAction(action: WorkbenchV2Action): boolean {
  switch (action.type) {
    case 'addDevice':
    case 'removeDevice':
    case 'updateDevice':
    case 'updateDeviceConfig':
    case 'renameDevice':
    case 'changeDeviceProtocol':
    case 'addRule':
    case 'removeRule':
    case 'updateRule':
    case 'renameRule':
    case 'toggleRuleEnabled':
    case 'updateRuleSkipped':
    case 'toggleRuleSkippedAddress':
    case 'toggleRuleShareEnabled':
    case 'updateRuleShareStart':
    case 'updateRuleShareStride':
    case 'updateMapping':
    case 'toggleMappingEnabled':
    case 'setAllMappingsEnabled':
    case 'bulkApplyTransform':
    case 'updateDbConnector':
    case 'upsertDbTarget':
    case 'updateDbTarget':
    case 'setAllDbTargetsEnabled':
    case 'autoAssignDbTargets':
    case 'setDbRowGroups':
      return true;
    default:
      return false;
  }
}
export function useStudioV2AutosaveState(
  enabled: boolean,
  locatorEntry?: StudioLocatorEntryRead,
  locatorEditorMounted = false,
  persistedMCWorkspaceScope?: PersistedMCWorkspaceScope,
  sourceWorkspaceId?: string,
) {
  const { t } = useTranslation('workbench-v2');
  const actions = useWorkbenchV2State();
  const devicesQuery = useStudioV2WorkspaceDevicesQuery(enabled);
  const rulesQuery = useStudioV2RulesQuery(enabled);
  const createDeviceMutation = useCreateStudioV2WorkspaceDeviceMutation();
  const updateDeviceMutation = useUpdateStudioV2WorkspaceDeviceMutation();
  const updateAvailabilityMutation = useUpdateStudioV2WorkspaceDeviceAvailabilityMutation();
  const deleteDeviceMutation = useDeleteStudioV2WorkspaceDeviceMutation();
  const stateRef = React.useRef(actions.state);
  const hydratedRef = React.useRef(false);
  const draftTrackingArmedRef = React.useRef(false);
  const [workspaceHydrated, setWorkspaceHydrated] = React.useState(false);
  const [draftLossWarning, setDraftLossWarning] = React.useState<string | null>(null);
  const saveMetaRef = React.useRef<Record<string, SaveMeta>>({});
  const appliedMCPersistedSignatureRef = React.useRef<string | null>(null);
  const locatorHydration = useStudioV2LocatorHydration(actions, stateRef, workspaceHydrated, locatorEntry, locatorEditorMounted);
  const sourceLease = useLocatorDeviceSourceLease(locatorEntry, sourceWorkspaceId, locatorHydration.ready);
  const rangeQueue = React.useRef<{ rule: (action: WorkbenchV2Action) => void; mapping: (previous: WorkbenchV2State, action: WorkbenchV2Action, next: WorkbenchV2State) => void }>({ rule: () => undefined, mapping: () => undefined });
  const sourceRanges = useRecordingRangeSubmission({ workspaceId: sourceWorkspaceId, enabled: locatorEntry ? locatorHydration.ready : enabled && workspaceHydrated,
    stateRef, apply: action => { stateRef.current = workbenchV2Reducer(stateRef.current, action); actions.dispatch(action); },
    queueRule: action => rangeQueue.current.rule(action), queueMapping: (previous, action, next) => rangeQueue.current.mapping(previous, action, next),
    prepareDevice: async id => {
      if (!sourceWorkspaceId) return;
      const device = await readSourceRangeDevice(sourceWorkspaceId, id, async () => {
        if (locatorEntry) { if (!sourceLease.allowed && !await sourceLease.acquire(sourceWorkspaceId)) return; return sourceLease.readback(); }
        const result = await devicesQuery.refetch(); return result.isError ? undefined : result.data;
      });
      if (!device) return;
      const readiness = await deviceAPI.checkReadiness(id);
      return readiness.device_id === id ? { device, activationAllowed: readiness.activation_allowed === true } : undefined;
    },
    diagnoseDevice: async id => {
      const result = await deviceAPI.testConnection(id);
      if (!result.success) return false;
      const readiness = await deviceAPI.checkReadiness(id);
      return readiness.device_id === id && readiness.activation_allowed === true;
    },
  });
  const persistedMCPoints = useStudioV2PersistedMCPoints(enabled, rulesQuery.data, devicesQuery.data, persistedMCWorkspaceScope, actions.state.rules);
  const persistedMCCanMutateRule = persistedMCPoints.canMutateRule;
  const canMutateRule = React.useCallback((ruleID: string): boolean => {
    if (!sourceRanges.canWrite(ruleID)) return false;
    if (sourceRanges.compatibleMC(ruleID)) return true;
    const localRule = stateRef.current.rules.find((rule) => rule.id === ruleID);
    if (localRule && isPersistedMCRule(localRule, stateRef.current.devices)) {
      return persistedMCCanMutateRule(ruleID, localRule.revision_id);
    }
    const localDevice = stateRef.current.devices.find((device) => device.id === localRule?.device_id);
    if (localRule?.persisted && localRule.mc_persisted_points && (!localDevice || localDevice.protocol === 'mc_3e')) return false;
    if (localRule) return true;
    return persistedMCCanMutateRule(ruleID);
  }, [persistedMCCanMutateRule, sourceRanges]);
  const ruleAutosave = useStudioV2RuleAutosave(actions, stateRef, canMutateRule, sourceRanges.ruleObservation, sourceWorkspaceId);
  const mappingAutosave = useStudioV2MappingAutosave(actions, stateRef, enabled, canMutateRule, sourceRanges.mappingObservation);
  rangeQueue.current = { rule: ruleAutosave.afterRuleAction, mapping: mappingAutosave.afterMappingAction };
  const databaseAutosave = useStudioV2DatabaseAutosave(actions, stateRef, enabled, canMutateRule);
  React.useEffect(() => {
    stateRef.current = actions.state;
  }, [actions.state]);
  React.useEffect(() => {
    if (
      !enabled ||
      !devicesQuery.isSuccess ||
      !rulesQuery.isSuccess ||
      !mappingAutosave.mappingsQuery.isSuccess ||
      !databaseAutosave.databaseConfigQuery.isSuccess ||
      !databaseAutosave.databaseTargetsQuery.isSuccess ||
      hydratedRef.current
    ) {
      return;
    }
    hydratedRef.current = true;
    const hydratedDevices = devicesQuery.data.map(hydrateStudioV2Device);
    const hydratedRules = withPersistedMCEvidence(rulesQuery.data
      .filter(canHydrateStudioV2Rule)
      .map(hydrateStudioV2Rule), persistedMCPoints.evidenceByRule);
    const deviceProtocolMap: Record<string, ProtocolType> = Object.fromEntries(hydratedDevices.map((d) => [d.id, d.protocol]));
    const enabledPoints = deriveAllPoints(hydratedRules, deviceProtocolMap).filter((point) => point.enabled && !point.skipped);
    const hydratedConnector = databaseAutosave.databaseConfigQuery.data
      ? hydrateStudioV2DatabaseConnector(databaseAutosave.databaseConfigQuery.data)
      : stateRef.current.db.connector;
    const hydratedRowGroups = databaseAutosave.databaseConfigQuery.data
      ? hydrateStudioV2DatabaseRowGroups(databaseAutosave.databaseConfigQuery.data.row_groups)
      : stateRef.current.db.row_groups ?? [];
    const baseState = workbenchV2Reducer(stateRef.current, {
      type: 'SET_STATE',
      payload: {
        devices: hydratedDevices,
        rules: hydratedRules,
        selectedRuleId: hydratedRules[0]?.id ?? null,
        points: [],
        mappings: {},
        db: {
          connector: hydratedConnector,
          row_groups: hydratedRowGroups,
          targets: {},
        },
      },
    });
    const hydratedMCState = mergePersistedMCState(
      { ...baseState, points: enabledPoints },
      persistedMCPoints.evidenceByRule,
      mappingAutosave.mappingsQuery.data,
      databaseAutosave.databaseTargetsQuery.data,
    );
    const hydratedProgress = inferHydratedProgress(hydratedDevices, hydratedMCState.rules, hydratedMCState.mappings);
    const hydratedPayload = {
      ...hydratedProgress,
      devices: hydratedDevices,
      rules: hydratedMCState.rules,
      selectedRuleId: hydratedRules[0]?.id ?? null,
      points: hydratedMCState.points,
      mappings: hydratedMCState.mappings,
      db: hydratedMCState.db,
    };
    const nextState = workbenchV2Reducer(baseState, {
      type: 'SET_STATE',
      payload: hydratedPayload,
    });
    appliedMCPersistedSignatureRef.current = persistedMCPoints.signature;
    stateRef.current = nextState;
    actions.dispatch({
      type: 'SET_STATE',
      payload: hydratedPayload,
    });
    try {
      if (typeof window !== 'undefined' && window.sessionStorage && window.sessionStorage.getItem(draftLossStorageKey) === '1') {
        setDraftLossWarning('偵測到上次重新整理前有未保存的本地草稿，系統已回復為最後一次成功保存的設定。');
        window.sessionStorage.removeItem(draftLossStorageKey);
      }
    } catch {
      setDraftLossWarning(null);
    }
    setWorkspaceHydrated(true);
  }, [
    actions,
    databaseAutosave.databaseConfigQuery.data,
    databaseAutosave.databaseConfigQuery.isSuccess,
    databaseAutosave.databaseTargetsQuery.data,
    databaseAutosave.databaseTargetsQuery.isSuccess,
    devicesQuery.data,
    devicesQuery.isSuccess,
    enabled,
    mappingAutosave.mappingsQuery.data,
    mappingAutosave.mappingsQuery.isSuccess,
    rulesQuery.data,
    rulesQuery.isSuccess,
    persistedMCPoints.evidenceByRule,
    persistedMCPoints.signature,
  ]);
  React.useEffect(() => {
    if (!workspaceHydrated || !hydratedRef.current || appliedMCPersistedSignatureRef.current === persistedMCPoints.signature) {
      return;
    }
    const merged = mergePersistedMCState(
      stateRef.current,
      persistedMCPoints.evidenceByRule,
      mappingAutosave.mappingsQuery.data,
      databaseAutosave.databaseTargetsQuery.data,
    );
    appliedMCPersistedSignatureRef.current = persistedMCPoints.signature;
    const payload = { rules: merged.rules, points: merged.points, mappings: merged.mappings, db: merged.db };
    stateRef.current = workbenchV2Reducer(stateRef.current, { type: 'SET_STATE', payload });
    actions.dispatch({ type: 'SET_STATE', payload });
  }, [
    actions,
    databaseAutosave.databaseTargetsQuery.data,
    mappingAutosave.mappingsQuery.data,
    persistedMCPoints.evidenceByRule,
    persistedMCPoints.signature,
    stateRef,
    workspaceHydrated,
  ]);
  React.useEffect(() => {
    if (!workspaceHydrated || !draftTrackingArmedRef.current) {
      return;
    }
    try {
      if (typeof window !== 'undefined' && window.sessionStorage) {
        if (hasLocalUnpersistedSetupDrafts(actions.state)) {
          window.sessionStorage.setItem(draftLossStorageKey, '1');
        } else {
          window.sessionStorage.removeItem(draftLossStorageKey);
        }
      }
    } catch {
      // 忽略 sessionStorage 例外
    }
  }, [actions.state, workspaceHydrated]);
  const applyDevicePatch = React.useCallback((deviceId: string, patch: Partial<Device>) => {
    const nextState = workbenchV2Reducer(stateRef.current, {
      type: 'updateDevice',
      deviceId,
      patch,
    });
    stateRef.current = nextState;
    actions.dispatch({
      type: 'updateDevice',
      deviceId,
      patch,
    });
  }, [actions]);
  const queueDeviceRef = React.useRef<(id: string) => void>(() => undefined);
  const deviceSubmission = useDeviceSourceSubmission({ workspaceId: sourceWorkspaceId, enabled: locatorEntry ? locatorHydration.ready : enabled && workspaceHydrated,
    canSubmit: () => !locatorEntry || sourceLease.allowed,
    stateRef, apply: action => { stateRef.current = workbenchV2Reducer(stateRef.current, action); actions.dispatch(action); },
    busy: id => Boolean(saveMetaRef.current[id]?.inFlight || saveMetaRef.current[id]?.pending), queue: id => queueDeviceRef.current(id),
    resume: id => { if (saveMetaRef.current[id]?.pending) { saveMetaRef.current[id].pending = false; queueDeviceRef.current(id); } },
    readback: async () => { if (locatorEntry) return sourceLease.readback(); const result = await devicesQuery.refetch(); return result.isError ? undefined : result.data; },
  });
  const flushDeviceSave = React.useCallback(async (deviceId: string) => {
    const meta = saveMetaFor(saveMetaRef.current, deviceId);
    const currentDevice = stateRef.current.devices.find((device) => device.id === deviceId);
    if (!currentDevice) {
      delete saveMetaRef.current[deviceId];
      return;
    }
    if (!isStudioV2DeviceValid(currentDevice)) {
      applyDevicePatch(deviceId, {
        save_state: 'draft-invalid',
        save_error: null,
      });
      if (currentDevice.persisted && currentDevice.running && currentDevice.availability_status !== 'unavailable') {
        void updateAvailabilityMutation.mutateAsync({
          deviceId,
          request: {
            availability_status: 'unavailable',
            availability_reason: invalidDeviceAvailabilityReason,
          },
        }).then((savedDevice) => {
          applyDevicePatch(deviceId, {
            availability_status: savedDevice.availability_status ?? 'unavailable',
            availability_reason: savedDevice.availability_reason ?? invalidDeviceAvailabilityReason,
            running: savedDevice.running ?? false,
          });
        }).catch((error) => {
          applyDevicePatch(deviceId, {
            save_error: errorMessageOf(error, t('errors.autosave_failed')),
          });
        });
      }
      meta.inFlight = false;
      meta.pending = false;
      return;
    }
    meta.inFlight = true;
    applyDevicePatch(deviceId, {
      save_state: 'saving',
      save_error: null,
    });
    try {
      const savedDevice = currentDevice.persisted
        ? await updateDeviceMutation.mutateAsync({
          deviceId,
          request: toStudioV2DeviceUpdateRequest(currentDevice),
        })
        : await createDeviceMutation.mutateAsync(toStudioV2DeviceCreateRequest(currentDevice));
      const observed = await deviceSubmission.acknowledge(deviceId, savedDevice);
      if (!deviceSubmission.canPublish(deviceId)) return;
      if (observed && savedDevice.id !== deviceId) {
        applyDevicePatch(deviceId, { save_state: 'save-error', save_error: t('errors.autosave_failed') });
        return;
      }
      const latestDevice = stateRef.current.devices.find((device) => device.id === deviceId);
      if (latestDevice) {
        const runtimeApplyFailed = savedDevice.runtime_apply_status === 'apply_failed';
        applyDevicePatch(deviceId, {
          name: latestDevice.name,
          description: latestDevice.description,
          config: latestDevice.config,
          protocol: latestDevice.protocol,
          persisted: true,
          status: savedDevice.status === 'active' ? 'active' : latestDevice.status,
          save_state: runtimeApplyFailed || observed?.status === 'unknown' ? 'save-error' : 'saved',
          save_error: runtimeApplyFailed ? savedDevice.runtime_apply_message ?? '執行中設定套用失敗' : null,
          runtime_apply_status: savedDevice.runtime_apply_status,
          runtime_apply_message: savedDevice.runtime_apply_message ?? null,
          availability_status: savedDevice.availability_status ?? 'available',
          availability_reason: savedDevice.availability_reason ?? null,
          running: savedDevice.running ?? latestDevice.running ?? false,
        });
      }
    } catch (error) {
      deviceSubmission.failed(deviceId);
      if (!deviceSubmission.canPublish(deviceId)) return;
      applyDevicePatch(deviceId, {
        save_state: 'save-error',
        save_error: errorMessageOf(error, t('errors.autosave_failed')),
      });
    } finally {
      meta.inFlight = false;
      if (meta.pending && deviceSubmission.canWrite(deviceId)) {
        meta.pending = false;
        void flushDeviceSave(deviceId);
      }
    }
  }, [applyDevicePatch, createDeviceMutation, deviceSubmission, t, updateAvailabilityMutation, updateDeviceMutation]);
  const queueDeviceSave = React.useCallback((deviceId: string) => {
    const meta = saveMetaFor(saveMetaRef.current, deviceId);
    if (!deviceSubmission.canWrite(deviceId)) { meta.pending = true; return; }
    if (meta.inFlight) {
      meta.pending = true;
      return;
    }
    void flushDeviceSave(deviceId);
  }, [deviceSubmission, flushDeviceSave]);
  queueDeviceRef.current = queueDeviceSave;
  const dispatch = React.useCallback((action: WorkbenchV2Action) => {
    const previousState = stateRef.current;
    if (action.type === 'removeDevice') {
      if (deviceSubmission.hasUnresolved(action.deviceId)) return;
      const targetDevice = previousState.devices.find((device) => device.id === action.deviceId);
      if (!targetDevice) {
        return;
      }
      if (!targetDevice.persisted) {
        const nextState = workbenchV2Reducer(previousState, action);
        stateRef.current = nextState;
        actions.dispatch(action);
        return;
      }
      applyDevicePatch(action.deviceId, {
        save_state: 'saving',
        save_error: null,
      });
      void deleteDeviceMutation.mutateAsync(action.deviceId)
        .then(() => {
          const nextState = workbenchV2Reducer(stateRef.current, action);
          stateRef.current = nextState;
          actions.dispatch(action);
          delete saveMetaRef.current[action.deviceId];
        })
        .catch((error) => {
          applyDevicePatch(action.deviceId, {
            save_state: 'save-error',
            save_error: errorMessageOf(error, t('errors.autosave_failed')),
          });
        });
      return;
    }
    if (ruleAutosave.interceptRuleAction(previousState, action)) {
      return;
    }
    const nextState = workbenchV2Reducer(previousState, action);
    stateRef.current = nextState;
    actions.dispatch(action);
    if (workspaceHydrated && isSetupMutationAction(action)) {
      draftTrackingArmedRef.current = true;
    }
    switch (action.type) {
      case 'addDevice':
      case 'updateDevice':
      case 'updateDeviceConfig':
      case 'renameDevice':
      case 'changeDeviceProtocol':
        queueDeviceSave(action.type === 'addDevice' ? action.device.id : action.deviceId);
        break;
      default:
        break;
    }
    ruleAutosave.afterRuleAction(action);
    mappingAutosave.afterMappingAction(previousState, action, nextState);
    databaseAutosave.afterDatabaseAction(action, nextState);
  }, [actions, applyDevicePatch, databaseAutosave, deleteDeviceMutation, deviceSubmission, mappingAutosave, queueDeviceSave, ruleAutosave, t, workspaceHydrated]);
  return {
    actions: {
      ...actions,
      dispatch,
    },
    deviceSource: locatorEntry ? { ...deviceSubmission.port, range: sourceRanges.port,
      prepare: async (id: string) => id === sourceWorkspaceId && (sourceLease.allowed ? Boolean(await sourceLease.readback()) : await sourceLease.acquire(id)),
      submit: async (request: DeviceSourceRequest) => sourceLease.allowed && await sourceLease.readback() ? deviceSubmission.port.submit(request) : { request, status: 'blocked' as const },
    } : { ...deviceSubmission.port, range: sourceRanges.port },
    databaseConfigQuery: databaseAutosave.databaseConfigQuery,
    databaseTargetsQuery: databaseAutosave.databaseTargetsQuery,
    devicesQuery,
    draftLossWarning,
    mappingsQuery: mappingAutosave.mappingsQuery,
    rulesQuery,
    state: actions.state,
    workspaceHydrated,
    autosaveBarrier: getStudioV2AutosaveBarrier(actions.state), locatorHydration,
  };
}
