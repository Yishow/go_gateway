# UI/UX 收斂計畫執行狀態

## 規劃檔案對照

**規劃檔案**：`.cursor/plans/ui-ux-audit-plan_e4905ba3.plan.md`

## Todos 完成狀態

| Todo ID | 內容 | 狀態 | 對應 Phase |
|---------|------|------|-----------|
| `legacy-audit` | 盤查 SmartDashboard 現役結構、舊路由、未引用元件與可裁剪 legacy 相容層 | ✅ complete | Phase 1 |
| `tdd-test-foundation` | 先補主要 UI 流程測試，再進行任何介面調整 | ✅ complete | Phase 1 |
| `design-system-foundation` | 建立最小可行的設計系統規範 | ✅ complete | Phase 2 |
| `smartdashboard-core-flow` | 重整 SmartDashboard 為單人友善主流程 | ✅ complete | Phase 3 |
| `testpage-style-only` | 保留 TestPage 為專用測試工具，只做風格一致化 | ✅ complete | Phase 4 |
| `a11y-responsive-validation` | 補齊可近用性、響應式與主要互動流程驗證基準 | ✅ complete | Phase 5 |

## Phase 完成狀態

| Phase | 狀態 | 完成日期 | 備註 |
|-------|------|---------|------|
| Phase 1: Legacy 盤查與 TDD 基線建立 | ✅ complete | 2026-03-08 | 25 tests passed |
| Phase 2: 建立最小 UI 規範底座 | ✅ complete | 2026-03-08 | `designSystem.ts` 已建立 |
| Phase 3: 重整 SmartDashboard 核心流程 | ✅ complete | 2026-03-08 | 流程指引、資料流向說明已加入 |
| Phase 4: 整理 LocalModbusWorkbench 與 TestPage | ✅ complete | 2026-03-08 | 區塊層級已重整 |
| Phase 5: 可近用性、響應式與品質驗證 | ✅ complete | 2026-03-08 | Accessibility 與響應式檢查完成 |

## 後續工作（Phase 2 後續實作）

以下項目在 Phase 2 中標記為「後續實作」，尚未開始：

- [ ] 逐步將現有元件遷移到使用 `designSystem.components.*`
- [ ] 補齊表單元件的 `autocomplete`、`inputmode`、`name` 屬性
- [ ] 統一 microcopy 使用 `designSystem.microcopy.*`
- [ ] 考慮建立統一的 FormInput、FormLabel、FormError 元件

**狀態**：這些是漸進式改進項目，不影響核心功能，可視需要逐步進行。

## 測試驗證結果

- ✅ SmartDashboard 單元測試：25 tests passed（4 test files）
- ✅ LocalModbusWorkbench 測試：4 tests passed
- ✅ UI Regression 測試：3 tests passed
- ✅ Lint 檢查：通過
- ✅ Build 檢查：通過

## 修改檔案清單

### Phase 1-2
- `frontend/src/styles/designSystem.ts`（新建）
- `frontend/tests/README.md`（新建）
- `frontend/tests/unit/**/*.test.ts(x)`（遷移與重組）

### Phase 3
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardCommitPanel.tsx`
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardSidebar.tsx`
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardWorkflowModal.tsx`
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardWorkspaceContent.tsx`

### Phase 4
- `frontend/src/pages/datalink/LocalModbusWorkbenchPage.tsx`
- `frontend/src/pages/datalink/__tests__/LocalModbusWorkbenchPage.test.tsx`

### Phase 5
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardSidebar.tsx`（accessibility）
- `frontend/src/pages/datalink/smart-dashboard/SmartDashboardCommitPanel.tsx`（accessibility）
- `frontend/src/pages/datalink/LocalModbusWorkbenchPage.tsx`（accessibility）
- `frontend/tests/integration/ui/smart-dashboard-regression.test.tsx`（新建）

## 工作記錄檔案

- `task_plan.md`：階段、狀態、決策、待辦與錯誤
- `findings.md`：盤查結果、研究發現、legacy 清單、UI/UX 問題
- `progress.md`：工作階段、測試結果、已修改檔案與驗證紀錄
