/**
 * Device Query Hooks
 *
 * 使用 TanStack Query 管理 Device 資源的 Server State。
 * 包含列表查詢、詳情查詢、CRUD 操作及樂觀更新。
 */
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { deviceAPI } from '../../services/datalink';
import { deviceKeys, mappingKeys, pointKeys } from './keys';
import type {
  Device,
  CreateDeviceRequest,
  UpdateDeviceRequest,
} from '../../types/datalink';

// =============================================================================
// 查詢參數型別
// =============================================================================

/** Device 列表篩選參數 */
export interface DeviceListFilters {
  protocol?: string;
  status?: string;
  limit?: number;
  offset?: number;
}

// =============================================================================
// 查詢 Hooks
// =============================================================================

/**
 * 取得 Device 列表
 *
 * @param filters - 可選的篩選條件
 * @returns Query 結果，包含 data, isLoading, error 等狀態
 */
export function useDevicesQuery(filters?: DeviceListFilters) {
  return useQuery({
    queryKey: deviceKeys.list(filters),
    queryFn: () => deviceAPI.list(filters),
  });
}

/**
 * 取得單一 Device 詳情
 *
 * @param id - Device ID
 * @returns Query 結果
 */
export function useDeviceQuery(id: string) {
  return useQuery({
    queryKey: deviceKeys.detail(id),
    queryFn: () => deviceAPI.get(id),
    enabled: !!id,
  });
}

// =============================================================================
// Mutation Hooks
// =============================================================================

/**
 * 建立 Device
 *
 * 成功後自動使列表快取失效。
 */
export function useCreateDeviceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateDeviceRequest) => deviceAPI.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() });
    },
  });
}

/**
 * 更新 Device
 *
 * 成功後自動使列表及該 Device 詳情快取失效。
 */
export function useUpdateDeviceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateDeviceRequest }) =>
      deviceAPI.update(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() });
      queryClient.invalidateQueries({ queryKey: deviceKeys.detail(id) });
    },
  });
}

/**
 * 刪除 Device
 *
 * 成功後自動使列表快取失效。
 */
export function useDeleteDeviceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deviceAPI.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() });
      queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
      queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });
    },
  });
}

/**
 * Check Device Readiness
 */
export function useCheckReadinessMutation() {
  return useMutation({
    mutationFn: (id: string) => deviceAPI.checkReadiness(id),
  });
}

/**
 * Test Connection (Legacy/Simple)
 */
export function useTestConnectionMutation() {
  return useMutation({
    mutationFn: (id: string) => deviceAPI.testConnection(id),
  });
}

/**
 * 切換 Device 狀態（啟用/停用）
 *
 * 使用樂觀更新：UI 立即反映，失敗時回退。
 */
export function useToggleDeviceStatusMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      currentStatus,
    }: {
      id: string;
      currentStatus: string;
    }) => {
      if (currentStatus === 'active') {
        return deviceAPI.disable(id);
      } else {
        return deviceAPI.activate(id);
      }
    },

    // 樂觀更新：立即反映在 UI
    onMutate: async ({ id, currentStatus }) => {
      // 取消進行中的查詢，避免覆蓋樂觀更新
      await queryClient.cancelQueries({ queryKey: deviceKeys.lists() });

      // 快照目前資料以便回退
      const previousDevices = queryClient.getQueryData<Device[]>(
        deviceKeys.list()
      );

      // 樂觀更新快取
      const newStatus = currentStatus === 'active' ? 'disabled' : 'active';
      queryClient.setQueryData<Device[]>(deviceKeys.list(), (old) =>
        old?.map((d) =>
          d.id === id ? { ...d, status: newStatus as Device['status'] } : d
        )
      );

      return { previousDevices };
    },

    // 錯誤時回退
    onError: (_, __, context) => {
      if (context?.previousDevices) {
        queryClient.setQueryData(deviceKeys.list(), context.previousDevices);
      }
    },

    // 無論成功或失敗，都重新驗證資料
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: deviceKeys.lists() });
    },
  });
}
