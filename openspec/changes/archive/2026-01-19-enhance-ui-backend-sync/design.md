# Design: enhance-ui-backend-sync

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────┐
│                      React Components                        │
│  (DevicesPage, TagsPage, MappingsPage, ...)                  │
└───────────────────────────┬──────────────────────────────────┘
                            │ useDevices(), useTags(), ...
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                   Custom Query Hooks                         │
│  src/hooks/datalink/                                         │
│    ├── useDevices.ts                                         │
│    ├── useTags.ts                                            │
│    ├── usePoints.ts                                          │
│    └── useMappings.ts                                        │
└───────────────────────────┬──────────────────────────────────┘
                            │ queryFn, mutationFn
                            ▼
┌──────────────────────────────────────────────────────────────┐
│               Existing API Service Layer                     │
│  src/services/datalink.ts                                    │
│  (deviceAPI, tagAPI, pointAPI, mappingAPI, ...)              │
└───────────────────────────┬──────────────────────────────────┘
                            │ axios
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                    Go Backend API                            │
│  /api/v1/datalink/*                                          │
└──────────────────────────────────────────────────────────────┘
```

## Key Design Decisions

### 1. Query Key 策略

使用結構化的 Query Key 以支援精準的 Cache Invalidation：

```typescript
// Query Keys Factory
export const deviceKeys = {
  all: ["devices"] as const,
  lists: () => [...deviceKeys.all, "list"] as const,
  list: (filters: DeviceFilters) => [...deviceKeys.lists(), filters] as const,
  details: () => [...deviceKeys.all, "detail"] as const,
  detail: (id: string) => [...deviceKeys.details(), id] as const,
};
```

### 2. Optimistic Updates

對於使用者期望即時反饋的操作（如 Enable/Disable），採用樂觀更新：

```typescript
const toggleMutation = useMutation({
  mutationFn: (id: string) => deviceAPI.disable(id),
  onMutate: async (id) => {
    await queryClient.cancelQueries({ queryKey: deviceKeys.all });
    const previous = queryClient.getQueryData(deviceKeys.lists());
    queryClient.setQueryData(deviceKeys.lists(), (old) =>
      old?.map((d) => (d.id === id ? { ...d, status: "disabled" } : d)),
    );
    return { previous };
  },
  onError: (err, id, context) => {
    queryClient.setQueryData(deviceKeys.lists(), context?.previous);
  },
  onSettled: () => {
    queryClient.invalidateQueries({ queryKey: deviceKeys.all });
  },
});
```

### 3. Error Handling

統一在 Query Client 層處理全域錯誤（如網路錯誤、401），個別 Query 可覆寫：

```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 2,
      staleTime: 30 * 1000, // 30 秒
      gcTime: 5 * 60 * 1000, // 5 分鐘
    },
  },
});
```

## Trade-offs

| 決策                     | 優點               | 缺點           |
| ------------------------ | ------------------ | -------------- |
| 保留 axios service layer | 低侵入性，逐步遷移 | 多一層抽象     |
| 手動 Query Keys          | 完全控制 Cache     | 需紀律維護     |
| 不導入 WebSocket         | 簡化實作           | 無真正即時推送 |

## Future Enhancements

1. **Phase 2**: 導入 WebSocket 即時狀態推送（`/ws/devices/status`）
2. **Phase 3**: 自動化 Go Struct -> TypeScript 型別產生工具
