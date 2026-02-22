/**
 * SmartDashboard 共用 Tailwind className 常數
 *
 * 設計原則：
 * - base：不含顏色，可組合到任何變體
 * - btn*：完整按鈕類，直接使用
 * - input：輸入框
 * - pill*：小型狀態標籤（badge）
 * - card：卡片容器
 */

// ─── 按鈕基底（不含顏色） ────────────────────────────────────────────────
const btnBase =
  'inline-flex items-center gap-1.5 rounded-lg border px-3 py-1.5 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 disabled:cursor-not-allowed disabled:opacity-40';

// ─── 色彩變體按鈕 ────────────────────────────────────────────────────────
export const btn = {
  /** 天空藍：Auto 分配等主要工具操作 */
  sky: `${btnBase} border-sky-500/40 bg-sky-500/15 text-sky-100 hover:bg-sky-500/25 focus-visible:ring-sky-500`,
  /** 翠綠：套用到 Grid、Commit 等確認操作 */
  emerald: `${btnBase} border-emerald-500/40 bg-emerald-500/20 text-emerald-100 hover:bg-emerald-500/30 focus-visible:ring-emerald-500`,
  /** 靛藍：儲存模板等次要建立操作 */
  indigo: `${btnBase} border-indigo-500/40 bg-indigo-500/20 text-indigo-100 hover:bg-indigo-500/30 focus-visible:ring-indigo-500`,
  /** 藍色：工作台連結、套用 Tag 等導航操作 */
  blue: `${btnBase} border-blue-500/40 bg-blue-500/15 text-blue-100 hover:bg-blue-500/25 focus-visible:ring-blue-500`,
  /** 琥珀：衝突過濾、Retry 等警示操作 */
  amber: `${btnBase} border-amber-500/40 bg-amber-500/20 text-amber-100 hover:bg-amber-500/30 focus-visible:ring-amber-500`,
  /** 石板灰：中性次要操作 */
  slate: `${btnBase} border-slate-700 bg-slate-800/70 text-slate-200 hover:bg-slate-700 focus-visible:ring-blue-500`,
  /** 玫瑰紅：刪除、Rollback 等危險操作 */
  rose: `${btnBase} border-rose-500/40 bg-rose-500/15 text-rose-100 hover:bg-rose-500/25 focus-visible:ring-rose-500`,
} as const;

// ─── 輸入框 ──────────────────────────────────────────────────────────────
export const input =
  'w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500';

/** 單行等寬字體輸入框（Address、Tag Key 等） */
export const inputMono = `${input} font-mono`;

// ─── 小型狀態 Pill（badge） ──────────────────────────────────────────────
const pillBase = 'inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px]';

export const pill = {
  emerald: `${pillBase} border-emerald-500/30 bg-emerald-500/10 text-emerald-300`,
  rose:    `${pillBase} border-rose-500/30    bg-rose-500/10    text-rose-300`,
  amber:   `${pillBase} border-amber-500/30   bg-amber-500/10   text-amber-300`,
  blue:    `${pillBase} border-blue-500/30    bg-blue-500/10    text-blue-300`,
  slate:   `${pillBase} border-white/10       bg-slate-800/60   text-slate-300`,
} as const;

// ─── 卡片容器 ────────────────────────────────────────────────────────────
/** 面板內的深色子卡片 */
export const card = 'rounded-2xl border border-white/10 bg-slate-900/60 p-4';

/** 更淡的內嵌子區塊 */
export const cardInner = 'rounded-lg border border-white/10 bg-slate-800/40 p-3';

// ─── 區塊標題行 ──────────────────────────────────────────────────────────
/** 用於小節 label，如「批次命名前綴」上方的說明文字 */
export const sectionLabel = 'text-[11px] font-semibold uppercase tracking-wider text-slate-500';
