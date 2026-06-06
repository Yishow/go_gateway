import * as React from 'react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger' | 'success';
  size?: 'sm' | 'md' | 'lg';
  icon?: React.ReactNode;
}

/**
 * Workbench V2 統一按鈕元件
 * 
 * 落地設計決策：「共用元件移植：1:1 對應原型 + TypeScript 化」
 * 支援 variant (primary/secondary/ghost/danger/success) 與 size (sm/md/lg)。
 */
export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  icon,
  children,
  className = '',
  ...props
}) => {
  const base =
    'inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all duration-200 active:scale-[0.98] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-300/60 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950 disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100';
  
  const sizes = {
    sm: 'px-2.5 py-1.5 text-xs',
    md: 'px-4 py-2 text-sm',
    lg: 'px-5 py-2.5 text-sm',
  };

  const variants = {
    primary: 'bg-cyan-400 text-slate-950 hover:bg-cyan-300 shadow-[0_14px_34px_rgba(34,211,238,0.18)]',
    secondary: 'bg-slate-900/80 text-slate-200 hover:bg-slate-800 border border-slate-700/60 shadow-[0_1px_0_rgba(255,255,255,0.04)_inset]',
    ghost: 'text-slate-300 hover:bg-slate-800/70 hover:text-slate-50',
    danger: 'bg-rose-600 text-white hover:bg-rose-500 shadow-[0_14px_34px_rgba(244,63,94,0.14)]',
    success: 'bg-emerald-500 text-slate-950 hover:bg-emerald-400 shadow-[0_14px_34px_rgba(16,185,129,0.16)]',
  };

  return (
    <button
      {...props}
      className={`${base} ${sizes[size]} ${variants[variant]} ${className}`}
    >
      {icon && <span className="inline-flex items-center">{icon}</span>}
      {children}
    </button>
  );
};
