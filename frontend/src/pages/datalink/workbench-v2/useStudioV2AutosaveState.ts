import * as React from 'react';
import { useCreateStudioV2WorkspaceDeviceMutation, useDeleteStudioV2WorkspaceDeviceMutation, useStudioV2WorkspaceDevicesQuery, useUpdateStudioV2WorkspaceDeviceAvailabilityMutation, useUpdateStudioV2WorkspaceDeviceMutation } from '../../../hooks/datalink/useStudioV2WorkspaceDevices';
import { hydrateStudioV2Device, isStudioV2DeviceValid, toStudioV2DeviceCreateRequest, toStudioV2DeviceUpdateRequest } from '../../../features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { hydrateStudioV2Rule, rebindDraftRulesToWorkspaceDevice } from '../../../features/datalink/workbench-v2/state/studioV2RuleAutosave';
import type { Device } from '../../../features/datalink/workbench-v2/state/types';
import { useWorkbenchV2State, workbenchV2Reducer, type WorkbenchV2Action } from '../../../features/datalink/workbench-v2/state/useWorkbenchV2State';
import { useStudioV2RulesQuery } from '../../../hooks/datalink/useStudioV2Rules';
import { useStudioV2RuleAutosave } from './useStudioV2RuleAutosave';
import { useStudioV2MappingAutosave } from './useStudioV2MappingAutosave';
import { useStudioV2DatabaseAutosave } from './useStudioV2DatabaseAutosave';

type SaveMeta = {
  inFlight: boolean;
  pending: boolean;
};

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
  const [workspaceHydrated, setWorkspaceHydrated] = React.useState(false);
  const saveMetaRef = React.useRef<Record<string, SaveMeta>>({});
  const ruleAutosave = useStudioV2RuleAutosave(actions, stateRef);
  const mappingAutosave = useStudioV2MappingAutosave(actions, stateRef, enabled);
  const databaseAutosave = useStudioV2DatabaseAutosave(actions, stateRef, enabled);

  React.useEffect(() => {
    stateRef.current = actions.state;
  }, [actions.state]);

  React.useEffect(() => {
    if (!enabled || !devicesQuery.isSuccess || !rulesQuery.isSuccess || hydratedRef.current) {
      return;
    }

    hydratedRef.current = true;
    const hydratedDevices = devicesQuery.data.map(hydrateStudioV2Device);
    const nextDevices = hydratedDevices.length > 0 ? hydratedDevices : stateRef.current.devices;
    const hydratedRules = rulesQuery.data.length > 0
      ? rulesQuery.data.map(hydrateStudioV2Rule)
      : nextDevices.length > 0
        ? rebindDraftRulesToWorkspaceDevice(stateRef.current.rules, nextDevices[0].id)
        : [];

    actions.dispatch({
      type: 'SET_STATE',
      payload: {
        devices: nextDevices,
        rules: hydratedRules,
        selectedRuleId: hydratedRules[0]?.id ?? null,
        points: [],
        mappings: {},
        db: stateRef.current.db,
      },
    });
    setWorkspaceHydrated(true);
  }, [actions, devicesQuery.data, devicesQuery.isSuccess, enabled, rulesQuery.data, rulesQuery.isSuccess]);

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
  }, [actions, applyDevicePatch, databaseAutosave, deleteDeviceMutation, mappingAutosave, queueDeviceSave, ruleAutosave]);

  return {
    actions: {
      ...actions,
      dispatch,
    },
    databaseConfigQuery: databaseAutosave.databaseConfigQuery,
    databaseTargetsQuery: databaseAutosave.databaseTargetsQuery,
    devicesQuery,
    mappingsQuery: mappingAutosave.mappingsQuery,
    rulesQuery,
    state: actions.state,
    workspaceHydrated,
  };
}
