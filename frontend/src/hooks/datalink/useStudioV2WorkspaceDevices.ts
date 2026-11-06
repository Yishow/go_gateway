import { observeRecordingHomeSourceMutation } from './useRecordingHomeSources';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from './keys';
import { modbusShareStatusKey } from './useModbusShareStatus';
import { studioV2WorkspaceDevicesAPI } from '../../services/studioV2WorkspaceDevices';
import type { CreateDeviceRequest, UpdateDeviceRequest } from '../../types/datalink';
import type { StudioV2AvailabilityRequest } from '../../types/studioV2Availability';

export function useStudioV2WorkspaceDevicesQuery(enabled: boolean) {
  return useQuery({
    queryKey: studioV2WorkspaceKeys.devices(),
    queryFn: () => studioV2WorkspaceDevicesAPI.list(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
    staleTime: 5000,
  });
}

export function useCreateStudioV2WorkspaceDeviceMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: (request: CreateDeviceRequest) => observeRecordingHomeSourceMutation(queryClient, () => studioV2WorkspaceDevicesAPI.create(request), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}

export function useUpdateStudioV2WorkspaceDeviceMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: ({ deviceId, request }: { deviceId: string; request: UpdateDeviceRequest }) =>
      observeRecordingHomeSourceMutation(queryClient, () => studioV2WorkspaceDevicesAPI.update(deviceId, request), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}

export function useDeleteStudioV2WorkspaceDeviceMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: (deviceId: string) => observeRecordingHomeSourceMutation(queryClient, () => studioV2WorkspaceDevicesAPI.remove(deviceId), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}

export function useUpdateStudioV2WorkspaceDeviceAvailabilityMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: ({ deviceId, request }: { deviceId: string; request: StudioV2AvailabilityRequest }) =>
      observeRecordingHomeSourceMutation(queryClient, () => studioV2WorkspaceDevicesAPI.updateAvailability(deviceId, request), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}

export function useUpdateStudioV2WorkspaceDeviceOrderMutation(workspaceId?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    retry: false,
    mutationFn: (orderedDeviceIDs: string[]) => observeRecordingHomeSourceMutation(queryClient, () => studioV2WorkspaceDevicesAPI.updateOrder(orderedDeviceIDs), workspaceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}
