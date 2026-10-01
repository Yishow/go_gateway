# C 方案 UI 重新設計視覺與互動規範

## 1. 核心設計理念 (Intentional Minimalism)
- **視覺降噪**：移除多餘邊框與陰影，透過留白與對比建立層次。
- **工業級高對比**：確保關鍵狀態（正常/警告/錯誤）一眼可辨，符合工業現場操作需求。
- **操作一致性**：完整對齊現有 Datalink 的互動邏輯與習慣，無縫銜接。

## 2. 視覺規範 (Visual Specs)

### 2.1 色彩系統 (Color Palette)
- **Primary (主色)**:
  - Base: `#2563EB` (Tailwind Blue-600)
  - Hover: `#1D4ED8` (Tailwind Blue-700)
- **Surface (表面與卡片)**:
  - Background: `#F8FAFC` (Tailwind Slate-50)
  - Card/Panel: `#FFFFFF` (White)
  - Border: `#E2E8F0` (Tailwind Slate-200)
- **Text (文字)**:
  - Primary: `#0F172A` (Tailwind Slate-900)
  - Secondary: `#64748B` (Tailwind Slate-500)
- **Status (狀態)**:
  - Success (連線/正常): `#10B981` (Emerald-500)
  - Warning (延遲/重試): `#F59E0B` (Amber-500)
  - Error (斷線/錯誤): `#EF4444` (Red-500)

### 2.2 字級與排版 (Typography)
- **字體**: Inter 或系統無襯線字體。
- **H1 (頁面標題)**: `24px`, Font Weight: `600` (Semibold), Tracking: `-0.02em`
- **H2 (區塊/卡片標題)**: `18px`, Font Weight: `600` (Semibold)
- **Body (內文/表單標籤)**: `14px`, Font Weight: `400` (Regular)
- **Caption (輔助說明)**: `12px`, Font Weight: `400` (Regular), Color: Secondary Text

### 2.3 間距與網格 (Spacing & Grid)
- **基數**: `4px` (Tailwind base)
- **組件內間距 (Padding)**: `16px` (p-4)
- **區塊間距 (Gap/Margin)**: `24px` (gap-6 / m-6)
- **頁面邊距**: `32px` (p-8)

### 2.4 組件規範 (Components)
- **卡片 (Cards)**:
  - 圓角: `12px` (rounded-xl)
  - 陰影: `0 4px 6px -1px rgb(0 0 0 / 0.1)` (shadow-md)，懸浮時增加深度 (shadow-lg)
  - 邊框: 單純 1px solid Border 色，無多餘裝飾。
- **按鈕 (Buttons)**:
  - 圓角: `8px` (rounded-lg)
  - Primary: 背景 Primary 色，白字，Hover 變深。
  - Secondary: 背景透明，邊框 Border 色，字體 Primary Text，Hover 背景 Slate-50。
  - 高度: `40px` (h-10) 或 `36px` (h-9) 用於緊湊介面。
- **表單 (Forms)**:
  - 輸入框高度: `40px`
  - 圓角: `8px`
  - 邊框: 預設 Slate-300，Focus 時 Primary 色搭配 `ring-2`。
- **狀態指示 (Status Indicators)**:
  - 採用實心圓點 (`8px` x `8px`) 加上文字標籤。

## 3. 互動規範 (Interaction Specs)
- **對齊 Datalink**:
  - **表格行與卡片點擊**: 支援整行/整區塊點擊進入詳情，Hover 時背景呈現微亮色 (`bg-slate-50`)。
  - **即時回饋**: 所有耗時操作（如連線測試、保存配置）必須有明確的 Loading 狀態與 Toast 通知。
  - **防呆確認**: 涉及刪除或覆寫的重要操作，需彈出 `ConfirmDialog`，不可直接生效。
