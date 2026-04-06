# Tasks: Workbench UX Operator Efficiency Overhaul

> 對應 proposal: `openspec/changes/workbench-ux-operator-efficiency/proposal.md`
> 對應 design: `openspec/changes/workbench-ux-operator-efficiency/design.md`
> 設計規範: `frontend/docs/DESIGN-PRINCIPLES.md`

**策略：新建頁面，保留舊頁面**
- 新頁面位於 `frontend/src/pages/studio/` 與 `frontend/src/pages/tools/modbus/`
- 舊頁面 `frontend/src/pages/datalink/workbench/` 保持不動
- 路由：新 `/studio`、舊 `/datalink/workbench` 平行存在

---

## Phase 0：設計系統基礎建設

> 目標：建立新頁面所需的設計系統基礎，確保後續所有 UI 元件一致。

- [ ] 0.1 建立 CSS/Tailwind 設計系統 tokens（依 `DESIGN-PRINCIPLES.md`）
  - 背景層次 tokens（`--bg-void`, `--bg-base`, `--bg-elevated`, `--bg-surface`）
  - 文字 tokens、accent、狀態色
  - 確保 `tailwind.config.ts` 完整對映設計原則顏色

- [ ] 0.2 建立 Studio 頁面骨架目錄結構
  - `frontend/src/pages/studio/` 依設計文件建立完整目錄
  - `frontend/src/pages/tools/modbus/`
  - `StudioPage.tsx`、`StudioProvider.tsx` 基本框架

- [ ] 0.3 新增路由
  - `/studio` 路由指向 `StudioPage`（含步驟 query params）
  - `/tools/modbus` 路由指向 `LocalModbusToolPage`
  - 在主導覽加入 Studio 入口連結

- [ ] 0.4 建立共用 Studio UI 基本元件
  - `StudioTopBar`：步驟指示器 + 診斷按鈕槽
  - `StudioStepRail`：4 步驟切換側欄（含 readiness dots）
  - 依 `DESIGN-PRINCIPLES.md` 顏色系統設計

- [ ] 0.5 補 Vitest 測試：`StudioStepRail` 步驟切換與 readiness 顯示

---

## Phase 1：Studio Step 1 — Device 頁面

> 目標：全新設計設備管理頁面，Connect / Probe 分層診斷清晰呈現。
> **新建**：`frontend/src/pages/studio/device/`

- [ ] 1.1 建立 `StudioDevicePage.tsx`（頁面入口）
- [ ] 1.2 建立 `DeviceListPanel`（設備列表、選取、篩選）
  - 依 DESIGN-PRINCIPLES.md 卡片樣式
- [ ] 1.3 建立 `DeviceEditForm`（建立/編輯/Clone 表單）
  - 從 `WorkbenchDeviceStep` 提取業務邏輯（保留舊元件）
- [ ] 1.4 建立 `DeviceTestConsole`（Connect + Probe 分層結果）
  - Connect 欄與 Probe 欄各自獨立狀態顯示
  - Connect 成功但 Probe 失敗時顯示「資料收集將被阻擋」標語
  - 測試記錄保留最近 3 次歷史
- [ ] 1.5 補 Vitest 測試：各子元件渲染 + 狀態組合

---

## Phase 2：Studio Step 2 — Source 頁面

> 目標：全新設計 Source Rule 管理頁面，消除重複輸入痛點。
> **新建**：`frontend/src/pages/studio/source/`

- [ ] 2.1 建立 `StudioSourcePage.tsx`（頁面入口）
- [ ] 2.2 建立 `SourceRulePanel`（規則清單 + 規則表單）
  - 從 `SourceCanvasSection` 提取規則管理業務邏輯
- [ ] 2.3 建立 `RuleTemplateQuickBar`（左側展開式模板列）
  - 懸停預覽起始地址、數量、命名前綴
  - 「套用」一鍵填入規則表單
  - 「從當前規則儲存」一鍵建立模板
  - 整合 `sourceTemplateStorage.ts`（擴充或重用）
- [ ] 2.4 建立 `AddressCanvasView`（畫布 + Plan/Live/Link Tab 列）
  - Tab 列取代原 toolbar icon
  - 從 `SourceCanvasSection` 提取畫布渲染邏輯
- [ ] 2.5 建立 `SourceLivePanel`（Live 資料流 + 收集控制）
- [ ] 2.6 更新 i18n 字典（模板快捷列、Tab 標籤）
- [ ] 2.7 補 Vitest 測試：套用模板後 SourceRulePanel 狀態、Tab 模式切換

---

## Phase 3：Studio Step 3 — Tag 頁面

> 目標：全新設計 Tag 批次操作頁面，批次操作從 5 步縮短到 2 步。
> **新建**：`frontend/src/pages/studio/tag/`

- [ ] 3.1 建立 `StudioTagPage.tsx`（頁面入口）
- [ ] 3.2 建立 `TagCandidateBoard`（候選 Point 列表 + 篩選）
  - `partial`：amber 底色；`blocked`：rose 底色
  - 從 `TagBindingStudio` 提取業務邏輯
- [ ] 3.3 建立強化版 `TagBatchActionBar`
  - 「一鍵全選」+ 策略選擇（Create / Link Existing）合併至行內 dropdown
  - Diff preview：清楚標明每個 Point 將建立/綁定/跳過
  - Apply 後成功橫幅 + CTA（「前往 Step 4 設定資料庫輸出 →」）
