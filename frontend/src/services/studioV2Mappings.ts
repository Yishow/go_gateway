import type { APIResponse, StudioV2WorkspaceMappingRecord } from '../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../types/studioV2RuntimeApply';
import { studioV2DatalinkApi } from './studioV2Workspace';

export interface StudioV2WorkspaceMappingRequest {
  rule_id: string;
  address: string;
  tag_key: string;
  display_name: string;
  unit: string;
  target_type: StudioV2WorkspaceMappingRecord['target_type'];
  scale: number;
  offset: number;
  enabled: boolean;
}

export const studioV2MappingsAPI = {
  async list(): Promise<StudioV2WorkspaceMappingRecord[]> {
    const res = await studioV2DatalinkApi.get<APIResponse<StudioV2WorkspaceMappingRecord[]>>('/studio-v2/workspace/mappings');
    return res.data.data ?? [];
  },

  async create(request: StudioV2WorkspaceMappingRequest): Promise<StudioV2RuntimeAppliedRecord<StudioV2WorkspaceMappingRecord>> {
    const res = await studioV2DatalinkApi.post<APIResponse<StudioV2RuntimeAppliedRecord<StudioV2WorkspaceMappingRecord>>>('/studio-v2/workspace/mappings', request);
    return res.data.data!;
  },

  async update(mappingId: string, request: StudioV2WorkspaceMappingRequest): Promise<StudioV2RuntimeAppliedRecord<StudioV2WorkspaceMappingRecord>> {
    const res = await studioV2DatalinkApi.put<APIResponse<StudioV2RuntimeAppliedRecord<StudioV2WorkspaceMappingRecord>>>(`/studio-v2/workspace/mappings/${mappingId}`, request);
    return res.data.data!;
  },

  async remove(mappingId: string): Promise<void> {
    await studioV2DatalinkApi.delete(`/studio-v2/workspace/mappings/${mappingId}`);
  },
};
