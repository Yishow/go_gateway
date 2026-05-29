import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { SettingsHeader } from '../../../src/features/datalink/workbench-v2/settings/SettingsHeader';
import { TimeseriesSection } from '../../../src/features/datalink/workbench-v2/settings/TimeseriesSection';
import { SchedulerSection } from '../../../src/features/datalink/workbench-v2/settings/SchedulerSection';
import { ModbusShareSection } from '../../../src/features/datalink/workbench-v2/settings/ModbusShareSection';
import { UiSection } from '../../../src/features/datalink/workbench-v2/settings/UiSection';
import { ApiSection } from '../../../src/features/datalink/workbench-v2/settings/ApiSection';
import { DiagnosticsSection } from '../../../src/features/datalink/workbench-v2/settings/DiagnosticsSection';
import { SaveBar } from '../../../src/features/datalink/workbench-v2/settings/SaveBar';
import { SettingsPage } from '../../../src/features/datalink/workbench-v2/settings/SettingsPage';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

// 模擬 i18next 避免語系載入錯誤
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, fallback?: string) => fallback || key,
  }),
}));

describe('Settings Components', () => {
  describe('SettingsHeader', () => {
    it('renders header title and subtitle', () => {
      render(<SettingsHeader />);
      expect(screen.getByText('系統設定')).toBeInTheDocument();
      expect(screen.getByText(/管理閘道全局連線/)).toBeInTheDocument();
    });
  });

  describe('TimeseriesSection', () => {
    it('renders timeseries inputs and changes value', () => {
      const onChange = vi.fn();
      render(<TimeseriesSection settings={INITIAL_STATE.settings.timeseries} onChange={onChange} />);
      
      const precisionSelect = screen.getByText('時間精度').closest('div')?.querySelector('select');
      const partitionSelect = screen.getByText('分區間隔').closest('div')?.querySelector('select');
      const retentionInput = screen.getByText('歷史保留天數').closest('div')?.querySelector('input');

      expect(precisionSelect).toHaveValue('millisecond');
      expect(partitionSelect).toHaveValue('daily');

      fireEvent.change(precisionSelect!, { target: { value: 'second' } });
      expect(onChange).toHaveBeenCalledWith({ write_precision: 'second' });

      fireEvent.change(retentionInput!, { target: { value: '180' } });
      expect(onChange).toHaveBeenCalledWith({ retention_days: 180 });
    });
  });

  describe('SchedulerSection', () => {
    it('renders scheduler inputs and triggers change', () => {
      const onChange = vi.fn();
      render(<SchedulerSection settings={INITIAL_STATE.settings.scheduler} onChange={onChange} />);

      const intervalInput = screen.getByText('預設輪詢週期').closest('div')?.querySelector('input');
      const retryInput = screen.getByText('連線失敗重試上限').closest('div')?.querySelector('input');
      const toggle = screen.getByText('開機自動啟動').closest('.justify-between')?.querySelector('button');

      expect(intervalInput).toHaveValue(1000);
      expect(retryInput).toHaveValue(3);

      fireEvent.change(intervalInput!, { target: { value: '2000' } });
      expect(onChange).toHaveBeenCalledWith({ default_interval_ms: 2000 });

      fireEvent.click(toggle!);
      expect(onChange).toHaveBeenCalledWith({ auto_start: false });
    });
  });

  describe('ModbusShareSection', () => {
    it('hides inputs when disabled and shows them when enabled', () => {
      const onChange = vi.fn();
      const settingsDisabled = { ...INITIAL_STATE.settings.modbus_share, enabled: false };
      const { rerender } = render(<ModbusShareSection settings={settingsDisabled} onChange={onChange} />);

      expect(screen.getByText(/Modbus 轉發服務已停用/)).toBeInTheDocument();

      rerender(<ModbusShareSection settings={INITIAL_STATE.settings.modbus_share} onChange={onChange} />);
      const baseRegisterInput = screen.getByText('起始暫存器位址').closest('div')?.querySelector('input');
      expect(baseRegisterInput).toHaveValue(40001);

      fireEvent.change(baseRegisterInput!, { target: { value: '40010' } });
      expect(onChange).toHaveBeenCalledWith({ base_register: 40010 });
    });
  });

  describe('UiSection / ApiSection / DiagnosticsSection', () => {
    it('renders ui section preferences and triggers change', () => {
      const onChange = vi.fn();
      render(<UiSection settings={INITIAL_STATE.settings.general} onChange={onChange} />);
      
      const themeSelect = screen.getByText('視覺主題').closest('div')?.querySelector('select');
      fireEvent.change(themeSelect!, { target: { value: 'light' } });
      expect(onChange).toHaveBeenCalledWith({ theme: 'light' });
    });

    it('renders api section settings and triggers change', () => {
      const onChange = vi.fn();
      render(<ApiSection settings={INITIAL_STATE.settings.general} onChange={onChange} />);

      const apiBaseInput = screen.getByText('API 伺服器位址 (API Base)').closest('div')?.querySelector('input');
      fireEvent.change(apiBaseInput!, { target: { value: 'http://127.0.0.1:9090' } });
      expect(onChange).toHaveBeenCalledWith({ api_base: 'http://127.0.0.1:9090' });
    });

    it('renders diagnostics settings and triggers change', () => {
      const onChange = vi.fn();
      render(<DiagnosticsSection settings={INITIAL_STATE.settings.general} onChange={onChange} />);

      const logLevelSelect = screen.getByText('日誌等級 (Log Level)').closest('div')?.querySelector('select');
      fireEvent.change(logLevelSelect!, { target: { value: 'debug' } });
      expect(onChange).toHaveBeenCalledWith({ log_level: 'debug' });
    });
  });

  describe('SaveBar', () => {
    beforeEach(() => {
      vi.spyOn(window, 'confirm');
    });

    it('calls onReset when confirm is accepted', () => {
      const onReset = vi.fn();
      const onSave = vi.fn();
      render(<SaveBar onReset={onReset} onSave={onSave} />);

      vi.mocked(window.confirm).mockReturnValue(true);
      fireEvent.click(screen.getByText('重設為預設'));
      expect(window.confirm).toHaveBeenCalled();
      expect(onReset).toHaveBeenCalled();
    });

    it('does not call onReset when confirm is rejected', () => {
      const onReset = vi.fn();
      const onSave = vi.fn();
      render(<SaveBar onReset={onReset} onSave={onSave} />);

      vi.mocked(window.confirm).mockReturnValue(false);
      fireEvent.click(screen.getByText('重設為預設'));
      expect(window.confirm).toHaveBeenCalled();
      expect(onReset).not.toHaveBeenCalled();
    });

    it('spies on Save All Settings', () => {
      const onReset = vi.fn();
      const onSave = vi.fn();
      const consoleWarnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});

      render(<SaveBar onReset={onReset} onSave={onSave} />);
      fireEvent.click(screen.getByText('儲存所有設定'));
      expect(onSave).toHaveBeenCalled();
      consoleWarnSpy.mockRestore();
    });
  });

  describe('SettingsPage Container', () => {
    it('renders all sections and triggers actions via dispatch', () => {
      const dispatch = vi.fn();
      render(<SettingsPage state={INITIAL_STATE} dispatch={dispatch} />);

      expect(screen.getByTestId('settings-page')).toBeInTheDocument();
      expect(screen.getByText('系統設定')).toBeInTheDocument();
      expect(screen.getByText('資料庫連接器池')).toBeInTheDocument();
      expect(screen.getByText('時序儲存策略')).toBeInTheDocument();
      expect(screen.getByText('排程器核心')).toBeInTheDocument();
      expect(screen.getByText('Local Modbus Share')).toBeInTheDocument();
      expect(screen.getByText('介面與格式偏好')).toBeInTheDocument();
      expect(screen.getByText('API 連線端點')).toBeInTheDocument();
      expect(screen.getByText('日誌與開發診斷')).toBeInTheDocument();
    });
  });
});
