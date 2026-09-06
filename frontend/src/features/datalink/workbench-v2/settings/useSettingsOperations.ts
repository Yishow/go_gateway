import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type {
  GeneralSettings,
  ModbusShareSettings,
  SchedulerSettings,
  SettingsConnector,
  Settings,
  TimeseriesSettings,
  WorkbenchV2State,
} from '../state/types';
import type { WorkbenchV2Action } from '../state/useWorkbenchV2State';
import {
  useCreateDbTargetConnectorMutation,
  useDeleteDbTargetConnectorMutation,
  useTestDbTargetConnectorMutation,
  useUpdateDbTargetConnectorMutation,
  useUpdateSettingKeyMutation,
} from '../../../../hooks/datalink/useSettings';
import { makeDefaultConnector } from '../state/settingsDefaults';
import { getSafeErrorMessage, type SafeErrorMessage } from '../../../../utils/typedErrors';
import {
  buildCreateConnectorRequest,
  buildPersistableSettingEntries,
  buildUpdateConnectorRequest,
  mapDatabaseConnectorToSettingsConnector,
} from './backendMappings';
import { createSettingsOperationOwnership } from './settingsOperationOwnership';

type SettingsRetryOperation =
  | { kind: 'save' }
  | { kind: 'add-connector' }
  | { kind: 'update-connector'; id: string; patch: Partial<SettingsConnector> }
  | { kind: 'delete-connector'; id: string }
  | { kind: 'test-connector'; id: string };

export interface SettingsOperations {
  operationError: SafeErrorMessage | null;
  isSaving: boolean;
  pending: boolean;
  onUpdateGeneral: (patch: Partial<GeneralSettings>) => void;
  onUpdateTimeseries: (patch: Partial<TimeseriesSettings>) => void;
  onUpdateScheduler: (patch: Partial<SchedulerSettings>) => void;
  onUpdateModbusShare: (patch: Partial<ModbusShareSettings>) => void;
  onAddConnector: () => Promise<void>;
  onUpdateConnector: (id: string, patch: Partial<SettingsConnector>) => void;
  onRemoveConnector: (id: string) => void;
  onTestConnector: (id: string) => void;
  onReset: () => void;
  onSave: () => Promise<void>;
  onRetry: () => void;
}

