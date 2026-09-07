import type { MappingPreviewResponse, RuntimeDeviceStatusEvent, RuntimeStreamStateEvent, RuntimeValueEvent } from '../types/datalink';
import type { RuntimeTruthState } from '../types/runtimeTruth';
import type { StudioV2ActivationResponse } from '../types/studioV2Activation';

export const MAX_SAFE_JSON_BYTES = 64 * 1024;
export const MAX_SAFE_JSON_DEPTH = 6;
export const MAX_SAFE_JSON_ARRAY_LENGTH = 256;
export const MAX_SAFE_JSON_OBJECT_KEYS = 64;
export const MAX_SAFE_JSON_STRING_LENGTH = 256;
export const MAX_TYPED_FIELD_LENGTH = 128;
const SAFE_ERROR_CODES = new Set([
  'preview_invalid_request', 'preview_unavailable', 'runtime_device_not_found', 'preview_stream_closed',
  'runtime_snapshot_unavailable', 'runtime_stream_unavailable', 'workspace_not_ready',
  'activation_failed', 'activation_request_invalid', 'readiness_blocked',
  'settings_unavailable', 'settings_update_failed', 'settings_invalid',
  'modbus_share_save_incomplete', 'modbus_share_revision_conflict',
  'modbus_share_hydration_required', 'modbus_share_disabled', 'modbus_share_range_collision',
  'modbus_share_capacity_exceeded', 'modbus_share_listener_bind_failed',
  'modbus_share_workspace_scope', 'modbus_share_reconcile_failed',
  'modbus_share_dirty_unknown', 'modbus_share_projection_required',
  'modbus_share_invalid_geometry', 'validation', 'revision_mismatch', 'not_found', 'internal',
]);
const SAFE_ERROR_ACTIONS = new Set([
  'retry', 'reconnect', 'reload', 'retry runtime stream', 'retry the runtime stream', 'retry the preview stream', 'retry preview',
]);

