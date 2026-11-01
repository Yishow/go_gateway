import { Toggle } from '../components/Toggle';
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
    <div className="space-y-4 bg-[var(--gw-surface)] border border-[var(--gw-border)] rounded-2xl p-5">
      <div>
        <h4 className="text-sm font-semibold text-[var(--gw-text)]">
          {t('settings.diagnostics_title', '日誌與開發診斷')}
        </h4>
        <span className="text-xs text-[var(--gw-muted)] block mt-0.5">
          {t('settings.diagnostics_subtitle', '控制系統運行日誌等級、逾時參數與開發者偵錯控制台。')}
        </span>
      </div>

      <div className="space-y-3">
        {/* 日誌等級 */}
        <div className="space-y-1.5">
          <label className="text-sm font-medium text-[var(--gw-text)] block">
            {t('settings.diagnostics_log_level', '日誌等級 (Log Level)')}
          </label>
          <select
            value={settings.log_level}
            onChange={(e) => onChange({ log_level: e.target.value as 'trace' | 'debug' | 'info' | 'warn' | 'error' })}
            className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-2 py-1.5 text-[var(--gw-text)] text-xs focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all cursor-pointer"
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
            <label className="text-sm font-medium text-[var(--gw-text)] block">
              {t('settings.diagnostics_timeout', '通訊逾時')}
            </label>
            <div className="relative">
              <input
                type="number"
                min={1}
                max={300}
                value={settings.timeout_seconds}
                onChange={(e) => onChange({ timeout_seconds: Math.max(1, parseInt(e.target.value) || 1) })}
                className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-2 py-1.5 text-[var(--gw-text)] text-xs focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all"
              />
              <span className="absolute right-2 top-1.5 text-xs text-[var(--gw-muted)] select-none">秒</span>
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-[var(--gw-text)] block">
              {t('settings.diagnostics_sse', 'SSE 心跳')}
            </label>
            <div className="relative">
              <input
                type="number"
                min={5}
                max={60}
                value={settings.sse_heartbeat_seconds}
                onChange={(e) => onChange({ sse_heartbeat_seconds: Math.max(5, parseInt(e.target.value) || 5) })}
                className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-2 py-1.5 text-[var(--gw-text)] text-xs focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all"
              />
              <span className="absolute right-2 top-1.5 text-xs text-[var(--gw-muted)] select-none">秒</span>
            </div>
          </div>
        </div>

        {/* Audit Log Switch */}
        <div className="flex items-center justify-between pt-1">
          <label className="text-xs font-medium text-[var(--gw-text)]">
            {t('settings.diagnostics_audit', '記錄操作審計日誌')}
          </label>
          <div className="shrink-0">
            <Toggle checked={settings.enable_audit_log} onChange={(checked) => onChange({ enable_audit_log: checked })} aria-label={t('settings.diagnostics_audit', '記錄操作審計日誌')} />
          </div>
        </div>

        {/* Debug Panel Switch */}
        <div className="flex items-center justify-between">
          <label className="text-xs font-medium text-[var(--gw-text)]">
            {t('settings.diagnostics_debug_panel', '顯示開發者調試面板')}
          </label>
          <div className="shrink-0">
            <Toggle checked={settings.enable_debug_panel} onChange={(checked) => onChange({ enable_debug_panel: checked })} aria-label={t('settings.diagnostics_debug_panel', '顯示開發者調試面板')} />
          </div>
        </div>
      </div>
    </div>
  );
}
