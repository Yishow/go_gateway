## 1. Workbench Foundation split

- [x] [P] 1.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 legacy-routes、operator-flow、device-list、device-editor、device-inspector 行為分別移至 DatalinkWorkbenchFoundation.legacy-routes.test.tsx、DatalinkWorkbenchFoundation.operator-flow.test.tsx、DatalinkWorkbenchFoundation.device-list.test.tsx、DatalinkWorkbenchFoundation.device-editor.test.tsx、DatalinkWorkbenchFoundation.device-inspector.test.tsx，將共用 typed setup 移至 DatalinkWorkbenchFoundation.testHarness.tsx，保留每個 title/assertion/mock/fixture/route，確認 6 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 2. Workbench Output split

- [x] [P] 2.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 modbus-surface、modbus-bindings、modbus-canvas、database-targets、database-schema、database-schema-columns、output-inspector、cross-step-handoff 行為分別移至 DatalinkWorkbenchOutputStep.modbus-surface.test.tsx、DatalinkWorkbenchOutputStep.modbus-bindings.test.tsx、DatalinkWorkbenchOutputStep.modbus-canvas.test.tsx、DatalinkWorkbenchOutputStep.database-targets.test.tsx、DatalinkWorkbenchOutputStep.database-schema.test.tsx、DatalinkWorkbenchOutputStep.database-schema-columns.test.tsx、DatalinkWorkbenchOutputStep.output-inspector.test.tsx、DatalinkWorkbenchOutputStep.cross-step-handoff.test.tsx，將共用 typed setup 移至 DatalinkWorkbenchOutputStep.testHarness.tsx，確認 9 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 3. Workbench Shell split

- [x] [P] 3.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 frame、context、inspector、summary、switching 行為分別移至 DatalinkWorkbenchShellUi.frame.test.tsx、DatalinkWorkbenchShellUi.context.test.tsx、DatalinkWorkbenchShellUi.inspector.test.tsx、DatalinkWorkbenchShellUi.summary.test.tsx、DatalinkWorkbenchShellUi.switching.test.tsx，將共用 typed setup 移至 DatalinkWorkbenchShellUi.testHarness.tsx，確認 6 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 4. Workbench Source split

- [x] [P] 4.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 planning、inspector、rule-actions、canvas-overlays、batch-selection、inline-editing、conflict-queue、datatype-span、visual-span 行為分別移至 DatalinkWorkbenchSourceStep.planning.test.tsx、DatalinkWorkbenchSourceStep.inspector.test.tsx、DatalinkWorkbenchSourceStep.rule-actions.test.tsx、DatalinkWorkbenchSourceStep.canvas-overlays.test.tsx、DatalinkWorkbenchSourceStep.batch-selection.test.tsx、DatalinkWorkbenchSourceStep.inline-editing.test.tsx、DatalinkWorkbenchSourceStep.conflict-queue.test.tsx、DatalinkWorkbenchSourceStep.datatype-span.test.tsx、DatalinkWorkbenchSourceStep.visual-span.test.tsx，將共用 typed setup 移至 DatalinkWorkbenchSourceStep.testHarness.tsx，確認 10 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 5. Workbench Tag split

- [x] [P] 5.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 review-overview、review-selection、tag-binding、tag-results 行為分別移至 DatalinkWorkbenchTagStep.review-overview.test.tsx、DatalinkWorkbenchTagStep.review-selection.test.tsx、DatalinkWorkbenchTagStep.tag-binding.test.tsx、DatalinkWorkbenchTagStep.tag-results.test.tsx，將共用 typed setup 移至 DatalinkWorkbenchTagStep.testHarness.tsx，確認 5 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 6. Output Incident Desk split

- [x] [P] 6.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 command 與 recovery 行為分別移至 MuiOutputIncidentDesk.reopen.command.test.tsx、MuiOutputIncidentDesk.reopen.recovery.test.tsx，將共用 typed bootstrap 移至 MuiOutputIncidentDesk.reopen.testHarness.tsx，確認 3 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 7. Database autosave split

