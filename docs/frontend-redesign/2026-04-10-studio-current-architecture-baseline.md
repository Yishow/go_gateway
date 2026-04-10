# 2026-04-10 Studio Current Architecture Baseline

## 目的

這份文件是 `/studio` 現況基準。

用途不是直接定義新 UI，而是固定記錄：

- repo 內目前真正存在的 `/studio` route、step 與 output target
- 現有 page shell、state ownership、service / API 對接方式
- accepted Stitch screens 與現況程式之間的對齊程度

之後做 `/studio` redesign、accepted screen 對接、Destination 分流與 frontend-backend reconnection 時，應先對照這份文件，而不是每次重新探索。

## 範圍

本文件只描述目前 repo 內 `/studio` 現況。

不處理：

- `/studio` 最終 route 重構
- Stitch 畫面直接落地實作
- backend contract 變更提案

## Route 與入口

### 前端 route

- `/studio`
  - 入口在 `frontend/src/App.tsx`
  - render：`<DatalinkWorkbenchPage />`
- legacy datalink routes 會 redirect 回 `/studio`
  - `/datalink`
  - `/datalink/workbench`
  - `/datalink/*`
  - `devices / points / polling-groups / tags / mappings / wizard / local-modbus` legacy variants

### Query state

目前 `/studio` 仍是單一路由 + query-driven shell：

- `step=device|source|tag|output`
- `target=modbus|database`

目前不存在真正獨立的：

- `Destination Hub` route
- `Database Workspace` route
- `Local Modbus Workspace` route
- `MQTT Workspace` route

### 頁面入口

- `frontend/src/pages/datalink/workbench/DatalinkWorkbenchPage.tsx`
  - `/studio` 真正頁面
  - 依 `step` 切換 `device / source / tag / output`
- `frontend/src/pages/datalink/workbench/MuiWorkbenchShell.tsx`
  - 最外層 shell
- `frontend/src/pages/datalink/workbench/MuiWorkbenchFrame.tsx`
  - 主工作台 frame

## 當前 shell 與 state ownership

### Shell owner

`WorkbenchProvider.tsx` 是 `/studio` 的主要 page-level state owner。

它目前管理：

- `activeStep`
- `selectedDeviceId`
- `inspectorSelection`
- `activeOutputTarget`
- `outputSelectionState`
- `crossStepContext`
- `sourcePlanningState`
- `devicePanelState`
- `recentDeviceTests`
- `sourceStepNotice`
- `sourceStepInspectorBanner`

### Route 同步

`DatalinkWorkbenchPage.tsx` 內的 `WorkbenchRouteStateSyncWithRouter()` 處理：

- mount 時 query params -> provider state
- hydration 後 provider state -> query params

這代表目前真正的 source of truth 是 `WorkbenchProvider`，不是 URL router 本身。

## 後端 API 現況

主要入口在 `internal/api/router.go` 的 `/api/v1/datalink/*`。

### `/studio` 直接關聯的 API 群組

- health
- dashboard
- runtime status / runtime stream / preview stream
- protocols
- devices CRUD / `test-draft` / `test` / `activate` / `disable` / `test-batch`
- polling-groups CRUD
- source-rules CRUD / `enable` / `disable` / `candidates` / `recompute`
- source-rules tag apply
- source-rules database output apply
- source-rules local modbus apply
- source-rules tag review decisions
- points CRUD / batch create / poll
- tags CRUD / activate / retire / batch / validate-key
- mappings CRUD / preview / validate-pipeline
- settings
- modbus-share status / lifecycle / mappings / sync / write-tag-value
- db-target connectors / test / tables / validate / schema / dry-run / mappings CRUD

### 重要限制

- router 目前沒有獨立的 MQTT destination planner API
- router 目前沒有獨立的 Destination Hub contract
- Database 與 Local Modbus 的正式對接能力存在，但仍被包在 `output` step 內

## 前端 integration 現況

### 主要 service layer

`frontend/src/services/datalink.ts`

目前主要封裝：

- `deviceAPI`
- `pointAPI`
- `sourceRuleAPI`
- `pollingGroupAPI`
- `tagAPI`
- `mappingAPI`
- `runtimeAPI`
- `modbusShareAPI`
- `dbTargetAPI`
- `protocolAPI`

### 主要 hook layer

`frontend/src/hooks/datalink/`

目前主要使用：

- `useDevicesQuery`
- `usePointsQuery`
- `useSourceRulesQuery`
- `useSourceRuleCandidatesQuery`
- `useSourceRuleTagReviewDecisionsQuery`
- `useTagsQuery`
- `useMappingsQuery`
- `usePollingGroupsQuery`
- `useRuntimeStream`

### 現有 integration chain

`route -> DatalinkWorkbenchPage -> step command deck / board -> hook -> service -> API`

