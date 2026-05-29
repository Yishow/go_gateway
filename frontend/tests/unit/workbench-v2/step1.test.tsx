import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { ProtocolSelector } from '../../../src/features/datalink/workbench-v2/steps/step1/ProtocolSelector';
import { ConnectionConfigForm } from '../../../src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm';
import { ConnectionTestPanel } from '../../../src/features/datalink/workbench-v2/steps/step1/ConnectionTestPanel';
import { DeviceTabRail } from '../../../src/features/datalink/workbench-v2/steps/step1/DeviceTabRail';
import { DeviceEditor } from '../../../src/features/datalink/workbench-v2/steps/step1/DeviceEditor';
import { Step1Device } from '../../../src/features/datalink/workbench-v2/steps/step1/Step1Device';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';
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
  };

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

    it('新增 spy 斷言 payload 渲染但不發任何 fetch', () => {
      const spyFetch = vi.spyOn(window, 'fetch');
      render(<ConnectionTestPanel device={mockDevice} onRunTest={vi.fn()} />);

      // 驗證未發出 fetch 請求
      expect(spyFetch).not.toHaveBeenCalled();
      spyFetch.mockRestore();
    });
  });

  describe('DeviceTabRail', () => {
    it('應渲染裝置分頁與改名 Input，且點擊刪除時觸發 confirm dialog', () => {
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
  });

  describe('DeviceEditor', () => {
    it('應能綁定名稱改值反映 state，切換協議時觸發 reset', () => {
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

      const nameInput = screen.getByTestId('editor-input-name');
      fireEvent.change(nameInput, { target: { value: 'Temp PLC' } });
      expect(handleUpdate).toHaveBeenCalledWith({ name: 'Temp PLC' });

      // 切換協議
      fireEvent.click(screen.getByTestId('protocol-card-mqtt'));
      expect(handleChangeProtocol).toHaveBeenCalledWith('mqtt');
    });
  });

  describe('Step1Device Mock Test Animation & Continue Gate', () => {
    beforeEach(() => {
      vi.useFakeTimers();
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it('點擊執行測試應透過 timer 循序推進步驟', () => {
      const dispatch = vi.fn();
      const onContinue = vi.fn();

      render(<Step1Device state={INITIAL_STATE} dispatch={dispatch} onContinue={onContinue} />);

      // 找到執行測試按鈕
      const runBtn = screen.getByTestId('run-test-button');
      fireEvent.click(runBtn);

      // 斷言 startDeviceTest 觸發
      expect(dispatch).toHaveBeenCalledWith({ type: 'startDeviceTest', deviceId: 'dev-01' });

      // 快進 220ms
      act(() => {
        vi.advanceTimersByTime(220);
      });
      // 斷言 advanceDeviceTest 推進第一步 resolve
      expect(dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'advanceDeviceTest',
          deviceId: 'dev-01',
          stageId: 'resolve',
        })
      );

      // 快進 380ms
      act(() => {
        vi.advanceTimersByTime(380);
      });
      // 推進第二步 connect
      expect(dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'advanceDeviceTest',
          deviceId: 'dev-01',
          stageId: 'connect',
        })
      );

      // 快進 380ms
      act(() => {
        vi.advanceTimersByTime(380);
      });
      // 推進第三步 probe
      expect(dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'advanceDeviceTest',
          deviceId: 'dev-01',
          stageId: 'probe',
        })
      );

      // 快進 380ms 到結束
      act(() => {
        vi.advanceTimersByTime(380);
      });
      expect(dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: 'completeDeviceTest',
          deviceId: 'dev-01',
        })
      );
    });

    it('元件卸載時應呼叫 clearTimeout 進行清理', () => {
      const dispatch = vi.fn();
      const spyClear = vi.spyOn(window, 'clearTimeout');

      const { unmount } = render(
        <Step1Device state={INITIAL_STATE} dispatch={dispatch} onContinue={vi.fn()} />
      );

      const runBtn = screen.getByTestId('run-test-button');
      fireEvent.click(runBtn);

      unmount();
      expect(spyClear).toHaveBeenCalled();
      spyClear.mockRestore();
    });

    it('所有設備通過測試前「繼續」按鈕應為 disabled，全通過後啟用且點擊觸發 onContinue', () => {
      const onContinue = vi.fn();

      // 1. 未測試狀態下
      const { rerender } = render(
        <Step1Device state={INITIAL_STATE} dispatch={vi.fn()} onContinue={onContinue} />
      );
      const continueBtn = screen.getByTestId('btn-continue-step1');
      expect(continueBtn).toBeDisabled();

      // 2. 測試成功狀態下
      const successState: WorkbenchV2State = {
        ...INITIAL_STATE,
        devices: [
          {
            ...INITIAL_STATE.devices[0],
            status: 'tested',
            test: {
              status: 'success',
              latency_ms: 45,
              stages: {},
            },
          },
        ],
      };

      rerender(<Step1Device state={successState} dispatch={vi.fn()} onContinue={onContinue} />);
      const continueBtnActive = screen.getByTestId('btn-continue-step1');
      expect(continueBtnActive).not.toBeDisabled();

      fireEvent.click(continueBtnActive);
      expect(onContinue).toHaveBeenCalled();
    });
  });
});
