import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Device } from '../../state/types';
import { getColorTheme, DEVICE_COLORS } from '../../state/deviceColors';
import { Icon } from '../../components/Icon';

export interface DeviceTabRailProps {
  /** 設備列表 */
  devices: Device[];
  /** 當前選中的設備 ID */
  selectedId: string;
  /** 當選中設備改變時的 callback */
  onSelect: (deviceId: string) => void;
  /** 新增設備的動作 callback */
  onAdd: () => void;
  /** 刪除設備的動作 callback */
  onDelete: (deviceId: string) => void;
  /** 重命名設備的動作 callback */
  onRename: (deviceId: string, name: string) => void;
}

/**
 * 裝置 Tab 標籤列元件 (DeviceTabRail)
 * 
 * 落地需求 **Multi-device tab management**。
 * 支持多設備的標籤切換、行內(inline)即時改名、顯示各設備的測試狀態，並提供防呆與連動清理的刪除操作。
 */
export const DeviceTabRail: React.FC<DeviceTabRailProps> = ({
  devices,
  selectedId,
  onSelect,
  onAdd,
  onDelete,
  onRename,
}) => {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="flex items-center gap-3 border-b border-slate-800 pb-3" data-testid="device-tab-rail">
      {/* Tab 清單區 (支援水平捲動) */}
      <div className="flex flex-1 items-center gap-2 overflow-x-auto pb-0.5 scrollbar-thin scrollbar-thumb-slate-800">
        {devices.map((device, index) => {
          const isSelected = device.id === selectedId;
          const color = getColorTheme(DEVICE_COLORS[index % 6]);

          // 刪除處理，包含原生 confirm 提示
          const handleDelete = (e: React.MouseEvent) => {
            e.stopPropagation();
            if (window.confirm(t('step1.dialogs.delete_confirm', { name: device.name }))) {
              onDelete(device.id);
            }
          };

          return (
            <div
              key={device.id}
              onClick={() => onSelect(device.id)}
              data-testid={`device-tab-${device.id}`}
              className={`group relative flex items-center gap-2.5 rounded-lg border px-3.5 py-2.5 cursor-pointer transition-all duration-200 ${
                isSelected
                  ? 'border-slate-700 bg-slate-900/60 shadow-inner'
                  : 'border-slate-800/80 bg-slate-950/20 hover:border-slate-800 hover:bg-slate-950/40'
              }`}
            >
              {/* 設備色彩指示點 */}
              <span
                className={`h-2.5 w-2.5 shrink-0 rounded-full transition-transform group-hover:scale-110 ${color.solid}`}
                data-testid={`device-color-dot-${device.id}`}
              />

              {/* Inline 改名 Input */}
              <input
                type="text"
                value={device.name}
                onChange={(e) => onRename(device.id, e.target.value)}
                onClick={(e) => e.stopPropagation()} // 防止 input 點擊引發 tab 重選
                className="w-28 bg-transparent border-none text-xs font-semibold text-slate-200 focus:outline-none focus:ring-0 p-0 placeholder-slate-600 transition-colors focus:text-white"
                data-testid={`input-rename-${device.id}`}
              />

              {/* 測試狀態 Chip */}
              {device.test?.status === 'success' && (
                <span
                  className="rounded-full bg-emerald-500/10 px-1.5 py-0.5 font-mono text-[9px] font-medium text-emerald-400 border border-emerald-500/10"
                  data-testid={`status-tested-${device.id}`}
                >
                  ✓ {device.test.latency_ms}ms
                </span>
              )}
              {device.test?.status === 'running' && (
                <span
                  className="flex items-center gap-1 rounded-full bg-blue-500/10 px-1.5 py-0.5 text-[9px] font-medium text-blue-400 border border-blue-500/10"
                  data-testid={`status-testing-${device.id}`}
                >
                  <Icon name="refresh" className="h-2.5 w-2.5 animate-spin" />
                  {t('step1.status.testing')}
                </span>
              )}
              {(!device.test || device.test.status === 'failed') && (
                <span
                  className="rounded-full bg-slate-800/60 px-1.5 py-0.5 text-[9px] font-medium text-slate-500"
                  data-testid={`status-untested-${device.id}`}
                >
                  {device.test?.status === 'failed' ? t('step1.status.failed') : t('step1.status.untested')}
                </span>
              )}

              {/* 刪除按鈕 (滑鼠移入時顯示，且至少有 2 個設備) */}
              {devices.length >= 2 && (
                <button
                  type="button"
                  onClick={handleDelete}
                  className="ml-1 opacity-0 group-hover:opacity-100 p-0.5 rounded-md hover:bg-slate-800 text-slate-500 hover:text-slate-200 transition-all"
                  data-testid={`btn-delete-device-${device.id}`}
                  title={t('step1.buttons.delete_device')}
                >
                  <Icon name="close" className="h-3 w-3 stroke-[2]" />
                </button>
              )}
            </div>
          );
        })}
      </div>

      {/* 新增設備按鈕 */}
      <button
        type="button"
        onClick={onAdd}
        className="flex items-center gap-1.5 rounded-lg border border-dashed border-blue-500/30 bg-blue-500/5 px-3.5 py-2.5 text-xs font-semibold text-blue-400 transition-all hover:bg-blue-500/10 hover:text-blue-300"
        data-testid="btn-add-device"
      >
        <Icon name="plus" className="h-3.5 w-3.5" />
        {t('step1.buttons.add_device')}
      </button>
    </div>
  );
};
