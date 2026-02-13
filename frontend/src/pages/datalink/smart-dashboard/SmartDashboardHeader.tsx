import type { RefObject } from 'react';
import { Bell, Box, Search, Settings } from 'lucide-react';

interface SmartDashboardHeaderProps {
  searchInputRef: RefObject<HTMLInputElement | null>;
  productSubtitle: string;
  searchLabel: string;
  searchPlaceholder: string;
  notificationsLabel: string;
  preferencesLabel: string;
  profileMenuLabel: string;
  onOpenLocalModbusWorkbench: () => void;
}

export default function SmartDashboardHeader({
  searchInputRef,
  productSubtitle,
  searchLabel,
  searchPlaceholder,
  notificationsLabel,
  preferencesLabel,
  profileMenuLabel,
  onOpenLocalModbusWorkbench,
}: SmartDashboardHeaderProps) {
  return (
    <header className="flex flex-wrap items-center justify-between gap-3 border-b border-white/5 px-4 py-3 sm:px-6">
      <div className="flex min-w-0 items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-600/20 text-blue-400 ring-1 ring-blue-500/30">
          <Box className="h-6 w-6" />
        </div>
        <div>
          <h1 className="bg-gradient-to-r from-white to-slate-400 bg-clip-text text-lg font-bold tracking-tight text-transparent sm:text-xl">
            GoGateway
          </h1>
          <p className="text-[10px] font-medium uppercase tracking-wider text-slate-400">
            {productSubtitle}
          </p>
        </div>
      </div>

      <div className="group relative order-3 w-full sm:order-2 sm:mx-4 sm:w-auto sm:max-w-lg sm:flex-1">
        <label htmlFor="dashboard-search" className="sr-only">
          {searchLabel}
        </label>
        <div className="absolute inset-0 rounded-full bg-blue-500/20 opacity-0 blur-xl transition-opacity group-focus-within:opacity-100" />
        <Search className="absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400 transition-colors group-focus-within:text-blue-400" />
        <input
          ref={searchInputRef}
          id="dashboard-search"
          type="text"
          placeholder={searchPlaceholder}
          className="w-full rounded-full border border-slate-700/50 bg-slate-900/50 py-2.5 pl-11 pr-4 text-sm text-slate-200 shadow-lg shadow-black/20 backdrop-blur-md placeholder:text-slate-500 transition-all focus:border-blue-500/50 focus:outline-none focus:ring-2 focus:ring-blue-500/50"
        />
      </div>

      <div className="order-2 flex items-center gap-3 sm:order-3">
        <button
          type="button"
          onClick={onOpenLocalModbusWorkbench}
          className="min-h-11 rounded-lg border border-cyan-400/40 bg-cyan-500/15 px-3 py-2 text-xs font-semibold text-cyan-100 hover:bg-cyan-500/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400"
        >
          Server Memory Grid
        </button>
        <div className="flex items-center gap-2 rounded-full border border-slate-700/50 bg-slate-800/50 p-1">
          <button
            type="button"
            aria-label={notificationsLabel}
            className="rounded-full p-2 text-slate-400 transition-all hover:bg-slate-700/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            <Bell className="h-5 w-5" />
          </button>
          <button
            type="button"
            aria-label={preferencesLabel}
            className="rounded-full p-2 text-slate-400 transition-all hover:bg-slate-700/50 hover:text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            <Settings className="h-5 w-5" />
          </button>
        </div>
        <button
          type="button"
          aria-label={profileMenuLabel}
          className="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 text-sm font-bold shadow-lg shadow-blue-500/25 ring-2 ring-slate-900 transition-all hover:ring-offset-2 hover:ring-offset-slate-900 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
        >
          Y
        </button>
      </div>
    </header>
  );
}
