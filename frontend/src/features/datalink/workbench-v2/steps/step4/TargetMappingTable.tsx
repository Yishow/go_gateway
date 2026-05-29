import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Mapping, DbTarget } from '../../state/types';
import type { DbColumn } from '../../state/dbSchemas';

/**
 * TargetMappingTable 元件屬性
 */
interface TargetMappingTableProps {
  points: Point[];
  mappings: Record<string, Mapping>;
  targets: Record<string, DbTarget>;
  columns: DbColumn[];
  onUpdateTarget: (pointId: string, patch: Partial<DbTarget>) => void;
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
  columns,
  onUpdateTarget,
  disabled = false
}: TargetMappingTableProps) {
  const { t } = useTranslation('workbench-v2');

  // 過濾出有 mapping 定義且啟用的點位
  const visiblePoints = useMemo(() => {
    return points.filter(p => p.enabled && mappings[p.id]);
  }, [points, mappings]);

  // 過濾排除 PK 的可用欄位
  const availableColumns = useMemo(() => {
    return columns.filter(c => !c.primary_key);
  }, [columns]);

  // 衝突檢測計算：只針對 visiblePoints 且 target 啟用的目標進行重複計數
  const conflictingColumns = useMemo(() => {
    const counts: Record<string, number> = {};
    visiblePoints.forEach(p => {
      const target = targets[p.id];
      if (target && target.enabled) {
        counts[target.column_name] = (counts[target.column_name] || 0) + 1;
      }
    });

    const conflicts = new Set<string>();
    Object.entries(counts).forEach(([col, count]) => {
      if (count > 1) {
        conflicts.add(col);
      }
    });
    return conflicts;
  }, [visiblePoints, targets]);

  // 是否有任何衝突
  const hasConflict = conflictingColumns.size > 0;

  return (
    <div className="space-y-4">
      {/* 標題與說明 */}
      <div>
        <h3 className="text-lg font-semibold text-white">
          {t('step4.mapping_title', '寫入欄位配置')}
        </h3>
        <p className="text-xs text-gray-500 mt-1">
          {t('step4.mapping_subtitle', '將採集點位對應到目標資料表的欄位，主鍵欄位 (ts) 將自動填入系統時間')}
        </p>
      </div>

      {/* 表格區塊 */}
      <div className="border border-gray-800 rounded-xl overflow-hidden bg-gray-900/10">
        <table className="w-full text-left border-collapse text-sm">
          <thead>
            <tr className="border-b border-gray-800 bg-gray-900/40 text-gray-400 font-medium text-xs select-none">
              <th className="p-3">{t('step4.col_tag_key', 'Tag 名稱')}</th>
              <th className="p-3">{t('step4.col_point', '來源點位')}</th>
              <th className="p-3 text-center w-12">→</th>
              <th className="p-3">{t('step4.col_db_field', '資料表欄位')}</th>
              <th className="p-3">{t('step4.col_field_type', '欄位型態')}</th>
              <th className="p-3 text-right pr-6 w-24">{t('step4.col_enabled', '啟用')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800/40">
            {visiblePoints.map(p => {
              const mapping = mappings[p.id];
              const target = targets[p.id] || { tag_id: `tag.${mapping.tag_key}`, column_name: '', enabled: true };
              
              // 取得目前所選資料表欄位定義以顯示型別
              const currentColumn = columns.find(c => c.name === target.column_name);
              const columnType = currentColumn ? currentColumn.type : '—';
              
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
                      <select
                        value={target.column_name}
                        disabled={disabled || !target.enabled}
                        onChange={(e) => onUpdateTarget(p.id, { column_name: e.target.value })}
                        className={`
                          bg-gray-950 border text-xs rounded-lg px-2 py-1.5 outline-none transition-all focus:ring-1 focus:ring-blue-500
                          ${isConflicting 
                            ? 'border-red-500/60 text-red-200' 
                            : 'border-gray-800 text-gray-300 disabled:text-gray-500'
                          }
                          disabled:opacity-50 disabled:cursor-not-allowed
                        `}
                      >
                        {availableColumns.map(c => (
                          <option key={c.name} value={c.name}>
                            {c.name}
                          </option>
                        ))}
                      </select>
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

                  {/* 啟用 Toggle */}
                  <td className="p-3 text-right pr-6">
                    <button
                      type="button"
                      disabled={disabled}
                      onClick={() => onUpdateTarget(p.id, { enabled: !target.enabled })}
                      className={`
                        relative inline-flex h-5 w-9 items-center rounded-full transition-colors duration-200 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed
                        ${target.enabled ? 'bg-blue-600' : 'bg-gray-800'}
                      `}
                    >
                      <span
                        className={`
                          inline-block h-3.5 w-3.5 transform rounded-full bg-white transition-transform duration-200
                          ${target.enabled ? 'translate-x-4.5' : 'translate-x-1'}
                        `}
                      />
                    </button>
                  </td>
                </tr>
              );
            })}

            {visiblePoints.length === 0 && (
              <tr>
                <td colSpan={6} className="p-8 text-center text-gray-500 text-xs">
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
