import { useTranslation } from 'react-i18next';

/**
 * 系統設定頁面頂部標題元件
 * 落地需求：「Settings page layout」之標題列
 */
export function SettingsHeader() {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="relative overflow-hidden rounded-2xl border border-blue-500/20 bg-gradient-to-r from-blue-600/15 via-blue-500/5 to-slate-950 p-6 md:p-8 shadow-lg shadow-blue-500/[0.02]">
      {/* 漸層光暈裝飾 */}
      <div className="absolute -right-16 -top-16 w-48 h-48 rounded-full bg-blue-500/10 blur-3xl pointer-events-none" />
      
      <div className="flex items-center gap-4">
        {/* Sliders Icon */}
        <div className="grid place-items-center w-12 h-12 rounded-xl bg-blue-500/20 border border-blue-500/30 text-blue-400 shadow-inner">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="w-5 h-5"
          >
            <line x1="4" y1="21" x2="4" y2="14" />
            <line x1="4" y1="10" x2="4" y2="3" />
            <line x1="12" y1="21" x2="12" y2="12" />
            <line x1="12" y1="8" x2="12" y2="3" />
            <line x1="20" y1="21" x2="20" y2="16" />
            <line x1="20" y1="12" x2="20" y2="3" />
            <line x1="1" y1="14" x2="7" y2="14" />
            <line x1="9" y1="8" x2="15" y2="8" />
            <line x1="17" y1="16" x2="23" y2="16" />
          </svg>
        </div>
        
        <div>
          <h1 className="text-xl md:text-2xl font-bold tracking-tight text-white">
            {t('settings.title', '系統設定')}
          </h1>
          <p className="text-xs md:text-sm text-gray-400 mt-1">
            {t('settings.subtitle', '管理閘道全局連線、儲存排程、位元組映射轉發以及開發人員診斷工具。')}
          </p>
        </div>
      </div>
    </div>
  );
}
