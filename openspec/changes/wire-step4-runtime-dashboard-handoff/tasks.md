## 1. Handoff 契約測試先行

- [ ] [P] 1.1 以 Red-Green 方式鎖定 `Commit completion card` 與 `Decision: Keep navigation orchestration in WorkbenchV2Shell`：先新增 shell routing 測試，完成後可觀察到成功卡按鈕在 shell integration 中不再只是 `console.log`，而是導向 runtime dashboard route；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/workbench-v2/shell.test.tsx`。
- [ ] [P] 1.2 以 Red-Green 方式鎖定 `Decision: Resolve the handoff device from existing workbench state with a deterministic order`：先新增 `resolveRuntimeDashboardDevice` 單元測試，完成後可觀察到 selected rule、single-device、unresolved 三種情況都回傳固定結果；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts`。
- [ ] [P] 1.3 以 Red-Green 方式鎖定 `Decision: Prefer explicit runtime route fallback over silent no-op`：先補 Step 4 commit flow 測試，完成後可觀察到無法解析 device 時仍會 handoff 到 `/studio/runtime`，而不是停在原頁或只留下 console side effect；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/workbench-v2/step4-commit.test.tsx`。

## 2. Device 解析與導航實作

- [ ] 2.1 落實 `Decision: Resolve the handoff device from existing workbench state with a deterministic order`：新增集中式 helper，依 `selectedRuleId`、single-rule-device、single-device、unresolved 的固定順序解析 handoff device，完成後 shell 不再散落手寫判斷；驗證：重新執行 `npm run test -- --run frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts`。
- [ ] 2.2 落實 `Decision: Keep navigation orchestration in WorkbenchV2Shell` 與 `Commit completion card`：在 `WorkbenchV2Shell` 以 route-aware `onCommit` 取代 placeholder callback，完成後成功卡 CTA 會導向 `/studio/runtime?device_id=<id>` 或 fallback route；驗證：重新執行 `npm run test -- --run frontend/tests/unit/workbench-v2/shell.test.tsx frontend/tests/unit/workbench-v2/step4-commit.test.tsx`。

## 3. Fallback 與回歸驗證

- [ ] 3.1 落實 `Decision: Prefer explicit runtime route fallback over silent no-op`：當 helper 無法解析單一 device 時，成功卡按鈕仍必須導向 `/studio/runtime`，完成後使用者至少會看到 route change 定義的 `missing-device-context`，而不是無反應；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/workbench-v2/step4-commit.test.tsx` 並人工確認不再出現 `console.log`-only handoff。
- [ ] 3.2 完成 `datalink-workbench-v2-step4-database` spec 的實作收斂：執行 `spectra analyze wire-step4-runtime-dashboard-handoff --json`、`spectra validate wire-step4-runtime-dashboard-handoff`、`cd frontend && npm run test -- --run frontend/tests/unit/workbench-v2/resolveRuntimeDashboardDevice.test.ts frontend/tests/unit/workbench-v2/shell.test.tsx frontend/tests/unit/workbench-v2/step4-commit.test.tsx`，確認無 Critical/Warning 且 handoff 行為可直接依賴。
