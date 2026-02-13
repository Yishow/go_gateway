# datalink-smart-dashboard Specification

## Purpose
TBD - created by archiving change optimize-point-configuration-ux. Update Purpose after archive.
## Requirements
### Requirement: 三欄式佈局 (REQ-DASH-001)

The Smart Dashboard SHALL use a three-column layout: device tree (left), dynamic content (center), and quick actions panel (right).

#### Scenario: 頁面載入

- **Given** 使用者開啟 SmartDashboard 頁面
- **When** 頁面完成載入
- **Then** 應顯示三欄佈局
- **And** 左側顯示設備樹狀導覽（240px）
- **And** 中間顯示設備選擇提示或記憶體格子
- **And** 右側顯示快速操作面板（280px）

#### Scenario: 響應式收合

- **Given** 視窗寬度小於 1280px
- **When** 頁面調整佈局
- **Then** 左側導覽應自動收合至圖標模式（64px）
- **And** 右側面板應改為浮動按鈕

---

### Requirement: 設備樹狀導覽 (REQ-DASH-002)

The UI SHALL provide a tree navigation that displays device hierarchy and SHALL support expand/collapse of child items.

#### Scenario: 設備列表顯示

- **Given** 系統有 3 個設備
- **When** 樹狀導覽載入
- **Then** 應顯示 3 個設備節點
- **And** 每個節點顯示設備名稱與圖標
- **And** 啟用設備顯示綠色指示點

#### Scenario: 展開設備子項目

- **Given** 設備 "PLC-001" 有 5 個點位和 3 個映射
- **When** 使用者點擊展開圖標
- **Then** 應顯示 "Points (5)" 子節點
- **And** 應顯示 "Mappings (3)" 子節點

#### Scenario: 選中設備

- **Given** 使用者點擊設備節點
- **When** 選中事件觸發
- **Then** 該節點應顯示選中樣式
- **And** 中間區域應載入該設備的記憶體格子

#### Scenario: 拖曳排序設備

- **Given** 使用者拖曳設備節點 "PLC-002"
- **When** 放置到 "PLC-001" 上方
- **Then** 設備順序應更新為 PLC-002, PLC-001, PLC-003
- **And** 應呼叫 `PATCH /api/v1/datalink/devices/reorder`
- **And** 順序應持久化到資料庫

#### Scenario: 拖曳視覺反饋

- **Given** 使用者開始拖曳設備節點
- **When** 拖曳進行中
- **Then** 被拖曳節點應顯示半透明效果
- **And** 目標位置應顯示藍色插入線
- **And** 拖曳手柄圖標應可見

---

### Requirement: 記憶體格子視覺化 (REQ-DASH-003)

The center area SHALL display a memory grid that visualizes PLC register usage.

#### Scenario: 格子渲染

- **Given** 設備協議為 MC3E，中心位址為 D100
- **When** 記憶體格子載入
- **Then** 應顯示 100 個格子（D050-D149）
- **And** 已使用格子顯示綠色
- **And** 可用格子顯示灰色
- **And** 每個格子顯示位址標籤

#### Scenario: 格子 Tooltip

- **Given** 格子 D100 已綁定點位 "Temp_Sensor"
- **When** 使用者懸停該格子
- **Then** 應顯示 Tooltip
- **And** Tooltip 包含點位名稱、最後讀取值、時間

#### Scenario: 單擊選取

- **Given** 使用者點擊可用格子 D105
- **When** 點擊事件觸發
- **Then** 該格子應顯示選中樣式
- **And** 側邊面板應開啟「新增點位」表單

#### Scenario: 範圍選取

- **Given** 使用者已選取 D100
- **When** 使用者 Shift+點擊 D109
- **Then** D100-D109 共 10 個格子應全部選中
- **And** 快速操作面板應顯示「選取 10 個位址」

---

### Requirement: 快速操作面板 (REQ-DASH-004)

The right panel SHALL provide quick action buttons and real-time status display.

#### Scenario: 批量建立按鈕

- **Given** 使用者選取 5 個格子
- **When** 快速操作面板顯示
- **Then** 「批量建立」按鈕應啟用
- **And** 按鈕應顯示「批量建立 (5)」

#### Scenario: 即時狀態顯示

- **Given** 設備 PLC-001 已連線
- **When** 快速操作面板載入
- **Then** 應顯示連線狀態指示燈（綠色）
- **And** 應顯示最後同步時間

---

### Requirement: 側邊滑出面板 (REQ-DASH-005)

The UI SHALL provide a slide-out panel from the right when the user selects a cell or action button.

#### Scenario: 開啟面板

- **Given** 使用者點擊「批量建立」按鈕
- **When** 面板開啟
- **Then** 面板應從右側滑入（300ms）
- **And** 背景應顯示半透明遮罩

#### Scenario: 關閉面板

- **Given** 面板已開啟
- **When** 使用者按 ESC 鍵
- **Then** 面板應滑出關閉
- **And** 遮罩應淡出

---

### Requirement: 批量點位建立 (REQ-DASH-006)

The batch creator SHALL support creating multiple points for consecutive addresses in one operation.

#### Scenario: 命名模板

- **Given** 起始位址 D100，數量 5，模板 "Sensor\_{index}"
- **When** 預覽生成
- **Then** 應顯示：Sensor_0, Sensor_1, Sensor_2, Sensor_3, Sensor_4
- **And** 對應位址：D100, D101, D102, D103, D104

#### Scenario: 批量建立成功

- **Given** 使用者填寫完整表單
- **When** 點擊「建立」按鈕
- **Then** 應呼叫 `POST /api/v1/datalink/points/batch`
- **And** 成功後顯示成功訊息
- **And** 記憶體格子應更新顯示新建立的點位

---

