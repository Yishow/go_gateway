import type {
  AcquisitionAnomalyDetailCode,
  AcquisitionAnomalyErrorCode,
  AcquisitionAnomalyFilter,
  AcquisitionAnomalyFilterField,
  AcquisitionAnomalyItem,
  AcquisitionAnomalyPage,
  AcquisitionAnomalySummary,
} from '../types/acquisitionAnomaly';
import { ACQUISITION_ANOMALY_DETAIL_CODES, ACQUISITION_ANOMALY_KINDS, ACQUISITION_ANOMALY_STATES } from '../types/acquisitionAnomaly';
import { api } from './datalinkClient';

/** Backend page bounds (design condition from the anomaly API). */
export const ACQUISITION_ANOMALY_DEFAULT_LIMIT = 50;
export const ACQUISITION_ANOMALY_MAX_LIMIT = 500;

/** Typed failure; carries only allowlisted classification, never Axios response, config or raw text. */
export class AcquisitionAnomalyError extends Error {
  readonly name = 'AcquisitionAnomalyError';
  constructor(readonly code: AcquisitionAnomalyErrorCode, readonly retryable: boolean, readonly field?: AcquisitionAnomalyFilterField) {
    super(code);
  }
}

export const acquisitionAnomalyErrorKey = (error: AcquisitionAnomalyError): string =>
  `recording-home:acquisitionAnomalies.errors.${error.code}`;

const FILTER_FIELDS: readonly AcquisitionAnomalyFilterField[] = ['state', 'kind', 'limit', 'offset', 'device_id'];
const BACKEND_CODES: Record<number, [string, AcquisitionAnomalyErrorCode, boolean]> = {
  400: ['acquisition_anomaly_filter_invalid', 'invalid_filter', false],
  404: ['acquisition_anomaly_device_not_found', 'device_not_found', false],
  500: ['acquisition_anomaly_unavailable', 'unavailable', true],
  503: ['acquisition_anomaly_capability_unavailable', 'capability_unavailable', false],
};

function record(value: unknown): Record<string, unknown> | undefined {
  return value !== null && typeof value === 'object' && !Array.isArray(value) ? value as Record<string, unknown> : undefined;
}

function classifyError(error: unknown): AcquisitionAnomalyError {
  if (error instanceof AcquisitionAnomalyError) return error;
  const response = record(record(error)?.response);
  const status = response?.status;
  const body = record(record(response?.data)?.error);
  const expected = typeof status === 'number' ? BACKEND_CODES[status] : undefined;
  if (!expected || body?.code !== expected[0]) return new AcquisitionAnomalyError('unknown', true);
  const field = FILTER_FIELDS.find(name => name === body.field);
  return new AcquisitionAnomalyError(expected[1], expected[2], expected[1] === 'invalid_filter' ? field : undefined);
}

class Malformed extends Error {}
function malformed(): never { throw new Malformed(); }
function integer(value: unknown, minimum = 0): number {
  return typeof value === 'number' && Number.isSafeInteger(value) && value >= minimum ? value : malformed();
}
function text(value: unknown): string { return typeof value === 'string' ? value : malformed(); }
function nullableText(value: unknown): string | null { return value === null ? null : text(value); }
function timestamp(value: unknown): string {
  return typeof value === 'string' && /^\d{4}-\d{2}-\d{2}T/.test(value) && Number.isFinite(Date.parse(value)) ? value : malformed();
}
function oneOf<T extends string>(value: unknown, allowed: readonly T[]): T {
  return allowed.find(item => item === value) ?? malformed();
}
function detailCodes(value: unknown): AcquisitionAnomalyDetailCode[] {
  return text(value).split(';').map(part => part.trim()).filter(Boolean)
    .map(part => ACQUISITION_ANOMALY_DETAIL_CODES.find(code => code === part) ?? 'unknown');
}
function successData(value: unknown): Record<string, unknown> {
  const body = record(value);
  return body?.success === true ? record(body.data) ?? malformed() : malformed();
}
function guarded<T>(read: () => T): T {
  try { return read(); } catch (error) {
    if (error instanceof Malformed) throw new AcquisitionAnomalyError('malformed_response', true);
    throw error;
  }
}

