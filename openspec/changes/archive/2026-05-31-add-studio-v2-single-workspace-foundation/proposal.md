## Why

`/studio/v2` 要承接逐步自動存、多台設備交接與 runtime 工作區視角，前提是後端必須先有一個真實存在的 V2 工作區。目前 repo 還沒有這個資料邊界，因此 V2 只能停留在前端本地 state，無法成為真正可持久化的主線。

## What Changes

- 新增唯一且可持久化的 V2 工作區資料模型
- 新增讀取即自動建立的 workspace bootstrap contract
- 定義 workspace 與 `/studio` legacy 主線隔離，不自動匯入舊資料
- 補齊 frontend 讀取 workspace boot 狀態所需的 service / hook 契約

## Capabilities

### New Capabilities

- `studio-v2-workspace`: 定義唯一 V2 工作區的建立、保留、隔離與 bootstrap 行為

### Modified Capabilities

- `datalink-api`: 新增 V2 工作區 bootstrap API 與回應契約

## Impact

- Affected specs: `studio-v2-workspace`, `datalink-api`
- Affected code:
  - New: `internal/datalink/workspace/`, `internal/api/handlers/studio_v2_workspace_handler.go`, `frontend/src/services/studioV2Workspace.ts`, `frontend/src/hooks/datalink/useStudioV2Workspace.ts`
  - Modified: `internal/api/router.go`, `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`, `frontend/src/types/datalink.ts`
  - Removed: none
