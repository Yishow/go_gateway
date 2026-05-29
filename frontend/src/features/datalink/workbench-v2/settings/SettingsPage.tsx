import * as React from 'react';
import type { WorkbenchV2State, GeneralSettings, TimeseriesSettings, SchedulerSettings, ModbusShareSettings } from '../state/types';
import type { WorkbenchV2Action } from '../state/useWorkbenchV2State';
import { SettingsHeader } from './SettingsHeader';
import { ConnectorPoolSection } from './ConnectorPoolSection';
import { TimeseriesSection } from './TimeseriesSection';
import { SchedulerSection } from './SchedulerSection';
import { ModbusShareSection } from './ModbusShareSection';
import { UiSection } from './UiSection';
import { ApiSection } from './ApiSection';
import { DiagnosticsSection } from './DiagnosticsSection';
import { SaveBar } from './SaveBar';

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

  const handleSave = React.useCallback(() => {
    // eslint-disable-next-line no-console
    console.warn('Saving all settings to backend...', state.settings);
  }, [state.settings]);

  return (
    <div className="space-y-6 pb-24" data-testid="settings-page">
      {/* 漸層標題卡 */}
      <SettingsHeader />
      
      {/* 連接器池區塊 */}
      <ConnectorPoolSection
        connectors={state.settings.connectors}
        dispatch={dispatch}
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
        <SaveBar onReset={handleReset} onSave={handleSave} />
      </div>
    </div>
  );
}
