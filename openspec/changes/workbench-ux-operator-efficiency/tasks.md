# Tasks: Workbench UX Operator Efficiency Overhaul

> 對應 proposal: `openspec/changes/workbench-ux-operator-efficiency/proposal.md`
> 對應 design: `openspec/changes/workbench-ux-operator-efficiency/design.md`

---

## Phase 1：元件拆分（重構）

> 目標：在不改變 UX 的前提下，建立乾淨的元件邊界，為後續 UX 強化奠定基礎。

- [ ] 1.1 將 `WorkbenchDeviceStep` (1713行) 拆分為：
  - `DeviceListPanel`（設備列表、選取）
  - `DeviceEditForm`（建立/編輯/Clone 表單）
  - `DeviceTestConsole`（Connect + Probe 測試結果顯示）
  - 補齊各子元件的 Vitest 測試

- [ ] 1.2 將 `SourceCanvasSection` (2688行) 拆分為：
  - `SourceRulePanel`（規則清單 + 規則表單）
  - `AddressCanvasView`（畫布渲染 + view mode 切換）
  - `SourceLivePanel`（Live 資料流 + 收集控制）
  - 補齊各子元件的 Vitest 測試

- [ ] 1.3 將 `TagBindingStudio` (1712行) 拆分為：
  - `TagCandidateBoard`（候選 Point 列表 + 篩選）
  - `TagBatchActionBar`（批次選取 + diff preview + apply）
  - 補齊各子元件的 Vitest 測試

- [ ] 1.4 確認拆分後現有測試（`DatalinkWorkbenchFoundation.test.tsx` 等 7 個）全部通過
- [ ] 1.5 更新 i18n 字典，確保所有文字鍵仍正確對應

---

## Phase 2：WorkbenchDiagnosticPanel（Slide-over 診斷面板）

> 目標：讓操作員點一下就能知道「卡在哪裡」並直接跳轉修復。

- [ ] 2.1 設計並實作 `WorkbenchDiagnosticPanel` 元件
  - Slide-over 右側推出，不覆蓋主工作區
  - 問題列表依步驟分組（device / source / tag / output）
  - 每個問題項目：描述 + 一鍵跳轉按鈕

- [ ] 2.2 實作診斷問題計算邏輯（`buildDiagnosticIssues`）
  - 涵蓋：Source Rule 衝突、Tag 未綁定、DB 欄位未映射、收集未啟動等

- [ ] 2.3 在 `WorkbenchFrame` 整合觸發入口：
  - 點擊任何 `partial`/`blocked` 色點 → 開啟診斷面板
  - 頂欄常駐「⚠ N 個問題」按鈕（有問題才顯示）

- [ ] 2.4 實作「一鍵跳轉並高亮」：跳轉步驟後，相關元素（規則、Tag、欄位）獲得視覺高亮
- [ ] 2.5 補 Vitest 測試：問題列表計算邏輯
- [ ] 2.6 補 Playwright E2E：觸發診斷面板 → 點跳轉 → 驗證高亮

---

## Phase 3：Step 2 RuleTemplateQuickBar + 畫布 Tab 列

> 目標：消除資深操作員每次重新輸入規則的重複勞動。

- [ ] 3.1 實作 `RuleTemplateQuickBar` 元件
  - 左側展開式固定欄位（桌面展開、小螢幕收合）
  - 顯示已儲存模板清單（名稱 + hover 預覽摘要）
  - 「套用」：一鍵填入規則表單（起始地址、數量、命名前綴）
  - 「儲存」：從當前 SourceRulePanel 狀態一鍵建立模板

- [ ] 3.2 模板儲存格式與現有 `sourceTemplateStorage.ts` 整合（擴充或重用）
- [ ] 3.3 將 Plan/Live/Link 模式切換改為畫布頂部明確 Tab 列，取代原 toolbar icon
- [ ] 3.4 更新 i18n 字典（模板快捷列、Tab 標籤）
- [ ] 3.5 補 Vitest 測試：套用模板後 SourceRulePanel 狀態驗證
- [ ] 3.6 補現有 `DatalinkWorkbenchSourceStep.test.tsx` 相關情境

---

## Phase 4：Step 3 TagBatchActionBar 強化

> 目標：將 Tag 批次操作從 5 步縮短到 2 步，Apply 後明確引導下一動作。

- [ ] 4.1 強化 `TagBatchActionBar`：
  - 新增「一鍵全選」，策略選擇（Create / Link Existing）合併至行內 dropdown
  - Apply 按鈕前顯示 diff preview（已有 `buildBatchDiffPreview`，強化視覺呈現）

- [ ] 4.2 Apply 成功後顯示進度完成橫幅：
  - 「✓ 已建立 N 個 Tag，前往 Step 4 設定資料庫輸出 →」（可點擊跳轉）

