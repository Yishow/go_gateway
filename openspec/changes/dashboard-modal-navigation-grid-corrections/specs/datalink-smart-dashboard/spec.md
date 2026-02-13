## MODIFIED Requirements

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
