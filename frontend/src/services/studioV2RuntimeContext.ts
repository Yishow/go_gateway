import type { EvidenceResult, EvidenceScope, RuleEvidenceView } from '../types/runtimeRuleEvidence';
import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';
import type { APIResponse } from '../types/datalink';
import type { RuntimeSourcePage, RuntimeSourcePageRequest, RuntimeSourceSelection, StudioV2RuntimeContext } from '../types/studioV2RuntimeContext';
import { projectRuntimeContext, projectRuntimeRuleEvidence } from './runtimeContextProjection';
import { studioV2DatalinkApi } from './studioV2Workspace';
import { normalizeRuntimeSourcePageRequest, parseRuntimeSourcePage } from './runtimeSourcePageParser';

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
  async getSourceEvidence(scope: EvidenceScope, selection: RuntimeSourceSelection, signal?: AbortSignal): Promise<EvidenceResult<RuleEvidenceView>> {
    // Capture only explicit fields before awaiting transport. Neither mutable caller
    // state nor a server-owned identity may change what this response is bound to.
    const params = {
      workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, device_id: scope.device_id,
      rule_id: selection.rule_id, rule_revision: selection.rule_revision, point_id: selection.point_id,
      ...(selection.mapping_id === undefined ? {} : { mapping_id: selection.mapping_id }),
    };
    if (Object.values(params).some((value) => !value || safeRuntimeIdentity(value) !== value)) return unavailable();
    try {
      const res = await studioV2DatalinkApi.get<APIResponse<unknown>>('/studio-v2/workspace/runtime-context', { params, signal });
      if (res.data.success !== true) return unavailable();
      const parsed = projectRuntimeRuleEvidence(res.data.data, params);
      if (!parsed.ok) return parsed;
      const source = parsed.value.source;
      if (source === null) return unavailable();
      if (source.rule_id !== params.rule_id || source.point_id !== params.point_id || source.saved_revision !== params.rule_revision ||
        source.mapping_id !== (params.mapping_id ?? null)) return { ok: false, reason: 'stale-reference', action: 'refresh' };
      if (params.mapping_id === undefined && (source.tag_id !== null || [source.saved, source.effective].some((facts) =>
        facts !== null && (facts.mapping_revision !== null || facts.result_type !== null ||
          facts.transforms.some((step) => step.kind === 'scale' || step.kind === 'cast'))))) return unavailable();
      return parsed;
    } catch (error) {
      const code = record(record(record(record(error).response).data).error).code;
      if (code === 'stale-reference') return { ok: false, reason: code, action: 'refresh' };
      if (code === 'foreign-scope') return { ok: false, reason: code, action: 'open-workspace-list' };
      return unavailable();
    }
  },
  async getSourcePage(scope: EvidenceScope, request: RuntimeSourcePageRequest, signal?: AbortSignal): Promise<EvidenceResult<RuntimeSourcePage>> {
    const observed = { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, device_id: scope.device_id };
    const selected = normalizeRuntimeSourcePageRequest(request);
    const failure: EvidenceResult<RuntimeSourcePage> = { ok: false, reason: 'unavailable', action: 'refresh' };
    if (selected === null || Object.values(observed).some((value) => !value || safeRuntimeIdentity(value) !== value)) return failure;
    const params = { ...observed, source_list: selected.kind, q: selected.query, limit: selected.limit,
      ...(selected.cursor === undefined ? {} : { cursor: selected.cursor }),
      ...(selected.rule_id === undefined ? {} : { rule_id: selected.rule_id, rule_revision: selected.rule_revision }),
    };
    try {
      const res = await studioV2DatalinkApi.get<APIResponse<unknown>>('/studio-v2/workspace/runtime-context', { params, signal });
      if (res.data.success !== true) return failure;
      const data = record(res.data.data);
      if (!safeRuntimeIdentity(data.workspace_id)) return failure;
      if (data.workspace_id !== observed.workspace_id) return { ok: false, reason: 'foreign-scope', action: 'open-workspace-list' };
      return parseRuntimeSourcePage(data.source_page, observed, selected);
    } catch (error) {
      const code = record(record(record(record(error).response).data).error).code;
      if (code === 'stale-reference') return { ok: false, reason: code, action: 'refresh' };
      if (code === 'foreign-scope') return { ok: false, reason: code, action: 'open-workspace-list' };
      return failure;
    }
  },
};
