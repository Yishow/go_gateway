# 獨立程式 review：shorten-recording-save-feedback

## 結論與範圍

完成 correctness、efficiency、simplification/reuse、convention 四面檢查；目前沒有未解決的高信心缺陷。這是程式範圍審查，不是 browser／field／release acceptance。

- Base與HEAD：`406490c8ff1f801c27ddc491d9f29debbd12d523`，加上當前 unstaged／untracked layers。
- Scope source：`touched_tracking`；status：`resolved`。
- Source snapshot：`998030ac80ea03b95498bc3955f168aab1b3ba57`；獨立 reviewer 的 check-snapshot 回傳 current。
- 審查51個可讀文字檔，沒有 binary 或未能檢查的 source。
- Spectra 標記 `database-save-feedback.test.tsx` 在後續 task baseline 前已 dirty；路徑歸屬不能單獨證明每個 hunk，完整新增檔已逐項審查。
- 排除原有未提交 `DESIGN.md`。Reviewer 沒有編輯 repository、stage、commit、archive 或 push。

## 已修復的發現

### Warning：唯讀恢復後 refresh failure 可能卡住

原 anchor：`recordingConfigRefresh.ts:98-108`；`Step4SupportPanels.tsx:61`。違反 delta spec「Background refresh remains scoped and observable」的可用 retry 契約。

有效保存後 background config GET 失敗，再以頁面既有 refetch 完成三類讀取，或第一次 lost ACK 後於相同 QueryClient 重開 editor，均可讓全部 queries success/idle，但 refresh unavailable 仍為 true。獨立重現2個案例先失敗，修復後2/2通過。

修復只在 stale-error、沒有 pending coordinator refresh，且三類 matching active reads 都 success／idle／未 invalidated 才恢復。沒有製造 ACK，也不清除 unknown-write guard。另有回歸防止 PUT 尚 pending 時，較早完成的 GET 解除鎖定；單獨 bootstrap 成功不足。

### Warning：owner A→B→A 可接受舊 ACK

Anchor：`useStudioV2DatabaseAutosave.ts:107-112,235-240,265,285`。違反「New edit or workspace wins」與 scope fencing 契約。

只比較返回後的 workspace/origin 值會讓 round trip 恢復相同值，誤接受原 owner 的 ACK。兩個新增回歸先失敗，之後以捕捉 owner object identity 修復；workspace、完整 baseURL 或 enabled 狀態改變即替換 token。成功與失敗回覆都檢查 token，保留既有 reducer field-generation fences。

## 其他核對

- config-only ACK 與並行刷新分離，捕捉事前 owner、取消與 generation 隔離，沒有第二份 persisted config/home cache。
- B 的 capture／complete mutation API被沿用；active current bootstrap owner優先，失敗不使用舊 cached owner。
- target／write-group／schema／activation／test-write helpers保留 awaited reconciliation/readback。
- transport只新增 captured baseURL／AbortSignal options，沒有改後端或 HTTP schema。
- locator測試僅等待非同步 enabled狀態；禁止 save/source/lifecycle、stale owner與 no-extra-read斷言全部保留。
- synthetic1ms＋三個200ms是 hook/query測試，不是實際使用延遲或 SQL證據。

完整 frontend QA 的最終結果見 `non-browser-evidence.md`。Mac browser、截圖、SQL／設備／field、archive與部署未執行。

## 審查檔案

- `frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx`
- `frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx`
- `frontend/src/hooks/datalink/recordingConfigRefresh.ts`
- `frontend/src/hooks/datalink/recordingHomeSourceCache.ts`
- `frontend/src/hooks/datalink/recordingHomeSourceMutation.ts`
- `frontend/src/hooks/datalink/useStudioV2Workspace.ts`
- `frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts`
- `frontend/src/i18n/locales/en/workbench-v2.json`
- `frontend/src/i18n/locales/zh-TW/workbench-v2.json`
- `frontend/src/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page.tsx`
- `frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts`
- `frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts`
- `frontend/src/services/studioV2Workspace.ts`
- `frontend/src/services/studioV2WorkspaceDatabase.ts`
- `frontend/tests/e2e/fixtures/database-save-feedback-api.ts`
- `frontend/tests/e2e/recording-database-save-feedback.spec.ts`
- `frontend/tests/unit/hooks/use-studio-locator-entry.test.tsx`
- `frontend/tests/unit/hooks/useStudioV2Mappings.test.ts`
- `frontend/tests/unit/hooks/useStudioV2Rules.test.ts`
- `frontend/tests/unit/hooks/useStudioV2Workspace.test.tsx`
- `frontend/tests/unit/hooks/useStudioV2WorkspaceDatabase.test.ts`
- `frontend/tests/unit/hooks/useStudioV2WorkspaceDevices.test.ts`
- `frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx`
- `frontend/tests/unit/workbench-v2/database-autosave-page.testHarness.tsx`
- `frontend/tests/unit/workbench-v2/database-autosave-page.validation.test.tsx`
- `frontend/tests/unit/workbench-v2/database-connector-input-race.test.tsx`
- `frontend/tests/unit/workbench-v2/database-save-feedback-api.test.ts`
- `frontend/tests/unit/workbench-v2/database-save-feedback-barriers.test.tsx`
- `frontend/tests/unit/workbench-v2/database-save-feedback-owner.test.tsx`
- `frontend/tests/unit/workbench-v2/database-save-feedback-page.test.tsx`
- `frontend/tests/unit/workbench-v2/database-save-feedback-presentation.test.tsx`
- `frontend/tests/unit/workbench-v2/database-save-feedback.test.tsx`
- `frontend/tests/unit/workbench-v2/database-save-refresh-scope.test.tsx`
- `frontend/tests/unit/workbench-v2/locator-own-save-source-context.test.tsx`
- `frontend/tests/unit/workbench-v2/mapping-autosave-page.hydration.test.tsx`
- `frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx`
- `frontend/tests/unit/workbench-v2/mc-persisted-hydration.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-draft-page-refresh.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-editor-range-mapping-recovery.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-guided-host.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-home-hooks.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-home-mc-locator.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-home-route.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-home-source-modal.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-test-write-recovery-host.test.tsx`
- `frontend/tests/unit/workbench-v2/share-status-page-stability.test.tsx`
- `frontend/tests/unit/workbench-v2/step4-first-activation.test.tsx`
- `frontend/tests/unit/workbench-v2/step4-share-activation.test.tsx`
- `frontend/tests/unit/workbench-v2/studio-locator-group-focus.test.tsx`
- `frontend/tests/unit/workbench-v2/studio-locator-receiver.test.tsx`
- `frontend/tests/unit/workbench-v2/workspace-boot.test.tsx`
