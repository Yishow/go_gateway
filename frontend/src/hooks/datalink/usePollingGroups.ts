/**
 * PollingGroup Query Hooks
 *
 * 使用 TanStack Query 管理 PollingGroup 資源的 Server State。
 * 包含列表查詢、CRUD 操作及樂觀更新。
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { pollingGroupAPI } from '../../services/datalink';
import { pollingGroupKeys } from './keys';
import type {
  CreatePollingGroupRequest,
  UpdatePollingGroupRequest,
} from '../../types/datalink';

// =============================================================================
// 查詢 Hooks
// =============================================================================

/**
 * 取得 PollingGroup 列表
 *
 * @returns Query 結果，包含 data, isLoading, error 等狀態
 */
export function usePollingGroupsQuery() {
  return useQuery({
    queryKey: pollingGroupKeys.lists(),
    queryFn: () => pollingGroupAPI.list(),
  });
}

// =============================================================================
// Mutation Hooks
// =============================================================================

/**
 * 建立 PollingGroup
 *
 * 成功後自動使列表快取失效。
 */
export function useCreatePollingGroupMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePollingGroupRequest) => pollingGroupAPI.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: pollingGroupKeys.lists() });
    },
  });
}

/**
 * 更新 PollingGroup
 *
 * 成功後自動使列表快取失效。
 */
export function useUpdatePollingGroupMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePollingGroupRequest }) =>
      pollingGroupAPI.update(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: pollingGroupKeys.lists() });
    },
  });
}

/**
 * 刪除 PollingGroup
 *
 * 成功後自動使列表快取失效。
 */
export function useDeletePollingGroupMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => pollingGroupAPI.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: pollingGroupKeys.lists() });
    },
  });
}
