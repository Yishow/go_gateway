import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { SettingsHeader } from '../../../src/features/datalink/workbench-v2/settings/SettingsHeader';
import { TimeseriesSection } from '../../../src/features/datalink/workbench-v2/settings/TimeseriesSection';
import { SchedulerSection } from '../../../src/features/datalink/workbench-v2/settings/SchedulerSection';
import { ModbusShareSection } from '../../../src/features/datalink/workbench-v2/settings/ModbusShareSection';
import { UiSection } from '../../../src/features/datalink/workbench-v2/settings/UiSection';
import { ApiSection } from '../../../src/features/datalink/workbench-v2/settings/ApiSection';
import { DiagnosticsSection } from '../../../src/features/datalink/workbench-v2/settings/DiagnosticsSection';
import { SaveBar } from '../../../src/features/datalink/workbench-v2/settings/SaveBar';
import { SettingsPage } from '../../../src/features/datalink/workbench-v2/settings/SettingsPage';
import { SettingsStatus } from '../../../src/features/datalink/workbench-v2/settings/SettingsStatus';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { settingsAPI, dbTargetAPI } from '../../../src/services/datalink';
import { modbusShareAPI } from '../../../src/services/modbusShare';

// 模擬 i18next 避免語系載入錯誤
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, fallback?: unknown) => {
      if (key === 'errors.modbus_share_revision_conflict') return 'Workspace configuration revision conflict.';
      if (key === 'errors.retry') return 'Retry';
      return typeof fallback === 'string' ? fallback : key;
    },
  }),
}));

vi.mock('../../../src/services/datalink', () => ({
  settingsAPI: {
    listItems: vi.fn(),
    updateKey: vi.fn(),
  },
  dbTargetAPI: {
    listConnectors: vi.fn(),
    createConnector: vi.fn(),
    updateConnector: vi.fn(),
    deleteConnector: vi.fn(),
    testConnector: vi.fn(),
  },
}));

vi.mock('../../../src/services/modbusShare', () => ({
  modbusShareAPI: {
    start: vi.fn(),
    stop: vi.fn(),
  },
}));

