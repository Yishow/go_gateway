## Why

即使 Step 1 到 Step 4、runtime、database output 都補齊，如果操作員仍然看不到資料在哪一段斷掉、哪一次 activation 失敗、哪個設定變更引發 drift、哪次 write flush 失敗，那系統仍然不算「可觀察、可記錄」。目前這些資訊大多散在 server logs 或根本沒有結構化記錄，無法形成 operator-facing 診斷與 audit trail。

## What Changes

- 建立 runtime flow diagnostics 契約，讓 collector、mapping、runtime、database delivery 的階段性狀態可查詢。
- 建立 workspace audit history 契約，記錄 setup save、activation、runtime transitions、database delivery failures 等關鍵事件。
- 讓 runtime-facing UI 與 workspace surfaces 能查看最近失敗、最近成功、最近 activation 與關鍵 audit history，而不是依賴原始 logs。

## Capabilities

### New Capabilities

- `runtime-flow-diagnostics`: 定義 collector -> mapping -> runtime -> database delivery 的 operator-facing diagnostics contract。
- `workspace-audit-history`: 定義 workspace setup、activation、runtime lifecycle、database delivery 關鍵事件的 audit contract。

### Modified Capabilities

- `runtime-dashboard-backend`: backend contract 必須支援 diagnostics summary 與最近失敗資訊。
- `runtime-dashboard`: dashboard 必須顯示 operator-facing diagnostics，而不是只顯示樂觀 summary。
- `studio-v2-workspace`: workspace surfaces 必須能查詢最近 activation 與 audit history。

## Impact

- Affected specs:
  - New: runtime-flow-diagnostics, workspace-audit-history
  - Modified: runtime-dashboard-backend, runtime-dashboard, studio-v2-workspace
- Affected code:
  - Modified: internal/datalink/runtime/service.go
  - Modified: internal/datalink/runtime/status.go
  - Modified: internal/datalink/dbtarget/writer.go
  - Modified: internal/api/handlers/runtime_handler.go
  - Modified: internal/api/handlers/runtime_stream_handler.go
  - Modified: internal/api/handlers/studio_v2_workspace_handler.go
  - Modified: internal/api/handlers/studio_v2_workspace_activation_handler.go
  - Modified: frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx
  - Modified: frontend/src/features/datalink/runtime-dashboard/RealtimeLogsPanel.tsx
  - Modified: frontend/src/features/datalink/runtime-dashboard/RuntimeSummaryPanel.tsx
  - Modified: frontend/src/features/datalink/runtime-dashboard/useRuntimeDashboardState.ts
  - Modified: frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx
  - New: internal/datalink/audit/
  - New: internal/api/handlers/studio_v2_workspace_audit_handler.go
