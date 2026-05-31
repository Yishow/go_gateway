import type { APIResponse, StudioV2RuntimeContext } from '../types/datalink';
import { studioV2DatalinkApi } from './studioV2Workspace';

export const studioV2RuntimeContextAPI = {
  async get(): Promise<StudioV2RuntimeContext> {
    const res = await studioV2DatalinkApi.get<APIResponse<StudioV2RuntimeContext>>(
      '/studio-v2/workspace/runtime-context',
    );
    return res.data.data!;
  },
};
