import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { dataTypeWidth } from '../../state/sourceRule';
import { addressParser } from '../../../../../utils/addressParser';
import type { ProtocolType } from '../../../../../types/datalink';

export interface RangeSummaryProps {
  startAddress: string;
  count: number;
  dataType: string;
  protocol?: ProtocolType;
}

/**
 * 規則位址範圍與功能碼摘要元件
 * 
 * 依起點前綴與協議推斷功能碼或暫存器類型，呈現相應的調性樣式及總長度摘要。
 */
export const RangeSummary: React.FC<RangeSummaryProps> = ({
  startAddress,
  count,
  dataType,
  protocol = 'modbus_tcp',
}) => {
  const { t } = useTranslation('workbench-v2');
  const stride = dataTypeWidth(dataType);
  const totalRegisters = count * stride;
  const startAddr = startAddress.trim();
  const isAddressValid = startAddr.length > 0 && addressParser.validate(startAddr, protocol).valid;
  const endAddr = isAddressValid
    ? addressParser.offset(startAddr, Math.max(0, count - 1) * stride, protocol)
    : '';

  // 推斷功能碼類型與色彩調性
  const getFuncMeta = (addr: string, proto: ProtocolType) => {
    const info = addressParser.getAreaInfo(addr, proto);

    if (proto === 'fatek_fbs') {
      return {
        label: t('step2.summary.function.fatek', {
          defaultValue: `FATEK ${info.typeLabel}`,
          label: info.typeLabel,
        }),
        bg: 'bg-emerald-500/10',
        border: 'border-emerald-500/20',
        text: 'text-emerald-400',
      };
    }

    if (proto === 'mc_3e') {
      return {
        label: t('step2.summary.function.mc', {
          defaultValue: `MC 3E ${info.typeLabel}`,
          label: info.typeLabel,
        }),
        bg: 'bg-indigo-500/10',
        border: 'border-indigo-500/20',
        text: 'text-indigo-400',
      };
    }

    if (info.typeLabel === 'coil') {
      return {
        label: t('step2.summary.function.coil', 'Coil (0x)'),
        bg: 'bg-amber-500/10',
        border: 'border-amber-500/20',
        text: 'text-amber-400',
      };
    }
    if (info.typeLabel === 'discrete_input') {
      return {
        label: t('step2.summary.function.discrete_input', 'Discrete Input (1x)'),
        bg: 'bg-slate-500/10',
        border: 'border-slate-500/20',
        text: 'text-slate-400',
      };
    }
    if (info.typeLabel === 'input_register') {
      return {
        label: t('step2.summary.function.input_register', 'Input Register (3x)'),
        bg: 'bg-cyan-500/10',
        border: 'border-cyan-500/20',
        text: 'text-cyan-400',
      };
    }
    return {
      label: t('step2.summary.function.holding_register', 'Holding Register (4x)'),
      bg: 'bg-blue-500/10',
      border: 'border-blue-500/20',
      text: 'text-blue-400',
    };
  };

  const meta = isAddressValid
    ? getFuncMeta(startAddr, protocol)
    : {
        label: t('step2.summary.invalid_address', '位址不符合所選協議格式。'),
        bg: 'bg-red-500/10',
        border: 'border-red-500/20',
        text: 'text-red-300',
      };

  return (
    <div
      className="flex flex-wrap items-center gap-3 p-3 rounded-lg bg-slate-900/60 border border-slate-800/80 text-xs text-slate-300 font-mono"
      aria-invalid={!isAddressValid}
      data-invalid={isAddressValid ? 'false' : 'true'}
      data-testid="range-summary"
    >
      {/* 範圍標籤 */}
      <span className="font-semibold text-slate-400">
        {t('step2.summary.range_label', 'Address range:')}
      </span>
      <span className="px-2 py-0.5 rounded bg-slate-800 border border-slate-700 text-slate-100">
        {startAddr || '—'} ~ {endAddr || '—'}
      </span>

      {/* 功能碼標籤 */}
      <span
        className={`px-2 py-0.5 rounded border ${meta.bg} ${meta.border} ${meta.text}`}
        data-testid="function-code-badge"
      >
        {meta.label}
      </span>

      {/* 總暫存器數說明 */}
      <span className="text-slate-400 ml-auto sm:ml-0">
        {t('step2.summary.register_count', {
          defaultValue: `${totalRegisters} registers (length ${stride})`,
          count: totalRegisters,
          stride,
        })}
      </span>

      {!isAddressValid && (
        <span className="basis-full text-red-300" data-testid="range-summary-error">
          {t('step2.summary.invalid_address', '位址不符合所選協議格式。')}
        </span>
      )}
    </div>
  );
};
