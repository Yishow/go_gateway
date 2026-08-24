import type { APIResponse, CreateSourceRuleRequest, SourceRuleRecord } from '../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../types/studioV2RuntimeApply';
import { studioV2DatalinkApi } from './studioV2Workspace';

export interface StudioV2UpdateSourceRuleRequest {
  device_id?: string;
  start_address?: string;
  count?: number;
  data_type?: SourceRuleRecord['data_type'];
  naming_prefix?: string;
  enabled?: boolean;
  skipped_addresses?: string[];
  scale_multiplier?: number;
  scale_offset?: number;
  data_format?: string;
  share_enabled?: boolean;
  share_start_register?: number | null;
  share_stride?: number | null;
}

export const studioV2RulesAPI = {
  async list(): Promise<SourceRuleRecord[]> {
    const res = await studioV2DatalinkApi.get<APIResponse<SourceRuleRecord[]>>('/studio-v2/workspace/source-rules');
    return res.data.data ?? [];
  },

  async create(request: CreateSourceRuleRequest): Promise<StudioV2RuntimeAppliedRecord<SourceRuleRecord>> {
    const res = await studioV2DatalinkApi.post<APIResponse<StudioV2RuntimeAppliedRecord<SourceRuleRecord>>>('/studio-v2/workspace/source-rules', request);
    return res.data.data!;
  },

  async update(ruleId: string, request: StudioV2UpdateSourceRuleRequest): Promise<StudioV2RuntimeAppliedRecord<SourceRuleRecord>> {
    const res = await studioV2DatalinkApi.put<APIResponse<StudioV2RuntimeAppliedRecord<SourceRuleRecord>>>(`/studio-v2/workspace/source-rules/${ruleId}`, request);
    return res.data.data!;
  },

  async remove(ruleId: string): Promise<void> {
    await studioV2DatalinkApi.delete(`/studio-v2/workspace/source-rules/${ruleId}`);
  },
};