- [ ] 4.3 `TagCandidateBoard` 候選項目底色區別：
  - `partial`：amber 底色
  - `blocked`：rose 底色
  - `bound`：無特殊底色（正常）

- [ ] 4.4 `TagBindingFailure` 詳情：每個失敗項加入修復建議文字（依 stage 類型）
- [ ] 4.5 更新 i18n 字典（橫幅、批次操作標籤）
- [ ] 4.6 補 Vitest 測試：`TagBatchActionBar` 操作路徑 + 成功橫幅顯示

---

## Phase 5：Step 4 Database Output 完整工作台

> 目標：從「嵌在主頁面的表單」升級為「完整的 DB 輸出工作台」。

- [ ] 5.1 `DBConnectorManager`：多 connector 管理
  - 支援同時維護多個 connector（SQLite + PostgreSQL）
  - 每個 connector 有獨立的 `ConnectorSetupWizard`（Modal，4 步）
  - connector 狀態徽章：已連線 / 草稿 / 失敗

- [ ] 5.2 `TableSchemaGenerator`：一鍵建表
  - 依選取 Tag 的 data_type 生成 SQL（preview → 確認 → 執行）
  - 確認後端 `/api/db-targets/schema/generate` 或等價 API 支援

- [ ] 5.3 `ColumnMappingBoard`：映射狀態標示
  - 已映射（✓）/ 未映射（○）/ 衝突（✗）三態
  - 點擊未映射 Tag → 右側面板選擇目標 Table + Column

- [ ] 5.4 Dry-run 預覽模態：
  - 顯示每個 Tag → DB 欄位的預期寫入值
  - 標示欄位不存在等問題，阻擋套用

- [ ] 5.5 寫入模式設定：Upsert / Insert-only / Overwrite（下拉選擇，每個 connector 獨立設定）

- [ ] 5.6 歷史記錄查看（`WritingHistoryPanel`）：
  - 最近 N 次寫入的結果（時間、筆數、錯誤訊息）
  - 確認後端是否需要新增寫入日誌 API

- [ ] 5.7 頁面整體佈局：Tag 候選列表（左）+ 映射設定面板（右），頂部 connector 管理欄
- [ ] 5.8 更新 i18n 字典
- [ ] 5.9 補 Vitest 測試：DBConnectorManager CRUD、ColumnMappingBoard 狀態
- [ ] 5.10 補 Playwright E2E：完整 DB 輸出設定流程

---

## Phase 6：Local Modbus 工具頁（/tools/modbus）

> 目標：將 Local Modbus 從主流程移出，建立獨立的完整工具頁。

- [ ] 6.1 在 `frontend/src/pages/tools/modbus/` 建立工具頁目錄結構
- [ ] 6.2 新增路由 `/tools/modbus`，並在主導覽加入入口連結
- [ ] 6.3 `RegisterOverview`：HR 地址總覽表
  - 顯示全部 HR 的映射狀態（已映射 ●、空閑 ○、衝突 ✗）
  - 支援分頁 / 虛擬滾動（HR 數量可能達 10,000+）

- [ ] 6.4 `AutoMapEngine` + Dry-run：
  - 策略選擇（sequential / reverse）
  - Dry-run 預覽（顯示預期分配結果）→ 確認 → 套用

- [ ] 6.5 `ConflictResolver`：
  - 自動偵測所有衝突
  - 每個衝突提供候選解決方案（移動至空閑 register / 解除某個映射）
  - 一鍵套用解決方案

- [ ] 6.6 Live 監控：SSE 實時讀值，支援 Freeze / 快照
- [ ] 6.7 分組管理：定義地址區間（HR0-99 A 組，HR100-199 B 組）
- [ ] 6.8 手動寫入 register（含確認提示，防誤觸）
- [ ] 6.9 從舊的 `LocalModbusBoard` 中遷移現有映射邏輯（`computeAutoMap`、`buildRegisterUsage`）
- [ ] 6.10 更新路由：移除主流程 Step 4 的 Modbus 切換入口，改為工具頁連結
- [ ] 6.11 更新 i18n 字典
- [ ] 6.12 補 Vitest 測試：ConflictResolver 算法、AutoMapEngine dry-run
- [ ] 6.13 補 Playwright E2E：工具頁 AutoMap + 衝突解決流程

---

## 驗收條件

- [ ] 全部 Vitest 測試通過（包含新增的子元件測試）
- [ ] `npm run lint` 無錯誤
- [ ] `npm run build` 成功
- [ ] Playwright E2E：4 步驟主流程全流程通過
- [ ] Playwright E2E：診斷面板觸發 + 跳轉通過
- [ ] Playwright E2E：Local Modbus 工具頁主要流程通過
- [ ] 後端 `go test ./...` 通過（若有 API 異動）
