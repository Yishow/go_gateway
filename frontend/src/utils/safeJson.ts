import { isSafeRuntimeSensorText, safeRuntimeIdentity, safeRuntimeAddress } from './runtimeOperatorSafety';
import type { MappingPreviewResponse, RuntimeDeviceStatusEvent, RuntimeStreamStateEvent, RuntimeValueEvent } from '../types/datalink';
import type { RuntimeTruthState } from '../types/runtimeTruth';
import type { StudioV2ActivationResponse } from '../types/studioV2Activation';
import { BACKEND_ERROR_CODES } from './backendErrorCodes';

export const MAX_SAFE_JSON_BYTES = 64 * 1024;
export const MAX_SAFE_JSON_DEPTH = 6;
export const MAX_SAFE_JSON_ARRAY_LENGTH = 256;
export const MAX_SAFE_JSON_OBJECT_KEYS = 64;
export const MAX_SAFE_JSON_STRING_LENGTH = 256;
export const MAX_TYPED_FIELD_LENGTH = 128;
const SAFE_ERROR_CODES = new Set<string>(BACKEND_ERROR_CODES);
const SAFE_ERROR_ACTIONS = new Set([
  'retry', 'reconnect', 'reload', 'retry runtime stream', 'retry the runtime stream',
  'retry the preview stream', 'retry preview', 'wait_for_supported_operation',
  'review_source', 'open_write_groups', 'review_request', 'review_group', 'disable_group',
  'select_recording_destination', 'select_new_managed_table_or_advanced_table',
  'select_plain_storage',
]);

function byteLength(value: string): number {
  return typeof TextEncoder === 'undefined' ? value.length : new TextEncoder().encode(value).byteLength;
}
export type BoundedJsonStringLengthLimit = (path: readonly (string | number)[]) => number | undefined;

