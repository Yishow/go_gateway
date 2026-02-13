## Why

目前新版 Dashboard 與既有側邊欄選項對應不完整，造成操作入口分裂：
- 使用者仍會從舊 sidebar 找功能，但新版主要流程已移入 dashboard。
- `Server Memory Grid` 獨立頁入口不夠明確，導流不足。
- Dashboard 中央 Grid 高度過高，壓縮上下文資訊。
- 型別占格實際行為與預期不一致（`float32`、`int64`）。

此外，需求已明確：除 `/test` 外，其餘 sidebar 功能全改為 modal。

## What Changes

1. 將原 sidebar 功能入口全部改為 dashboard 內 modal（僅 `/test` 保持獨立頁）。
2. 補齊 `Server Memory Grid` 導流路徑（context bar、Quick Actions、Grid 卡片三處一致入口）。
3. 下修 Dashboard Memory Grid 視覺高度，提升同屏資訊密度。
4. 修正並強制驗證 typed occupancy 規則：
   - `float32` = 2 連續格
   - `int64` = 4 連續格
5. 以 OpenSpec 任務化 20 項 UX 改進建議，並排除指定項目（原提案 10/15/16 不納入）。

## Scope Constraints

- 不做原提案第 10、15、16 項。
- `/test` 保留獨立頁，不轉 modal。
- 本變更優先修正入口一致性、可達性與占格正確性，不擴大到新資料模型。

## Capabilities

### Modified Capabilities

- `datalink-ui`
- `datalink-smart-dashboard`
- `local-modbus-memory-workbench`

## Impact

- Frontend routing and modal orchestration
- SmartDashboard layout and navigation actions
- MemoryGrid typed occupancy rendering/validation
- Local Modbus workbench entry discoverability
