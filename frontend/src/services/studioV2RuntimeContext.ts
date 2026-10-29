import type { EvidenceResult, EvidenceScope, RuleEvidenceView } from '../types/runtimeRuleEvidence';
import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';
import type { APIResponse } from '../types/datalink';
import type { StudioV2RuntimeContext } from '../types/studioV2RuntimeContext';
import { projectRuntimeContext, projectRuntimeRuleEvidence } from './runtimeContextProjection';
import { studioV2DatalinkApi } from './studioV2Workspace';

const record = (value: unknown): Record<string, unknown> => value !== null && typeof value === 'object' && !Array.isArray(value) ? value as Record<string, unknown> : {};
const unavailable = (): EvidenceResult<RuleEvidenceView> => ({ ok: false, reason: 'unavailable', action: 'refresh' });

export const studioV2RuntimeContextAPI = {
  async get(): Promise<StudioV2RuntimeContext> {
    const res = await studioV2DatalinkApi.get<APIResponse<StudioV2RuntimeContext>>(
      '/studio-v2/workspace/runtime-context',
    );
    return projectRuntimeContext(res.data.data);
  },
  async getRuleEvidence(scope: EvidenceScope, signal?: AbortSignal): Promise<EvidenceResult<RuleEvidenceView>> {
    const params = { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, device_id: scope.device_id };
    if (Object.values(params).some((value) => !value || safeRuntimeIdentity(value) !== value)) return unavailable();
    try {
      const res = await studioV2DatalinkApi.get<APIResponse<unknown>>('/studio-v2/workspace/runtime-context', { params, signal });
      if (res.data.success !== true) return unavailable();
      return projectRuntimeRuleEvidence(res.data.data, params);
    } catch (error) {
      const code = record(record(record(record(error).response).data).error).code;
      if (code === 'stale-reference') return { ok: false, reason: code, action: 'refresh' };
      if (code === 'foreign-scope') return { ok: false, reason: code, action: 'open-workspace-list' };
      return unavailable();
    }
  },
};
