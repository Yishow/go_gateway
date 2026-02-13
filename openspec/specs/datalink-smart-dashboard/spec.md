# datalink-smart-dashboard Specification

## Purpose
TBD - created by archiving change optimize-point-configuration-ux. Update Purpose after archive.
## Requirements
### Requirement: 三欄式佈局 (REQ-DASH-001)

The Smart Dashboard SHALL use a single-screen Pipeline Studio layout with three persistent regions: source planning rail (left), memory grid workspace (center), and tag/commit rail (right).

#### Scenario: 頁面載入
- **Given** 使用者開啟 SmartDashboard 頁面
- **When** 頁面完成載入
- **Then** 應顯示 Pipeline Studio 三區佈局
- **And** 左側顯示來源規劃與模板
- **And** 中間顯示記憶體格與占用狀態
- **And** 右側顯示 Tag 編輯與 DB 提交面板

#### Scenario: 響應式收合
- **Given** 視窗寬度小於 1280px
- **When** 頁面調整佈局
- **Then** 三區功能仍可在單頁完成
- **And** 不得出現影響主流程的水平捲軸

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

The center area SHALL display a memory grid that visualizes typed occupancy based on selected source data type and count.

Typed span policy:
- bool/int16/uint16: 1 cell per source
- int32/uint32/float32: 2 adjacent cells per source
- int64/uint64/float64: 4 adjacent cells per source

#### Scenario: int 類型占格
- **Given** 使用者選擇 int16 類型且來源數量為 5
- **When** 產生配置預覽
- **Then** 記憶體格應標示 5 格規劃占用

#### Scenario: float 類型占格
- **Given** 使用者選擇 float32 類型且來源數量為 10
- **When** 產生配置預覽
- **Then** 記憶體格應標示 20 格占用
- **And** 視覺上顯示為 10 組雙格配對

#### Scenario: 衝突占用提示
- **Given** 規劃占用與既有點位或已綁定區段衝突
- **When** 使用者嘗試提交
- **Then** 衝突格應以明確狀態標示
- **And** 提供修正建議與跳轉定位

#### Scenario: 自動配置連續區段
- **Given** 使用者設定來源類型與數量
- **When** 使用者啟用自動配置
- **Then** 系統應尋找最近可用連續區段進行配置
- **And** 若無可用區段則回報不可配置原因

#### Scenario: 多格型別群組視覺
- **Given** 使用者配置 float32 或 int32 類型來源
- **When** 記憶體格渲染
- **Then** 每個來源的雙格群組應以明確群組邊框顯示
- **And** 群組內應標示順序索引

#### Scenario: 只看衝突格
- **Given** 畫面存在多個衝突占格
- **When** 使用者切換「只看衝突格」過濾器
- **Then** 記憶體格應僅顯示衝突項目與必要上下文

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

