import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Mapping, Device } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { MappingRow } from './MappingRow';
import type { RuntimeStreamConnectionState, RuntimeStreamRecovery } from '../../../../../types/datalink';

export interface MappingTableProps {
  points: Point[];
  mappings: Record<string, Mapping>;
  selectedIdx: number | null;
  setSelectedIdx: (idx: number | null) => void;
  devices: Device[];
  rawValues: Record<string, unknown>;
  connectionByDevice: Record<string, RuntimeStreamConnectionState>;
  recoveryByDevice?: Record<string, RuntimeStreamRecovery>;
  workspaceId?: string;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

/**
 * 點位映射表格殼層元件
 * 
 * 落地設計決策：「拆檔策略：6 個元件 + 2 個 state module」
 * 渲染 mappings 編輯表格，並於表尾提供批次套用線性轉換與型態的連結。
 */
export const MappingTable: React.FC<MappingTableProps> = ({
  points,
  mappings,
  selectedIdx,
  setSelectedIdx,
  devices,
  rawValues,
  connectionByDevice,
  recoveryByDevice = {},
  workspaceId,
  dispatch,
}) => {
  const { t } = useTranslation('workbench-v2');

  const selectedPoint = selectedIdx !== null ? points[selectedIdx] : null;
  const selectedMapping = selectedPoint ? mappings[selectedPoint.id] : null;
  const tagShort = selectedMapping ? (selectedMapping.tag_key || selectedPoint?.name || '') : '';

  const handleBulkApply = (fields: ('scale' | 'offset' | 'target_type' | 'unit')[]) => {
    if (!selectedPoint) return;
    dispatch({
      type: 'bulkApplyTransform',
      fromPointId: selectedPoint.id,
      fields,
    });
  };

  return (
    <div className="flex flex-col h-full bg-slate-950/40 border border-slate-900 rounded-lg overflow-hidden" data-testid="mapping-table-wrapper">
      <div className="flex items-center justify-between gap-3 border-b border-slate-900 bg-slate-950/40 px-3 py-2">
        <div className="text-[11px] text-slate-500">
          {t('step3.table.enabledSummary', {
            defaultValue: '已載入 {{count}} 個映射點位',
            count: points.length,
          })}
        </div>
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => dispatch({ type: 'setAllMappingsEnabled', enabled: false })}
            disabled={points.length === 0}
            className="rounded-md border border-slate-800 px-2.5 py-1 text-[11px] text-slate-300 transition-colors hover:bg-slate-900 disabled:cursor-not-allowed disabled:opacity-50"
            data-testid="btn-disable-all-mappings"
          >
            {t('step3.table.disableAll', { defaultValue: '全部停用' })}
          </button>
          <button
            type="button"
            onClick={() => dispatch({ type: 'setAllMappingsEnabled', enabled: true })}
            disabled={points.length === 0}
            className="rounded-md border border-blue-700/70 px-2.5 py-1 text-[11px] text-blue-300 transition-colors hover:bg-blue-950/40 disabled:cursor-not-allowed disabled:opacity-50"
            data-testid="btn-enable-all-mappings"
          >
            {t('step3.table.enableAll', { defaultValue: '全部啟用' })}
          </button>
        </div>
      </div>
      <div className="flex-1 overflow-auto max-h-[520px] scrollbar-thin">
        <table className="w-full border-collapse text-left">
          <thead>
            <tr className="border-b border-slate-900 bg-slate-950/60 sticky top-0 z-10">
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-24">
                {t('step3.columns.source', { defaultValue: '來源' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-20">
                {t('step3.columns.register', { defaultValue: '暫存器' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans">
                {t('step3.columns.tagKey', { defaultValue: 'Tag Key' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans">
                {t('step3.columns.displayName', { defaultValue: '顯示名稱' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-16">
                {t('step3.columns.unit', { defaultValue: '單位' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-28">
                {t('step3.columns.targetType', { defaultValue: '目標型態' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-36">
                {t('step3.columns.transform', { defaultValue: '× Scale + Offset' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-24">
                {t('step3.columns.deviceValue', { defaultValue: '讀值' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-24">
                {t('step3.columns.scaleResult', { defaultValue: '縮放值' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-24">
                {t('step3.columns.castResult', { defaultValue: '轉型值' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-24">
                {t('step3.columns.finalResult', { defaultValue: '最終值' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans w-20 text-center">
                {t('step3.columns.payload', { defaultValue: 'API payload' })}
              </th>
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans text-center w-14">
                {t('step3.columns.enabled', { defaultValue: '啟用' })}
              </th>
            </tr>
          </thead>
          <tbody>
            {points.length === 0 ? (
              <tr>
                <td colSpan={13} className="p-8 text-center text-slate-500 text-xs">
                  {t('step3.table.noPoints', { defaultValue: '無啟用的點位，請回上一步新增或啟用規則。' })}
                </td>
              </tr>
            ) : (
              points.map((p, idx) => {
                const m = mappings[p.id];
                if (!m) return null;
                return (
                  <MappingRow
                    key={p.id}
                    point={p}
                    mapping={m}
                    isSelected={selectedIdx === idx}
                    devices={devices}
                    liveValue={rawValues[p.id]}
                    connectionState={connectionByDevice[p.device_id] ?? 'disconnected'}
                    streamRecovery={recoveryByDevice[p.device_id]}
                    workspaceId={workspaceId}
                    onSelect={() => setSelectedIdx(idx)}
                    dispatch={dispatch}
                  />
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {/* 表尾批次套用 */}
      <div className="p-3 border-t border-slate-900 bg-slate-950/20 flex items-center justify-between gap-3">
        {selectedMapping ? (
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-[11px] text-slate-500 font-mono">
              {t('step3.table.bulkApplySource', {
                defaultValue: '來源：{{tag}}',
                tag: tagShort,
              })}
            </span>
            <button
              type="button"
              onClick={() => handleBulkApply(['scale', 'offset', 'target_type', 'unit'])}
              className="text-xs text-blue-400 hover:text-blue-300 font-semibold transition-colors"
              data-testid="btn-bulk-apply-all"
            >
              {t('step3.table.bulkApplyAll', {
                defaultValue: '套用轉換+單位到全部列',
                tag: tagShort,
              })}
            </button>
            <button
              type="button"
              onClick={() => handleBulkApply(['scale'])}
              className="text-xs text-slate-400 hover:text-slate-200 transition-colors"
              data-testid="btn-bulk-apply-scale"
            >
              {t('step3.table.bulkApplyScale', { defaultValue: 'Scale' })}
            </button>
            <button
              type="button"
              onClick={() => handleBulkApply(['offset'])}
              className="text-xs text-slate-400 hover:text-slate-200 transition-colors"
              data-testid="btn-bulk-apply-offset"
            >
              {t('step3.table.bulkApplyOffset', { defaultValue: 'Offset' })}
            </button>
            <button
              type="button"
              onClick={() => handleBulkApply(['target_type'])}
              className="text-xs text-slate-400 hover:text-slate-200 transition-colors"
              data-testid="btn-bulk-apply-target-type"
            >
              {t('step3.table.bulkApplyTargetType', { defaultValue: '型態' })}
            </button>
            <button
              type="button"
              onClick={() => handleBulkApply(['unit'])}
              className="text-xs text-slate-400 hover:text-slate-200 transition-colors"
              data-testid="btn-bulk-apply-unit"
            >
              {t('step3.table.bulkApplyUnit', { defaultValue: '單位' })}
            </button>
          </div>
        ) : (
          <span className="text-xs text-slate-600 flex items-center gap-1 cursor-not-allowed" data-testid="btn-bulk-apply-all-disabled">
            {t('step3.table.bulkApplyPlaceholder', {
              defaultValue: '點擊任一列後，才能套用轉換或單位到全部啟用列',
            })}
          </span>
        )}
        <span className="text-[10px] text-slate-600 font-mono">
          total {points.length}
        </span>
      </div>
    </div>
  );
};
