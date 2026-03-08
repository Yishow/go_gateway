# UI/UX 收斂任務計畫

## 任務目標
- 將專案規範來源整理為一致的層級結構。
- 明確要求 Agent 同時遵守 `AGENTS.md`、agent 專屬文件與 `.github/instructions/`。
- 以單人操作情境重整前端 UI/UX 改造方向。
- 所有後續 UI 變更採 TDD 先行，並保留完整工作記錄。

## 規範來源層級
- `AGENTS.md`
  專案共通目標、目錄分工、測試、TDD、文件化工作流。
- `CLAUDE.md` / `GEMINI.md`
  Agent 工作流程、架構脈絡、執行邊界與補充說明。
- `.github/instructions/*.md`
  依檔案類型套用的語言與框架實作規範。

## 階段
- [complete] Phase 1: 更新專案規範文件
- [complete] Phase 2: 建立工作記錄檔
- [in_progress] Phase 3: 後續 UI/UX 實作前置準備

## 待辦
- [x] 在 `AGENTS.md` 宣告 `.github/instructions/` 為正式規範來源
- [x] 在 `CLAUDE.md` 補充文件分工與閱讀順序
- [x] 在 `GEMINI.md` 對齊相同規範來源層級
- [x] 建立 `findings.md`
- [x] 建立 `progress.md`
- [ ] 規劃 SmartDashboard / LocalModbusWorkbench / TestPage 的 TDD 基線
- [x] 補上 `TestPage` 第一批基線測試
- [x] 完成 `TestPage` 第一輪風格一致化收斂
- [x] 盤點前端測試分布，將正式測試入口統一到 `frontend/tests/` 並分類
- [x] 更新 `AGENTS.md`，明確規定前端測試入口與 Go 測試例外
- [ ] 盤點 `.github/instructions/` 在後續 UI/UX 實作中的具體套用點

## 錯誤紀錄
- `tests/unit/pages/datalink.test.ts`
  過度聚合 `SmartDashboard` 相關測試後，`vi.mock` 隔離被模組快取污染，導致 `SmartDashboardGridOverlaysSection` 的 mock 未生效；改為更細粒度的 root wrapper 後解決。
- `test:gateway:unit`
  `GatewayQuickSetupPage` 舊測試在 root wrapper 驗證時暴露出既有不穩定斷言；改成重新查詢啟用後的提交按鈕並以獨立 wrapper 執行 page suites 後恢復穩定。
