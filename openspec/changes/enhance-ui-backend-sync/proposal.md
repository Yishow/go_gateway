# Proposal: enhance-ui-backend-sync

## Summary

導入 **TanStack Query (React Query)** 及 **Shared TypeScript Types** 機制，提升 Datalink UI 與後端 API 的連結品質與開發者體驗。

## Problem Statement

目前 Datalink UI 使用 `useState` + `useEffect` 進行 API 呼叫，存在以下問題：

1. **冗餘的樣板代碼**：每個頁面都需重複處理 `loading`、`error`、`refetch` 邏輯。
2. **無緩存機制**：頁面切換後數據重新載入，使用者體驗不流暢。
3. **無自動重試**：網路閃斷時直接顯示錯誤，需手動重新整理。
4. **型別同步風險**：前端 TypeScript 與後端 Go Struct 手動維護，可能不一致。

## Proposed Solution

1. **導入 TanStack Query v5**：管理所有 CRUD 資源的 Server State。
2. **建立 Custom Hooks**：為 Devices、Points、Tags、Mappings 建立專用 Hooks。
3. **Shared Types 機制**：保持現有手動維護的 TypeScript types。

## Scope

- **In Scope**:
  - 為 `DevicesPage.tsx` 示範改寫
  - 建立可複用的 Query Hooks 結構
  - 加入 Optimistic Updates 範例

- **Out of Scope**:
  - WebSocket 即時推送（留待後續）
  - 自動化 Go -> TypeScript 型別產生（留待後續）

## Related Specs

- `datalink-ui`: 需新增 Query 狀態管理需求
- `datalink-api`: 無修改（純前端變更）

## Risks & Mitigations

| Risk        | Mitigation                        |
| ----------- | --------------------------------- |
| 學習曲線    | 提供完整範例代碼與註解            |
| Bundle 增大 | TanStack Query ~12KB gzip，可接受 |
