## Why

目前 Studio V2 的 activation 與 live apply 缺少明確的 workspace readiness 契約，導致 device、source rule、mapping、database target 之間只要有任何缺口，就可能在 activation 後、runtime 中、甚至資料寫入階段才晚爆。這讓操作員很難在 Step 1 到 Step 4 完成時判斷「現在到底能不能安全啟動」。

## What Changes

- 建立 Studio V2 workspace readiness 契約，統一描述 Step 1 到 Step 4 的 blocking 與 warning issues。
- 將 activation 與 live apply 的 gating 移到同一套 readiness issue model，而不是各自臨時判斷。
- 讓 shell、summary、Step 4 activation surface 顯示相同的 readiness 結論與 issue code。
- 明確區分「可保存但不可啟動」與「可啟動但有警告」兩種狀態。

## Capabilities

### New Capabilities

- `studio-v2-workspace-readiness`: 定義 workspace completeness、issue severity、activation gating、live-apply gating 與 step-scoped readiness 回報。

### Modified Capabilities

- `studio-v2-workspace`: workspace bootstrap 與 shell 需暴露同一份 readiness summary。
- `source-rule-runtime`: source rule、derived points、tags、mappings、database targets 的缺口需在 activation 前轉成 readiness issues。
- `datalink-workbench-v2-step4-database`: Step 4 activation surface 需顯示 readiness blockers 與 warning。
- `studio-v2-live-config-apply`: live apply 需遵守 readiness gating，而不是在 runtime 中靜默接受不完整設定。

## Impact

- Affected specs:
  - New: studio-v2-workspace-readiness
  - Modified: studio-v2-workspace, source-rule-runtime, datalink-workbench-v2-step4-database, studio-v2-live-config-apply
- Affected code:
  - Modified: internal/datalink/workspace/service.go
  - Modified: internal/datalink/workspace/service_activation.go
  - Modified: internal/datalink/sourcerule/activation_readiness.go
  - Modified: internal/api/handlers/studio_v2_workspace_activation_handler.go
  - Modified: internal/api/handlers/studio_v2_workspace_handler.go
  - Modified: internal/api/handlers/studio_v2_runtime_apply.go
  - Modified: frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step4/CommitSummary.tsx
  - Modified: frontend/src/hooks/datalink/useStudioV2WorkspaceActivation.ts