function byteLength(value: string): number {
  return typeof TextEncoder === 'undefined' ? value.length : new TextEncoder().encode(value).byteLength;
}
function cloneBoundedJson(value: unknown, depth: number): unknown | null {
  if (depth > MAX_SAFE_JSON_DEPTH) return null;
  if (value === null || typeof value === 'boolean') return value;
  if (typeof value === 'number') return Number.isFinite(value) ? value : null;
  if (typeof value === 'string') return value.length <= MAX_SAFE_JSON_STRING_LENGTH ? value : null;
  if (Array.isArray(value)) {
    if (value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
    const result: unknown[] = [];
    for (const item of value) {
      const cloned = cloneBoundedJson(item, depth + 1);
      if (cloned === null && item !== null) return null;
      result.push(cloned);
    }
    return result;
  }
  if (typeof value !== 'object') return null;
  const entries = Object.entries(value as Record<string, unknown>);
  if (entries.length > MAX_SAFE_JSON_OBJECT_KEYS) return null;
  const result: Record<string, unknown> = {};
  for (const [key, item] of entries) {
    if (key.length > MAX_SAFE_JSON_STRING_LENGTH) return null;
    const cloned = cloneBoundedJson(item, depth + 1);
    if (cloned === null && item !== null) return null;
    result[key] = cloned;
  }
  return result;
}
/** Parse untrusted JSON into a bounded, detached JSON value without throwing. */
export function parseBoundedJson(value: unknown): unknown | null {
  if (typeof value === 'string' && byteLength(value) > MAX_SAFE_JSON_BYTES) return null;
  let parsed = value;
  if (typeof value === 'string') {
    try {
      parsed = JSON.parse(value) as unknown;
    } catch {
      return null;
    }
  }
  return cloneBoundedJson(parsed, 0);
}
/** Parse untrusted JSON input into a bounded plain object without throwing. */
export function parseJsonRecord(value: unknown): Record<string, unknown> | null {
  const parsed = parseBoundedJson(value);
  if (typeof parsed !== 'object' || parsed === null || Array.isArray(parsed)) return null;
  return parsed as Record<string, unknown>;
}
/** Read an EventSource message as a bounded detached object without trusting its shape. */
export function parseMessageEventRecord(event: Event): Record<string, unknown> | null {
  return parseJsonRecord((event as MessageEvent<unknown>).data);
}

export interface NormalizedTypedEnvelope {
  code?: string;
  action?: string;
  requestId?: string;
  retryable?: boolean;
}

function boundedField(value: unknown): string | undefined {
  if (typeof value !== 'string') return undefined;
  const normalized = value.trim();
  return normalized.length > 0 && normalized.length <= MAX_TYPED_FIELD_LENGTH ? normalized : undefined;
}

/** Rebuild typed error metadata from an allowlist; raw code/action text is never retained. */
export function normalizeTypedEnvelope(value: unknown): NormalizedTypedEnvelope {
  const root = typeof value === 'object' && value !== null && !Array.isArray(value)
    ? value as Record<string, unknown>
    : {};
  const response = typeof root.response === 'object' && root.response !== null && !Array.isArray(root.response)
    ? root.response as Record<string, unknown>
    : undefined;
  const responseData = response && typeof response.data === 'object' && response.data !== null && !Array.isArray(response.data)
    ? response.data as Record<string, unknown>
    : undefined;
  const record = responseData?.error && typeof responseData.error === 'object' && !Array.isArray(responseData.error)
    ? responseData.error as Record<string, unknown>
    : root.error && typeof root.error === 'object' && !Array.isArray(root.error)
      ? root.error as Record<string, unknown>
      : root;
  const nested = typeof record.error === 'object' && record.error !== null && !Array.isArray(record.error)
    ? record.error as Record<string, unknown>
    : {};
  const code = boundedField(record.code) ?? boundedField(nested.code);
  const action = boundedField(record.action) ?? boundedField(nested.action);
  const requestId = boundedField(record.request_id) ?? boundedField(record.requestID) ??
    boundedField(nested.request_id) ?? boundedField(nested.requestID);
  return {
    code: code && SAFE_ERROR_CODES.has(code) ? code : undefined,
    action: action && SAFE_ERROR_ACTIONS.has(action) ? action : undefined,
    requestId,
    retryable: typeof record.retryable === 'boolean'
      ? record.retryable
      : typeof nested.retryable === 'boolean' ? nested.retryable : undefined,
  };
}

export function boundedString(value: unknown): string | undefined {
  return boundedField(value);
}

function boundedRuntimeValue(value: unknown): unknown | undefined {
  if (value === null || typeof value === 'boolean') return value;
  if (typeof value === 'string') return value.length <= MAX_SAFE_JSON_STRING_LENGTH ? value : undefined;
  return typeof value === 'number' && Number.isFinite(value) ? value : undefined;
}

/** Rebuilds only the scalar Step 3 preview projection needed by the UI. */
export function parseMappingPreviewResponse(value: unknown): MappingPreviewResponse | null {
  const bounded = parseBoundedJson(value);
  if (typeof bounded !== 'object' || bounded === null || Array.isArray(bounded)) return null;
  const root = bounded as Record<string, unknown>;
  const record = (typeof root.data === 'object' && root.data !== null && !Array.isArray(root.data))
    ? (root.data as Record<string, unknown>)
    : root;
  const rawValue = boundedRuntimeValue(record.raw_value);
  const finalValue = boundedRuntimeValue(record.final_value);
  if (rawValue === undefined || finalValue === undefined || !Array.isArray(record.step_results)) return null;
  const stepResults = record.step_results.map((step) => {
    if (typeof step !== 'object' || step === null || Array.isArray(step)) return null;
    const item = step as Record<string, unknown>;
    const stepIndex = item.step_index;
    const stepType = boundedString(item.step_type);
    const rawIn = item.input_value !== undefined ? item.input_value : item.input;
    const rawOut = item.output_value !== undefined ? item.output_value : (item.output !== undefined ? item.output : (item.error ? null : undefined));
    const inputValue = boundedRuntimeValue(rawIn);
    const outputValue = boundedRuntimeValue(rawOut);
    const error = item.error === undefined || item.error === '' ? '' : boundedString(item.error);
    if (!Number.isInteger(stepIndex) || (stepIndex as number) < 0 || !stepType ||
      !['decode', 'cast', 'scale', 'lookup', 'conditional', 'formula'].includes(stepType) ||
      inputValue === undefined || outputValue === undefined || error === undefined) return null;
    return {
      step_index: stepIndex as number,
      step_type: stepType as MappingPreviewResponse['step_results'][number]['step_type'],
      input_value: inputValue,
      output_value: outputValue,
      error,
    };
  });
  if (stepResults.some((step) => step === null)) return null;
  return {
    raw_value: rawValue,
    final_value: finalValue,
    step_results: stepResults as MappingPreviewResponse['step_results'],
    ...(boundedString(record.error) ? { error: boundedString(record.error) } : {}),
  };
}

/** Rebuilds the allowlisted Step 4 activation response and typed result metadata. */
export function parseStudioV2ActivationResponse(value: unknown): StudioV2ActivationResponse | null {
  const bounded = parseBoundedJson(value);
  if (typeof bounded !== 'object' || bounded === null || Array.isArray(bounded)) return null;
  const record = bounded as Record<string, unknown>;
  const workspaceId = boundedString(record.workspace_id);
  if (!workspaceId || !Array.isArray(record.results)) return null;
  const results = record.results.map((item) => {
    if (typeof item !== 'object' || item === null || Array.isArray(item)) return null;
    const result = item as Record<string, unknown>;
    const deviceId = boundedString(result.device_id);
    const status = boundedString(result.status);
    if (!deviceId || !status || !['pending', 'success', 'failed', 'skipped'].includes(status)) return null;
    const envelope = normalizeTypedEnvelope(result);
    return {
      device_id: deviceId,
      status: status as StudioV2ActivationResponse['results'][number]['status'],
      ...(envelope.code ? { code: envelope.code } : {}),
      ...(envelope.action ? { action: envelope.action } : {}),
      ...(envelope.requestId ? { request_id: envelope.requestId } : {}),
      ...(envelope.retryable !== undefined ? { retryable: envelope.retryable } : {}),
    };
  });
  if (results.some((result) => result === null)) return null;
  const envelope = normalizeTypedEnvelope(record);
  return {
    workspace_id: workspaceId,
    results: results as StudioV2ActivationResponse['results'],
    ...(envelope.code ? { code: envelope.code } : {}),
    ...(envelope.action ? { action: envelope.action } : {}),
    ...(envelope.requestId ? { request_id: envelope.requestId } : {}),
    ...(envelope.retryable !== undefined ? { retryable: envelope.retryable } : {}),
  };
}

/** Rebuilds the allowlisted runtime value event shape from bounded JSON. */
export function parseRuntimeValueRecord(value: unknown): RuntimeValueEvent | null {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) return null;
  const record = value as Record<string, unknown>;
  const deviceId = boundedString(record.device_id);
  const pointId = boundedString(record.point_id);
  const address = boundedString(record.address);
  const timestamp = boundedString(record.timestamp);
  const rawValue = boundedRuntimeValue(record.raw_value);
  const transformedValue = boundedRuntimeValue(record.transformed_value);
  if (!deviceId || !pointId || !address || !timestamp || rawValue === undefined || transformedValue === undefined) return null;
  if (record.quality !== 'good' && record.quality !== 'bad' && record.quality !== 'uncertain') return null;
  if (typeof record.stale !== 'boolean') return null;
  return {
    device_id: deviceId,
    point_id: pointId,
    address,
    raw_value: rawValue,
    transformed_value: transformedValue,
    quality: record.quality,
    stale: record.stale,
    timestamp,
  };
}

