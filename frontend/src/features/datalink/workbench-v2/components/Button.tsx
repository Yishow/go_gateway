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
    'inline-flex items-center justify-center gap-2 font-medium rounded-md transition-colors duration-150 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--gw-focus)] disabled:cursor-not-allowed disabled:bg-[var(--gw-ground)] disabled:text-[var(--gw-muted)] disabled:border-[var(--gw-border)]';
  
  const sizes = {
    sm: 'min-h-8 px-2.5 py-1.5 text-xs',
    md: 'min-h-10 px-4 py-2 text-sm',
    lg: 'min-h-10 px-5 py-2.5 text-sm',
  };

  const variants = {
    primary: 'bg-[var(--gw-primary)] text-[var(--gw-on-primary)] enabled:hover:bg-[var(--gw-primary-hover)]',
    secondary: 'bg-[var(--gw-surface)] text-[var(--gw-text)] enabled:hover:bg-[var(--gw-active)] border border-[var(--gw-control)]',
    ghost: 'text-[var(--gw-text)] enabled:hover:bg-[var(--gw-active)]',
    danger: 'border border-[var(--gw-danger-border)] bg-[var(--gw-danger-bg)] text-[var(--gw-danger-fg)]',
    success: 'border border-[var(--gw-success-border)] bg-[var(--gw-success-bg)] text-[var(--gw-success-fg)]',
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
