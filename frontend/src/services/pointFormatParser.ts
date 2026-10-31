import type { PointFormatMutationResult, PointFormatOperation, PointFormatResolution, PointFormatScope, PointFormatSnapshot } from '../types/pointFormat';
import { PointFormatError, record } from './pointFormatErrors';
const invalid = (): never => { throw new PointFormatError('point_format_unknown', 'unknown'); };
export const wordFormats = ['ABCD', 'BADC', 'CDAB', 'DCBA'] as const;
export function text(value: unknown, allowEmpty = true): string {
  if (typeof value !== 'string' || value.length > 256 || (!allowEmpty && !value.trim()) || value !== value.trim()) return invalid();
  return value;
}
function oneOf<T extends string | number>(value: unknown, choices: readonly T[]): T {
  if (!choices.includes(value as T)) return invalid();
  return value as T;
}
function boolean(value: unknown): boolean { return typeof value === 'boolean' ? value : invalid(); }
export function parseResolution(value: unknown): PointFormatResolution {
  const r = record(value) ?? invalid();
  const result: PointFormatResolution = {
    resolved: boolean(r.resolved), format_mode: oneOf(r.format_mode, ['legacy', 'inherit', 'override']),
    configured_format: oneOf(r.configured_format, ['', ...wordFormats]), effective_format: oneOf(r.effective_format, ['', ...wordFormats]),
    effective_source: oneOf(r.effective_source, ['point', 'connection', 'native', 'rule', 'override']),
    effective_policy: oneOf(r.effective_policy, ['legacy_modbus', 'legacy_mc3e', 'legacy_fatek_components', 'inherit_scalar', 'override_scalar']),
    native_compatibility: boolean(r.native_compatibility), applicable: boolean(r.applicable), format_version: oneOf(r.format_version, [0, 1]),
    point_revision: text(r.point_revision), workspace_revision: text(r.workspace_revision, false),
    rule_revision: text(r.rule_revision), connection_revision: text(r.connection_revision),
  };
  if (!result.resolved) return invalid();
  if (result.format_mode !== 'override' && result.configured_format !== '') return invalid();
  if (result.format_mode === 'override' && (!result.configured_format || !result.applicable || result.effective_policy !== 'override_scalar' || result.effective_format !== result.configured_format)) return invalid();
  return result;
}
export function responseData(value: unknown): Record<string, unknown> {
  const envelope = record(value);
  return envelope?.success === true ? record(envelope.data) ?? invalid() : invalid();
}
export function parseSnapshot(value: unknown, scope: PointFormatScope): PointFormatSnapshot {
  const r = record(value) ?? invalid();
  const result: PointFormatSnapshot = {
    workspace_id: text(r.workspace_id, false), device_id: text(r.device_id, false), point_id: text(r.point_id, false), rule_id: text(r.rule_id),
    workspace_revision: text(r.workspace_revision, false), workspace_revision_origin: oneOf(r.workspace_revision_origin, ['initial', 'persisted']),
    rule_revision: text(r.rule_revision), point_revision: text(r.point_revision), connection_revision: text(r.connection_revision),
    editable: boolean(r.editable), resolution: parseResolution(r.resolution),
  };
  for (const key of ['workspace_id', 'device_id', 'point_id'] as const) if (result[key] !== scope[key]) return invalid();
  if (scope.rule_id !== undefined && result.rule_id !== scope.rule_id) return invalid();
  for (const key of ['workspace_revision', 'point_revision', 'rule_revision', 'connection_revision'] as const) if (result[key] !== result.resolution[key]) return invalid();
  return result;
}
export function parseMutation(value: unknown, op: PointFormatOperation): PointFormatMutationResult {
  const r = record(value) ?? invalid(); const p = record(r.point) ?? invalid();
  const point: PointFormatMutationResult['point'] = {
    id: text(p.id, false), device_id: text(p.device_id, false), revision_id: text(p.revision_id, false),
    format_mode: oneOf(p.format_mode, ['inherit', 'override']), format_override: oneOf(p.format_override, ['', ...wordFormats]),
    format_version: oneOf(p.format_version, [1]), data_format: text(p.data_format),
  };
  const resolution = parseResolution(r.resolution);
  const workspace_revision = text(r.workspace_revision, false);
  if (resolution.rule_revision !== (op.expected_rule_revision ?? '')) return invalid();
  if (point.id !== op.point_id || point.device_id !== op.device_id || point.format_mode !== op.format_mode || point.format_override !== op.configured_format || resolution.format_mode !== op.format_mode || resolution.configured_format !== op.configured_format || resolution.point_revision !== point.revision_id || resolution.workspace_revision !== workspace_revision || !resolution.resolved) return invalid();
  return { point, resolution, workspace_revision, runtime_apply_status: oneOf(r.runtime_apply_status, ['pending']), invalidation_status: oneOf(r.invalidation_status, ['pending', 'invalidated']) };
}
