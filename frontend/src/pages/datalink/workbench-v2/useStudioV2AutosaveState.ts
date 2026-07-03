import * as React from 'react';
import { useCreateStudioV2WorkspaceDeviceMutation, useDeleteStudioV2WorkspaceDeviceMutation, useStudioV2WorkspaceDevicesQuery, useUpdateStudioV2WorkspaceDeviceAvailabilityMutation, useUpdateStudioV2WorkspaceDeviceMutation } from '../../../hooks/datalink/useStudioV2WorkspaceDevices';
import { hydrateStudioV2Device, isStudioV2DeviceValid, toStudioV2DeviceCreateRequest, toStudioV2DeviceUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { canHydrateStudioV2Rule, hydrateStudioV2Rule } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { hydrateStudioV2Mapping } from '../../../features/datalink/workbench-v2/state/studioV2MappingAutosave';
import type { DbTarget, Device, Mapping, Point, Rule, WorkbenchV2State } from '../../../features/datalink/workbench-v2/state/types';
import { useWorkbenchV2State, workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { useStudioV2RulesQuery } from '../../../hooks/datalink/useStudioV2Rules';
import { useStudioV2RuleAutosave } from './useStudioV2RuleAutosave';
import { useStudioV2MappingAutosave } from './useStudioV2MappingAutosave';
import { useStudioV2DatabaseAutosave } from './useStudioV2DatabaseAutosave';
import { hydrateStudioV2DatabaseConnector, hydrateStudioV2DatabaseRowGroups, hydrateStudioV2DatabaseTarget } from '../../../features/datalink/workbench-v2/state/studioV2DatabaseAutosave';
import { deriveAllPoints } from '../../../features/datalink/workbench-v2/state/sourceRule';
import type { StudioV2WorkspaceMappingRecord, StudioV2WorkspaceDatabaseTargetRecord } from '../../../types/datalink';

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

function errorMessageOf(error: unknown): string {
  return error instanceof Error ? error.message : '儲存失敗';
}

const invalidDeviceAvailabilityReason = 'device form is invalid';

function mappingRowKey(ruleId: string, address: string): string {
  return `${ruleId}::${address}`;
}

function buildPersistedPointAliases(mappings: Record<string, Mapping>): Map<string, string> {
  const aliases = new Map<string, string>();
  Object.values(mappings).forEach((mapping) => {
    aliases.set(mapping.point_id, mapping.point_id);
    if (mapping.persisted_point_id) {
      aliases.set(mapping.persisted_point_id, mapping.point_id);
    }
  });
  return aliases;
}

function buildBootstrapMappings(
  baseState: WorkbenchV2State,
  points: Point[],
  records: StudioV2WorkspaceMappingRecord[] | undefined,
): { mappings: Record<string, Mapping>; persistedPointAliases: Map<string, string> } {
  const recordByKey = new Map((records ?? []).map((record) => [mappingRowKey(record.rule_id, record.address), record]));
  const nextState = workbenchV2Reducer(baseState, {
    type: 'initMappingsForPoints',
    points,
  });

  const hydratedMappings = { ...nextState.mappings };
  points.forEach((point) => {
    const record = recordByKey.get(mappingRowKey(point.rule_id, point.address));
    const current = hydratedMappings[point.id];
    if (!record || !current) {
      return;
    }
    hydratedMappings[point.id] = hydrateStudioV2Mapping(point, record, current);
  });

  return {
    mappings: hydratedMappings,
    persistedPointAliases: buildPersistedPointAliases(hydratedMappings),
  };
}

function buildBootstrapTargets(
  persistedPointAliases: Map<string, string>,
  records: StudioV2WorkspaceDatabaseTargetRecord[] | undefined,
): Record<string, DbTarget> {
  return (records ?? []).reduce<Record<string, DbTarget>>((accumulator, record) => {
    const currentPointId = persistedPointAliases.get(record.point_id);
    if (!currentPointId) {
      return accumulator;
    }
    accumulator[currentPointId] = hydrateStudioV2DatabaseTarget(
      { ...record, point_id: currentPointId },
      accumulator[currentPointId],
    );
    return accumulator;
  }, {});
}

function inferHydratedProgress(
  devices: Device[],
  rules: Rule[],
  mappings: Record<string, Mapping>,
): Pick<WorkbenchV2State, 'current' | 'completed'> {
  const completed = new Set<number>();
  if (devices.length > 0) completed.add(1);
  if (rules.length > 0) completed.add(2);
  if (Object.values(mappings).some((mapping) => mapping.persisted && mapping.enabled && mapping.tag_id)) completed.add(3);
  const current: WorkbenchV2State['current'] = completed.has(3) ? 4 : completed.has(2) ? 3 : completed.has(1) ? 2 : 1;
  return { current, completed };
}

function isDraftPending(saveState?: string): boolean {
  return saveState === 'saving' || saveState === 'save-error' || saveState === 'draft-invalid';
}

function isUnpersistedNonIdle(saveState?: string, persisted?: boolean): boolean {
  return !persisted && saveState !== 'idle';
}

function hasLocalUnpersistedSetupDrafts(state: WorkbenchV2State): boolean {
  if (state.devices.some((device) => isUnpersistedNonIdle(device.save_state, device.persisted) || device.save_state === 'saving' || device.save_state === 'draft-invalid' || (device.save_state === 'save-error' && device.persisted !== true && device.runtime_apply_status !== 'apply_failed'))) {
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

export function useStudioV2AutosaveState(enabled: boolean) {
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
  const ruleAutosave = useStudioV2RuleAutosave(actions, stateRef);
  const mappingAutosave = useStudioV2MappingAutosave(actions, stateRef, enabled);
  const databaseAutosave = useStudioV2DatabaseAutosave(actions, stateRef, enabled);

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
    const hydratedRules = rulesQuery.data
      .filter(canHydrateStudioV2Rule)
      .map(hydrateStudioV2Rule);
    const fallbackDeviceId = hydratedDevices[0]?.id || 'dev-01';
    const enabledPoints = deriveAllPoints(hydratedRules, fallbackDeviceId).filter((point) => point.enabled && !point.skipped);
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

    const { mappings: hydratedMappings, persistedPointAliases } = buildBootstrapMappings(baseState, enabledPoints, mappingAutosave.mappingsQuery.data);
    const hydratedTargets = buildBootstrapTargets(persistedPointAliases, databaseAutosave.databaseTargetsQuery.data);
    const hydratedProgress = inferHydratedProgress(hydratedDevices, hydratedRules, hydratedMappings);
    const nextState = workbenchV2Reducer(baseState, {
      type: 'SET_STATE',
      payload: {
        ...hydratedProgress,
        points: enabledPoints,
        mappings: hydratedMappings,
        db: {
          connector: hydratedConnector,
          row_groups: hydratedRowGroups,
          targets: hydratedTargets,
        },
      },
    });

    stateRef.current = nextState;
    actions.dispatch({
      type: 'SET_STATE',
      payload: {
        ...hydratedProgress,
        devices: hydratedDevices,
        rules: hydratedRules,
        selectedRuleId: hydratedRules[0]?.id ?? null,
        points: enabledPoints,
        mappings: hydratedMappings,
        db: {
          connector: hydratedConnector,
          row_groups: hydratedRowGroups,
          targets: hydratedTargets,
        },
      },
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
            save_error: errorMessageOf(error),
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
          save_state: runtimeApplyFailed ? 'save-error' : 'saved',
          save_error: runtimeApplyFailed ? savedDevice.runtime_apply_message ?? '執行中設定套用失敗' : null,
          runtime_apply_status: savedDevice.runtime_apply_status,
          runtime_apply_message: savedDevice.runtime_apply_message ?? null,
          availability_status: savedDevice.availability_status ?? 'available',
          availability_reason: savedDevice.availability_reason ?? null,
          running: savedDevice.running ?? latestDevice.running ?? false,
        });
      }
    } catch (error) {
      applyDevicePatch(deviceId, {
        save_state: 'save-error',
        save_error: errorMessageOf(error),
      });
    } finally {
      meta.inFlight = false;
      if (meta.pending) {
        meta.pending = false;
        void flushDeviceSave(deviceId);
      }
    }
  }, [applyDevicePatch, createDeviceMutation, updateAvailabilityMutation, updateDeviceMutation]);

  const queueDeviceSave = React.useCallback((deviceId: string) => {
    const meta = saveMetaFor(saveMetaRef.current, deviceId);
    if (meta.inFlight) {
      meta.pending = true;
      return;
    }

    void flushDeviceSave(deviceId);
  }, [flushDeviceSave]);

  const dispatch = React.useCallback((action: WorkbenchV2Action) => {
    const previousState = stateRef.current;

    if (action.type === 'removeDevice') {
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
            save_error: errorMessageOf(error),
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
  }, [actions, applyDevicePatch, databaseAutosave, deleteDeviceMutation, mappingAutosave, queueDeviceSave, ruleAutosave, workspaceHydrated]);

  return {
    actions: {
      ...actions,
      dispatch,
    },
    databaseConfigQuery: databaseAutosave.databaseConfigQuery,
    databaseTargetsQuery: databaseAutosave.databaseTargetsQuery,
    devicesQuery,
    draftLossWarning,
    mappingsQuery: mappingAutosave.mappingsQuery,
    rulesQuery,
    state: actions.state,
    workspaceHydrated,
  };
}
