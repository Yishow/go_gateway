# `/studio` Pre-delete Mainline Inventory

> This document is retained as historical pre-delete evidence only. `/studio` is no longer a product fallback or supported dedicated surface. The immediate deletion scope, preserved surfaces, and Git rollback limits are recorded in [the retirement record](../../releases/retire-legacy-studio-and-polish-v2.md).

## 定位

- Route: `/studio`
- 受眾: 進階操作者、維運人員、需要完整控制的人
- 狀態: legacy dedicated surface，待立即刪除
- 問題: 功能完整，但資訊密度高、心智負擔高，後續仍需要大幅整理
- 目前策略: **立即刪除 dedicated route 與 proven legacy-only graph members**
- 未來目標: 不重新建立 `/studio`；刪除後與任意 unknown route 共用既有 generic unknown-route policy

## Pre-delete 說明

這份文件保留 `/studio` 的完整功能盤點，作為刪除前的 route/import/owner 追溯證據，避免刪除時誤移除 `/studio/v2`、`/studio/runtime`、`/test` 或 `/gateway/*` 共用的模組。它不再描述現行產品能力或未來重整計畫。

## Route / Shell / 共用狀態

| 區塊 | 前端實作 | 行為 | API / 狀態來源 |
| --- | --- | --- | --- |
| Route sync | `DatalinkWorkbenchPage.tsx` | 以 `?step=`、`?target=` 維持深連結 | 無直接 API；同步 `WorkbenchProvider` state |
| Cross-step context | `WorkbenchProvider.tsx` | 管理 `selectedDeviceId`、`focusedRuleId`、`activeOutputTarget`、inspector selection | 前端 context |
| Summary / diagnostics | `useWorkbenchSummary.ts`、`useWorkbenchShellDiagnostics.ts` | 跨步驟 readiness、刷新狀態、blocker 摘要 | `devices`、`points`、`tags`、`mappings` queries |

## Step 1: Device

### 使用者可見動作

- 列出設備
- 搜尋 / 依 protocol 篩選
- 建立設備
- 編輯 / clone 設備
- 測試 draft connection
- 觀察 recent test 與 diagnostics

### 前端主要檔案

- `frontend/src/pages/datalink/workbench/MuiDeviceStep.tsx`
- `frontend/src/pages/datalink/workbench/WorkbenchDeviceStep.tsx`
- `frontend/src/pages/datalink/workbench/MuiDeviceEditor.tsx`

### 對應 hooks / API

| 類型 | 前端 | 後端 API | 狀態 |
| --- | --- | --- | --- |
| list devices | `useDevicesQuery()` | `GET /api/v1/datalink/devices` | `wired` |
| create device | `useCreateDeviceMutation()` | `POST /api/v1/datalink/devices` | `wired` |
| update device | `useUpdateDeviceMutation()` | `PUT /api/v1/datalink/devices/:id` | `wired` |
| test draft connection | `useTestDraftConnectionMutation()` | `POST /api/v1/datalink/devices/test-draft` | `wired` |
| test persisted connection | `useTestConnectionMutation()` | `POST /api/v1/datalink/devices/:id/test` | `exists-not-wired` in主要表單，部分輔助區塊會用 |
| activate / disable device | `useToggleDeviceStatusMutation()` | `POST /api/v1/datalink/devices/:id/activate` / `disable` | `wired`，但主要在 source/runtime cluster 使用 |
| readiness check | `useCheckReadinessMutation()` | `POST /api/v1/datalink/devices/:id/readiness` | `exists-not-wired` |
| protocol list | `protocolAPI.list()` | `GET /api/v1/datalink/protocols` | `exists-not-wired` |

### 維護觀察

- `draft test` 與 `persisted test` 仍是兩套語意，文件與 UI 必須明確區分。
- `readiness` API 已存在；這是 legacy pre-delete 主線的歷史缺口，不代表刪除後仍有 `/studio` owner。

## Step 2: Source

### 使用者可見動作

- 選取設備後建立 / 編輯 / 刪除 source rules
- enable / disable rule
- 以格狀畫布規劃位址
- 建立 / 刪除 points
- 觀看 runtime live values 與 collector 摘要
- 管理 template、覆蓋區段、衝突與 coverage

### 前端主要檔案

- `frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx`
- `frontend/src/pages/datalink/workbench/WorkbenchSourceRuntimeCluster.tsx`
- `frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx`

