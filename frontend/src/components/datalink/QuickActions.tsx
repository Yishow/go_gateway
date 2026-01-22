

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
  Activity: () => (
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

  // Mock status for now
  const isConnected = device.status === 'active';

  return (
    <div className="space-y-6" data-testid="quick-actions">
      {/* Device Status Card */}
      <div className="bg-white dark:bg-slate-800 rounded-xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold text-slate-800 dark:text-slate-100">設備狀態</h3>
          <div className={`flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium ${
            isConnected 
              ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300' 
              : 'bg-slate-100 text-slate-600 dark:bg-slate-700 dark:text-slate-300'
          }`}>
            <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500 animate-pulse' : 'bg-slate-400'}`}></span>
            {isConnected ? '已連線' : '未連線'}
          </div>
        </div>
        
        <div className="space-y-3">
          <div className="flex justify-between text-sm">
            <span className="text-slate-500 dark:text-slate-400">協議</span>
            <span className="font-medium text-slate-700 dark:text-slate-200 font-mono">{device.protocol}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-slate-500 dark:text-slate-400">最後同步</span>
            <span className="font-medium text-slate-700 dark:text-slate-200">
                {device.last_test_at ? new Date(device.last_test_at).toLocaleTimeString() : '-'}
            </span>
          </div>
        </div>
      </div>

      {/* Actions */}
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
