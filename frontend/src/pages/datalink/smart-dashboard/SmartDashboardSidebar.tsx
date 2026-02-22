import { Download, Keyboard, Redo2, Undo2, Upload, LayoutGrid, Plus } from 'lucide-react';
import type { TFunction } from 'i18next';
import type { Device } from '../../../types/datalink';
import SmartDashboardPlanningTab, {
  type SmartDashboardPlanningTabProps,
} from './SmartDashboardPlanningTab';
import SmartDashboardTagPanel, {
  type SmartDashboardTagPanelProps,
} from './SmartDashboardTagPanel';
import SmartDashboardModbusPanel, {
  type SmartDashboardModbusPanelProps,
} from './SmartDashboardModbusPanel';
import SmartDashboardCommitPanel, {
  type SmartDashboardCommitPanelProps,
} from './SmartDashboardCommitPanel';

export type SidebarTab = 'plan' | 'tag' | 'modbus' | 'commit';

interface MiniToolbarProps {
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

export interface SmartDashboardSidebarProps
  extends MiniToolbarProps {
  sidebarTab: SidebarTab;
  setSidebarTab: (tab: SidebarTab) => void;
  /** 選取格位非空時，Tag Tab 顯示通知點 */
  hasTagSelection: boolean;
  planningTabProps: SmartDashboardPlanningTabProps;
  tagPanelProps: SmartDashboardTagPanelProps;
  modbusPanelProps: SmartDashboardModbusPanelProps;
  commitPanelProps: SmartDashboardCommitPanelProps;
}

/** Tab 標籤定義 */
const TABS = [
  { id: 'plan', label: '規劃' },
  { id: 'tag', label: 'Tag' },
  { id: 'modbus', label: 'Modbus' },
  { id: 'commit', label: '提交' },
] as const;

/** 右側欄：Mini Toolbar + 三分頁（規劃 / Tag / 提交） */
export default function SmartDashboardSidebar({
  sidebarTab,
  setSidebarTab,
  hasTagSelection,
  planningTabProps,
  tagPanelProps,
  modbusPanelProps,
  commitPanelProps,
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
}: SmartDashboardSidebarProps) {
  return (
    <div className="flex h-full flex-col overflow-hidden">
      {/* ── Mini Toolbar ── */}
      <div className="flex shrink-0 flex-wrap items-center gap-1 border-b border-white/5 px-3 py-2">
        {selectedDevice && (
          <button
            type="button"
            onClick={onBatchCreate}
            disabled={selectedAddressesCount === 0}
            title={`批量建立點位${selectedAddressesCount > 0 ? ` (${selectedAddressesCount})` : ''}`}
            className="relative inline-flex h-8 w-8 items-center justify-center rounded-lg border border-blue-500/30 bg-blue-500/10 text-blue-200 hover:bg-blue-500/20 disabled:cursor-not-allowed disabled:opacity-40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            <Plus className="h-4 w-4" />
            {selectedAddressesCount > 0 && (
              <span className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-blue-500 text-[9px] font-bold text-white">
                {selectedAddressesCount > 9 ? '9+' : selectedAddressesCount}
              </span>
            )}
          </button>
        )}
        <button
          type="button"
          onClick={onOpenWorkbench}
          title="Server Memory Grid"
          className="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-slate-700/50 text-slate-400 hover:bg-slate-800/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
        >
          <LayoutGrid className="h-4 w-4" />
        </button>
        <div className="mx-1 h-4 w-px bg-white/10" />
        {selectedDevice && (
          <>
            <button
              type="button"
              onClick={onOpenImport}
              title={t('smartDashboard.import')}
              aria-keyshortcuts="Control+I"
              className="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-slate-700/50 text-slate-400 hover:bg-slate-800/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              <Upload className="h-4 w-4" />
            </button>
            <button
              type="button"
              onClick={onOpenExport}
              disabled={!canExport}
              title={t('smartDashboard.export')}
              aria-keyshortcuts="Control+E"
              className="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-slate-700/50 text-slate-400 hover:bg-slate-800/50 hover:text-white disabled:cursor-not-allowed disabled:opacity-40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              <Download className="h-4 w-4" />
            </button>
            <div className="mx-1 h-4 w-px bg-white/10" />
          </>
        )}
        <button
          type="button"
          onClick={onUndo}
          disabled={!canUndo}
          title={undoDescription}
          aria-keyshortcuts="Control+Z"
          className="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-slate-700/50 text-slate-400 hover:bg-slate-800/50 hover:text-white disabled:cursor-not-allowed disabled:opacity-30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
        >
          <Undo2 className="h-4 w-4" />
        </button>
        <button
          type="button"
          onClick={onRedo}
          disabled={!canRedo}
          title={redoDescription}
          aria-keyshortcuts="Control+Y"
          className="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-slate-700/50 text-slate-400 hover:bg-slate-800/50 hover:text-white disabled:cursor-not-allowed disabled:opacity-30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
        >
          <Redo2 className="h-4 w-4" />
        </button>
        <button
          type="button"
          onClick={onOpenShortcuts}
          title={t('smartDashboard.shortcuts')}
          className="ml-auto inline-flex h-8 items-center justify-center gap-1 rounded-lg border border-slate-700/50 px-2 text-slate-400 hover:bg-slate-800/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
        >
          <Keyboard className="h-4 w-4" />
          <span className="text-[10px] font-mono opacity-60">?</span>
        </button>
      </div>

      {/* ── Tab Bar ── */}
      <div className="flex shrink-0 border-b border-white/5">
        {TABS.map((tab) => {
          const isActive = sidebarTab === tab.id;
          const hasBadge = tab.id === 'tag' && hasTagSelection;
          return (
            <button
              key={tab.id}
              type="button"
              onClick={() => setSidebarTab(tab.id)}
              className={`relative flex flex-1 items-center justify-center gap-1 py-2.5 text-xs font-medium transition-colors focus-visible:outline-none focus-visible:ring-inset focus-visible:ring-1 focus-visible:ring-blue-500 ${
                isActive
                  ? 'border-b-2 border-blue-400 text-blue-200'
                  : 'border-b-2 border-transparent text-slate-400 hover:text-slate-200'
              }`}
            >
              {tab.label}
              {hasBadge && (
                <span className="h-1.5 w-1.5 rounded-full bg-blue-400" aria-hidden />
              )}
            </button>
          );
        })}
      </div>

      {/* ── Tab Content ── */}
      <div className="scrollbar-thin scrollbar-thumb-slate-700/50 scrollbar-track-transparent flex-1 overflow-y-auto">
        {sidebarTab === 'plan' && (
          <SmartDashboardPlanningTab {...planningTabProps} />
        )}

        {sidebarTab === 'tag' && (
          <SmartDashboardTagPanel {...tagPanelProps} />
        )}

        {sidebarTab === 'modbus' && (
          <SmartDashboardModbusPanel {...modbusPanelProps} />
        )}

        {sidebarTab === 'commit' && (
          <SmartDashboardCommitPanel {...commitPanelProps} />
        )}
      </div>
    </div>
  );
}
