# 2026-04-09 API Inventory and Integration

## 目標

盤點目前 backend API、frontend integration layer，以及 AI 產版後如何接回現有 backend。

## 現有 backend API inventory

主要入口在 [router.go](/Users/yishow/prj/go_gateway/internal/api/router.go)。

### `/api/v1/test/*`

服務 `/test`：

- `POST /test/connect`
- `POST /test/disconnect`
- `GET /test/status`
- `POST /test/read`
- `POST /test/write`
- `POST /test/batch`
- `POST /test/script`
- `GET /test/scripts`
- `POST /test/scripts`
- `DELETE /test/scripts/:id`
- `POST /test/monitor/start`
- `POST /test/monitor/stop`
- `GET /test/monitor/stream`

### `/api/v1/debug/*`

- `GET /debug/packets`
- `GET /debug/logs`
- `DELETE /debug/clear`
- `POST /debug/send-raw`
- `GET /debug/analyze/:packetId`

### `/api/v1/datalink/*`

服務 `/studio`：

- health
- dashboard stats / device statuses
- runtime status / stream
- protocols
- polling groups CRUD
- devices CRUD / `test-draft` / `test` / `activate` / `disable` / `test-batch`
- source-rules CRUD / `enable` / `disable` / `candidates` / `recompute`
- source-rules tag apply
- source-rules database output apply
- source-rules local modbus apply
- source-rules tag review decisions
- points CRUD / batch create / poll
- tags CRUD / activate / retire / batch / validate-key
- mappings CRUD / preview / validate-pipeline
- settings
- preview stream
- modbus-share status / lifecycle / mappings / sync / write-tag-value
- db-target connectors / test / tables / validate / schema generate / dry-run / write-history / mappings CRUD

## 現有 frontend integration

### `/studio`

主要 service layer：

- [datalink.ts](/Users/yishow/prj/go_gateway/frontend/src/services/datalink.ts)

主要 hook layer：

- [useDevices.ts](/Users/yishow/prj/go_gateway/frontend/src/hooks/datalink/useDevices.ts)
- [usePoints.ts](/Users/yishow/prj/go_gateway/frontend/src/hooks/datalink/usePoints.ts)
- [useTags.ts](/Users/yishow/prj/go_gateway/frontend/src/hooks/datalink/useTags.ts)
- [useMappings.ts](/Users/yishow/prj/go_gateway/frontend/src/hooks/datalink/useMappings.ts)
- [useSourceRules.ts](/Users/yishow/prj/go_gateway/frontend/src/hooks/datalink/useSourceRules.ts)

現有 integration chain：

`route -> page -> hook -> service -> API`

### `/test`

主要 service layer：

- [api.ts](/Users/yishow/prj/go_gateway/frontend/src/services/api.ts)

現有 API wrapper：

- `useTestAPI()`
- `useDebugAPI()`

## 目前 service 覆蓋摘要

### `datalink.ts`

已封裝：

- `deviceAPI`
- `pointAPI`
- `sourceRuleAPI`
- `pollingGroupAPI`
- `tagAPI`
- `mappingAPI`
- `settingsAPI`
- `protocolAPI`
- `healthAPI`
- `dashboardAPI`
- `runtimeAPI`
- `modbusShareAPI`
- `dbTargetAPI`

### `api.ts`

已封裝：

- connect / disconnect / status
- read / write / batch
- monitor start / stop
- debug packets / logs / clear

## 目前可以直接重用的 backend 契約

我建議優先保留這些契約不動：

- device CRUD 與 test 流
- source-rule CRUD 與 candidates
- tag CRUD / validate-key
- mapping preview / validate-pipeline
- modbus-share mapping lifecycle
- db-target connector / table / validate / schema / dry-run
- `/test` 的 connect/read/write/batch/monitor/debug

## AI 產版後的接回原則

- 優先保留現有 backend contract
- 不要讓 AI 產出的 screen 直接寫 raw `fetch`
- 透過既有 `services/*.ts` 與 TanStack Query hooks 承接
- 若新 IA 與舊資料結構不對齊，先加 view-model adapter
- 只有 adapter 無法合理承接時，才提 backend change

## 建議 adapter 層

### `/studio`

- `DeviceCapabilityAdapter`
- `SourceCoverageAdapter`
- `SemanticRefinementAdapter`
- `DeliveryGroupSuggestionAdapter`
- `DatabaseRowPlannerAdapter`
- `LocalModbusBlockPlannerAdapter`
- `MqttMessagePlannerAdapter`

### `/test`

- `EngineerConnectionSessionAdapter`
- `PacketLogAdapter`
- `OperationResultAdapter`
- `MonitorStateAdapter`

## 接回步驟

1. AI 先產 route shell 與 screen skeleton
2. 明確 screen 使用哪些 domain objects
3. 建 adapter 層
4. adapter 接現有 hooks / services
5. 補 mutation、dry-run、error/retry wiring
6. 跑 targeted tests、typecheck、build

## `/studio` screen-to-contract 建議

### Device

直接重用：

- `deviceAPI`
- `useDevicesQuery`
- create / update / delete / test mutations

### Source

直接重用：

- `sourceRuleAPI`
- `pointAPI`
- runtime stream

可能新增 adapter：

- source rule planner to lattice view-model

### Tag

直接重用：

- `tagAPI`
- `mappingAPI`
- source-rule candidates
- tag review decisions

需要 adapter：

- candidate -> semantic editor model
- tag groups
- destination-readiness summary

### Destination / Database

直接重用：

- `dbTargetAPI`
- source-rule candidates where applicable

需要 adapter：

- delivery group -> row planner model
- schema preview / dry-run summary

### Destination / Local Modbus

直接重用：

- `modbusShareAPI`

需要 adapter：

- delivery group -> register block plan
- conflict heatmap summary

### Destination / MQTT

目前 router 中沒有獨立 MQTT destination UI contract。

建議：

- 第一輪先在設計文件中規劃 MQTT workspace
- 實作前再盤點目前 backend 是否已有可重用 publish / connector contract
- 若沒有，再以 OpenSpec change 補正式需求與 API

## `/test` screen-to-contract 建議

### 可直接重用

- `useTestAPI()`
- `useDebugAPI()`
- monitor stream

### 需要 adapter

- connection session summary
- command history / quick action groups
- packet/log correlation view

## 成功標準

AI 產版後，不是看畫面長得像不像，而是看：

- route 能否成立
- 現有 API 是否可承接
- mutation / dry-run / retry 是否落地
- typecheck / build / tests 是否通過

也就是 skill 的終點應該定義為：

`frontend generated + backend reconnected + verification completed`
