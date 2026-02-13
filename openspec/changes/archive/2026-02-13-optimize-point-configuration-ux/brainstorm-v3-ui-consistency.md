# 🧠 Brainstorm: UI 一致性與進階功能優化

## Context

用戶要求進一步優化：

1. **新增拖曳排序功能**（設備樹節點）
2. **確保全系統 UI 一致性**
3. **TestPage 為獨立系統不可刪除**，但需保持 UI 風格一致

**現有架構**：

- `DatalinkLayout.tsx`：Sidebar + Header + Outlet（9 個 datalink 子頁面共用）
- `TestPage.tsx`：獨立頁面（20KB，協議測試工具），不使用 DatalinkLayout

---

## 優化空間分析

### 1️⃣ 設計系統一致性（Critical）

**問題**：TestPage 與 Datalink 頁面可能存在 UI 不一致

| 元素     | Datalink 頁面         | TestPage |
| -------- | --------------------- | -------- |
| 背景色   | `#0F172A` (Slate-900) | 待確認   |
| 卡片樣式 | 圓角 + 邊框 + 半透明  | 待確認   |
| 按鈕樣式 | 統一 Token            | 可能自訂 |
| 圖標     | Lucide (inline SVG)   | 待確認   |

### 2️⃣ 拖曳排序功能

**問題**：設備樹節點目前不支援拖曳排序

### 3️⃣ 共用組件提取

**問題**：可能存在重複的 UI 組件

---

## Option A: 統一設計 Token + 共用組件庫

建立統一的設計 Token 檔案，讓所有頁面引用同一套樣式。

```
frontend/src/
├── styles/
│   ├── tokens.ts        # 顏色、間距、字型 Token
│   ├── components.css   # 共用組件樣式
│   └── animations.css   # 共用動畫
├── components/
│   ├── ui/              # 統一 UI 組件
│   │   ├── Button.tsx
│   │   ├── Card.tsx
│   │   ├── Input.tsx
│   │   ├── Select.tsx
│   │   ├── SlidePanel.tsx
│   │   └── Tooltip.tsx
│   └── ...
```

✅ **Pros:**

- 所有頁面樣式一致
- 修改 Token 即可全局生效
- TestPage 可引用相同組件

❌ **Cons:**

- 需重構現有組件
- 增加維護成本

📊 **Effort:** Medium

---

## Option B: 統一 Layout 包裝器

建立通用 Layout 包裝器，TestPage 也使用類似框架。

```
frontend/src/layouts/
├── DatalinkLayout.tsx    # 現有 Datalink 佈局
├── TestLayout.tsx        # TestPage 專用佈局（相同風格）
└── BaseLayout.tsx        # 共用基礎佈局（Header + 基本樣式）
```

```
/* TestPage 使用 TestLayout */
<TestLayout>
  <TestPage />
</TestLayout>
```

✅ **Pros:**

- 保持 TestPage 獨立性
- Header、背景、字型統一
- 最小化修改範圍

❌ **Cons:**

- 需建立新 Layout
- 可能有重複程式碼

📊 **Effort:** Low

---

## Option C: 設備樹拖曳排序功能 ⭐ 新增

使用 `@dnd-kit/core` 實現設備樹節點拖曳排序。

```typescript
import { DndContext, closestCenter, DragEndEvent } from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
  useSortable,
} from "@dnd-kit/sortable";

interface SortableTreeItem {
  id: string;
  type: "device" | "point" | "mapping";
  order: number;
  parentId?: string;
}
```

**功能規格**：

- 設備節點可上下拖曳重新排序
- 拖曳時顯示半透明指示器
- 放下後自動儲存順序到後端
- 支援鍵盤排序（Accessibility）

✅ **Pros:**

- 提升使用者體驗
- 讓使用者自訂設備順序
- dnd-kit 輕量且無障礙友好

❌ **Cons:**

- 需新增依賴
- 需後端 API 儲存順序

📊 **Effort:** Medium

---

## Option D: 完整設計系統文件 + Storybook ⭐ 進階

建立完整的設計系統文件與 Storybook 組件庫。

