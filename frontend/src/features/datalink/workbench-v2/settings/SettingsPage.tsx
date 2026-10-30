import { useOptionalTheme } from '../../../../contexts/ThemeContext';
import { themeWireValue } from '../../../../contexts/themePreference';
import * as React from 'react';
import { useQueryClient } from '@tanstack/react-query';
import type { WorkbenchV2State } from '../state/types';
import type { WorkbenchV2Action } from '../state/useWorkbenchV2State';
import { useDbTargetConnectorsQuery, useSettingsItemsQuery } from '../../../../hooks/datalink/useSettings';
import { SettingsSections } from './SettingsSections';
import { SettingsBackendStatus, SettingsStatus } from './SettingsStatus';
import { mapDatabaseConnectorToSettingsConnector, mapSettingItemsToWorkbenchSettings } from './backendMappings';
import type { ModbusShareStatus } from '../../../../types/modbusShare';
import { useSettingsOperations } from './useSettingsOperations';
import { modbusShareStatusKey } from '../../../../hooks/datalink/useModbusShareStatus';
import { studioV2WorkspaceKeys } from '../../../../hooks/datalink/keys';

interface SettingsPageProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  shareStatus?: ModbusShareStatus | null;
}

/** System settings container; async mutation ownership lives in useSettingsOperations. */
export function SettingsPage({ state, dispatch, shareStatus }: SettingsPageProps) {
  const theme = useOptionalTheme();
  const preference = theme?.preference;
  React.useEffect(() => {
    if (preference) dispatch({ type: 'updateSettingsSection', section: 'general', patch: { theme: themeWireValue(preference) } });
  }, [dispatch, preference]);
  const settingsItemsQuery = useSettingsItemsQuery();
  const connectorsQuery = useDbTargetConnectorsQuery();
  const [hydrated, setHydrated] = React.useState(false);
  const operations = useSettingsOperations(state, dispatch);
  const queryClient = useQueryClient();
  // Start/Stop 會改變後端 Share 狀態；重新讀取讓 Step 2／Step 4 與啟用前檢查看到同一份事實。
  const refreshShareStatus = React.useCallback(() => {
    void queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    void queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
  }, [queryClient]);
  const persistedBaseRegister = state.rules
    .filter((rule) => rule.persisted === true && rule.share_enabled && rule.share_start_register !== null)
    .map((rule) => rule.share_start_register as number)
    .sort((left, right) => left - right)[0] ?? null;

  React.useEffect(() => {
    if (hydrated || !settingsItemsQuery.isSuccess || !connectorsQuery.isSuccess) return;
    dispatch({
      type: 'updateSettings',
      patch: {
        ...mapSettingItemsToWorkbenchSettings(settingsItemsQuery.data),
        ...(preference ? { general: { ...mapSettingItemsToWorkbenchSettings(settingsItemsQuery.data).general, theme: themeWireValue(preference) } } : {}),
        connectors: connectorsQuery.data.map((connector) => mapDatabaseConnectorToSettingsConnector(connector)),
      },
    });
    setHydrated(true);
  }, [connectorsQuery.data, connectorsQuery.isSuccess, dispatch, hydrated, preference, settingsItemsQuery.data, settingsItemsQuery.isSuccess]);

  if (!hydrated) {
    return <SettingsBackendStatus hasError={settingsItemsQuery.isError || connectorsQuery.isError} error={settingsItemsQuery.error ?? connectorsQuery.error} />;
  }

  return (
    <div className="space-y-6 pb-24" data-testid="settings-page">
      <SettingsStatus operationError={operations.operationError} onRetry={operations.onRetry} />
      <SettingsSections
        state={state}
        isSaving={operations.isSaving}
        persistedBaseRegister={persistedBaseRegister}
        shareStatus={shareStatus}
        onShareLifecycleSettled={refreshShareStatus}
        onUpdateGeneral={operations.onUpdateGeneral}
        onUpdateTimeseries={operations.onUpdateTimeseries}
        onUpdateScheduler={operations.onUpdateScheduler}
        onUpdateModbusShare={operations.onUpdateModbusShare}
        onAddConnector={operations.onAddConnector}
        onUpdateConnector={operations.onUpdateConnector}
        onRemoveConnector={operations.onRemoveConnector}
        onTestConnector={operations.onTestConnector}
        onReset={operations.onReset}
        onSave={operations.onSave}
        pending={operations.pending}
      />
    </div>
  );
}
