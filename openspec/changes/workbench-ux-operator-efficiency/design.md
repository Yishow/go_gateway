# Design: Workbench UX Operator Efficiency Overhaul

## 設計目標

以**資深操作員效率**為核心，重構 Datalink Workbench 的元件架構，同步強化各步驟的 UX 功能，使操作員能快速完成配置、批次處理、診斷問題。

---

## 架構設計

### 主流程：4 步驟 Workbench

```
Step 1: Device     — 建立 / 選擇資料來源設備，分層診斷 Connect / Probe
Step 2: Source     — 規劃 Source Rule，格狀畫布查看 planned/used/unmanaged/conflict
Step 3: Tag        — 批次 review 系統自動建立的 Tag / Mapping
Step 4: Database   — 設定 DB 輸出連線、映射、寫入模式
```

**Local Modbus** 獨立為工具頁，路由：`/tools/modbus`，從主導覽或 Tag 頁面的快速入口進入。

---

## Step 1：DeviceTestConsole 設計

### 元件拆分

```
WorkbenchDeviceStep (1713行)
  ├── DeviceListPanel      設備列表、選取、篩選
  ├── DeviceEditForm       建立 / 編輯 / Clone 表單
  └── DeviceTestConsole    Connect + Probe 分層結果
```

### DeviceTestConsole 佈局

```
┌─────────────────────────────────────────┐
│ [Connect 階段]          ●成功 / ✗失敗   │
│  主機：192.168.1.10:502                 │
│  耗時：45ms                            │
├─────────────────────────────────────────┤
│ [Probe 階段]            ●成功 / ✗失敗   │
│  協議：Modbus TCP                       │
│  Slave ID 1：讀取 HR0 → 1234           │
│  ⚠ 若 Connect 成功但 Probe 失敗：      │
│    「資料收集將被阻擋，直到 Probe 通過」 │
├─────────────────────────────────────────┤
│ 歷史：[2025-04-05 ✓] [2025-04-04 ✗]  │
└─────────────────────────────────────────┘
```

**設計原則：**
- Connect 欄與 Probe 欄視覺重量不同（Connect 是前提，字體較大）
- 兩個欄位各有獨立的「重新測試」按鈕
- 測試記錄保留最近 3 次，操作員可快速對比

---

## Step 2：SourceRulePanel + RuleTemplateQuickBar 設計

### 元件拆分

```
SourceCanvasSection (2688行)
  ├── SourceRulePanel       規則清單 + 規則表單
  ├── RuleTemplateQuickBar  模板快捷列（新增）
  ├── AddressCanvasView     畫布渲染 + Tab 切換
  └── SourceLivePanel       Live 資料流 + 收集控制
```

### RuleTemplateQuickBar 佈局（左側展開式）

```
┌──────────────────────────────┐
│  📋 規則模板                 │
│  ─────────────────────────── │
│  [Fatek D 區段]              │
│   D100-D120, uint16, d_       │（hover 顯示預覽）
│                [套用] [刪除]  │
│  [Modbus 標準 HR]            │
│   HR0-HR99, float32, hr_      │
│                [套用] [刪除]  │
│  ─────────────────────────── │
│  [+ 從當前規則儲存]          │
└──────────────────────────────┘
```

**設計原則：**
- 模板列預設展開（桌面）、收合（小螢幕）
- Hover 顯示起始地址、數量、命名前綴預覽
- 「從當前規則儲存」：一鍵，不需進入 toolbar 次級選單

### 畫布視模切換：Tab 列（取代 toolbar icon）

```
[Plan 模式] [Live 模式] [Link 模式]
```
清楚的 Tab 取代原本 toolbar 內的小 icon，每個模式有簡短說明文字。

---

## Step 3：TagBatchActionBar 強化設計

### 元件拆分

```
TagBindingStudio (1712行)
  ├── TagCandidateBoard    候選 Point 列表 + 篩選（unbound/partial/bound）
  └── TagBatchActionBar    批次操作工具列（強化）
```

### TagBatchActionBar 操作路徑

```
一鍵全選 → 選策略（Create / Link Existing）→ Diff Preview → Apply → 成功橫幅 + CTA
```

**Apply 後的成功橫幅：**
```
✓ 已建立 15 個 Tag，5 個綁定成功    [前往 Step 4 設定資料庫輸出 →]
```

