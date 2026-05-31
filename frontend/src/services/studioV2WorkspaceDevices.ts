import type { CreateDeviceRequest, Device, UpdateDeviceRequest, APIResponse, StudioV2Workspace } from '../types/datalink';
import type { StudioV2Availability, StudioV2AvailabilityRequest } from '../types/studioV2Availability';
import type { StudioV2RuntimeApply } from '../types/studioV2RuntimeApply';
import { studioV2DatalinkApi } from './studioV2Workspace';

export type StudioV2WorkspaceDeviceRecord = Device & Partial<StudioV2Availability> & Partial<StudioV2RuntimeApply>;

export const studioV2WorkspaceDevicesAPI = {
  async list(): Promise<StudioV2WorkspaceDeviceRecord[]> {
    const res = await studioV2DatalinkApi.get<APIResponse<StudioV2WorkspaceDeviceRecord[]>>('/studio-v2/workspace/devices');
    return res.data.data ?? [];
  },

  async create(request: CreateDeviceRequest): Promise<StudioV2WorkspaceDeviceRecord> {
    const res = await studioV2DatalinkApi.post<APIResponse<StudioV2WorkspaceDeviceRecord>>('/studio-v2/workspace/devices', request);
    return res.data.data!;
  },

  async update(deviceId: string, request: UpdateDeviceRequest): Promise<StudioV2WorkspaceDeviceRecord> {
    const res = await studioV2DatalinkApi.put<APIResponse<StudioV2WorkspaceDeviceRecord>>(`/studio-v2/workspace/devices/${deviceId}`, request);
    return res.data.data!;
  },

  async updateAvailability(deviceId: string, request: StudioV2AvailabilityRequest): Promise<StudioV2WorkspaceDeviceRecord> {
    const res = await studioV2DatalinkApi.put<APIResponse<StudioV2WorkspaceDeviceRecord>>(`/studio-v2/workspace/devices/${deviceId}/availability`, request);
    return res.data.data!;
  },

  async remove(deviceId: string): Promise<void> {
    await studioV2DatalinkApi.delete(`/studio-v2/workspace/devices/${deviceId}`);
  },

  async updateOrder(orderedDeviceIDs: string[]): Promise<StudioV2Workspace> {
    const res = await studioV2DatalinkApi.put<APIResponse<StudioV2Workspace>>('/studio-v2/workspace/device-order', {
      ordered_device_ids: orderedDeviceIDs,
    });
    return res.data.data!;
  },
};
