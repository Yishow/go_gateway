import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Mapping, Device } from '../../state/types';
import { MappingRow } from './MappingRow';

export interface MappingTableProps {
  points: Point[];
  mappings: Record<string, Mapping>;
  selectedIdx: number | null;
  setSelectedIdx: (idx: number | null) => void;
  devices: Device[];
  dispatch: React.Dispatch<any>;
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
  dispatch,
}) => {
  const { t } = useTranslation('workbench-v2');

  const selectedPoint = selectedIdx !== null ? points[selectedIdx] : null;
  const selectedMapping = selectedPoint ? mappings[selectedPoint.id] : null;
  const tagShort = selectedMapping ? (selectedMapping.tag_key || selectedPoint?.name || '') : '';

  const handleBulkApply = () => {
    if (!selectedPoint) return;
    dispatch({
      type: 'bulkApplyTransform',
      fromPointId: selectedPoint.id,
      fields: ['scale', 'offset', 'target_type'],
    });
  };

  return (
    <div className="flex flex-col h-full bg-slate-950/40 border border-slate-900 rounded-lg overflow-hidden" data-testid="mapping-table-wrapper">
      <div className="flex-1 overflow-y-auto max-h-[500px] scrollbar-thin">
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
              <th className="p-3 text-xs font-semibold text-slate-400 font-sans text-center w-14">
                {t('step3.columns.enabled', { defaultValue: '啟用' })}
              </th>
            </tr>
          </thead>
          <tbody>
            {points.length === 0 ? (
              <tr>
                <td colSpan={8} className="p-8 text-center text-slate-500 text-xs">
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
      <div className="p-3 border-t border-slate-900 bg-slate-950/20 flex items-center justify-between">
        {selectedMapping ? (
          <button
            type="button"
            onClick={handleBulkApply}
            className="text-xs text-blue-400 hover:text-blue-300 font-semibold transition-colors flex items-center gap-1 focus:outline-none"
            data-testid="btn-bulk-apply"
          >
            <span>⚡</span>
            {t('step3.table.bulkApplyText', {
              defaultValue: '將此轉換套用至全部：使用 {{tag}} 的 Scale / Offset / 型態',
              tag: tagShort,
            })}
          </button>
        ) : (
          <span className="text-xs text-slate-600 flex items-center gap-1 cursor-not-allowed" data-testid="btn-bulk-apply-disabled">
            <span>⚡</span>
            {t('step3.table.bulkApplyPlaceholder', {
              defaultValue: '點擊任一列以將該列的 Scale / Offset / 型態套用到全部點位',
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