**設計原則：**
- `partial` / `blocked` 候選項用不同底色區別（amber/rose），不只靠文字
- 批次操作從 5 步縮短為 2 步（選策略 + 確認）
- Diff preview 清楚標明每個 Point 將會發生的事（建立/綁定/跳過）

---

## Step 4：Database Output 完整工作台設計

### 頁面區域佈局

```
┌──────────────────────────────────────────────────────┐
│  DB Connector 管理欄（頂部）                          │
│  [SQLite: dev.db ●] [PostgreSQL: prod ●] [+ 新增]   │
├───────────────────────┬──────────────────────────────┤
│  Tag 候選列表         │  映射設定面板                 │
│  ─────────────────── │  ─────────────────────────── │
│  ▪ temp_1  [未映射]  │  表：sensor_data             │
│  ▪ temp_2  [已映射]  │  欄位：temperature (float)  │
│  ▪ press_1 [衝突⚠]  │                               │
│                       │  [Upsert ▾] [Dry-run] [套用] │
│                       │  ─────────────────────────── │
│                       │  寫入記錄（最近 5 次）        │
│                       │  2025-04-05 12:00 ✓ 15 筆   │
│                       │  2025-04-05 11:59 ✗ 錯誤…  │
└───────────────────────┴──────────────────────────────┘
```

### DBConnectorManager

- 支援多 connector（SQLite + PostgreSQL 同時存在）
- 每個 connector 有獨立的 `ConnectorSetupWizard`（Modal，4 步：選類型→填設定→測試→選 schema）
- connector 狀態徽章：● 連線中 / ○ 草稿 / ✗ 連線失敗

### TableSchemaGenerator

- 選取一組 Tag → 點「自動建表」→ 預覽 SQL → 確認執行
- 生成欄位依 Tag 的 `data_type` 映射（`float64` → `REAL`，`int32` → `INTEGER`，etc.）

### Dry-run 預覽模態

```
┌────────────────────────────────┐
│  Dry-run 預覽（不會真正寫入）  │
│  ──────────────────────────── │
│  temp_1 → temperature  ✓ 25.3│
│  press_1 → pressure    ✓ 1.01│
│  flow_1  → flow_rate   ✗ 欄位不存在 │
│                       [取消] [確認套用] │
└────────────────────────────────┘
```

---

## WorkbenchDiagnosticPanel（Slide-over）

### 觸發方式

1. 點擊任何步驟的 `partial` / `blocked` 色點
2. 頂欄常駐按鈕：「⚠ 3 個問題」（有問題時顯示，無問題時隱藏）

### Slide-over 結構

```
┌────────────────────────────────────┐（右側推出，不覆蓋主畫面）
│  ⚠ 工作台診斷                     │
│  ──────────────────────────────── │
│  Step 1 設備                   ✓  │
│  Step 2 來源規則               ⚠  │
│    ├ Rule #2 衝突（影響 3 Points） │
│    │                [前往並高亮 →] │
│    └ 收集未啟動                   │
│                     [啟動收集 →]  │
│  Step 3 Tag 綁定               ⚠  │
│    └ 5 個 Point 尚未綁定 Tag     │
│                   [前往 Step 3 →] │
│  Step 4 資料庫輸出             ●  │
│    └ 3 個 Tag 尚未映射 DB 欄位   │
│                   [前往 Step 4 →] │
└────────────────────────────────────┘
```

**設計原則：**
- 問題清單自動刷新（每次步驟切換後重新計算）
- 每個問題項目有：描述、一鍵跳轉、跳轉後高亮相關元素
- Slide-over 不阻擋操作員繼續操作主畫面

---

## Local Modbus 工具頁（/tools/modbus）

### 頁面佈局

```
┌─────────────────────────────────────────────────────┐
│  🔧 Local Modbus 工具                               │
│  [Register 總覽] [分組管理] [Live 監控] [手動寫入]  │
├──────────────────────────────────────────────────────┤
│  Register 總覽（主要視圖）                           │
│  [AutoMap ▾] [Dry-run] [偵測衝突] [套用]           │
│  ─────────────────────────────────────────────────── │
│  HR0001  temp_1    float32  25.3  ●已映射           │
│  HR0003  press_1   float32   1.01 ●已映射           │
│  HR0005  ─────────────────────── ○空閑             │
│  HR0006  flow_1 [衝突⚠]  temp_2  ✗衝突            │
│  ─────────────────────────────────────────────────── │
│  衝突詳情：HR0006 被 flow_1 + temp_2 同時映射        │
│  建議解法：[移動 temp_2 → HR0010] [解除 temp_2 映射]│
└──────────────────────────────────────────────────────┘
```

