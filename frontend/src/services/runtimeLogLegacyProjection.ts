import { parsePreciseInstant } from '../utils/preciseInstant';
import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';
import type { RuntimeLogLevel, RuntimeLogRecord } from '../types/runtimeLogs';

// Closed consumer catalogue of diagnostics/projection.go; not a new event authority.
const templates = {
  'startup.begin': ['info', 'startup'],
  'startup.ready': ['info', 'startup'],
  'startup.failed': ['error', 'startup'],
  'startup.config_invalid': ['error', 'startup'],
  'startup.path_ambiguous': ['error', 'startup'],
  'startup.path_unusable': ['error', 'startup'],
  'startup.owner_busy': ['error', 'startup'],
  'startup.log_unavailable': ['error', 'startup'],
  'startup.tray_failed': ['error', 'startup'],
  'startup.bind_failed': ['error', 'startup'],
  'startup.assets_missing': ['error', 'startup'],
  'startup.database_linked': ['error', 'startup'],
  'startup.db_open_failed': ['error', 'startup'],
  'startup.db_ping_failed': ['error', 'startup'],
  'startup.migration_failed': ['error', 'startup'],
  'startup.service_failed': ['error', 'startup'],
  'startup.pipeline_failed': ['error', 'startup'],
  'runtime.started': ['info', 'runtime'],
  'runtime.degraded': ['warn', 'runtime'],
  'runtime.share_degraded': ['warn', 'runtime'],
  'runtime.group_failed': ['error', 'runtime'],
  'runtime.stopped': ['info', 'runtime'],
  'shutdown.begin': ['info', 'shutdown'],
  'shutdown.finalizing': ['info', 'shutdown'],
  'shutdown.complete': ['info', 'shutdown'],
  'shutdown.timeout': ['warn', 'shutdown'],
  'shutdown.failed': ['error', 'shutdown'],
  'http.access': ['info', 'http'],
  'http.request': ['info', 'http'],
  'http.recovered': ['error', 'http'],
  'http.recovery': ['error', 'http'],
  'http.server_error': ['error', 'http'],
  'raw.suppressed': ['info', 'application'],
} as const;
export type RuntimeLogLegacyCode = keyof typeof templates;
export type LegacySource = 'startup' | 'runtime' | 'shutdown' | 'http' | 'application' | 'standard' | 'slog';
export type LegacyPhase = 'startup' | 'http' | 'runtime' | 'pipeline' | 'share' | 'connections' | 'database' | 'diagnostics' | 'tray' | 'owner';
export interface LegacyPresentation {
  kind: 'legacy'; code: RuntimeLogLegacyCode; category: 'process' | 'http' | 'suppressed';
  timestamp: string; level: RuntimeLogLevel; source: LegacySource; cursor: string; truncated: boolean;
  method: string | null; status: number | null; durationMs: number | null; requestID: string | null; phase: LegacyPhase | null;
}
export type RuntimeLogLegacyProjection = LegacyPresentation | { kind: 'suppressed'; timestamp: string | null; cursor: string | null; truncated: boolean };
const envelopeKeys = ['instance_id', 'sequence', 'timestamp', 'level', 'source', 'code', 'message', 'fields', 'truncated'];
const methods = ['GET', 'HEAD', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS', 'CONNECT', 'TRACE'] as const;
const phases: readonly LegacyPhase[] = ['startup', 'http', 'runtime', 'pipeline', 'share', 'connections', 'database', 'diagnostics', 'tray', 'owner'];
const suppressedSources: readonly LegacySource[] = ['application', 'standard', 'slog', 'http', 'runtime'];
const encoder = new TextEncoder();
function refuse(): never { throw new Error('legacy_presentation_suppressed'); }
/** No getters, inherited state, symbols, unknown keys or input spread enter the detached result. */
function object(value: unknown, keys: readonly string[]): Record<string, unknown> {
  if (!value || typeof value !== 'object' || Object.getPrototypeOf(value) !== Object.prototype) return refuse();
  const result: Record<string, unknown> = {};
  for (const key of Reflect.ownKeys(value)) {
    const descriptor = Object.getOwnPropertyDescriptor(value, key);
    if (typeof key !== 'string' || !keys.includes(key) || !descriptor || !('value' in descriptor)) return refuse();
    result[key] = descriptor.value;
  }
  return result;
}
function utc(value: unknown): string | null {
  if (typeof value !== 'string' || !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,9})?Z$/.test(value)) return null;
  const at = parsePreciseInstant(value);
  // Keep original nanoseconds; reject calendar normalization performed by Date.parse.
  return at && new Date(at.seconds * 1000).toISOString().slice(0, 19) === value.slice(0, 19) ? value : null;
}
function transportCursor(instance: unknown, sequence: unknown): string | null {
  if (typeof instance !== 'string' || !/^[A-Za-z0-9_-]{1,128}$/.test(instance) || safeRuntimeIdentity(instance) !== instance ||
      typeof sequence !== 'string' || !/^[1-9][0-9]{0,19}$/.test(sequence) || BigInt(sequence) > 18446744073709551615n) return null;
  return `${instance}:${sequence}`;
}
function enumeration<T extends string>(value: unknown, allowed: readonly T[]): T {
  return allowed.find(item => item === value) ?? refuse();
}
function optional<T>(fields: Record<string, unknown>, key: string, parse: (value: unknown) => T): T | null {
  return Object.hasOwn(fields, key) ? parse(fields[key]) : null;
}
function numeric(value: unknown, max: number, min = 0, integer = false): number {
  if (typeof value !== 'number' || !Number.isFinite(value) || value < min || value > max || (integer && !Number.isInteger(value))) return refuse();
  return value;
}
function requestID(value: unknown): string {
  return typeof value === 'string' && /^(?:[0-9a-fA-F]{32}|[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$/.test(value) ? value : refuse();
}
/** Pure display projection only. No request, operation, current-health or delivery authority. */
export function projectRuntimeLogLegacy(record: RuntimeLogRecord): RuntimeLogLegacyProjection {
  let timestamp: string | null = null, cursor: string | null = null, truncated = false;
  try {
    const input = object(record, envelopeKeys);
    timestamp = utc(input.timestamp); cursor = transportCursor(input.instance_id, input.sequence); truncated = input.truncated === true;
    if (!timestamp || !cursor || typeof input.truncated !== 'boolean' || typeof input.message !== 'string' ||
        typeof input.code !== 'string' || !Object.hasOwn(templates, input.code)) return refuse();
    const code = input.code as RuntimeLogLegacyCode, [level, fixedSource] = templates[code];
    const source = code === 'raw.suppressed' ? enumeration(input.source, suppressedSources) : fixedSource;
    if (input.level !== level || input.source !== source) return refuse();
    const category = code === 'raw.suppressed' ? 'suppressed' : code.startsWith('http.') ? 'http' : 'process';
    const fields = object(input.fields, category === 'http' ? ['method', 'route', 'status', 'duration_ms', 'request_id'] :
      code.startsWith('shutdown.') ? ['phase'] : []);
    // Registered route ownership remains exclusively in Go. Validate only its producer shape and never retain its bytes.
    if (Object.hasOwn(fields, 'route') && (typeof fields.route !== 'string' || encoder.encode(fields.route).byteLength > 1024)) return refuse();
    return { kind: 'legacy', code, category, timestamp, level, source, cursor, truncated,
      method: optional(fields, 'method', value => enumeration(value, methods)),
      status: optional(fields, 'status', value => numeric(value, 599, 100, true)),
      durationMs: optional(fields, 'duration_ms', value => numeric(value, 86400000)), requestID: optional(fields, 'request_id', requestID),
      phase: optional(fields, 'phase', value => enumeration(value, phases)) };
  } catch { return { kind: 'suppressed', timestamp, cursor, truncated }; }
}
