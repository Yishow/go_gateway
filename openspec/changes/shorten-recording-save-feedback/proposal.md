## Summary

把已確認設定儲存的 ACK 與後續背景刷新分開，縮短等待但不提前宣稱 readback、啟用或寫入成功。

## Motivation

useUpdateStudioV2DatabaseConfigMutation 在 onSuccess 依序等待三類 invalidation。雲端合成測試設定 ACK 約 1 ms、每類刷新 200 ms，完整等待約 604 ms；這只證明程式等待鏈，不是實際 DB 或使用者瀏覽器儲存延遲。

## Proposed Solution

- 只先改 database configuration save：有效 scope/revision ACK 完成該次保存，背景刷新另顯示 refreshing／stale／error。
- 獨立刷新並行、同 key 重複需求合併；保持 newer local intent、CAS、workspace 切換與 unknown write recovery。
- 保留設定保存、fresh readback、activation、schema apply、lifecycle 與 SQL receipt 的不同完成條件；不對共用 write-group helper 全面取消等待。
- 納入首頁摘要的精確 invalidation，避免快速保存後返回看到被標成 fresh 的舊資料。

## Impact

- Affected specs: studio-v2-workspace
- Affected code:
  - Modified: frontend/src/hooks/datalink/useStudioV2WorkspaceDatabase.ts, frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts, frontend/src/hooks/datalink/keys.ts, frontend/src/hooks/datalink/useStudioV2Workspace.ts, frontend/src/hooks/datalink/recordingHomeSourceMutation.ts, frontend/src/hooks/datalink/recordingHomeSourceCache.ts, frontend/src/services/studioV2Workspace.ts, frontend/src/services/studioV2WorkspaceDatabase.ts, frontend/src/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page.tsx, frontend/src/features/datalink/workbench-v2/steps/step4/Step4SupportPanels.tsx, frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx, frontend/src/features/datalink/workbench-v2/steps/step4/ConnectorSection.tsx, frontend/src/i18n/locales/en/workbench-v2.json, frontend/src/i18n/locales/zh-TW/workbench-v2.json
  - New: frontend/tests/unit/workbench-v2/database-save-feedback.test.tsx, frontend/src/hooks/datalink/recordingConfigRefresh.ts
  - Removed: 無。
- 依賴 reuse-recording-home-snapshots 的 query key／invalidation 契約後再整合；舊 autosave、reload recovery 與 schema/readback contracts 不改語意。必要 caller 測試與 en/zh-TW 保存狀態文案一併覆蓋。
