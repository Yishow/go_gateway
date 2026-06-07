import type { APIResponse } from '../types/datalink';
import type { StudioV2WorkspaceAuditHistory } from '../types/studioV2WorkspaceAudit';
import { studioV2DatalinkApi } from './studioV2Workspace';

export const studioV2WorkspaceAuditAPI = {
  async list(limit = 10): Promise<StudioV2WorkspaceAuditHistory> {
    const res = await studioV2DatalinkApi.get<APIResponse<StudioV2WorkspaceAuditHistory>>(
      '/studio-v2/workspace/audit-history',
      { params: { limit } },
    );
    return res.data.data ?? { entries: [] };
  },
};
