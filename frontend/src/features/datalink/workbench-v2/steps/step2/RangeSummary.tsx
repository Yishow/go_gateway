import * as React from 'react';
import { dataTypeWidth } from '../../state/sourceRule';

export interface RangeSummaryProps {
  startAddress: string;
  count: number;
  dataType: string;
}

/**
 * 規則位址範圍與功能碼摘要元件
 * 
 * 落地設計決策：「Modbus function code inference」
 * 依起點前綴推斷功能碼類型，呈現相應的調性樣式 (Amber, Slate, Cyan, Blue) 及總長度摘要。
 */
export const RangeSummary: React.FC<RangeSummaryProps> = ({
  startAddress,
  count,
  dataType,
}) => {
  const startNum = parseInt(startAddress.replace(/[^0-9]/g, ''), 10) || 40001;
  const stride = dataTypeWidth(dataType);
  const totalRegisters = count * stride;
  const endNum = startNum + Math.max(0, count - 1) * stride;

  // 推斷功能碼類型與色彩調性
  const getFuncMeta = (addr: string) => {
    const cleanAddr = addr.trim();
    if (cleanAddr.startsWith('0')) {
      return {
        label: 'Coil (0x)',
        bg: 'bg-amber-500/10',
        border: 'border-amber-500/20',
        text: 'text-amber-400',
      };
    }
    if (cleanAddr.startsWith('1')) {
      return {
        label: 'Discrete Input (1x)',
        bg: 'bg-slate-500/10',
        border: 'border-slate-500/20',
        text: 'text-slate-400',
      };
    }
    if (cleanAddr.startsWith('3')) {
      return {
        label: 'Input Register (3x)',
        bg: 'bg-cyan-500/10',
        border: 'border-cyan-500/20',
        text: 'text-cyan-400',
      };
    }
    return {
      label: 'Holding Register (4x)',
      bg: 'bg-blue-500/10',
      border: 'border-blue-500/20',
      text: 'text-blue-400',
    };
  };

  const meta = getFuncMeta(startAddress);

  return (
    <div
      className="flex flex-wrap items-center gap-3 p-3 rounded-lg bg-slate-900/60 border border-slate-800/80 text-xs text-slate-300 font-mono"
      data-testid="range-summary"
    >
      {/* 範圍標籤 */}
      <span className="font-semibold text-slate-400">位址範圍：</span>
      <span className="px-2 py-0.5 rounded bg-slate-800 border border-slate-700 text-slate-100">
        {startNum} → {endNum}
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
        共 <strong className="text-slate-200">{totalRegisters}</strong> 個暫存器 (長度 {stride})
      </span>
    </div>
  );
};
