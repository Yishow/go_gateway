import { Box } from 'lucide-react';
import type { Device } from '../../../types/datalink';
import type { DashboardModalIntent } from '../../../features/datalink/legacyRoutes';

type DashboardTab = 'overview' | 'devices' | 'settings';

interface SmartDashboardControlBarProps {
  tabs: Array<{ key: DashboardTab; label: string }>;
  activeTab: DashboardTab;
  onSelectTab: (tab: DashboardTab) => void;
  selectedDevice: Device | null;
  selectedDeviceStateView: {
    label: string;
    className: string;
    dotClass: string;
  };
  lastSwitchedAt: string | null;
  modalQuickLinks: Array<{ key: DashboardModalIntent; label: string }>;
  modalIntent: DashboardModalIntent | null;
  onOpenWorkflowModal: (modal: DashboardModalIntent) => void;
  onChooseDevice: () => void;
  onCreateDevice: () => void;
  isSwitchingDevice: boolean;
  hasUnsavedChanges: boolean;
  lastSyncLabel: string;
  cycleTimeLabel: string;
}

export default function SmartDashboardControlBar({
  tabs,
  activeTab,
  onSelectTab,
  selectedDevice,
  selectedDeviceStateView,
  lastSwitchedAt,
  modalQuickLinks,
  modalIntent,
  onOpenWorkflowModal,
  onChooseDevice,
  onCreateDevice,
  isSwitchingDevice,
  hasUnsavedChanges,
  lastSyncLabel,
  cycleTimeLabel,
}: SmartDashboardControlBarProps) {
  return (
    <section className="mx-3 mt-3 rounded-2xl border border-white/10 bg-gradient-to-r from-slate-900/90 via-slate-800/70 to-slate-900/90 shadow-xl shadow-black/20 ring-1 ring-white/5 backdrop-blur-md sm:mx-4">
      <div className="flex flex-col gap-3 px-4 py-3 sm:flex-row sm:items-center sm:justify-between sm:px-5 sm:py-3">
        <div className="flex shrink-0 items-center gap-1.5">
          {tabs.map((tab) => {
            const isActive = activeTab === tab.key;
            return (
              <button
                key={tab.key}
                type="button"
                onClick={() => onSelectTab(tab.key)}
                className={`min-h-9 cursor-pointer rounded-lg px-3.5 py-2 text-sm font-semibold transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 ${
                  isActive
                    ? 'bg-blue-500/90 text-white shadow-lg shadow-blue-500/25'
                    : 'border border-transparent bg-white/5 text-slate-300 hover:border-white/10 hover:bg-white/10 hover:text-white'
                }`}
              >
                {tab.label}
              </button>
            );
          })}
        </div>

        <div className="flex min-w-0 flex-1 flex-wrap items-center gap-x-4 gap-y-2 sm:min-w-0 sm:flex-initial sm:gap-x-5 sm:border-l sm:border-white/10 sm:pl-4">
          {selectedDevice ? (
            <>
              <div className="flex min-w-0 shrink-0 items-center gap-2.5">
                <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-white/10 bg-white/[0.07]">
                  <Box className="h-4 w-4 text-slate-400" aria-hidden />
                </div>
                <div className="flex min-w-0 items-center gap-2">
                  <span className="max-w-[140px] truncate text-sm font-semibold text-white sm:max-w-[220px]">
                    {selectedDevice.name}
                  </span>
                  <span
                    className={`inline-flex shrink-0 items-center gap-1 rounded-md px-1.5 py-0.5 text-[10px] font-bold uppercase tracking-wider ${selectedDeviceStateView.className}`}
                  >
                    <span className={`h-1 w-1 rounded-full ${selectedDeviceStateView.dotClass}`} aria-hidden />
                    {selectedDeviceStateView.label}
                  </span>
                </div>
              </div>
              <div className="flex flex-wrap items-center gap-x-3 gap-y-1 text-[11px]">
                <span className="shrink-0 font-mono text-slate-400">
                  ID:{selectedDevice.id.slice(0, 8)} · {selectedDevice.protocol}
                </span>
                <span className="shrink-0 text-slate-500">·</span>
                <span className="shrink-0 text-slate-400">
                  <span className="text-slate-500">{lastSyncLabel}</span>{' '}
                  <span className="font-mono text-slate-200">
                    {selectedDevice.last_test_at
                      ? new Date(selectedDevice.last_test_at).toLocaleTimeString()
                      : '-'}
                  </span>
                </span>
                <span className="shrink-0 text-slate-500">·</span>
                <span className="shrink-0">
                  <span className="text-slate-500">{cycleTimeLabel}</span>{' '}
                  <span className="font-mono text-blue-400">100 ms</span>
                </span>
                <span className="hidden shrink-0 text-slate-500 md:inline">·</span>
                <span className="hidden shrink-0 text-slate-400 md:inline">
                  最近切換:{' '}
                  <span className="text-slate-300">
                    {lastSwitchedAt ? new Date(lastSwitchedAt).toLocaleTimeString() : '-'}
                  </span>
                </span>
              </div>
              <div className="flex flex-wrap items-center gap-1.5">
                {modalQuickLinks.map((item) => (
                  <button
                    key={item.key}
                    type="button"
                    onClick={() => onOpenWorkflowModal(item.key)}
                    className={`min-h-8 rounded-md border px-2 py-1 text-[10px] font-semibold uppercase tracking-wide transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 ${
                      modalIntent === item.key
                        ? 'border-indigo-400/60 bg-indigo-500/30 text-indigo-100'
                        : 'border-slate-600 bg-slate-800/70 text-slate-300 hover:bg-slate-700'
                    }`}
                  >
                    {item.label}
                  </button>
                ))}
              </div>
              <div className="ml-auto flex shrink-0 gap-2 sm:ml-0">
                <button
                  type="button"
                  onClick={onChooseDevice}
                  disabled={isSwitchingDevice}
                  className="min-h-9 cursor-pointer rounded-lg border border-blue-400/50 bg-blue-500/25 px-3 py-2 text-xs font-semibold text-blue-100 shadow-sm transition-all duration-200 hover:border-blue-400/60 hover:bg-blue-500/35 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  <span className="inline-flex items-center gap-1.5">
                    切換設備
                    {hasUnsavedChanges && (
                      <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-rose-400" aria-hidden />
                    )}
                    {isSwitchingDevice && <span className="text-[10px] text-blue-200">切換中...</span>}
                  </span>
                </button>
                <button
                  type="button"
                  onClick={onCreateDevice}
                  className="min-h-9 cursor-pointer rounded-lg bg-blue-600 px-3 py-2 text-xs font-semibold text-white shadow-md shadow-blue-500/25 transition-all duration-200 hover:bg-blue-500 hover:shadow-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                >
                  新增設備
                </button>
              </div>
            </>
          ) : (
            <>
              <div className="flex items-center gap-2.5">
                <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border border-white/10 bg-white/[0.07]">
                  <Box className="h-4 w-4 text-slate-500" aria-hidden />
                </div>
                <div>
                  <p className="text-sm font-semibold text-white">尚未選擇設備</p>
                  <p className="text-[11px] text-slate-400">請先選擇或新增設備以開始規劃</p>
                </div>
              </div>
              <div className="ml-auto flex shrink-0 gap-2 sm:ml-0">
                <button
                  type="button"
                  onClick={onChooseDevice}
                  className="min-h-9 cursor-pointer rounded-lg border border-slate-600 bg-slate-800/80 px-3 py-2 text-xs font-semibold text-slate-200 transition-all duration-200 hover:border-slate-500 hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                >
                  選擇設備
                </button>
                <button
                  type="button"
                  onClick={onCreateDevice}
                  className="min-h-9 cursor-pointer rounded-lg bg-blue-600 px-3 py-2 text-xs font-semibold text-white shadow-md shadow-blue-500/25 transition-all duration-200 hover:bg-blue-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                >
                  新增設備
                </button>
              </div>
            </>
          )}
          <div className="w-full border-t border-white/10 pt-2">
            <div className="flex flex-wrap items-center gap-1.5">
              <span className="mr-1 text-[10px] uppercase tracking-wider text-slate-500">原頁面入口</span>
              {modalQuickLinks.map((item) => (
                <button
                  key={`global-${item.key}`}
                  type="button"
                  onClick={() => onOpenWorkflowModal(item.key)}
                  className={`min-h-8 rounded-md border px-2 py-1 text-[10px] font-semibold uppercase tracking-wide transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 ${
                    modalIntent === item.key
                      ? 'border-indigo-400/60 bg-indigo-500/30 text-indigo-100'
                      : 'border-slate-600 bg-slate-800/70 text-slate-300 hover:bg-slate-700'
                  }`}
                >
                  {item.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
