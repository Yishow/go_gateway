import { studioV2DatalinkApi } from './studioV2Workspace';
import { parseWriteGroupAPIData } from '../utils/studioV2WriteGroupJson';
import { parseWriteGroupLifecycle } from '../utils/studioV2WriteGroupLifecycleJson';
import { requireLifecycleMatch, safeLifecycleRequest } from './studioV2LifecycleTransport';
export const studioV2WriteGroupLifecycleAPI = {
  lifecycle(groupId: string) {
    return safeLifecycleRequest(async () => {
      const response = await studioV2DatalinkApi.get<unknown>(`/studio-v2/workspace/write-groups/${encodeURIComponent(groupId)}/lifecycle`);
      const result = parseWriteGroupAPIData(response.data, 'recording lifecycle', parseWriteGroupLifecycle, 12);
      return requireLifecycleMatch(result, result.group_id === groupId);
    });
  },
};
