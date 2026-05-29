import { useTranslation } from 'react-i18next';
import type { TimeseriesSettings } from '../state/types';

/**
 * TimeseriesSection 元件屬性
 */
interface TimeseriesSectionProps {
  settings: TimeseriesSettings;
  onChange: (patch: Partial<TimeseriesSettings>) => void;
}

/**
 * 時序資料庫儲存策略設定元件
 * 落地需求：「Timeseries strategy fields」
 */
export function TimeseriesSection({ settings, onChange }: TimeseriesSectionProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="space-y-6 bg-gray-900/10 border border-gray-800 rounded-2xl p-6 backdrop-blur-sm">
      {/* 標題與說明 */}
      <div>
        <h3 className="text-lg font-semibold text-white">
          {t('settings.timeseries_title', '時序儲存策略')}
        </h3>
        <p className="text-xs text-gray-500 mt-1">
          {t('settings.timeseries_subtitle', '調整時序資料寫入時的解析度、物理磁碟分區割裂頻率以及保存期限。')}
        </p>
      </div>

      {/* 表單網格 */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* 時間精度 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {t('settings.timeseries_precision', '時間精度')}
          </label>
          <select
            value={settings.write_precision}
            onChange={(e) => onChange({ write_precision: e.target.value as 'second' | 'millisecond' })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all cursor-pointer"
          >
            <option value="millisecond">{t('settings.precision_ms', '毫秒 (millisecond)')}</option>
            <option value="second">{t('settings.precision_s', '秒 (second)')}</option>
          </select>
        </div>

        {/* 分區間隔 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {t('settings.timeseries_partition', '分區間隔')}
          </label>
          <select
            value={settings.partition_interval}
            onChange={(e) => onChange({ partition_interval: e.target.value as 'daily' | 'weekly' | 'monthly' })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all cursor-pointer"
          >
            <option value="daily">{t('settings.partition_daily', '每日分區 (daily)')}</option>
            <option value="weekly">{t('settings.partition_weekly', '每週分區 (weekly)')}</option>
            <option value="monthly">{t('settings.partition_monthly', '每月分區 (monthly)')}</option>
          </select>
        </div>

        {/* 批次寫入大小 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {t('settings.timeseries_batch', '批次寫入限制條數')}
          </label>
          <input
            type="number"
            min={1}
            max={5000}
            value={settings.batch_size}
            onChange={(e) => onChange({ batch_size: Math.max(1, parseInt(e.target.value) || 1) })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all"
            placeholder="500"
          />
        </div>

        {/* 保留天數 */}
        <div className="space-y-1.5">
          <label className="text-xs font-medium text-gray-400">
            {t('settings.timeseries_retention', '歷史保留天數')}
          </label>
          <div className="relative">
            <input
              type="number"
              min={1}
              max={3650}
              value={settings.retention_days}
              onChange={(e) => onChange({ retention_days: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all"
              placeholder="90"
            />
            <span className="absolute right-3 top-2 text-xs text-gray-500 select-none">
              {t('settings.days_unit', '天')}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
