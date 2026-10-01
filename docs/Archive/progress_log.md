# go_gateway 進度總覽（彙整版）

更新時間：2026-03-03（Asia/Taipei）
資料來源：
- `docs/implement-datalink-ui-test-summary.md`
- `docs/datalink_ticket_backlog.md`
- `docs/merge_postcheck_2026-03-02.md`
- `docs/next_session_execution_plan_2026-03-02.md`
- `docs/test_delivery_checklist.md`

---

## 任務追蹤表（站會版）

| 任務代號 | 狀態 | 分支 | PR | 負責代理 | ETA | 阻塞原因 | 最新驗證 | 最後更新 |
|---|---|---|---|---|---|---|---|---|
| GG-20260302-001（SmartDashboard Batch1~4） | 已完成 | `feat/datalink-smartdashboard-refactor-batch1` | 待補連結 | Codex + Gemini | 已完成 | 無 | 依 merge 後驗收：frontend test/lint/build PASS | 2026-03-02 |
| GG-20260302-002（SmartDashboard Batch5） | 進行中 | `feat/datalink-smartdashboard-batch5-i18n` | 待補連結 | Codex + Gemini | 待補 ETA | 前端測試/建置需再收斂，待環境與案例修正 | checklist 曾出現 frontend test/build fail，需再收斂 | 2026-03-02 |
| GG-20260302-003（Datalink Ticket Backlog P0） | 待辦 | 待建立 | - | 待指定 | 待排程 | 尚未展開任務分派與時程 | 尚未開始（backlog 多數未勾選） | 2026-03-02 |

---

## 已完成（Done）

### A. Datalink 測試覆蓋
- SQL Repository 測試：6 模組，約 120 case，標記完成
- API Handler 測試：7 檔，約 85 case，標記完成
- 合計約 205 測試案例，報告標示覆蓋率約 90~95%
- 來源：`docs/implement-datalink-ui-test-summary.md`

### B. Merge 後驗收（2026-03-02）
- `main` 已完成 fetch/pull 同步
- Frontend 驗證：`npm test -- --run`、`npm run lint`、`npm run build` 全 PASS（該次驗收）
- 來源：`docs/merge_postcheck_2026-03-02.md`

### C. SmartDashboard 重構基線（計畫層敘述）
- 計畫文件記錄「Batch 1~4 已合併」
- 來源：`docs/next_session_execution_plan_2026-03-02.md`

---

## 進行中（In Progress）

### D. SmartDashboard Batch5（i18n + 四態體驗）
目標：
1. 統一 Loading/Empty/Error/Success 四態
2. 優化 Commit/Retry/Rollback/Workflow 回饋
3. 清理硬編碼字串改 i18n key（zh-TW/en）
4. 維持 test/lint/build 全綠

- 來源：`docs/next_session_execution_plan_2026-03-02.md`

---

## 待辦（Todo）

### E. Datalink Ticket Backlog（多數未勾選）
- P0：DLK-001~005（readiness 契約、poll 真實化、SmartDashboard 測試基座、IA 重整、統一狀態元件）
- P1：DLK-006~008（Commit-Runtime 聯動、i18n 全鍵值化、React Query 收斂）
- P2：DLK-009~010（新手引導、KPI 儀表板）
- 來源：`docs/datalink_ticket_backlog.md`

### F. implement-datalink-ui 測試報告中未完成項
- SQLite/Postgres 完整整合測試
- 前端手動驗證（Wizard 流程、SSE、草稿儲存/恢復）
- 來源：`docs/implement-datalink-ui-test-summary.md`

---

## 阻塞/風險（依文件）

1. API 雙軌風險（`internal/api` vs `internal/datalink/api`）
2. Runtime 事件若仍 mock，會影響觀測聯動真實性
3. 部分驗證環境缺少 `go` / `golangci-lint`（在 checklist 紀錄）
4. Batch5 檢核紀錄中曾出現前端 test/build fail（需再收斂）

來源：
- `docs/datalink_ticket_backlog.md`
- `docs/test_delivery_checklist.md`

---

## 建議維護方式
- 每完成一批（Batch / PR）就補 1 筆：
  - 任務代號
  - 分支/PR
  - 測試結果
  - 是否可上線
- 每日只保留「最新狀態 + 連結」，細節留在對應報告檔，避免進度檔過胖。
