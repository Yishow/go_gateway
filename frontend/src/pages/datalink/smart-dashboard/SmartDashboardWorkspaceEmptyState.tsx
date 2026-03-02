import { Sparkles } from 'lucide-react';
import type { TFunction } from 'i18next';

interface SmartDashboardWorkspaceEmptyStateProps {
  t: TFunction;
  onChooseDevice: () => void;
  onCreateDevice: () => void;
}

export default function SmartDashboardWorkspaceEmptyState({
  t,
  onChooseDevice,
  onCreateDevice,
}: SmartDashboardWorkspaceEmptyStateProps) {
  return (
    <div className="relative flex flex-1 flex-col items-center justify-center overflow-hidden p-6 text-slate-400 sm:p-8">
      <div className="absolute inset-0 bg-gradient-to-t from-blue-500/5 to-transparent opacity-50" />
      <div className="relative z-10 flex flex-col items-center text-center">
        <div className="mb-6 flex h-24 w-24 items-center justify-center rounded-3xl border border-white/5 bg-gradient-to-br from-blue-500/20 to-purple-500/20 shadow-[0_0_50px_rgba(59,130,246,0.2)] sm:mb-8 sm:h-32 sm:w-32">
          <Sparkles className="h-10 w-10 text-blue-400 sm:h-12 sm:w-12" />
        </div>
        <h3 className="mb-3 text-xl font-bold text-white sm:text-2xl">{t('smartDashboard.welcomeTitle')}</h3>
        <p className="mb-6 max-w-md leading-relaxed text-slate-300">{t('smartDashboard.welcomeDescription')}</p>
        <div className="flex flex-wrap items-center justify-center gap-3">
          <button
            type="button"
            onClick={onChooseDevice}
            className="min-h-11 cursor-pointer rounded-lg border border-slate-600 bg-slate-800 px-4 py-2 text-sm font-semibold text-slate-200 hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            {t('smartDashboard.actions.chooseDevice')}
          </button>
          <button
            type="button"
            onClick={onCreateDevice}
            className="min-h-11 cursor-pointer rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            {t('smartDashboard.actions.createDevice')}
          </button>
        </div>
      </div>
    </div>
  );
}