describe('Settings Components', () => {
  beforeEach(() => {
    vi.mocked(settingsAPI.listItems).mockResolvedValue([]);
    vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([]);
  });

  const renderSettingsPage = () => {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });

    return render(
      <QueryClientProvider client={queryClient}>
        <SettingsPage state={INITIAL_STATE} dispatch={vi.fn()} />
      </QueryClientProvider>,
    );
  };

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

      rerender(<ModbusShareSection settings={{ ...INITIAL_STATE.settings.modbus_share, enabled: true }} onChange={onChange} />);
      expect(screen.getByTestId('share-base-register-value')).toHaveTextContent('自動由已保存候選配置');
      expect(screen.queryByTestId('share-base-register-input')).not.toBeInTheDocument();

      const capacityInput = screen.getByText('暫存器容量').closest('div')?.querySelector('input');
      expect(capacityInput).toHaveValue(32768);
      fireEvent.change(capacityInput!, { target: { value: '4096' } });
      expect(onChange).toHaveBeenCalledWith({ capacity_registers: 4096 });
    });

    it('exposes pressed state and prevents duplicate activation while pending', () => {
      const onChange = vi.fn();
      render(
        <ModbusShareSection
          settings={{ ...INITIAL_STATE.settings.modbus_share, enabled: false }}
          onChange={onChange}
          pending={true}
        />,
      );

      const toggle = screen.getByRole('button', { name: '未啟用' });
      expect(toggle).toHaveAttribute('aria-pressed', 'false');
      expect(toggle).toBeDisabled();
      fireEvent.click(toggle);
      expect(onChange).not.toHaveBeenCalled();
    });

    it('shows pending hydration before a failed runtime status', () => {
      render(
        <ModbusShareSection
          settings={{ ...INITIAL_STATE.settings.modbus_share, enabled: true }}
          onChange={vi.fn()}
          shareStatus={{
            enabled: false,
            port: 0,
            address: '',
            bind_state: 'fail',
            mapping_count: 0,
            hydration_state: 'pending',
          }}
        />,
      );

      expect(screen.getByTestId('share-runtime-status')).toHaveTextContent('pending');
      expect(screen.getByTestId('share-runtime-status')).not.toHaveTextContent('failed');
    });

    it('sends the current settings revision from Start/Stop controls', async () => {
      vi.mocked(modbusShareAPI.stop).mockResolvedValue({
        enabled: true,
        running: false,
        port: 5020,
        address: '127.0.0.1',
        bind_state: 'stopped',
        mapping_count: 0,
        settings_revision: 'settings-9',
      });
      render(
        <ModbusShareSection
          settings={{ ...INITIAL_STATE.settings.modbus_share, enabled: true, settings_revision: 'settings-9' }}
          onChange={vi.fn()}
          shareStatus={{
            enabled: true,
            running: true,
            port: 5020,
            address: '127.0.0.1',
            bind_state: 'pass',
            mapping_count: 0,
            settings_revision: 'settings-9',
          }}
        />,
      );

      fireEvent.click(screen.getByRole('button', { name: 'Stop' }));
      await waitFor(() => expect(modbusShareAPI.stop).toHaveBeenCalledWith('settings-9'));
    });

    it('renders a typed 409 conflict safely without backend details', async () => {
      vi.mocked(modbusShareAPI.stop).mockRejectedValue({
        response: {
          status: 409,
          data: {
            error: {
              code: 'modbus_share_revision_conflict',
              message: 'current revision settings-secret-10',
              action: 'reload secret internal state',
              retryable: true,
              request_id: 'req-settings-conflict-1',
            },
          },
        },
      });
      render(
        <ModbusShareSection
          settings={{ ...INITIAL_STATE.settings.modbus_share, enabled: true, settings_revision: 'settings-9' }}
          onChange={vi.fn()}
          shareStatus={{
            enabled: true,
            running: true,
            port: 5020,
            address: '127.0.0.1',
            bind_state: 'pass',
            mapping_count: 0,
            settings_revision: 'settings-9',
          }}
        />,
      );

      fireEvent.click(screen.getByRole('button', { name: 'Stop' }));
      const error = await screen.findByTestId('share-runtime-operation-error');
      expect(error).toHaveTextContent('Workspace configuration revision conflict.');
      expect(error).toHaveTextContent('req-settings-conflict-1');
      expect(error).not.toHaveTextContent('settings-secret-10');
      expect(error).not.toHaveTextContent('reload secret internal state');
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
    it('keeps typed request id and localized retry action for backend save errors', () => {
      const onRetry = vi.fn();
      render(
        <SettingsStatus
          operationError={{
            code: 'settings_update_failed',
            title: 'settings_update_failed',
            message: '系統設定保存失敗。',
            action: '重試',
            requestId: 'req-settings-7',
            retryable: true,
          }}
          onRetry={onRetry}
        />,
      );

      expect(screen.getByTestId('settings-operation-error')).toHaveTextContent('系統設定保存失敗。');
      expect(screen.getByTestId('settings-operation-error')).toHaveTextContent('req-settings-7');
      expect(screen.getByRole('button', { name: '重試' })).toBeInTheDocument();
      expect(screen.queryByText(/raw|backend/i)).not.toBeInTheDocument();
      fireEvent.click(screen.getByRole('button', { name: '重試' }));
      expect(onRetry).toHaveBeenCalledTimes(1);
    });

    it('retries the failed connector add instead of invoking Save All Settings', async () => {
      vi.mocked(dbTargetAPI.createConnector)
        .mockRejectedValueOnce({
          error: {
            code: 'settings_update_failed',
            retryable: true,
            request_id: 'req-connector-add-1',
          },
        })
        .mockResolvedValueOnce({
          id: 'connector-1',
          name: 'SQLite Connector 1',
          kind: 'sqlite',
          connection_config: { path: 'gateway.db' },
          status: 'ready',
          last_check_error: '',
          enabled: true,
          default_write_interval_seconds: 5,
          created_at: '2026-08-25T00:00:00Z',
          updated_at: '2026-08-25T00:00:00Z',
        });
      const updateKey = vi.mocked(settingsAPI.updateKey);

      renderSettingsPage();
      await screen.findByTestId('settings-page');
      fireEvent.click(screen.getByRole('button', { name: /新增連接器/ }));
      await screen.findByTestId('settings-operation-error');
      fireEvent.click(screen.getByRole('button', { name: 'Retry' }));

      await waitFor(() => expect(dbTargetAPI.createConnector).toHaveBeenCalledTimes(2));
      expect(updateKey).not.toHaveBeenCalled();
    });

    it('renders all sections and triggers actions via dispatch', async () => {
      renderSettingsPage();

      expect(await screen.findByTestId('settings-page')).toBeInTheDocument();
      expect(screen.getByText('系統設定')).toBeInTheDocument();
      expect(screen.getByText('時序儲存策略')).toBeInTheDocument();
      expect(screen.getByText('排程器核心')).toBeInTheDocument();
      expect(screen.getByText('Local Modbus Share')).toBeInTheDocument();
      expect(screen.queryByTestId('share-base-register-input')).not.toBeInTheDocument();
      expect(screen.getByTestId('share-base-register-value')).toHaveTextContent('自動由已保存候選配置');
      expect(screen.getByText('介面與格式偏好')).toBeInTheDocument();
      expect(screen.getByText('API 連線端點')).toBeInTheDocument();
      expect(screen.getByText('日誌與開發診斷')).toBeInTheDocument();
    });

    it('marks Share settings dirty so activation cannot bypass an unsaved edit', async () => {
      const dispatch = vi.fn();
      const queryClient = new QueryClient({
        defaultOptions: {
          queries: { retry: false },
          mutations: { retry: false },
        },
      });

      render(
        <QueryClientProvider client={queryClient}>
          <SettingsPage state={INITIAL_STATE} dispatch={dispatch} />
        </QueryClientProvider>,
      );

      await screen.findByTestId('settings-page');
      fireEvent.click(screen.getByRole('button', { name: '啟用' }));

      // 編輯只更新資料與清除既有錯誤，save_state 留給實際發出的存檔請求。
      expect(dispatch).toHaveBeenLastCalledWith({
        type: 'updateSettingsSection',
        section: 'modbus_share',
        patch: expect.objectContaining({
          enabled: false,
          save_error: null,
        }),
      });
      expect(dispatch).not.toHaveBeenCalledWith(expect.objectContaining({
        patch: expect.objectContaining({ save_state: 'saving' }),
      }));
    });
  });
});
