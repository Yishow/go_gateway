import type { EvidenceResult, EvidenceScope, RuleEvidenceView } from '../types/runtimeRuleEvidence';
import { parseRuntimeRuleEvidence } from './runtimeRuleEvidenceParser';
import { safeRuntimeIdentity as id, safeRuntimeLabel as label, safeRuntimeAddress as address, safeRuntimeTagKey, isSafeRuntimeSensorText, safeRuntimeProtocol, safeRuntimeAlignment } from '../utils/runtimeOperatorSafety';
import type { DataType } from '../types/datalink';
import type {
  StudioV2RuntimeContext, StudioV2RuntimeContextDevice, StudioV2RuntimeSetupContext,
  StudioV2RuntimeSetupDatabaseConfig, StudioV2RuntimeSetupDatabaseTarget,
  StudioV2RuntimeSetupMapping, StudioV2RuntimeSetupSourceRule,
} from '../types/studioV2RuntimeContext';
import type { StudioV2WorkspaceReadinessSummary } from '../types/studioV2WorkspaceReadiness';

type RecordValue = Record<string, unknown>;
const record = (value: unknown): RecordValue => value && typeof value === 'object' && !Array.isArray(value) ? value as RecordValue : {};
const list = (value: unknown): RecordValue[] => Array.isArray(value) ? value.map(record) : [];
const text = (value: unknown, pattern: RegExp): string => typeof value === 'string' && pattern.test(value) ? value : '';
const number = (value: unknown): number => typeof value === 'number' && Number.isFinite(value) && value >= 0 ? value : 0;
const choice = <T extends string>(value: unknown, allowed: readonly T[], fallback: T): T => allowed.find((item) => item === value) ?? fallback;
const datatype = (value: unknown): DataType | '' => choice(value, ['bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string', ''], '');
function device(value: RecordValue): StudioV2RuntimeContextDevice {
  return {
    device_id: id(value.device_id), name: label(value.name), protocol: safeRuntimeProtocol(value.protocol),
    running: value.running === true && Boolean(id(value.device_id)) && value.availability_status === 'available' && safeRuntimeProtocol(value.protocol) !== 'unavailable',
    availability_status: choice(value.availability_status, ['available', 'unavailable'], 'unavailable'),
    projection_alignment: safeRuntimeAlignment(value.projection_alignment),
    runtime_projection_version: id(value.runtime_projection_version), workspace_projection_version: id(value.workspace_projection_version),
  };
}
function rule(value: RecordValue): StudioV2RuntimeSetupSourceRule {
  return {
    id: id(value.id), device_id: id(value.device_id), start_address: address(value.start_address), count: number(value.count),
    data_type: datatype(value.data_type), naming_prefix: label(value.naming_prefix), enabled: value.enabled === true,
    revision_id: id(value.revision_id), target_data_type: value.target_data_type == null ? null : datatype(value.target_data_type),
  };
}
function mapping(value: RecordValue): StudioV2RuntimeSetupMapping {
  return {
    id: id(value.id), rule_id: id(value.rule_id), device_id: id(value.device_id), point_id: id(value.point_id), address: address(value.address),
    tag_id: id(value.tag_id), tag_key: safeRuntimeTagKey(value.tag_key), display_name: label(value.display_name),
    unit: typeof value.unit === 'string' && isSafeRuntimeSensorText(value.unit) ? text(value.unit, /^[\p{L}\p{N}°%/²³µΩ]{0,32}$/u) : '', target_type: datatype(value.target_type),
    enabled: value.enabled === true, status: choice(value.status, ['draft', 'validated', 'active', 'out_of_sync', 'error', ''], ''),
  };
}
function database(value: unknown): StudioV2RuntimeSetupDatabaseConfig | null {
  if (value == null) return null;
  const v = record(value);
  return {
    id: id(v.id), workspace_id: id(v.workspace_id), kind: choice(v.kind, ['sqlite', 'postgres', 'mysql', 'sqlserver', 'unavailable'], 'unavailable'),
    name: label(v.name), schema: label(v.schema), table: label(v.table),
    write_mode: choice(v.write_mode, ['insert', 'upsert', 'unavailable'], 'unavailable'), write_interval_seconds: number(v.write_interval_seconds),
    timestamp_column: label(v.timestamp_column), status: choice(v.status, ['ready', 'draft', 'disabled', 'error', 'unavailable'], 'unavailable'),
  };
}
function target(v: RecordValue): StudioV2RuntimeSetupDatabaseTarget {
  return {
    id: id(v.id), tag_id: id(v.tag_id), connector_id: id(v.connector_id), table_schema: label(v.table_schema),
    table_name: label(v.table_name), column_name: label(v.column_name), write_mode: choice(v.write_mode, ['insert', 'upsert', 'unavailable'], 'unavailable'),
    timestamp_column: label(v.timestamp_column), group_key: id(v.group_key),
    write_interval_seconds: v.write_interval_seconds == null ? null : number(v.write_interval_seconds), enabled: v.enabled === true,
  };
}
function readiness(value: unknown): StudioV2WorkspaceReadinessSummary | null {
  if (value == null) return null;
  const v = record(value);
  return {
    ready: v.ready === true, blocking_count: number(v.blocking_count), warning_count: number(v.warning_count),
    issues: list(v.issues).filter((issue) => ['Step 1', 'Step 2', 'Step 3', 'Step 4'].includes(String(issue.step))).map((issue) => ({
      code: id(issue.code), severity: choice(issue.severity, ['blocking', 'warning'], 'blocking'),
      step: choice(issue.step, ['Step 1', 'Step 2', 'Step 3', 'Step 4'], 'Step 1'), scope: id(issue.scope),
      // Legacy messages may contain driver exceptions. Render typed code/action instead.
      message: '',
    })),
  };
}
function setup(value: unknown): StudioV2RuntimeSetupContext | null {
  if (value == null) return null;
  const v = record(value);
  return { evidence_state: choice<'saved-only' | 'unavailable'>(v.evidence_state, ['saved-only', 'unavailable'], 'unavailable'), readiness_summary: readiness(v.readiness_summary), source_rules: list(v.source_rules).map(rule),
    mappings: list(v.mappings).map(mapping), database_config: database(v.database_config), database_targets: list(v.database_targets).map(target) };
}
/** Legacy runtime compatibility projection, not an effective-config authority. Unknown fields are never copied. */
export function projectRuntimeContext(value: unknown): StudioV2RuntimeContext {
  const v = record(value);
  if (!id(v.workspace_id)) throw new Error('runtime-context-unavailable');
  const devices = list(v.devices).map(device);
  const defaultDevice = devices.find((item) => item.device_id && item.device_id === id(v.default_device_id) && item.availability_status === 'available' && item.protocol !== 'unavailable');
  return { workspace_id: id(v.workspace_id), devices,
    default_device_id: defaultDevice?.device_id ?? null, setup: setup(v.setup) };
}

/** Only the caller's requested scope may bind effective evidence, never legacy response metadata. */
export function projectRuntimeRuleEvidence(value: unknown, scope: EvidenceScope): EvidenceResult<RuleEvidenceView> {
  const context = record(value);
  if (!id(context.workspace_id)) return { ok: false, reason: 'unavailable', action: 'refresh' };
  if (context.workspace_id !== scope.workspace_id) return { ok: false, reason: 'foreign-scope', action: 'open-workspace-list' };
  const parsed = parseRuntimeRuleEvidence(context.rule_evidence, scope);
  // Apply the same operator text boundary as legacy projections after the closed-schema parser.
  // This visits detached, typed facts only, never arbitrary connector payloads or getters.
  const safeText = (item: unknown): boolean => {
    if (typeof item === 'string') return isSafeRuntimeSensorText(item);
    if (Array.isArray(item)) return item.every(safeText);
    return item === null || typeof item !== 'object' || Object.values(item).every(safeText);
  };
  return parsed.ok && !safeText(parsed.value) ? { ok: false, reason: 'unavailable', action: 'refresh' } : parsed;
}
