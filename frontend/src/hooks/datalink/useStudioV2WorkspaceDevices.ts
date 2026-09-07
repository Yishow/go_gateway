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
  });
}

export function useCreateStudioV2WorkspaceDeviceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: CreateDeviceRequest) => studioV2WorkspaceDevicesAPI.create(request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}

export function useUpdateStudioV2WorkspaceDeviceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ deviceId, request }: { deviceId: string; request: UpdateDeviceRequest }) =>
      studioV2WorkspaceDevicesAPI.update(deviceId, request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}

export function useDeleteStudioV2WorkspaceDeviceMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (deviceId: string) => studioV2WorkspaceDevicesAPI.remove(deviceId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}

export function useUpdateStudioV2WorkspaceDeviceAvailabilityMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ deviceId, request }: { deviceId: string; request: StudioV2AvailabilityRequest }) =>
      studioV2WorkspaceDevicesAPI.updateAvailability(deviceId, request),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}

export function useUpdateStudioV2WorkspaceDeviceOrderMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (orderedDeviceIDs: string[]) => studioV2WorkspaceDevicesAPI.updateOrder(orderedDeviceIDs),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.bootstrap() });
      await queryClient.invalidateQueries({ queryKey: studioV2WorkspaceKeys.devices() });
      await queryClient.invalidateQueries({ queryKey: modbusShareStatusKey });
    },
  });
}
