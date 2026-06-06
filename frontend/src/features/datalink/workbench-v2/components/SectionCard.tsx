import * as React from 'react';

export interface SectionCardProps {
  title?: string;
  subtitle?: string;
  icon?: React.ReactNode;
  aside?: React.ReactNode;
  children: React.ReactNode;
  className?: string;
  contentClassName?: string;
}

/**
 * Workbench V2 統一區塊卡片元件
 * 
 * 落地設計決策：「共用元件移植：1:1 對應原型 + TypeScript 化」
 * 支援 title、subtitle、icon、aside 以及自訂的 class。
 * 當 title 或 aside 任一存在時會自動渲染 header。
 */
export const SectionCard: React.FC<SectionCardProps> = ({
  title,
  subtitle,
  icon,
  aside,
  children,
  className = '',
  contentClassName = '',
}) => {
  const showHeader = !!title || !!aside;
  const titleId = React.useId();

  return (
    <section
      role={title ? 'region' : undefined}
      aria-labelledby={title ? titleId : undefined}
      className={`wbv2-surface rounded-2xl border border-slate-700/60 bg-slate-900/60 backdrop-blur-sm shadow-xl shadow-black/30 ${className}`}
    >
      {showHeader && (
        <header className="flex items-start justify-between gap-4 border-b border-slate-700/60 px-5 py-4">
          <div className="flex items-start gap-3">
            {icon && (
              <div className="grid place-items-center w-9 h-9 rounded-lg bg-blue-500/10 border border-blue-500/20 text-blue-300">
                {icon}
              </div>
            )}
            <div>
              {title && (
                <h3 id={titleId} className="text-sm font-semibold text-slate-100 leading-tight">
                  {title}
                </h3>
              )}
              {subtitle && (
                <p className="text-xs text-slate-400 mt-0.5">{subtitle}</p>
              )}
            </div>
          </div>
          {aside && <div className="flex items-center gap-2">{aside}</div>}
        </header>
      )}
      <div className={`p-5 ${contentClassName}`}>{children}</div>
    </section>
  );
};
