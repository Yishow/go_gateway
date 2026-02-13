# Change: Update Datalink Visual Flow V2

## Why

目前 Datalink UI 雖然已有 Smart Dashboard 與 Wizard，但「來源資料 → Memory Grid → Tag → 寫入資料庫」的主流程仍分散在不同區塊與心智模型中，造成：

1. 操作者難以快速判斷目前資料流卡在哪一段（來源、映射、轉換、寫入）。
2. Memory Grid 選位址後，與 Tag 及寫入目標的關聯缺乏持續可視化上下文。
3. 新手建立第一條可運作資料流需要跨頁/跨面板比對，學習成本偏高。

本次變更目標是把核心操作聚焦為「Flow-first」工作流，以可視化流程為主體，讓配置、驗證、啟用更可預期。

## What Changes

- 新增 Flow-first 工作區規範：同屏呈現 Source、Memory Grid、Tag、Write Target 四段資料流。
- 修改 `datalink-ui` 的引導流程需求，要求每一步都維持可視化流程上下文。
- 新增 `mapping-pipeline` 規範：映射需有明確 flow state（draft/validated/active/error）與逐段可診斷結果。
- 建立可觀測性要求：資料流每段需顯示最近值、品質、時間戳與錯誤位置。
- 定義無障礙與響應式底線：桌面優先但必須在 1024px 以上無水平捲軸、完整鍵盤操作。

## Impact

- Affected specs:
  - `datalink-ui`
  - `mapping-pipeline`
- Affected code (planned):
  - `frontend/src/pages/datalink/SmartDashboard.tsx`
  - `frontend/src/components/datalink/MemoryGrid.tsx`
  - `frontend/src/components/datalink/QuickActions.tsx`
  - `frontend/src/components/datalink/PointDetailPanel.tsx`
  - `frontend/src/components/datalink/MappingCanvas.tsx`
  - `frontend/src/components/datalink/TagTable.tsx`
  - `frontend/src/components/datalink/wizard/*`
  - related query hooks and UI state composition

## Dependencies / Coordination

- 需與進行中的 `optimize-point-configuration-ux` 對齊，避免 UI 規範重複或衝突。
- 不改動通訊協定採集核心與儲存引擎 schema；本次以 UI/流程層與可視化契約為主。