/** Rebuilds the allowlisted runtime device status shape from bounded JSON. */
export function parseRuntimeStatusRecord(value: unknown): RuntimeDeviceStatusEvent | null {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) return null;
  const record = value as Record<string, unknown>;
  const deviceId = boundedString(record.device_id);
  const status = boundedString(record.status);
  const breakerState = boundedString(record.breaker_state);
  if (!deviceId || !status || !breakerState ||
    !['idle', 'running', 'warning', 'error'].includes(status) ||
    ![record.points_total, record.points_healthy, record.points_stale, record.points_error].every(
      (item) => typeof item === 'number' && Number.isFinite(item),
    )) return null;
  const lastReadAt = record.last_read_at;
  const lastError = record.last_error;
  if (lastReadAt !== undefined && lastReadAt !== null && !boundedString(lastReadAt)) return null;
  if (lastError !== undefined && lastError !== null && !boundedString(lastError)) return null;
  return {
    device_id: deviceId,
    status: status as RuntimeDeviceStatusEvent['status'],
    points_total: record.points_total as number,
    points_healthy: record.points_healthy as number,
    points_stale: record.points_stale as number,
    points_error: record.points_error as number,
    breaker_state: breakerState,
    last_read_at: lastReadAt === null ? null : boundedString(lastReadAt),
    last_error: lastError === null ? null : boundedString(lastError),
  };
}

/** Rebuilds the allowlisted runtime readiness envelope from bounded JSON. */
export function parseRuntimeStreamStateRecord(value: unknown): (RuntimeStreamStateEvent & { stream_state: RuntimeTruthState }) | null {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) return null;
  const record = value as Record<string, unknown>;
  const stateValue = record.stream_state;
  if (typeof stateValue !== 'object' || stateValue === null || Array.isArray(stateValue)) return null;
  const state = stateValue as Record<string, unknown>;
  const stateName = boundedString(state.state);
  if (!stateName || typeof state.empty !== 'boolean' || typeof state.degraded !== 'boolean' ||
    typeof state.unavailable !== 'boolean' || typeof state.stale !== 'boolean') return null;
  const reason = state.reason === undefined ? undefined : boundedString(state.reason);
  if (state.reason !== undefined && !reason) return null;
  const deviceId = record.device_id === undefined ? undefined : boundedString(record.device_id);
  const timestamp = record.timestamp === undefined ? undefined : boundedString(record.timestamp);
  if (record.device_id !== undefined && !deviceId) return null;
  if (record.timestamp !== undefined && !timestamp) return null;
  const envelope = normalizeTypedEnvelope(record);
  return {
    ...(deviceId ? { device_id: deviceId } : {}),
    ...(timestamp ? { timestamp } : {}),
    ...(envelope.code ? { code: envelope.code } : {}),
    ...(envelope.action ? { action: envelope.action } : {}),
    ...(envelope.requestId ? { request_id: envelope.requestId } : {}),
    ...(envelope.retryable !== undefined ? { retryable: envelope.retryable } : {}),
    stream_state: {
      state: stateName,
      empty: state.empty,
      degraded: state.degraded,
      unavailable: state.unavailable,
      stale: state.stale,
      ...(reason ? { reason } : {}),
    },
  };
}
