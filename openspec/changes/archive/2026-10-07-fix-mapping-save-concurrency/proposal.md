## Problem

同一來源規則的八列 mapping 同時儲存時，現場 PUT 回傳 409、404、500，讀值與 preview 卻正常。隔離真 SQL 重現依序八列皆成功、並行失敗；錯誤不能被泛稱為使用者映射內容不合法。

## Root Cause

前端按 point 限制重入，批次會並行；workspace save 與 confirmation 同步其他 mapping 並以分開 DELETE/INSERT 重建整組 source_rule_links，缺乏共同規則協調與原子替換。舊連結替換競爭與新確認 CAS 路徑共同受影響。

## Proposed Solution

同規則儲存、source mutation、derived sync、candidate apply 共用可重入且有序的協調；SQL/Memory 原子替換連結；workspace save→confirm 位於協調內，保留真 stale revision/accepted hash/CAS 拒絕。前端按 rule 排隊、保留最新草稿及可明確重試的失敗；顯示安全 status/code 與可行動訊息。
Modified Capabilities: datalink-workbench-v2-step3-mapping。

## Success Criteria

真 SQL 八列並行全數成功且身份/pipeline/accepted hash 不丟失；DELETE 窗口消失；失敗替換回滾；不同 rule 可進展；save 與 edit/sync/reapply/delete 無死鎖且 stale source 仍拒絕。UI 多次批次、快速連改、失敗重試與正常 preview/persist readback 通過。完整 repo gates 與 exact source review 完成。

## Impact

- Affected code:
  - Modified: internal/datalink/sourcerule/, internal/datalink/mapping/sql_repo.go, internal/datalink/mapping/workspace_save_rollback.go, internal/datalink/tag/workspace_save_rollback.go, internal/api/handlers/studio_v2_workspace_mapping_request.go, internal/api/handlers/studio_v2_workspace_mappings_handler.go, internal/api/handlers/studio_v2_workspace_mappings_recovery.go, internal/api/handlers/studio_v2_workspace_mappings_delete.go, frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts, frontend/src/features/datalink/workbench-v2/steps/step3/MappingRow.tsx, frontend/src/utils/backendErrorCodes.ts, frontend state/hooks/services safe error metadata, frontend/src/i18n/config.ts（新增 mapping-errors namespace；不修改歷史超長 workbench-v2.json）
  - New: scoped Go/React regression tests, rule-scoped coordinator/queue helpers, safe MappingSaveFailure, en/zh-TW mapping-errors.json, and isolated mapping UI regression runner.
  - Removed: none.
