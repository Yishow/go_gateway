import { Download, Keyboard, Redo2, Undo2, Upload } from 'lucide-react';
import type { TFunction } from 'i18next';
import type { Device } from '../../../types/datalink';
import { QuickActions } from '../../../components/datalink/QuickActions';

interface SmartDashboardSidebarToolsProps {
  selectedDevice: Device | null;
  selectedAddressesCount: number;
  onBatchCreate: () => void;
  onOpenWorkbench: () => void;
  onOpenImport: () => void;
  onOpenExport: () => void;
  canExport: boolean;
  onUndo: () => void;
  onRedo: () => void;
  canUndo: boolean;
  canRedo: boolean;
  undoDescription: string;
  redoDescription: string;
  onOpenShortcuts: () => void;
  t: TFunction;
}

export default function SmartDashboardSidebarTools({
  selectedDevice,
  selectedAddressesCount,
  onBatchCreate,
  onOpenWorkbench,
  onOpenImport,
  onOpenExport,
  canExport,
  onUndo,
  onRedo,
  canUndo,
  canRedo,
  undoDescription,
  redoDescription,
  onOpenShortcuts,
  t,
}: SmartDashboardSidebarToolsProps) {
  return (
    <>
      <QuickActions
        device={selectedDevice}
        selectedCount={selectedAddressesCount}
        onBatchCreate={onBatchCreate}
        onQuickMapping={() => {}}
        onTestConnection={() => {}}
        onOpenWorkbench={onOpenWorkbench}
      />
      {selectedDevice && (
        <div className="border-t border-white/5 px-4 py-2">
          <div className="flex gap-2">
            <button
              onClick={onOpenImport}
              aria-keyshortcuts="Control+I"
              className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-slate-700/50 px-3 py-2 text-xs text-slate-300 transition-colors hover:bg-slate-800/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
            >
              <Upload className="h-4 w-4" />
              <span>{t('smartDashboard.import')}</span>
            </button>
            <button
              onClick={onOpenExport}
              disabled={!canExport}
              aria-keyshortcuts="Control+E"
              className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-slate-700/50 px-3 py-2 text-xs text-slate-300 transition-colors hover:bg-slate-800/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <Download className="h-4 w-4" />
              <span>{t('smartDashboard.export')}</span>
            </button>
          </div>
        </div>
      )}
      <div className="border-t border-white/5 px-4 py-2">
        <div className="flex gap-2">
          <button
            onClick={onUndo}
            disabled={!canUndo}
            aria-keyshortcuts="Control+Z"
            className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-slate-700/50 px-3 py-2 text-xs text-slate-300 transition-colors hover:bg-slate-800/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:cursor-not-allowed disabled:opacity-30"
            title={undoDescription}
          >
            <Undo2 className="h-4 w-4" />
            <span>{t('smartDashboard.undo')}</span>
          </button>
          <button
            onClick={onRedo}
            disabled={!canRedo}
            aria-keyshortcuts="Control+Y"
            className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-slate-700/50 px-3 py-2 text-xs text-slate-300 transition-colors hover:bg-slate-800/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:cursor-not-allowed disabled:opacity-30"
            title={redoDescription}
          >
            <Redo2 className="h-4 w-4" />
            <span>{t('smartDashboard.redo')}</span>
          </button>
        </div>
      </div>
      <div className="border-t border-white/5 p-4">
        <button
          onClick={onOpenShortcuts}
          className="flex w-full items-center justify-center gap-2 rounded-lg px-4 py-2 text-xs text-slate-300 transition-colors hover:bg-slate-800/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
        >
          <Keyboard className="h-4 w-4" />
          <span>{t('smartDashboard.shortcuts')}</span>
          <span className="ml-auto text-[10px] font-mono opacity-60">?</span>
        </button>
      </div>
    </>
  );
}
