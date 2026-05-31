import * as React from 'react';
import type {
  WorkbenchV2State,
  GeneralSettings,
  TimeseriesSettings,
  SchedulerSettings,
  ModbusShareSettings,
  SettingsConnector,
} from '../state/types';
import type { WorkbenchV2Action } from '../state/useWorkbenchV2State';
import {
  useCreateDbTargetConnectorMutation,
  useDbTargetConnectorsQuery,
  useDeleteDbTargetConnectorMutation,
  useSettingsItemsQuery,
  useTestDbTargetConnectorMutation,
  useUpdateDbTargetConnectorMutation,
  useUpdateSettingKeyMutation,
} from '../../../../hooks/datalink/useSettings';
import { SettingsHeader } from './SettingsHeader';
import { ConnectorPoolSection } from './ConnectorPoolSection';
import { TimeseriesSection } from './TimeseriesSection';
import { SchedulerSection } from './SchedulerSection';
import { ModbusShareSection } from './ModbusShareSection';
import { UiSection } from './UiSection';
import { ApiSection } from './ApiSection';
import { DiagnosticsSection } from './DiagnosticsSection';
import { SaveBar } from './SaveBar';
import { makeDefaultConnector } from '../state/settingsDefaults';
import {
  buildCreateConnectorRequest,
  buildPersistableSettingEntries,
  buildUpdateConnectorRequest,
  mapDatabaseConnectorToSettingsConnector,
  mapSettingItemsToWorkbenchSettings,
} from './backendMappings';

/**
 * SettingsPage 元件屬性
 */
interface SettingsPageProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

/**
 * 系統設定頁面容器元件
 * 落地需求：「Settings page layout」與任務 3.1 規格
 */
