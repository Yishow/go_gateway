# 2026-04-09 Test Console Analysis and IA

## 目的

分析目前 `/test` 頁面實際在做什麼，並定義下一版更合理的資訊架構。

## 我目前看到的頁面結構

目前 `/test` 是一個大型單頁 console。

主要入口在 [TestPage.tsx](/Users/yishow/prj/go_gateway/frontend/src/pages/TestPage.tsx)。

它把以下能力都放在同一頁：

- Profile 管理
- protocol family 與 connection mode 選擇
- connection config 與 connect / disconnect
- 單次 read / write
- batch operation
- polling
- monitor
- device scan
- debug packets / logs
- RTU polling

實際畫面結構大致是：

- 上方：頁首、ProfileSelector、ProtocolSelector
- 左側主區：Config、Operations、Scanner、Monitor、RTU Polling
- 右側：DebugPanel
- 底部浮動：MinimizedCardsBar

## 目前頁面優點

### 1. 工具能力很完整

對現場工程師來說，常見動作大多已經涵蓋。

### 2. DebugPanel 方向正確

[DebugPanel.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/DebugPanel.tsx) 已經很接近真正的工程除錯面板：

- packets / logs 分頁
- filter
- search
- display mode
- export
- clear

### 3. Profile 機制有保留價值

[ProfileSelector.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/ProfileSelector.tsx) 與 [useProfiles.ts](/Users/yishow/prj/go_gateway/frontend/src/hooks/useProfiles.ts) 很適合工廠現場快速切不同設備或站點設定。

## 目前頁面主要問題

### 1. 主角不是 session，而是卡片

現在的心智模型偏向：

`很多工具卡片放在同一頁`

而不是：

`我現在正在測試哪一條連線 session`

這會導致使用者知道工具很多，但不容易持續掌握：

- 現在是哪個 profile
- 現在是哪個 protocol / mode
- 當前連線是否健康
- 正在監控什麼
- 目前封包 / 錯誤是否持續出現

### 2. 操作節奏混在一起

目前以下幾種行為被放在同一層：

- 建立連線
- 執行命令
- 長時間監看
- 掃描設備
- 協定除錯

實際上這些是不同工作節奏，不該全部平鋪同等呈現。

### 3. polling / monitor / RTU polling 邏輯重疊

目前有三類看起來都像持續讀值：

- `TestOperations` 裡的 polling
- `MonitorControl` 裡的 monitor
- `RTUPollingCard` 的 RTU polling

這對工程師不夠直覺，因為三者都像「持續觀測」，但實際用途不同。

### 4. 最小化 card 心智太強

[MinimizedCardsBar.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/MinimizedCardsBar.tsx) 讓頁面更像自訂 dashboard，而不是現場 debug console。

這不是不能用，但它現在已經壓過了真正該突出的 session 狀態。

### 5. 後端能力尚未形成完整前端工作台

[router.go](/Users/yishow/prj/go_gateway/internal/api/router.go) 裡 `/test` 與 `/debug` 還有這些能力：

- `/test/script`
- `/test/scripts`
- `/debug/send-raw`
- `/debug/analyze/:packetId`
- `/connections`
- `/config/presets`
- `/templates`

但目前主頁面主要只吃到基礎 connect/read/write/batch/monitor/debug refresh。

也就是說，頁面已經很大，但還不是完整能力中心。

## 我建議的新 IA

### 核心心智模型

把 `/test` 改成：

`Connection Session-centered Field Engineer Debug Console`

不是很多卡片的集合，而是一條當前 session 的控制台。

### 建議保留單一路由

我建議先保留：

- `/test`

不要第一輪就拆很多 route。

原因：

- 現場工程師常需要快速切換不同工作區
- 單頁更符合 debug console 節奏
- 可先用內部分頁 / segmented workspace 解決複雜度

## 建議頁面骨架

### 1. Session Header

固定顯示：

- current profile
- protocol
- mode
- endpoint
- connection state
- connection id
- monitor state
- latest packet / error summary
- quick connect / disconnect

### 2. Left Setup Rail

保留但收斂成 setup 區：

- profile manager
- protocol selector
- connection config
- connection presets

### 3. Main Workspace

主工作區採 segmented mode，不同工作節奏不要同時攤平：

- `Command Lab`
- `Live Monitor`
- `Scan`
- `RTU Polling`

### 4. Right Diagnostics Rail

固定保留 diagnostics：

- packet stream
- log stream
- selected packet detail
- raw send
- analyze result

## 各 workspace 的角色

### Command Lab

用於臨時命令驗證。

包含：

- single read / write
- batch operation
- short polling
- result history

### Live Monitor

用於長時間觀測。

包含：

- monitor item config
- chart
- current values
- monitor log

### Scan

用於找設備，不應再和命令執行混成同一塊。

包含：

- station / IP scan setup
- scan progress
- statistics
- result list

### RTU Polling

保留為專用進階工具。

包含：

- polling job list
- per-job operation setup
- enable / disable
- run status

## 建議的命名重整

為了避免三個「持續讀值」混淆，我建議改名：

- `Polling` -> `Command Polling`
- `Monitor` -> `Live Monitor`
- `RTU Polling` -> `RTU Polling Jobs`

## 設計語氣建議

`/test` 應與 `/studio` 同品牌，但語氣更像：

- engineer console
- dense tool surface
- fast diagnosis
- operational detail first

不應像：

- 精簡版 `/studio`
- onboarding flow
- generic card dashboard

## 對現有元件的保留建議

### 建議保留並升級

- [ProfileSelector.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/ProfileSelector.tsx)
- [ProtocolSelector.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/ProtocolSelector.tsx)
- [ConfigForm.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/ConfigForm.tsx)
- [DebugPanel.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/DebugPanel.tsx)
- [DeviceScanner.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/DeviceScanner.tsx)

### 建議重組而不是直接保留原樣

- [TestOperations.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/TestOperations.tsx)
- [MonitorControl.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/MonitorControl.tsx)
- [RTUPollingCard.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/RTUPollingCard.tsx)

### 建議降級為次要輔助機制

- [MinimizedCardsBar.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/MinimizedCardsBar.tsx)
- [CardMinimizeButton.tsx](/Users/yishow/prj/go_gateway/frontend/src/components/CardMinimizeButton.tsx)

## 第一輪 redesign 的重點

第一輪不是補更多功能，而是先做這三件事：

1. 把主角改成 current session
2. 把主工作區改成分明的四類 workspace
3. 把 diagnostics 升成固定右側診斷柱

## 成功標準

下一版 `/test` 成功的判準不是更漂亮，而是：

- 使用者一眼知道現在在測哪條 session
- 命令、監控、掃描、除錯不再互相搶主畫面
- packet / log / error 對當前 session 的關聯更清楚
- UI 結構能自然接回現有 `/api/v1/test/*` 與 `/api/v1/debug/*`
