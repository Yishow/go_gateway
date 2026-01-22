# Change: 優化資料點配置流程 - 智慧 Dashboard 模式

## Why

當前資料點（Point）設定流程存在以下關鍵問題：

### 🔴 核心痛點

1. **頁面過多**：9 個獨立頁面（Dashboard, Devices, Points, Tags, Mappings 等），流程分散
2. **逐一輸入繁瑣**：設定 100 個連續點位需要手動輸入 100 次表單
3. **無視覺化回饋**：使用者無法直觀看到 PLC 記憶體區域的使用狀況
4. **上下文切換頻繁**：設備、點位、映射分散在不同頁面

### 💡 解決方案：智慧 Dashboard 模式（Option C）

採用類 IDE 的「一站式工作台」設計：

- **極簡頁面結構**（9 → 2 頁）
- **三欄佈局**：設備樹 | 動態內容區 | 快速操作面板
- **記憶體格子視覺化**：直觀呈現 PLC 暫存器使用狀況
- **側邊滑出面板**：無需跳頁即可完成詳細設定

---

## What Changes

### 頁面結構

```
舊結構 (9 頁)                  新結構 (2 頁)
─────────────────              ─────────────────
Dashboard                      SmartDashboard (一站式工作台)
Devices                           ├── 設備樹狀導覽
Points                            ├── 動態內容區 (記憶體格子)
Tags                              └── 快速操作面板 + 側邊抽屜
Mappings
PollingGroups                  Settings
MappingWizard
Settings
DeviceOnboarding
```

### 智慧 Dashboard 架構

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  🔍 Search...                                              [PLC-001 ▼]     │
├──────────────┬──────────────────────────────────────────┬──────────────────┤
│              │                                          │                  │
│  設備導覽     │         動態內容區                        │  快速操作        │
│              │                                          │                  │
│  ▼ PLC-001   │  ┌─────┬─────┬─────┬─────┬─────┬────┐   │  [+ 批量建立]    │
│    Points(5) │  │D100 │D101 │D102 │D103 │D104 │... │   │  [+ 快速映射]    │
│    Maps(3)   │  │ 🟢 │ 🟢 │ ⚪ │ ⚪ │ ⚪ │    │   │  [測試連線]      │
│  ▶ PLC-002   │  └─────┴─────┴─────┴─────┴─────┴────┘   │                  │
│  ▶ PLC-003   │                                          │  ─────────────   │
│              │  點選格子 → 開啟右側面板                    │  即時狀態        │
│              │                                          │  ● CPU: 12%     │
│              │                                          │  ● 記憶體: 45%   │
│              │                                          │  ● 連線: 正常    │
│              │                                          │                  │
└──────────────┴──────────────────────────────────────────┴──────────────────┘
                                                                    │
                                                                    ▼
                                                          ┌──────────────────┐
                                                          │ Point Details    │
                                                          │ 側邊滑出面板      │
                                                          │                  │
                                                          │ Name: ________   │
                                                          │ Address: D100    │
                                                          │ Type: int16 ▼    │
                                                          │                  │
                                                          │ [+ 新增映射]      │
                                                          │ ─────────────    │
                                                          │ Scale: 0.1       │
                                                          │ Offset: 0        │
                                                          │                  │
                                                          │ [取消] [儲存]    │
                                                          └──────────────────┘
```

---

## UI/UX 設計規格

### 設計系統

| 屬性         | 值                     |
| ------------ | ---------------------- |
| **主題**     | Dark Mode (工控標準)   |
| **主色**     | `#3B82F6` (Blue-500)   |
| **背景**     | `#0F172A` (Slate-900)  |
| **卡片**     | `#1E293B` (Slate-800)  |
| **邊框**     | `#334155` (Slate-700)  |
| **文字**     | `#F1F5F9` (Slate-100)  |
| **次要文字** | `#94A3B8` (Slate-400)  |
| **成功**     | `#22C55E` (Green-500)  |
| **警告**     | `#EAB308` (Yellow-500) |
| **錯誤**     | `#EF4444` (Red-500)    |

### 排版系統

| 元素        | 字體           | 大小 | 行高 |
| ----------- | -------------- | ---- | ---- |
| **標題 H1** | Inter Bold     | 24px | 1.2  |
| **標題 H2** | Inter SemiBold | 18px | 1.3  |
| **正文**    | Inter Regular  | 14px | 1.5  |
| **小字**    | Inter Regular  | 12px | 1.5  |
| **程式碼**  | JetBrains Mono | 13px | 1.4  |

