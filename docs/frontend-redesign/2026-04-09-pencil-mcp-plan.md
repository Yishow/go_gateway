# 2026-04-09 Pencil MCP Plan

## 目標

使用 Pencil MCP 重建 `/studio` 與 `/test` 的新 IA 與 screen structure，並保留現有 backend contract。

本文件重點是 page structure 與 layout planning，不直接等同 visual polish spec。

## 建議檔案策略

- 新建一個獨立 `.pen` 檔作 redesign exploration
- 不直接覆寫既有設計資產
- 先做 screen map，再做單頁細化

建議檔名：

- `designs/studio-redesign-2026-04-09.pen`

## 全域設計規則

- `/studio` 與 `/test` 共用同一品牌語言
- `/studio` workflow-oriented
- `/test` tool-oriented
- 深色工業控制台風格
- 高資訊密度，但保留清楚層級

## 先建立的頂層 screens

建議第一批建立 7 個 frame：

1. `Studio Shell`
2. `Destination Hub`
3. `Database Workspace`
4. `Share Hub`
5. `Local Modbus Workspace`
6. `MQTT Workspace`
7. `Test Debug Console`

## Screen 1: Studio Shell

### 主要區塊

- top context bar
- 4-step rail
- primary work area
- inspector panel
- bottom summary bar

### Pencil 結構建議

- 根 frame: desktop canvas frame
- child A: context bar frame
- child B: step rail frame
- child C: main content frame
- child D: inspector frame
- child E: bottom summary frame

## Screen 2: Destination Hub

### 主要區塊

- step title / summary
- current tag readiness summary
- destination recommendation notice
- destination cards

### 主要卡片

- `Database`
- `Share / Publish`

每張卡片包含：

- short explanation
- recommended badge optional
- next action CTA
- current readiness summary

## Screen 3: Database Workspace

### 結構

- left rail: `Delivery Groups`
- center panel: `Row Planner`
- right panel: `Schema / Apply Preview / Validation`

### 需要的元件類型

- group list cards
- row template toolbar
- column chips / editable rows
- schema summary card
- apply preview card
- validation issue list

## Screen 4: Share Hub

### 結構

- page intro
- two capability cards

### 卡片

- `Local Modbus`
- `MQTT`

### 不要做的內容

- 不做巨大 KPI dashboard
- 不做跨協定複雜總覽
- 不做所有配置直接堆在 hub 上

## Screen 5: Local Modbus Workspace

### 結構

- left rail: `Delivery Groups`
- center panel: `Register Block Planner`
- right panel: `Register Map Preview / Conflict / Apply`

### 視覺重點

- block 作為主操作單位
- register heatmap / occupancy strip
- contiguous preview clearly visible

### 建議元件

- group list
- base register input
- auto-pack strategy buttons
- block width summary
- conflict list
- dry run result panel

## Screen 6: MQTT Workspace

### 結構

- left rail: `Delivery Groups`
- center panel: `Payload / Topic Planner`
- right panel: `Topic Preview / Payload Preview / Validation / Apply`

### 建議元件

- group list
- topic pattern editor
- payload mode switch
- metadata toggles
- sample JSON preview
- validation issue panel

## Screen 7: Test Debug Console

### 結構

- top `Session Header`
- left `Setup Rail`
- center `Main Workspace`
- right `Diagnostics Rail`

### 設計要求

- 比 `/studio` 更高密度
- 更多 monospace
- 更多 inline status
- 主角應是 current connection session，不是很多獨立卡片

### Main Workspace 建議模式

- `Command Lab`
- `Live Monitor`
- `Scan`
- `RTU Polling Jobs`

### 結構意圖

- `Session Header` 持續顯示 profile、protocol、endpoint、connection state
- `Setup Rail` 集中 profile、protocol、mode、config、preset
- `Main Workspace` 一次只突出一種工作節奏
- `Diagnostics Rail` 固定承接 packet、log、detail、export、clear

## Pencil 操作順序建議

### Phase 1

- 建 document-level screen map
- 排好 7 個 screens
- 確認 route-to-screen mapping

### Phase 2

- 先細化 `Studio Shell`
- 再細化 `Destination Hub`
- 再細化 `Database Workspace`

### Phase 3

- 細化 `Share Hub`
- 細化 `Local Modbus Workspace`
- 細化 `MQTT Workspace`

### Phase 4

- 細化 `Test Debug Console`

## 建議的 Pencil 驗證檢查點

- 是否一眼看出 `/studio` 與 `/test` 同品牌、不同語氣
- `Destination` 是否明確是 hub 而不是直接跳頁
- `Delivery Group` 是否成為 destination 主列表單位
- `Database` 是否看起來是 row planner，而不是 tag-column CRUD
- `Local Modbus` 是否看起來是 block planner，而不是逐筆 register 綁定
- `MQTT` 是否看起來是 message planner，而不是單純欄位表單

## 與 backend integration 的設計提醒

Pencil 規畫階段要保留這些 adapter 掛點：

- device capability summary
- source coverage summary
- semantic refinement summary
- delivery group suggestion summary
- database row preview
- modbus block preview
- mqtt message preview
- test console connection state

也就是畫面上要預留「接 view-model adapter」的位置，而不是假設 UI 直接讀 API response 原型。
