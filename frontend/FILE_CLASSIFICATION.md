# Frontend 檔案分類

更新時間：2026-08-24

> **目前 route contract**：`/studio/v2` 是唯一產品 setup 入口；`/studio/runtime` 是 setup 後 observer；`/test` 與 `/gateway/*` 保留。`/studio` dedicated surface 已移除，之後依 generic unknown-route policy 處理。本文件下方標示為 historical pre-delete 的 legacy 路徑，不代表現行產品檔案。

## 1) 路由與頁面分類

### 核心產品頁（現行）
- `src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`（`/studio/v2` setup）
- `src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx`（`/studio/runtime` observer）

### 相容／unknown-route policy
- `/studio`：沒有 dedicated route；依 generic unknown-route policy 收斂到 `/studio/v2`。
- `/datalink`、`/datalink/workbench`、`/datalink/workbench/*`、`/datalink/dashboard-legacy`：相容入口收斂到 V2 generic entry，不再指向 `/studio`。
- `/datalink/devices-legacy`：保留 legacy intent 相容路徑，不以 `/studio` 為目標。

### 協議測試工具頁（獨立於 Datalink 主流程）
- `src/pages/TestPage.tsx`
- `src/pages/TestPageShell.tsx`

### 現行頁面測試
- `tests/unit/workbench-v2/**`
- `tests/unit/runtime-dashboard/**`

### Historical pre-delete 檔案（已移除，非現行）
- `src/pages/datalink/workbench/DatalinkWorkbenchPage.tsx` 與 `src/pages/datalink/workbench/**`
- 舊 workbench 測試路徑 `tests/unit/pages/datalink/workbench/**`

## 2) UI 與元件分類

### Datalink 領域元件
- 目錄：`src/components/datalink/`
- 用途：MemoryGrid、QuickActions、DeviceForm、Import/Export、Wizard 等核心互動元件。

### 協議測試工具元件
- 檔案：`src/components/ProtocolSelector.tsx`、`src/components/ConfigForm.tsx`、`src/components/TestOperations.tsx`、`src/components/DeviceScanner.tsx`、`src/components/MonitorControl.tsx`、`src/components/RTUPollingCard.tsx` 等。

### 共用基礎元件
- 檔案：`src/components/Toast.tsx`、`src/components/ThemeToggle.tsx`、`src/components/ui/*`

## 3) 狀態與領域邏輯分類

### Feature 層（純邏輯/策略）
- 目錄：`src/features/datalink/`、`src/features/flow/`
- 用途：commit lifecycle、route intent、allocation strategy、motion guidance、state machine 等。

### Hooks 層
- 目錄：`src/hooks/datalink/`：API Query/Mutation hooks
- 目錄：`src/hooks/`：快捷鍵、歷史、profiles、preview stream

### Service 層
- 目錄：`src/services/`
- 用途：HTTP API 呼叫封裝（`api.ts`、`datalink.ts`）

### 型別與工具
- 型別：`src/types/`
- 工具：`src/utils/`
- 共用函式：`src/lib/`

## 4) 全域基礎設施分類

- 入口與路由：`src/main.tsx`、`src/App.tsx`
- i18n：`src/i18n/`
- Context：`src/contexts/`
- 樣式：`src/index.css`、`src/styles/tokens.ts`
- 測試初始化：`src/setupTests.ts`
- 建置設定：`vite.config.ts`、`tsconfig.json`、`tailwind.config.cjs`、`eslint.config.js`

## 5) 測試檔案分類

- 頁面互動測試：`tests/unit/pages/**`
- 元件測試：`tests/unit/components/**`
- 邏輯測試：`tests/unit/features/**`
- Hook 測試：`tests/unit/hooks/**`
- Utils 測試：`tests/unit/utils/**` 與 `*.test.ts`

## 6) 歷史去重已刪除檔案

以下檔案已確認無路由與無引用，已移除：

- `src/pages/datalink/DeviceOnboardingPage.tsx`
- `src/pages/datalink/PointsPage.tsx`
- `src/pages/datalink/PollingGroupsPage.tsx`
- `src/pages/datalink/TagsPage.tsx`
- `src/pages/datalink/MappingsPage.tsx`
- `src/pages/datalink/MappingWizardPage.tsx`
- `src/pages/datalink/SettingsPage.tsx`
- `src/pages/datalink/Dashboard.tsx`
- `src/pages/datalink/DevicesPage.tsx`
- `src/pages/datalink/SmartDashboard.tsx`
- `src/pages/datalink/SmartDashboardPage.tsx`
- `src/pages/datalink/LocalModbusWorkbenchPage.tsx`
- `src/pages/TemplatesPage.tsx`
- `src/pages/HistoryPage.tsx`
- `src/pages/ComparePage.tsx`
- `src/pages/AnalyzerPage.tsx`
- `src/components/Layout.tsx`
