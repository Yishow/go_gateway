## Why

目前 Step 4 只支援「單一目標表 + 每個 point 對應唯一欄位」的配置方式。當多條規則需要寫入同一張表，且希望重複使用同一組業務欄位、改以不同列或不同群組寫入時，現有模型會直接把它判定為欄位衝突，無法表達實際需求。

這個需求已經超出單純 UI 重排，會影響 workspace 持久化模型、Step 4 規劃語意，以及後續資料寫入契約；需要先以 spec 明確定義，避免前後端各自延伸出不相容的「共表」解讀。

## What Changes

- 在 Studio V2 Step 4 引入 row-group / write-group 規劃語意，讓不同規則可寫入同一張表，並能明確指定哪些 point 屬於同一列群組。
- 定義「可共表、可重用欄位、但必須透過 row-group 區分列實例」的合法配置規則，取代目前單純以欄位名稱重複視為衝突的判定。
- 擴充 workspace database 設定與 target 持久化契約，保存 row-group 所需的識別、群組鍵與欄位重用資訊。
- 調整 Step 4 readiness / validation / apply 流程，讓 shared-column row groups 的錯誤能以 row-group 語意呈現，而不是只顯示 generic column conflict。
- 盤點 runtime / database delivery 寫入流程對 row-group 的影響，明確定義哪些模式支援 insert、upsert，哪些情況仍然應阻擋。

## Capabilities

### New Capabilities

- `workspace-database-row-groups`: 定義 Studio V2 workspace 如何描述同表多列群組、欄位重用與群組鍵約束。

### Modified Capabilities

- `datalink-workbench-v2-step4-database`: Step 4 的欄位配置與衝突規則將從「單欄唯一綁定」擴充為支援 row-group 的共表規劃。
- `studio-v2-workspace`: workspace 的 database config / target 持久化需求將擴充，以保存 row-group 規劃狀態。

## Impact

- Affected specs: `workspace-database-row-groups`, `datalink-workbench-v2-step4-database`, `studio-v2-workspace`
- Affected code:
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx`
  - Modified: `frontend/src/features/datalink/workbench-v2/state/types.ts`
  - Modified: `frontend/src/pages/datalink/workbench-v2/useStudioV2DatabaseAutosave.ts`
  - Modified: `frontend/src/services/studioV2WorkspaceDatabase.ts`
  - Modified: `frontend/src/types/datalink.ts`
  - Modified: `internal/api/handlers/studio_v2_workspace_database_handler.go`
  - Modified: `internal/api/handlers/studio_v2_workspace_database_handler_helpers.go`
  - Modified: `internal/datalink/workspace/service_database.go`
  - Modified: `internal/datalink/dbtarget/types.go`
  - Modified: `internal/datalink/dbtarget/service.go`
  - Modified: `internal/datalink/runtime/service.go`
  - New: `openspec/specs/workspace-database-row-groups/spec.md`