export function SettingsPage({ state, dispatch }: SettingsPageProps) {
  const settingsItemsQuery = useSettingsItemsQuery();
  const connectorsQuery = useDbTargetConnectorsQuery();
  const updateSettingKeyMutation = useUpdateSettingKeyMutation();
  const createConnectorMutation = useCreateDbTargetConnectorMutation();
  const updateConnectorMutation = useUpdateDbTargetConnectorMutation();
  const deleteConnectorMutation = useDeleteDbTargetConnectorMutation();
  const testConnectorMutation = useTestDbTargetConnectorMutation();

  const [hydrated, setHydrated] = React.useState(false);
  const [operationError, setOperationError] = React.useState<string | null>(null);
  const [isSaving, setIsSaving] = React.useState(false);

  React.useEffect(() => {
    if (hydrated || !settingsItemsQuery.isSuccess || !connectorsQuery.isSuccess) {
      return;
    }

    dispatch({
      type: 'updateSettings',
      patch: {
        ...mapSettingItemsToWorkbenchSettings(settingsItemsQuery.data),
        connectors: connectorsQuery.data.map((connector) =>
          mapDatabaseConnectorToSettingsConnector(connector),
        ),
      },
    });
    setHydrated(true);
    setOperationError(null);
  }, [connectorsQuery.data, connectorsQuery.isSuccess, dispatch, hydrated, settingsItemsQuery.data, settingsItemsQuery.isSuccess]);

  const handleUpdateGeneral = React.useCallback((patch: Partial<GeneralSettings>) => {
    dispatch({ type: 'updateSettingsSection', section: 'general', patch });
  }, [dispatch]);

  const handleUpdateTimeseries = React.useCallback((patch: Partial<TimeseriesSettings>) => {
    dispatch({ type: 'updateSettingsSection', section: 'timeseries', patch });
  }, [dispatch]);

  const handleUpdateScheduler = React.useCallback((patch: Partial<SchedulerSettings>) => {
    dispatch({ type: 'updateSettingsSection', section: 'scheduler', patch });
  }, [dispatch]);

  const handleUpdateModbusShare = React.useCallback((patch: Partial<ModbusShareSettings>) => {
    dispatch({ type: 'updateSettingsSection', section: 'modbus_share', patch });
  }, [dispatch]);

  const handleReset = React.useCallback(() => {
    dispatch({ type: 'resetSettingsToDefaults' });
  }, [dispatch]);

  const handleSave = React.useCallback(async () => {
    setIsSaving(true);
    setOperationError(null);
    try {
      for (const entry of buildPersistableSettingEntries(state.settings)) {
        await updateSettingKeyMutation.mutateAsync(entry);
      }
    } catch (error) {
      setOperationError(error instanceof Error ? error.message : '儲存設定失敗');
    } finally {
      setIsSaving(false);
    }
  }, [state.settings, updateSettingKeyMutation]);

  const handleAddConnector = React.useCallback(async () => {
    const draftConnector = makeDefaultConnector(state.settings.connectors.length + 1);
    setOperationError(null);
    try {
      const created = await createConnectorMutation.mutateAsync(
        buildCreateConnectorRequest(draftConnector),
      );
      dispatch({
        type: 'updateSettings',
        patch: {
          connectors: [
            ...state.settings.connectors,
            mapDatabaseConnectorToSettingsConnector(created, draftConnector.password ?? ''),
          ],
        },
      });
    } catch (error) {
      setOperationError(error instanceof Error ? error.message : '建立連接器失敗');
    }
  }, [createConnectorMutation, dispatch, state.settings.connectors]);

  const handleUpdateConnector = React.useCallback((id: string, patch: Partial<SettingsConnector>) => {
    dispatch({ type: 'updateConnector', id, patch });
    const currentConnector = state.settings.connectors.find((connector) => connector.id === id);
    if (!currentConnector) {
      return;
    }

    const nextConnector = { ...currentConnector, ...patch };
    setOperationError(null);
    void updateConnectorMutation
      .mutateAsync({
        id,
        data: buildUpdateConnectorRequest(nextConnector),
      })
      .then((savedConnector) => {
        dispatch({
          type: 'updateConnector',
          id,
          patch: mapDatabaseConnectorToSettingsConnector(
            savedConnector,
            nextConnector.password ?? '',
          ),
        });
      })
      .catch((error) => {
        setOperationError(error instanceof Error ? error.message : '更新連接器失敗');
      });
  }, [dispatch, state.settings.connectors, updateConnectorMutation]);

  const handleRemoveConnector = React.useCallback((id: string) => {
    setOperationError(null);
    void deleteConnectorMutation
      .mutateAsync(id)
      .then(() => {
        dispatch({ type: 'removeConnector', id });
      })
      .catch((error) => {
        setOperationError(error instanceof Error ? error.message : '刪除連接器失敗');
      });
  }, [deleteConnectorMutation, dispatch]);

  const handleTestConnector = React.useCallback((id: string) => {
    const currentConnector = state.settings.connectors.find((connector) => connector.id === id);
    dispatch({ type: 'startConnectorTest', id });
    setOperationError(null);
    void testConnectorMutation
      .mutateAsync(id)
      .then((testedConnector) => {
        const connector = mapDatabaseConnectorToSettingsConnector(
          testedConnector,
          currentConnector?.password ?? '',
        );
        dispatch({
          type: 'completeConnectorTest',
          id,
          result: {
            status: connector.status === 'testing' ? 'unknown' : connector.status,
            last_check_at: connector.last_check_at ?? new Date().toISOString(),
            last_check_error: connector.last_check_error,
          },
        });
        dispatch({
          type: 'updateConnector',
          id,
          patch: connector,
        });
      })
      .catch((error) => {
        const message = error instanceof Error ? error.message : '測試連接器失敗';
        dispatch({
          type: 'completeConnectorTest',
          id,
          result: {
            status: 'unreachable',
            last_check_at: new Date().toISOString(),
            last_check_error: message,
          },
        });
        setOperationError(message);
      });
  }, [dispatch, state.settings.connectors, testConnectorMutation]);

  if (!hydrated) {
    if (settingsItemsQuery.isError || connectorsQuery.isError) {
      return (
        <div
          className="rounded-2xl border border-red-500/20 bg-red-500/5 px-5 py-4 text-sm text-red-200"
          data-testid="settings-backend-error"
        >
          無法載入系統設定。請檢查後端資料庫連線後重試。
        </div>
      );
    }

    return (
      <div
        className="rounded-2xl border border-slate-700/60 bg-slate-900/40 px-5 py-4 text-sm text-slate-300"
        data-testid="settings-backend-loading"
      >
        載入系統設定...
      </div>
    );
  }

  return (
    <div className="space-y-6 pb-24" data-testid="settings-page">
      {operationError && (
        <div
          className="rounded-2xl border border-red-500/20 bg-red-500/5 px-5 py-4 text-sm text-red-200"
          data-testid="settings-operation-error"
        >
          {operationError}
        </div>
      )}
      {/* 漸層標題卡 */}
      <SettingsHeader />
      
      {/* 連接器池區塊 */}
      <ConnectorPoolSection
        connectors={state.settings.connectors}
        onAddConnector={handleAddConnector}
        onUpdateConnector={handleUpdateConnector}
        onRemoveConnector={handleRemoveConnector}
        onTestConnector={handleTestConnector}
      />

      {/* 時序策略與排程器並排 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <TimeseriesSection
          settings={state.settings.timeseries}
          onChange={handleUpdateTimeseries}
        />
        <SchedulerSection
          settings={state.settings.scheduler}
          onChange={handleUpdateScheduler}
        />
      </div>

      {/* Local Modbus Share */}
      <ModbusShareSection
        settings={state.settings.modbus_share}
        onChange={handleUpdateModbusShare}
      />

      {/* 偏好、API 與診斷三欄並排 */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <UiSection
          settings={state.settings.general}
          onChange={handleUpdateGeneral}
        />
        <ApiSection
          settings={state.settings.general}
          onChange={handleUpdateGeneral}
        />
        <DiagnosticsSection
          settings={state.settings.general}
          onChange={handleUpdateGeneral}
        />
      </div>

      {/* Sticky 底部控制列 */}
      <div className="sticky bottom-4 z-10">
        <SaveBar onReset={handleReset} onSave={isSaving ? () => undefined : handleSave} />
      </div>
    </div>
  );
}