**功能規格：**
- **Register 總覽**：顯示全部 HR，狀態：已映射（●）/ 空閑（○）/ 衝突（✗）
- **AutoMap + Dry-run**：sequential/reverse 策略，干跑預覽後才套用
- **衝突解決**：自動偵測所有衝突，逐一顯示候選解決方案（移動/解除）
- **Live 監控**：SSE 即時值，可 Freeze / 快照
- **分組管理**：定義地址區間（HR0-99 給 A 組，HR100-199 給 B 組）
- **手動寫入**：輸入 register + 值 → 寫入（測試用，有確認提示）

---

## 資料流

```
[Step 1] DeviceListPanel → selectedDeviceId (WorkbenchProvider)
         DeviceTestConsole ← ConnectionTestResult (connect/probe 分層)

[Step 2] SourceRulePanel → SourceRule[] (DB 持久化)
         RuleTemplateQuickBar → applyTemplate → SourceRule draft
         AddressCanvasView ← buildAddressCanvasItems (planned/used/unmanaged/conflict)
         SourceLivePanel ← useRuntimeStream (SSE)

[Step 3] TagCandidateBoard ← points (unbound/partial/bound)
         TagBatchActionBar → buildBatchDiffPreview → apply → batchSummary → CTA

[Step 4] DBConnectorManager → 多 connector 管理 (SQLite + PostgreSQL)
         TableSchemaGenerator → 依 Tag 自動建表
         ColumnMappingBoard → Tag → DB column 綁定 (dry-run → apply)
         WritingHistoryPanel → 最近寫入記錄 / 錯誤日誌

[/tools/modbus] RegisterOverview ← modbusShareAPI
                AutoMapEngine → dry-run → apply
                ConflictResolver → 偵測衝突 → 候選解決方案
                LiveRegisterReader ← SSE
```

---

## 錯誤處理

| 場景 | 處理方式 |
|------|---------|
| Connect 失敗（網路不通） | DeviceTestConsole Connect 欄顯示錯誤，不執行 Probe |
| Connect 成功但 Probe 失敗 | 分欄顯示，標示「資料收集被阻擋」，允許儲存設備 |
| Source Rule 衝突 | conflict 格子 tooltip + 診斷面板問題項 |
| Tag Batch Apply 失敗 | TagBindingFailure[] 逐項說明 stage + 修復建議 |
| DB 連線測試失敗 | 行內顯示錯誤，允許 draft 模式儲存 connector |
| Modbus Register 衝突 | RegisterOverview 高亮 + ConflictResolver 候選解法 |
| Dry-run 發現欄位不存在 | 預覽模態標紅，阻擋套用，提示建表 |

---

## 測試策略

### 單元測試（Vitest）

| 元件/模組 | 測試重點 |
|----------|---------|
| `DeviceTestConsole` | Connect/Probe 狀態組合渲染 |
| `RuleTemplateQuickBar` | 套用/儲存模板的狀態變化 |
| `AddressCanvasView` | plan/live/link 模式切換 |
| `TagBatchActionBar` | diff preview 計算、batch apply 結果 |
| `DBConnectorManager` | 多 connector CRUD |
| `WorkbenchDiagnosticPanel` | 問題列表計算邏輯 |
| `ConflictResolver` | 衝突偵測算法 + 候選解法生成 |

### E2E 測試（Playwright）

- 完整 4 步驟主流程（從建立設備到 DB 映射套用）
- 診斷面板觸發與一鍵跳轉
- Local Modbus 工具頁 AutoMap + Dry-run 流程

---

## 實作分期建議

| Phase | 內容 | 說明 |
|-------|------|------|
| P1 | 元件拆分（重構）| 不改 UX，先建立乾淨的元件邊界 |
| P2 | WorkbenchDiagnosticPanel | 跨步驟診斷，立即提升操作員可見性 |
| P3 | Step 2 RuleTemplateQuickBar + Tab 列 | 解決最高頻的重複輸入痛點 |
| P4 | Step 3 TagBatchActionBar 強化 | 批次操作效率 |
| P5 | Step 4 Database Output 完整工作台 | 完整 DB 輸出功能 |
| P6 | Local Modbus 工具頁（/tools/modbus）| 獨立工具，完整 6 大功能 |
