import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { getRowGroupColumnConflicts } from '../../state/rowGroupValidation';
import type { Point, Mapping, DbRowGroup, DbTarget } from '../../state/types';
import type { DbColumn } from '../../state/dbSchemas';
import { TargetMappingRow } from './TargetMappingRow';

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

  // 過濾出有 mapping 定義且啟用的點位
  const visiblePoints = useMemo(() => {
    return points.filter(p => p.enabled && mappings[p.id]);
  }, [points, mappings]);

  const conflictingColumns = useMemo(() => {
    return getRowGroupColumnConflicts(visiblePoints, mappings, targets, rowGroups);
  }, [visiblePoints, mappings, targets, rowGroups]);

  // 是否有任何衝突
  const hasConflict = conflictingColumns.size > 0;

  return (
    <div className="space-y-4">
      {/* 標題與說明 */}
      <div>
        <h3 className="text-lg font-semibold text-white">
          {t('step4.mapping_title')}
        </h3>
        <p className="text-xs text-gray-500 mt-1">
          {t('step4.mapping_subtitle')}
        </p>
        <p className="mt-1 text-xs text-slate-400">
          {t('step4.mapping_location_hint')}
        </p>
      </div>

      {/* 表格區塊 */}
      <div className="border border-gray-800 rounded-xl overflow-hidden bg-gray-900/10">
        <div className="flex items-center justify-between gap-3 border-b border-gray-800 bg-gray-950/40 px-3 py-2">
          <div className="text-[11px] text-slate-500">
            {t('step4.mapping_count', {
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
              {t('step4.disable_all_targets')}
            </button>
            <button
              type="button"
              onClick={() => onSetAllEnabled?.(true)}
              disabled={disabled || visiblePoints.length === 0}
              className="rounded-md border border-blue-700/70 px-2.5 py-1 text-[11px] text-blue-300 transition-colors hover:bg-blue-950/40 disabled:cursor-not-allowed disabled:opacity-50"
              data-testid="btn-enable-all-db-targets"
            >
              {t('step4.enable_all_targets')}
            </button>
          </div>
        </div>
        <table className="w-full text-left border-collapse text-sm">
          <thead>
            <tr className="border-b border-gray-800 bg-gray-900/40 text-gray-400 font-medium text-xs select-none">
              <th className="p-3">{t('step4.col_tag_key')}</th>
              <th className="p-3">{t('step4.col_point')}</th>
              <th className="p-3 text-center w-12">→</th>
              <th className="p-3">{t('step4.col_db_field')}</th>
              <th className="p-3">{t('step4.col_field_type')}</th>
              <th className="p-3">{t('step4.col_row_group')}</th>
              <th className="p-3 text-right pr-6 w-24">{t('step4.col_enabled')}</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800/40">
            {visiblePoints.map(p => {
              const mapping = mappings[p.id];
              const target = targets[p.id] || { tag_id: `tag.${mapping.tag_key}`, column_name: '', enabled: true };

              // 是否處於衝突狀態 (已啟用且欄位名稱重複)
              const isConflicting = target.enabled && conflictingColumns.has(target.column_name);

              return (
                <TargetMappingRow
                  key={p.id}
                  point={p}
                  mapping={mapping}
                  target={target}
                  rowGroups={rowGroups}
                  columns={columns}
                  isConflicting={isConflicting}
                  onUpdateTarget={onUpdateTarget}
                  disabled={disabled}
                />
              );
            })}

            {visiblePoints.length === 0 && (
              <tr>
                <td colSpan={7} className="p-8 text-center text-gray-500 text-xs">
                  {t('step4.no_points')}
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
              {t('step4.conflict_banner_title')}
            </h4>
            <p className="mt-0.5 text-red-400/80">
              {t('step4.conflict_banner_desc')}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}
