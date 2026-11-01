import * as React from 'react';
import { Icon } from './Icon';
import { ParameterHelp } from './ParameterHelp';
import type { ParameterHelpKey } from './parameterHelpContent';

export interface FieldProps {
  label: string;
  hint?: string;
  help?: ParameterHelpKey;
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
  help,
  error,
  required,
  children,
  className = '',
}) => {
  const field = (
    <label className={`block ${help ? 'min-w-0 flex-1' : className}`}>
      <span className={`label ${required ? 'after:content-["*"] after:text-[var(--gw-danger-fg)] after:ml-0.5' : ''}`}>
        {label}
      </span>
      {children}
      {error ? (
        <span className="mt-1 flex items-center gap-1 text-xs text-[var(--gw-danger-fg)]">
          <Icon name="alert" className="w-3 h-3" />
          {error}
        </span>
      ) : hint ? (
        <span className="mt-1 block text-xs text-[var(--gw-muted)]">{hint}</span>
      ) : null}
    </label>
  );
  return help ? <div className={`flex min-w-0 items-start gap-2 ${className}`}>{field}<ParameterHelp help={help} /></div> : field;
};
