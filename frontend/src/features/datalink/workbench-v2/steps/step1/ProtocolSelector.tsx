import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { PROTOCOLS } from '../../state/protocols';
import type { ProtocolId } from '../../state/types';
import { Icon } from '../../components/Icon';
import type { IconName } from '../../components/Icon';

export interface ProtocolSelectorProps {
  /** 當前選中的協議 ID */
  value: ProtocolId;
  /** 當選中協議改變時的 callback 函數 */
  onChange: (protocol: ProtocolId) => void;
}

/**
 * 依協議 ID 取得對應的 Icon 名稱
 * 
 * @param protocol 協議 ID
 * @returns 圖示名稱
 */
function getIconForProtocol(protocol: ProtocolId): IconName {
  switch (protocol) {
    case 'modbus_rtu':
      return 'sliders';
    case 'mqtt':
      return 'flow';
    case 'fatek_fbs':
    case 'mc_3e':
      return 'bolt';
    case 'modbus_tcp':
    case 'modbus_udp':
    default:
      return 'cable';
  }
}

/**
 * 協議卡片選擇器元件 (ProtocolSelector)
 * 
 * 落地需求 **Protocol switching with reset**。
 * 提供 6 種協議的格狀卡片介面，選中時呈現高亮與打勾圖示。
 */
export const ProtocolSelector: React.FC<ProtocolSelectorProps> = ({ value, onChange }) => {
  const { t } = useTranslation('workbench-v2');

  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
      {PROTOCOLS.map((protocol) => {
        const isSelected = protocol.id === value;
        const iconName = getIconForProtocol(protocol.id);

        return (
          <button
            key={protocol.id}
            type="button"
            onClick={() => onChange(protocol.id)}
            data-testid={`protocol-card-${protocol.id}`}
            className={`relative flex flex-col items-start rounded-xl border p-3.5 text-left transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500/30 ${
              isSelected
                ? 'border-blue-500 bg-blue-600/10 shadow-lg shadow-blue-500/5'
                : 'border-slate-800 bg-slate-900/30 hover:border-slate-700 hover:bg-slate-900/50'
            }`}
          >
            {/* 右上角選中標記 */}
            {isSelected && (
              <div className="absolute right-2.5 top-2.5 flex h-4 w-4 items-center justify-center rounded-full bg-blue-500 text-slate-950">
                <Icon name="check" className="h-3 w-3 stroke-[2.5]" />
              </div>
            )}

            {/* 圖示 */}
            <div
              className={`mb-2.5 flex h-8 w-8 items-center justify-center rounded-lg ${
                isSelected ? 'bg-blue-500/15 text-blue-400' : 'bg-slate-800/60 text-slate-400'
              }`}
            >
              <Icon name={iconName} className="h-4.5 w-4.5" />
            </div>

            {/* 協議名稱 */}
            <span className="text-sm font-semibold text-slate-100">
              {t(protocol.nameKey)}
            </span>

            {/* 協議描述 */}
            <span className="mt-1 text-xs text-slate-400/80 leading-normal line-clamp-2">
              {t(protocol.descKey)}
            </span>
          </button>
        );
      })}
    </div>
  );
};
