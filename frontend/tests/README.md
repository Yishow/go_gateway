# Frontend Tests

前端測試統一由 `frontend/tests/` 作為正式入口，依測試性質分為：

- `tests/unit/`
  Vitest 單元與頁面互動測試入口。
- `tests/integration/`
  Vitest 跨模組整合測試入口。
- `tests/e2e/`
  Playwright 端對端測試。

## 實體遷移狀態

已完成實體遷移至 `frontend/tests/` 的模組：
- `tests/unit/utils/`
- `tests/unit/hooks/`
- `tests/unit/features/`（datalink、gateway、flow）
- `tests/unit/components/`（datalink、monitoring）
- `tests/integration/features/flow/`
- `tests/unit/pages/test-page.test.tsx`

仍以 wrapper 匯入 `src/pages/.../__tests__/` 的頁面測試：
- Gateway（QuickSetupPage、ExpertWorkbenchPage）
- SmartDashboard（interaction、grid-overlays、state hooks）、LocalModbusWorkbenchPage

上述頁面測試後續可視需要做實體搬遷；新增測試一律放在 `frontend/tests/` 對應分類。

## 分類原則

- `tests/unit/components/`
- `tests/unit/features/`
- `tests/unit/hooks/`
- `tests/unit/pages/`
- `tests/unit/utils/`
- `tests/integration/features/`
- `tests/e2e/`
