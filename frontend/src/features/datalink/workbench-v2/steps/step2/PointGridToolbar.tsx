import * as React from 'react';

export interface PointGridToolbarProps {
  gridSelection: Set<string>;
  totalCount: number;
  stride: number;
  onAction: (
    actionType:
      | 'all_enabled'
      | 'all_skipped'
      | 'reverse'
      | 'skip_selected'
      | 'enable_selected'
      | 'clear_selection'
  ) => void;
}

/**
 * 點位網格批次控制工具列
 * 
 * 落地設計決策：「Point grid batch toolbar」
 * 渲染 6 個批次按鈕（包含啟用、略過、反轉及對選取區域的操作），
 * 並於右側呈現該規則的點位統計 (total points 與型別 width / stride)。
 */
export const PointGridToolbar: React.FC<PointGridToolbarProps> = ({
  gridSelection,
  totalCount,
  stride,
  onAction,
}) => {
  const selectionCount = gridSelection.size;
  const isSelectionEmpty = selectionCount === 0;

  return (
    <div
      className="flex flex-wrap items-center justify-between gap-3 p-3 rounded-lg bg-slate-900/40 border border-slate-800/80"
      data-testid="point-grid-toolbar"
    >
      {/* 批次操作按鈕群組 */}
      <div className="flex flex-wrap items-center gap-1.5">
        <button
          type="button"
          onClick={() => onAction('all_enabled')}
          className="border border-slate-800 bg-slate-900 hover:bg-slate-800 text-[11px] font-medium px-2.5 py-1 rounded-md text-slate-300 hover:text-slate-100 transition-colors"
          data-testid="btn-all-enabled"
        >
          全部啟用
        </button>
        <button
          type="button"
          onClick={() => onAction('all_skipped')}
          className="border border-slate-800 bg-slate-900 hover:bg-slate-800 text-[11px] font-medium px-2.5 py-1 rounded-md text-slate-300 hover:text-slate-100 transition-colors"
          data-testid="btn-all-skipped"
        >
          全部略過
        </button>
        <button
          type="button"
          onClick={() => onAction('reverse')}
          className="border border-slate-800 bg-slate-900 hover:bg-slate-800 text-[11px] font-medium px-2.5 py-1 rounded-md text-slate-300 hover:text-slate-100 transition-colors"
          data-testid="btn-reverse"
        >
          反轉啟用
        </button>
        <div className="h-4 w-px bg-slate-800" />
        <button
          type="button"
          onClick={() => onAction('skip_selected')}
          disabled={isSelectionEmpty}
          className="border border-slate-800 bg-slate-900 hover:bg-slate-800 text-[11px] font-medium px-2.5 py-1 rounded-md text-slate-300 hover:text-slate-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          data-testid="btn-skip-selected"
        >
          略過選取
        </button>
        <button
          type="button"
          onClick={() => onAction('enable_selected')}
          disabled={isSelectionEmpty}
          className="border border-slate-800 bg-slate-900 hover:bg-slate-800 text-[11px] font-medium px-2.5 py-1 rounded-md text-slate-300 hover:text-slate-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          data-testid="btn-enable-selected"
        >
          啟用選取
        </button>

        {/* 條件渲染：僅在有選取時才顯示清除選取按鈕 */}
        {!isSelectionEmpty && (
          <button
            type="button"
            onClick={() => onAction('clear_selection')}
            className="text-[11px] text-blue-400 hover:text-blue-300 font-medium px-2 py-1 ml-1"
            data-testid="btn-clear-selection"
          >
            清除選取 ({selectionCount})
          </button>
        )}
      </div>

      {/* 統計說明 */}
      <div className="text-[11px] font-mono text-slate-500" data-testid="grid-summary">
        total {totalCount} · stride {stride}
      </div>
    </div>
  );
};
