## Why

目前 Step 4 database 設定、schema ensure、runtime dbtarget writer、寫入結果回報之間沒有形成一條可驗證的 delivery contract。結果就是你會看到「有在採集，但完全沒寫入到資料庫」，而且操作員很難知道是 connector 不可用、schema 沒建立、mapping 沒生效、flush 失敗，還是 hidden stale target 阻擋整體流程。

## What Changes

- 定義 database output delivery 契約，從 Step 4 connector/targets、schema ensure、runtime writer 到 last-write outcome 一條到底。
- 規範 schema ensure 與 writer 只處理 live、enabled、可解析的 mappings，不讓 hidden stale targets 阻斷整體流程。
- 補 last schema ensure、last write success、last write failure、flush outcome 等 delivery truth，讓 Step 4 與 runtime 可看見。

## Capabilities

### New Capabilities

- `database-output-delivery`: 定義 database connector readiness、schema ensure scope、writer delivery outcome、last-write status 與 failure visibility。

### Modified Capabilities

- `datalink-workbench-v2-step4-database`: Step 4 必須顯示 connector readiness、schema outcome、delivery truth。
- `database-target-workbench`: database target workflow 必須以 live visible target set 為準，不能被 stale hidden mappings 影響。
- `source-rule-runtime`: rule-derived database target 關係變化不得用 hidden stale mappings 阻斷無關 activation 或 write path。

## Impact

- Affected specs:
  - New: database-output-delivery
  - Modified: datalink-workbench-v2-step4-database, database-target-workbench, source-rule-runtime
- Affected code:
  - Modified: internal/datalink/dbtarget/service.go
  - Modified: internal/datalink/dbtarget/tooling_service.go
  - Modified: internal/datalink/dbtarget/writer.go
  - Modified: internal/api/handlers/studio_v2_workspace_database_handler.go
  - Modified: internal/api/handlers/studio_v2_workspace_activation_handler.go
  - Modified: internal/datalink/sourcerule/service.go
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - Modified: frontend/src/services/studioV2WorkspaceDatabase.ts
  - Modified: frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
