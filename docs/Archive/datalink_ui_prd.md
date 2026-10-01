# Datalink UI 改善 PRD

## 1. 背景與目標
目前 Datalink 介面已具備 Device/Point/Mapping/Tag 的核心操作，但主流程集中於 `SmartDashboardPage.tsx`、前後端 API 存在雙軌差異、SSE 與 Runtime 真值落差、測試與 i18n 一致性不足，導致導入與維運成本偏高。

本 PRD 目標：在不破壞既有資料模型前提下，建立可擴充、可測試、可觀測、可國際化的 Datalink UI 主流程。

## 2. 範圍
### In Scope
- Datalink 主流程資訊架構重整（導覽、頁面分區、任務導向流程）
- Smart Dashboard 操作與狀態體驗優化
- API 契約一致性對齊（以主線 Gin API 為準）
- 錯誤處理、空狀態、載入狀態、回饋訊息標準化
- RWD、A11y、i18n 補強
- 可量測 KPI 與驗收標準建立

### Out of Scope
- 新增全新資料來源協定（非 Modbus）
- 大幅改動後端領域模型（僅做契約對齊與補齊）

## 3. 目標使用者
- 系統整合工程師（建立 Device/Point/Mapping）
- 維運工程師（監看健康狀態、排查錯誤）
- 測試人員（驗證流程、回歸）

## 4. 核心問題（現況痛點）
1. **流程分散**：Wizard、Dashboard、Workbench 邏輯重疊，使用者不易判斷下一步。
2. **契約不一致**：前端呼叫 readiness/poll 類 API 時，主線能力與預期有落差。
3. **狀態不可預期**：Commit lifecycle 多為前端模擬，與後端真實執行結果關聯不足。
4. **可觀測性不足**：SSE 事件目前偏示意，難以作為操作決策依據。
5. **體驗不一致**：部分元件仍為直接 fetch，錯誤/重試/快取策略不一致。
6. **國際化不完整**：Smart Dashboard 區塊存在硬編碼字串。

## 5. 產品目標與成功指標
### 產品目標
- 使用者可在單一路徑完成「建置→驗證→提交→監看」閉環。
- UI 反映後端真實狀態（readiness、preview、polling、commit 結果）。
- 關鍵操作都有一致的錯誤與復原引導。

### KPI（上線後 4 週觀察）
- 新手完成首條 datalink 建立時間：**-30%**
- Mapping 驗證失敗後一次修復成功率：**+25%**
- 主要流程中斷率（取消/返回/卡住）：**-40%**
- 同一問題重複提問客服/開發：**-30%**
- e2e 關鍵路徑測試穩定度：**95%+**

## 6. 使用者旅程（To-Be）
1. 進入 Datalink 首頁 → 系統提示目前成熟度與待辦
2. 建立/選擇 Device → readiness 檢查
3. 建立 Point（單筆/批次）→ 即時格式與衝突檢查
4. 建立 Mapping + Preview（Gate C）→ 顯示可修復建議
5. Commit（含 dry-run 摘要）→ 顯示進度、成功/失敗與 rollback 指引
6. 進入 Runtime 監看頁 → 健康狀態、最近事件、資料延遲

## 7. Roadmap
### P0（必做）
- 主流程 IA 重整（任務導向）
- API 契約對齊（readiness/poll/preview）
- 主要元件錯誤/載入/空狀態標準化
- 關鍵流程測試修復（含 QueryClientProvider）

### P1（重要）
- Commit 與 Runtime 觀測面板聯動
- i18n 全面外掛字串
- 行動版流程優化（抽屜、分步）

### P2（加值）
- 互動引導（新手模式）
- 進階診斷報表與操作建議

## 8. 驗收條件（PRD 層）
- [ ] 任務流程可在 1 個主路徑完成（不需在多頁來回）
- [ ] UI 所有關鍵狀態皆有對應設計（loading/empty/error/success/partial）
- [ ] API 文件與前端 service 契約一致
- [ ] 主要頁面字串可切換中/英
- [ ] 關鍵 KPI 事件可上報並可追蹤

---

## 來源
- `docs/datalink_flow_gate_v1_plan.md`
- `docs/datalink_runtime_wbs.md`
- `docs/implement-datalink-ui-test-summary.md`
- `frontend/datalink_wireframe.html`
- `frontend/src/pages/datalink/SmartDashboardPage.tsx`
