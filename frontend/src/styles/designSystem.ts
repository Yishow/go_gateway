/**
 * 統一設計系統規範
 * 
 * 本檔案整合 `tokens.ts`、`styleTokens.ts` 與 Tailwind 類別，作為專案單一可信來源。
 * 後續 UI 變更應優先使用此設計系統，而非直接手寫 Tailwind class。
 * 
 * 使用方式：
 * - 色彩、間距、圓角、陰影：使用 `designSystem.tokens.*`
 * - 元件樣式類別：使用 `designSystem.components.*`
 * - 表單規範：參考 `designSystem.forms.*`
 */

import { tokens } from './tokens';

export const designSystem = {
  /**
   * 設計 tokens（色彩、間距、圓角、陰影、字體）
   * 與 `tokens.ts` 對齊，作為單一來源
   */
  tokens,

  /**
   * 元件樣式類別（Tailwind 組合）
   * 提供可重用的元件樣式組合，避免重複手寫 class
   */
  components: {
    /**
     * 按鈕樣式
     * 優先使用 `components/ui/button.tsx`，此處提供直接 Tailwind 組合（用於非 Button 元件場景）
     */
    button: {
      primary: 'inline-flex items-center justify-center px-4 py-2 bg-blue-500 text-white text-sm font-medium rounded-lg hover:bg-blue-600 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed',
      secondary: 'inline-flex items-center justify-center px-4 py-2 bg-slate-700 text-slate-200 text-sm font-medium rounded-lg hover:bg-slate-600 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed',
      danger: 'inline-flex items-center justify-center px-4 py-2 bg-red-500 text-white text-sm font-medium rounded-lg hover:bg-red-600 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-red-500 focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed',
      ghost: 'inline-flex items-center justify-center px-4 py-2 text-slate-300 text-sm font-medium rounded-lg hover:bg-slate-800/50 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2',
      icon: 'inline-flex items-center justify-center p-2 text-slate-400 rounded-lg hover:bg-slate-800/50 hover:text-white transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2',
    },

    /**
     * 卡片與容器
     */
    card: {
      base: 'rounded-2xl border border-white/10 bg-slate-900/70 p-6 shadow-xl shadow-black/20 transition-colors',
      surface: 'rounded-xl border border-slate-700 bg-slate-800/50 p-4 shadow-lg',
      panel: 'rounded-lg border border-white/5 bg-slate-900/50 p-4',
    },

    /**
     * Section Header
     */
    sectionHeader: {
      base: 'flex items-center justify-between pb-4 border-b border-white/10',
      title: 'text-sm font-semibold uppercase tracking-[0.2em] text-slate-400',
      subtitle: 'text-xs text-slate-500 mt-1',
    },

    /**
     * Badge
     */
    badge: {
      default: 'inline-flex items-center rounded-full border border-slate-700/50 bg-slate-800/50 px-3 py-1 text-xs font-medium text-slate-300',
      success: 'inline-flex items-center rounded-full border border-emerald-500/30 bg-emerald-500/10 px-3 py-1 text-xs font-medium text-emerald-200',
      warning: 'inline-flex items-center rounded-full border border-amber-500/30 bg-amber-500/10 px-3 py-1 text-xs font-medium text-amber-200',
      error: 'inline-flex items-center rounded-full border border-red-500/30 bg-red-500/10 px-3 py-1 text-xs font-medium text-red-200',
      info: 'inline-flex items-center rounded-full border border-blue-500/30 bg-blue-500/10 px-3 py-1 text-xs font-medium text-blue-200',
    },
  },

  /**
   * 表單規範
   * 定義統一的 label、placeholder、錯誤訊息、disabled/loading、autocomplete/name/inputmode 策略
   */
  forms: {
    /**
     * 輸入框基礎樣式
     */
    input: {
      base: 'block w-full px-3 py-2 text-sm border border-slate-700/50 rounded-lg bg-slate-900/50 text-slate-200 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors',
      error: 'border-red-500/50 focus:ring-red-500 focus:border-red-500',
    },

    /**
     * Label 樣式
     */
    label: {
      base: 'block text-sm font-medium text-slate-300 mb-1.5',
      required: 'after:content-["*"] after:ml-0.5 after:text-red-400',
    },

    /**
     * 錯誤訊息樣式
     */
    error: {
      base: 'mt-1.5 text-xs text-red-400 flex items-center gap-1',
      icon: 'h-3.5 w-3.5',
    },

    /**
     * 成功訊息樣式
     */
    success: {
      base: 'mt-1.5 text-xs text-emerald-400 flex items-center gap-1',
      icon: 'h-3.5 w-3.5',
    },

    /**
     * 提示訊息樣式
     */
    hint: {
      base: 'mt-1.5 text-xs text-slate-500',
    },

    /**
     * autocomplete 屬性對應表
     * 參考 MDN autocomplete 規範：https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/autocomplete
     */
    autocomplete: {
      // 文字輸入
      text: 'off', // 一般文字輸入預設關閉 autocomplete
      name: 'name', // 姓名
      organization: 'organization', // 組織名稱
      // 數字輸入
      number: 'off', // 數字輸入預設關閉
      port: 'off', // 連接埠
      address: 'address-line1', // 地址
      // 其他
      email: 'email',
      url: 'url',
      tel: 'tel',
    },

    /**
     * inputmode 屬性對應表
     * 參考 MDN inputmode 規範：https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/inputmode
     */
    inputmode: {
      text: 'text', // 一般文字
      numeric: 'numeric', // 純數字（0-9）
      decimal: 'decimal', // 小數（含 .）
      tel: 'tel', // 電話號碼
      email: 'email', // 電子郵件
      url: 'url', // URL
      search: 'search', // 搜尋
      none: 'none', // 無虛擬鍵盤
    },

    /**
     * name 屬性命名規範
     * 使用 kebab-case，語意明確
     */
    name: {
      deviceName: 'device-name',
      deviceProtocol: 'device-protocol',
      deviceAddress: 'device-address',
      devicePort: 'device-port',
      tagKey: 'tag-key',
      tagDisplayName: 'tag-display-name',
      tagUnit: 'tag-unit',
      tagDescription: 'tag-description',
      pointAddress: 'point-address',
      pointDataType: 'point-data-type',
      modbusRegister: 'modbus-register',
      modbusPort: 'modbus-port',
      searchQuery: 'search-query',
      templateName: 'template-name',
      batchNamePrefix: 'batch-name-prefix',
    },
  },

  /**
   * Microcopy 規範
   * 統一使用繁中語氣、`…`、按鈕命名、狀態回饋語氣
   */
  microcopy: {
    /**
     * 載入狀態
     */
    loading: {
      default: '載入中…',
      saving: '儲存中…',
      deleting: '刪除中…',
      connecting: '連線中…',
      validating: '驗證中…',
      committing: '提交中…',
    },

    /**
     * 按鈕文字
     */
    button: {
      save: '儲存',
      cancel: '取消',
      delete: '刪除',
      edit: '編輯',
      create: '建立',
      update: '更新',
      submit: '提交',
      confirm: '確認',
      close: '關閉',
      back: '返回',
      next: '下一步',
      previous: '上一步',
      retry: '重試',
      reset: '重置',
    },

    /**
     * 狀態回饋
     */
    feedback: {
      success: {
        saved: '已儲存',
        created: '已建立',
        updated: '已更新',
        deleted: '已刪除',
        connected: '連線成功',
        validated: '驗證通過',
        committed: '提交成功',
      },
      error: {
        saveFailed: '儲存失敗，請重試',
        createFailed: '建立失敗，請檢查輸入',
        updateFailed: '更新失敗，請重試',
        deleteFailed: '刪除失敗，請重試',
        connectionFailed: '連線失敗，請檢查設定',
        validationFailed: '驗證失敗，請檢查輸入',
        commitFailed: '提交失敗，請重試',
      },
      warning: {
        unsavedChanges: '您有未儲存的變更',
        conflict: '發生衝突，請檢查',
        timeout: '操作逾時，請重試',
      },
    },

    /**
     * 省略號規範
     * 統一使用 `…`（U+2026），而非 `...`（三個點）
     */
    ellipsis: '…',
  },
} as const;
