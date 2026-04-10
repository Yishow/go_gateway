# 2026-04-10 Test Console Current Architecture Baseline

## 目的

這份文件是 `/test` 現況基準。

用途不是直接定義新 UI，而是固定記錄：

- 目前 repo 內真正存在的 `/test` 能力
- 現有 route / page / component / service / API 對接關係
- 現有 state ownership
- accepted Stitch screens 與現況程式的落差

之後做 `/test` redesign、accepted screen 對接、adapter 規劃時，應先對照這份文件，而不是每次重新做一次探索。

## 範圍

本文件只描述目前 repo 內 `/test` 現況。

不處理：

- `/studio` 主產品流程
- `/test` 的最終重構實作
- backend contract 變更提案

## Route 與入口

### 前端 route

- `/test`
  - 入口在 `frontend/src/App.tsx`
  - render：`<TestPageShell><TestPage /></TestPageShell>`
- legacy 測試工具 route 會 redirect 到 `/test`
  - `/templates`
  - `/history`
  - `/compare`
  - `/analyzer`

### 頁面入口

- `frontend/src/pages/TestPage.tsx`
  - `/test` 真正頁面
- `frontend/src/pages/TestPageShell.tsx`
  - 目前只提供背景與 padding
  - 尚未形成 session shell / setup rail / diagnostics rail

## 後端 API 現況

### `/api/v1/test/*`

定義於 `internal/api/router.go`：

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

### 其他周邊 API

- `/api/v1/config/*`
  - presets
- `/api/v1/templates/*`
  - 測試模板
- `/api/v1/connections/*`
  - 連線池資訊

## 前端 API 封裝現況

### `frontend/src/services/api.ts`

目前 `/test` 主要透過 `useTestAPI()` 與 `useDebugAPI()` 呼叫：

- `connect`
- `disconnect`
- `getStatus`
- `read`
- `write`
- `batch`
- `startMonitor`
- `stopMonitor`
- `refresh debug packets/logs`
- `clear debug data`

現況問題：

- `script`、`scripts`、`send-raw`、`analyze`、`connections`、`presets`、`templates` 尚未在主頁形成完整工作流
- API 有，但 `/test` 主頁沒有把它們提升成正式 workspace 能力

## 頁面結構現況

### 頂部區

`TestPage.tsx` 頂部目前包含：

- 頁首標題
- `ProfileSelector`
- `ProtocolSelector`
- 縮小版 `Config` 摘要

### 主區塊

主區分成兩欄：

- 左側主欄
  - `ConfigForm`
  - `TestOperations`
  - `DeviceScanner`
  - `MonitorControl`
  - `RTUPollingCard`
- 右側欄
  - `DebugPanel`

### 補助區

- `MinimizedCardsBar`
  - 管理被縮小的卡片

## 功能模組盤點

### 1. Profile 管理

相關檔案：

- `frontend/src/components/ProfileSelector.tsx`
- `frontend/src/hooks/useProfiles.ts`

現有能力：

- create profile
- rename profile
- delete profile
- switch profile
- localStorage persist
- 依 `connectionMode` 分開保存 config

價值：

- 很適合現場工程師在不同設備、站點、通訊模式間快速切換
- 是 `/test` 新 shell 應保留的高價值能力

### 2. Protocol 與連線配置

相關檔案：

- `frontend/src/components/ProtocolSelector.tsx`
- `frontend/src/components/ConfigForm.tsx`
- `frontend/src/pages/TestPage.tsx`

現有能力：

- protocol family 切換
- connection mode 切換
- tcp / udp / serial config
- connect / disconnect
- 連線成功後自動縮小 config card

現況問題：

- setup 區雖然存在，但仍被實作成其中一張卡
- 尚未成為 accepted `/test` 所需的固定 `Setup Rail`

### 3. Command Lab 類能力

相關檔案：

- `frontend/src/components/TestOperations.tsx`

現有能力：

- single read / write
- batch queue
- short polling
- read / write result 顯示
- protocol-specific operation 切換

現況判定：

- 能力本身夠做 accepted `Command Lab`
- 但目前與 monitor / scan / RTU polling 同層平鋪，工作節奏不夠清楚

### 4. Live Monitor 類能力

相關檔案：

- `frontend/src/components/MonitorControl.tsx`

現有能力：

- monitor item 配置
- interval 設定
- SSE stream 接收
- chart 與最新資料
- reconnect logic
- monitor config 存回 profile

現況判定：

- backend contract 足夠
- 前端能力也足夠形成 accepted `Live Monitor`
- 真正缺的是 session shell 與 diagnostics 聯動

