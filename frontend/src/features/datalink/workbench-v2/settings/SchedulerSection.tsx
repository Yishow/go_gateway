import { Toggle } from '../components/Toggle';
import { useTranslation } from 'react-i18next';
import type { SchedulerSettings } from '../state/types';

/**
 * SchedulerSection 元件屬性
 */
interface SchedulerSectionProps {
  settings: SchedulerSettings;
  onChange: (patch: Partial<SchedulerSettings>) => void;
}

/**
 * 排程器預設參數設定元件
 * 落地需求：「Scheduler defaults fields」
 */
export function SchedulerSection({ settings, onChange }: SchedulerSectionProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="space-y-6 bg-[var(--gw-surface)] border border-[var(--gw-border)] rounded-2xl p-6">
      {/* 標題與說明 */}
      <div>
        <h3 className="text-lg font-semibold text-[var(--gw-text)]">
          {t('settings.scheduler_title', '排程器核心')}
        </h3>
        <p className="text-xs text-[var(--gw-muted)] mt-1">
          {t('settings.scheduler_subtitle', '排程採集重試策略、熔斷器臨界值以及開機啟動選項設定。')}
        </p>
      </div>

      {/* 表單網格 */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* 預設輪詢間隔 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-[var(--gw-muted)]">
            {t('settings.scheduler_interval', '預設輪詢週期')}
          </label>
          <div className="relative">
            <input
              type="number"
              min={100}
              max={60000}
              value={settings.default_interval_ms}
              onChange={(e) => onChange({ default_interval_ms: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-3 py-2 text-[var(--gw-text)] text-sm focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all"
              placeholder="1000"
            />
            <span className="absolute right-3 top-2 text-xs text-[var(--gw-muted)] select-none">
              毫秒 (ms)
            </span>
          </div>
        </div>

        {/* 預設重試次數 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-[var(--gw-muted)]">
            {t('settings.scheduler_retry', '連線失敗重試上限')}
          </label>
          <div className="relative">
            <input
              type="number"
              min={0}
              max={10}
              value={settings.default_retry_count}
              onChange={(e) => onChange({ default_retry_count: Math.max(0, parseInt(e.target.value) || 0) })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-3 py-2 text-[var(--gw-text)] text-sm focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all"
              placeholder="3"
            />
            <span className="absolute right-3 top-2 text-xs text-[var(--gw-muted)] select-none">
              次
            </span>
          </div>
        </div>

        {/* 重試延遲 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-[var(--gw-muted)]">
            {t('settings.scheduler_retry_delay', '重試延遲時間')}
          </label>
          <div className="relative">
            <input
              type="number"
              min={100}
              max={5000}
              value={settings.default_retry_delay_ms}
              onChange={(e) => onChange({ default_retry_delay_ms: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-3 py-2 text-[var(--gw-text)] text-sm focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all"
              placeholder="500"
            />
            <span className="absolute right-3 top-2 text-xs text-[var(--gw-muted)] select-none">
              毫秒 (ms)
            </span>
          </div>
        </div>

        {/* 熔斷閥值 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-[var(--gw-muted)]">
            {t('settings.scheduler_breaker', '熔斷計數器閥值')}
          </label>
          <div className="relative">
            <input
              type="number"
              min={1}
              max={100}
              value={settings.breaker_threshold}
              onChange={(e) => onChange({ breaker_threshold: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-3 py-2 text-[var(--gw-text)] text-sm focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all"
              placeholder="10"
            />
            <span className="absolute right-3 top-2 text-xs text-[var(--gw-muted)] select-none">
              次
            </span>
          </div>
        </div>

        {/* 開機啟動 (跨兩欄) */}
        <div className="md:col-span-2 pt-2 flex items-center justify-between">
          <div>
            <label className="text-xs font-medium text-[var(--gw-text)] block">
              {t('settings.scheduler_auto_start', '開機自動啟動')}
            </label>
            <span className="text-xs text-[var(--gw-muted)] block mt-0.5">
              {t('settings.scheduler_auto_start_desc', '當閘道核心服務啟動時，自動激活所有已部署的採集排程器。')}
            </span>
          </div>
          <div className="shrink-0">
            <Toggle checked={settings.auto_start} onChange={(checked) => onChange({ auto_start: checked })} aria-label={t('settings.scheduler_auto_start', '開機自動啟動')} />
          </div>
        </div>
      </div>
    </div>
  );
}
