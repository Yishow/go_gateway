# Proposal: Workbench UX Operator Efficiency Overhaul

## Why

目前 Datalink Workbench 的 5 步驟流程（Device → Source → Tag → Output）在功能上已具備完整的 E2E 資料採集設定能力，但在**操作員效率**上存在多項系統性缺陷：

**核心痛點：**

1. **重複輸入問題（Step 2）**：資深操作員每次建立 Source Rule 都要重新輸入起始地址、數量、命名前綴，沒有 Template 快速套用入口（Template 功能藏在 toolbar 次級選單）。

2. **元件過度膨脹**：`SourceCanvasSection`（2688行）、`WorkbenchDeviceStep`（1713行）、`TagBindingStudio`（1712行）三個核心元件單一職責嚴重違反，導致維護困難、測試覆蓋不完整。

3. **就緒狀態不可操作**：StepRail 與 BottomSummaryBar 的 `partial`/`blocked` 色點沒有說明「為什麼卡住」，操作員需自行逐步驟排查原因。

4. **Step 4 資訊密度過高**：DB connector 建立、Table schema 設定、Column 映射全部嵌在同一頁面，首次設定體驗極差；Local Modbus 雖然是可選功能，但目前與 Database Output 強制並排。

5. **跨目標映射不透明**：操作員不清楚哪些 Tag 已映射到 Modbus、哪些到 Database，兩個輸出目標的狀態沒有統一視圖。

## Design System

所有新頁面以 **Linear Design System** 為設計語言基礎，適配工業 Workbench 情境：

- **設計規範來源**：`frontend/docs/DESIGN-PRINCIPLES.md`（唯一設計參考文件）
- **設計語言參考**：`frontend/docs/LINEAR-DESIGN-REFERENCE.md`（Linear 原語對照）
- **核心原則**：深色優先、資訊密集、精準工程感
- **Accent 色**：Cyan `#06b6d4`（取代 Linear 的 Indigo，與現有 slate/cyan 架構一致）
- **背景層次**：`#08090a` void → `#0f1011` base → `#191a1b` elevated → `#28282c` surface

## Implementation Strategy

### 新建頁面，保留舊頁面

- **新頁面路徑**：`frontend/src/pages/studio/`（全新設計，基於 DESIGN-PRINCIPLES.md）
- **路由**：`/studio`（新主流程入口）、`/tools/modbus`（Local Modbus 工具頁）
- **舊頁面保留**：`frontend/src/pages/datalink/workbench/` 保持不動，舊路由 `/datalink/workbench` 維持運作
- **平行存在期**：新舊頁面同時可用，切換穩定後再考慮廢棄舊路由

## What Changes

### 流程架構重組

**主流程縮減為 4 步驟**，Local Modbus 獨立為工具頁：

```
主流程: Device → Source → Tag → Database Output
工具頁: /tools/modbus（Local Modbus 完整工具，可選）
```

### 全新頁面設計（基於 Linear 設計系統）

| 新頁面 | 路由 | 對應舊元件（邏輯遷移來源） |
|--------|------|--------------------------|
| `StudioDevicePage` | `/studio?step=device` | `WorkbenchDeviceStep` (1713行) |
| `StudioSourcePage` | `/studio?step=source` | `SourceCanvasSection` (2688行) |
| `StudioTagPage` | `/studio?step=tag` | `TagBindingStudio` (1712行) |
| `StudioDatabasePage` | `/studio?step=output` | `DatabaseTargetBoard` (1191行) |
| `LocalModbusToolPage` | `/tools/modbus` | `LocalModbusBoard` (1084行) |

舊元件**不修改**，新頁面從 `frontend/src/pages/studio/` 全新建立，遵守 `DESIGN-PRINCIPLES.md`。

### 新增 UX 功能

**Step 1 改善：**
- `DeviceTestConsole`：Connect / Probe 分欄卡片，視覺層次清楚；測試記錄保留最近 3 次歷史
- Connect 成功但 Probe 失敗時顯示「資料收集將被阻擋」標語

**Step 2 改善：**
- `RuleTemplateQuickBar`：左側固定展開式模板列，支援懸停預覽、一鍵套用、一鍵儲存
- `Plan/Live/Link` 模式移到畫布頂部明確 Tab 列