### 元件規格

#### 1. 設備樹狀導覽 (`DeviceTreeNav`)

| 屬性       | 規格                                      |
| ---------- | ----------------------------------------- |
| 寬度       | 240px (可收合至 64px)                     |
| 項目高度   | 36px                                      |
| 縮排       | 16px per level                            |
| 圖標       | Lucide Icons 20x20                        |
| Hover 狀態 | `bg-slate-700/50`                         |
| 選中狀態   | `border-l-2 border-blue-500 bg-slate-700` |

#### 2. 記憶體格子視覺化 (`MemoryGrid`)

| 屬性     | 規格                   |
| -------- | ---------------------- |
| 格子大小 | 64x48px                |
| 間距     | 4px                    |
| 圓角     | 6px                    |
| 已使用   | `bg-green-500/80`      |
| 可用     | `bg-slate-700`         |
| 選中     | `ring-2 ring-blue-500` |
| 待建立   | `bg-yellow-500/60`     |
| Tooltip  | 顯示點位名稱、值、時間 |

#### 3. 快速操作面板 (`QuickActions`)

| 屬性     | 規格                              |
| -------- | --------------------------------- |
| 寬度     | 280px                             |
| 按鈕高度 | 40px                              |
| 按鈕間距 | 8px                               |
| 主按鈕   | `bg-blue-600 hover:bg-blue-500`   |
| 次按鈕   | `bg-slate-700 hover:bg-slate-600` |

#### 4. 側邊滑出面板 (`SlidePanel`)

| 屬性       | 規格                               |
| ---------- | ---------------------------------- |
| 寬度       | 400px                              |
| 動畫       | `slide-in-from-right duration-300` |
| 背景遮罩   | `bg-black/50 backdrop-blur-sm`     |
| 表單間距   | 16px                               |
| 輸入框高度 | 40px                               |

---

## Impact

### Affected Specs

- `datalink-ui/spec.md` - 完整重新設計 UI 架構

### Affected Code

**前端（新增）**：

- `frontend/src/pages/datalink/SmartDashboard.tsx` - 一站式工作台主頁
- `frontend/src/components/datalink/DeviceTreeNav.tsx` - 設備樹狀導覽
- `frontend/src/components/datalink/MemoryGrid.tsx` - 記憶體格子視覺化
- `frontend/src/components/datalink/QuickActions.tsx` - 快速操作面板
- `frontend/src/components/datalink/SlidePanel.tsx` - 側邊滑出面板
- `frontend/src/components/datalink/BatchPointCreator.tsx` - 批量建立器
- `frontend/src/utils/addressParser.ts` - 位址解析工具

**前端（移除/整合）**：

- ~~`PointsPage.tsx`~~ → 整合至 SmartDashboard
- ~~`TagsPage.tsx`~~ → 整合至側邊面板
- ~~`MappingsPage.tsx`~~ → 整合至側邊面板
- ~~`MappingWizardPage.tsx`~~ → 整合至側邊面板
- ~~`DeviceOnboardingPage.tsx`~~ → 整合至引導流程

**後端（新增）**：

- `POST /api/v1/datalink/points/batch` - 批量建立端點

---

## Risks

1. **開發工作量大**：完整重構前端 UI
   - 緩解：分階段實施，優先交付核心功能

2. **小螢幕適應性**：三欄佈局在小螢幕受限
   - 緩解：響應式設計，小螢幕改為抽屜式導覽

3. **使用者適應成本**：UI 變動較大
   - 緩解：提供使用教學，保留 Settings 頁面

---

## Migration Plan

### Phase 1（P0）：核心框架

1. 建立 `SmartDashboard` 三欄佈局
2. 實作 `DeviceTreeNav` 元件
3. 實作基礎 `MemoryGrid` 元件
4. 保留舊頁面，新舊並行

### Phase 2（P1）：功能完善

1. 實作 `SlidePanel` 側邊面板
2. 實作 `BatchPointCreator` 批量建立
3. 後端批量建立 API
4. 整合標籤映射設定

### Phase 3（P2）：遷移與優化

1. 移除舊頁面
2. 響應式優化
3. 鍵盤快捷鍵
4. 匯入/匯出功能
