import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getRowGroupColumnConflicts } from '../../state/rowGroupValidation';
import type { Point, Mapping, DbRowGroup, DbTarget } from '../../state/types';
import type { DbColumn } from '../../state/dbSchemas';

/**
 * TargetMappingTable 元件屬性
 */
interface TargetMappingTableProps {
  points: Point[];
  mappings: Record<string, Mapping>;
  targets: Record<string, DbTarget>;
  rowGroups?: DbRowGroup[];
  columns: DbColumn[];
  onUpdateTarget: (pointId: string, patch: Partial<DbTarget>) => void;
  onSetAllEnabled?: (enabled: boolean) => void;
  disabled?: boolean;
}

/**
 * 點位至資料庫欄位映射配置表格
 * 落地需求：「Column conflict detection」之配置表格與衝突展示
 */
export function TargetMappingTable({
  points,
  mappings,
  targets,
  rowGroups = [],
  columns,
  onUpdateTarget,
  onSetAllEnabled,
  disabled = false
}: TargetMappingTableProps) {
  const { t } = useTranslation('workbench-v2');
  const [columnDrafts, setColumnDrafts] = useState<Record<string, string>>({});

  // 過濾出有 mapping 定義且啟用的點位
  const visiblePoints = useMemo(() => {
    return points.filter(p => p.enabled && mappings[p.id]);
  }, [points, mappings]);

  // 過濾排除 PK 的可用欄位
  const availableColumns = useMemo(() => {
    return columns.filter(c => !c.primary_key);
  }, [columns]);

  const conflictingColumns = useMemo(() => {
    return getRowGroupColumnConflicts(visiblePoints, mappings, targets, rowGroups);
  }, [visiblePoints, mappings, targets, rowGroups]);

  // 是否有任何衝突
  const hasConflict = conflictingColumns.size > 0;

  useEffect(() => {
    setColumnDrafts((current) => {
      const next: Record<string, string> = {};
      visiblePoints.forEach((point) => {
        next[point.id] = targets[point.id]?.column_name ?? '';
      });
      const currentKeys = Object.keys(current);
      const nextKeys = Object.keys(next);
      const isSame = currentKeys.length === nextKeys.length &&
        nextKeys.every((key) => current[key] === next[key]);
      return isSame ? current : next;
    });
  }, [targets, visiblePoints]);

  const commitColumnName = (pointId: string, currentValue: string, draftValue: string) => {
    const nextValue = draftValue.trim();
    if (!nextValue) {
      setColumnDrafts((current) => ({ ...current, [pointId]: currentValue }));
      return;
    }
    if (nextValue !== currentValue) {
      onUpdateTarget(pointId, { column_name: nextValue });
    }
  };

  return (
    <div className="space-y-4">
      {/* 標題與說明 */}
      <div>
        <h3 className="text-lg font-semibold text-white">
          {t('step4.mapping_title', '寫入欄位配置')}
        </h3>
        <p className="text-xs text-gray-500 mt-1">
          {t('step4.mapping_subtitle', { defaultValue: '將採集點位對應到目標資料表的欄位，主鍵欄位 (ts) 將自動填入系統時間' })}
        </p>
        <p className="mt-1 text-xs text-slate-400">
          {t('step4.mapping_location_hint', { defaultValue: '欄位清單來自下方資料庫設定的 schema/table；若目標表還沒準備好，請先用「資料表結構」預覽 DDL 或建立資料表。' })}
        </p>
      </div>

      {/* 表格區塊 */}
      <div className="border border-gray-800 rounded-xl overflow-hidden bg-gray-900/10">
        <div className="flex items-center justify-between gap-3 border-b border-gray-800 bg-gray-950/40 px-3 py-2">
          <div className="text-[11px] text-slate-500">
            {t('step4.mapping_count', {
              defaultValue: '已載入 {{count}} 個寫入欄位',
              count: visiblePoints.length,
            })}
            <span className="ml-1 font-mono text-slate-300" data-testid="step4-target-count">
              {visiblePoints.length}
            </span>
          </div>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => onSetAllEnabled?.(false)}
              disabled={disabled || visiblePoints.length === 0}
              className="rounded-md border border-slate-800 px-2.5 py-1 text-[11px] text-slate-300 transition-colors hover:bg-slate-900 disabled:cursor-not-allowed disabled:opacity-50"
              data-testid="btn-disable-all-db-targets"
            >
              {t('step4.disable_all_targets', { defaultValue: '全部停用' })}
            </button>
            <button
              type="button"
              onClick={() => onSetAllEnabled?.(true)}
              disabled={disabled || visiblePoints.length === 0}
              className="rounded-md border border-blue-700/70 px-2.5 py-1 text-[11px] text-blue-300 transition-colors hover:bg-blue-950/40 disabled:cursor-not-allowed disabled:opacity-50"
              data-testid="btn-enable-all-db-targets"
            >
              {t('step4.enable_all_targets', { defaultValue: '全部啟用' })}
            </button>
          </div>
        </div>
        <table className="w-full text-left border-collapse text-sm">
          <thead>
            <tr className="border-b border-gray-800 bg-gray-900/40 text-gray-400 font-medium text-xs select-none">
              <th className="p-3">{t('step4.col_tag_key', 'Tag 名稱')}</th>
              <th className="p-3">{t('step4.col_point', '來源點位')}</th>
              <th className="p-3 text-center w-12">→</th>
              <th className="p-3">{t('step4.col_db_field', '資料表欄位')}</th>
              <th className="p-3">{t('step4.col_field_type', '欄位型態')}</th>
              <th className="p-3">{t('step4.col_row_group', 'Row group')}</th>
              <th className="p-3 text-right pr-6 w-24">{t('step4.col_enabled', '啟用')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800/40">
            {visiblePoints.map(p => {
              const mapping = mappings[p.id];
              const target = targets[p.id] || { tag_id: `tag.${mapping.tag_key}`, column_name: '', enabled: true };
              const columnDraft = columnDrafts[p.id] ?? target.column_name;
              
              // 取得目前所選資料表欄位定義以顯示型別
              const currentColumn = columns.find(c => c.name === target.column_name);
              const columnType = currentColumn ? currentColumn.type : 'custom';
              
              // 是否處於衝突狀態 (已啟用且欄位名稱重複)
              const isConflicting = target.enabled && conflictingColumns.has(target.column_name);

              return (
                <tr
                  key={p.id}
                  className={`
                    transition-colors duration-150
                    ${isConflicting ? 'bg-red-500/5 hover:bg-red-500/10 border-l-2 border-l-red-500' : 'hover:bg-gray-800/20'}
                  `}
                >
                  {/* Tag 名稱 */}
                  <td className="p-3 font-mono text-xs text-gray-300">
                    {mapping.tag_key}
                  </td>

                  {/* 來源點位 */}
                  <td className="p-3 text-xs text-gray-400">
                    {p.name} <span className="text-gray-600 font-mono">({p.address})</span>
                  </td>

                  {/* 箭頭 */}
                  <td className="p-3 text-center text-gray-600">
                    →
                  </td>

                  {/* 資料表欄位選擇 */}
                  <td className="p-3">
                    <div className="flex items-center gap-2">
                      <input
                        type="text"
                        value={columnDraft}
                        list={`step4-target-column-options-${p.id}`}
                        data-testid={`step4-target-column-${p.id}`}
                        disabled={disabled || !target.enabled}
                        onChange={(e) => setColumnDrafts((current) => ({ ...current, [p.id]: e.target.value }))}
                        onBlur={() => commitColumnName(p.id, target.column_name, columnDraft)}
                        onKeyDown={(e) => {
                          if (e.key === 'Enter') {
                            e.preventDefault();
                            commitColumnName(p.id, target.column_name, columnDraft);
                            e.currentTarget.blur();
                          }
                        }}
                        className={`
                          w-44 bg-gray-950 border text-xs rounded-lg px-2 py-1.5 outline-none transition-all focus:ring-1 focus:ring-blue-500
                          ${isConflicting 
                            ? 'border-red-500/60 text-red-200' 
                            : 'border-gray-800 text-gray-300 disabled:text-gray-500'
                          }
                          disabled:opacity-50 disabled:cursor-not-allowed
                        `}
                      />
                      <datalist id={`step4-target-column-options-${p.id}`}>
                        {availableColumns.map(c => (
                          <option key={c.name} value={c.name} />
                        ))}
                      </datalist>
                      {isConflicting && (
                        <span className="text-red-400 text-xs flex items-center gap-1 select-none" title={t('step4.conflict_tooltip', '欄位與其他點位重複')}>
                          ⚠️
                        </span>
                      )}
                    </div>
                  </td>

                  {/* 欄位型態 */}
                  <td className="p-3 text-xs text-gray-500 font-mono">
                    {target.enabled ? columnType : '—'}
                  </td>

                  {/* Row group */}
                  <td className="p-3">
                    <select
                      value={target.row_group_id ?? ''}
                      disabled={disabled || !target.enabled || rowGroups.length === 0}
                      data-testid={`step4-target-row-group-${p.id}`}
                      onChange={(e) => onUpdateTarget(p.id, { row_group_id: e.target.value || undefined })}
                      className="w-40 rounded-lg border border-gray-800 bg-gray-950 px-2 py-1.5 text-xs text-gray-300 outline-none transition-all focus:ring-1 focus:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      <option value="">-</option>
                      {rowGroups.map((group) => (
                        <option key={group.id} value={group.id}>
                          {group.id}
                        </option>
                      ))}
                    </select>
                  </td>

                  {/* 啟用 Toggle */}
                  <td className="p-3 text-right pr-6">
                    <input
                      type="checkbox"
                      checked={target.enabled}
                      disabled={disabled}
                      onChange={() => onUpdateTarget(p.id, { enabled: !target.enabled })}
                      className="h-3.5 w-3.5 rounded border-slate-800 bg-slate-950 text-blue-600 focus:ring-blue-500 focus:ring-offset-slate-950 cursor-pointer disabled:cursor-not-allowed disabled:opacity-50"
                    />
                  </td>
                </tr>
              );
            })}

            {visiblePoints.length === 0 && (
              <tr>
                <td colSpan={7} className="p-8 text-center text-gray-500 text-xs">
                  {t('step4.no_points', '尚無已配置的映射點位，請先完成 Step 3')}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* 衝突警告 Banner */}
      {hasConflict && (
        <div className="flex items-center gap-3 p-4 rounded-xl border border-red-500/20 bg-red-500/5 text-red-400 text-xs">
          <span className="text-base select-none">🚨</span>
          <div className="flex-1">
            <h4 className="font-semibold text-red-200">
              {t('step4.conflict_banner_title', '欄位配置重複')}
            </h4>
            <p className="mt-0.5 text-red-400/80">
              {t('step4.conflict_banner_desc', '有複數點位映射至同一個資料表欄位。請調整下拉選單以消除紅框衝突，否則無法提交設定。')}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
