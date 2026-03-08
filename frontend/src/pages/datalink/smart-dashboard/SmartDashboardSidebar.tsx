import { Download, Keyboard, Redo2, Undo2, Upload, LayoutGrid } from 'lucide-react';
import type { TFunction } from 'i18next';
import type { Device } from '../../../types/datalink';
import { designSystem } from '../../../styles/designSystem';
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
        <button
          type="button"
          onClick={onOpenWorkbench}
          title="Server Memory Grid"
          className={`${designSystem.components.button.icon} h-8 w-8 border border-slate-700/50`}
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
          className={`${designSystem.components.button.icon} ml-auto h-8 gap-1 px-2 border border-slate-700/50`}
        >
          <Keyboard className="h-4 w-4" />
          <span className="text-[10px] font-mono opacity-60">?</span>
        </button>
      </div>

      {/* ── 流程指引 ── */}
      {selectedDevice && (
        <div className="shrink-0 border-b border-white/5 px-3 py-2" role="region" aria-label={t('smartDashboard.sidebar.workflowGuide', { defaultValue: '工作流程' })}>
          <div className="rounded-lg border border-blue-500/20 bg-blue-500/5 p-2">
            <p className="text-[10px] font-semibold uppercase tracking-wider text-blue-200 mb-1.5">
              {t('smartDashboard.sidebar.workflowGuide', { defaultValue: '工作流程' })}
            </p>
            <nav aria-label={t('smartDashboard.sidebar.workflowSteps', { defaultValue: '流程步驟' })} className="flex items-center gap-1 text-[10px] text-slate-300">
              <span className={sidebarTab === 'plan' ? 'text-blue-300 font-medium' : ''} aria-current={sidebarTab === 'plan' ? 'step' : undefined}>1. 規劃</span>
              <span className="text-slate-500" aria-hidden="true">→</span>
              <span className={sidebarTab === 'tag' ? 'text-blue-300 font-medium' : ''} aria-current={sidebarTab === 'tag' ? 'step' : undefined}>2. Tag</span>
              <span className="text-slate-500" aria-hidden="true">→</span>
              <span className={sidebarTab === 'modbus' ? 'text-blue-300 font-medium' : ''} aria-current={sidebarTab === 'modbus' ? 'step' : undefined}>3. Modbus</span>
              <span className="text-slate-500" aria-hidden="true">→</span>
              <span className={sidebarTab === 'commit' ? 'text-blue-300 font-medium' : ''} aria-current={sidebarTab === 'commit' ? 'step' : undefined}>4. 提交</span>
            </nav>
          </div>
        </div>
      )}

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