**Step 3 改善：**
- `TagBatchActionBar`：一鍵全部 Auto-Bind（含策略選擇）→ diff preview → apply
- Apply 後顯示進度完成橫幅，引導跳轉至 Step 4

**Step 4（Database Output）新增：**
- `DBConnectorManager`：多 connector 管理（同時寫入 SQLite + PostgreSQL）
- `TableSchemaGenerator`：依 Tag 自動建表（一鍵生成 schema）
- `ColumnMappingBoard`：Tag → DB column 映射狀態（已映射/未映射/衝突）
- Dry-run 預覽（真實寫入前可預覽結果）
- 寫入模式設定（Upsert / Insert-only / Overwrite）
- 歷史記錄查看（最近寫入的數據 / 錯誤日誌）

**新增：WorkbenchDiagnosticPanel（Slide-over）：**
- 觸發：點擊 `partial`/`blocked` 色點，或頂欄常駐「⚠ N 個問題」按鈕
- 內容：按步驟分組的問題清單，每項附一鍵跳轉+主畫面高亮
- Slide-over 右側推出，不遮擋主工作區，可邊診斷邊操作

**新增：Local Modbus 工具頁（/tools/modbus）：**
- Register 地址總覽表（全部 HR 的已映射/空閒/衝突狀態）
- AutoMap + Dry-run 預覽
- 衝突自動偵測與一鍵解決（提供候選解決方案）
- Live 實時讀值顯示（SSE）
- 分組管理（不同 Tag 組分配到不同地址區間）
- 手動寫入 register（測試用途）

## Capabilities

### New Capabilities

- `workbench-diagnostic-panel`：跨步驟診斷 Slide-over，問題清單 + 一鍵跳轉
- `rule-template-quick-bar`：Source Rule 模板快捷欄，支援套用/儲存/預覽
- `database-output-workbench`：完整 DB 輸出工作台（多 connector、自動建表、dry-run、寫入歷史）
- `local-modbus-tool-page`：獨立的 Local Modbus 工具頁（/tools/modbus）

### Modified Capabilities

- `datalink-workbench-desktop`：
  - 主流程縮減為 4 步驟（移除 Local Modbus 步驟）
  - 元件拆分（SourceCanvasSection、WorkbenchDeviceStep、TagBindingStudio）
  - Step 1：DeviceTestConsole 分層顯示
  - Step 2：Plan/Live/Link Tab 列；RuleTemplateQuickBar 整合
  - Step 3：TagBatchActionBar 強化；Apply 後引導 CTA
  - WorkbenchDiagnosticPanel 整合至 WorkbenchFrame

- `local-modbus-memory-workbench`：移出主流程，重組為獨立工具頁

## Non-Goals

- 不改動後端 API 介面（除非 DB Output 新功能需要）
- 不改動 TestPage 的工程測試功能
- 不重設計 i18n 字典結構（僅新增 key）
- 不引入新的狀態管理函式庫

## Impact

### 前端

- `frontend/src/pages/studio/`（新建）：全新 Studio 主流程頁面
- `frontend/src/pages/tools/modbus/`（新建）：Local Modbus 工具頁
- `frontend/src/router/`：新增 `/studio` 與 `/tools/modbus` 路由
- `frontend/src/i18n/locales/`：新增診斷面板、模板欄、DB 輸出相關 key
- `frontend/docs/DESIGN-PRINCIPLES.md`（已建立）：所有新頁面的設計規範來源
- 舊 `frontend/src/pages/datalink/workbench/`：**保持不動**

### 後端（可能需要）

- DB Output 的歷史記錄查看：可能需要新增寫入日誌 API
- 多 connector 同時寫入：確認現有 `dbtarget` service 是否支援

### 測試

- 各拆分元件補獨立 Vitest 測試
- `WorkbenchDiagnosticPanel`：問題列表計算邏輯測試
- `RuleTemplateQuickBar`：模板套用/儲存狀態變化
- `DBConnectorManager`：多 connector CRUD
- Playwright E2E：完整 4 步驟主流程驗收
