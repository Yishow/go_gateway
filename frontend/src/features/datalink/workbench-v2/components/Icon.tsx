import * as React from 'react';

/**
 * 支援的圖示名稱
 */
export type IconName =
  | 'device'
  | 'rule'
  | 'map'
  | 'db'
  | 'check'
  | 'chevron'
  | 'plus'
  | 'play'
  | 'refresh'
  | 'cable'
  | 'sliders'
  | 'tag'
  | 'table'
  | 'alert'
  | 'info'
  | 'spark'
  | 'arrow'
  | 'close'
  | 'save'
  | 'eye'
  | 'bolt'
  | 'flow'
  | 'grid';

export interface IconProps {
  name: IconName;
  className?: string;
}

/**
 * Workbench V2 統一圖示元件
 * 
 * 落地設計決策：「共用元件移植：1:1 對應原型 + TypeScript 化」
 * 不使用外部 Icon 庫，採內建 23 種 SVG Path 字典並禁用 emoji。
 */
export const Icon: React.FC<IconProps> = ({ name, className = 'w-4 h-4' }) => {
  const paths: Record<IconName, React.ReactNode> = {
    device: (
      <>
        <rect x="4" y="3" width="16" height="18" rx="2" />
        <path d="M9 7h6M9 11h6M9 15h3" />
        <circle cx="17" cy="15" r="1" fill="currentColor" />
      </>
    ),
    rule: (
      <>
        <path d="M4 6h12M4 12h16M4 18h8" />
        <circle cx="20" cy="6" r="2" fill="currentColor" />
        <circle cx="14" cy="18" r="2" fill="currentColor" />
      </>
    ),
    map: (
      <>
        <path d="M3 6l6-2 6 2 6-2v14l-6 2-6-2-6 2V6z" />
        <path d="M9 4v16M15 6v16" />
      </>
    ),
    db: (
      <>
        <ellipse cx="12" cy="5" rx="8" ry="3" />
        <path d="M4 5v6c0 1.7 3.6 3 8 3s8-1.3 8-3V5" />
        <path d="M4 11v6c0 1.7 3.6 3 8 3s8-1.3 8-3v-6" />
      </>
    ),
    check: <path d="M4 12l5 5L20 6" />,
    chevron: <path d="M9 6l6 6-6 6" />,
    plus: (
      <>
        <path d="M12 5v14M5 12h14" />
      </>
    ),
    play: <path d="M6 4l14 8-14 8V4z" fill="currentColor" />,
    refresh: <path d="M4 4v6h6M20 20v-6h-6M5 13a8 8 0 0014.5 4.5M19 11a8 8 0 00-14.5-4.5" />,
    cable: (
      <>
        <path d="M5 8a3 3 0 016 0v8a3 3 0 006 0V8" />
        <path d="M5 4v8M19 12v8" />
      </>
    ),
    sliders: (
      <>
        <path d="M4 6h10M18 6h2M4 12h2M10 12h10M4 18h14M18 18h2" />
        <circle cx="16" cy="6" r="2" fill="currentColor" />
        <circle cx="8" cy="12" r="2" fill="currentColor" />
        <circle cx="16" cy="18" r="2" fill="currentColor" />
      </>
    ),
    tag: (
      <>
        <path d="M3 12V5a2 2 0 012-2h7l9 9-9 9-9-9z" />
        <circle cx="8" cy="8" r="1.5" fill="currentColor" />
      </>
    ),
    table: (
      <>
        <rect x="3" y="4" width="18" height="16" rx="1.5" />
        <path d="M3 10h18M3 16h18M10 4v16M16 4v16" />
      </>
    ),
    alert: (
      <>
        <path d="M12 9v4M12 17h.01" />
        <path d="M10.3 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.7 3.86a2 2 0 00-3.4 0z" />
      </>
    ),
    info: (
      <>
        <circle cx="12" cy="12" r="9" />
        <path d="M12 8h.01M11 12h1v5h1" />
      </>
    ),
    spark: (
      <>
        <path d="M5 12l4-4 3 6 4-9 3 7" />
      </>
    ),
    arrow: <path d="M5 12h14M13 5l7 7-7 7" />,
    close: <path d="M6 6l12 12M18 6L6 18" />,
    save: (
      <>
        <path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z" />
        <path d="M17 21v-8H7v8M7 3v5h8" />
      </>
    ),
    eye: (
      <>
        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8S1 12 1 12z" />
        <circle cx="12" cy="12" r="3" />
      </>
    ),
    bolt: <path d="M13 2L3 14h7l-1 8 10-12h-7l1-8z" fill="currentColor" />,
    flow: (
      <>
        <path d="M3 6h6l3 6 3-6h6" />
        <circle cx="3" cy="6" r="1.5" fill="currentColor" />
        <circle cx="21" cy="6" r="1.5" fill="currentColor" />
      </>
    ),
    grid: (
      <>
        <rect x="3" y="3" width="7" height="7" rx="1" />
        <rect x="14" y="3" width="7" height="7" rx="1" />
        <rect x="3" y="14" width="7" height="7" rx="1" />
        <rect x="14" y="14" width="7" height="7" rx="1" />
      </>
    ),
  };

  return (
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.75"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      {paths[name] ?? null}
    </svg>
  );
};
