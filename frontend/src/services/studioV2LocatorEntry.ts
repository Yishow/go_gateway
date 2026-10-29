import { DATALINK_DATA_TYPES, type Device, type SourceRuleRecord, type StudioV2WorkspaceMappingRecord, type StudioV2WorkspaceDatabaseConfigRecord, type StudioV2WorkspaceDatabaseTargetRecord } from '../types/datalink';
import type { EntityLocator } from '../types/runtimeRuleEvidence';
import type { StudioLocatorEntryResult, StudioV2LocatorEntrySnapshot } from '../types/studioV2LocatorEntry';
import type { StudioV2Workspace } from '../types/studioV2Workspace';
import { resolveStudioLocatorEntrySelection } from '../features/datalink/workbench-v2/state/studioLocatorEntry';
import { parseBoundedJson } from '../utils/safeJson';
import { parseWriteGroupListData } from '../utils/studioV2WriteGroupJson';
import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';
import { buildStudioEntityLocatorTarget, parseStudioEntityLocator } from './studioEntityLocator';
import { studioV2DatalinkApi } from './studioV2Workspace';

type RecordValue = Record<string, unknown>;
type Predicate = (value: unknown) => boolean;
export interface StudioLocatorEditorExpectation { workspaceRevision: string; groupRevision: string }
type ReadOptions = { editor?: boolean; expected?: StudioLocatorEditorExpectation };
class Rejected extends Error { constructor(readonly reason: 'unavailable' | 'stale-reference' | 'foreign-scope' = 'unavailable') { super(reason); } }
function reject(reason?: 'stale-reference' | 'foreign-scope'): never { throw new Rejected(reason); }
function object(value: unknown, keys: readonly string[]): RecordValue {
  if (!value || typeof value !== 'object' || Array.isArray(value) || Object.getPrototypeOf(value) !== Object.prototype || Object.keys(value).some(key => !keys.includes(key))) reject();
  return value as RecordValue;
}
const text: Predicate = value => typeof value === 'string';
const boolean: Predicate = value => typeof value === 'boolean';
const number: Predicate = value => typeof value === 'number' && Number.isFinite(value);
const identity: Predicate = value => !!safeRuntimeIdentity(value);
const nullable = (check: Predicate): Predicate => value => value === null || check(value);
const oneOf = (values: readonly unknown[]): Predicate => value => values.includes(value);
const dataType = oneOf(DATALINK_DATA_TYPES);
const strings: Predicate = value => Array.isArray(value) && value.every(text);
function record(value: unknown, required: Record<string, Predicate>, optional: Record<string, Predicate> = {}): RecordValue {
  const result = object(value, [...Object.keys(required), ...Object.keys(optional)]);
  if (Object.entries(required).some(([key, check]) => !Object.hasOwn(result, key) || !check(result[key])) ||
      Object.entries(optional).some(([key, check]) => Object.hasOwn(result, key) && !check(result[key]))) reject();
  return result;
}
function ids(value: unknown): string[] {
  if (!Array.isArray(value) || !value.every(identity) || new Set(value).size !== value.length) reject();
  return [...value] as string[];
}
const texts = (...keys: string[]) => Object.fromEntries(keys.map(key => [key, text]));
function device(value: unknown): Device {
  return record(value, { ...texts('name', 'description', 'connection_config', 'created_at', 'updated_at'), last_test_error: oneOf(['']), id: identity,
    protocol: oneOf(['modbus_tcp', 'modbus_rtu', 'modbus_udp', 'fatek_fbs', 'mc_3e', 'mqtt']), status: oneOf(['draft', 'active', 'disabled']),
    last_test_at: nullable(text), last_test_success: nullable(boolean) }, { readiness_status: oneOf(['ready', 'not_ready', 'error']),
    availability_status: oneOf(['available', 'unavailable']), availability_reason: oneOf(['']), running: boolean }) as unknown as Device;
}
function rule(value: unknown): SourceRuleRecord {
  return record(value, { ...texts('start_address', 'naming_prefix', 'created_at', 'updated_at'), id: identity, workspace_id: identity, device_id: identity,
    revision_id: identity, count: v => number(v) && Number.isInteger(v) && (v as number) > 0, data_type: dataType, enabled: boolean, locked: boolean,
    origin: oneOf(['manual', 'template']), skipped_addresses: strings }, { template_name: text, save_state: oneOf(['saved']), target_data_type: dataType,
    scale_multiplier: number, scale_offset: number, data_format: text, share_enabled: boolean, share_start_register: nullable(number), share_stride: nullable(number) }) as unknown as SourceRuleRecord;
}
function mapping(value: unknown): StudioV2WorkspaceMappingRecord {
  return record(value, { ...texts('address', 'tag_key', 'display_name', 'unit', 'created_at', 'updated_at'), id: identity, workspace_id: identity, device_id: identity,
    point_id: identity, rule_id: identity, tag_id: identity, target_type: dataType, scale: number, offset: number, enabled: boolean }, { save_state: oneOf(['saved']) }) as unknown as StudioV2WorkspaceMappingRecord;
}
function collection<T>(input: unknown, parse: (value: unknown) => T, editor = false): T[] {
  if (!Array.isArray(input) || (!editor && input.length !== 1)) reject();
  return input.map(parse);
}
function databaseConfig(input: unknown): StudioV2WorkspaceDatabaseConfigRecord {
  const raw = object(input, ['id', 'identity_revision', 'setup_revision', 'workspace_id', 'kind', 'name', 'host', 'port', 'database', 'username', 'schema', 'table',
    'write_mode', 'write_interval_seconds', 'timestamp_column', 'status', 'row_groups', 'save_state', 'created_at', 'updated_at',
    'last_schema_ensure_at', 'last_schema_ensure_status', 'last_schema_ensure_error', 'last_write_at', 'last_write_status', 'last_write_error', 'last_flush_at', 'last_flush_status', 'last_flush_error']);
  const copy = { ...raw };
  for (const key of Object.keys(copy).filter(key => key.startsWith('last_'))) {
    if (key.endsWith('_error') ? copy[key] !== '' : !nullable(text)(copy[key])) reject(); delete copy[key];
  }
  if (copy.row_groups !== undefined) {
    if (!Array.isArray(copy.row_groups)) reject();
    copy.row_groups = copy.row_groups.map(value => record(value, { id: identity, table_name: text, member_point_ids: strings },
      { connector_id: identity, table_schema: text, group_key_columns: strings, unique_key_columns: strings }));
  }
  return record(copy, { ...texts('name', 'host', 'database', 'username', 'schema', 'table', 'timestamp_column', 'status', 'created_at', 'updated_at'),
    id: identity, workspace_id: identity, kind: oneOf(['sqlite', 'postgres', 'mysql', 'sqlserver']), port: number, write_mode: oneOf(['insert', 'upsert']), write_interval_seconds: number },
    { identity_revision: identity, setup_revision: text, row_groups: Array.isArray, save_state: oneOf(['saved']) }) as unknown as StudioV2WorkspaceDatabaseConfigRecord;
}
function databaseTarget(input: unknown): StudioV2WorkspaceDatabaseTargetRecord {
  return record(input, { ...texts('column_name', 'created_at', 'updated_at'), id: identity, workspace_id: identity, point_id: identity, tag_id: identity, enabled: boolean },
    { row_group_id: identity, setup_revision: text, save_state: oneOf(['saved']) }) as unknown as StudioV2WorkspaceDatabaseTargetRecord;
}
function heads(input: unknown): Record<string, string> {
  if (!input || typeof input !== 'object' || Array.isArray(input)) reject();
  const entries = Object.entries(input);
  if (entries.some(([id, revision]) => !identity(id) || !identity(revision))) reject();
  return Object.fromEntries(entries) as Record<string, string>;
}
function equalIDs(actual: string[], expected: string[]) { if (new Set(actual).size !== actual.length || actual.length !== expected.length || actual.some(id => !expected.includes(id))) reject('stale-reference'); }
function parseEntry(raw: unknown, locator: EntityLocator, expected?: StudioLocatorEditorExpectation): { workspace: StudioV2Workspace; entry: StudioV2LocatorEntrySnapshot } {
  const data = object(parseBoundedJson(raw, 16, path => path.at(-1) === 'connection_config' ? 65536 : undefined),
    ['id', 'kind', 'status', 'ordered_device_ids', 'created_at', 'updated_at', 'readiness_summary', 'locator_entry']);
  if (Object.hasOwn(data, 'readiness_summary') && data.readiness_summary !== null) reject();
  if (!identity(data.id) || !oneOf(['single'])(data.kind) || !oneOf(['empty', 'ready'])(data.status) || !text(data.created_at) || !text(data.updated_at)) reject();
  if (data.id !== locator.workspace_id) reject('foreign-scope');
  const ordered = ids(data.ordered_device_ids);
  const rawEntry = object(data.locator_entry, ['workspace_revision', 'hydration_scope', 'devices', 'source_rules', 'mappings', 'write_groups', 'rule_heads', 'group_heads', 'revalidated_mapping_ids', ...(expected ? ['editor_hydrated', 'database_config', 'database_targets'] : [])]);
  if (typeof rawEntry.workspace_revision !== 'string' || (rawEntry.workspace_revision !== '' && !identity(rawEntry.workspace_revision))) reject();
  const scope = object(rawEntry.hydration_scope, ['complete_workspace', 'device_ids', 'rule_ids', 'mapping_ids', 'group_ids']);
  if (scope.complete_workspace !== false) reject();
  const entry: StudioV2LocatorEntrySnapshot = { workspace_revision: rawEntry.workspace_revision,
    hydration_scope: { complete_workspace: false, device_ids: ids(scope.device_ids), rule_ids: ids(scope.rule_ids), mapping_ids: ids(scope.mapping_ids), group_ids: ids(scope.group_ids) },
    rule_heads: heads(rawEntry.rule_heads), group_heads: heads(rawEntry.group_heads), revalidated_mapping_ids: ids(rawEntry.revalidated_mapping_ids) };
  if (expected && (rawEntry.editor_hydrated !== true || !Object.hasOwn(rawEntry, 'database_config') || !Object.hasOwn(rawEntry, 'database_targets'))) reject();
  const needed = { devices: !!expected || locator.device_id !== null, source_rules: !!expected || locator.step === 2 || locator.step === 3, mappings: !!expected || locator.step === 3, write_groups: locator.step === 4 };
  for (const [key, required] of Object.entries(needed)) if (Object.hasOwn(rawEntry, key) !== required) reject();
  if (needed.devices) entry.devices = collection(rawEntry.devices, device, !!expected);
  if (needed.source_rules) entry.source_rules = collection(rawEntry.source_rules, rule, !!expected);
  if (needed.mappings) entry.mappings = collection(rawEntry.mappings, mapping, !!expected);
  if (needed.write_groups) {
    object(rawEntry.write_groups, ['workspace_id', 'workspace_revision', 'groups']);
    const groups = parseWriteGroupListData(rawEntry.write_groups); if (!groups || groups.groups.length !== 1) reject(); entry.write_groups = groups;
  }
  const deviceIDs = entry.devices?.map(item => item.id) ?? [];
  const ruleIDs = entry.source_rules?.map(item => item.id) ?? [];
  const mappingIDs = entry.mappings?.map(item => item.id) ?? [];
  const groupIDs = entry.write_groups?.groups.map(item => item.id) ?? [];
  equalIDs(entry.hydration_scope.device_ids, deviceIDs); equalIDs(entry.hydration_scope.rule_ids, ruleIDs);
  equalIDs(entry.hydration_scope.mapping_ids, mappingIDs); equalIDs(entry.hydration_scope.group_ids, groupIDs);
  if (!expected) { equalIDs(deviceIDs, locator.device_id === null ? [] : [locator.device_id]); equalIDs(mappingIDs, locator.mapping_id === null ? [] : [locator.mapping_id]); }
  equalIDs(groupIDs, locator.group_id === null ? [] : [locator.group_id]);
  equalIDs(ruleIDs, expected ? [...new Set(entry.mappings!.map(item => item.rule_id))] : locator.step === 2 ? [locator.rule_id!] : locator.step === 3 ? [entry.mappings![0].rule_id] : []);
  equalIDs(Object.keys(entry.rule_heads), ruleIDs); equalIDs(Object.keys(entry.group_heads), groupIDs); equalIDs(entry.revalidated_mapping_ids, mappingIDs);
  for (const row of [...(entry.source_rules ?? []), ...(entry.mappings ?? []), ...(entry.write_groups?.groups ?? [])]) if (row.workspace_id !== locator.workspace_id) reject('foreign-scope');
  if (entry.devices?.some(item => !ordered.includes(item.id)) || entry.source_rules?.some(item => (!expected && item.device_id !== locator.device_id) || item.revision_id !== entry.rule_heads[item.id]) ||
      entry.mappings?.some(item => !expected && item.device_id !== locator.device_id) || entry.write_groups?.groups.some(item => item.revision !== entry.group_heads[item.id]) ||
      (entry.write_groups && (entry.write_groups.workspace_id !== data.id || entry.write_groups.workspace_revision !== entry.workspace_revision))) reject('stale-reference');
  if (expected) {
    const group = entry.write_groups!.groups[0];
    if (entry.workspace_revision !== expected.workspaceRevision || group.revision !== expected.groupRevision || group.status === 'deleted') reject('stale-reference');
    for (const row of entry.mappings!) {
      if (!group.members.some(member => member.device_id === row.device_id && member.point_id === row.point_id && member.tag_id === row.tag_id) ||
          !entry.source_rules!.some(rule => rule.id === row.rule_id && rule.device_id === row.device_id) || !deviceIDs.includes(row.device_id)) reject('stale-reference');
    }
    const usedDevices = [...new Set([...entry.mappings!.map(row => row.device_id), ...(locator.device_id ? [locator.device_id] : [])])];
    equalIDs(deviceIDs, usedDevices);
    entry.editor_hydrated = true; entry.database_config = rawEntry.database_config === null ? null : databaseConfig(rawEntry.database_config);
    entry.database_targets = collection(rawEntry.database_targets, databaseTarget, true);
    if (entry.database_config && entry.database_config.workspace_id !== locator.workspace_id) reject('foreign-scope');
    if (new Set(entry.database_targets.map(row => row.id)).size !== entry.database_targets.length || entry.database_targets.some(row => row.workspace_id !== locator.workspace_id ||
        !group.members.some(member => member.point_id === row.point_id && member.tag_id === row.tag_id)) || (!entry.database_config && entry.database_targets.length)) reject('stale-reference');
  }
  const workspace: StudioV2Workspace = { id: data.id as string, kind: 'single', status: data.status as StudioV2Workspace['status'],
    ordered_device_ids: ordered, created_at: data.created_at as string, updated_at: data.updated_at as string };
  return { workspace, entry };
}

