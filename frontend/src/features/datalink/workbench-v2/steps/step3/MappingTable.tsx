import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Mapping, Device } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { MappingBulkApply } from './MappingBulkApply';
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
  const [advanced, setAdvanced] = React.useState(false);
  const [targets, setTargets] = React.useState(new Set<string>());

  return (
    <div className="flex flex-col gap-4 p-5 rounded-xl border border-slate-800 bg-slate-900/20" data-testid="mapping-table-wrapper">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="text-xs text-slate-500">
          {t('step3.table.enabledSummary', {
            defaultValue: '已載入 {{count}} 個映射點位',
            count: points.length,
          })}
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <button type="button" className="rounded-md border border-slate-800 px-2.5 py-1 text-[11px] text-slate-300 transition-colors hover:bg-slate-900" aria-expanded={advanced} data-testid="mapping-advanced-toggle" onClick={() => setAdvanced(value => !value)}>{t(advanced ? 'step3.bulk.basic' : 'step3.bulk.advanced')}</button>
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
      <div className="relative max-h-[520px] overflow-auto overscroll-y-contain border border-slate-800 rounded-lg bg-slate-950/40">
        <table className="w-full table-auto text-left border-collapse text-xs">
          <thead className="sticky top-0 bg-slate-900/90 backdrop-blur-sm z-10 border-b border-slate-800 text-slate-400 font-medium">
            <tr>
              <th className={`p-3 text-xs font-semibold text-slate-400 font-sans ${advanced ? 'w-24' : ''}`}>
                {t('step3.columns.source', { defaultValue: '來源' })}
              </th>
              <th hidden={!advanced} className="p-3 text-xs font-semibold text-slate-400 font-sans w-20">
                {t('step3.columns.register', { defaultValue: '暫存器' })}
              </th>
              <th className="min-w-48 p-3 text-xs font-semibold text-slate-400 font-sans">
                {t('step3.columns.displayName', { defaultValue: '顯示名稱' })}
              </th>
              <th className={`p-3 text-xs font-semibold text-slate-400 font-sans ${advanced ? 'w-16' : ''}`}>
                {t('step3.columns.unit', { defaultValue: '單位' })}
              </th>
              <th hidden={!advanced} className="p-3 text-xs font-semibold text-slate-400 font-sans w-28">
                {t('step3.columns.targetType', { defaultValue: '目標型態' })}
              </th>
              <th hidden={!advanced} className="p-3 text-xs font-semibold text-slate-400 font-sans w-36">
                {t('step3.columns.transform', { defaultValue: '× Scale + Offset' })}
              </th>
              <th hidden={!advanced} className="p-3 text-xs font-semibold text-slate-400 font-sans w-24">
                {t('step3.columns.deviceValue', { defaultValue: '讀值' })}
              </th>
              <th hidden={!advanced} className="p-3 text-xs font-semibold text-slate-400 font-sans w-24">
                {t('step3.columns.scaleResult', { defaultValue: '縮放值' })}
              </th>
              <th hidden={!advanced} className="p-3 text-xs font-semibold text-slate-400 font-sans w-24">
                {t('step3.columns.castResult', { defaultValue: '轉型值' })}
              </th>
              <th className={`p-3 text-xs font-semibold text-slate-400 font-sans ${advanced ? 'w-24' : ''}`}>
                {t('step3.columns.finalResult', { defaultValue: '最終值' })}
              </th>
              <th hidden={!advanced} className="p-3 text-xs font-semibold text-slate-400 font-sans w-20 text-center">
                {t('step3.columns.payload', { defaultValue: 'API payload' })}
              </th>
              <th className={`p-3 text-xs font-semibold text-slate-400 font-sans text-center ${advanced ? 'w-14' : ''}`}>
                {t('step3.columns.enabled', { defaultValue: '啟用' })}
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/60 text-slate-300">
            {points.length === 0 ? (
              <tr>
                <td colSpan={12} className="p-8 text-center text-slate-500 text-xs">
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
                    advanced={advanced}
                    point={p}
                    mapping={m}
                    isSelected={selectedIdx === idx}
                    devices={devices}
                    bulkChecked={targets.has(p.id)}
                    onBulkToggle={() => setTargets(previous => { const next = new Set(previous); if (next.has(p.id)) next.delete(p.id); else next.add(p.id); return next; })}
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

      <MappingBulkApply points={points} mappings={mappings} devices={devices} source={selectedPoint} targets={targets} setTargets={setTargets} dispatch={dispatch} />
    </div>
  );
};
