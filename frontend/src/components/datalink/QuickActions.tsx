

import type { Device } from '../../types/datalink';

// Icons
const Icons = {
  Plus: () => (
    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
    </svg>
  ),
  Zap: () => (
    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
    </svg>
  ),
  Link: () => (
    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
    </svg>
  )
};

export interface QuickActionsProps {
  device: Device | null;
  selectedCount: number;
  onBatchCreate: () => void;
  onQuickMapping: () => void;
  onTestConnection: () => void;
}

export function QuickActions({
  device,
  selectedCount,
  onBatchCreate,
  onQuickMapping,
  onTestConnection
}: QuickActionsProps) {
  
  if (!device) {
    return (
      <div className="bg-white dark:bg-slate-800 rounded-xl shadow-sm p-6 text-center" data-testid="quick-actions">
        <p className="text-slate-500 font-medium">請選擇設備以檢視操作</p>
      </div>
    );
  }

  return (
    <div className="space-y-6" data-testid="quick-actions">
      {/* 快速操作（設備狀態已合併至 SmartDashboard 標題列） */}
      <div className="bg-white dark:bg-slate-800 rounded-xl shadow-sm border border-slate-200 dark:border-slate-700 p-6 space-y-4">
        <h3 className="font-semibold text-slate-800 dark:text-slate-100 mb-4">快速操作</h3>
        
        <button
          onClick={onBatchCreate}
          disabled={selectedCount === 0}
          className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors font-medium shadow-sm"
        >
          <Icons.Plus />
          批量建立點位
          {selectedCount > 0 && <span className="bg-blue-500 px-2 py-0.5 rounded-full text-xs">{selectedCount}</span>}
        </button>

        <button
          onClick={onQuickMapping}
          disabled={selectedCount === 0}
          className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 disabled:opacity-50 disabled:cursor-not-allowed text-slate-700 dark:text-slate-200 rounded-lg transition-colors font-medium"
        >
          <Icons.Link />
          快速映射標籤
        </button>

        <button
          onClick={onTestConnection}
          className="w-full flex items-center justify-center gap-2 px-4 py-3 bg-slate-100 hover:bg-slate-200 dark:bg-slate-700 dark:hover:bg-slate-600 text-slate-700 dark:text-slate-200 rounded-lg transition-colors font-medium"
        >
          <Icons.Zap />
          測試連線
        </button>
      </div>
    </div>
  );
}
