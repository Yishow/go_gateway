# C 方案實作步驟與驗收清單 (Tasklist)

## 階段 1：基礎 Token 與樣式導入
**目標**：將 C 方案視覺規範抽象為統一的 Style Tokens，供各頁面引用，不更動邏輯。
- [ ] **任務 1.1**：建立 `styleTokens.ts`，定義共享 CSS classes 或 Tailwind 組合。
  - *驗收*：檔案存在且無 Type Error。
- [ ] **任務 1.2**：全域 CSS 或 Tailwind Config 檢視（若需擴充自訂色或字型設定）。
  - *驗收*：跑 `pnpm run build` 或 `tsc` 正常通過。

## 階段 2：Gateway Entry Page (入口頁) 翻新
**目標**：套用新卡片設計與字級間距，優化「快速設定」與「專家模式」的入口視覺。
- [ ] **任務 2.1**：替換 `GatewayEntryPage.tsx` 中的卡片容器為 `styleTokens.card` 樣式。
- [ ] **任務 2.2**：更新標題與輔助文字為 `styleTokens.text.h1` / `body` / `caption`。
- [ ] **任務 2.3**：調整入口卡片的 Hover 動效，對齊「互動規範」。
  - *驗收*：UI 渲染正常，現有 e2e 測試（如路由跳轉）不中斷。

## 階段 3：Gateway Quick Setup Page (快速設定頁) 翻新
**目標**：讓表單體驗更清晰，層次分明，符合工業級高對比要求。
- [ ] **任務 3.1**：表單輸入框與下拉選單套用 `styleTokens.form` 相關類別。
- [ ] **任務 3.2**：優化連線狀態指示器（Status），套用 `styleTokens.status` 色系。
- [ ] **任務 3.3**：主要與次要按鈕更新為 `styleTokens.button.primary` / `secondary`。
  - *驗收*：表單提交功能正常，連線測試按鈕狀態反饋正確。

## 階段 4：Gateway Expert Workbench Page (專家工作台頁) 翻新
**目標**：重構資料密集型介面的間距與表格/列表樣式，對齊 Datalink 操作體驗。
- [ ] **任務 4.1**：頂部控制列（Control Bar）與工具區套用新間距 `styleTokens.layout.sectionGap`。
- [ ] **任務 4.2**：數據列表或樹狀結構 Hover 效果套用 `styleTokens.interaction.hoverRow`。
- [ ] **任務 4.3**：更新配置卡片（Protocol, Device, Points）的邊框與圓角。
  - *驗收*：複雜組件的滾動與點擊互動正常，不影響任何狀態管理邏輯。

## 階段 5：全面測試與回歸驗證
- [ ] **任務 5.1**：執行所有單元測試 `pnpm test` 或相關 vitest 指令。
- [ ] **任務 5.2**：執行所有 E2E 測試 `pnpm playwright test`。
- [ ] **任務 5.3**：手動檢查深色模式/淺色模式（若專案有支援）是否與新 Token 相容。
  - *驗收*：零迴歸錯誤，符合 API 契約與現有邏輯。