### 對應 hooks / API

| 類型 | 前端 | 後端 API | 狀態 |
| --- | --- | --- | --- |
| list source rules | `useSourceRulesQuery()` | `GET /api/v1/datalink/source-rules` | `wired` |
| create source rule | `useCreateSourceRuleMutation()` | `POST /api/v1/datalink/source-rules` | `wired` |
| update source rule | `useUpdateSourceRuleMutation()` | `PUT /api/v1/datalink/source-rules/:id` | `wired` |
| delete source rule | `useDeleteSourceRuleMutation()` | `DELETE /api/v1/datalink/source-rules/:id` | `wired` |
| enable / disable rule | `useEnableSourceRuleMutation()` / `useDisableSourceRuleMutation()` | `POST /api/v1/datalink/source-rules/:id/enable` / `disable` | `wired` |
| list points | `usePointsQuery()` | `GET /api/v1/datalink/points` | `wired` |
| create point | `useCreatePointMutation()` | `POST /api/v1/datalink/points` | `wired` |
| delete point | `useDeletePointMutation()` | `DELETE /api/v1/datalink/points/:id` | `wired` |
| runtime snapshot | `runtimeAPI.getStatus()` via `useQuery` | `GET /api/v1/datalink/runtime/status` | `wired` |
| runtime stream | `useRuntimeStream()` | `GET /api/v1/datalink/runtime/stream` | `wired` |
| toggle device collection | `useToggleDeviceStatusMutation()` | `POST /api/v1/datalink/devices/:id/activate` / `disable` | `wired` |

### 維護觀察

- Source step 已經跨進 runtime / collector 領域，不只是規劃頁。
- 規則規劃、point 生成、runtime cluster 曾在 legacy 同一步驟裡交纏；若保留的 V2/runtime surface 需要相同能力，另以現行 route contract 拆分。

## Step 3: Tag

### 使用者可見動作

- review source rule candidates
- 建立 / 刪除 tags
- 建立 / 刪除 mappings
- 批次建立 tags
- 做 tag review decision 與 recover / retry

### 前端主要檔案

- `frontend/src/pages/datalink/workbench/TagBindingStudio.tsx`
- `frontend/src/pages/datalink/workbench/SourceRuleTagReviewSurface.tsx`
- `frontend/src/pages/datalink/workbench/MuiTagIncidentDesk.tsx`

### 對應 hooks / API

| 類型 | 前端 | 後端 API | 狀態 |
| --- | --- | --- | --- |
| list tags | `useTagsQuery()` | `GET /api/v1/datalink/tags` | `wired` |
| create tag | `useCreateTagMutation()` | `POST /api/v1/datalink/tags` | `wired` |
| delete tag | `useDeleteTagMutation()` | `DELETE /api/v1/datalink/tags/:id` | `wired` |
| batch create tags | `tagAPI.batchCreate()` | `POST /api/v1/datalink/tags/batch` | `wired` |
| list mappings | `useMappingsQuery()` | `GET /api/v1/datalink/mappings` | `wired` |
| create mapping | `useCreateMappingMutation()` | `POST /api/v1/datalink/mappings` | `wired` |
| delete mapping | `useDeleteMappingMutation()` | `DELETE /api/v1/datalink/mappings/:id` | `wired` |
| source rule candidates | `useSourceRuleCandidatesQuery()` | `GET /api/v1/datalink/source-rules/:id/candidates` | `wired` |
| review decisions list | `useSourceRuleTagReviewDecisionsQuery()` | `GET /api/v1/datalink/source-rules/:id/tag-review-decisions` | `wired` |
| review decision upsert | `useUpsertSourceRuleTagReviewDecisionMutation()` | `POST /api/v1/datalink/source-rules/:id/tag-review-decisions` | `wired` |
| apply tags | 目前以 review 與 tag/mapping CRUD 組裝 | `POST /api/v1/datalink/source-rules/:id/tags/apply` | `exists-not-wired` |

### 維護觀察

- Tag step 同時存在「直接 CRUD」與「review/apply contract」兩套路徑。
- `apply tags` API 已存在，但 legacy 前端曾大量用明細 CRUD 組裝流程；若 V2 需要，另以 V2 contract 定義，不重新建立 `/studio`。

## Step 4: Output

### 使用者可見動作

