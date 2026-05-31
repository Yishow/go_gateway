import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ProtocolSelector } from '../../../src/features/datalink/workbench-v2/steps/step1/ProtocolSelector';
import { ConnectionConfigForm } from '../../../src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm';
import { ConnectionTestPanel } from '../../../src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel';
import { DeviceTabRail } from '../../../src/features/datalink/workbench-v2/steps/step1/DeviceTabRail';
import { DeviceEditor } from '../../../src/features/datalink/workbench-v2/steps/step1/DeviceEditor';
import { Step1Device } from '../../../src/features/datalink/workbench-v2/steps/step1/Step1Device';
import { deviceAPI } from '../../../src/services/datalink';
import type { ConnectionTestResult } from '../../../src/types/datalink';
import {
  INITIAL_STATE,
  workbenchV2Reducer,
} from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Device, WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { DeviceListContext } from '../../../src/features/datalink/workbench-v2/state/deviceColors';

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: any) => {
      if (key === 'step1.panels.test_success_desc') {
        return `Successfully tested. Latency ${options?.latency} ms`;
      }
      return key;
    },
  }),
}));

describe('Step 1 Components & Integration', () => {
  const mockDevice: Device = {
    id: 'dev-01',
    name: 'PLC-生產線-01',
    description: 'Modbus TCP PLC',
    protocol: 'modbus_tcp',
    config: { host: '192.168.1.100', port: 502, slave_id: 1, timeout: 5 },
    status: 'draft',
    test: null,
    persisted: true,
    save_state: 'saved',
    save_error: null,
    availability_status: 'available',
    availability_reason: null,
    running: true,
  };

  const renderWithQueryClient = (ui: React.ReactElement) => {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    });

    return render(
      <QueryClientProvider client={queryClient}>
        {ui}
      </QueryClientProvider>,
    );
  };

  const Step1Harness = ({
    initialState = INITIAL_STATE,
    onContinue = vi.fn(),
  }: {
    initialState?: WorkbenchV2State;
    onContinue?: () => void;
  }) => {
    const [state, dispatch] = React.useReducer(workbenchV2Reducer, initialState);
    return <Step1Device state={state} dispatch={dispatch} onContinue={onContinue} />;
  };

  const buildSuccessResult = (): ConnectionTestResult => ({
    success: true,
    error: '',
    latency_ms: 55,
    connect: {
      status: 'success',
      message: 'connect ok',
      latency_ms: 21,
    },
    probe: {
      status: 'success',
      message: 'probe ok',
      latency_ms: 34,
    },
    can_activate: true,
    can_collect: true,
  });

  describe('ProtocolSelector', () => {
    it('應渲染 6 張協議卡片，且點擊時能正確觸發 onChange', () => {
      const handleChange = vi.fn();
      render(<ProtocolSelector value="modbus_tcp" onChange={handleChange} />);

      // 檢查是否渲染了 6 個卡片按鈕
      expect(screen.getByTestId('protocol-card-modbus_tcp')).toBeInTheDocument();
      expect(screen.getByTestId('protocol-card-modbus_rtu')).toBeInTheDocument();
      expect(screen.getByTestId('protocol-card-mqtt')).toBeInTheDocument();

      // 選中的卡片應有高亮樣式
      const selectedCard = screen.getByTestId('protocol-card-modbus_tcp');
      expect(selectedCard.className).toContain('border-blue-500');

      // 點擊卡片觸發 callback
      fireEvent.click(screen.getByTestId('protocol-card-mqtt'));
      expect(handleChange).toHaveBeenCalledWith('mqtt');
    });
  });

  describe('ConnectionConfigForm', () => {
    it('TCP 協議應渲染 host 與 port 欄位', () => {
      const handleChange = vi.fn();
      render(
        <ConnectionConfigForm
          protocol="modbus_tcp"
          config={mockDevice.config}
          onChange={handleChange}
        />
      );

      expect(screen.getByTestId('input-host')).toBeInTheDocument();
      expect(screen.getByTestId('input-port')).toBeInTheDocument();
      expect(screen.getByTestId('input-slave-id')).toBeInTheDocument();
    });

    it('RTU 協議應渲染 serial_port 與 baud_rate 欄位', () => {
      const handleChange = vi.fn();
      render(
        <ConnectionConfigForm
          protocol="modbus_rtu"
          config={{ port: '/dev/ttyUSB0', baud: 9600 }}
          onChange={handleChange}
        />
      );

      expect(screen.getByTestId('input-rtu-port')).toBeInTheDocument();
      expect(screen.getByTestId('select-baud')).toBeInTheDocument();
      expect(screen.getByTestId('select-parity')).toBeInTheDocument();
    });

    it('MQTT 協議應渲染 broker 與 client_id 欄位', () => {
      const handleChange = vi.fn();
      render(
        <ConnectionConfigForm
          protocol="mqtt"
          config={{ broker: 'mqtt://localhost:1883', client_id: 'g-01' }}
          onChange={handleChange}
        />
      );

      expect(screen.getByTestId('input-mqtt-broker')).toBeInTheDocument();
      expect(screen.getByTestId('input-mqtt-client-id')).toBeInTheDocument();
    });
  });

  describe('ConnectionTestPanel', () => {
    it('應展示 payload 預覽，且按鈕能正確控制與觸發測試', () => {
      const handleRunTest = vi.fn();
      render(<ConnectionTestPanel device={mockDevice} onRunTest={handleRunTest} />);

      // payload 預覽中應包含 host 等配置
      const preview = screen.getByTestId('payload-preview');
      expect(preview.textContent).toContain('modbus_tcp');
      expect(preview.textContent).toContain('192.168.1.100');

      // 點擊執行測試
      const btn = screen.getByTestId('run-test-button');
      expect(btn).not.toBeDisabled();
      fireEvent.click(btn);
      expect(handleRunTest).toHaveBeenCalledWith('dev-01');
    });

    it('測試成功後應顯示綠色成功 readiness panel', () => {
      const handleRunTest = vi.fn();
      const testedDevice: Device = {
        ...mockDevice,
        status: 'tested',
        test: {
          status: 'success',
          latency_ms: 55,
          stages: {},
        },
      };

      render(<ConnectionTestPanel device={testedDevice} onRunTest={handleRunTest} />);
      expect(screen.getByTestId('success-readiness-card')).toBeInTheDocument();
      expect(screen.getByText('Successfully tested. Latency 55 ms')).toBeInTheDocument();
    });
  });

  describe('DeviceTabRail', () => {
    it('應渲染裝置分頁、save chip 與改名 Input，且點擊刪除時觸發 confirm dialog', () => {
      const handleSelect = vi.fn();
      const handleAdd = vi.fn();
      const handleDelete = vi.fn();
      const handleRename = vi.fn();

      const devices = [
        mockDevice,
        { ...mockDevice, id: 'dev-02', name: '設備 2' },
      ];

      // 提供 DeviceListContext 確保 useDeviceColor 運作
      render(
        <DeviceListContext.Provider value={devices}>
          <DeviceTabRail
            devices={devices}
            selectedId="dev-01"
            onSelect={handleSelect}
            onAdd={handleAdd}
            onDelete={handleDelete}
            onRename={handleRename}
          />
        </DeviceListContext.Provider>
      );

      // 改名 input
      const renameInput = screen.getByTestId('input-rename-dev-01');
      expect(renameInput).toHaveValue('PLC-生產線-01');
      expect(screen.getByTestId('device-save-chip-dev-01')).toHaveTextContent('已儲存');
      fireEvent.change(renameInput, { target: { value: 'New Name' } });
      expect(handleRename).toHaveBeenCalledWith('dev-01', 'New Name');

      // 刪除對話框
      const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true);
      const delBtn = screen.getByTestId('btn-delete-device-dev-02');
      fireEvent.click(delBtn);

      expect(confirmSpy).toHaveBeenCalled();
      expect(handleDelete).toHaveBeenCalledWith('dev-02');
      confirmSpy.mockRestore();
    });

    it('應保留 unavailable 裝置並顯示可辨識狀態 chip', () => {
      const unavailableDevice: Device = {
        ...mockDevice,
        availability_status: 'unavailable',
        availability_reason: 'device form is invalid',
        running: false,
      };

      render(
        <DeviceListContext.Provider value={[unavailableDevice]}>
          <DeviceTabRail
            devices={[unavailableDevice]}
            selectedId="dev-01"
            onSelect={vi.fn()}
            onAdd={vi.fn()}
            onDelete={vi.fn()}
            onRename={vi.fn()}
          />
        </DeviceListContext.Provider>
      );

      expect(screen.getByTestId('device-tab-dev-01')).toBeInTheDocument();
      expect(screen.getByTestId('device-availability-chip-dev-01')).toHaveTextContent('step1.status.unavailable');
    });
  });

  describe('DeviceEditor', () => {
    it('應能顯示 save banner、綁定名稱改值反映 state，切換協議時觸發 reset', () => {
      const handleUpdate = vi.fn();
      const handleUpdateConfig = vi.fn();
      const handleChangeProtocol = vi.fn();

      render(
        <DeviceEditor
          device={mockDevice}
          onUpdate={handleUpdate}
          onUpdateConfig={handleUpdateConfig}
          onChangeProtocol={handleChangeProtocol}
        />
      );

      expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'saved');
      expect(screen.getByTestId('device-save-banner')).toHaveTextContent('已儲存');

      const nameInput = screen.getByTestId('editor-input-name');
      fireEvent.change(nameInput, { target: { value: 'Temp PLC' } });
      expect(handleUpdate).toHaveBeenCalledWith({ name: 'Temp PLC' });

      // 切換協議
      fireEvent.click(screen.getByTestId('protocol-card-mqtt'));
      expect(handleChangeProtocol).toHaveBeenCalledWith('mqtt');
    });

    it('應在 save-error 時保留錯誤訊息 marker', () => {
      render(
        <DeviceEditor
          device={{ ...mockDevice, save_state: 'save-error', save_error: 'save failed' }}
          onUpdate={vi.fn()}
          onUpdateConfig={vi.fn()}
          onChangeProtocol={vi.fn()}
        />
      );

      expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'save-error');
      expect(screen.getByTestId('device-save-error-text')).toHaveTextContent('save failed');
    });

    it('應在 unavailable 時顯示原因且不移除 editor 狀態', () => {
      render(
        <DeviceEditor
          device={{
            ...mockDevice,
            availability_status: 'unavailable',
            availability_reason: 'device form is invalid',
            running: false,
          }}
          onUpdate={vi.fn()}
          onUpdateConfig={vi.fn()}
          onChangeProtocol={vi.fn()}
        />
      );

      expect(screen.getByTestId('device-save-banner')).toBeInTheDocument();
      expect(screen.getByTestId('device-availability-reason')).toHaveTextContent(
        'step1.status.unavailable: device form is invalid',
      );
    });
  });

  describe('Step1Device Live Diagnostics & Continue Gate', () => {
    afterEach(() => {
      vi.restoreAllMocks();
    });

    it('點擊執行測試應發出 draft diagnostics request，並在成功後允許繼續', async () => {
      const onContinue = vi.fn();
      let resolveResult: ((value: ConnectionTestResult) => void) | undefined;
      const diagnosticsPromise = new Promise<ConnectionTestResult>((resolve) => {
        resolveResult = resolve;
      });
      const diagnosticsSpy = vi
        .spyOn(deviceAPI, 'testDraftConnection')
        .mockReturnValue(diagnosticsPromise);

      renderWithQueryClient(<Step1Harness onContinue={onContinue} />);

      const runBtn = screen.getByTestId('run-test-button');
      const continueBtn = screen.getByTestId('btn-continue-step1');
      expect(continueBtn).toBeDisabled();

      fireEvent.click(runBtn);

      await waitFor(() => {
        expect(diagnosticsSpy).toHaveBeenCalledWith({
          protocol: 'modbus_tcp',
          connection_config: {
            host: '192.168.1.100',
            port: 502,
            slave_id: 1,
            timeout: 5,
          },
        });
      });
      expect(runBtn).toBeDisabled();
      expect(screen.getByText('step1.buttons.testing')).toBeInTheDocument();

      if (!resolveResult) {
        throw new Error('expected diagnostics promise resolver to be set');
      }
      resolveResult(buildSuccessResult());

      await waitFor(() => {
        expect(screen.getByTestId('success-readiness-card')).toBeInTheDocument();
      });
      expect(screen.getByTestId('btn-continue-step1')).not.toBeDisabled();

      fireEvent.click(screen.getByTestId('btn-continue-step1'));
      expect(onContinue).toHaveBeenCalled();
    });

    it('diagnostics failure 不得讓 device 自動通過 continue gate', async () => {
      vi.spyOn(deviceAPI, 'testDraftConnection').mockResolvedValue({
        success: false,
        error: '讀取探測失敗: bad register',
        latency_ms: 18,
        connect: {
          status: 'success',
          message: 'connect ok',
          latency_ms: 12,
        },
        probe: {
          status: 'failed',
          error: 'bad register',
          latency_ms: 6,
        },
        can_activate: false,
        can_collect: false,
      });

      renderWithQueryClient(
        <Step1Harness
          initialState={{
            ...INITIAL_STATE,
            devices: [{ ...INITIAL_STATE.devices[0], test: null, status: 'draft' }],
          }}
        />,
      );

      const runBtn = screen.getByTestId('run-test-button');
      fireEvent.click(runBtn);

      await waitFor(() => {
        expect(screen.getByText('bad register')).toBeInTheDocument();
      });
      expect(screen.getByTestId('btn-continue-step1')).toBeDisabled();
      expect(screen.queryByTestId('success-readiness-card')).not.toBeInTheDocument();
    });
  });
});
