# Go-Gateway Workbench — UI Design Principles

> **基礎設計系統**：以 [Linear Design System](./LINEAR-DESIGN-REFERENCE.md) 為設計原語參考，在此基礎上針對工業數據採集閘道的操作員情境做適配。
>
> **這份文件是所有新頁面設計的唯一設計規範來源**，包含新版 Workbench（4步驟主流程）、Local Modbus 工具頁與所有新增元件。

---

## 核心設計哲學

> **「工業精準感」**：深色畫布上，每個元素都有其計算過的層次；操作員在高壓環境下能快速識別狀態、執行批次操作、診斷問題。

三個設計原則：
1. **暗色原生（Dark-Mode Native）**：深色背景是預設，資訊從黑暗中浮現
2. **密度可控（Controlled Density）**：資訊密集但層次清晰，不是雜亂
3. **狀態可讀（Status Legible）**：任何就緒/錯誤/進行中狀態，操作員一眼可判斷

---

## 顏色系統

### 背景層次

| Token | 值 | 用途 |
|-------|-----|------|
| `bg-void` | `#08090a` | 最深背景（頁面底層） |
| `bg-base` | `#0f1011` | 主要頁面背景 / 側欄 |
| `bg-elevated` | `#191a1b` | 卡片 / 面板背景 |
| `bg-surface` | `#28282c` | hover state / 輕度 elevated |

→ Tailwind 對應：`slate-950`, `slate-900`, `slate-800/90`, `slate-800/50`

### 文字層次

| Token | 值 | 用途 |
|-------|-----|------|
| `text-primary` | `#f7f8f8` | 主標題、重要數值 |
| `text-secondary` | `#d0d6e0` | 說明文字、次要資訊 |
| `text-muted` | `#8a8f98` | 佔位符、metadata |
| `text-faint` | `#62666d` | 時間戳、停用狀態 |

→ Tailwind 對應：`slate-50`, `slate-200/300`, `slate-400`, `slate-500`

### 主色調（品牌 Accent）

本專案延用現有的 **Cyan** 作為主操作色（取代 Linear 的 Indigo）：

| Token | 值 | 用途 |
|-------|-----|------|
| `accent-primary` | `#06b6d4` （cyan-500） | 主要按鈕、active state |
| `accent-hover` | `#67e8f9` （cyan-300） | hover、focus ring |
| `accent-muted` | `rgba(6,182,212,0.15)` | 選中底色、背景 wash |

### 狀態顏色系統

這是本專案最關鍵的顏色層：

| 狀態 | 色點 | 背景 wash | 邊框 | 語意 |
|------|------|----------|------|------|
| `ready` / `applied` | `#34d399` (emerald-400) | `emerald-500/10` | `emerald-500/40` | 正常運作 |
| `partial` | `#fbbf24` (amber-400) | `amber-500/10` | `amber-500/40` | 部分完成，需注意 |
| `blocked` | `#f87171` (rose-400) | `rose-500/10` | `rose-500/40` | 阻擋，需修復 |
| `draft` | `#475569` (slate-600) | `slate-700/30` | `slate-700/50` | 草稿，未啟用 |
| `active` / `running` | `#38bdf8` (sky-400) | `sky-500/10` | `sky-500/30` | 進行中 |

### 邊框系統

```css
/* 主要邊框 */
border: 1px solid rgba(255,255,255,0.06);

/* 強調邊框（卡片、面板） */
border: 1px solid rgba(255,255,255,0.08);

/* Active/選中邊框 */
border: 1px solid rgba(6,182,212,0.45);  /* cyan-500/45 */

/* 危險邊框 */
border: 1px solid rgba(248,113,113,0.40); /* rose-400/40 */
```

---

## 字體系統

**字體家族**：`Inter Variable`（已在現有 Tailwind 設定中）
- OpenType features：`"cv01", "ss03"` — 在 `index.css` 全域啟用
- 備用：`-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif`

**技術 Monospace**（用於地址、數值、代碼）：`JetBrains Mono, ui-monospace, 'SF Mono', Menlo`

### 字體層次

| 角色 | 大小 | 字重 | 間距 | 用途 |
|------|------|------|------|------|
| Display | 32px (2rem) | 510 | -0.6px | 頁面主標題 |
| Heading 1 | 24px (1.5rem) | 600 | -0.4px | 區塊標題 |
| Heading 2 | 18px (1.125rem) | 600 | -0.2px | 小節標題 |
| Body Default | 14px (0.875rem) | 400 | 0 | 一般說明 |
| Body Small | 13px (0.8125rem) | 400 | 0 | 次要內容 |
| Label | 11px (0.6875rem) | 600 | +0.18em | 全大寫標籤 |
| Mono | 13px | 400 | 0 | 地址、數值、代碼 |

---

## 間距系統

**基礎單位**：4px

```
4px   → micro gap（圖示與文字間距）
8px   → compact padding（小元件內間距）
12px  → comfortable gap（元件間距）
16px  → section gap（區塊內間距）
20px  → panel padding（面板內間距）
24px  → major gap（主要區域分隔）
32px  → section break（大區塊分隔）
```

---

## 元件設計規範

### 按鈕

**Primary Button（主要操作）：**
```css
background: linear-gradient(to bottom, cyan-500/15, cyan-500/8);
border: 1px solid cyan-500/40;
color: cyan-50;
border-radius: 12px;
padding: 8px 20px;
font-weight: 600;
font-size: 12px;

/* hover */
background: linear-gradient(to bottom, cyan-500/20, cyan-500/10);
border-color: cyan-400/55;
```

