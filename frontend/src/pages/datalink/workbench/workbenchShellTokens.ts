/**
 * Datalink Workbench 外殼視覺 token。
 *
 * 設計意圖（第二輪收斂）：
 * - **實心背景**取代多區塊各自漸層，降低長時盯屏疲勞（B2B 儀表板）。
 * - **單一邊線 + 輕陰影**，避免 `gradient + ring-inset + 重外陰影` 疊加造成的「數位炫光」感。
 * - 五區（Context / Rail / Main / Inspector / Summary）共用同一底層，再以內距與內容層次區分角色。
 */

/** 外殼區塊共用：頂欄、步驟軌、底欄、檢查面板容器。不含 padding，由各元件自行加上。 */
export const WB_SHELL_SURFACE =
  'rounded-2xl border border-slate-800/70 bg-slate-900/95 shadow-sm shadow-black/30';
