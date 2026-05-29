import * as React from 'react';
import { Icon } from './Icon';

export interface FieldProps {
  label: string;
  hint?: string;
  error?: string;
  required?: boolean;
  children: React.ReactNode;
  className?: string;
}

/**
 * Workbench V2 統一表單欄位包裝元件
 * 
 * 落地設計決策：「共用元件移植：1:1 對應原型 + TypeScript 化」
 * required 時出現紅色 * 標記，有 error 時會隱藏 hint 並優先顯示 alert 圖示與紅字。
 */
export const Field: React.FC<FieldProps> = ({
  label,
  hint,
  error,
  required,
  children,
  className = '',
}) => {
  return (
    <label className={`block ${className}`}>
      <span className={`label ${required ? 'after:content-["*"] after:text-red-400 after:ml-0.5' : ''}`}>
        {label}
      </span>
      {children}
      {error ? (
        <span className="mt-1 flex items-center gap-1 text-xs text-red-300">
          <Icon name="alert" className="w-3 h-3" />
          {error}
        </span>
      ) : hint ? (
        <span className="mt-1 block text-[11px] text-slate-500">{hint}</span>
      ) : null}
    </label>
  );
};
