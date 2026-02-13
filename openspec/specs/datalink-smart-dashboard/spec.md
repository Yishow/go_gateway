# datalink-smart-dashboard Specification

## Purpose
TBD - created by archiving change optimize-point-configuration-ux. Update Purpose after archive.
## Requirements
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

### Requirement: 記憶體格子視覺化 (REQ-DASH-003)

The center area SHALL display a memory grid with compact height constraints and correct typed occupancy spans.

Typed span policy (authoritative):
- bool/int16/uint16: 1 cell per source
- int32/uint32/float32: 2 adjacent cells per source
- int64/uint64/float64: 4 adjacent cells per source

#### Scenario: Compact grid height in dashboard
- **Given** 使用者位於 `/datalink` dashboard
- **When** 中央 Memory Grid 渲染
- **Then** Grid 區塊高度應低於現行版本基線
- **And** 上下文與右側操作區資訊在常見桌面解析度可同屏看到更多內容

#### Scenario: float32 span correctness
- **Given** 使用者設定 `float32` 且來源數量為 10
- **When** 系統產生規劃占格
- **Then** 應產生 20 格占用
- **And** 每個來源以 2 格連續群組顯示

#### Scenario: int64 span correctness
- **Given** 使用者設定 `int64` 且來源數量為 5
- **When** 系統產生規劃占格
- **Then** 應產生 20 格占用
- **And** 每個來源以 4 格連續群組顯示

### Requirement: 快速操作面板 (REQ-DASH-004)

The right panel SHALL provide commit-oriented actions for validate, commit, and recovery, with pending allocation summary.

#### Scenario: 提交前摘要
- **Given** 使用者完成來源、格位、Tag 指派
- **When** 右側面板更新
- **Then** 應顯示待提交數量、衝突數、受影響全域 Tag 數量
- **And** 應顯示預估輪詢負載變化量

#### Scenario: 提交成功回饋
- **Given** 驗證與提交成功
- **When** 系統回傳成功
- **Then** 右側面板顯示成功狀態
- **And** 記憶體格即時更新為已連結狀態

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

### Requirement: 動畫引導流程 (REQ-DASH-007)

The dashboard SHALL provide motion cues to guide next actions during key transitions: plan confirmed, grid allocated, tag linked, and DB committed.

#### Scenario: 來源規劃完成引導
- **Given** 使用者完成來源規劃
- **When** 按下下一步或確認
- **Then** 系統以短動畫引導焦點到記憶體格區域

#### Scenario: 提交完成引導
- **Given** 使用者完成 DB 提交
- **When** 提交結果為成功
- **Then** 系統以短動畫標示完成節點
- **And** 在可存取模式下提供等效靜態提示

### Requirement: 提交安全與可追蹤性 (REQ-DASH-008)

The dashboard SHALL support snapshot-based restore, commit traceability, and two-stage validation before DB commit.

#### Scenario: 兩階段驗證
- **Given** 使用者準備提交配置
- **When** 執行驗證
- **Then** 系統先執行結構驗證
- **And** 結構驗證通過後再執行可執行驗證

#### Scenario: 顯示批次追蹤 ID
- **Given** 使用者提交配置
- **When** 提交流程開始
- **Then** 系統產生批次追蹤 ID
- **And** 面板可檢視該 ID 對應的提交結果

#### Scenario: 快照還原
- **Given** 使用者建立多次規劃修改
- **When** 使用者選擇還原到先前快照
- **Then** 系統應恢復對應的來源、格位、Tag 規劃狀態

### Requirement: Tag 全域編輯差異預覽 (REQ-DASH-009)

The dashboard SHALL present field-level before/after diff before saving global tag edits in linkage context.

#### Scenario: 顯示差異後確認儲存
- **Given** 使用者在連結流程中編輯全域 Tag
- **When** 使用者按下儲存
- **Then** 系統應先顯示欄位差異預覽
- **And** 使用者確認後才真正寫入全域 Tag

### Requirement: 動畫可讀性檢核 (REQ-DASH-010)

The dashboard SHALL enforce animation readability checks to ensure motion does not hide or delay critical decision information.

#### Scenario: 可讀性檢核不通過
- **Given** 某段動畫遮蔽關鍵資訊或超過規範時長
- **When** UI 驗證執行
- **Then** 該動畫應被標記為不合格
- **And** 需調整後方可通過交付檢核

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

