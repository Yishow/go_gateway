## Problem

SQL Server 真資料庫補測顯示：schema-confirm 已成功並保存相同 proof 後，目的端離線時完全相同的確認重放仍觸發新 probe，回傳 422 而不是持久化成功結果。

## Root Cause

HTTP handler 的 validateCompleted 已識別 proofSaved，但仍呼叫 h.scope，經 SuccessorSchemaScope 進入外部能力 probe；既有 service 的 ConfirmSuccessorSchema 已有完整、唯讀的精確重放驗證，handler卻未使用。

## Proposed Solution

proofSaved 分支改呼叫既有 ConfirmSuccessorSchema，重新核對持久化 ledger、reservation、scope、current revision、connector revision 與 exact proof；驗證成功才重放原結果，不再 probe 或 DDL。其他未保存 proof 的路徑保持原樣。

## Non-Goals (optional)

不放寬驗證、不新增第二套proof authority、不改第一次schema-confirm或activation、不由成功schema proof啟用採集、不更動已封存change。

## Success Criteria

已保存的完全相同確認在目的端離線時仍回原成功結果，內部/外部狀態不變；alias token、connector/current revision不匹配必須拒絕。真SQLServer回歸與本機負向回歸通過。

## Impact

- Affected code:
  - Modified: internal/api/handlers/studio_v2_workspace_write_group_successor_schema.go
  - Modified: internal/api/router_studio_v2_successor_sqlserver_token_test.go
  - New: none
  - Removed: none
- Affected specs: recording-lifecycle-management（新增精確schema-confirm重放要求）
