import type {
  GeneralSettings,
  ModbusShareSettings,
  SchedulerSettings,
  SettingsConnector,
  TimeseriesSettings,
  WorkbenchV2State,
} from '../state/types';
import type { ModbusShareStatus } from '../../../../types/modbusShare';
import { SettingsHeader } from './SettingsHeader';
import { ConnectorPoolSection } from './ConnectorPoolSection';
import { TimeseriesSection } from './TimeseriesSection';
import { SchedulerSection } from './SchedulerSection';
import { ModbusShareSection } from './ModbusShareSection';
import { UiSection } from './UiSection';
import { ApiSection } from './ApiSection';
import { DiagnosticsSection } from './DiagnosticsSection';
import { SaveBar } from './SaveBar';

interface SettingsSectionsProps {
  state: WorkbenchV2State;
  isSaving: boolean;
  persistedBaseRegister: number | null;
  shareStatus?: ModbusShareStatus | null;
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
  pending: boolean;
}

export function SettingsSections({
  state,
  isSaving,
  persistedBaseRegister,
  shareStatus,
  onUpdateGeneral,
  onUpdateTimeseries,
  onUpdateScheduler,
  onUpdateModbusShare,
  onAddConnector,
  onUpdateConnector,
  onRemoveConnector,
  onTestConnector,
  onReset,
  onSave,
  pending,
}: SettingsSectionsProps) {
  return (
    <>
      <SettingsHeader />
      <ConnectorPoolSection
        connectors={state.settings.connectors}
        onAddConnector={onAddConnector}
        onUpdateConnector={onUpdateConnector}
        onRemoveConnector={onRemoveConnector}
        onTestConnector={onTestConnector}
      />
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <TimeseriesSection settings={state.settings.timeseries} onChange={onUpdateTimeseries} />
        <SchedulerSection settings={state.settings.scheduler} onChange={onUpdateScheduler} />
      </div>
      <ModbusShareSection
        settings={state.settings.modbus_share}
        onChange={onUpdateModbusShare}
        pending={pending}
        persistedBaseRegister={persistedBaseRegister}
        shareStatus={shareStatus}
      />
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <UiSection settings={state.settings.general} onChange={onUpdateGeneral} />
        <ApiSection settings={state.settings.general} onChange={onUpdateGeneral} />
        <DiagnosticsSection settings={state.settings.general} onChange={onUpdateGeneral} />
      </div>
      <div className="sticky bottom-4 z-10">
        <SaveBar onReset={onReset} onSave={isSaving ? () => undefined : onSave} />
      </div>
    </>
  );
}