## 步驟結構現況

### 1. Device

主要組件：

- `MuiDeviceStep`
- `MuiDeviceList`
- `MuiDeviceDetail`
- `MuiDeviceEditor`
- `MuiDeviceDiagnosticsPanel`

現有能力：

- device list / detail / create / edit / clone
- connect 與 probe diagnostics
- readiness / activation flow

現況判定：

- 已有可用的 `Device Workspace` 主體
- 但仍是單一 `/studio` shell 下的一步，不是獨立 page route

### 2. Source

主要組件：

- `MuiSourceCommandDeck`
- `SourceCanvasSection`
- `AddressCanvas`
- `SourceCanvasStatusState`
- `WorkbenchSourceRuntimeCluster`

現有能力：

- source rule 規劃
- address lattice / coverage surface
- rule candidate refresh
- runtime / preview context

現況判定：

- 已有 `Source Workspace` 主體
- 但 rule coverage 與 planner 模型仍被內嵌在 step shell 中

### 3. Tag

主要組件：

- `MuiTagCommandDeck`
- `TagBindingStudio`
- `SourceRuleTagReviewSurface`
- `SourceRuleTagReviewCandidateRow`

現有能力：

- review-first tag 建立與綁定
- candidate review decision
- mapping / tag batch 操作
- cross-step traceability 到 output

現況判定：

- 已有 `Tag Workspace` 主體
- 但目前 repo 內還沒有 Stitch 精煉版那種明確的 `Tag Group / handoff bundle / delivery hint` page-level 語意

### 4. Output

主要組件：

- `MuiOutputCommandDeck`
- `LocalModbusBoard`
- `DatabaseTargetBoard`
- `WorkbenchInspectorPanel`

現有能力：

- `target=modbus|database` 切換
- Local Modbus mapping / sync / value write / occupancy handling
- Database connector / table / schema / validation / mapping / dry-run

現況判定：

- repo 內其實已同時有 Database 與 Local Modbus 兩條實作線
- 但它們目前仍被壓在同一個 `output` step 裡
- `Destination Hub`、`Database Workspace`、`Local Modbus Workspace` 尚未成為獨立 route / workspace page
- MQTT 目前仍是設計規劃，不是現有 code route

## Accepted Stitch screens 對照

目前應作為 `/studio` 比對基準的 accepted screens：

- `3fd873d663c046edb98afc64bdfebf6e`
  - `[已接受候選] 裝置配置工作台 (Device Workspace)`
- `7aa29d10a3f640f39830f15e13e24212`
  - `[已接受候選] 來源規則規劃 (Source Workspace)`
- `93c37ae5c52b4da19b1f908effc95ff8`
  - `[已接受候選] 語義精煉工作台 (Tag Workspace - 精煉版)`
- `64f12472cee74a65bd55029eb3ea840e`
  - `[已接受候選] 目標傳輸中心 (Destination Hub)`
- `6d040a45908143248f4bca91e2c7ca57`
  - `[已接受候選] 目標輸出工作台 (Database Workspace)`
- `9bca8e1e643c4baf94e7ced109eeef3c`
  - `[已接受候選] 目標輸出工作台 (Local Modbus Workspace)`
- `8489c68de1aa4c8da6970e94a01a639d`
  - `[已接受候選] 目標發佈工作台 (MQTT Workspace)`

## 核心落差

### 目前 code 已存在但 route ownership 尚未分離

- Device
- Source
- Tag
- Output / Database surface
- Output / Local Modbus surface

### 目前只有 Stitch 規劃，尚無 code route

- Destination Hub
- MQTT Workspace

### 目前最大結構問題

- `/studio` 真正的 code 仍是 `device -> source -> tag -> output`
- 但 accepted Stitch 已把 `output` 拆成：
  - `Destination Hub`
  - `Database Workspace`
  - `Local Modbus Workspace`
  - `MQTT Workspace`

所以後續若要做逐頁對照，不可以再只把 `/studio` 當成四步驟頁，而要用「主線 step + destination family pages」雙層結構來比。

## 驗證結論

`/studio` 目前已具備：

- 主產品唯一入口
- 穩定的 step shell
- 可重用的 datalink services / hooks
- Database 與 Local Modbus 兩條正式 output contract

`/studio` 目前尚未具備：

- page-per-workspace route ownership
- Destination Hub 的正式 page shell
- MQTT workspace 的現成 backend contract

因此後續驗證與規劃應採：

1. 先以這份 baseline 固定 current code reality
2. 再以 accepted Stitch screens 作為 page-level target
3. 對缺少 code route 的頁面，直接標記為 `Stitch-planned baseline`
4. 等 UI 全部確認完畢，再做 route / adapter / backend reconnection
