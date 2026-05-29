import { useTranslation } from 'react-i18next';
import type { GeneralSettings } from '../state/types';

/**
 * UiSection 元件屬性
 */
interface UiSectionProps {
  settings: GeneralSettings;
  onChange: (patch: Partial<GeneralSettings>) => void;
}

/**
 * 介面偏好設定元件
 * 落地需求：「UI / API / Diagnostics options」
 */
export function UiSection({ settings, onChange }: UiSectionProps) {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="space-y-4 bg-gray-900/10 border border-gray-800 rounded-2xl p-5 backdrop-blur-sm">
      <div>
        <h4 className="text-sm font-semibold text-white">
          {t('settings.ui_title', '介面與格式偏好')}
        </h4>
        <span className="text-[10px] text-gray-500 block mt-0.5">
          {t('settings.ui_subtitle', '設定工作台外觀樣式、語系以及暫存器位址顯示格式。')}
        </span>
      </div>

      <div className="space-y-3">
        {/* 視覺主題 */}
        <div className="space-y-1.5">
          <label className="text-[11px] font-medium text-gray-400 block">
            {t('settings.ui_theme', '視覺主題')}
          </label>
          <select
            value={settings.theme}
            onChange={(e) => onChange({ theme: e.target.value as 'dark' | 'light' | 'auto' })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-2 py-1.5 text-white text-xs focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all cursor-pointer"
          >
            <option value="dark">{t('settings.theme_dark', '深色模式 (Dark)')}</option>
            <option value="light">{t('settings.theme_light', '淺色模式 (Light)')}</option>
            <option value="auto">{t('settings.theme_auto', '跟隨系統 (Auto)')}</option>
          </select>
        </div>

        {/* 顯示語系 */}
        <div className="space-y-1.5">
          <label className="text-[11px] font-medium text-gray-400 block">
            {t('settings.ui_locale', '顯示語系')}
          </label>
          <select
            value={settings.locale}
            onChange={(e) => onChange({ locale: e.target.value as 'zh-TW' | 'en' })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-2 py-1.5 text-white text-xs focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all cursor-pointer"
          >
            <option value="zh-TW">繁體中文</option>
            <option value="en">English</option>
          </select>
        </div>

        {/* 位址格式 */}
        <div className="space-y-1.5">
          <label className="text-[11px] font-medium text-gray-400 block">
            {t('settings.ui_addr_format', '暫存器位址格式')}
          </label>
          <select
            value={settings.addr_format}
            onChange={(e) => onChange({ addr_format: e.target.value as 'modbus' | 'hex' | 'raw' })}
            className="w-full bg-gray-950 border border-gray-800 rounded-lg px-2 py-1.5 text-white text-xs focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all cursor-pointer"
          >
            <option value="modbus">{t('settings.addr_modbus', 'Modbus 式 (e.g. 40001)')}</option>
            <option value="hex">{t('settings.addr_hex', '十六進位 (e.g. 0x0000)')}</option>
            <option value="raw">{t('settings.addr_raw', '十進位 (e.g. 0)')}</option>
          </select>
        </div>
      </div>
    </div>
  );
}
