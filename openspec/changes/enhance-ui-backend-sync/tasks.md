# Tasks: enhance-ui-backend-sync

## Phase 1: 基礎建設

- [ ] **1.1** 安裝依賴：`npm install @tanstack/react-query`
- [ ] **1.2** 建立 `src/lib/queryClient.ts`，配置 QueryClient 預設選項
- [ ] **1.3** 在 `App.tsx` 或 `main.tsx` 加入 `<QueryClientProvider>`

## Phase 2: Query Keys & Hooks 建立

- [ ] **2.1** 建立 `src/hooks/datalink/keys.ts`，定義所有 Query Key Factories
- [ ] **2.2** 建立 `src/hooks/datalink/useDevices.ts`：
  - `useDevicesQuery()`: 列表查詢
  - `useDeviceQuery(id)`: 單一查詢
  - `useCreateDeviceMutation()`: 建立
  - `useUpdateDeviceMutation()`: 更新
  - `useDeleteDeviceMutation()`: 刪除
  - `useTestConnectionMutation()`: 測試連線
  - `useToggleDeviceStatusMutation()`: 切換狀態（含 Optimistic Update）
- [ ] **2.3** 建立 `src/hooks/datalink/index.ts`，統一匯出

## Phase 3: 頁面改寫示範

- [ ] **3.1** 改寫 `DevicesPage.tsx`，使用新的 Query Hooks
- [ ] **3.2** 移除舊的 `useState`/`useEffect` 樣板代碼
- [ ] **3.3** 驗證功能：列表、建立、編輯、刪除、測試連線

## Phase 4: 驗證與文件

- [ ] **4.1** 執行前端開發伺服器驗證 (`npm run dev`)
- [ ] **4.2** 確認 Network tab 中 Cache 行為正確
- [ ] **4.3** 更新 `openspec/specs/datalink-ui/spec.md`，新增 Query 狀態管理需求

## Dependencies

- Task 1.x 必須完成後才能進行 Task 2.x
- Task 2.x 必須完成後才能進行 Task 3.x
