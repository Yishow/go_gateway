## MODIFIED Requirements
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

## ADDED Requirements
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