### 5. Device Scan 類能力

相關檔案：

- `frontend/src/components/DeviceScanner.tsx`

現有能力：

- Modbus / Fatek station scan
- MC Protocol IP range scan
- protocol-specific operation / symbol 設定
- statistics / progress / result filtering
- 可用 base config 建立臨時掃描連線

現況判定：

- 功能實際上比表面看起來多
- 已足夠支撐 accepted `Device Scan`
- 但目前只是大頁中的一張卡，還不是獨立 workspace

### 6. RTU Polling Jobs 類能力

相關檔案：

- `frontend/src/components/RTUPollingCard.tsx`
- `frontend/src/components/PollingOperationCard.tsx`

現有能力：

- 多個 polling operation
- protocol-specific operation options
- enable / disable
- per-operation setup

現況問題：

- 視覺語言仍是舊卡片風格
- 與 accepted `/test` 家族不一致
- state 仍內聚在 card 內，尚未提升為 page workspace 層

### 7. Diagnostics 類能力

相關檔案：

- `frontend/src/components/DebugPanel.tsx`
- `frontend/src/services/api.ts`

現有能力：

- packets / logs 分頁
- auto refresh
- direction filter
- search
- selected packet
- clear
- hex / ascii / parsed display mode

現況判定：

- 方向正確
- 這塊很適合直接承接 accepted `Diagnostics Rail`
- 但現在仍是右側一張 debug card，不是整個 shell 的固定診斷柱

## State Ownership 現況

### Page-level state

在 `TestPage.tsx`：

- `selectedProtocol`
- `connectionMode`
- `config`
- `connectionId`
- `isConfigMinimized`

這些 state 代表「目前 session 的上層上下文」。

### Component-level state

各工具卡仍大量自管 state：

- `TestOperations`
  - 單次命令、batch queue、polling records
- `MonitorControl`
  - monitor items、interval、monitorData、SSE lifecycle
- `DeviceScanner`
  - scan config、scan results、statistics、filter
- `RTUPollingCard`
  - operations list
- `DebugPanel`
  - active tab、search、selected packet、display mode

現況結論：

- `/test` 現在不是一個 session-centered shell
- 而是「一個頁層 connection context + 多張自帶內部狀態的工具卡」

這表示後續接 accepted `/test` 時，需要先建立：

- `test-console-shell view-model`
- `test-console-session summary`
- `active workspace adapter`

不能只做 JSX 重排。

## 現況整合鏈

### Route -> Page -> Component -> Service -> API

#### Session / config

- `/test`
- `TestPage`
- `ProfileSelector` / `ProtocolSelector` / `ConfigForm`
- `useProfiles` / `useTestAPI`
- `/test/connect` `/test/disconnect` `/test/status`

#### Command Lab

- `/test`
- `TestPage`
- `TestOperations`
- `useTestAPI`
- `/test/read` `/test/write` `/test/batch`

#### Live Monitor

- `/test`
- `TestPage`
- `MonitorControl`
- `useTestAPI`
- `/test/monitor/start` `/test/monitor/stop` `/test/monitor/stream`

#### Device Scan

- `/test`
- `TestPage`
- `DeviceScanner`
- `useTestAPI`
- 主要仍複用 `/test/connect` `/test/read` `/test/disconnect`

#### RTU Polling

- `/test`
- `TestPage`
- `RTUPollingCard`
- 透過 polling operation card 與既有測試能力運作
- 目前尚未被提升成獨立 shell mode

#### Diagnostics

- `/test`
- `TestPage`
- `DebugPanel`
- `useDebugAPI`
- `/debug/packets` `/debug/logs` `/debug/clear`

## 與 accepted Stitch screens 的對照

### 可直接對應的 accepted screen

- `Command Lab`
  - 現有能力足夠
- `Live Monitor`
  - 現有能力足夠
- `Device Scan`
  - 現有能力足夠
- `RTU Polling Jobs`
  - 現有能力存在，但 UI 和 state ownership 還較弱

### 真正缺口

缺的不是主要 backend contract，而是：

- `/test` shell
- mode switching choreography
- session summary
- diagnostics rail 固定化
- component state 往 page-level adapter 收口

## 作為後續比對依據的使用方式

之後每次討論 `/test` 時，先用這份文件回答三件事：

1. 現有 repo 已有什麼能力
2. accepted screen 想表達的是不是已有能力，只是殼不對
3. 哪些地方真的需要新 adapter 或 backend amendment

如果沒有新的後端 contract 變更，就不要重新從零分析 `/test`。

先更新這份 baseline，再在這份 baseline 上做比對即可。
