/**
 * Gateway C-Style Refresh Tokens
 * 遵循 Intentional Minimalism 理念，為 Gateway 相關頁面提供統一的樣式類別組合。
 * 這些 Tokens 直接映射到 Tailwind CSS 工具類別，確保不改變既有邏輯即可套用新設計。
 */

export const gatewayStyleTokens = {
  // 色彩與文字排版
  text: {
    h1: 'text-2xl font-semibold text-slate-900 tracking-tight',
    h2: 'text-lg font-semibold text-slate-900',
    body: 'text-sm font-normal text-slate-700',
    caption: 'text-xs font-normal text-slate-500',
    label: 'text-sm font-medium text-slate-700',
  },

  // 頁面與區塊佈局
  layout: {
    pageContainer: 'p-8 max-w-7xl mx-auto space-y-6',
    sectionGap: 'space-y-6',
    itemGap: 'space-y-4',
  },

  // 卡片與容器
  card: {
    base: 'bg-white rounded-xl border border-slate-200 shadow-sm',
    hoverable: 'bg-white rounded-xl border border-slate-200 shadow-sm transition-shadow duration-200 hover:shadow-md cursor-pointer',
    padding: 'p-6',
    header: 'flex items-center justify-between mb-4',
  },

  // 按鈕
  button: {
    primary: 'inline-flex items-center justify-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2',
    secondary: 'inline-flex items-center justify-center px-4 py-2 bg-transparent text-slate-700 text-sm font-medium rounded-lg border border-slate-300 hover:bg-slate-50 transition-colors focus:outline-none focus:ring-2 focus:ring-slate-500 focus:ring-offset-2',
    danger: 'inline-flex items-center justify-center px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-lg hover:bg-red-700 transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2',
    icon: 'p-2 text-slate-500 hover:text-slate-700 hover:bg-slate-100 rounded-md transition-colors',
  },

  // 表單元件
  form: {
    input: 'block w-full px-3 py-2 text-sm border border-slate-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
    select: 'block w-full px-3 py-2 text-sm border border-slate-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
    group: 'space-y-1.5',
  },

  // 狀態與指示器
  status: {
    success: 'flex items-center text-sm text-emerald-600 font-medium',
    warning: 'flex items-center text-sm text-amber-600 font-medium',
    error: 'flex items-center text-sm text-red-600 font-medium',
    indicator: {
      success: 'w-2 h-2 rounded-full bg-emerald-500 mr-2',
      warning: 'w-2 h-2 rounded-full bg-amber-500 mr-2',
      error: 'w-2 h-2 rounded-full bg-red-500 mr-2',
      offline: 'w-2 h-2 rounded-full bg-slate-400 mr-2',
    }
  },

  // 互動效果
  interaction: {
    hoverRow: 'hover:bg-slate-50 transition-colors',
    focusRing: 'focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2',
  }
};
