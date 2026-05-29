import * as React from 'react';

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {}

/**
 * 統一文字輸入框元件
 */
export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className = '', ...props }, ref) => {
    return (
      <input
        ref={ref}
        {...props}
        className={`w-full rounded-lg border border-slate-700/60 bg-slate-950/40 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 transition-colors focus:border-blue-500/70 focus:outline-none focus:ring-2 focus:ring-blue-500/20 ${className}`}
      />
    );
  }
);

Input.displayName = 'Input';

export interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {}

/**
 * 統一選擇器元件
 */
export const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
  ({ className = '', children, ...props }, ref) => {
    return (
      <div className="relative w-full">
        <select
          ref={ref}
          {...props}
          className={`w-full appearance-none rounded-lg border border-slate-700/60 bg-slate-950/40 px-3 py-2 pr-9 text-sm text-slate-100 transition-colors focus:border-blue-500/70 focus:outline-none focus:ring-2 focus:ring-blue-500/20 ${className}`}
        >
          {children}
        </select>
        <svg
          viewBox="0 0 24 24"
          className="pointer-events-none absolute right-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.75"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <path d="M6 9l6 6 6-6" />
        </svg>
      </div>
    );
  }
);

Select.displayName = 'Select';

export interface TextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {}

/**
 * 統一多行文字輸入框元件
 */
export const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className = '', ...props }, ref) => {
    return (
      <textarea
        ref={ref}
        {...props}
        className={`w-full rounded-lg border border-slate-700/60 bg-slate-950/40 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 transition-colors focus:border-blue-500/70 focus:outline-none focus:ring-2 focus:ring-blue-500/20 ${className}`}
      />
    );
  }
);

Textarea.displayName = 'Textarea';
