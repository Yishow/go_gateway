## Context

使用者要求將舊 sidebar 與新版 dashboard 的資訊架構統一：
- 除 `/test` 外，舊 sidebar 對應功能全部改為 modal。
- `Server Memory Grid`（`/datalink/local-modbus`）需有清楚可見的固定入口。
- Dashboard Memory Grid 需降低高度以提升操作區與狀態區可視範圍。
- typed occupancy 必須正確且可驗證：`float32=2`、`int64=4`。

## Goals / Non-Goals

**Goals**
- 將 sidebar 功能遷移為 modal-first 導覽模型（test 例外）。
- 建立 local-modbus 的多點一致入口，避免深層路由遺失。
- 重新平衡 dashboard 版面高度與資訊密度。
- 以規格與測試鎖定 typed occupancy 規則，避免回歸。
- 交付 20 項可實作 UX 改進建議。

**Non-Goals**
- 不重做 `/test` 頁面資訊架構。
- 不導入新設計系統。
- 不納入原建議第 10、15、16 項（依使用者指示排除）。

## Decisions

### D1. Sidebar 全部 modal 化（test 例外）
- 內容：`devices/settings/points/mappings/wizard` 統一以 dashboard modal 承載。
- 例外：`/test` 保持獨立頁。

### D2. Local-modbus 三點導流
- 內容：在 context bar、Quick Actions、Dashboard Grid 卡片都提供 `Server Memory Grid` 入口。
- 規則：入口文案一致，皆導向 `/datalink/local-modbus`。

### D3. Grid 高度下修
- 內容：Dashboard 主 Grid 區塊改為較低高度（以可用空間 + 斷點策略限制最大高度）。
- 目標：提升上方 Context 與右側 Queue 同屏可視比例。

### D4. Typed occupancy 強制規格化
- 內容：明確驗收 `float32` 兩格連續、`int64` 四格連續，且 UI 顯示群組連續性。
- 驗證：單元測試 + 互動測試 + 視覺檢查。

## 20 項建議（落地版）

1. `Device Center Modal`：整合選擇/搜尋/切換/建立。
2. `Settings Modal`：分區 tabs（通訊、寫入、安全、通知）。
3. `Point Manager Modal`：批次建立、匯入匯出。
4. `Mapping Modal`：來源→Tag→Sink 的單窗編輯。
5. `Flow Wizard Modal`：新手流程入口。
6. modal 入口統一由 context bar 觸發。
7. dashboard 顶部加入「管理入口」按鈕群。
8. Quick Actions 固定顯示 `Server Memory Grid` CTA。
9. Grid 卡片右上固定顯示 `完整工作台` 連結。
10. (排除) conflict 深連結 query 導航，不納入。
11. local-modbus 支援回跳 dashboard 並保留 section。
12. local-modbus 顯示目前來源設備 context 摘要。
13. dashboard Grid 卡摘要顯示 bind_state/mapping/conflict。
14. dashboard Grid 卡提供 start/stop 迷你控制（完整設定仍在獨立頁）。
15. (排除) 最近寫入事件列表，不納入。
16. (排除) Top3 衝突 register 展示，不納入。
17. Preflight Gate 視覺統一（PASS/FAIL + 固定色階）。
18. 異常訊息提供 `重試 + 查看詳情 + 前往工作台`。
19. 手機維持側滑 panel modal，底部固定操作列。
20. 舊 sidebar 路由一律 redirect 到 `/datalink` 並打開對應 modal。

## Risks / Trade-offs

- modal 數量增加可能使狀態管理複雜化。
- grid 高度下降可能增加捲動需求。
- 路由兼容邏輯增加，需避免重定向迴圈。

## Mitigation

- 將 modal state 集中管理（query + local state 協作）。
- 響應式斷點下採不同高度上限。
- 對 legacy route 做明確白名單映射與測試。

## Rollout

1. 路由/入口一致化。
2. modal 化與 dashboard 導流補齊。
3. Grid 高度調整。
4. typed occupancy 修正與測試。
5. 驗證與回歸測試。
