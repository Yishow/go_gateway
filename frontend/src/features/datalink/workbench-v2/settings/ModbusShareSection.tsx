import { useTranslation } from 'react-i18next';
import type { ModbusShareSettings } from '../state/types';

/**
 * ModbusShareSection 元件屬性
 */
interface ModbusShareSectionProps {
  settings: ModbusShareSettings;
  onChange: (patch: Partial<ModbusShareSettings>) => void;
}

/**
 * Local Modbus Share 設定區塊元件
 * 落地需求：「Local Modbus Share configuration」
 */
export function ModbusShareSection({ settings, onChange }: ModbusShareSectionProps) {
  const { t } = useTranslation('workbench-v2');

  const { enabled, bind_address, port, slave_id, base_register } = settings;

  return (
    <div className="space-y-6 bg-gray-900/10 border border-gray-800 rounded-2xl p-6 backdrop-blur-sm">
      {/* 標題與側欄開關 */}
      <div className="flex justify-between items-start">
        <div>
          <h3 className="text-lg font-semibold text-white">
            {t('settings.modbus_share_title', 'Local Modbus Share')}
          </h3>
          <p className="text-xs text-gray-500 mt-1">
            {t('settings.modbus_share_subtitle', '將採集資料映射並轉發至本地虛擬 Modbus Server 暫存器，供上層 SCADA/HMI 直接抓取。')}
          </p>
        </div>
        
        {/* 總開關 */}
        <div className="flex items-center gap-2">
          <span className="text-xs text-gray-400 select-none">
            {enabled ? t('status.active', '啟用') : t('status.disabled', '未啟用')}
          </span>
          <button
            type="button"
            onClick={() => onChange({ enabled: !enabled })}
            className={`
              relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-200 focus:outline-none cursor-pointer
              ${enabled ? 'bg-blue-600' : 'bg-gray-800'}
            `}
          >
            <span
              className={`
                inline-block h-4 w-4 transform rounded-full bg-white transition-transform duration-200
                ${enabled ? 'translate-x-6' : 'translate-x-1'}
              `}
            />
          </button>
        </div>
      </div>

      {/* 依據開關渲染內容 */}
      {enabled ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 p-4 rounded-xl border border-gray-800/60 bg-gray-950/20 sweep-in">
          {/* 綁定位址 */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_bind_address', '監聽本機位址')}
            </label>
            <input
              type="text"
              value={bind_address}
              onChange={(e) => onChange({ bind_address: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all font-mono"
              placeholder="0.0.0.0"
            />
          </div>

          {/* 通訊埠 */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_port', '監聽通訊埠 (Port)')}
            </label>
            <input
              type="number"
              min={1}
              max={65535}
              value={port}
              onChange={(e) => onChange({ port: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all font-mono"
              placeholder="5020"
            />
          </div>

          {/* Slave ID */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_slave_id', '本機站號 (Slave ID)')}
            </label>
            <input
              type="number"
              min={1}
              max={247}
              value={slave_id}
              onChange={(e) => onChange({ slave_id: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all font-mono"
              placeholder="1"
            />
          </div>

          {/* 起始暫存器位址 */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_base_register', '起始暫存器位址')}
            </label>
            <input
              type="number"
              min={1}
              max={99999}
              value={base_register}
              onChange={(e) => onChange({ base_register: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all font-mono"
              placeholder="40001"
            />
          </div>
        </div>
      ) : (
        <div className="flex items-center justify-center p-6 border border-dashed border-gray-800 rounded-xl bg-gray-950/10 text-gray-500 text-xs">
          <span>{t('settings.modbus_share_disabled_msg', 'Modbus 轉發服務已停用。開啟後點位分配將會與 Step 2 的 Modbus Share 暫存器配置連動。')}</span>
        </div>
      )}
    </div>
  );
}
