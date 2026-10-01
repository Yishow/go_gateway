# 2026-04-10 All Pages Verification Baseline

## 目的

這份文件把 `/test` 的 baseline 驗證方式擴展到所有頁面。

之後不論是：

- review 現有前端
- 比對 accepted Stitch screens
- 規劃尚未存在於 code 的新頁
- 準備後續 frontend-backend reconnection

都應先對照這份文件，而不是每次重新探索 repo 與 Stitch 專案。

## 驗證規則

每一頁都要固定記錄：

- page identity
- current route ownership
- current page / state owner
- current service / API basis
- Stitch accepted page basis
- 驗證結論
- 下一步規劃方式

若 repo 內還沒有該頁：

- 不可標成「不存在就略過」
- 必須以 accepted Stitch page 作為 `Stitch-planned baseline`

## Page Inventory Summary

| Page | 當前狀態 | 當前 owner | Stitch basis | 驗證結論 |
| --- | --- | --- | --- | --- |
| Device Workspace | code 已存在 | `/studio?step=device` | `3fd873d663c046edb98afc64bdfebf6e` | 可直接比對 |
| Source Workspace | code 已存在 | `/studio?step=source` | `7aa29d10a3f640f39830f15e13e24212` | 可直接比對 |
| Tag Workspace | code 已存在 | `/studio?step=tag` | `93c37ae5c52b4da19b1f908effc95ff8` | 可直接比對，但以精煉版為準 |
| Destination Hub | Stitch-only | 無正式 code route | `64f12472cee74a65bd55029eb3ea840e` | 以 Stitch 規劃為基線 |
| Database Workspace | code 已存在但仍內嵌 output | `/studio?step=output&target=database` | `6d040a45908143248f4bca91e2c7ca57` | 可比對，但 route 尚未分離 |
| Local Modbus Workspace | code 已存在但仍內嵌 output | `/studio?step=output&target=modbus` | `9bca8e1e643c4baf94e7ced109eeef3c` | 可比對，但 route 尚未分離 |
| MQTT Workspace | Stitch-only | 無正式 code route | `8489c68de1aa4c8da6970e94a01a639d` | 以 Stitch 規劃為基線 |
| Command Lab | code 已存在 | `/test` 單殼內 command area | `72df01403b9a4fe2942e4cdc48b400ef` | 可直接比對 |
| Live Monitor | code 已存在 | `/test` 單殼內 monitor area | `0c95dcd8c89a46a691317f200800c322` | 可直接比對 |
| Device Scan | code 已存在 | `/test` 單殼內 scan area | `1c56fb38b066459397286e695adc35c6` | 可直接比對 |
| RTU Polling Jobs | code 已存在 | `/test` 單殼內 polling area | `521bb617cef24f8baebcc90a7b5a8a48` | 可直接比對 |

## `/studio` Pages

### 1. Device Workspace

- Current owner：
  - route：`/studio?step=device`
  - page：`DatalinkWorkbenchPage -> MuiDeviceStep`
  - state：`WorkbenchProvider.selectedDeviceId / devicePanelState / recentDeviceTests`
- Current API basis：
  - `deviceAPI`
  - `protocolAPI`
  - `useDevicesQuery`
  - device connection / readiness mutations
- Stitch basis：
  - `3fd873d663c046edb98afc64bdfebf6e`
  - `[已接受候選] 裝置配置工作台 (Device Workspace)`
- 驗證結論：
  - 這頁已存在於 code，可直接做 UI / IA 對照
  - 真正缺的是 page-level shell 整理，不是 backend contract

### 2. Source Workspace

- Current owner：
  - route：`/studio?step=source`
  - page：`DatalinkWorkbenchPage -> MuiSourceCommandDeck`
  - state：`WorkbenchProvider.sourcePlanningState / sourceStepNotice / sourceStepInspectorBanner`
- Current API basis：
  - `sourceRuleAPI`
  - `pointAPI`
  - `useSourceRulesQuery`
  - `useSourceRuleCandidatesQuery`
  - `useRuntimeStream`
- Stitch basis：
  - `7aa29d10a3f640f39830f15e13e24212`
  - `[已接受候選] 來源規則規劃 (Source Workspace)`
- 驗證結論：
  - 這頁已存在於 code，可直接做規則畫布與 shell 對照
  - 後續主要工作會落在 planner projection 與 coverage presentation

### 3. Tag Workspace

- Current owner：
  - route：`/studio?step=tag`
  - page：`DatalinkWorkbenchPage -> MuiTagCommandDeck`
  - state：`WorkbenchProvider.crossStepContext / inspectorSelection`
- Current API basis：
  - `tagAPI`
  - `mappingAPI`
  - `useTagsQuery`
  - `useMappingsQuery`
  - `useSourceRuleTagReviewDecisionsQuery`
- Stitch basis：
  - `93c37ae5c52b4da19b1f908effc95ff8`
  - `[已接受候選] 語義精煉工作台 (Tag Workspace - 精煉版)`
- 驗證結論：
  - 這頁已存在於 code，可直接做 review-first 對照
  - 後續比對應以精煉版為準，不再回退到舊 accepted tag screen

### 4. Destination Hub

- Current owner：
  - 目前 repo 內沒有正式 route
  - 現實上被 `step=output` 取代