**Ghost Button（輔助操作）：**
```css
background: slate-950/60;
border: 1px solid slate-700;
color: slate-200;
border-radius: 12px;
padding: 8px 12px;

/* hover */
border-color: slate-500;
color: slate-50;
```

**Danger Button（危險操作）：**
```css
background: rose-500/10;
border: 1px solid rose-500/40;
color: rose-200;

/* hover */
background: rose-500/15;
border-color: rose-400/55;
```

### 輸入框

```css
background: slate-950/70;
border: 1px solid slate-700;
border-radius: 12px;
padding: 8px 12px;
color: slate-100;
font-size: 14px;

/* focus */
border-color: cyan-400;
ring: 2px ring cyan-400/30;

/* error */
border-color: rose-500;
ring: 2px ring rose-400/30;
```

### 狀態色點（Readiness Dot）

```css
/* ready */
.dot { background: emerald-400; border-radius: 50%; width: 8px; height: 8px; }
/* partial */
.dot { background: amber-400; }
/* blocked */
.dot { background: rose-400; }
/* draft */
.dot { background: slate-600; }
```

### 卡片 / 面板

```css
background: slate-900/80;
border: 1px solid rgba(255,255,255,0.06);
border-radius: 16px;
box-shadow: 0 1px 0 0 rgba(255,255,255,0.04) inset;
```

### 全大寫標籤（Section Label）

```css
font-size: 10px (10.5px);
font-weight: 700;
letter-spacing: 0.18em to 0.22em;
color: slate-400 to slate-500;
text-transform: uppercase;
```

---

## 佈局系統

### Workbench 主框架（Desktop First, 1280px+）

```
┌─────────────────────────────────────────────────────────────┐
│  WorkbenchContextBar（頂欄）                    bg-base      │
│  [StepRail: 1→2→3→4] [設備名稱+狀態] [⚠ 問題] [主按鈕]   │
├─────────────────────────────────────────────────────────────┤
│                                          │                   │
│  PrimaryWorkArea                         │  InspectorPanel   │
│  （各步驟主要內容）                       │  （選中項目詳情）  │
│                              flex-1      │  w-80 ~ w-96      │
│                                          │                   │
├─────────────────────────────────────────────────────────────┤
│  WorkbenchBottomSummaryBar（底欄）                           │
│  [Points: N] [Tags: N] [Outputs: N] ● Device ● Source ...  │
└─────────────────────────────────────────────────────────────┘
```

**工具頁（/tools/modbus）：**
```
┌─────────────────────────────────────────────────────────────┐
│  PageHeader（頁面標題 + Tab 列）                             │
├─────────────────────────────────────────────────────────────┤
│  ActionBar（工具列：AutoMap / Dry-run / 偵測衝突 / 篩選）   │
├─────────────────────────────────────────────────────────────┤
│                                          │                   │
│  MainContent（主要內容：Register 表格等）│  DetailPanel      │
│                                          │  （選中詳情）     │
└─────────────────────────────────────────────────────────────┘
```

---

## WorkbenchDiagnosticPanel（Slide-over）

```css
/* 整體 Slide-over 容器 */
position: fixed;
top: 0; right: 0;
width: 360px; /* sm: 400px */
height: 100vh;
background: #0f1011;
border-left: 1px solid rgba(255,255,255,0.08);
box-shadow: -8px 0 32px rgba(0,0,0,0.6);
z-index: 50;
/* animation: slide-in from right */
```

---

## 設計禁止事項（Anti-Patterns）

1. ❌ **不使用純白背景**（#ffffff）：所有頁面都是深色
2. ❌ **不使用純黑文字在深色背景**：最深文字用 `slate-50`
3. ❌ **狀態色點只用顏色傳達**：必須同時有文字標籤（無障礙）
4. ❌ **表單驗證只在送出後觸發**：必須有即時行內驗證
5. ❌ **按鈕文字超過 4 個中文字**：標籤保持簡短明確
6. ❌ **資訊密度無層次**：必須透過字重、顏色、間距建立清楚的視覺層次
7. ❌ **Toast 與 Inline notice 並存**：同一操作只用一種反饋方式
8. ❌ **元件超過 500 行不拆分**：超過 500 行必須考慮拆分子元件

---

## 設計驗收標準

每個新設計的頁面或元件，必須通過：

- [ ] **暗色一致性**：背景使用本文件的背景層次，不使用其他顏色
- [ ] **狀態可讀性**：所有狀態（ready/partial/blocked/draft）有色點 + 文字標籤
- [ ] **間距對齊**：使用 4px 基礎單位的倍數，不使用任意像素值
- [ ] **字體層次**：使用本文件定義的字體角色，不自訂任意 font-size
- [ ] **按鈕語意**：primary/ghost/danger 按鈕對應正確的操作語意
- [ ] **無障礙**：focus ring 可見（`focus-visible:ring-2 ring-cyan-400/40`），顏色對比 ≥ 4.5:1
- [ ] **響應式基準**：Desktop（1280px+）為主設計，sm（640px+）有最低可用性

---

## 參考原語

- **Linear Design System** → `frontend/docs/LINEAR-DESIGN-REFERENCE.md`
- **現有 Tailwind 設定** → `frontend/tailwind.config.js`（如有）
- **現有 Workbench 元件** → `frontend/src/pages/datalink/workbench/`（舊版，供對比用）
- **OpenSpec 變更** → `openspec/changes/workbench-ux-operator-efficiency/design.md`
