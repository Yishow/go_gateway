## MODIFIED Requirements

### Requirement: 三欄式佈局 (REQ-DASH-001)

The Smart Dashboard SHALL use a single-screen Pipeline Studio layout with persistent operator regions, with top controls integrated inside the dashboard container.

#### Scenario: 頁面載入
- **Given** 使用者開啟 SmartDashboard 頁面
- **When** 頁面完成載入
- **Then** 應顯示 dashboard-first 單頁佈局
- **And** 上方顯示整合控制列（分頁導覽 + 設備 context）
- **And** 主區顯示記憶體格工作區與來源規劃
- **And** 右側顯示 Tag 編輯與 DB 提交面板

#### Scenario: 響應式收合
- **Given** 視窗寬度小於 1280px
- **When** 頁面調整佈局
- **Then** 整合控制列可由水平改為上下堆疊
- **And** 三區核心功能仍可在單頁完成
- **And** 不得出現影響主流程的水平捲軸

### Requirement: 設備樹狀導覽 (REQ-DASH-002)

The UI SHALL provide device hierarchy navigation via on-demand drawer and SHALL support expand/collapse of child items.

#### Scenario: 設備列表顯示
- **Given** 系統有 3 個設備
- **When** 使用者開啟切換設備 drawer
- **Then** 應顯示 3 個設備節點
- **And** 每個節點顯示設備名稱與圖標
- **And** 啟用設備顯示狀態指示

#### Scenario: 展開設備子項目
- **Given** 設備 "PLC-001" 有 5 個點位和 3 個映射
- **When** 使用者在 drawer 中點擊展開圖標
- **Then** 應顯示 "Points (5)" 子節點
- **And** 應顯示 "Mappings (3)" 子節點

#### Scenario: 選中設備
- **Given** 使用者在 drawer 內點擊設備節點
- **When** 選中事件觸發並確認切換
- **Then** context bar 應更新設備資訊
- **And** 中間區域應載入該設備的記憶體格子

#### Scenario: 非切換時不常駐設備欄
- **Given** 使用者已完成設備選擇
- **When** 使用者執行日常規劃與提交操作
- **Then** 設備樹不應以常駐左欄佔用主版面
- **And** 僅在使用者主動切換設備時展開

## ADDED Requirements

### Requirement: Top control alignment and adaptive behavior
The dashboard SHALL align search and top control groups on a consistent grid at desktop breakpoints.

#### Scenario: Desktop alignment consistency
- **WHEN** the dashboard is rendered at desktop width
- **THEN** search, navigation tabs, and device context controls align to the same grid rhythm
- **AND** visual jumps between top control rows are minimized

### Requirement: Empty-state guided entry
The dashboard SHALL provide a guided empty state when no device is selected.

#### Scenario: No selected device guidance
- **WHEN** the operator opens dashboard without selected device
- **THEN** the workspace shows empty-state guidance
- **AND** provides only two primary actions: select existing device or create new device

### Requirement: Quick actions prioritize by device state
The dashboard SHALL dynamically prioritize quick actions based on selected device state.

#### Scenario: Offline prioritization
- **WHEN** selected device is offline
- **THEN** reconnect/test actions are prioritized above commit actions
- **AND** unavailable actions are clearly disabled with reasons
