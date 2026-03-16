# Datalink Workbench — Phase 1 UI/互動細規

- Parent Spec: `2026-03-15-datalink-workbench-design.md`
- Scope: Release 1 — `Source → Tag → Local Modbus` 完整主路徑
- 目標讀者: Opus UI 子代理（可直接依此實作）

---

## 1. 頁面版面配置

### 1.1 整體結構

```
┌─────────────────────────────────────────────────────────┐
│  App shell sidebar（workbench 外部容器，非本 spec 範圍）│
│  ┌───────────────────────────────────────────────────┐  │
│  │  DatalinkWorkbenchPage (flex-1, 佔滿 main area)   │  │
│  │  ┌──┬──────────────────────────────────────────┐  │  │
│  │  │  │ ContextBar (sticky top, h-14)            │  │  │
│  │  │  ├──────────────────────────────┬───────────┤  │  │
│  │  │S │                              │           │  │  │
│  │  │t │   PrimaryWorkArea            │ Inspector │  │  │
│  │  │e │   (SourceCanvas / TagBinder  │ Panel     │  │  │
│  │  │p │    / OutputTargets)          │ (w-80~96) │  │  │
│  │  │  │                              │           │  │  │
│  │  │R │   視 step 切換顯示           │ 視 step   │  │  │
│  │  │a │                              │ 切換內容  │  │  │
│  │  │i │                              │           │  │  │
│  │  │l │                              │           │  │  │
│  │  │  ├──────────────────────────────┴───────────┤  │  │
│  │  │  │ BottomSummaryBar (sticky bottom, h-12)   │  │  │
│  │  └──┴──────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 1.2 區塊尺寸與權重

| 區塊 | 寬度策略 | 高度策略 | z-index 層級 |
|------|---------|---------|-------------|
| StepRail | `w-14` (icon-only), xl 以上 `w-48` (icon+label) | 填滿 workbench 高度 | z-10 |
| ContextBar | `flex-1` (剩餘寬度) | `h-14` 固定 | z-20 (sticky) |
| PrimaryWorkArea | `flex-1` (扣除 StepRail + Inspector) | `min-h-0 flex-1` 填滿中間 | z-0 |
| InspectorPanel | `w-80` (320px), 可收合為 `w-0` | 與 PrimaryWorkArea 同高 | z-10 |
| BottomSummaryBar | 全寬 | `h-12` 固定 | z-20 (sticky) |

### 1.3 Sticky 行為

- **ContextBar**: `sticky top-0`，隨 PrimaryWorkArea 捲動時保持可見
- **BottomSummaryBar**: `sticky bottom-0`，始終顯示驗證/計數摘要
- **StepRail**: 固定於左側，不捲動（`sticky top-0 self-start h-full`）
- **InspectorPanel**: 內部可獨立滾動 (`overflow-y-auto`)，外框不隨主區滾動

### 1.4 響應式斷點行為

| 斷點 | StepRail | Inspector | Grid 欄數 | ContextBar |
|------|----------|-----------|----------|------------|
| `< md` (768) | 底部水平 tab bar 替代 | 全螢幕 slide-over | 5 欄 | 簡化：僅設備名 + 狀態燈 |
| `md–lg` (768–1024) | `w-14` icon-only | overlay drawer (右側滑出) | 8 欄 | 省略 quick actions，合併至 ⋯ 選單 |
| `lg–xl` (1024–1280) | `w-14` icon-only | inline `w-72` | 10 欄 | 完整 |
| `≥ xl` (1280+) | `w-48` icon+label | inline `w-80` | 10 欄（寬裕間距） | 完整 + 連線測試結果時間戳 |

### 1.5 外層容器樣式基準

```
容器: min-h-[calc(100vh-11rem)]
背景: bg-gradient-to-br from-[#0B0F19] via-[#111827] to-[#0F172A]
文字: text-slate-100 font-sans
圓角: rounded-2xl overflow-hidden
選取: selection:bg-blue-500/30
```

> 採用新的 dark-first 工作台視覺基底，不承接舊 SmartDashboard 的元件樣式。

---

## 2. 各區塊元件責任與內容

### 2.1 StepRail

**責任**: 持久導航、步驟狀態指示、步驟跳轉

**內容**:
| # | Step | Icon | 狀態指示 |
|---|------|------|---------|
| 1 | Device | `Server` / `Cpu` | ● 綠=已選擇有效設備、○ 灰=未選 |
| 2 | Source Grid | `Grid3x3` / `LayoutGrid` | badge 數字=已建 point 數、⚠ 有衝突 |
| 3 | Tag Binding | `Tags` / `Link` | badge=已綁定 / 總數、⚠ 有未綁 |
| 4 | Output Targets | `ArrowRightFromLine` / `Database` | ● 綠=全映射完成、◐ 半=部分完成 |

**互動規則**:
- 點擊任何已解鎖步驟可跳轉（無 wizard-style 線性鎖定）
- Step 1 (Device) 未完成前，Step 2–4 顯示為 dimmed 但可點擊（跳轉後顯示 empty state 引導回 Step 1）
- 當前步驟高亮底色 `bg-blue-600/20 border-l-2 border-blue-500`
- Hover: `bg-slate-700/50` 過渡 150ms

**視覺細節**:
- Icon 尺寸: `w-5 h-5`
- Step 標籤字體: `text-xs font-medium` (xl 以上展開時)
- 分隔線: 步驟間用 `border-b border-slate-700/50` 微分
- 底部保留空間放「返回舊版」連結: `text-xs text-slate-500 hover:text-slate-400`

### 2.2 ContextBar

**責任**: 全局上下文展示、快速操作入口

**結構** (左→右):
```
[ 設備名稱 + 協定 badge ] [ 連線狀態燈 ] [ 最後測試結果 ]  ——彈性空白——  [ 測試連線 btn ] [ 重新整理 btn ] [ 切換設備 btn ]
```

**內容細節**:

| 元素 | 規格 |
|------|------|
| 設備名稱 | `text-sm font-semibold text-slate-100` truncate max-w-[200px] |
| 協定 badge | `text-xs px-2 py-0.5 rounded-full bg-slate-700 text-slate-300` 例: `Modbus TCP` |
| 連線狀態 | 圓點 `w-2 h-2 rounded-full` + 文字 `text-xs`。綠=connected, 紅=disconnected, 黃=warning, 灰=unknown |
| 最後測試 | `text-xs text-slate-400` 顯示 "上次測試: 2 分鐘前 ✓" 或 "✗ timeout" |
| 測試連線 | Ghost button `text-sm`，按下後顯示 spinner → 結果 toast |
| 重新整理 | Icon button (RefreshCw)，旋轉動畫 during loading |
| 切換設備 | Icon button (ChevronDown) → 展開設備選擇下拉 popover |

**Step 1 未選設備時**:
- 顯示: `「請選擇設備以開始」` placeholder + 設備選擇 button (primary style)

### 2.3 SourceCanvas (Step 2 主區域)

**責任**: 地址空間視覺化、Point 建立/選取、Grid/Table 模式切換

**子結構**:
```
┌─ SourceCanvas ──────────────────────────────────────┐
│ ┌─ ToolBar ────────────────────────────────────────┐│
│ │ [Grid模式|Table模式]  起始地址  資料型態  數量    ││
│ │ [套用規劃]  [批次建立]  [匯入]  [匯出]           ││
│ └──────────────────────────────────────────────────┘│
│                                                     │
│ ┌─ GridView / TableView ──────────────────────────┐ │
│ │                                                 │ │
│ │   (Grid 或 Table 依模式切換顯示)                 │ │
│ │                                                 │ │
│ └─────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

#### 2.3.1 ToolBar

**模式切換**: Segmented control (兩態 toggle)
- `Grid` (icon: Grid3x3) — 預設
- `Table` (icon: List)
- 樣式: `bg-slate-800 rounded-lg p-0.5` 內部 button `px-3 py-1.5 text-xs rounded-md`
- 選中態: `bg-blue-600 text-white`
- 未選中: `text-slate-400 hover:text-slate-200`

**Source Planner 控制列** (Grid 模式專屬):
- 起始地址: `input text-sm w-24` (例: `D0`, `40001`)
- 資料型態: `select text-sm w-32` (int16 / uint16 / int32 / float32 / int64)
- 數量: `input type=number text-sm w-16`
- 套用規劃: Primary button `bg-blue-600 hover:bg-blue-700 text-sm px-3 py-1.5`
- 批次建立: Secondary button (選擇多個空格後啟用)
- 匯入/匯出: Ghost buttons with icon

#### 2.3.2 GridView（主角）

**以全新的 `AddressCanvas` 核心實作以下能力**:

**Cell 格線配置**:
- 欄數: 10（≥lg）、8（md）、5（<md）
- Row header: 顯示基底地址（如 `D0`, `D10`, `D20`…），`text-xs text-slate-500 font-mono w-12`
- Column header: 偏移量 `+0` `+1` … `+9`，`text-xs text-slate-500 text-center`

**Cell 視覺狀態** (優先序由高到低):

| 狀態 | 背景 | 邊框 | 文字 | 附加 |
|------|------|------|------|------|
| conflict | `bg-rose-900/40` | `ring-2 ring-rose-500` | `text-rose-200` | 閃爍 icon ⚠ |
| selected | `bg-blue-600/30` | `ring-2 ring-blue-400` | `text-blue-100` | — |
| linked (已綁 tag) | `bg-violet-900/30` | `border border-violet-600/50` | `text-violet-200` | 小 tag icon |
| used (已建 point) | `bg-emerald-900/30` | `border border-emerald-700/40` | `text-emerald-200` | — |
| planned (規劃中) | `bg-sky-900/20` | `border border-sky-700/30 border-dashed` | `text-sky-300` | — |
| hover | `bg-slate-700/40` | — | — | cursor pointer |
| available (空) | `bg-slate-800/50` | `border border-slate-700/20` | `text-slate-500` | 顯示地址 |
| blocked | `bg-slate-900/60` | — | `text-slate-600` | 斜線或 lock icon |

**Cell 內容顯示**:
- 空格: 地址文字 `text-xs font-mono text-center`
- 有 Point: Point 名稱 (truncate) `text-xs font-medium`
- 多 cell 合併: 依據資料型態 span (int16=1, int32/float32=2, int64=4)，合併 cell 顯示連續背景 + 單一標籤置中

**選取互動**:
- 單擊: 選取單格 → Inspector 顯示詳情
- Ctrl/Cmd+Click: 加選/減選（toggle）
- Shift+Click: 範圍選取（矩形區域）
- 拖曳選取: mousedown → mousemove → mouseup 框選
- 右鍵: context menu（建立 Point / 編輯 / 刪除 / 快速綁 Tag）

#### 2.3.3 TableView（輔助）

**用途**: 稠密審計、批次編輯、排序篩選

**欄位定義**:

| 欄位 | 寬度 | 可排序 | 可篩選 |
|------|------|--------|--------|
| ☐ (checkbox) | w-10 | — | — |
| 地址 | w-24 | ✓ | ✓ (搜尋) |
| 名稱 | flex-1 min-w-[120px] | ✓ | ✓ |
| 資料型態 | w-28 | ✓ | ✓ (下拉) |
| Span | w-16 | ✓ | — |
| 當前值 | w-24 | — | — |
| 轉換狀態 | w-20 | — | ✓ |
| Tag 綁定 | w-32 | ✓ | ✓ (綁定/未綁) |
| 輸出就緒 | w-20 | — | ✓ |
| 操作 | w-20 | — | — |

**樣式基準**:
- 表頭: `sticky top-0 bg-slate-800 text-xs text-slate-400 font-medium uppercase tracking-wider`
- 行: `bg-slate-900/50 hover:bg-slate-800/70 border-b border-slate-700/30`
- 交替行: `even:bg-slate-850/30`（微弱區隔）
- 已選行: `bg-blue-900/20 ring-1 ring-blue-500/30`
- 行高: `h-10` (40px)，緊湊但可讀

**批次操作 toolbar** (顯示於選取行 > 0 時):
- 浮動 bar: `sticky top-[40px] bg-blue-900/90 backdrop-blur px-4 py-2 rounded-lg`
- 內容: `已選 {n} 項` + `[批次建立 Tag]` `[批次刪除]` `[修改資料型態]`

### 2.4 InspectorPanel

**責任**: 選取物件的深度檢視、屬性編輯、轉換預覽、衝突診斷

**通用框架**:
```
┌─ Inspector ─────────────────────┐
│ ┌─ Header ────────────────────┐ │
│ │ 標題      [收合 btn] [⋯]   │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─ Content (scrollable) ─────┐ │
│ │                             │ │
│ │  (依 step + 選取狀態填充)    │ │
│ │                             │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─ Footer Actions ───────────┐ │
│ │ [主要操作 btn] [次要操作]   │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

**依 Step 切換的內容**:

#### Step 1 (Device) Inspector:
- 設備詳情: 名稱、協定、IP/Port、unit ID
- 連線參數摘要
- 最近連線測試歷史（最近 3 筆）
- Footer: `[編輯設備]` `[測試連線]`

#### Step 2 (Source Grid) Inspector:

**無選取時**: 空狀態提示「點選 Grid 上的地址格來檢視詳情」

**選取單格（空格）**:
- 地址: `D100` (mono, large)
- 狀態: `可用`
- Footer: `[建立 Point]`

**選取單格（已有 Point）**:
- Section: 地址資訊
  - 地址、資料型態、Span 長度
- Section: Point 屬性
  - 名稱（可編輯 inline）
  - 描述（可編輯）
  - 當前值（顯示 "—" 如無 polling, Release 2 顯示 live value）
- Section: 轉換預覽
  - 套用的 transform steps 清單
  - 預覽結果（如有值）
- Section: 綁定狀態
  - Tag: `已綁定 → sensor.temp.01` 或 `未綁定`
  - 輸出: `Modbus: 40001` 或 `未映射`
- Section: 衝突 (如有)
  - 衝突描述 + 修正建議
- Footer: `[編輯]` `[綁定 Tag]` `[刪除]`

**選取多格**:
- 選取摘要: `已選 {n} 個地址`
- 細分: `{a} 空格 / {b} 已建 Point / {c} 已綁 Tag`
- 衝突摘要（如有）
- Footer: `[批次建立 Point]` `[批次綁定 Tag]`

#### Step 3 (Tag Binding) Inspector:
- 選取的 Point(s) 摘要
- Tag key 編輯（批次自動命名模板）
- display_name / unit / description 欄位
- 已存在相同 key 的衝突警告
- Footer: `[套用綁定]` `[取消]`

#### Step 4 (Output Targets) Inspector:
- 選取的 mapping 詳情
- Register 地址、資料型態
- 衝突偵測結果
- Footer: `[儲存映射]` `[移除映射]`

### 2.5 TagBindingStudio (Step 3 主區域)

**責任**: 批次 tag 建立/連結的專屬工作面

**版面**:
```
┌─ TagBindingStudio ──────────────────────────────────┐
│ ┌─ 左半: 來源 Points ──────┐ ┌─ 右半: Tag 設定 ──┐ │
│ │ 已選 Points 清單          │ │ 命名模板           │ │
│ │ (checkbox + 地址 + 名稱)  │ │ [prefix].[auto]    │ │
│ │                           │ │                    │ │
│ │ 篩選: 全部/未綁定/已綁定   │ │ 批次欄位:          │ │
│ │                           │ │  display_name 模板 │ │
│ │ ┌─────────────────────┐  │ │  unit 統一設定      │ │
│ │ │ point 列表…          │  │ │  description 模板   │ │
│ │ │                     │  │ │                    │ │
│ │ └─────────────────────┘  │ │ 預覽清單:           │ │
│ └──────────────────────────┘ │ │  (生成的 tag keys) │ │
│                              │ └────────────────────┘│
│ ┌─ 底部: 操作 ──────────────────────────────────────┐│
│ │ [全選] [反選]   已選 {n}/{total}    [套用綁定]     ││
│ └───────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────┘
```

**命名模板系統**:
- 格式: `{device}.{area}.{index:03d}`
- 可用變數: `{device}`, `{protocol}`, `{address}`, `{index}`, `{name}`, `{type}`
- 即時預覽: 每個 point 對應的生成 key 即時顯示
- 衝突檢查: 重複 key 即時標記紅色

**左半 Point 列表**:
- 每行: `☐ D100 | temperature_raw | int16 | 未綁定`
- 狀態 badge: 綠=已綁、灰=未綁、紅=衝突
- Shift+Click 多選

### 2.6 OutputTargets (Step 4 主區域)

**責任**: 輸出目標的統一管理面板

**版面**:
```
┌─ OutputTargets ─────────────────────────────────────┐
│ ┌─ Tab: [Local Modbus] [Database (即將推出)] ──────┐│
│ └──────────────────────────────────────────────────┘│
│                                                     │
│ ┌─ LocalModbusBoard (Phase 1 完整) ───────────────┐ │
│ │ ┌─ 伺服器狀態列 ──────────────────────────────┐ │ │
│ │ │ 狀態: 🟢 運行中 | Port: 5020 | [啟動/停止]  │ │ │
│ │ └─────────────────────────────────────────────┘ │ │
│ │                                                 │ │
│ │ ┌─ Mapping 工作區 ────────────────────────────┐ │ │
│ │ │ ToolBar: [新增映射] [自動映射] [匯入] [匯出] │ │ │
│ │ │ [同步全部]                                   │ │ │
│ │ │                                              │ │ │
│ │ │ Mapping 表格:                                │ │ │
│ │ │  Tag Key | Register | 型態 | 衝突 | 操作     │ │ │
│ │ │  ────────────────────────────────────────── │ │ │
│ │ │  sensor.temp.01 | 40001 | int16 | — | ✏🗑  │ │ │
│ │ │  sensor.temp.02 | 40002 | int16 | — | ✏🗑  │ │ │
│ │ │  …                                          │ │ │
│ │ └─────────────────────────────────────────────┘ │ │
│ │                                                 │ │
│ │ ┌─ 衝突治理面板 (條件顯示) ───────────────────┐ │ │
│ │ │ ⚠ 3 個 register 衝突                        │ │ │
│ │ │ 40005: sensor.pressure ↔ motor.speed         │ │ │
│ │ │ [自動解決] [逐一檢視]                        │ │ │
│ │ └─────────────────────────────────────────────┘ │ │
│ └─────────────────────────────────────────────────┘ │
│                                                     │
│ ┌─ DatabaseTargetCard (Phase 1: 停用) ───────────┐ │
│ │ 🔒 資料庫目標輸出                               │ │
│ │ 此功能需要後端契約完成後才能使用。                 │ │
│ │ 預計在 Release 2 提供。                          │ │
│ └─────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

**LocalModbusBoard 完整功能**:
- 伺服器控制: 啟動/停止 toggle + port 設定
- 映射 CRUD: 新增/編輯/刪除 tag-to-register mapping
- 自動映射: 依據已綁定 tag 順序自動分配 register
- 衝突偵測: 相同 register 被多個 tag 使用時即時警告
- 同步: 將 draft 映射寫入後端
- 測試寫入: 選擇一個映射 → 寫入測試值 → 確認讀回
- 匯入/匯出: JSON 格式

**DatabaseTargetCard (Phase 1)**:
- 視覺: `opacity-60` 整體降低
- 內容: 鎖定 icon + 說明文字 + 預計時間
- 不可互動，但佔位以表達完整 IA

### 2.7 BottomSummaryBar

**責任**: 全局狀態摘要、驗證就緒度、一鍵操作入口

**結構**:
```
[Step badge]  |  Points: 24  |  Tags: 18/24  |  Modbus: 15  |  DB: —  |  驗證: ✓ 就緒  |  [部署/同步 btn]
```

**各指標規格**:

| 指標 | 格式 | 狀態色 |
|------|------|--------|
| Points | `{n} 個點位` | 灰=0, 白>0 |
| Tags | `{bound}/{total}` | 綠=100%, 黃=部分, 灰=0% |
| Modbus | `{mapped} 映射` | 綠=全映射, 黃=部分, 紅=有衝突 |
| Database | `—` (Phase 1) | 灰 |
| 驗證 | `✓ 就緒` / `⚠ {n} 問題` / `✗ 阻塞` | 綠/黃/紅 |

**樣式**:
- 底色: `bg-slate-900/95 backdrop-blur-sm border-t border-slate-700/50`
- 指標間分隔: `border-l border-slate-700/30`
- 文字: `text-xs font-medium`
- 部署按鈕: Primary `bg-blue-600` 或 Disabled `bg-slate-700 cursor-not-allowed`

---

## 3. 關鍵互動狀態

### 3.1 狀態定義與視覺映射

| 狀態 | 觸發條件 | 背景 | 邊框 | Icon/Badge | 動畫 |
|------|---------|------|------|-----------|------|
| **hover** | 滑鼠懸停任何可互動元素 | 該元素 bg 亮一階 | — | — | `transition-colors duration-150` |
| **selected** | 使用者點選 | `bg-blue-600/25` | `ring-2 ring-blue-400` | — | `transition-all duration-100` |
| **conflict** | 地址重疊/register 重複/tag key 重複 | `bg-rose-900/30` | `ring-2 ring-rose-500` | ⚠ 三角 | `animate-pulse` (邊框) 2s 週期 |
| **linked** | Point 已綁定 Tag | `bg-violet-900/25` | `border border-violet-600/40` | 🔗 小 tag icon | — |
| **blocked** | 前置條件未滿足 | `bg-slate-900/60` | — | 🔒 lock | `opacity-50 cursor-not-allowed` |
| **saving** | 任何 mutation 進行中 | 不變 | — | spinner 替換操作 icon | `animate-spin` on spinner |
| **success** | mutation 成功 | 短暫 `bg-emerald-900/20` flash | `ring-2 ring-emerald-400` | ✓ checkmark | flash 1s → fade out |
| **empty** | 無資料 / 未開始 | `bg-slate-800/30` | `border border-dashed border-slate-600` | illustration + 引導文字 | — |
| **error** | API 錯誤 / 驗證失敗 | `bg-red-900/20` | `ring-1 ring-red-500/50` | ✗ 圓叉 | toast 彈出 |
| **draft** | 有未儲存變更 | 不變 | `border-l-2 border-amber-500` | 小圓點 ● amber | — |

### 3.2 狀態轉換場景

**Grid Cell 生命週期**:
```
available → (規劃) → planned → (建立) → used → (綁定 tag) → linked → (映射 Modbus) → linked+output
                                    ↘ (衝突) → conflict
```

**Mutation 操作回饋流程**:
```
idle → (click action) → saving (spinner + disable btn) → success (flash + toast) → idle
                                                       → error (shake + toast + inline msg)
```

**批次操作結果**:
```
saving → mixed result:
  ├── 成功 items: 短暫 success 高亮 → 正常狀態
  ├── 失敗 items: 持續 error 高亮 + Inspector 顯示逐項失敗原因
  └── Summary toast: "已完成 18/20，2 項失敗 [查看詳情]"
```

### 3.3 選取模式的優先序

當 cell 同時具有多個狀態時，視覺優先序:
1. `conflict` (最高 — 安全性)
2. `selected` (互動中)
3. `error`
4. `linked`
5. `used`
6. `planned`
7. `hover`
8. `available` (最低)

---

## 4. Grid 與表格的任務分工

### 4.1 Grid 承擔的任務（主角）

| 任務 | 說明 |
|------|------|
| **空間理解** | 使用者看到地址空間的二維佈局，理解哪些地址已用、哪些可用 |
| **連續 span 視覺化** | int32/float32 合併 2 格、int64 合併 4 格，直觀看出資料型態佔用 |
| **衝突定位** | 衝突 cell 以紅色高亮，在空間上立即可見（無須掃描表格） |
| **範圍規劃** | 使用者框選一段空地址 → 套用規劃 → 預覽結果 → 確認建立 |
| **快速巡覽** | 一眼掃過 100+ 地址的使用狀態（密度是表格的 5–10 倍） |
| **綁定狀態總覽** | 紫色=已綁 Tag、綠色=已建 Point、灰色=空，一目了然 |
| **context 保持** | 操作 Inspector 時 Grid 保持可見，不跳頁 |

### 4.2 Table 承擔的任務（輔助）

| 任務 | 說明 |
|------|------|
| **稠密審計** | 逐行檢視每個 Point 的完整屬性（地址、名稱、型態、值、tag、輸出） |
| **排序篩選** | 依地址/名稱/型態/綁定狀態排序，快速找到特定 point |
| **批次編輯** | Checkbox 多選 → toolbar 批次操作（改型態、刪除、綁定） |
| **數值顯示** | 當前值/轉換值的精確數字顯示（Grid cell 空間不足） |
| **匯出準備** | 表格是匯出 CSV/JSON 前的最後確認介面 |
| **無障礙替代** | 螢幕閱讀器使用者可用表格的語義結構操作 |

### 4.3 兩者的協同

- **選取同步**: Grid 選取的 cell 在切換到 Table 時，對應行自動高亮（反之亦然）
- **模式記憶**: 切換模式不清除選取狀態
- **Inspector 共享**: 兩種模式共用同一個 Inspector panel
- **快捷鍵一致**: `Ctrl+A` 全選在兩種模式下行為一致

---

## 5. Production-Ready UI 細節

### 5.1 視覺層次

**四層深度系統**:

| 層 | 用途 | 背景色參考 | elevation |
|----|------|-----------|-----------|
| L0 (page bg) | 頁面底色 | gradient `#0B0F19 → #0F172A` | — |
| L1 (sections) | 各區塊面板 | `bg-slate-800/80` | `shadow-sm` |
| L2 (cards/inputs) | 卡片、輸入框 | `bg-slate-700/60` | `shadow-md` |
| L3 (popovers/menus) | 浮層、下拉 | `bg-slate-700` | `shadow-xl` + `ring-1 ring-slate-600/50` |

**文字層次**:
| 層級 | Class | 用途 |
|------|-------|------|
| 主標題 | `text-lg font-semibold text-slate-50` | 區塊標題 |
| 次標題 | `text-sm font-medium text-slate-200` | Section 標題 |
| 正文 | `text-sm text-slate-300` | 說明文字、值顯示 |
| 輔助 | `text-xs text-slate-400` | 提示、時間戳、計數 |
| 禁用 | `text-xs text-slate-500` | 不可操作項目 |
| 數值 | `text-sm font-mono text-slate-100` | 地址、數值、register |

### 5.2 回饋機制

**操作回饋矩陣**:

| 操作 | 即時回饋 | 載入中 | 成功 | 失敗 |
|------|---------|--------|------|------|
| 建立 Point | button disable + spinner | — | toast success + cell 變綠 flash | toast error + inline message |
| 批次建立 | progress bar (n/total) | 逐項更新進度 | summary toast + grid 批次 flash | partial failure: 成功項 flash + 失敗項 error 高亮 |
| 綁定 Tag | button disable + spinner | — | cell 變紫 + toast | toast + Inspector 顯示失敗原因 |
| 刪除 | confirm dialog → spinner | — | cell 還原空 + toast | toast + 復原建議 |
| 測試連線 | ContextBar spinner | 狀態燈 → 黃色閃爍 | 狀態燈 → 綠 + 時間戳更新 | 狀態燈 → 紅 + error tooltip |
| 儲存映射 | button spinner | — | row 綠 flash + toast | row 紅 + inline error |
| 匯入 | modal progress | 逐項驗證進度 | summary: 成功/跳過/失敗 計數 | 逐項錯誤清單 |
| 匯出 | — | — | 檔案下載 + toast | toast |

### 5.3 批次操作結果顯示

**原則**: 每次批次操作後，使用者必須清楚知道「哪些成功、哪些失敗、為什麼」

**結果 toast 格式**:
```
✓ 批次建立完成: 18/20 成功
  ├── 2 項失敗: D105 (地址衝突), D112 (型態不相容)
  └── [查看詳情] [重試失敗項]
```

**Inspector 批次結果面板**:
- 成功列表: 收合式，預設收合
- 失敗列表: 預設展開，每項含失敗原因 + 建議操作
- 持續顯示直到使用者關閉或開始新操作

### 5.4 空狀態設計

每個主區域都需要有意義的空狀態:

| 區域 | 空狀態內容 |
|------|-----------|
| Device (未選) | 插畫 (Server icon, muted) + 「選擇一個設備開始您的 Datalink 工作流程」+ `[選擇設備]` button |
| Source Grid (無 Point) | Grid 顯示空白地址格 + 頂部提示 bar: 「使用上方工具列規劃地址空間，或匯入現有配置」 |
| Tag Binding (無 unbound points) | 「所有點位都已綁定 Tag ✓」成功狀態 或 「尚未建立點位，請先完成 Source Grid 步驟」 |
| Output / Modbus (無 mapping) | 「將已綁定的 Tag 映射到 Modbus Register 以開始資料輸出」+ `[自動映射]` button |
| Inspector (無選取) | 「點選 Grid 上的地址格或表格中的行以檢視詳情」+ 鍵盤提示 |
| Table (無資料) | 表格 header 仍顯示 + 單行訊息「尚無點位。切換到 Grid 模式建立」 |

### 5.5 無障礙 (a11y)

| 項目 | 規格 |
|------|------|
| 鍵盤導航 | Grid: arrow keys 移動焦點、Space 選取、Enter 開啟 context menu、Tab 離開 Grid |
| ARIA roles | Grid: `role="grid"`, row: `role="row"`, cell: `role="gridcell"` + `aria-selected` |
| 螢幕閱讀器 | 每個 cell announce: 「地址 D100, 已建立點位 temperature, 已綁定 Tag」 |
| 焦點可見 | focus ring: `focus-visible:ring-2 focus-visible:ring-blue-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900` |
| 色彩對比 | 所有文字-背景組合 WCAG AA (4.5:1 最小)，重要操作 AAA |
| 動畫減少 | 偵測 `prefers-reduced-motion` → 停用 pulse、flash、spin 動畫 |
| Live region | BottomSummaryBar: `aria-live="polite"` 計數變化時 announce |
| StepRail | `role="tablist"`, 每個 step: `role="tab"` + `aria-selected` + `aria-controls` |

### 5.6 快捷鍵（延伸現有 useKeyboardShortcuts）

| 組合鍵 | 動作 | 適用 Step |
|--------|------|----------|
| `Ctrl+B` | 批次建立 Points | Step 2 |
| `Ctrl+T` | 批次綁定 Tag | Step 2, 3 |
| `Ctrl+M` | 新增 Modbus 映射 | Step 4 |
| `Ctrl+G` | 切換 Grid/Table 模式 | Step 2 |
| `Ctrl+I` | 匯入 | Step 2, 4 |
| `Ctrl+E` | 匯出 | Step 2, 4 |
| `Ctrl+Enter` | 執行主要操作 (依 context) | 全域 |
| `Escape` | 關閉 Inspector / 取消選取 | 全域 |
| `Ctrl+A` | 全選 (Grid 或 Table) | Step 2 |
| `Delete` | 刪除選取項 (需確認) | Step 2, 3, 4 |
| `1`–`4` | 切換步驟 (非輸入框時) | 全域 |
| `Ctrl+[` | 收合/展開 Inspector | 全域 |
| `?` | 顯示快捷鍵面板 | 全域 |

---

## 6. Domain 沿用邊界與全新 UI 建構

### 6.1 直接沿用（非 UI 層）

| 模組/能力 | 路徑參考 | 沿用方式 |
|----------|---------|---------|
| `hooks/datalink/*` (React Query hooks) | `hooks/datalink/` | 直接 import，不動 |
| `services/datalink.ts` | `services/` | 直接 import，不動 |
| `types/datalink.ts` | `types/` | 直接 import，必要時以新型別擴充 |
| `useToast` | `hooks/useToast.ts` | 沿用通知機制 |
| `useKeyboardShortcuts` | `hooks/useKeyboardShortcuts.ts` | 沿用快捷鍵註冊能力，workbench 自己定義 key map |
| i18n 字典結構 | `i18n/locales/` | 新增 `workbench` namespace 與新文案 |

### 6.2 只能參考邏輯，不直接沿用 UI

| 來源 | 可參考內容 | 不可沿用內容 |
|------|-----------|-------------|
| `LocalModbusWorkbenchPage` | server 控制流程、mapping CRUD、衝突檢查、匯入匯出邏輯 | page shell、卡片版型、表單 UI |
| `useSmartDashboardTagLinking` | tag 綁定流程與 domain 規則 | hook API 與 SmartDashboard 假設 |
| `useSmartDashboardCommitFlow` | validation / readiness / chunking 思路 | 任何頁面耦合狀態與文案 |
| `SmartDashboardWorkspaceSection` | source planner 與選取流程概念 | workspace 版型、tabs、overlay 行為 |

### 6.3 必須全新建立

| 元件 / 區塊 | 說明 |
|------------|------|
| `DatalinkWorkbenchPage` | 新路由入口，全新 page shell |
| `WorkbenchProvider` | 新 workflow context / state |
| `WorkbenchStepNavigator` | 全新 step rail / responsive step nav |
| `WorkbenchHeaderBar` | 全新裝置/狀態/快速操作區 |
| `AddressCanvas` | 全新地址視覺化元件（取代舊 grid） |
| `AddressLedger` | 全新 table companion |
| `SelectionInspector` | 全新 inspector |
| `TagBindingStudio` | 全新批次 tag 工作面 |
| `OutputStudio` | 全新輸出步驟容器 |
| `LocalModbusBoard` | 全新 Local Modbus 體驗 |
| `WorkbenchActionDock` | 全新底部摘要/驗證區 |

---

## 7. Opus UI 子代理切法

### 7.1 建議切為 6 個獨立 UI Slices

每個 slice 應獨立可開發、可測試、無跨 slice 依賴（透過 WorkbenchProvider context 解耦）。

---

#### Slice 0: Foundation (必須最先完成)

**範圍**:
- `DatalinkWorkbenchPage` (route + page shell + 外層容器)
- `WorkbenchProvider` (context: step state, device selection, grid selection state, draft state)
- `WorkbenchStepNavigator`
- 路由註冊 (`/datalink/workbench`)
- i18n 新增 workbench namespace 骨架

**交付物**: 可進入頁面、可切換 4 個 step（內容為 placeholder）、StepRail 正常運作

**測試重點**: step 切換、provider state 持久化、路由存取

---

#### Slice 1: HeaderBar + Device Step

**範圍**:
- `WorkbenchHeaderBar`
- Device Step 的設備選擇 UI (list/cards + 搜尋/篩選)
- 設備 Inspector (Step 1 版本)
- 連線測試互動

**依賴**: Slice 0 (WorkbenchProvider)

**交付物**: 可選擇設備、ContextBar 顯示設備資訊、測試連線

**測試重點**: 設備選擇 → provider 更新、連線測試回饋、空狀態

---

#### Slice 2: Source Experience (Canvas + Ledger + Toolbar)

**範圍**:
- `SourceCanvasSection`（含 toolbar、mode switch，內含全新的 `AddressCanvas` 與 `AddressLedger`）
- `AddressCanvas`
- `AddressLedger`
- Source planner 控制 (起始地址、型態、數量)
- Grid/Table 選取同步

**依賴**: Slice 0, Slice 1 (需要 selected device)

**交付物**: Grid 顯示地址空間、可規劃/建立 Point、可切換 Table 模式

**測試重點**: Cell 狀態渲染、選取模式 (single/multi/range)、模式切換保留狀態、批次建立

---

#### Slice 3: SelectionInspector

**範圍**:
- `SelectionInspector` (通用框架 + Step 2/3/4 內容切換)
- 單選/多選 Inspector 內容
- 收合/展開互動
- 窄螢幕 drawer 模式（全新實作）

**依賴**: Slice 0, Slice 2 (需要 selection state)

**交付物**: 選取 Grid cell 或 Table row → Inspector 顯示詳情、可 inline 編輯

**測試重點**: 選取→顯示同步、多選摘要、收合/展開、響應式

---

#### Slice 4: TagBindingStudio (Step 3)

**範圍**:
- `TagBindingStudio` 完整 UI
- 命名模板系統
- 衝突檢查
- 批次綁定操作 + 結果顯示
- Inspector Step 3 內容

**依賴**: Slice 0, Slice 2 (需要 points)

**交付物**: 可批次選取 Points → 批次產生 Tag keys → 綁定 → 結果顯示

**測試重點**: 命名模板預覽、衝突偵測、批次結果 (partial failure)

---

#### Slice 5: OutputStudio + ActionDock (Step 4)

**範圍**:
- `OutputStudio`
- `LocalModbusBoard`
- `DatabaseTargetCard` (placeholder)
- `WorkbenchActionDock` (全域 summary)
- Inspector Step 4 內容

**依賴**: Slice 0, Slice 4 (需要 bound tags)

**交付物**: 完整 Local Modbus 映射工作流 + 底部驗證摘要

**測試重點**: Server 控制、映射 CRUD、衝突偵測、匯入/匯出、BottomBar 計數正確性

---

### 7.2 Slice 依賴圖

```
Slice 0 (Foundation)
  ├── Slice 1 (HeaderBar + Device)
  ├── Slice 2 (Source Experience) ← depends on Slice 1
  │     ├── Slice 3 (SelectionInspector) ← depends on Slice 2
  │     └── Slice 4 (TagBindingStudio) ← depends on Slice 2
  │           └── Slice 5 (OutputStudio + ActionDock) ← depends on Slice 4
```

**平行可行性**:
- Slice 1 + Slice 2 可平行（各自 mock device selection / grid data）
- Slice 3 + Slice 4 可平行（各自 mock selection state / points）
- Slice 5 需等 Slice 4 完成（tag binding 是前置）

### 7.3 每個 Slice 的估計 UI 元件數

| Slice | 新元件 | 修改元件 | 新 hooks | i18n 新 keys |
|-------|--------|---------|----------|-------------|
| 0 | 3 | 0 | 1 (provider) | ~15 |
| 1 | 2 | 0 | 0 (沿用 domain hooks) | ~25 |
| 2 | 4 | 0 | 1 (useSourcePlanner) | ~45 |
| 3 | 1 | 0 | 0 | ~20 |
| 4 | 2 | 0 | 1 (useTagBinding or equivalent new domain hook) | ~35 |
| 5 | 3 | 0 | 1 (useLocalModbusTarget or equivalent new domain hook) | ~45 |
| **合計** | **15** | **0** | **4** | **~185** |

---

## 附錄 A: 色彩速查

| 語義 | Tailwind class | 用途 |
|------|---------------|------|
| 互動/進行中 | `blue-500/600` | 選取、主按鈕、進度 |
| 有效/健康 | `emerald-500/600` | 連線成功、已建立、已綁定 |
| 警告/衝突 | `amber-500/600` | 部分完成、衝突可解決 |
| 阻塞/失敗 | `rose-500/600` | 嚴重衝突、API 錯誤、連線失敗 |
| 已綁定 Tag | `violet-500/600` | Tag 關聯視覺 |
| 規劃中 | `sky-500/600` | 尚未確認的規劃 |
| 中性/禁用 | `slate-500/600/700` | 背景、邊框、禁用 |

## 附錄 B: 設計決策紀錄

1. **StepRail 非線性**: 選擇非 wizard-style 線性鎖定，因為進階使用者需要在 step 間自由跳轉。空 state 引導取代硬鎖。
2. **Inspector 內嵌而非 modal**: 減少 context switch，保持 Grid 可見。窄螢幕退化為 drawer 而非全螢幕 modal。
3. **BottomBar 常駐**: 驗證摘要必須隨時可見，不依賴特定 step 或 scroll 位置。
4. **DatabaseTargetCard 佔位而非隱藏**: 表達完整 IA 意圖，使用者知道功能即將到來，避免 Release 2 時 UI 結構性調整。
5. **Grid 10 欄**: 對齊 Modbus 10 進位地址習慣 (D0–D9, D10–D19…)，工程師直覺。
6. **暗色主題優先**: 工控場景常在控制室使用，暗色減少眩光且突出狀態色彩。
