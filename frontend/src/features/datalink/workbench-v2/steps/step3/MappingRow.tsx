import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Mapping, Device } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { useDeviceColor } from '../../state/deviceColors';
import type { RuntimeStreamConnectionState, RuntimeStreamRecovery } from '../../../../../types/datalink';
import { MappingPreviewCells } from './MappingPreviewCells';
import { MappingSaveFailure } from './MappingSaveFailure';
import { MappingPayloadDialog } from './MappingPayloadDialog';

export interface MappingRowProps {
  point: Point;
  advanced?: boolean;
  mapping: Mapping;
  isSelected: boolean;
  bulkChecked?: boolean;
  onBulkToggle?: () => void;
  devices: Device[];
  liveValue?: unknown;
  connectionState?: RuntimeStreamConnectionState;
  streamRecovery?: RuntimeStreamRecovery;
  workspaceId?: string;
  onSelect: () => void;
  dispatch: React.Dispatch<WorkbenchV2Action>;
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
 * 提供完整行內編輯；聚焦編輯欄位同步來源列，批次對象勾選獨立。
 */
export const MappingRow: React.FC<MappingRowProps> = ({
  point,
  advanced = true,
  mapping,
  isSelected,
  devices,
  bulkChecked = false,
  onBulkToggle,
  liveValue,
  connectionState = 'disconnected',
  streamRecovery,
  workspaceId,
  onSelect,
  dispatch,
}) => {
  const { t } = useTranslation('workbench-v2');
  const devTheme = useDeviceColor(point.device_id);

  const devName = devices.find((d) => d.id === point.device_id)?.name || t('step3.unknownDevice');
  // 點擊不重複觸發列事件；聚焦編輯欄位仍同步來源。
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
      : t(`step3.connectionStates.${connectionState}`, {
        defaultValue: t('step3.connectionStates.disconnected'),
      });

  return (
    <tr
      onClick={onSelect}
      onFocus={onSelect}
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
          {onBulkToggle && <input type="checkbox" aria-label={t('step3.bulk.target', { tag: `${devName} ${point.address}` })} checked={bulkChecked} onChange={onBulkToggle} onFocus={e => e.stopPropagation()} onClick={e => e.stopPropagation()} data-testid={`bulk-target-${point.id}`} />}
          <span className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${devTheme.solid}`} />
          <span className={`font-sans truncate ${advanced ? 'max-w-[80px]' : 'max-w-[30vw]'} ${devTheme.text}`}>
            {devName}
          </span>
        </div>
        <span hidden={advanced} className="font-mono text-slate-400">{point.address}</span>
      </td>

      {/* 2. 暫存器 (只讀) */}
      <td hidden={!advanced} className="p-3 text-xs font-mono text-slate-400">
        {point.address}
      </td>

      {/* 4. 顯示名稱 (可編輯) */}
      <td className="min-w-48 p-2">
        <input
          type="text"
          value={mapping.display_name}
          onChange={(e) => handleTextChange('display_name', e.target.value)}
          onClick={preventPropagation}
          onMouseDown={preventPropagation}
          placeholder={t('step3.placeholders.displayName')}
          className="w-full bg-slate-900/50 hover:bg-slate-900 border border-slate-800 hover:border-slate-700 focus:border-blue-500 text-xs px-2 py-1 rounded text-slate-200 focus:outline-none transition-all"
          data-testid={`input-display-name-${point.id}`}
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
            {t(`step3.saveStates.${mapping.save_state}`)}
            {mapping.save_state === 'save-error' && <MappingSaveFailure pointId={point.id} mapping={mapping} dispatch={dispatch} />}
          </div>
        )}
      </td>

      {/* 5. 單位 (可編輯) */}
      <td className="p-2">
        <input
          type="text"
          value={mapping.unit}
          onChange={(e) => handleTextChange('unit', e.target.value)}
          onClick={preventPropagation}
          onMouseDown={preventPropagation}
          placeholder={t('step3.placeholders.unit')}
          className={`${advanced ? 'w-14' : 'w-full min-w-14'} bg-slate-900/50 hover:bg-slate-900 border border-slate-800 hover:border-slate-700 focus:border-blue-500 text-xs px-2 py-1 rounded text-slate-200 focus:outline-none transition-all`}
          data-testid={`input-unit-${point.id}`}
        />
      </td>

      {/* 6. 目標型態 (可編輯) */}
      <td hidden={!advanced} className="p-2">
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

      {/* 7. Scale / Offset (可編輯) */}
      <td hidden={!advanced} className="p-2">
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

      {/* 8. 讀值 (raw/decode) */}
      <td hidden={!advanced} className="p-3 text-xs font-mono text-slate-300">
        <span data-testid={`device-live-value-${point.id}`}>{deviceValueLabel}</span>
        {streamRecovery && (
          <div className="mt-1 space-y-0.5 font-sans text-[10px] text-amber-300" data-testid={`runtime-recovery-${point.id}`} role="alert">
            <div>{t(`errors.${streamRecovery.code ?? 'runtime_stream_unavailable'}`, {
              defaultValue: t('errors.runtime_stream_unavailable'),
            })}</div>
            <div>{t('step3.recovery.retryAction')}</div>
            {streamRecovery.requestId && (
              <div data-testid={`runtime-recovery-request-id-${point.id}`}>
                {t('errors.request_id')}: {streamRecovery.requestId}
              </div>
            )}
          </div>
        )}
      </td>

      <MappingPreviewCells
        advanced={advanced}
        point={point}
        mapping={mapping}
        rawValue={liveValue}
        connectionState={connectionState}
        workspaceId={workspaceId}
      />

      <MappingPayloadDialog hidden={!advanced} point={point} mapping={mapping} />

      {/* 9. 啟用 Toggle */}
      <td className="p-3 text-center">
        <input
          type="checkbox"
          checked={mapping.enabled}
          onChange={() => dispatch({ type: 'toggleMappingEnabled', pointId: point.id })}
          onFocus={(e) => e.stopPropagation()}
          onClick={preventPropagation}
          onMouseDown={preventPropagation}
          className="h-3.5 w-3.5 rounded border-slate-800 bg-slate-950 text-blue-600 focus:ring-blue-500 focus:ring-offset-slate-950 cursor-pointer"
          data-testid={`checkbox-enabled-${point.id}`}
        />
      </td>
    </tr>
  );
};
