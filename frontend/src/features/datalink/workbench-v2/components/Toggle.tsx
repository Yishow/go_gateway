import * as React from 'react';
import { useTranslation } from 'react-i18next';

export interface ToggleProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  label?: string;
  size?: 'sm' | 'md';
  disabled?: boolean;
  'aria-label'?: string;
}

/**
 * Workbench V2 統一開關切換元件
 * 
 * 落地設計決策：「共用元件移植：1:1 對應原型 + TypeScript 化」
 * 採用精確像素計算實作，並曝露 role="switch" 和 aria-checked 等 ARIA 屬性。
 */
export const Toggle: React.FC<ToggleProps> = ({
  checked,
  onChange,
  label,
  size = 'md',
  disabled = false,
  'aria-label': ariaLabel,
}) => {
  const { t } = useTranslation('workbench-v2');
  const dim = size === 'sm'
    ? { w: 32, h: 16, dot: 12, pad: 2 }
    : { w: 40, h: 20, dot: 16, pad: 2 };
  
  const tx = checked ? dim.w - dim.pad * 2 - dim.dot : 0;

  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      aria-label={ariaLabel || label || t('buttons.toggle')}
      disabled={disabled}
      onClick={() => {
        if (!disabled) {
          onChange(!checked);
        }
      }}
      className={`inline-flex items-center gap-2 group focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 rounded ${
        disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
      }`}
    >
      <span
        className={`relative inline-block rounded-full transition-colors duration-200 flex-shrink-0 ${
          checked ? 'bg-blue-600' : 'bg-slate-700'
        }`}
        style={{ width: dim.w, height: dim.h }}
      >
        <span
          className="absolute rounded-full bg-white shadow-sm transition-transform duration-200 will-change-transform"
          style={{
            width: dim.dot,
            height: dim.dot,
            top: dim.pad,
            left: dim.pad,
            transform: `translateX(${tx}px)`,
          }}
        />
      </span>
      {label && <span className="text-xs text-slate-300 select-none">{label}</span>}
    </button>
  );
};
