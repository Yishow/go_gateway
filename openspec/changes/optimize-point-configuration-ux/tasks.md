# Tasks: 智慧 Dashboard 模式 - TDD 任務清單

## 概要

採用 TDD（Test-Driven Development）工作流程，先定義測試再實作程式碼。

---

## Phase 1: 核心框架（P0）🎯

### 1.0 UI 一致性與設計系統統一

- [x] **建立設計 Token**
  - [x] 建立 `frontend/src/styles/tokens.ts`
  - [x] 定義顏色、間距、圓角、陰影、動畫 Token
  - [x] 更新 `tailwind.config.js` 引用 Token

- [x] **TestPage UI 一致性**
  - [x] 審查 TestPage 樣式與 Datalink 頁面差異
  - [x] 統一背景色、卡片樣式、按鈕樣式
  - [x] 確保圖標風格一致（Lucide Icons）

### 1.1 設備樹拖曳排序（TDD）

- [x] **安裝依賴**
  - [x] `npm install @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities`

- [x] **RED**: 撰寫拖曳排序測試
  - [x] 測試拖曳開始/結束事件
  - [x] 測試順序更新
  - [x] 測試鍵盤無障礙支援
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作拖曳排序
  - [x] 更新 `DeviceTreeNav.tsx` 整合 dnd-kit
  - [x] 實作 `SortableDeviceItem` 組件
  - [x] 後端新增 `PATCH /devices/reorder` API
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化拖曳體驗
  - [x] 調整拖曳視覺反饋
  - [x] 確保測試仍通過

### 1.2 設計系統與 Token

### 1.2 位址解析器（TDD）

- [x] **RED**: 撰寫 AddressParser 單元測試
  - [x] 測試 Modbus 位址解析 (`40001` → `{area: 'HR', start: 1}`)
  - [x] 測試 FATEK 位址解析 (`D0100` → `{area: 'D', start: 100}`)
  - [x] 測試 MC3E 位址解析 (`D100` → `{area: 'D', start: 100}`)
  - [x] 測試位址範圍展開 (`expand('D100', 10)` → 10 個位址)
  - [x] 測試位址驗證 (`validate('D100')` → `{valid: true}`)
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 AddressParser
  - [x] 建立 `frontend/src/utils/addressParser.ts`
  - [x] 實作各協議解析邏輯
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化程式碼
  - [x] 確保測試仍通過

### 1.3 設備樹狀導覽（TDD）

- [x] **RED**: 撰寫 DeviceTreeNav 組件測試
  - [x] 測試設備列表渲染
  - [x] 測試展開/收合子項目
  - [x] 測試選中狀態
  - [x] 測試收合模式
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 DeviceTreeNav
  - [x] 建立 `frontend/src/components/datalink/DeviceTreeNav.tsx`
  - [x] 實作樹狀結構渲染
  - [x] 實作展開/收合邏輯
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化組件結構
  - [x] 提取 TreeNode 子組件
  - [x] 確保測試仍通過

### 1.4 記憶體格子視覺化（TDD）

- [x] **RED**: 撰寫 MemoryGrid 組件測試
  - [x] 測試格子數量渲染（100 個）
  - [x] 測試已使用/可用狀態顯示
  - [x] 測試單擊選取
  - [x] 測試 Shift+Click 範圍選取
  - [x] 測試 Tooltip 顯示
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 MemoryGrid
  - [x] 建立 `frontend/src/components/datalink/MemoryGrid.tsx`
  - [x] 實作格子渲染
  - [x] 實作選取邏輯
  - [x] 實作 Tooltip
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化效能
  - [x] 使用 React.memo 避免不必要的重渲染
  - [x] 確保測試仍通過

### 1.5 快速操作面板（TDD）

- [x] **RED**: 撰寫 QuickActions 組件測試
  - [x] 測試按鈕渲染
  - [x] 測試點擊事件
  - [x] 測試禁用狀態（無選中設備）
  - [x] 測試即時狀態顯示
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 QuickActions
  - [x] 建立 `frontend/src/components/datalink/QuickActions.tsx`
  - [x] 整合 Lucide Icons
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化視覺效果
  - [x] 確保測試仍通過

### 1.6 側邊滑出面板（TDD）

