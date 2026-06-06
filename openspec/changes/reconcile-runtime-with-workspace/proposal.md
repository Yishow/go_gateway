## Why

目前 activation 後的 runtime 與 persisted workspace 之間缺少正式的 reconciliation 契約，導致 device、source rule、mapping、database target 在保存後可能已經改變，但 runtime 還在用舊投影、舊 scheduler 狀態、舊 mapping cache。這會直接造成「Runtime 顯示不含工作台設定」與「採集還在跑，但跑的不是現在這份配置」。

## What Changes

- 定義 runtime workspace projection 與 reconciliation 契約，讓 activation、restart、live apply 都使用同一份 persisted workspace projection。
- 規範 persisted config 變更後哪些可以 targeted reconcile，哪些必須顯式標記 pending restart 或 pending reactivation。
- 讓 runtime snapshot/stream 與 workspace runtime view 能看出目前是最新投影、延遲投影、或失敗投影。

## Capabilities

### New Capabilities

- `runtime-workspace-reconciliation`: 定義 persisted workspace projection、reconciliation triggers、targeted refresh、pending restart/deferred apply 與 reconciliation outcome。

### Modified Capabilities

- `studio-v2-live-config-apply`: live apply 回應必須說明 reconcile outcome，而不是只說成功保存。
- `studio-runtime-workspace-view`: workspace-scoped runtime view 必須反映當前 runtime projection 是否與 workspace 對齊。
- `source-rule-runtime`: rule lifecycle 變更必須觸發與 runtime projection 一致的 reconcile 結果。
- `runtime-dashboard-backend`: runtime snapshot/stream contract 必須攜帶 projection/reconciliation 狀態，避免 UI 看不到 drift。

## Impact

- Affected specs:
  - New: runtime-workspace-reconciliation
  - Modified: studio-v2-live-config-apply, studio-runtime-workspace-view, source-rule-runtime, runtime-dashboard-backend
- Affected code:
  - Modified: internal/datalink/runtime/service.go
  - Modified: internal/datalink/runtime/service_device_sync.go
  - Modified: internal/datalink/workspace/service_activation.go
  - Modified: internal/api/handlers/studio_v2_runtime_apply.go
  - Modified: internal/api/handlers/runtime_handler.go
  - Modified: internal/api/handlers/runtime_stream_handler.go
  - Modified: frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - Modified: frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/CommitSuccessCard.tsx
  - Modified: frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
