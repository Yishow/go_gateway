# 前端 UI/UX 審查報告

審查日期：2025-02-25  
範圍：`frontend/src` 版面、元件、無障礙、主題與互動

---

## 一、整體優點

### 1. 無障礙（a11y）亮點
- **DatalinkLayout**：具備「跳過至主內容」連結（`sr-only` + `focus:not-sr-only`）、`main` 設 `id="datalink-main-content"` 與 `tabIndex={-1}`，側欄開關與通知按鈕有 `aria-label`，導航使用 `focus-visible:ring-2`。
- **SmartDashboardHeader**：搜尋框有 `sr-only` 的 `<label htmlFor="dashboard-search">`，圖示按鈕皆有 `aria-label`，按鈕有 `focus-visible:ring`。
- **Toast**：錯誤用 `role="alert"`、其餘用 `role="status"`，容器有 `aria-live="polite"`、`aria-atomic="false"`，關閉按鈕有 `aria-label="關閉通知"`。
- **ConfirmDialog / SlidePanel**：`role="dialog"`、`aria-modal="true"`、`aria-labelledby` / `aria-describedby` 正確。
- **index.css**：支援 `prefers-reduced-motion: reduce` 關閉動畫、`color-scheme` 與 `html/body` 高度設定得當。

### 2. 設計與主題
- **tokens.ts**：色彩、間距、圓角、陰影、字體等設計 token 集中管理，利於一致性。
- **ThemeContext**：尊重系統 `prefers-color-scheme` 並可寫入 localStorage，切換流暢。
- **Tailwind + 深色模式**：Layout / DatalinkLayout 的 sidebar、header、main 區塊皆有 light/dark 樣式，視覺一致。

### 3. 互動與回饋
- **Toast**：依類型區分樣式與圖示，具進度條與 hover 暫停，關閉按鈕明確。
- **表單與按鈕**：多數可聚焦控制項使用 `focus-visible:outline-none focus-visible:ring-2`，鍵盤可操作。
- **DeviceTreeNav、QuickActions、MemoryGrid** 等具 `aria-label` / `aria-pressed` / `role`，語意清楚。

### 4. 國際化與語意
- Datalink 區塊使用 `react-i18next`（`t('nav.*')`、`t('layout.*')` 等），利於多語與維護。
- 區塊標題與流程狀態使用 `aria-label` / `aria-labelledby` 綁定 i18n 字串，螢幕閱讀器友善。

---

## 二、待改進項目

### 1. Layout.tsx（舊版協議測試工具版面）

| 項目 | 現狀 | 建議 |
|------|------|------|
| 跳過主內容 | 無 | 在 `<main>` 前加入「跳過至主內容」連結，並給 main 設 `id`，與 DatalinkLayout 一致。 |
| 選單按鈕 | 僅 `focus:outline-none`，無 `aria-label` | 加上 `aria-label`（如「開啟選單」/「關閉選單」），並改用 `focus-visible:ring-2` 取代單純 outline-none。 |
| 導航區塊 | `<nav>` 無 `aria-label` | 加上 `aria-label="主要導航"` 或等同 i18n key。 |
| 主內容區 | `<main>` 無 id | 設 `id="main-content"`（或與 skip 連結對應的 id），供 skip link 與錨點使用。 |
| 導航連結 | 無 focus-visible 環 | 與 DatalinkLayout 一致，為導航連結加上 `focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2`。 |

### 2. 主題切換按鈕（ThemeToggle.tsx）

- **aria-label** 目前為英文 `"Toggle theme"`，若產品以繁中為主，建議改為 i18n（例如 `t('layout.toggleTheme')`），並在 `frontend/src/i18n/locales` 補上對應 key。

### 3. 字體與 index.css

- `index.css` 中 `code` 使用 `source-code-pro, Menlo, Monaco, Consolas, 'Courier New'`，而 `@import` 與 tokens 為 **Inter** + **JetBrains Mono**。若希望程式碼統一為 JetBrains Mono，建議將 `code` 的 `font-family` 改為與 tokens 一致（例如 `var(--font-mono)` 或 JetBrains Mono）。

### 4. Toast 關閉按鈕文案

- 關閉按鈕為 `aria-label="關閉通知"`（繁中），若其他 UI 已全面 i18n，可改為 `t('common.closeNotification')` 等 key，以支援多語與一致用詞。

### 5. 行動版側欄遮罩（Layout.tsx）

- 遮罩為 `<div onClick={...}>`，語意上為「關閉側欄」的控件。可考慮改為 `<button type="button" aria-label="關閉選單">` 並在點擊時關閉側欄，以提升鍵盤與無障礙一致性；若維持 div，建議加上 `aria-hidden="true"` 並確保焦點不會落入遮罩（例如用 focus trap 在側欄內）。

---

## 三、建議優先順序

1. **高**：Layout.tsx 的 skip link、main 的 id、選單按鈕的 aria-label 與 focus-visible。
2. **中**：ThemeToggle 與 Toast 關閉的 i18n、Layout 的 nav `aria-label` 與導航連結 focus-visible。
3. **低**：index.css 的 code 字體統一、遮罩改為 button 或補 aria-hidden。

---

## 四、小結

- **Datalink 相關版面**（DatalinkLayout、SmartDashboard、表單、對話框、Toast）在無障礙與互動上已達一定水準，skip link、ARIA、focus-visible、role 使用正確。
- **Layout.tsx** 為主要短板：缺少 skip link、main id、選單與導航的 aria 與 focus-visible，建議依上表補齊。
- 其餘為 i18n 與字體一致性的小調整，不影響基本可用性，但能提升體驗與維護性。

完成上述高優先項後，整體可符合常見的 Web 無障礙與 UX 最佳實踐。
