import type { ReactNode } from 'react';
import { designSystem } from '../../styles/designSystem';

export interface FormLabelProps {
  /** Label 文字 */
  children: ReactNode;
  /** 是否為必填欄位 */
  required?: boolean;
  /** 對應的 input id（用於 accessibility） */
  htmlFor?: string;
  /** 額外的 className */
  className?: string;
}

/**
 * 統一的表單 Label 元件
 * 使用 designSystem.forms.label.* 樣式
 */
export default function FormLabel({
  children,
  required = false,
  htmlFor,
  className = '',
}: FormLabelProps) {
  return (
    <label
      htmlFor={htmlFor}
      className={`${designSystem.forms.label.base} ${required ? designSystem.forms.label.required : ''} ${className}`}
    >
      {children}
    </label>
  );
}
