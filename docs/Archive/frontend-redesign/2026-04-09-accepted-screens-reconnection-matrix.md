# 2026-04-09 Accepted Screens Reconnection Matrix

## 目的

整理目前 Google Stitch 專案中已標記為 `[已接受候選]` 的 `/studio` 與 `/test` 畫面，判斷它們與本 repo 現有 backend contract 的對齊程度，並定義前端實作時的接線順序。

這份文件不是在評估美術，而是在回答：

- 哪些 accepted screen 已可視為穩定方向
- 哪些可以直接重用現有 API
- 哪些需要 adapter 或 view-model
- 哪些目前仍只有設計規劃，不能直接進 implementation

## Accepted screen set

### `/studio`

| Screen | Stitch screen id | 建議 route |
| --- | --- | --- |
| Device | `3fd873d663c046edb98afc64bdfebf6e` | `/studio?step=device` |
| Source | `7aa29d10a3f640f39830f15e13e24212` | `/studio?step=source` |
| Tag | `93c37ae5c52b4da19b1f908effc95ff8` | `/studio?step=tag` |
| Destination Hub | `64f12472cee74a65bd55029eb3ea840e` | `planned: /studio/destination` |
| Database Workspace | `6d040a45908143248f4bca91e2c7ca57` | `current: /studio?step=output&target=database` |
| Local Modbus Workspace | `9bca8e1e643c4baf94e7ced109eeef3c` | `current: /studio?step=output&target=modbus` |
| MQTT Workspace | `8489c68de1aa4c8da6970e94a01a639d` | `planned: /studio/share/mqtt` |

### `/test`

| Screen | Stitch screen id | 建議 route / mode |
| --- | --- | --- |
| Command Lab | `72df01403b9a4fe2942e4cdc48b400ef` | `/test?mode=command-lab` |
| Live Monitor | `0c95dcd8c89a46a691317f200800c322` | `/test?mode=live-monitor` |
| Device Scan | `1c56fb38b066459397286e695adc35c6` | `/test?mode=scan` |
| RTU Polling Jobs | `521bb617cef24f8baebcc90a7b5a8a48` | `/test?mode=rtu-polling` |

## 分類規則

- `A. 可直接進接線`
  - 可直接重用現有 service / hook / API
  - 只需要正常的 page shell 與元件搬遷
- `B. 可進接線，但需要 adapter`
  - 現有 API 足夠，但 Stitch 畫面的主操作模型與現有資料結構不完全一致
  - 需要先建 view-model 或 planner adapter
- `C. 先保留設計，不直接落地`
  - backend contract 還沒有 ready
  - 或 route / workspace ownership 還不完整

## Screen-by-screen matrix

| Route / screen | 主操作物件 | 直接重用 API / hook | 仍需 adapter / bridge | 判定 | 備註 |
| --- | --- | --- | --- | --- | --- |
| `/studio?step=device` Device | `Device` | `deviceAPI`, `protocolAPI`, `useDevicesQuery`, create/update/test mutations | `DeviceCapabilityAdapter` | B | accepted screen 已穩，主要工作是把 connect、probe、traits、readiness 聚合成 page-level capability view |
| `/studio?step=source` Source | `Source Rule`, `Point`, coverage lattice | `sourceRuleAPI`, `pointAPI`, runtime stream / preview stream, `useSourceRules*`, `usePoints*` | `SourceCoverageAdapter` | B | accepted screen 與 code 主題一致，但 lattice 與 coverage 仍需要 planner projection |
| `/studio?step=tag` Tag | `Tag`, `Mapping`, `Tag Group`, handoff bundle | `tagAPI`, `mappingAPI`, source-rule candidates, tag review decisions, `useTags*`, `useMappings*` | `SemanticRefinementAdapter`, `DeliveryGroupSuggestionAdapter` | B | 以精煉版 accepted screen 為準；code 已有 review-first 雛型，但還缺 tag-group / delivery hint 聚合層 |
| `planned: /studio/destination` Destination Hub | `Destination family`, handoff hub | `dbTargetAPI`, `modbusShareAPI`, source-rule database output apply, source-rule local modbus apply | `DestinationHubAdapter` | C | 目前 repo 內沒有正式 hub route；應視為 Stitch-planned baseline，不可誤判為已實作 |
| `current: /studio?step=output&target=database` Database Workspace | database connector, schema planner, row mapping | `dbTargetAPI`, source-rule database output apply | `DatabaseRowPlannerAdapter` | B | code 能力已存在，但仍內嵌在 output step 中，尚未獨立成 page route |
| `current: /studio?step=output&target=modbus` Local Modbus Workspace | register planner, block occupancy, share runtime | `modbusShareAPI`, source-rule local modbus apply | `LocalModbusBlockPlannerAdapter` | B | code 能力已存在，但 route / workspace identity 尚未分離 |
| `planned: /studio/share/mqtt` MQTT Workspace | publish plan, topic mapping, share routing | existing connector reality only, no dedicated planner API | `MqttMessagePlannerAdapter` | C | 現有 repo 有 MQTT 方向需求，但 router 尚無獨立 planner contract；先保留設計與文件 |
| `/test?mode=command-lab` Command Lab | `Connection Session`, command execution | `useTestAPI()`, `useDebugAPI()`, `/test/connect`, `/disconnect`, `/status`, `/read`, `/write`, `/batch`, debug packets/logs | `EngineerConnectionSessionAdapter`, `OperationResultAdapter`, `PacketLogAdapter` | B | shell 已穩，可進前端重構；需把現有卡片模型重組成 session-centered console |
| `/test?mode=live-monitor` Live Monitor | `Monitor Session`, live value watchlist | `useTestAPI()`, `/test/monitor/start`, `/test/monitor/stop`, `/test/monitor/stream`, debug packets/logs | `MonitorStateAdapter`, `EngineerConnectionSessionAdapter`, `PacketLogAdapter` | B | 畫面方向已對；需把 current values、事件流、watchlist、packet correlation 串起來 |
| `/test?mode=scan` Device Scan | scan job, scan result, scan diagnostics | 現有 `DeviceScanner` 相關前端流程與 `/test` 連線能力 | `ScanSessionAdapter` | B | accepted screen 與現有工具相容，但要把掃描流程正式掛入新的 `/test` shell 模式切換 |
| `/test?mode=rtu-polling` RTU Polling Jobs | polling job, scheduler state, result timeline | 現有 `RTUPollingCard` 前端能力與 `/test` 連線上下文 | `PollingJobAdapter`, `EngineerConnectionSessionAdapter` | B | 可做 mode refactor；但仍要盤點目前 RTU polling 的實際資料來源與狀態模型 |