export function useSettingsOperations(
  state: WorkbenchV2State,
  dispatch: React.Dispatch<WorkbenchV2Action>,
): SettingsOperations {
  const { t } = useTranslation('workbench-v2');
  const updateSettingKeyMutation = useUpdateSettingKeyMutation();
  const createConnectorMutation = useCreateDbTargetConnectorMutation();
  const updateConnectorMutation = useUpdateDbTargetConnectorMutation();
  const deleteConnectorMutation = useDeleteDbTargetConnectorMutation();
  const testConnectorMutation = useTestDbTargetConnectorMutation();
  const [operationError, setOperationError] = React.useState<SafeErrorMessage | null>(null);
  const [retryOperation, setRetryOperation] = React.useState<SettingsRetryOperation | null>(null);
  const [isSaving, setIsSaving] = React.useState(false);
  const ownership = React.useRef(createSettingsOperationOwnership()).current;
  const activeSaveRequests = React.useRef(0);
  const saveQueue = React.useRef<{ snapshot: Settings; waiters: Array<() => void> } | null>(null);
  const saveRunning = React.useRef(false);
  const addQueue = React.useRef(Promise.resolve());
  const pendingAdds = React.useRef(0);
  const successfulAdds = React.useRef<SettingsConnector[]>([]);
  const deletedConnectorIds = React.useRef(new Set<string>());
  const pendingConnectorDrafts = React.useRef(new Map<string, SettingsConnector>());
  const latestState = React.useRef(state);
  latestState.current = state;
  const safeErrorMessage = React.useCallback((error: unknown): SafeErrorMessage => (
    getSafeErrorMessage(error, (key, options) => t(key, options))
  ), [t]);

  const clearOperationError = React.useCallback(() => {
    setOperationError(null);
    setRetryOperation(null);
  }, []);
  const onUpdateGeneral = React.useCallback((patch: Partial<GeneralSettings>) => {
    dispatch({ type: 'updateSettingsSection', section: 'general', patch });
  }, [dispatch]);
  const onUpdateTimeseries = React.useCallback((patch: Partial<TimeseriesSettings>) => {
    dispatch({ type: 'updateSettingsSection', section: 'timeseries', patch });
  }, [dispatch]);
  const onUpdateScheduler = React.useCallback((patch: Partial<SchedulerSettings>) => {
    dispatch({ type: 'updateSettingsSection', section: 'scheduler', patch });
  }, [dispatch]);
  // 編輯只更新資料與清除既有錯誤；save_state 一律留給實際發出的存檔請求設定。
  const onUpdateModbusShare = React.useCallback((patch: Partial<ModbusShareSettings>) => {
    dispatch({ type: 'updateSettingsSection', section: 'modbus_share', patch: { ...patch, save_error: null } });
  }, [dispatch]);
  const onReset = React.useCallback(() => {
    dispatch({ type: 'resetSettingsToDefaults' });
  }, [dispatch]);

  const persistSettingsSnapshot = React.useCallback(async (snapshot: Settings) => {
    try {
      for (const entry of buildPersistableSettingEntries(snapshot)) {
        const saved = await updateSettingKeyMutation.mutateAsync(entry);
        if (entry.key === 'modbus_share' && typeof saved.value === 'object' && saved.value !== null) {
          const value = saved.value as Record<string, unknown>;
          const revision = value.new_settings_revision ?? value.settings_revision;
          // settings_revision 是伺服器指派的 CAS 憑證，不是使用者可編輯的欄位：
          // 一律採用最新回應。漏掉它會讓後續每次存檔都撞上 revision 衝突，
          // 且設定頁只在初次載入 hydrate，除了重新整理沒有恢復路徑。
          if (typeof revision === 'string') {
            dispatch({ type: 'updateSettingsSection', section: 'modbus_share', patch: { settings_revision: revision, expected_settings_revision: revision } });
          }
        }
      }
      // 終態一律寫入：不論期間使用者是否編輯其他欄位，存檔都不得停在 saving。
      // 但佇列中還有後續存檔時不得提前宣告 saved，否則 autosave barrier 會把
      // 尚未寫入的設定判成已完成，讓啟用越過它。
      if (!saveQueue.current) {
        dispatch({ type: 'updateSettingsSection', section: 'modbus_share', patch: { save_state: 'saved', save_error: null } });
      }
    } catch (error) {
      const safeError = safeErrorMessage(error);
      setOperationError(safeError);
      setRetryOperation({ kind: 'save' });
      dispatch({ type: 'updateSettingsSection', section: 'modbus_share', patch: { save_state: 'save-error', save_error: safeError.message } });
    }
  }, [dispatch, safeErrorMessage, updateSettingKeyMutation]);

  const drainSaveQueue = React.useCallback(async () => {
    if (saveRunning.current) return;
    saveRunning.current = true;
    try {
      while (saveQueue.current) {
        const item = saveQueue.current;
        saveQueue.current = null;
        // 每一批開始前重新宣告 saving，讓 barrier 在整段佇列期間都看得到進行中。
        dispatch({ type: 'updateSettingsSection', section: 'modbus_share', patch: { save_state: 'saving' } });
        await persistSettingsSnapshot(item.snapshot);
        activeSaveRequests.current = Math.max(0, activeSaveRequests.current - item.waiters.length);
        setIsSaving(activeSaveRequests.current > 0);
        item.waiters.forEach((resolve) => resolve());
      }
    } finally {
      saveRunning.current = false;
      if (saveQueue.current) void drainSaveQueue();
    }
  }, [dispatch, persistSettingsSnapshot]);

  const onSave = React.useCallback(() => {
    activeSaveRequests.current += 1;
    setIsSaving(true);
    clearOperationError();
    dispatch({ type: 'updateSettingsSection', section: 'modbus_share', patch: { save_state: 'saving', save_error: null } });
    return new Promise<void>((resolve) => {
      if (saveQueue.current) {
        saveQueue.current.snapshot = state.settings;
        saveQueue.current.waiters.push(resolve);
      } else {
        saveQueue.current = { snapshot: state.settings, waiters: [resolve] };
      }
      void drainSaveQueue();
    });
  }, [clearOperationError, dispatch, drainSaveQueue, state.settings]);

  const onAddConnector = React.useCallback(() => {
    const scope = 'connector-add';
    const token = ownership.begin(scope);
    const draft = makeDefaultConnector(latestState.current.settings.connectors.length + pendingAdds.current + 1);
    pendingAdds.current += 1;
    clearOperationError();
    const run = addQueue.current.then(async () => {
      try {
        const created = await createConnectorMutation.mutateAsync(buildCreateConnectorRequest(draft));
        const row = mapDatabaseConnectorToSettingsConnector(created, draft.password ?? '');
        // 後端重新配發同一個 id 時，該 id 重新成為存活項目。
        deletedConnectorIds.current.delete(row.id);
        const existingIds = new Set(latestState.current.settings.connectors.map((connector) => connector.id));
        // 已寫回列表或已被刪除的項目都算結案，不再參與後續的樂觀合併。
        const isSettled = (added: SettingsConnector) => existingIds.has(added.id) || deletedConnectorIds.current.has(added.id);
        successfulAdds.current = successfulAdds.current.filter((added) => !isSettled(added));
        successfulAdds.current.push(row);
        dispatch({
          type: 'updateSettings',
          patch: { connectors: [...latestState.current.settings.connectors, ...successfulAdds.current.filter((added) => !isSettled(added))] },
        });
      } catch (error) {
        if (ownership.isCurrent(token, scope)) { setOperationError(safeErrorMessage(error)); setRetryOperation({ kind: 'add-connector' }); }
      } finally {
        pendingAdds.current = Math.max(0, pendingAdds.current - 1);
      }
    });
    addQueue.current = run.catch(() => undefined);
    return run;
  }, [clearOperationError, createConnectorMutation, dispatch, ownership, safeErrorMessage]);

  const onUpdateConnector = React.useCallback((id: string, patch: Partial<SettingsConnector>) => {
    const scope = `connector:${id}`;
    const token = ownership.begin(scope);
    dispatch({ type: 'updateConnector', id, patch });
    const current = pendingConnectorDrafts.current.get(id)
      ?? latestState.current.settings.connectors.find((connector) => connector.id === id);
    if (!current) return;
    const next = { ...current, ...patch };
    pendingConnectorDrafts.current.set(id, next);
    clearOperationError();
    void updateConnectorMutation.mutateAsync({ id, data: buildUpdateConnectorRequest(next) }).then((saved) => {
      if (ownership.isCurrent(token, scope)) {
        pendingConnectorDrafts.current.delete(id);
        dispatch({ type: 'updateConnector', id, patch: mapDatabaseConnectorToSettingsConnector(saved, next.password ?? '') });
      }
    }).catch((error) => {
      if (ownership.isCurrent(token, scope)) { setOperationError(safeErrorMessage(error)); setRetryOperation({ kind: 'update-connector', id, patch }); }
    });
  }, [clearOperationError, dispatch, ownership, safeErrorMessage, updateConnectorMutation]);

  const onRemoveConnector = React.useCallback((id: string) => {
    const scope = `connector:${id}`;
    const token = ownership.begin(scope);
    pendingConnectorDrafts.current.delete(id);
    clearOperationError();
    void deleteConnectorMutation.mutateAsync(id).then(() => {
      if (ownership.isCurrent(token, scope)) {
        deletedConnectorIds.current.add(id);
        successfulAdds.current = successfulAdds.current.filter((added) => added.id !== id);
        dispatch({ type: 'removeConnector', id });
      }
    }).catch((error) => {
      if (ownership.isCurrent(token, scope)) { setOperationError(safeErrorMessage(error)); setRetryOperation({ kind: 'delete-connector', id }); }
    });
  }, [clearOperationError, deleteConnectorMutation, dispatch, ownership, safeErrorMessage]);

  const onTestConnector = React.useCallback((id: string) => {
    const scope = `connector:${id}`;
    const token = ownership.begin(scope);
    pendingConnectorDrafts.current.delete(id);
    const current = latestState.current.settings.connectors.find((connector) => connector.id === id);
    dispatch({ type: 'startConnectorTest', id });
    clearOperationError();
    void testConnectorMutation.mutateAsync(id).then((tested) => {
      if (!ownership.isCurrent(token, scope)) return;
      const connector = mapDatabaseConnectorToSettingsConnector(tested, current?.password ?? '');
      dispatch({ type: 'completeConnectorTest', id, result: { status: connector.status === 'testing' ? 'unknown' : connector.status, last_check_at: connector.last_check_at ?? new Date().toISOString(), last_check_error: connector.last_check_error } });
      dispatch({ type: 'updateConnector', id, patch: connector });
    }).catch((error) => {
      if (!ownership.isCurrent(token, scope)) return;
      const safeError = safeErrorMessage(error);
      dispatch({ type: 'completeConnectorTest', id, result: { status: 'unreachable', last_check_at: new Date().toISOString(), last_check_error: safeError.message } });
      setOperationError(safeError);
      setRetryOperation({ kind: 'test-connector', id });
    });
  }, [clearOperationError, dispatch, ownership, safeErrorMessage, testConnectorMutation]);

  const onRetry = React.useCallback(() => {
    if (!retryOperation) return;
    if (retryOperation.kind === 'save') void onSave();
    if (retryOperation.kind === 'add-connector') void onAddConnector();
    if (retryOperation.kind === 'update-connector') onUpdateConnector(retryOperation.id, retryOperation.patch);
    if (retryOperation.kind === 'delete-connector') onRemoveConnector(retryOperation.id);
    if (retryOperation.kind === 'test-connector') onTestConnector(retryOperation.id);
  }, [onAddConnector, onRemoveConnector, onSave, onTestConnector, onUpdateConnector, retryOperation]);

  return { operationError, isSaving, pending: isSaving || updateSettingKeyMutation.isPending, onUpdateGeneral, onUpdateTimeseries, onUpdateScheduler, onUpdateModbusShare, onAddConnector, onUpdateConnector, onRemoveConnector, onTestConnector, onReset, onSave, onRetry };
}