```
frontend/
├── .storybook/
│   └── main.ts
├── src/
│   ├── components/
│   │   └── ui/
│   │       ├── Button/
│   │       │   ├── Button.tsx
│   │       │   ├── Button.stories.tsx
│   │       │   └── Button.test.tsx
│   │       └── ...
│   └── styles/
│       └── design-system.md
```

✅ **Pros:**

- 組件可視化文件
- 團隊協作友好
- 測試覆蓋完整

❌ **Cons:**

- 初始建置工作量大
- 需持續維護

📊 **Effort:** Very High

---

## 💡 Recommendation

### 短期（P0）：Option A + B + C

1. **統一設計 Token**（Option A）
   - 建立 `tokens.ts` 定義顏色、間距、動畫
   - 所有頁面引用統一 Token

2. **TestPage 佈局統一**（Option B）
   - 建立 `TestLayout.tsx` 或直接套用 Token
   - 確保背景、字型、邊框與 Datalink 一致

3. **設備樹拖曳排序**（Option C）
   - 安裝 `@dnd-kit/core` 和 `@dnd-kit/sortable`
   - DeviceTreeNav 支援拖曳排序
   - 後端新增 `PATCH /devices/:id/order` API

### 中期（P1）：進階優化

4. **共用組件庫**
   - 提取 Button、Card、Input、Select 等共用組件
   - 所有頁面引用統一組件

5. **動畫一致性**
   - 統一 transition 設定
   - 統一 hover 效果

### 長期（P2）：Option D

6. **Storybook**
   - 建立組件文件
   - 視覺回歸測試

---

## 統一設計 Token 定義

```typescript
// frontend/src/styles/tokens.ts

export const tokens = {
  colors: {
    // 背景
    bgPrimary: "#0F172A", // Slate-900
    bgSecondary: "#1E293B", // Slate-800
    bgTertiary: "#334155", // Slate-700
    bgHover: "#3F4F6A", // Slate-600/70

    // 文字
    textPrimary: "#F1F5F9", // Slate-100
    textSecondary: "#94A3B8", // Slate-400
    textMuted: "#64748B", // Slate-500

    // 強調色
    primary: "#3B82F6", // Blue-500
    primaryHover: "#2563EB", // Blue-600
    primaryLight: "#3B82F620", // Blue-500/12

    // 狀態
    success: "#22C55E",
    warning: "#EAB308",
    error: "#EF4444",

    // 邊框
    borderDefault: "#334155",
    borderHover: "#475569",
  },

  spacing: {
    xs: "4px",
    sm: "8px",
    md: "16px",
    lg: "24px",
    xl: "32px",
  },

  radius: {
    sm: "6px",
    md: "8px",
    lg: "12px",
    xl: "16px",
    full: "9999px",
  },

  shadow: {
    sm: "0 1px 2px rgba(0, 0, 0, 0.05)",
    md: "0 4px 6px -1px rgba(0, 0, 0, 0.1)",
    lg: "0 10px 15px -3px rgba(0, 0, 0, 0.1)",
    glow: {
      blue: "0 0 20px rgba(59, 130, 246, 0.3)",
      green: "0 0 20px rgba(34, 197, 94, 0.3)",
    },
  },

  transition: {
    fast: "150ms ease",
    normal: "200ms ease",
    slow: "300ms ease",
  },

  font: {
    family: {
      sans: "Inter, system-ui, sans-serif",
      mono: "JetBrains Mono, monospace",
    },
    size: {
      xs: "12px",
      sm: "14px",
      md: "16px",
      lg: "18px",
      xl: "24px",
    },
  },
} as const;
```

---

## 更新的任務清單

### Phase 1 新增項目

- [ ] **1.0 設計系統統一**
  - [ ] 建立 `frontend/src/styles/tokens.ts`
  - [ ] 更新 `tailwind.config.js` 引用 Token
  - [ ] 審查 TestPage 樣式一致性

- [ ] **1.X 設備樹拖曳排序**
  - [ ] 安裝 `@dnd-kit/core` `@dnd-kit/sortable`
  - [ ] 撰寫拖曳排序單元測試
  - [ ] 實作 `SortableDeviceTreeNav` 組件
  - [ ] 後端新增 `PATCH /devices/:id/order` API

---

**請問是否採用此優化方案？**