- [ ] 3.4 `TagBindingFailure` 詳情：依 stage 類型加入修復建議文字
- [ ] 3.5 更新 i18n 字典（橫幅、批次操作標籤）
- [ ] 3.6 補 Vitest 測試：`TagBatchActionBar` 操作路徑 + 成功橫幅顯示

---

## Phase 4：WorkbenchDiagnosticPanel（Slide-over 診斷面板）

> 目標：讓操作員點一下就能知道「卡在哪裡」並直接跳轉修復。
> **新建**：`frontend/src/pages/studio/components/WorkbenchDiagnosticPanel.tsx`

- [ ] 4.1 設計並實作 `WorkbenchDiagnosticPanel` 元件
  - Slide-over 右側推出，不覆蓋主工作區
  - 問題列表依步驟分組（device / source / tag / output）
  - 每個問題項目：描述 + 一鍵跳轉按鈕

- [ ] 4.2 實作診斷問題計算邏輯（`buildDiagnosticIssues`）
  - 涵蓋：Source Rule 衝突、Tag 未綁定、DB 欄位未映射、收集未啟動等

- [ ] 4.3 在 `StudioTopBar` 整合觸發入口：
  - 頂欄常駐「⚠ N 個問題」按鈕（有問題才顯示）
  - 點擊任何 `partial`/`blocked` 色點 → 開啟診斷面板

- [ ] 4.4 實作「一鍵跳轉並高亮」：跳轉步驟後，相關元素獲得視覺高亮
- [ ] 4.5 補 Vitest 測試：問題列表計算邏輯
- [ ] 4.6 補 Playwright E2E：觸發診斷面板 → 點跳轉 → 驗證高亮

---

## Phase 5：Studio Step 4 — Database Output 頁面

> 目標：全新設計完整 DB 輸出工作台。
> **新建**：`frontend/src/pages/studio/output/`

- [ ] 5.1 建立 `StudioDatabasePage.tsx`（頁面入口）
- [ ] 5.2 `DBConnectorManager`：多 connector 管理
  - 支援同時維護多個 connector（SQLite + PostgreSQL）
  - 每個 connector 有獨立的 `ConnectorSetupWizard`（Modal，4 步）
  - connector 狀態徽章：已連線 / 草稿 / 失敗

- [ ] 5.3 `TableSchemaGenerator`：一鍵建表
  - 依選取 Tag 的 data_type 生成 SQL（preview → 確認 → 執行）
  - 確認後端 `/api/db-targets/schema/generate` 或等價 API 支援

- [ ] 5.4 `ColumnMappingBoard`：映射狀態標示
  - 已映射（✓）/ 未映射（○）/ 衝突（✗）三態

- [ ] 5.5 Dry-run 預覽模態
- [ ] 5.6 寫入模式設定：Upsert / Insert-only / Overwrite
- [ ] 5.7 `WritingHistoryPanel`：最近 N 次寫入記錄
- [ ] 5.8 更新 i18n 字典
- [ ] 5.9 補 Vitest 測試：DBConnectorManager CRUD、ColumnMappingBoard 狀態
- [ ] 5.10 補 Playwright E2E：完整 DB 輸出設定流程

---

## Phase 6：Local Modbus 工具頁（/tools/modbus）

> 目標：將 Local Modbus 從主流程移出，建立獨立的完整工具頁。
> **新建**：`frontend/src/pages/tools/modbus/`

- [ ] 6.1 建立 `LocalModbusToolPage.tsx`（頁面入口）
- [ ] 6.2 `RegisterOverview`：HR 地址總覽表（已映射 ●、空閑 ○、衝突 ✗）
  - 支援分頁 / 虛擬滾動（HR 數量可能達 10,000+）
- [ ] 6.3 `AutoMapEngine` + Dry-run：策略選擇 → preview → 套用
- [ ] 6.4 `ConflictResolver`：偵測衝突 + 候選解決方案 + 一鍵套用
- [ ] 6.5 Live 監控：SSE 實時讀值，支援 Freeze / 快照
- [ ] 6.6 分組管理：定義地址區間
- [ ] 6.7 手動寫入 register（含確認提示）
- [ ] 6.8 從 `LocalModbusBoard` 遷移現有映射邏輯（`computeAutoMap`、`buildRegisterUsage`）
- [ ] 6.9 在主導覽加入 `/tools/modbus` 入口連結
- [ ] 6.10 更新 i18n 字典
- [ ] 6.11 補 Vitest 測試：ConflictResolver 算法、AutoMapEngine dry-run
- [ ] 6.12 補 Playwright E2E：工具頁 AutoMap + 衝突解決流程

---

## 驗收條件

### 設計規範驗收

- [ ] 所有新頁面元件通過 `DESIGN-PRINCIPLES.md` 中的設計驗收標準
- [ ] 顏色系統：使用設計 token，不出現 hardcoded hex（除 token 定義本身）
- [ ] 深色模式：所有新頁面在深色背景下可讀

### 功能驗收

- [ ] 全部 Vitest 測試通過（包含新增的子元件測試）
- [ ] `npm run lint` 無錯誤
- [ ] `npm run build` 成功
- [ ] Playwright E2E：4 步驟主流程全流程通過
- [ ] Playwright E2E：診斷面板觸發 + 跳轉通過
- [ ] Playwright E2E：Local Modbus 工具頁主要流程通過
- [ ] 後端 `go test ./...` 通過（若有 API 異動）

### 相容性驗收

- [ ] 舊路由 `/datalink/workbench` 維持可正常使用
- [ ] 新舊路由平行存在無衝突
