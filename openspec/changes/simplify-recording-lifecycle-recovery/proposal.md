## Summary

在簡易記錄首頁／編輯頁完成已保存記錄的修改、複製草稿、停止、重新檢查後恢復、移除與資料表前向變更，以及有證據的結果檢視與異常修復入口。

## Motivation

現行 SavedRecordingActions 同時列出多個技術操作，RecordingChangeDialog 已有安全 state machine，但使用者仍不易理解何時能改、舊資料會怎樣、關閉視窗後操作是否已送出。應簡化呈現而保留所有 scope、revision、schema、ownership 與 operation 保護。

## Proposed Solution

- 新增 simple-recording-lifecycle 能力：依 server action availability 顯示主要操作與有理由的更多選單。
- 複製僅預填新未啟用草稿，不複製 readiness、operation、receipt、ownership proof；managed profile 若既有 ensure 去重不得假稱已建立另一組。
- 寬／長表、來源或目的地的受保護變更重用 successor review/reserve/schema/confirm/start-or-cutover，說明新表、舊歷史與切換邊界；不原地改歷史。
- 提供部分來源、過期、DB失敗及unknown的具名修復導航／查狀態；runtime/logs細節按需展開，不創造通用SQL查詢能力。

## Impact

- Affected specs: simple-recording-lifecycle (new)
- Affected code:
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingActions.tsx; frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeDialog.tsx; frontend/src/features/datalink/workbench-v2/state/recordingWorkflowView.ts
  - New: frontend/src/features/datalink/workbench-v2/recording-home/RecordingActions.tsx; frontend/tests/unit/workbench-v2/simple-recording-lifecycle.test.tsx
  - Removed: none
- 依賴 streamline-recording-task-editor；重用已完成 manage-recording-lifecycle-safely、fix-recording-cutover-boundaries、expose-runtime-recording-rules、make-runtime-logs-actionable。既有後端契約不重寫。

本次更新納入 `DESIGN.md` 語意 tokens、system/light/dark 主題、RWD、簡潔文案與字型階層；新彈出設定為 viewport 寬高各75%的置中 modal。已選手機近滿版／小確認框較小例外，分類依 DESIGN.md；舊 demo 不構成此次確認。
