import type { AcquisitionPlanPreview, AcquisitionPlanPreviewResult, AcquisitionPreviewScope } from '../types/acquisitionPlanPreview';
import { studioV2DatalinkApi } from './studioV2Workspace';
import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';

type Row = Record<string, unknown>;
const row = (v: unknown): v is Row => !!v && typeof v === 'object' && !Array.isArray(v);
const count = (v: unknown) => typeof v === 'number' && Number.isSafeInteger(v) && v >= 0;
const number = (v: unknown) => typeof v === 'number' && Number.isFinite(v) && v >= 0;
const identity = (v: unknown) => typeof v === 'string' && safeRuntimeIdentity(v) !== '';
// connector.ResolveAcquisitionConfig emits this namespaced SHA-256 revision.
const capabilityRevision = (v: unknown) => typeof v === 'string' && /^acquisition-config-v1:[a-f0-9]{64}$/.test(v);
const code = (v: unknown) => typeof v === 'string' && /^[a-z][a-z0-9_-]{0,100}$/.test(v);
const list = (v: unknown, validate: (item: unknown) => boolean) => Array.isArray(v) && v.every(validate);
function actual(v: unknown): boolean {
  if (!row(v)) return false;
  if (v.status === 'unavailable') return code(v.reason);
  if (v.status !== 'available' || !identity(v.acquisition_id) || typeof v.observed_at !== 'string' || !Number.isFinite(Date.parse(v.observed_at)) ||
      !identity(v.declared_plan_revision) || !identity(v.initial_plan_revision) ||
      !list(v.phases, p => row(p) && ['initial', 'safe', 'fallback'].includes(String(p.phase)) && identity(p.plan_revision)) ||
      !count(v.wire_attempts) || !count(v.retry_wire_attempts) || !count(v.fallback_wire_attempts) ||
      Number(v.retry_wire_attempts) > Number(v.wire_attempts) || Number(v.fallback_wire_attempts) > Number(v.wire_attempts) ||
      !number(v.cycle_latency_ms) || typeof v.overrun !== 'boolean' || !(v.last_good_age_ms === null || number(v.last_good_age_ms)) || !row(v.logical_requests)) return false;
  return v.logical_requests.status === 'not_started' || (v.logical_requests.status === 'observed' &&
    count(v.logical_requests.succeeded) && count(v.logical_requests.failed) && count(v.logical_requests.canceled));
}
function valid(v: unknown, scope: AcquisitionPreviewScope): v is AcquisitionPlanPreview {
  if (!row(v) || v.workspace_id !== scope.workspace_id || v.device_id !== scope.device_id || v.workspace_revision !== scope.expected_workspace_revision) return false;
  if (v.status === 'unavailable') return code(v.reason);
  return v.status === 'available' && identity(v.preview_revision) && (!scope.expected_preview_revision || v.preview_revision === scope.expected_preview_revision) && v.estimate === true &&
    list(v.groups, g => row(g) && identity(g.polling_group_id) && count(g.poll_interval_ms) && Number(g.poll_interval_ms) > 0 && identity(g.plan_revision) && identity(g.declared_plan_revision) &&
      capabilityRevision(g.capability_revision) && ['declared_no_connection', 'declared_current_connection', 'safe_current_connection'].includes(String(g.basis)) && count(g.acquisition_point_count) && count(g.planned_requests) &&
      list(g.unique_wire_units, u => row(u) && typeof u.area === 'string' && /^[A-Za-z0-9_-]+$/.test(u.area) && typeof u.function === 'string' && /^[A-Za-z0-9_-]*$/.test(u.function) && (u.address_unit_bits === 1 || u.address_unit_bits === 16) && count(u.unique_address_units)) &&
      (g.request_limits === undefined || list(g.request_limits, r => row(r) && (r.kind === 'contiguous' || r.kind === 'random') && typeof r.area === 'string' && /^[A-Za-z0-9_-]*$/.test(r.area) && typeof r.function === 'string' && /^[A-Za-z0-9_-]*$/.test(r.function) && count(r.native_unit_limit) && Number(r.native_unit_limit) > 0 && count(r.request_frame_bytes) && count(r.response_frame_bytes) && typeof r.frame_cost_known === 'boolean')) &&
      list(g.point_errors, p => row(p) && identity(p.point_id) && code(p.code)) && list(g.fallback_reasons, code) && actual(g.actual)) &&
    row(v.output_loss) && (v.output_loss.status === 'unavailable' ? code(v.output_loss.reason) : v.output_loss.status === 'partial' && v.output_loss.history_complete === false &&
      list(v.output_loss.retained_detail, d => row(d) && identity(d.point_id) && count(d.losses)));
}

/** One coherent read; never installs, probes, or mutates bootstrap caches. */
export async function readAcquisitionPlanPreview(scope: AcquisitionPreviewScope, signal?: AbortSignal): Promise<AcquisitionPlanPreviewResult> {
  const captured = { ...scope };
  if (!identity(captured.workspace_id) || !identity(captured.device_id) || typeof captured.expected_workspace_revision !== 'string' ||
      (captured.expected_workspace_revision !== '' && !identity(captured.expected_workspace_revision)) ||
      (captured.expected_preview_revision !== undefined && !identity(captured.expected_preview_revision))) return { ok: false, reason: 'unavailable' };
  try {
    const params = { workspace_id: captured.workspace_id, device_id: captured.device_id, expected_workspace_revision: captured.expected_workspace_revision,
      ...(captured.expected_preview_revision === undefined ? {} : { expected_preview_revision: captured.expected_preview_revision }) };
    const response = await studioV2DatalinkApi.get('/studio-v2/workspace/acquisition-plan-preview', { params, signal });
    if (signal?.aborted || response.data?.success !== true || !valid(response.data.data, captured)) return { ok: false, reason: 'unavailable' };
    return { ok: true, value: response.data.data };
  } catch (error) {
    const status = (error as { response?: { status?: unknown } } | null)?.response?.status;
    return { ok: false, reason: status === 403 ? 'foreign-scope' : status === 409 ? 'stale-reference' : 'unavailable' };
  }
}
