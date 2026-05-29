/**
 * Datalink Workbench V2 設計系統 Token
 * 
 * 落地設計決策：「Feature 目錄結構：frontend/src/features/datalink/workbench-v2/」
 * 提供 Dark Industrial Telemetry 風格的設計 Token。
 */

export const COLORS = {
  bg: '#0b1220',
  success: '#10b981',
  primary: '#3b82f6',
  warning: '#f59e0b',
  danger: '#ef4444',
  info: '#3b82f6',
  slate: {
    100: '#f1f5f9',
    200: '#e2e8f0',
    300: '#cbd5e1',
    400: '#94a6b8',
    500: '#64748b',
    600: '#475569',
    700: '#334155',
    800: '#1e293b',
    900: '#0f172a',
    950: '#020617',
  },
};

export const FONTS = {
  sans: "'Inter', system-ui, sans-serif",
  mono: "'JetBrains Mono', ui-monospace, monospace",
};

export const RADII = {
  none: '0px',
  sm: '0.125rem',
  default: '0.25rem',
  md: '0.375rem',
  lg: '0.5rem',
  xl: '0.75rem',
  '2xl': '1rem',
  '3xl': '1.5rem',
  full: '9999px',
};

export const SHADOWS = {
  sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
  default: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',
  md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)',
  lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)',
  xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
  '2xl': '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
  inner: 'inset 0 2px 4px 0 rgba(0, 0, 0, 0.05)',
  none: 'none',
};

export const ANIMATIONS = {
  pulseDot: 'pulseDot 1.6s ease-in-out infinite',
  sweepIn: 'sweepIn 240ms ease-out both',
  shimmer: 'shimmer 1.6s ease-in-out infinite',
};
