# Datalink UI Spec（可交付設計與前端實作）

## 1. 資訊架構（IA）

## 1.1 導覽結構
- Datalink Overview（總覽）
- Device Setup（設備）
- Point Setup（點位）
- Mapping & Preview（映射與預覽）
- Commit & Runtime（提交與執行監看）
- Settings（系統設定）

## 1.2 主畫面區塊
1. **流程狀態條（頂部）**：顯示 Gate A/B/C/D 進度與阻塞原因
2. **左側任務導覽**：可切 section，顯示未完成項目數
3. **中央工作區**：根據 section 顯示表單/表格/預覽
4. **右側上下文面板**：目前資源摘要、錯誤清單、建議修復動作

## 2. 元件規格
### 2.1 FlowStatusBar
- Props: `steps`, `currentStep`, `blockedReasons[]`
- 狀態：done / active / blocked / pending
- 互動：點擊 step 可跳轉；blocked step 顯示 tooltip

### 2.2 DeviceReadinessCard
- 顯示 device 基本資訊、最近測試時間、readiness 結果
- CTA：`測試連線`、`啟用設備`
- 錯誤提示包含可執行修復（例如 Port、Unit ID）

### 2.3 PointBatchEditor
- 支援 CSV/JSON 匯入、欄位映射、批次驗證
- 驗證結果分級：error / warning / info
- 批次結果可下載檢核報告

### 2.4 MappingPreviewPanel
- 顯示 Preview 統計：成功/失敗/衝突
- 顯示每筆失敗原因與修復建議
- 支援篩選：只看失敗 / 只看未對應 tag

### 2.5 CommitPanel（真實狀態）
- 操作前顯示 dry-run 摘要
- Commit 後以事件時間軸呈現（queued/running/success/fail）
- 失敗時顯示 rollback 指引與重試按鈕

### 2.6 RuntimeHealthPanel
- 指標：資料延遲、最近成功寫入、錯誤率、breaker 狀態
- 事件流：與 SSE 真實事件對齊（非 mock ticker）

## 3. 互動與狀態規範
### 3.1 統一 UI 狀態
- Loading：Skeleton + 進度文案
- Empty：說明 + 下一步 CTA
- Error：錯誤碼、人類可讀原因、修復動作
- Success：摘要 + 下一步
- Partial：成功/失敗混合摘要

### 3.2 非同步行為
- 所有資料操作改由 React Query 管理
- mutation 失敗顯示可重試與錯誤追蹤 ID
- optimistic update 僅限低風險欄位

## 4. RWD 策略
### Desktop（>=1280）
- 三欄式（導航 / 工作區 / 上下文）

### Tablet（768~1279）
- 兩欄式（工作區 + 可收合側欄）

### Mobile（<768）
- 單欄步驟流程，右側面板改底部抽屜
- 高頻按鈕固定底部（sticky action bar）

## 5. 無障礙（A11y）
- 支援鍵盤導覽（Tab 順序、Esc 關閉 modal）
- ARIA 標籤覆蓋所有互動元件
- 對比符合 WCAG AA
- 保留 skip-link 並套用到主路由

## 6. 國際化（i18n）
- Smart Dashboard 所有硬編碼字串改 key-based
- 文案命名：`datalink.{module}.{context}.{key}`
- 必備語系：`zh-TW`、`en`

## 7. 可觀測性規格
- 事件追蹤：
  - `datalink_flow_step_entered`
  - `datalink_validation_failed`
  - `datalink_commit_started/succeeded/failed`
  - `datalink_runtime_alert_seen`
- 每個事件至少包含：`user_id`, `project_id`, `device_id`, `step`, `timestamp`

## 8. 驗收清單（UI Spec 層）
- [ ] 設計稿與元件 props 對齊
- [ ] 所有 API 錯誤在 UI 有明確對應
- [ ] 手機版可完成完整建立流程
- [ ] e2e 腳本覆蓋關鍵路徑
- [ ] i18n key 覆蓋率 100%

---

## 來源
- `frontend/src/pages/datalink/SmartDashboardPage.tsx`
- `frontend/src/components/datalink/*`
- `frontend/src/hooks/datalink/*`
- `internal/api/router.go`
- `internal/datalink/runtime/service.go`
