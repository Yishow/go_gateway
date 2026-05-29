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
 * 支援 tone (neutral/info/success/warning/error/draft)，info 與 success 帶 pulse-dot 動態效果。
 */
export const StatusChip: React.FC<StatusChipProps> = ({
  tone = 'neutral',
  dot = true,
  children,
  className = '',
}) => {
  const tones = {
    neutral: 'border-slate-700/70 bg-slate-800/60 text-slate-300',
    info: 'border-blue-500/30 bg-blue-500/10 text-blue-200',
    success: 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200',
    warning: 'border-amber-500/30 bg-amber-500/10 text-amber-200',
    error: 'border-red-500/30 bg-red-500/10 text-red-200',
    draft: 'border-slate-500/30 bg-slate-700/40 text-slate-300',
  };

  const dots = {
    neutral: 'bg-slate-400',
    info: 'bg-blue-400',
    success: 'bg-emerald-400',
    warning: 'bg-amber-400',
    error: 'bg-red-400',
    draft: 'bg-slate-400',
  };

  const isPulse = tone === 'info' || tone === 'success';

  return (
    <span className={`chip ${tones[tone]} ${className}`}>
      {dot && (
        <span
          className={`w-1.5 h-1.5 rounded-full ${dots[tone]} ${
            isPulse ? 'pulse-dot' : ''
          }`}
        />
      )}
      {children}
    </span>
  );
};
