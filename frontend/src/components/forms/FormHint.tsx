import type { ReactNode } from 'react';
import { designSystem } from '../../styles/designSystem';

export interface FormHintProps {
  /** 提示訊息 */
  children: ReactNode;
  /** 額外的 className */
  className?: string;
  /** id 屬性（用於 aria-describedby） */
  id?: string;
}

/**
 * 統一的表單提示訊息元件
 * 使用 designSystem.forms.hint.* 樣式
 */
export default function FormHint({
  children,
  className = '',
  id,
}: FormHintProps) {
  return (
    <p id={id} className={`${designSystem.forms.hint.base} ${className}`}>
      {children}
    </p>
  );
}
