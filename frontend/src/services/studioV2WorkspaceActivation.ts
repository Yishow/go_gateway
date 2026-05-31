import type { APIResponse } from '../types/datalink';
import type { StudioV2ActivationResponse } from '../types/studioV2Activation';
import { studioV2DatalinkApi } from './studioV2Workspace';

export const studioV2WorkspaceActivationAPI = {
  async activate(): Promise<StudioV2ActivationResponse> {
    const res = await studioV2DatalinkApi.post<APIResponse<StudioV2ActivationResponse>>(
      '/studio-v2/workspace/activate',
    );
    return res.data.data!;
  },
};
