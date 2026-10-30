import type { RuntimeLogRecord } from '@/types/runtimeLogs';
import type { RuntimeLogLegacyCode } from '@/services/runtimeLogLegacyProjection';

/** Static expected cases from the Go legacy projector; never used by product code. */
export const legacyCases = [
  ['startup.begin', 'info', 'startup'],
  ['startup.ready', 'info', 'startup'],
  ['startup.failed', 'error', 'startup'],
  ['startup.config_invalid', 'error', 'startup'],
  ['startup.path_ambiguous', 'error', 'startup'],
  ['startup.path_unusable', 'error', 'startup'],
  ['startup.owner_busy', 'error', 'startup'],
  ['startup.log_unavailable', 'error', 'startup'],
  ['startup.tray_failed', 'error', 'startup'],
  ['startup.bind_failed', 'error', 'startup'],
  ['startup.assets_missing', 'error', 'startup'],
  ['startup.database_linked', 'error', 'startup'],
  ['startup.db_open_failed', 'error', 'startup'],
  ['startup.db_ping_failed', 'error', 'startup'],
  ['startup.migration_failed', 'error', 'startup'],
  ['startup.service_failed', 'error', 'startup'],
  ['startup.pipeline_failed', 'error', 'startup'],
  ['runtime.started', 'info', 'runtime'],
  ['runtime.degraded', 'warn', 'runtime'],
  ['runtime.share_degraded', 'warn', 'runtime'],
  ['runtime.group_failed', 'error', 'runtime'],
  ['runtime.stopped', 'info', 'runtime'],
  ['shutdown.begin', 'info', 'shutdown'],
  ['shutdown.finalizing', 'info', 'shutdown'],
  ['shutdown.complete', 'info', 'shutdown'],
  ['shutdown.timeout', 'warn', 'shutdown'],
  ['shutdown.failed', 'error', 'shutdown'],
  ['http.access', 'info', 'http'],
  ['http.request', 'info', 'http'],
  ['http.recovered', 'error', 'http'],
  ['http.recovery', 'error', 'http'],
  ['http.server_error', 'error', 'http'],
  ['raw.suppressed', 'info', 'application'],
] as const;
export const legacyInstant = '2026-10-29T10:00:00.123456789Z';
export const legacyRequestID = '01234567-89ab-cdef-0123-456789abcdef';
/** Fields overrides merge; tests may delete keys explicitly to represent missing optional producer fields. */
export function legacyRecord(code: RuntimeLogLegacyCode = 'runtime.started', overrides: Partial<RuntimeLogRecord> = {}): RuntimeLogRecord {
  const row = legacyCases.find(item => item[0] === code);
  if (!row) throw new Error('Unknown legacy fixture code');
  const fields: RuntimeLogRecord['fields'] = code.startsWith('http.') ? { method: 'GET', status: 200, duration_ms: 0, request_id: legacyRequestID } :
    code.startsWith('shutdown.') ? { phase: 'diagnostics' } : {};
  return { instance_id: 'instance-a', sequence: '9007199254740993', timestamp: legacyInstant, code, level: row[1], source: row[2],
    message: 'Candidate message must not drive presentation.', truncated: false, ...overrides, fields: { ...fields, ...overrides.fields } };
}
