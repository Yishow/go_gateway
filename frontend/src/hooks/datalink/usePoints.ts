/**
 * Point Query Hooks
 *
 * 使用 TanStack Query 管理 Point 資源的 Server State。
 * 包含列表查詢、詳情查詢、CRUD 操作及即時輪詢。
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { pointAPI } from '../../services/datalink';
import { pointKeys } from './keys';
import type {
  Point,
  CreatePointRequest,
  UpdatePointRequest,
} from '../../types/datalink';

// =============================================================================
// 查詢參數型別
// =============================================================================

/** Point 列表篩選參數 */
export interface PointListFilters {
  device_id?: string;
  polling_group_id?: string;
  enabled?: boolean;
  data_type?: string;
  limit?: number;
  offset?: number;
}

// =============================================================================
// 查詢 Hooks
// =============================================================================

/**
 * 取得 Point 列表
 *
 * @param filters - 可選的篩選條件
 * @returns Query 結果
 */
export function usePointsQuery(filters?: PointListFilters) {
  return useQuery({
    queryKey: pointKeys.list(filters),
    queryFn: () => pointAPI.list(filters),
  });
}

/**
 * 取得單一 Point 詳情
 *
 * @param id - Point ID
 * @returns Query 結果
 */
export function usePointQuery(id: string) {
  return useQuery({
    queryKey: pointKeys.detail(id),
    queryFn: () => pointAPI.get(id),
    enabled: !!id,
  });
}

// =============================================================================
// Mutation Hooks
// =============================================================================

/**
 * 建立 Point
 */
export function useCreatePointMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreatePointRequest) => pointAPI.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
    },
  });
}

/**
 * 更新 Point
 */
export function useUpdatePointMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePointRequest }) =>
      pointAPI.update(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
      queryClient.invalidateQueries({ queryKey: pointKeys.detail(id) });
    },
  });
}

/**
 * 刪除 Point
 */
export function useDeletePointMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => pointAPI.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
    },
  });
}

/**
 * 立即輪詢 Point
 */
export function usePollPointMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => pointAPI.pollNow(id),
    onSuccess: (data, id) => {
        // Optionally update the point cache if the result contains the latest value
        // For now, we just invalidate details to force a refresh if needed
        queryClient.invalidateQueries({ queryKey: pointKeys.detail(id) });
    }
  });
}

/**
 * 批量輪詢 Points
 */
export function useBatchPollPointsMutation() {
    return useMutation({
        mutationFn: (ids: string[]) => pointAPI.batchPoll(ids),
    });
}
