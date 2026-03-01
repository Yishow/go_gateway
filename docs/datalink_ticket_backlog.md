# Datalink UI 改善 Ticket Backlog

> 優先級定義：P0（必做/阻塞上線）、P1（重要/提升品質）、P2（加值）

## P0

### DLK-001｜對齊 Device Readiness API（前後端契約）
- 類型：Backend + Frontend
- 說明：補齊主線 API readiness endpoint，前端 service 改用一致契約
- 交付：
  - Gin router/handler 支援 readiness
  - 前端 `deviceAPI.checkReadiness` 整合新回應格式
- 驗收條件：
  - [ ] `/api/v1/datalink/devices/:id/readiness` 可用
  - [ ] UI 顯示 ready/not_ready/error 與原因
  - [ ] 單元測試 + 整合測試通過
- 估時：1.5 人日

### DLK-002｜Poll / PollBatch 真實化
- 類型：Backend
- 說明：移除 placeholder，串接 scheduler poll now/batch 行為
- 驗收條件：
  - [ ] Poll 端點回傳真實資料
  - [ ] 錯誤碼與訊息標準化
  - [ ] 壓測下錯誤率符合既有門檻
- 估時：2 人日

### DLK-003｜SmartDashboard 測試基座修復
- 類型：Frontend Test
- 說明：interaction test 補 `QueryClientProvider` 與必要 mock
- 驗收條件：
  - [ ] `SmartDashboard.interaction.test.tsx` 全綠
  - [ ] CI 可重現穩定通過
- 估時：0.5 人日

### DLK-004｜主流程 IA 重整（Overview→Commit）
- 類型：Design + Frontend
- 說明：導覽與區塊改任務導向，減少跨頁來回
- 驗收條件：
  - [ ] 一條主路徑可完成建立→驗證→提交
  - [ ] 每步驟都有阻塞原因與下一步引導
- 估時：3 人日

### DLK-005｜統一錯誤/載入/空狀態元件
- 類型：Frontend
- 說明：建立全域狀態元件與使用規範
- 驗收條件：
  - [ ] Device/Point/Mapping/Commit 全覆蓋
  - [ ] 每種狀態有一致文案與 CTA
- 估時：1.5 人日

## P1

### DLK-006｜Commit 與 Runtime 觀測聯動
- 類型：Backend + Frontend
- 說明：Commit panel 顯示真實 runtime event timeline
- 驗收條件：
  - [ ] Commit 事件與 runtime 指標可對照
  - [ ] SSE 來源改為真實事件流
- 估時：3 人日

### DLK-007｜i18n 全面鍵值化
- 類型：Frontend
- 說明：移除 Smart Dashboard 硬編碼字串
- 驗收條件：
  - [ ] zh-TW / en 可完整切換
  - [ ] 新增字串皆走 i18n key
- 估時：1 人日

### DLK-008｜Batch/Point 元件改 React Query
- 類型：Frontend
- 說明：移除直接 fetch，統一快取、錯誤、重試行為
- 驗收條件：
  - [ ] `BatchPointCreator`、`PointDetailPanel` 不再直接 fetch
  - [ ] mutation 錯誤有追蹤 ID
- 估時：1.5 人日

## P2

### DLK-009｜新手引導模式（Guided Setup）
- 類型：Design + Frontend
- 驗收條件：
  - [ ] 首次使用者可按步驟完成設定
  - [ ] 可關閉/重啟引導
- 估時：2 人日

### DLK-010｜操作品質儀表板（KPI）
- 類型：Frontend + Data
- 驗收條件：
  - [ ] 可查看流程完成率、失敗率、平均耗時
  - [ ] 可依版本比對趨勢
- 估時：2 人日

---

## 依賴與風險
- API 雙軌（internal/api vs internal/datalink/api）需先決策唯一主線
- runtime 事件來源若仍 mock，P1 觀測聯動會失真
- 測試環境需固定 Node/npm 版本，避免 CI 漂移

## 建議排程（兩週衝刺）
- Sprint 1：DLK-001~005（P0）
- Sprint 2：DLK-006~008（P1）
- Sprint 3：DLK-009~010（P2）

## 來源
- `docs/datalink_flow_gate_v1_plan.md`
- `docs/datalink_runtime_wbs.md`
- `docs/implement-datalink-ui-test-summary.md`
- `frontend/src/pages/datalink/SmartDashboardPage.tsx`
- `internal/api/handlers/*.go`
