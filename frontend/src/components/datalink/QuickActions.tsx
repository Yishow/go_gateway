import { Plus, Zap, Link2, LayoutGrid } from 'lucide-react';
import type { Device } from '../../types/datalink';

export interface QuickActionsProps {
  device: Device | null;
  selectedCount: number;
  onBatchCreate: () => void;
  onQuickMapping: () => void;
  onTestConnection: () => void;
  onOpenWorkbench: () => void;
}

export function QuickActions({
  device,
  selectedCount,
  onBatchCreate,
  onQuickMapping,
  onTestConnection,
  onOpenWorkbench,
}: QuickActionsProps) {
  const isPrimaryOpsReady = device?.status === 'active';
  const prioritizedActions = isPrimaryOpsReady
    ? ['batch', 'mapping', 'test'] as const
    : ['test', 'batch', 'mapping'] as const;

  if (!device) {
    return (
      <div className="border-b border-white/5 px-4 py-3" data-testid="quick-actions">
        <p className="text-xs text-slate-500">請選擇設備以檢視操作</p>
      </div>
    );
  }

  return (
    <div className="border-b border-white/5 px-4 py-3" data-testid="quick-actions">
      <p className="mb-2 text-[11px] font-semibold uppercase tracking-wider text-slate-500">快速操作</p>
      <div className="space-y-1.5">
        {prioritizedActions.map((action) => {
          if (action === 'batch') {
            return (
              <button
                key={action}
                type="button"
                onClick={onBatchCreate}
                disabled={selectedCount === 0}
                className="flex w-full items-center gap-2 rounded-lg border border-blue-500/30 bg-blue-500/15 px-3 py-2 text-xs font-medium text-blue-100 transition-colors hover:bg-blue-500/25 disabled:cursor-not-allowed disabled:opacity-40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              >
                <Plus className="h-3.5 w-3.5 shrink-0" />
                <span>批量建立點位</span>
                {selectedCount > 0 && (
                  <span className="ml-auto rounded-full bg-blue-500/40 px-1.5 py-0.5 text-[10px] font-bold">
                    {selectedCount}
                  </span>
                )}
              </button>
            );
          }
          if (action === 'mapping') {
            return (
              <button
                key={action}
                type="button"
                onClick={onQuickMapping}
                disabled={selectedCount === 0}
                className="flex w-full items-center gap-2 rounded-lg border border-slate-700/50 bg-slate-800/60 px-3 py-2 text-xs text-slate-300 transition-colors hover:bg-slate-800 hover:text-white disabled:cursor-not-allowed disabled:opacity-40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              >
                <Link2 className="h-3.5 w-3.5 shrink-0" />
                <span>快速映射標籤</span>
              </button>
            );
          }
          return (
            <button
              key={action}
              type="button"
              onClick={onTestConnection}
              className="flex w-full items-center gap-2 rounded-lg border border-slate-700/50 bg-slate-800/60 px-3 py-2 text-xs text-slate-300 transition-colors hover:bg-slate-800 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              <Zap className="h-3.5 w-3.5 shrink-0" />
              <span>測試連線</span>
            </button>
          );
        })}
        <button
          type="button"
          onClick={onOpenWorkbench}
          className="flex w-full items-center gap-2 rounded-lg border border-cyan-500/30 bg-cyan-500/10 px-3 py-2 text-xs font-medium text-cyan-200 transition-colors hover:bg-cyan-500/20 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500"
        >
          <LayoutGrid className="h-3.5 w-3.5 shrink-0" />
          <span>Server Memory Grid</span>
        </button>
      </div>
    </div>
  );
}