- [x] **RED**: 撰寫 SlidePanel 組件測試
  - [x] 測試開啟/關閉動畫
  - [x] 測試背景遮罩
  - [x] 測試 ESC 鍵關閉
  - [x] 測試內容渲染
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 SlidePanel
  - [x] 建立 `frontend/src/components/datalink/SlidePanel.tsx`
  - [x] 實作滑動動畫
  - [x] 實作鍵盤事件
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化動畫效能
  - [x] 使用 CSS transform
  - [x] 確保測試仍通過

### 1.7 SmartDashboard 主頁面

- [x] **組合各組件**
  - [x] 建立 `frontend/src/pages/datalink/SmartDashboard.tsx`
  - [x] 實作三欄佈局
  - [x] 整合 DeviceTreeNav
  - [x] 整合 MemoryGrid
  - [x] 整合 QuickActions
  - [x] 整合 SlidePanel
  - [x] 建立 SmartDashboardContext

- [x] **路由配置**
  - [x] 更新 `App.tsx` 路由
  - [x] SmartDashboard 設為預設首頁
  - [x] 保留舊頁面（過渡期）

---

## Phase 2: 功能完善（P1）

### 2.1 批量建立器（TDD）

- [x] **RED**: 撰寫 BatchPointCreator 組件測試
- [x] **GREEN**: 實作 BatchPointCreator
- [x] **REFACTOR**: 優化使用者體驗

### 2.2 後端批量建立 API（TDD）

- [x] **RED**: 撰寫 API 測試
- [x] **GREEN**: 實作 API
- [x] **REFACTOR**: 優化效能

### 2.3 點位詳情面板

- [x] 實作 PointDetailPanel 組件
- [x] 顯示點位資訊
- [x] 編輯點位屬性
- [x] 顯示關聯映射
- [x] 新增/編輯映射

### 2.4 快速映射功能

- [x] 實作 QuickMappingPanel 組件
- [x] 選取點位 → 選取/建立標籤
- [x] 設定轉換規則
- [x] 預覽轉換結果

---

## Phase 3: 遷移與優化（P2）

### 3.1 舊頁面遷移

- [ ] 移除 PointsPage（功能已整合）
- [ ] 移除 MappingsPage（功能已整合）
- [ ] 移除 MappingWizardPage（功能已整合）
- [ ] 移除 DeviceOnboardingPage（功能已整合）
- [ ] 更新導航結構

### 3.2 響應式優化

- [ ] 平板佈局（1024px - 1280px）
  - [ ] 隱藏 QuickActions，改為浮動按鈕
  - [ ] TreeNav 預設收合
- [ ] 大螢幕優化（1920px+）
  - [ ] 加大格子尺寸
  - [ ] 顯示更多資訊

### 3.3 進階功能

- [ ] 鍵盤快捷鍵
  - [ ] `Ctrl+B` 批量建立
  - [ ] `Ctrl+M` 快速映射
  - [ ] `Escape` 關閉面板
- [ ] 匯入/匯出功能
  - [ ] CSV 匯入點位
  - [ ] JSON 匯出設定
- [ ] 撤銷/重做
  - [ ] 記錄操作歷史
  - [ ] `Ctrl+Z` / `Ctrl+Y`

---

## 驗證檢查清單

### 每個 Phase 完成前

- [ ] 所有測試通過 (`npm test`)
- [ ] TypeScript 無錯誤 (`npm run typecheck`)
- [ ] ESLint 無警告 (`npm run lint`)
- [ ] 專案可建置 (`npm run build`)

### Phase 1 完成驗證

- [ ] SmartDashboard 頁面可正常載入
- [ ] 設備樹狀導覽正常運作
- [ ] 記憶體格子正確顯示
- [ ] 側邊面板可開啟/關閉
- [ ] 響應式佈局正常（1280px, 1440px, 1920px）

### Phase 2 完成驗證

- [x] 批量建立功能正常
- [ ] 點位詳情面板正常
- [ ] 快速映射功能正常
- [ ] 後端 API 正常

### Phase 3 完成驗證

- [ ] 舊頁面已移除
- [ ] 無功能迴歸
- [ ] 鍵盤快捷鍵正常
- [ ] 匯入/匯出正常
