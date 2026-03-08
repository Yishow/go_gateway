import type { ReactNode } from 'react';
import { CheckCircle2 } from 'lucide-react';
import { designSystem } from '../../styles/designSystem';

export interface FormSuccessProps {
  /** 成功訊息 */
  children: ReactNode;
  /** 是否顯示圖示 */
  showIcon?: boolean;
  /** 額外的 className */
  className?: string;
  /** id 屬性（用於 aria-describedby） */
  id?: string;
}

/**
 * 統一的表單成功訊息元件
 * 使用 designSystem.forms.success.* 樣式
 */
export default function FormSuccess({
  children,
  showIcon = true,
  className = '',
  id,
}: FormSuccessProps) {
  return (
    <p id={id} className={`${designSystem.forms.success.base} ${className}`}>
      {showIcon && (
        <CheckCircle2 className={designSystem.forms.success.icon} aria-hidden="true" />
      )}
      {children}
    </p>
  );
}
