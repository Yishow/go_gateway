/**
 * Tag Query Hooks
 *
 * 使用 TanStack Query 管理 Tag 資源的 Server State。
 * 包含列表查詢、詳情查詢、CRUD 操作及狀態切換。
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { tagAPI } from '../../services/datalink';
import { tagKeys } from './keys';
import type {
  Tag,
  CreateTagRequest,
  UpdateTagRequest,
} from '../../types/datalink';

// =============================================================================
// 查詢參數型別
// =============================================================================

/** Tag 列表篩選參數 */
export interface TagListFilters {
  status?: string;
  data_type?: string;
  key_prefix?: string;
  limit?: number;
  offset?: number;
}

// =============================================================================
// 查詢 Hooks
// =============================================================================

/**
 * 取得 Tag 列表
 *
 * @param filters - 可選的篩選條件
 * @returns Query 結果
 */
export function useTagsQuery(filters?: TagListFilters) {
  return useQuery({
    queryKey: tagKeys.list(filters),
    queryFn: () => tagAPI.list(filters),
  });
}

/**
 * 取得單一 Tag 詳情
 *
 * @param id - Tag ID
 * @returns Query 結果
 */
export function useTagQuery(id: string) {
  return useQuery({
    queryKey: tagKeys.detail(id),
    queryFn: () => tagAPI.get(id),
    enabled: !!id,
  });
}

// =============================================================================
// Mutation Hooks
// =============================================================================

/**
 * 建立 Tag
 */
export function useCreateTagMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTagRequest) => tagAPI.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
    },
  });
}

/**
 * 更新 Tag
 */
export function useUpdateTagMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTagRequest }) =>
      tagAPI.update(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
      queryClient.invalidateQueries({ queryKey: tagKeys.detail(id) });
    },
  });
}

/**
 * 刪除 Tag
 */
export function useDeleteTagMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => tagAPI.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
    },
  });
}

/**
 * 啟用 Tag（樂觀更新）
 */
export function useActivateTagMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => tagAPI.activate(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: tagKeys.lists() });
      const previousTags = queryClient.getQueryData<Tag[]>(tagKeys.list());
      
      queryClient.setQueryData<Tag[]>(tagKeys.list(), (old) =>
        old?.map((t) => (t.id === id ? { ...t, status: 'active' as const } : t))
      );
      
      return { previousTags };
    },
    onError: (_, __, context) => {
      if (context?.previousTags) {
        queryClient.setQueryData(tagKeys.list(), context.previousTags);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
    },
  });
}

/**
 * 退役 Tag（樂觀更新）
 */
export function useRetireTagMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => tagAPI.retire(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: tagKeys.lists() });
      const previousTags = queryClient.getQueryData<Tag[]>(tagKeys.list());
      
      queryClient.setQueryData<Tag[]>(tagKeys.list(), (old) =>
        old?.map((t) => (t.id === id ? { ...t, status: 'retired' as const } : t))
      );
      
      return { previousTags };
    },
    onError: (_, __, context) => {
      if (context?.previousTags) {
        queryClient.setQueryData(tagKeys.list(), context.previousTags);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
    },
  });
}

/**
 * 批量建立 Tags
 */
export function useBatchCreateTagsMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (tags: CreateTagRequest[]) => tagAPI.batchCreate(tags),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
    },
  });
}