- 在 Local Modbus 上規劃 register
- 啟停 local modbus share server
- 同步 mapping 到 modbus share
- 寫入 tag value 做測試
- 建立 / 編輯 / 測試 DB connector
- 檢查 schema / tables / validation
- 建立 / 更新 / 刪除 DB target mappings
- dry-run database mappings

### 前端主要檔案

- `frontend/src/pages/datalink/workbench/LocalModbusBoard.tsx`
- `frontend/src/pages/datalink/workbench/DatabaseTargetBoard.tsx`
- `frontend/src/pages/datalink/workbench/SourceRuleDatabaseTargetBoard.tsx`
- `frontend/src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx`

### 對應 hooks / API

| 類型 | 前端 | 後端 API | 狀態 |
| --- | --- | --- | --- |
| modbus share status | `modbusShareAPI.status()` | `GET /api/v1/datalink/modbus-share/status` | `wired` |
| start / stop share | `modbusShareAPI.start()` / `stop()` | `POST /api/v1/datalink/modbus-share/start` / `stop` | `wired` |
| list share mappings | `modbusShareAPI.listMappings()` | `GET /api/v1/datalink/modbus-share/mappings` | `wired` |
| upsert share mapping | `modbusShareAPI.upsertMapping()` | `PUT /api/v1/datalink/modbus-share/mappings/:tagId` | `wired` |
| delete share mapping | `modbusShareAPI.deleteMapping()` | `DELETE /api/v1/datalink/modbus-share/mappings/:tagId` | `wired` |
| sync share mappings | `modbusShareAPI.sync()` | `POST /api/v1/datalink/modbus-share/sync` | `wired` |
| write tag value | `modbusShareAPI.writeTagValue()` | `POST /api/v1/datalink/modbus-share/write-tag-value` | `wired` |
| list DB connectors | `dbTargetAPI.listConnectors()` | `GET /api/v1/datalink/db-targets/connectors` | `wired` |
| create / update / delete connector | `dbTargetAPI.createConnector()` 等 | `POST/PUT/DELETE /api/v1/datalink/db-targets/connectors...` | `wired` |
| test connector | `dbTargetAPI.testConnector()` | `POST /api/v1/datalink/db-targets/connectors/:id/test` | `wired` |
| list tables | `dbTargetAPI.listTables()` | `GET /api/v1/datalink/db-targets/connectors/:id/tables` | `wired` |
| validate connector | `dbTargetAPI.validateConnector()` | `GET /api/v1/datalink/db-targets/connectors/:id/validate` | `wired` |
| generate schema | `dbTargetAPI.generateSchema()` | `POST /api/v1/datalink/db-targets/connectors/:id/schema/generate` | `wired` |
| dry-run mappings | `dbTargetAPI.dryRunMappings()` | `POST /api/v1/datalink/db-targets/connectors/:id/mappings/dry-run` | `wired` |
| list / create / update / delete DB mappings | `dbTargetAPI.*Mapping()` | `GET/POST/PUT/DELETE /api/v1/datalink/db-targets/mappings...` | `wired` |
| apply database outputs | review surfaces only | `POST /api/v1/datalink/source-rules/:id/database-outputs/apply` | `exists-not-wired` |
| apply local modbus outputs | review surfaces only | `POST /api/v1/datalink/source-rules/:id/local-modbus/apply` | `exists-not-wired` |

### 維護觀察

- Output step 是 `/studio` 裡最接近正式 operator tooling 的區段。
- 但 review/apply 類 source-rule output endpoints 尚未完全成為主線，頁面仍以細粒度 CRUD 操作為主。

## `/studio` 後續整理時必記的事情

1. `step` 與 `target` query params 仍是主要 deep-link 契約。
2. `selectedDeviceId` 是跨步驟最重要的上下文。
3. Source / Tag / Output 不只共享裝置，還共享 `focusedRuleId`、candidate snapshot 與 review decisions。
4. Runtime live values 已深入 Source 與 Output 觀察面，不可把 runtime 視為獨立頁才有的能力。
5. `/studio` 的重整不應先碰 CRUD 後端，而應先整理 page boundary、shared context 與 review/apply contract。

## 暫停期間建議不要做的事

1. 不要為了配合 `/studio/v2` 而先大改 `/studio` 主線資訊架構。
2. 不要在沒有重新定義 step boundary 前，把更多 runtime 或 gateway 概念硬塞進 `/studio`。
3. 不要把 `/studio` 的複雜度問題誤判成「缺少 API」；它主要是 page responsibility 與 interaction density 問題。
