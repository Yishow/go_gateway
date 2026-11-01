import { useTranslation } from 'react-i18next';
import type { GeneralSettings } from '../state/types';

/**
 * ApiSection 元件屬性
 */
interface ApiSectionProps {
  settings: GeneralSettings;
  onChange: (patch: Partial<GeneralSettings>) => void;
}

/**
 * API 連線參數設定元件
 * 落地需求：「UI / API / Diagnostics options」
 */
export function ApiSection({ settings, onChange }: ApiSectionProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="space-y-4 bg-[var(--gw-surface)] border border-[var(--gw-border)] rounded-2xl p-5">
      <div>
        <h4 className="text-sm font-semibold text-[var(--gw-text)]">
          {t('settings.api_title', 'API 連線端點')}
        </h4>
        <span className="text-xs text-[var(--gw-muted)] block mt-0.5">
          {t('settings.api_subtitle', '設定後端閘道 API 的基礎連線位址與協議版本。')}
        </span>
      </div>

      <div className="space-y-3">
        {/* API 基礎位址 */}
        <div className="space-y-1.5">
          <label className="text-sm font-medium text-[var(--gw-text)] block">
            {t('settings.api_base', 'API 伺服器位址 (API Base)')}
          </label>
          <input
            type="text"
            value={settings.api_base}
            onChange={(e) => onChange({ api_base: e.target.value })}
            className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-2.5 py-1.5 text-[var(--gw-text)] text-xs focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all font-mono"
            placeholder="http://localhost:8080"
          />
        </div>

        {/* API 版本 */}
        <div className="space-y-1.5">
          <label className="text-sm font-medium text-[var(--gw-text)] block">
            {t('settings.api_version', 'API 路由版本')}
          </label>
          <select
            value={settings.api_version}
            onChange={(e) => onChange({ api_version: e.target.value as 'v1' | 'v2' })}
            className="w-full bg-[var(--gw-surface)] border border-[var(--gw-control)] rounded-lg px-2 py-1.5 text-[var(--gw-text)] text-xs focus:border-[var(--gw-focus)] focus:ring-1 focus:ring-[var(--gw-focus)] outline-none transition-all cursor-pointer"
          >
            <option value="v1">v1 (stable)</option>
            <option value="v2">v2 (beta)</option>
          </select>
        </div>
      </div>
    </div>
  );
}