/** Only this detached entry query reads locator ownership; never writes canonical bootstrap caches. */
export async function readStudioLocatorEntry(locator: EntityLocator, signal?: AbortSignal, options?: ReadOptions): Promise<StudioLocatorEntryResult> {
  try {
    const target = buildStudioEntityLocatorTarget(locator); if (!target) return { ok: false, reason: 'unavailable' };
    const intent = parseStudioEntityLocator(target.slice('/studio/v2'.length)); if (intent.kind !== 'locator') return { ok: false, reason: 'unavailable' };
    const selected = intent.locator;
    const expected = options?.editor ? options.expected : undefined;
    if (options?.editor && (selected.step !== 4 || !expected || typeof expected.workspaceRevision !== 'string' ||
        (expected.workspaceRevision !== '' && !identity(expected.workspaceRevision)) || !identity(expected.groupRevision))) return { ok: false, reason: 'unavailable' };
    const captured = expected ? { workspaceRevision: expected.workspaceRevision, groupRevision: expected.groupRevision } : undefined;
    const params = Object.fromEntries(new URLSearchParams(target.split('?')[1]));
    if (captured) Object.assign(params, { editor: '1', expected_workspace_revision: captured.workspaceRevision, expected_group_revision: captured.groupRevision });
    const response = await studioV2DatalinkApi.get('/studio-v2/workspace', { params, signal });
    if (signal?.aborted || response.data?.success !== true) return { ok: false, reason: 'unavailable' };
    const { workspace, entry } = parseEntry(response.data.data, selected, captured);
    const selection = resolveStudioLocatorEntrySelection(selected, workspace, entry);
    if (selection.state === 'unavailable' && !entry.write_groups?.groups.some(group => group.id === selected.group_id && group.status === 'deleted')) reject('stale-reference');
    return { ok: true, value: { locator: selected, workspace, entry, selection } };
  } catch (error) {
    if (error instanceof Rejected) return { ok: false, reason: error.reason };
    const status = (error as { response?: { status?: unknown } } | null)?.response?.status;
    return { ok: false, reason: status === 403 ? 'foreign-scope' : status === 409 ? 'stale-reference' : 'unavailable' };
  }
}

/** Explicit native Open/Edit hydration; historical revision hints never supply these expectations. */
export function readStudioLocatorEditorEntry(locator: EntityLocator, expected: StudioLocatorEditorExpectation, signal?: AbortSignal): Promise<StudioLocatorEntryResult> {
  return readStudioLocatorEntry(locator, signal, { editor: true, expected });
}
