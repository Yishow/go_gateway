import { useTranslation } from 'react-i18next';
import type { GeneralSettings } from '../state/types';

/**
 * DiagnosticsSection 元件屬性
 */
interface DiagnosticsSectionProps {
  settings: GeneralSettings;
  onChange: (patch: Partial<GeneralSettings>) => void;
}

/**
 * 系統診斷與日誌設定元件
 * 落地需求：「UI / API / Diagnostics options」
 */
export function DiagnosticsSection({ settings, onChange }: DiagnosticsSectionProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="space-y-4 bg-gray-900/10 border border-gray-800 rounded-2xl p-5 backdrop-blur-sm">
      <div>
        <h4 className="text-sm font-semibold text-white">
          {t('settings.diagnostics_title', '日誌與開發診斷')}
        </h4>
        <span className="text-[10px] text-gray-500 block mt-0.5">
          {t('settings.diagnostics_subtitle', '控制系統運行日誌等級、逾時參數與開發者偵錯控制台。')}
        </span>
      </div>

      <div className="space-y-3">
        {/* 日誌等級 */}
        <div className="space-y-1.5">
          <label className="text-[11px] font-medium text-gray-400 block">
            {t('settings.diagnostics_log_level', '日誌等級 (Log Level)')}
          </label>
          <select
            value={settings.log_level}
            onChange={(e) => onChange({ log_level: e.target.value as 'trace' | 'debug' | 'info' | 'warn' | 'error' })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-2 py-1.5 text-white text-xs focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all cursor-pointer"
          >
            <option value="trace">TRACE</option>
            <option value="debug">DEBUG</option>
            <option value="info">INFO</option>
            <option value="warn">WARN</option>
            <option value="error">ERROR</option>
          </select>
        </div>

        {/* 逾時設定 */}
        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1.5">
            <label className="text-[10px] font-medium text-gray-400 block">
              {t('settings.diagnostics_timeout', '通訊逾時')}
            </label>
            <div className="relative">
              <input
                type="number"
                min={1}
                max={300}
                value={settings.timeout_seconds}
                onChange={(e) => onChange({ timeout_seconds: Math.max(1, parseInt(e.target.value) || 1) })}
                className="w-full bg-gray-950 border border-gray-800 rounded-lg px-2 py-1.5 text-white text-xs focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all"
              />
              <span className="absolute right-2 top-1.5 text-[10px] text-gray-500 select-none">秒</span>
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-[10px] font-medium text-gray-400 block">
              {t('settings.diagnostics_sse', 'SSE 心跳')}
            </label>
            <div className="relative">
              <input
                type="number"
                min={5}
                max={60}
                value={settings.sse_heartbeat_seconds}
                onChange={(e) => onChange({ sse_heartbeat_seconds: Math.max(5, parseInt(e.target.value) || 5) })}
                className="w-full bg-gray-950 border border-gray-800 rounded-lg px-2 py-1.5 text-white text-xs focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all"
              />
              <span className="absolute right-2 top-1.5 text-[10px] text-gray-500 select-none">秒</span>
            </div>
          </div>
        </div>

        {/* Audit Log Switch */}
        <div className="flex items-center justify-between pt-1">
          <label className="text-[11px] font-medium text-gray-300">
            {t('settings.diagnostics_audit', '記錄操作審計日誌')}
          </label>
          <button
            type="button"
            onClick={() => onChange({ enable_audit_log: !settings.enable_audit_log })}
            className={`
              relative inline-flex h-5 w-9 items-center rounded-full transition-colors duration-200 focus:outline-none cursor-pointer
              ${settings.enable_audit_log ? 'bg-blue-600' : 'bg-gray-800'}
            `}
          >
            <span
              className={`
                inline-block h-3 w-3 transform rounded-full bg-white transition-transform duration-200
                ${settings.enable_audit_log ? 'translate-x-5' : 'translate-x-1'}
              `}
            />
          </button>
        </div>

        {/* Debug Panel Switch */}
        <div className="flex items-center justify-between">
          <label className="text-[11px] font-medium text-gray-300">
            {t('settings.diagnostics_debug_panel', '顯示開發者調試面板')}
          </label>
          <button
            type="button"
            onClick={() => onChange({ enable_debug_panel: !settings.enable_debug_panel })}
            className={`
              relative inline-flex h-5 w-9 items-center rounded-full transition-colors duration-200 focus:outline-none cursor-pointer
              ${settings.enable_debug_panel ? 'bg-blue-600' : 'bg-gray-800'}
            `}
          >
            <span
              className={`
                inline-block h-3 w-3 transform rounded-full bg-white transition-transform duration-200
                ${settings.enable_debug_panel ? 'translate-x-5' : 'translate-x-1'}
              `}
            />
          </button>
        </div>
      </div>
    </div>
  );
}
