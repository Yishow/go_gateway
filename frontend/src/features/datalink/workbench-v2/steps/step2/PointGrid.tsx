import * as React from 'react';
import { useState, useEffect } from 'react';
import type { Point, ShareLayout } from '../../state/types';
import { useTranslation } from 'react-i18next';
import { Icon } from '../../components';
import { pointConflictKey } from '../../state/sourceRule';

export interface PointGridProps {
  ruleId: string;
  points: Point[];
  conflictAddrs: Set<string>;
  shareLayout: ShareLayout | null;
  gridSelection: Set<string>;
  setGridSelection: React.Dispatch<React.SetStateAction<Set<string>>>;
  onToggleSkipAddress: (address: string) => void;
  onBatchToggleSkipAddresses: (addresses: string[], shouldSkip: boolean) => void;
}

/**
 * 點位互動網格元件
 * 
 * 落地設計決策：「點位網格互動：點擊 + Shift + Ctrl/⌘ 三模式」與「衝突偵測」
 * 依容器可用寬度排欄。點擊與修飾鍵只選取；明確勾選才修改採集。
 * 略過點位呈現半透明，衝突點位標示紅圈與 alert 圖示，轉發地址保留在內容流內。
 */
export const PointGrid: React.FC<PointGridProps> = ({
  ruleId,
  points,
  conflictAddrs,
  shareLayout,
  gridSelection,
  setGridSelection,
  onToggleSkipAddress,
}) => {
  const { t } = useTranslation('workbench-v2');
  const [lastClickedIdx, setLastClickedIdx] = useState<number | null>(null);

  // 切換規則時，重置選取狀態與最後點選索引
  useEffect(() => {
    setGridSelection(new Set());
    setLastClickedIdx(null);
  }, [ruleId, setGridSelection]);

  // 事先計算每個啟用點位的 Modbus Share 地址
  const shareAddresses: (number | null)[] = [];
  let enabledCount = 0;
  points.forEach((p) => {
    if (p.skipped || !p.enabled || !shareLayout) {
      shareAddresses.push(null);
    } else {
      const addr = shareLayout.start + enabledCount * shareLayout.stride;
      shareAddresses.push(addr);
      enabledCount++;
    }
  });

  const handleCellClick = (e: React.MouseEvent, point: Point, idx: number) => {
    e.preventDefault();

    if (e.shiftKey && lastClickedIdx !== null) {
      const start = Math.min(lastClickedIdx, idx);
      const end = Math.max(lastClickedIdx, idx);
      setGridSelection(new Set(points.slice(start, end + 1).map(p => p.id)));
      setLastClickedIdx(idx);
    } else if (e.metaKey || e.ctrlKey) {
      // Ctrl/Cmd + Click: 多選 toggle
      const next = new Set(gridSelection);
      if (next.has(point.id)) {
        next.delete(point.id);
      } else {
        next.add(point.id);
      }
      setGridSelection(next);
      setLastClickedIdx(idx);
    } else {
      // 一般點擊只選取，採集由明確勾選控制。
      setGridSelection(new Set([point.id]));
      setLastClickedIdx(idx);
    }
  };

  return (
    <div
      className="grid min-w-0 grid-cols-[repeat(auto-fill,minmax(min(100%,8.5rem),1fr))] gap-3"
      data-testid="point-grid"
    >
      {points.map((point, idx) => {
        const isSelected = gridSelection.has(point.id);
        const isConflict = conflictAddrs.has(pointConflictKey(point));
        const shareAddr = shareAddresses[idx];

        return (
          <div
            key={point.id}
            onClick={(e) => handleCellClick(e, point, idx)}
            role="button"
            tabIndex={0}
            aria-pressed={isSelected}
            onKeyDown={(e) => { if (e.target === e.currentTarget && (e.key === 'Enter' || e.key === ' ')) {
              e.preventDefault(); setGridSelection(new Set([point.id])); setLastClickedIdx(idx);
            } }}
            className={`relative flex min-w-0 flex-col justify-between gap-1 p-2.5 rounded-lg border transition-all select-none min-h-[72px] cursor-pointer ${
              point.skipped
                ? 'bg-slate-950/20 border-slate-900 text-slate-600 opacity-40 hover:opacity-60'
                : isConflict
                ? 'bg-red-500/5 border-red-500 text-red-200 hover:bg-red-500/10'
                : 'bg-slate-950/40 border-slate-800 text-slate-100 hover:border-slate-700/60'
            } ${
              isSelected
                ? 'ring-2 ring-blue-500/80 border-transparent'
                : isConflict
                ? 'ring-1 ring-red-500/50'
                : ''
            }`}
            data-testid={`point-cell-${point.id}`}
            data-address={point.address}
            data-skipped={point.skipped}
            data-selected={isSelected}
            data-conflict={isConflict}
          >
            {/* Top row: Address & Conflict Badge */}
            <div className="flex min-w-0 flex-wrap items-center justify-between gap-x-2 gap-y-1">
              <div className="flex min-w-0 items-center gap-1">
                <span className="break-all font-mono text-xs font-semibold">{point.address}</span>
                {isConflict && (
                  <span className="shrink-0 text-red-400" title={t('step2.table.conflict')}>
                    <Icon name="alert" className="w-3.5 h-3.5" />
                  </span>
                )}
              </div>
              <label
                className="inline-flex min-h-8 min-w-8 shrink-0 cursor-pointer items-center justify-center gap-1 whitespace-nowrap text-[10px] [@media(pointer:coarse)]:min-h-11 [@media(pointer:coarse)]:min-w-11"
                onClick={e => e.stopPropagation()}
              >
                <input className="shrink-0" type="checkbox" checked={!point.skipped} onChange={() => onToggleSkipAddress(point.address)} data-testid={`collect-${point.id}`} />
                {t('step2.grid.collect')}
              </label>
            </div>

            {/* Bottom row: Name & share address label */}
            <div className="flex min-w-0 flex-wrap items-end justify-between gap-x-2 gap-y-1">
              <span className="min-w-0 flex-1 truncate text-[10px] text-slate-400" title={point.name}>
                {point.name}
              </span>

              {shareAddr !== null && (
                <span
                  className="shrink-0 text-[9px] font-mono text-blue-400/80"
                  data-testid={`share-addr-${point.id}`}
                >
                  →{shareAddr}
                </span>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
};
