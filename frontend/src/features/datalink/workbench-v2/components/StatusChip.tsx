import * as React from 'react';

export interface StatusChipProps {
  tone?: 'neutral' | 'info' | 'success' | 'warning' | 'error' | 'draft';
  dot?: boolean;
  children: React.ReactNode;
  className?: string;
}

/**
 * Workbench V2 統一狀態徽章元件
 * 
 * 落地設計決策：「共用元件移植：1:1 對應原型 + TypeScript 化」
 * 支援 tone (neutral/info/success/warning/error/draft)，狀態以文字與圖示呈現，不依賴動畫。
 */
export const StatusChip: React.FC<StatusChipProps> = ({
  tone = 'neutral',
  dot = true,
  children,
  className = '',
}) => {
  const tones = {
    neutral: 'border-[var(--gw-border)] bg-[var(--gw-surface)] text-[var(--gw-text)]',
    info: 'border-[var(--gw-focus)] bg-[var(--gw-active)] text-[var(--gw-accent)]',
    success: 'border-[var(--gw-success-border)] bg-[var(--gw-success-bg)] text-[var(--gw-success-fg)]',
    warning: 'border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] text-[var(--gw-warning-fg)]',
    error: 'border-[var(--gw-danger-border)] bg-[var(--gw-danger-bg)] text-[var(--gw-danger-fg)]',
    draft: 'border-[var(--gw-border)] bg-[var(--gw-surface)] text-[var(--gw-text)]',
  };


  return (
    <span className={`chip ${tones[tone]} ${className}`}>
      {dot && (
        <span
          className="w-1.5 h-1.5 rounded-full bg-current"
        />
      )}
      {children}
    </span>
  );
};
