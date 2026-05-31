import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WorkbenchV2State, Device, ProtocolId } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { useTestDraftConnectionMutation } from '../../../../../hooks/datalink/useDevices';
import { DeviceTabRail } from './DeviceTabRail';
import { DeviceEditor } from './DeviceEditor';
import { ConnectionTestPanel } from './ConnectionTestPanel';
import { DeviceListContext } from '../../state/deviceColors';
import { Icon } from '../../components/Icon';
import { Button } from '../../components/Button';

export interface Step1DeviceProps {
  /** Workbench V2 全域狀態 */
  state: WorkbenchV2State;
  /** Dispatch 函數 */
  dispatch: React.Dispatch<WorkbenchV2Action>;
  /** 當前 Step 1 完成，準備進行下一步的 callback 函數 */
  onContinue: () => void;
}

/**
 * Workbench V2 第一步：裝置設定主元件 (Step1Device)
 * 
 * 落地設計決策：「Step 1 diagnostics 直接使用 backend test-draft 契約」
 * 實作需求 **Multi-device tab management**、**Readiness check with explicit connect / probe separation**、**Continue gate**。
 * 容納多裝置列表、設備參數編輯、連線測試進度以及流轉控制 footer。
 */
export const Step1Device: React.FC<Step1DeviceProps> = ({ state, dispatch, onContinue }) => {
  const { t } = useTranslation('workbench-v2');
  const testDraftConnectionMutation = useTestDraftConnectionMutation();

  // 1. 本地所選設備狀態
  const [selectedId, setSelectedId] = React.useState<string>(() => {
    return state.devices[0]?.id || '';
  });

  // 當設備刪除時，自動修正選中狀態
  React.useEffect(() => {
    if (state.devices.length > 0) {
      const exists = state.devices.some((d) => d.id === selectedId);
      if (!exists) {
        setSelectedId(state.devices[0].id);
      }
    } else {
      setSelectedId('');
    }
  }, [state.devices, selectedId]);

  // 當前選中的設備
  const activeDevice = React.useMemo(() => {
    return state.devices.find((d) => d.id === selectedId) || null;
  }, [state.devices, selectedId]);

  // 3. 使用真實 backend diagnostics 測試目前 draft
  const handleRunTest = async (deviceId: string) => {
    const dev = state.devices.find((d) => d.id === deviceId);
    if (!dev) return;

    dispatch({ type: 'startDeviceTest', deviceId });
    try {
      const result = await testDraftConnectionMutation.mutateAsync({
        protocol: dev.protocol,
        connection_config: dev.config,
      });
      dispatch({ type: 'resolveDeviceTest', deviceId, result });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'diagnostics request failed';
      dispatch({ type: 'failDeviceTest', deviceId, stageId: 'connect', message });
    }
  };

  // 4. 新增設備
  const handleAddDevice = () => {
    const newIndex = state.devices.length + 1;
    const newId = `dev-${Date.now()}-${Math.floor(Math.random() * 1000)}`;
    const newHost = newIndex === 2 ? '192.168.1.101' : `192.168.1.10${newIndex - 1}`;

    const newDevice: Device = {
      id: newId,
      name: `${t('step1.labels.device_prefix')} ${newIndex}`,
      description: '',
      protocol: 'modbus_tcp',
      config: {
        host: newHost,
        port: 502,
        slave_id: 1,
        timeout: 5,
      },
      status: 'draft',
      test: null,
      persisted: false,
      save_state: 'idle',
      save_error: null,
    };

    dispatch({ type: 'addDevice', device: newDevice });
    setSelectedId(newId);
  };

  // 5. 刪除設備
  const handleDeleteDevice = (deviceId: string) => {
    dispatch({ type: 'removeDevice', deviceId });
  };

  // 6. 更新設備基本資料
  const handleUpdateDevice = (patch: Partial<Device>) => {
    dispatch({ type: 'updateDevice', deviceId: selectedId, patch });
  };

  // 7. 更新配置
  const handleUpdateConfig = (patch: Record<string, any>) => {
    dispatch({ type: 'updateDeviceConfig', deviceId: selectedId, patch });
  };

  // 8. 更改協議
  const handleChangeProtocol = (protocol: ProtocolId) => {
    dispatch({ type: 'changeDeviceProtocol', deviceId: selectedId, protocol });
  };

  // 9. 重命名
  const handleRenameDevice = (deviceId: string, name: string) => {
    dispatch({ type: 'renameDevice', deviceId, name });
  };

  // 10. 繼續的閘口檢查 (Continue Gate)
  const totalCount = state.devices.length;
  const testedCount = state.devices.filter((d) => d.test?.status === 'success').length;
  const allTested = totalCount > 0 && testedCount === totalCount;

  return (
    <DeviceListContext.Provider value={state.devices}>
      <div className="flex min-h-[calc(100vh-140px)] flex-col justify-between p-6" data-testid="step1-device-view">
        <div className="space-y-6">
          {/* Tab 選擇列 */}
          <DeviceTabRail
            devices={state.devices}
            selectedId={selectedId}
            onSelect={setSelectedId}
            onAdd={handleAddDevice}
            onDelete={handleDeleteDevice}
            onRename={handleRenameDevice}
          />

          {/* 中層編輯區與測試診斷區 */}
          {activeDevice ? (
            <div className="grid grid-cols-1 gap-6 lg:grid-cols-12">
              {/* 左側：編輯表單 (Col 7) */}
              <div className="lg:col-span-7">
                <DeviceEditor
                  device={activeDevice}
                  onUpdate={handleUpdateDevice}
                  onUpdateConfig={handleUpdateConfig}
                  onChangeProtocol={handleChangeProtocol}
                />
              </div>

              {/* 右側：診斷測試 (Col 5) */}
              <div className="lg:col-span-5">
                <ConnectionTestPanel device={activeDevice} onRunTest={handleRunTest} />
              </div>
            </div>
          ) : (
            <div className="flex h-64 items-center justify-center rounded-xl border border-dashed border-slate-800 bg-slate-900/5 text-sm text-slate-500">
              {t('step1.hints.no_devices')}
            </div>
          )}
        </div>

        {/* 底部繼續按鈕列 (Continue Gate) */}
        <div className="mt-8 flex items-center justify-between border-t border-slate-800/80 pt-5">
          {/* 左側：統計與警告 chip */}
          <div className="flex items-center gap-3">
            <span className="text-xs text-slate-400">
              {t('step1.footer.progress_summary', { total: totalCount, tested: testedCount })}
            </span>
            {!allTested && (
              <div
                className="flex items-center gap-1 rounded bg-amber-500/10 px-2 py-1 text-[11px] font-medium text-amber-400 border border-amber-500/10"
                data-testid="warning-untested-chip"
              >
                <Icon name="alert" className="h-3 w-3" />
                {t('step1.footer.warning_untested')}
              </div>
            )}
          </div>

          {/* 右側：動作按鈕 */}
          <Button
            onClick={onContinue}
            disabled={!allTested}
            variant="primary"
            data-testid="btn-continue-step1"
          >
            <span className="flex items-center gap-1.5">
              {t('step1.buttons.continue')}
              <Icon name="arrow" className="h-3.5 w-3.5" />
            </span>
          </Button>
        </div>
      </div>
    </DeviceListContext.Provider>
  );
};