- Current API basis：
  - 無獨立 hub API
  - 實際可承接的是下游 `dbTargetAPI` 與 `modbusShareAPI`
- Stitch basis：
  - `64f12472cee74a65bd55029eb3ea840e`
  - `[已接受候選] 目標傳輸中心 (Destination Hub)`
- 驗證結論：
  - 這頁不能拿 code route 直接比
  - 目前應視為 `Stitch-planned baseline`
  - 後續用途是承接 output family 導流，不應誤當成已實作頁

### 5. Database Workspace

- Current owner：
  - route：`/studio?step=output&target=database`
  - page：`DatalinkWorkbenchPage -> MuiOutputCommandDeck -> DatabaseTargetBoard`
  - state：`WorkbenchProvider.activeOutputTarget / outputSelectionState`
- Current API basis：
  - `dbTargetAPI`
  - connector / table / validate / schema / dry-run / mapping flows
- Stitch basis：
  - `6d040a45908143248f4bca91e2c7ca57`
  - `[已接受候選] 目標輸出工作台 (Database Workspace)`
- 驗證結論：
  - code 能力已存在
  - 但 route 與 page identity 仍被壓在 output step 裡，尚未成為獨立 workspace

### 6. Local Modbus Workspace

- Current owner：
  - route：`/studio?step=output&target=modbus`
  - page：`DatalinkWorkbenchPage -> MuiOutputCommandDeck -> LocalModbusBoard`
  - state：`WorkbenchProvider.activeOutputTarget / outputSelectionState`
- Current API basis：
  - `modbusShareAPI`
  - status / lifecycle / mapping / sync / write-tag-value
- Stitch basis：
  - `9bca8e1e643c4baf94e7ced109eeef3c`
  - `[已接受候選] 目標輸出工作台 (Local Modbus Workspace)`
- 驗證結論：
  - code 能力已存在
  - 但 route 與 page identity 仍未獨立
  - 後續應把「共享 / 對外提供」語意做得比目前更清楚

### 7. MQTT Workspace

- Current owner：
  - 目前 repo 內沒有正式 route
  - 也沒有獨立 planner shell
- Current API basis：
  - 目前 router 沒有獨立 MQTT destination planner contract
- Stitch basis：
  - `8489c68de1aa4c8da6970e94a01a639d`
  - `[已接受候選] 目標發佈工作台 (MQTT Workspace)`
- 驗證結論：
  - 這頁目前只能視為 `Stitch-planned baseline`
  - 可用來固定未來 page 需求，但暫時不可宣稱已能接回 code

## `/test` Pages

`/test` 的細部基線仍以
`2026-04-10-test-console-current-architecture-baseline.md`
為準，這裡只保留全頁驗證摘要。

### 8. Command Lab

- Current owner：
  - route：`/test`
  - page：`TestPage -> TestOperations`
- Current API basis：
  - `useTestAPI()` 的 `connect / disconnect / status / read / write / batch`
  - `useDebugAPI()` 的 packet / log refresh
- Stitch basis：
  - `72df01403b9a4fe2942e4cdc48b400ef`
- 驗證結論：
  - 功能存在
  - 主要缺口是 shell / session choreography

### 9. Live Monitor

- Current owner：
  - route：`/test`
  - page：`TestPage -> MonitorControl`
- Current API basis：
  - `useTestAPI()` 的 `startMonitor / stopMonitor`
  - `/test/monitor/stream`
- Stitch basis：
  - `0c95dcd8c89a46a691317f200800c322`
- 驗證結論：
  - 功能存在
  - 主要缺口是 watchlist、diagnostics rail 與 session shell

### 10. Device Scan

- Current owner：
  - route：`/test`
  - page：`TestPage -> DeviceScanner`
- Current API basis：
  - 依當前測試連線上下文運作
  - 使用 `/test` 連線能力與 debug context
- Stitch basis：
  - `1c56fb38b066459397286e695adc35c6`
- 驗證結論：
  - 功能存在
  - 可直接作為獨立 workspace 規劃

### 11. RTU Polling Jobs

- Current owner：
  - route：`/test`
  - page：`TestPage -> RTUPollingCard`
- Current API basis：
  - 依 `/test` 連線上下文與 polling job card state 運作
- Stitch basis：
  - `521bb617cef24f8baebcc90a7b5a8a48`
- 驗證結論：
  - 功能存在
  - 主要缺口是 page-level job workspace 與 schedule summary

## 總結

目前全頁驗證後，可以把頁面分成三類：

### A. code 已存在，可直接比對

- Device
- Source
- Tag
- Database
- Local Modbus
- Command Lab
- Live Monitor
- Device Scan
- RTU Polling Jobs

### B. code 已有部分能力，但 page identity 尚未獨立

- Database Workspace
- Local Modbus Workspace
- `/test` 四個 workspace mode

### C. 目前只能用 Stitch 規劃作為基線

- Destination Hub
- MQTT Workspace

## 後續文件使用規則

接下來若要 review 任一頁，應先同時對照：

1. 該頁對應的 current baseline
2. 這份 all-pages verification baseline
3. 該頁的 accepted Stitch screen

這樣之後就不需要每次重新分析 route、state、API 與 page ownership。
