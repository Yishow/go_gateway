import * as React from 'react';
import type { Point, Mapping, Device } from '../../state/types';
import { useDeviceColor } from '../../state/deviceColors';
import type { RuntimeStreamConnectionState } from '../../../../../types/datalink';

const MAPPING_SAVE_STATE_LABELS = {
  saving: '儲存中',
  saved: '已保存',
  'save-error': '保存失敗',
  'draft-invalid': '本地草稿',
} as const;

export interface MappingRowProps {
  point: Point;
  mapping: Mapping;
  isSelected: boolean;
  devices: Device[];
  liveValue?: unknown;
  connectionState?: RuntimeStreamConnectionState;
  onSelect: () => void;
  dispatch: React.Dispatch<any>;
}

function formatDeviceValue(value: unknown): string {
  if (value === null || value === undefined) {
    return '--';
  }
  if (typeof value === 'string' || typeof value === 'number' || typeof value === 'boolean') {
    return String(value);
  }
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}

/**
 * 點位映射表的單一資料列元件
 * 
 * 落地設計決策：「行內編輯 + 列選取共存：stopPropagation 策略」
 * 提供完整的行內編輯控件，點擊控件不觸發列選取。
 */
export const MappingRow: React.FC<MappingRowProps> = ({
  point,
  mapping,
  isSelected,
  devices,
  liveValue,
  connectionState = 'disconnected',
  onSelect,
  dispatch,
}) => {
  const devTheme = useDeviceColor(point.device_id);

  const devName = devices.find((d) => d.id === point.device_id)?.name || '未知裝置';

  // 阻止事件傳播，避免點擊輸入框時觸發整列選取
  const preventPropagation = (e: React.MouseEvent | React.TouchEvent) => {
    e.stopPropagation();
  };

  const handleTextChange = (field: keyof Mapping, val: string) => {
    dispatch({
      type: 'updateMapping',
      pointId: point.id,
      patch: { [field]: val },
    });
  };

  const handleNumChange = (field: keyof Mapping, val: string) => {
    const num = parseFloat(val);
    dispatch({
      type: 'updateMapping',
      pointId: point.id,
      patch: { [field]: isNaN(num) ? 0 : num },
    });
  };

  const deviceValueLabel =
    liveValue !== null && liveValue !== undefined
      ? formatDeviceValue(liveValue)
      : connectionState === 'connecting'
        ? '連線中'
        : connectionState === 'error'
          ? '串流錯誤'
          : '尚未收到';

  return (
    <tr
      onClick={onSelect}
      className={`transition-all duration-150 cursor-pointer ${
        isSelected
          ? 'bg-blue-500/5 outline outline-1 outline-blue-500/50'
          : 'border-b border-slate-900 hover:bg-slate-900/20'
      } ${!mapping.enabled ? 'opacity-40 bg-slate-950/10' : ''}`}
      data-testid={`mapping-row-${point.id}`}
    >
      {/* 1. 來源設備 */}
      <td className="p-3 text-xs">
        <div className="flex items-center gap-1.5" title={`${devName} · ${point.rule_name}`}>
          <span className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${devTheme.solid}`} />
          <span className={`font-sans truncate max-w-[80px] ${devTheme.text}`}>
            {devName}
          </span>
        </div>
      </td>

      {/* 2. 暫存器 (只讀) */}
      <td className="p-3 text-xs font-mono text-slate-400">
        {point.address}
      </td>

      {/* 3. 裝置即時值 */}
      <td className="p-3 text-xs font-mono text-slate-300">
        <span data-testid={`device-live-value-${point.id}`}>{deviceValueLabel}</span>
      </td>

      {/* 4. Tag Key (可編輯) */}
      <td className="p-2">
        <input
          type="text"
          value={mapping.tag_key}
          onChange={(e) => handleTextChange('tag_key', e.target.value)}
          onClick={preventPropagation}
          onMouseDown={preventPropagation}
          placeholder="temp.inlet"
          className="w-full bg-slate-900/50 hover:bg-slate-900 border border-slate-800 hover:border-slate-700 focus:border-blue-500 text-xs px-2 py-1 rounded text-slate-200 focus:outline-none transition-all font-mono"
          data-testid={`input-tag-key-${point.id}`}
        />
        {mapping.save_state && mapping.save_state !== 'idle' && (
          <div
            className={`mt-1 text-[10px] ${
              mapping.save_state === 'save-error'
                ? 'text-rose-400'
                : mapping.save_state === 'saving'
                  ? 'text-amber-400'
                  : mapping.save_state === 'draft-invalid'
                    ? 'text-amber-300'
                    : 'text-emerald-400'
            }`}
            data-testid={`mapping-save-state-${point.id}`}
          >
            {MAPPING_SAVE_STATE_LABELS[mapping.save_state]}
            {mapping.save_state === 'save-error' && mapping.save_error ? ` · ${mapping.save_error}` : ''}
          </div>
        )}
      </td>

      {/* 5. 顯示名稱 (可編輯) */}
      <td className="p-2">
        <input
          type="text"
          value={mapping.display_name}
          onChange={(e) => handleTextChange('display_name', e.target.value)}
          onClick={preventPropagation}
          onMouseDown={preventPropagation}
          placeholder="顯示名稱"
          className="w-full bg-slate-900/50 hover:bg-slate-900 border border-slate-800 hover:border-slate-700 focus:border-blue-500 text-xs px-2 py-1 rounded text-slate-200 focus:outline-none transition-all"
          data-testid={`input-display-name-${point.id}`}
        />
      </td>

      {/* 6. 單位 (可編輯) */}
      <td className="p-2">
        <input
          type="text"
          value={mapping.unit}
          onChange={(e) => handleTextChange('unit', e.target.value)}
          onClick={preventPropagation}
          onMouseDown={preventPropagation}
          placeholder="無"
          className="w-14 bg-slate-900/50 hover:bg-slate-900 border border-slate-800 hover:border-slate-700 focus:border-blue-500 text-xs px-2 py-1 rounded text-slate-200 focus:outline-none transition-all"
          data-testid={`input-unit-${point.id}`}
        />
      </td>

      {/* 7. 目標型態 (可編輯) */}
      <td className="p-2">
        <select
          value={mapping.target_type}
          onChange={(e) => handleTextChange('target_type', e.target.value)}
          onClick={preventPropagation}
          onMouseDown={preventPropagation}
          className="w-24 bg-slate-900/50 hover:bg-slate-900 border border-slate-800 hover:border-slate-700 focus:border-blue-500 text-xs px-1.5 py-1 rounded text-slate-200 focus:outline-none transition-all font-mono"
          data-testid={`select-target-type-${point.id}`}
        >
          <option value="bool">bool</option>
          <option value="int16">int16</option>
          <option value="int32">int32</option>
          <option value="int64">int64</option>
          <option value="uint16">uint16</option>
          <option value="uint32">uint32</option>
          <option value="uint64">uint64</option>
          <option value="float32">float32</option>
          <option value="float64">float64</option>
          <option value="string">string</option>
        </select>
      </td>

      {/* 8. Scale / Offset (可編輯) */}
      <td className="p-2">
        <div className="flex items-center gap-1">
          <input
            type="number"
            step="any"
            value={mapping.scale}
            onChange={(e) => handleNumChange('scale', e.target.value)}
            onClick={preventPropagation}
            onMouseDown={preventPropagation}
            className="w-14 bg-slate-900/50 hover:bg-slate-900 border border-slate-800 hover:border-slate-700 focus:border-blue-500 text-xs px-1 py-1 rounded text-slate-200 focus:outline-none transition-all font-mono"
            data-testid={`input-scale-${point.id}`}
          />
          <span className="text-slate-600 text-[10px] font-mono">+</span>
          <input
            type="number"
            step="any"
            value={mapping.offset}
            onChange={(e) => handleNumChange('offset', e.target.value)}
            onClick={preventPropagation}
            onMouseDown={preventPropagation}
            className="w-14 bg-slate-900/50 hover:bg-slate-900 border border-slate-800 hover:border-slate-700 focus:border-blue-500 text-xs px-1 py-1 rounded text-slate-200 focus:outline-none transition-all font-mono"
            data-testid={`input-offset-${point.id}`}
          />
        </div>
      </td>

      {/* 9. 啟用 Toggle */}
      <td className="p-3 text-center">
        <input
          type="checkbox"
          checked={mapping.enabled}
          onChange={() => dispatch({ type: 'toggleMappingEnabled', pointId: point.id })}
          onClick={preventPropagation}
          onMouseDown={preventPropagation}
          className="h-3.5 w-3.5 rounded border-slate-800 bg-slate-950 text-blue-600 focus:ring-blue-500 focus:ring-offset-slate-950 cursor-pointer"
          data-testid={`checkbox-enabled-${point.id}`}
        />
      </td>
    </tr>
  );
};