function cloneBoundedJson(
  value: unknown, depth: number, maxDepth: number, path: readonly (string | number)[],
  stringLengthLimitForPath?: BoundedJsonStringLengthLimit,
): unknown | null {
  if (depth > maxDepth) return null;
  if (value === null || typeof value === 'boolean') return value;
  if (typeof value === 'number') return Number.isFinite(value) ? value : null;
  if (typeof value === 'string') {
    const stringLimit = stringLengthLimitForPath?.(path) ?? MAX_SAFE_JSON_STRING_LENGTH;
    return value.length <= stringLimit ? value : null;
  }
  if (Array.isArray(value)) {
    if (value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
    const result: unknown[] = [];
    for (const item of value) {
      const itemPath = stringLengthLimitForPath ? [...path, result.length] : path;
      const cloned = cloneBoundedJson(item, depth + 1, maxDepth, itemPath, stringLengthLimitForPath);
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
    const itemPath = stringLengthLimitForPath ? [...path, key] : path;
    const cloned = cloneBoundedJson(item, depth + 1, maxDepth, itemPath, stringLengthLimitForPath);
    if (cloned === null && item !== null) return null;
    result[key] = cloned;
  }
  return result;
}
/** Parse untrusted JSON into a bounded, detached JSON value without throwing. */
export function parseBoundedJson(
  value: unknown, maxDepth: number = MAX_SAFE_JSON_DEPTH,
  stringLengthLimitForPath?: BoundedJsonStringLengthLimit,
): unknown | null {
  if (typeof value === 'string' && byteLength(value) > MAX_SAFE_JSON_BYTES) return null;
  let parsed = value;
  if (typeof value === 'string') {
    try {
      parsed = JSON.parse(value) as unknown;
    } catch {
      return null;
    }
  }
  return cloneBoundedJson(parsed, 0, maxDepth, [], stringLengthLimitForPath);
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
  operationId?: string;
  retryable?: boolean;
}

function boundedField(value: unknown): string | undefined {
  if (typeof value !== 'string') return undefined;
  const normalized = value.trim();
  return normalized.length > 0 && normalized.length <= MAX_TYPED_FIELD_LENGTH ? normalized : undefined;
}

function findFirstField(sources: (Record<string, unknown> | undefined)[], keys: string[]): string | undefined {
  for (const src of sources) {
    if (!src) continue;
    for (const key of keys) {
      const val = boundedField(src[key]);
      if (val) return val;
    }
  }
  return undefined;
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
  const responseDataData = responseData && typeof responseData.data === 'object' && responseData.data !== null && !Array.isArray(responseData.data)
    ? responseData.data as Record<string, unknown>
    : undefined;
  const rootData = typeof root.data === 'object' && root.data !== null && !Array.isArray(root.data)
    ? root.data as Record<string, unknown>
    : undefined;
  const rootDataData = rootData && typeof rootData.data === 'object' && rootData.data !== null && !Array.isArray(rootData.data)
    ? rootData.data as Record<string, unknown>
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
  const requestId = findFirstField([record, nested], ['request_id', 'requestID']);
  const operationId = findFirstField(
    [record, nested, responseData, responseDataData, rootData, rootDataData, root],
    ['operation_id', 'operationID'],
  );
  return {
    code: code && SAFE_ERROR_CODES.has(code) ? code : undefined,
    action: action && SAFE_ERROR_ACTIONS.has(action) ? action : undefined,
    requestId,
    ...(operationId ? { operationId } : {}),
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

function envelopeFields(env: NormalizedTypedEnvelope) {
  return {
    ...(env.operationId ? { operation_id: env.operationId } : {}),
    ...(env.code ? { code: env.code } : {}),
    ...(env.action ? { action: env.action } : {}),
    ...(env.requestId ? { request_id: env.requestId } : {}),
    ...(env.retryable !== undefined ? { retryable: env.retryable } : {}),
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
    return {
      device_id: deviceId,
      status: status as StudioV2ActivationResponse['results'][number]['status'],
      ...envelopeFields(normalizeTypedEnvelope(result)),
    };
  });
  if (results.some((result) => result === null)) return null;
  return {
    workspace_id: workspaceId,
    results: results as StudioV2ActivationResponse['results'],
    ...envelopeFields(normalizeTypedEnvelope(record)),
  };
}

// Runtime-only scalar guard; shared preview/activation parsers retain their own contracts.
function safeRuntimeScalar(value: unknown): unknown | undefined {
  const scalar = boundedRuntimeValue(value);
  if (typeof scalar !== 'string') return scalar;
  if (!isSafeRuntimeSensorText(scalar)) return undefined;
  return scalar;
}
const runtimeIdentity = (value: unknown): string | undefined => safeRuntimeIdentity(value) || undefined;

const runtimeTimestamp = (value: unknown): string | undefined =>
  typeof value === 'string' && /^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(?:\.\d+)?(?:Z|[+-]\d\d:\d\d)$/.test(value) && Number.isFinite(Date.parse(value)) ? value : undefined;

/** Rebuilds the allowlisted runtime value event shape from bounded JSON. */
export function parseRuntimeValueRecord(value: unknown): RuntimeValueEvent | null {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) return null;
  const record = value as Record<string, unknown>;
  const deviceId = runtimeIdentity(record.device_id);
  const pointId = runtimeIdentity(record.point_id);
  const address = safeRuntimeAddress(record.address);
  const timestamp = runtimeTimestamp(record.timestamp);
  const rawValue = safeRuntimeScalar(record.raw_value);
  const transformedValue = safeRuntimeScalar(record.transformed_value);
  if (!deviceId || !pointId || !address || !timestamp || rawValue === undefined || transformedValue === undefined) return null;
  if (!['good', 'bad', 'uncertain', 'missing', 'stale', 'invalid'].includes(String(record.quality))) return null;
  if (record.value_state !== undefined && record.value_state !== 'unsupported') return null;
  if (typeof record.stale !== 'boolean') return null;
  return {
    device_id: deviceId,
    point_id: pointId,
    address,
    raw_value: rawValue,
    transformed_value: transformedValue,
    quality: record.quality as RuntimeValueEvent['quality'],
    ...(record.value_state === 'unsupported' ? { value_state: 'unsupported' as const } : {}),
    stale: record.stale,
    timestamp,
  };
}

/** Rebuilds the allowlisted runtime device status shape from bounded JSON. */
export function parseRuntimeStatusRecord(value: unknown): RuntimeDeviceStatusEvent | null {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) return null;
  const record = value as Record<string, unknown>;
  const deviceId = runtimeIdentity(record.device_id);
  const status = boundedString(record.status);
  const breakerState = boundedString(record.breaker_state);
  if (!deviceId || !status || !breakerState ||
    !['idle', 'running', 'warning', 'error'].includes(status) ||
    !['closed', 'open', 'half-open'].includes(breakerState) ||
    ![record.points_total, record.points_healthy, record.points_stale, record.points_error].every(
      (item) => typeof item === 'number' && Number.isFinite(item),
    )) return null;
  const lastReadAt = record.last_read_at;
  const lastError = record.last_error;
  if (lastReadAt !== undefined && lastReadAt !== null && !runtimeTimestamp(lastReadAt)) return null;
  if (lastError !== undefined && lastError !== null && !boundedString(lastError)) return null;
  return {
    device_id: deviceId,
    status: status as RuntimeDeviceStatusEvent['status'],
    points_total: record.points_total as number,
    points_healthy: record.points_healthy as number,
    points_stale: record.points_stale as number,
    points_error: record.points_error as number,
    breaker_state: breakerState,
    last_read_at: lastReadAt === null ? null : runtimeTimestamp(lastReadAt),
    last_error: lastError ? 'runtime-error' : null,
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
  if (!stateName || !['ready', 'empty', 'degraded', 'unavailable', 'stale'].includes(stateName) || typeof state.empty !== 'boolean' || typeof state.degraded !== 'boolean' ||
    typeof state.unavailable !== 'boolean' || typeof state.stale !== 'boolean') return null;
  const reason = state.reason === undefined ? undefined :
    typeof state.reason === 'string' && SAFE_ERROR_CODES.has(state.reason) ? state.reason : 'unavailable';
  if (state.reason !== undefined && !reason) return null;
  const deviceId = record.device_id === undefined ? undefined : runtimeIdentity(record.device_id);
  const timestamp = record.timestamp === undefined ? undefined : runtimeTimestamp(record.timestamp);
  if (record.device_id !== undefined && !deviceId) return null;
  if (record.timestamp !== undefined && !timestamp) return null;
  const envelope = normalizeTypedEnvelope(record);
  return {
    ...(deviceId ? { device_id: deviceId } : {}),
    ...(timestamp ? { timestamp } : {}),
    ...envelopeFields({ ...envelope, requestId: runtimeIdentity(envelope.requestId), operationId: runtimeIdentity(envelope.operationId) }),
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