function parseItem(value: unknown, deviceId: string): AcquisitionAnomalyItem {
  const v = record(value) ?? malformed();
  const state = oneOf(v.state, ACQUISITION_ANOMALY_STATES);
  const startedAt = timestamp(v.started_at);
  const endedAt = v.ended_at === null ? null : timestamp(v.ended_at);
  const duration = v.duration_ms === null ? null : integer(v.duration_ms);
  if ((state === 'ended') !== (endedAt !== null) || (endedAt === null) !== (duration === null)) malformed();
  if (endedAt !== null && Date.parse(endedAt) < Date.parse(startedAt)) malformed();
  if (v.device_id !== deviceId) malformed();
  return {
    id: integer(v.id, 1), device_id: deviceId, point_id: nullableText(v.point_id), point_name: nullableText(v.point_name),
    point_address: nullableText(v.point_address), kind: oneOf(v.kind, ACQUISITION_ANOMALY_KINDS), state,
    started_at: startedAt, ended_at: endedAt, duration_ms: duration, last_value: nullableText(v.last_value), detail_codes: detailCodes(v.detail),
  };
}

/** Parses one device page; anything off-contract is a failure, never an empty list. */
export function parseAcquisitionAnomalyPage(value: unknown, deviceId: string, request: { limit: number; offset: number }): AcquisitionAnomalyPage {
  return guarded(() => {
    const data = successData(value);
    if (!Array.isArray(data.items)) malformed();
    const items = data.items.map(item => parseItem(item, deviceId));
    const total = integer(data.total), limit = integer(data.limit, 1), offset = integer(data.offset);
    // total and items come from separate backend queries, so total is not cross-checked against the page.
    if (limit !== request.limit || offset !== request.offset || items.length > limit) malformed();
    return { items, total, limit, offset };
  });
}

/** Parses per-device active counts; an explicit 0 is kept distinct from a missing device. */
export function parseAcquisitionAnomalySummary(value: unknown): AcquisitionAnomalySummary {
  return guarded(() => {
    const data = successData(value);
    if (!Array.isArray(data.items)) malformed();
    const seen = new Set<string>();
    const items = data.items.map(entry => {
      const v = record(entry) ?? malformed();
      const deviceId = text(v.device_id);
      if (!deviceId || seen.has(deviceId)) malformed();
      seen.add(deviceId);
      return { device_id: deviceId, active_count: integer(v.active_count) };
    });
    const total = integer(data.total_active);
    if (total !== items.reduce((sum, item) => sum + item.active_count, 0)) malformed();
    return { items, total_active: total };
  });
}

function normalizeFilter(deviceId: string, filter: AcquisitionAnomalyFilter) {
  const reject = (field: AcquisitionAnomalyFilterField): never => { throw new AcquisitionAnomalyError('invalid_filter', false, field); };
  if (typeof deviceId !== 'string' || deviceId.trim() === '') reject('device_id');
  if (filter.state !== undefined && !ACQUISITION_ANOMALY_STATES.includes(filter.state)) reject('state');
  if (filter.kind !== undefined && !ACQUISITION_ANOMALY_KINDS.includes(filter.kind)) reject('kind');
  const limit = filter.limit ?? ACQUISITION_ANOMALY_DEFAULT_LIMIT, offset = filter.offset ?? 0;
  if (!Number.isSafeInteger(limit) || limit < 1 || limit > ACQUISITION_ANOMALY_MAX_LIMIT) reject('limit');
  if (!Number.isSafeInteger(offset) || offset < 0) reject('offset');
  return { ...(filter.state ? { state: filter.state } : {}), ...(filter.kind ? { kind: filter.kind } : {}), limit, offset };
}

/** Lists one device's events newest first. Failures reject with AcquisitionAnomalyError. */
export async function listDeviceAcquisitionAnomalies(deviceId: string, filter: AcquisitionAnomalyFilter): Promise<AcquisitionAnomalyPage> {
  const params = normalizeFilter(deviceId, filter);
  const response = await api.get<unknown>(`/devices/${encodeURIComponent(deviceId)}/acquisition-anomalies`, { params })
    .catch((error: unknown) => { throw classifyError(error); });
  return parseAcquisitionAnomalyPage(response.data, deviceId, params);
}

/** Reads every device's active event count. Failures reject; they never read as zero. */
export async function getAcquisitionAnomalySummary(): Promise<AcquisitionAnomalySummary> {
  const response = await api.get<unknown>('/acquisition-anomalies/summary').catch((error: unknown) => { throw classifyError(error); });
  return parseAcquisitionAnomalySummary(response.data);
}

/** Active count for one device, or null when the summary is missing or does not list it (status unknown). */
export function activeAnomalyCount(summary: AcquisitionAnomalySummary | undefined, deviceId: string): number | null {
  return summary?.items.find(item => item.device_id === deviceId)?.active_count ?? null;
}