- [x] [P] 7.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 hydration、validation、row-groups、bulk-toggle 行為分別移至 database-autosave-page.hydration.test.tsx、database-autosave-page.validation.test.tsx、database-autosave-page.row-groups.test.tsx、database-autosave-page.bulk-toggle.test.tsx，將共用 typed setup 移至 database-autosave-page.testHarness.tsx，確認 5 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 8. Device autosave split

- [x] [P] 8.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 hydration、draft-validation、failure-recovery 行為分別移至 device-autosave-page.hydration.test.tsx、device-autosave-page.draft-validation.test.tsx、device-autosave-page.failure-recovery.test.tsx，將共用 typed setup 移至 device-autosave-page.testHarness.tsx，確認 4 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 9. Mapping autosave split

- [x] [P] 9.1 依「以行為邊界拆分 monolith」與「test-only typed harness 與行為檔分離」先對 frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx 執行 characterization focused test，記錄 line count、完整 test title set 與既有失敗；再將 hydration、validation、reconciliation、bulk-toggle 行為分別移至 mapping-autosave-page.hydration.test.tsx、mapping-autosave-page.validation.test.tsx、mapping-autosave-page.reconciliation.test.tsx、mapping-autosave-page.bulk-toggle.test.tsx，將共用 typed setup 移至 mapping-autosave-page.testHarness.tsx，確認 5 個 destination 各自不超過 300 行、focused destination tests 通過且 title set 相等後才移除原檔。

## 10. Discovery and contract parity

- [x] 10.1 完成「保留測試探索與契約」及 Implementation Contract 的 Observable behavior、Exact split ownership、File and discovery contract 驗證：收集 53 個 destination glob，確認所有 .testHarness.tsx 沒有 top-level test registration、每個原始 suite 的 title set 與 destination title set 完全相等且每個 title 只出現一次、9 個原始 monolith 均已移除；以 frontend 的 Vitest discovery 與各群組 focused command 留存可重跑輸出。

## 11. Line-limit gate

- [x] 11.1 完成「Line-limit and frontend quality gates remain enforceable」與 Implementation Contract 的 Acceptance criteria：對 frontend/tests/unit 下所有變更的 TS/TSX（含 harness）執行物理行數檢查，確認全部不超過 300 行、未修改 .line-limit-ignore，並執行 make check-lines 或 repo 支援的等價 line gate；若任何 destination 超限，依 Failure modes and boundaries 在刪除原檔前再按行為切分。

## 12. Frontend quality and timeout evidence

- [x] 12.1 依「characterization-first 與分層驗證」及「Full-suite timeout is measured separately from the split」完成前端品質驗證：執行 cd frontend && npm run lint、cd frontend && npm run test、cd frontend && npm run build，記錄 pre-split/post-split full-suite duration、Workbench timeout 或既有失敗的 exact command/output，將 focused split 結果與 baseline 分開；確認沒有提高 timeout、skip test、修改 discovery/config 或引入 unrelated workaround。

## 13. Final audit

- [x] 13.1 依「Acceptance criteria」完成 final diff/audit review：確認新 capability frontend-test-file-maintainability 的四項 requirements（Oversized test suites are split into discoverable behavior files、Shared test-only harness preserves the test contract、Line-limit and frontend quality gates remain enforceable、Full-suite timeout is measured separately from the split）均有對應證據，proposal/design/spec/tasks 路徑完整，所有變更僅在核准的 frontend/tests/unit 與 OpenSpec artifact 範圍，無 production、dependency、timeout、ignore-list 或其他 test/config 變更，並留下已知 full-suite timeout 後續需另案處理的報告；另逐項核對 design.md 的 decision: 以行為邊界拆分 monolith、decision: test-only typed harness 與行為檔分離、decision: 保留測試探索與契約、decision: characterization-first 與分層驗證。
