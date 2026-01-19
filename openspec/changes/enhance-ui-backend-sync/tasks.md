# Tasks: enhance-ui-backend-sync

## Phase 1: 基礎建設

- [x] **1.1** 安裝依賴：`npm install @tanstack/react-query` _(已存在 v5.90.17)_
- [x] **1.2** 建立 `src/lib/queryClient.ts`，配置 QueryClient 預設選項
- [x] **1.3** 在 `main.tsx` 加入 `<QueryClientProvider>`

## Phase 2: Query Keys & Hooks 建立

- [x] **2.1** 建立 `src/hooks/datalink/keys.ts`，定義所有 Query Key Factories
- [x] **2.2** 建立 `src/hooks/datalink/useDevices.ts`：
  - `useDevicesQuery()`: 列表查詢
  - `useDeviceQuery(id)`: 單一查詢
  - `useCreateDeviceMutation()`: 建立
  - `useUpdateDeviceMutation()`: 更新
  - `useDeleteDeviceMutation()`: 刪除
  - `useTestConnectionMutation()`: 測試連線
  - `useToggleDeviceStatusMutation()`: 切換狀態（含 Optimistic Update）
- [x] **2.3** 建立 `src/hooks/datalink/index.ts`，統一匯出

## Phase 3: 頁面改寫示範

- [x] **3.1** 改寫 `DevicesPage.tsx`，使用新的 Query Hooks
- [x] **3.2** 移除舊的 `useState`/`useEffect` 樣板代碼
- [x] **3.3** 驗證功能：TypeScript 編譯通過

## Phase 4: 驗證與文件

- [x] **4.1** 執行 TypeScript 類型檢查通過 (`npx tsc --noEmit`)
- [x] **4.2** 確認 Cache 行為已配置（staleTime: 30s, gcTime: 5m）
- [x] **4.3** 更新 `openspec/specs/datalink-ui/spec.md`，新增 Query 狀態管理需求

## Dependencies

- Task 1.x 必須完成後才能進行 Task 2.x ✅
- Task 2.x 必須完成後才能進行 Task 3.x ✅
