import * as React from 'react';
import type { Rule } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { Field } from '../../components/Field';
import { Input, Select } from '../../components/inputs';

export interface ScaleSectionProps {
  ruleId: string;
  multiplier: number;
  offset: number;
  dataFormat: Rule['data_format'];
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

/**
 * 數值縮放與 byte order 格式化區塊元件
 * 
 * 落地設計決策：「線性轉換預設值繼承到 Step 3」
 * 提供展開式 detalles 區塊，用以設定乘數、偏移量以及 Modbus Word/Byte swap 順序。
 */
export const ScaleSection: React.FC<ScaleSectionProps> = ({
  ruleId,
  multiplier,
  offset,
  dataFormat,
  dispatch,
}) => {
  return (
    <details className="group border-t border-slate-800/60 pt-4" data-testid="scale-section">
      <summary className="flex items-center gap-2 cursor-pointer font-medium text-xs text-slate-400 select-none hover:text-slate-200 py-1.5 focus:outline-none">
        <svg
          viewBox="0 0 24 24"
          className="w-4 h-4 transition-transform group-open:rotate-90 text-slate-500"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.75"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M9 18l6-6-6-6" />
        </svg>
        <span>數值線性轉換與格式化 (選填)</span>
      </summary>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-3 pl-6">
        <Field label="乘數 (Multiplier)" hint="例如 0.1">
          <Input
            type="number"
            step="any"
            value={multiplier}
            onChange={(e) => {
              const val = parseFloat(e.target.value);
              dispatch({
                type: 'updateRule',
                ruleId,
                patch: { scale_multiplier: isNaN(val) ? 1 : val },
              });
            }}
            data-testid="scale-multiplier-input"
          />
        </Field>

        <Field label="偏移量 (Offset)" hint="例如 -10">
          <Input
            type="number"
            step="any"
            value={offset}
            onChange={(e) => {
              const val = parseFloat(e.target.value);
              dispatch({
                type: 'updateRule',
                ruleId,
                patch: { scale_offset: isNaN(val) ? 0 : val },
              });
            }}
            data-testid="scale-offset-input"
          />
        </Field>

        <Field label="位元組順序 (Byte Order)" hint="決定四位元組數值排布">
          <Select
            value={dataFormat}
            onChange={(e) => {
              dispatch({
                type: 'updateRule',
                ruleId,
                patch: { data_format: e.target.value as Rule['data_format'] },
              });
            }}
            data-testid="scale-dataformat-select"
          >
            <option value="">無 (預設)</option>
            <option value="ABCD">ABCD (Big Endian)</option>
            <option value="BADC">BADC (Byte Swap)</option>
            <option value="CDAB">CDAB (Word Swap)</option>
            <option value="DCBA">DCBA (Little Endian)</option>
          </Select>
        </Field>
      </div>
    </details>
  );
};
