import { safeRuntimeIdentity as id, safeRuntimeLabel as label, safeRuntimeProtocol, safeRuntimeAlignment } from '../utils/runtimeOperatorSafety';
import type { RuntimeStatusWithDiagnostics } from '../types/runtimeDiagnostics';

type RecordValue = Record<string, unknown>;
const record = (value: unknown): RecordValue => value && typeof value === 'object' && !Array.isArray(value) ? value as RecordValue : {};
const list = (value: unknown): RecordValue[] => Array.isArray(value) ? value.map(record) : [];
const numeric = (value: unknown): number => typeof value === 'number' && Number.isFinite(value) && value >= 0 ? value : 0;
const choice = <T extends string>(value: unknown, allowed: readonly T[], fallback: T): T => allowed.find((item) => item === value) ?? fallback;
const timestamp = (value: unknown): string | null => typeof value === 'string' && /^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(?:\.\d+)?Z$/.test(value) && Number.isFinite(Date.parse(value)) ? value : null;
const stages = ['collector', 'mapping', 'runtime_projection', 'database_delivery', 'modbus_share_delivery', 'collected', 'mapped', 'db_write', 'db_write_failed', 'db_write_succeeded', 'accepted', 'queued', 'sql-committed', 'unknown'] as const;
const stage = (value: unknown): string => choice(value, stages, 'unknown');
const optionalStage = (value: unknown): string | undefined => value == null || value === '' ? undefined : stage(value);
const failureCode = (value: unknown): string | undefined => ['runtime_snapshot_unavailable', 'modbus_share_delivery', 'database_delivery'].find((item) => item === value);
function scope(value: unknown): string {
  if (value === 'runtime') return value;
  if (typeof value !== 'string') return 'unknown';
  const match = /^(device|point|tag):(.+)$/.exec(value);
  return match && id(match[2]) ? `${match[1]}:${id(match[2])}` : 'unknown';
}
/** Detached legacy operator snapshot. It contains health evidence, never driver messages or connector data. */
export function projectRuntimeStatus(value: unknown): RuntimeStatusWithDiagnostics {
  const v = record(value); const metrics = record(v.metrics); const truth = record(v.snapshot_state);
  const state = choice(truth.state, ['ready', 'empty', 'degraded', 'unavailable', 'stale'], 'unavailable');
  return {
    running: v.running === true, uptime_seconds: numeric(v.uptime_seconds),
    snapshot_state: { state, empty: truth.empty === true, degraded: truth.degraded === true, unavailable: state === 'unavailable' || truth.unavailable === true, stale: truth.stale === true },
    metrics: { collected_total: numeric(metrics.collected_total), write_success_total: numeric(metrics.write_success_total),
      write_error_total: numeric(metrics.write_error_total), mapping_error_total: numeric(metrics.mapping_error_total), point_state_error_total: numeric(metrics.point_state_error_total) },
    collectors: list(v.collectors).map((c) => ({
      device_id: id(c.device_id), device_name: label(c.device_name), protocol: safeRuntimeProtocol(c.protocol),
      status: choice(c.status, ['idle', 'running', 'warning', 'error', 'unavailable'], 'unavailable'), running: c.running === true && Boolean(id(c.device_id)) && c.availability_status === 'available' && ['idle', 'running', 'warning', 'error'].includes(String(c.status)) && safeRuntimeProtocol(c.protocol) !== 'unavailable',
      availability_status: choice(c.availability_status, ['available', 'unavailable'], 'unavailable'),
      points_total: numeric(c.points_total), points_healthy: numeric(c.points_healthy), points_stale: numeric(c.points_stale), points_error: numeric(c.points_error),
      last_read_at: timestamp(c.last_read_at), last_error: c.last_error ? 'runtime-error' : null,
      breaker_state: choice(c.breaker_state, ['closed', 'open', 'half-open'], 'unavailable'),
      projection_alignment: safeRuntimeAlignment(c.projection_alignment),
      runtime_projection_version: id(c.runtime_projection_version), workspace_projection_version: id(c.workspace_projection_version),
    })),
    diagnostics: list(v.diagnostics).map((d) => ({
      scope: scope(d.scope), device_id: id(d.device_id), point_id: id(d.point_id), tag_id: id(d.tag_id),
      last_success_at: timestamp(d.last_success_at), last_failure_at: timestamp(d.last_failure_at),
      ...(optionalStage(d.latest_successful_stage) ? { latest_successful_stage: optionalStage(d.latest_successful_stage) } : {}),
      ...(optionalStage(d.failure_stage) ? { failure_stage: optionalStage(d.failure_stage) } : {}),
      ...(failureCode(d.failure_code) ? { failure_code: failureCode(d.failure_code) } : {}),
      stages: list(d.stages).map((s) => ({ stage: stage(s.stage), status: choice(s.status, ['success', 'failed', 'unknown'], 'unknown'), observed_at: timestamp(s.observed_at) })),
    })),
    database_delivery: list(v.database_delivery).map((d) => ({
      device_id: id(d.device_id), point_id: id(d.point_id), tag_id: id(d.tag_id), status: choice(d.status, ['succeeded', 'failed', 'disabled', 'accepted', 'queued', 'unknown'], 'unknown'),
      stages: Array.isArray(d.stages) ? d.stages.map(stage) : [], ...(optionalStage(d.failed_stage) ? { failed_stage: optionalStage(d.failed_stage) } : {}), observed_at: timestamp(d.observed_at) ?? '',
      last_success_at: timestamp(d.last_success_at), last_failure_at: timestamp(d.last_failure_at),
    })),
    modbus_share_delivery: list(v.modbus_share_delivery).map((d) => ({
      device_id: id(d.device_id), point_id: id(d.point_id), tag_id: id(d.tag_id), status: choice(d.status, ['succeeded', 'failed', 'disabled'], 'unavailable'),
      stage: stage(d.stage), observed_at: timestamp(d.observed_at) ?? '',
    })),
  };
}
