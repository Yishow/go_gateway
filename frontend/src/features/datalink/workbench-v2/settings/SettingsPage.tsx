import * as React from 'react';
import type { WorkbenchV2State } from '../state/types';
import type { WorkbenchV2Action } from '../state/useWorkbenchV2State';
import { useDbTargetConnectorsQuery, useSettingsItemsQuery } from '../../../../hooks/datalink/useSettings';
import { SettingsSections } from './SettingsSections';
import { SettingsBackendStatus, SettingsStatus } from './SettingsStatus';
import { mapDatabaseConnectorToSettingsConnector, mapSettingItemsToWorkbenchSettings } from './backendMappings';
import type { ModbusShareStatus } from '../../../../types/modbusShare';
import { useSettingsOperations } from './useSettingsOperations';

interface SettingsPageProps {
  state: WorkbenchV2State;
  dispatch: React.Dispatch<WorkbenchV2Action>;
  shareStatus?: ModbusShareStatus | null;
}

/** System settings container; async mutation ownership lives in useSettingsOperations. */
export function SettingsPage({ state, dispatch, shareStatus }: SettingsPageProps) {
  const settingsItemsQuery = useSettingsItemsQuery();
  const connectorsQuery = useDbTargetConnectorsQuery();
  const [hydrated, setHydrated] = React.useState(false);
  const operations = useSettingsOperations(state, dispatch);
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
        connectors: connectorsQuery.data.map((connector) => mapDatabaseConnectorToSettingsConnector(connector)),
      },
    });
    setHydrated(true);
  }, [connectorsQuery.data, connectorsQuery.isSuccess, dispatch, hydrated, settingsItemsQuery.data, settingsItemsQuery.isSuccess]);

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
