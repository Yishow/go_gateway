import type { ReactNode } from 'react';
import { AlertCircle } from 'lucide-react';
import { designSystem } from '../../styles/designSystem';

export interface FormErrorProps {
  /** 錯誤訊息 */
  children: ReactNode;
  /** 是否顯示圖示 */
  showIcon?: boolean;
  /** 額外的 className */
  className?: string;
  /** id 屬性（用於 aria-describedby） */
  id?: string;
}

/**
 * 統一的表單錯誤訊息元件
 * 使用 designSystem.forms.error.* 樣式
 */
export default function FormError({
  children,
  showIcon = true,
  className = '',
  id,
}: FormErrorProps) {
  return (
    <p id={id} className={`${designSystem.forms.error.base} ${className}`}>
      {showIcon && (
        <AlertCircle className={designSystem.forms.error.icon} aria-hidden="true" />
      )}
      {children}
    </p>
  );
}