## 目前不建議直接落地的部分

### Destination Hub / MQTT workspace

判定：`C. 先保留設計，不直接落地`

原因：

- `Destination Hub` 目前沒有正式 code route 或 page shell
- `MQTT Workspace` 目前沒有獨立 planner contract
- 現有 repo 中雖有 output 實作與部分 MQTT 方向需求，但不等於已具備可直接承接 accepted page 的 workspace 契約

建議：

- 先保留設計、文件與 prompt
- 先完成 all-pages baseline 對照
- 實作前再盤點是否以現有 connector / share contract 擴成 destination contract
- 若不夠，再走 OpenSpec change

## 實作優先順序

### Wave 1：先做 shell 與最穩定主線

1. `/test` shell
2. `Command Lab`
3. `Live Monitor`

原因：

- `/test` accepted 畫面方向已穩
- 可直接改善目前 `TestPage.tsx` 的大單頁卡片問題
- 先把 `session-centered shell` 落地，後面 `Scan` 與 `RTU Polling` 只要接同一個骨架

### Wave 2：補 `/studio` 主線骨架

1. `Device`
2. `Source`
3. `Tag`

原因：

- 這三步是主產品最核心的 workflow
- 現有 datalink hook / service 已可重用
- 主要工作會落在 adapter 與工作台 state ownership

### Wave 3：Destination 分流

1. `Destination Hub`
2. `Database Workspace`
3. `Local Modbus Workspace`
4. `MQTT Workspace` planning baseline

原因：

- 先把 accepted page family 跟 current output step 的對照固定下來
- Database / Local Modbus 先走 page identity 拆分規劃
- Destination Hub / MQTT 維持 Stitch-planned baseline，等 UI 確認完再談接線

### Wave 4：補工程進階模式

1. `Device Scan`
2. `RTU Polling Jobs`
3. `scripts / templates / raw send / packet analyze` 的 diagnostics extension

原因：

- 這些能力屬 `/test` 進階工程工具
- 不應阻塞主 shell 與主產品主線先落地

## 每張 accepted screen 的建議接線鏈

### `/studio`

建議統一採：

`route -> page shell -> screen adapter -> TanStack Query hook -> service -> API`

對應 adapter：

- Device: `DeviceCapabilityAdapter`
- Source: `SourceCoverageAdapter`
- Tag: `SemanticRefinementAdapter`
- Destination Hub: `DestinationHubAdapter`
- Database: `DatabaseRowPlannerAdapter`
- Local Modbus: `LocalModbusBlockPlannerAdapter`
- MQTT: `MqttMessagePlannerAdapter`

### `/test`

建議統一採：

`route -> test shell -> mode adapter -> service layer / stream state -> API`

對應 adapter：

- `/test` 共用：`EngineerConnectionSessionAdapter`
- Command Lab: `OperationResultAdapter`
- Live Monitor: `MonitorStateAdapter`
- 共用 diagnostics：`PacketLogAdapter`
- Device Scan: `ScanSessionAdapter`
- RTU Polling: `PollingJobAdapter`

## 建議的 route ownership

### `/studio`

- `/studio`
  - step shell
- `/studio?step=device`
- `/studio?step=source`
- `/studio?step=tag`
- `planned: /studio/destination`
- `planned: /studio/destination/database`
- `planned: /studio/destination/local-modbus`
- `planned: /studio/share/mqtt`

### `/test`

- `/test`
  - 保留單一路由
  - 用內部 mode switch 切 `command-lab`, `live-monitor`, `scan`, `rtu-polling`

原因：

- 這與目前文件定義的 `Field Engineer Debug Console` 一致
- 也最符合 accepted `/test` screens 的同殼切模式方向

## 是否還要繼續修 UI

目前判斷：

- 不建議再做大方向重生
- 只建議做 focused refinement

僅建議保留兩個 UI refinement 任務：

1. `Tag` 再補強 `Tag Group / Delivery Group suggestion / readiness`
2. `Destination Hub` 補家族導流與下游任務層級感

其餘 accepted screens 可視為進入實作規劃階段。

## 成功標準

accepted screen 不等於已完成。

真正可進 implementation 的標準是：

- route ownership 穩定
- accepted screen 對應的 domain object 明確
- service / hook / API 可承接
- adapter 需求已被列出
- 無 backend contract 的部分被明確標示為設計保留

目前這 11 張 accepted screen 已達到：

- 可作為前端重構目標
- 可作為 route shell 與 workspace skeleton 基準
- 可作為 backend reconnection 排程依據

但尚未達到：

- 直接下載 HTML 即可上線
- 不經 adapter 就可完整接線
- 所有 destination 分支都已有完整 backend contract
