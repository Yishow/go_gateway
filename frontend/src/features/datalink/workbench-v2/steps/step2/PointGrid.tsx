import * as React from 'react';
import { useState, useEffect } from 'react';
import type { Point, ShareLayout } from '../../state/types';
import { Icon } from '../../components';

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
 * 自適應 4/6/8 欄位 Grid。雙向綁定修飾鍵以進行多選或批次略過操作。
 * 視覺上，略過點位呈現半透明，衝突點位標示紅圈與 alert 圖示，轉發點位在右下角標註 Share 地址。
 */
export const PointGrid: React.FC<PointGridProps> = ({
  ruleId,
  points,
  conflictAddrs,
  shareLayout,
  gridSelection,
  setGridSelection,
  onToggleSkipAddress,
  onBatchToggleSkipAddresses,
}) => {
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
      // Shift + Click: 連續區段批次設定略過/啟用
      const start = Math.min(lastClickedIdx, idx);
      const end = Math.max(lastClickedIdx, idx);
      const targetPoints = points.slice(start, end + 1);
      const targetAddresses = targetPoints.map((p) => p.address);
      const shouldSkip = !point.skipped;

      onBatchToggleSkipAddresses(targetAddresses, shouldSkip);
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
      // 一般 Click: 切換單一略過狀態
      onToggleSkipAddress(point.address);
      setLastClickedIdx(idx);
    }
  };

  return (
    <div
      className="grid grid-cols-4 sm:grid-cols-6 lg:grid-cols-8 gap-3"
      data-testid="point-grid"
    >
      {points.map((point, idx) => {
        const isSelected = gridSelection.has(point.id);
        const isConflict = conflictAddrs.has(point.address);
        const shareAddr = shareAddresses[idx];

        return (
          <div
            key={point.id}
            onClick={(e) => handleCellClick(e, point, idx)}
            className={`relative flex flex-col justify-between p-2.5 rounded-lg border transition-all select-none h-[72px] cursor-pointer ${
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
            <div className="flex items-center justify-between">
              <span className="font-mono text-xs font-semibold">{point.address}</span>
              {isConflict && (
                <span className="text-red-400" title="此地址與其他規則衝突">
                  <Icon name="alert" className="w-3.5 h-3.5" />
                </span>
              )}
            </div>

            {/* Bottom row: Name & share address label */}
            <div className="flex items-end justify-between mt-1">
              <span className="text-[10px] text-slate-400 truncate max-w-[80%]">
                {point.name}
              </span>

              {shareAddr !== null && (
                <span
                  className="absolute bottom-1 right-2 text-[9px] font-mono text-blue-400/80"
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
