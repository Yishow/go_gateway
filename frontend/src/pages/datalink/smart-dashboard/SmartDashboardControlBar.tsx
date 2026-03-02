import type { TFunction } from 'i18next';

/**
 * SmartDashboard 控制列：僅提供選擇設備與新增設備按鈕。
 */

interface SmartDashboardControlBarProps {
  onChooseDevice: () => void;
  onCreateDevice: () => void;
  t: TFunction;
}

export default function SmartDashboardControlBar({
  onChooseDevice,
  onCreateDevice,
  t,
}: SmartDashboardControlBarProps) {
  return (
    <section className="mx-3 mt-3 rounded-2xl border border-white/10 bg-gradient-to-r from-slate-900/90 via-slate-800/70 to-slate-900/90 shadow-xl shadow-black/20 ring-1 ring-white/5 backdrop-blur-md sm:mx-4">
      <div className="flex flex-col gap-3 px-4 py-3 sm:flex-row sm:items-center sm:justify-end sm:gap-2 sm:px-5 sm:py-3">
        <button
          type="button"
          onClick={onChooseDevice}
          className="min-h-9 cursor-pointer rounded-lg border border-slate-600 bg-slate-800/80 px-3 py-2 text-xs font-semibold text-slate-200 transition-all duration-200 hover:border-slate-500 hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
        >
          {t('smartDashboard.actions.chooseDevice')}
        </button>
        <button
          type="button"
          onClick={onCreateDevice}
          className="min-h-9 cursor-pointer rounded-lg bg-blue-600 px-3 py-2 text-xs font-semibold text-white shadow-md shadow-blue-500/25 transition-all duration-200 hover:bg-blue-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
        >
          {t('smartDashboard.actions.createDevice')}
        </button>
      </div>
    </section>
  );
}
