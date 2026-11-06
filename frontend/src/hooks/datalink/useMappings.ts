import { observeRecordingHomeSourceMutation, captureRecordingHomeSourceMutation, completeRecordingHomeSourceMutation } from './useRecordingHomeSources';
/**
 * Mapping Query Hooks
 *
 * 使用 TanStack Query 管理 Mapping 資源的 Server State。
 * 包含列表查詢、CRUD 操作及狀態切換。
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { mappingAPI } from '../../services/datalink';
import { mappingKeys } from './keys';
import type {
  Mapping,
  CreateMappingRequest,
  UpdateMappingRequest,
  MappingPreviewRequest,
} from '../../types/datalink';

// =============================================================================
// 查詢參數型別
// =============================================================================

/** Mapping 列表篩選參數 */
export interface MappingListFilters {
  point_id?: string;
  tag_id?: string;
  enabled?: boolean;
  limit?: number;
  offset?: number;
}

// =============================================================================
// 查詢 Hooks
// =============================================================================

/**
 * 取得 Mapping 列表
 *
 * @param filters - 可選的篩選條件
 * @returns Query 結果
 */
export function useMappingsQuery(filters?: MappingListFilters) {
  return useQuery({
    queryKey: mappingKeys.list(filters),
    queryFn: () => mappingAPI.list(filters),
  });
}

/**
 * 取得單一 Mapping 詳情
 *
 * @param id - Mapping ID
 * @returns Query 結果
 */
export function useMappingQuery(id: string) {
  return useQuery({
    queryKey: mappingKeys.detail(id),
    queryFn: () => mappingAPI.get(id),
    enabled: !!id,
  });
}

// =============================================================================
// Mutation Hooks
// =============================================================================

/**
 * 建立 Mapping
 */
export function useCreateMappingMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateMappingRequest) => observeRecordingHomeSourceMutation(queryClient, () => mappingAPI.create(data)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
    },
  });
}

/**
 * 更新 Mapping
 */
export function useUpdateMappingMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMappingRequest }) =>
      observeRecordingHomeSourceMutation(queryClient, () => mappingAPI.update(id, data)),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
      queryClient.invalidateQueries({ queryKey: mappingKeys.detail(id) });
    },
  });
}

/**
 * 刪除 Mapping
 */
export function useDeleteMappingMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => observeRecordingHomeSourceMutation(queryClient, () => mappingAPI.delete(id)),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
    },
  });
}

/**
 * 切換 Mapping 啟用狀態（樂觀更新）
 */
export function useToggleMappingStatusMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, enabled }: { id: string; enabled: boolean }) => {
      return mappingAPI.update(id, { enabled: !enabled });
    },
    onMutate: async ({ id, enabled }) => {
      const sourceObservation = captureRecordingHomeSourceMutation(queryClient);
      await queryClient.cancelQueries({ queryKey: mappingKeys.lists() });
      const previousMappings = queryClient.getQueryData<Mapping[]>(
        mappingKeys.list()
      );

      queryClient.setQueryData<Mapping[]>(mappingKeys.list(), (old) =>
        old?.map((m) => (m.id === id ? { ...m, enabled: !enabled } : m))
      );

      return { previousMappings, sourceObservation };
    },
    onSuccess: (_data, _variables, context) => {
      completeRecordingHomeSourceMutation(queryClient, context?.sourceObservation);
    },
    onError: (_, __, context) => {
      if (context?.previousMappings) {
        queryClient.setQueryData(mappingKeys.list(), context.previousMappings);
      }
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
    },
  });
}

/**
 * 預覽 Mapping 轉換結果
 */
export function usePreviewMappingMutation() {
  return useMutation({
    mutationFn: (data: MappingPreviewRequest) => mappingAPI.preview(data),
  });
}

/**
 * 驗證轉換管線
 */
export function useValidatePipelineMutation() {
  return useMutation({
    mutationFn: (pipeline: MappingPreviewRequest['transform_pipeline']) =>
      mappingAPI.validatePipeline(pipeline),
  });
}
