import type { EvidenceResult, EvidenceScope } from '../types/runtimeRuleEvidence';
import type { DataType } from '../types/datalink';
import type { RuntimeSourcePage, RuntimeSourcePageRequest } from '../types/studioV2RuntimeContext';
import { isSafeRuntimeSensorText, safeRuntimeAddress, safeRuntimeIdentity, safeRuntimeLabel, safeRuntimeTagKey } from '../utils/runtimeOperatorSafety';

type FailureReason = 'unavailable' | 'stale-reference' | 'foreign-scope';
class Rejected extends Error {
  constructor(readonly reason: FailureReason = 'unavailable') { super(reason); }
}
function reject(reason?: FailureReason): never { throw new Rejected(reason); }
function id(value: unknown): string { return typeof value === 'string' && value !== '' && safeRuntimeIdentity(value) === value ? value : reject(); }
function safeText(value: unknown): string { return typeof value === 'string' && isSafeRuntimeSensorText(value) ? value : reject(); }
function label(value: unknown): string { return value === '' ? '' : typeof value === 'string' && safeRuntimeLabel(value) === value ? value : reject(); }
function address(value: unknown): string { return typeof value === 'string' && value !== '' && safeRuntimeAddress(value) === value ? value : reject(); }
function boolean(value: unknown): boolean { return typeof value === 'boolean' ? value : reject(); }
function integer(value: unknown, minimum = 0, maximum = Number.MAX_SAFE_INTEGER): number {
  return typeof value === 'number' && Number.isSafeInteger(value) && value >= minimum && value <= maximum ? value : reject();
}
function dataType(value: unknown): DataType | '' {
  return typeof value === 'string' && ['', 'bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string'].includes(value) ? value as DataType | '' : reject();
}
function object(value: unknown, keys: readonly string[]): Record<string, unknown> {
  if (value === null || typeof value !== 'object' || Array.isArray(value) || ![Object.prototype, null].includes(Object.getPrototypeOf(value))) return reject();
  for (const key of Reflect.ownKeys(value)) {
    if (typeof key !== 'string' || !keys.includes(key) || !('value' in Object.getOwnPropertyDescriptor(value, key)!)) return reject();
  }
  return value as Record<string, unknown>;
}
function cursor(value: unknown): string {
  if (typeof value !== 'string' || value.length < 1 || value.length > 1024 || !/^[A-Za-z0-9_-]+$/.test(value)) return reject();
  try {
    const decoded = atob(value.replace(/-/g, '+').replace(/_/g, '/'));
    if (decoded.length > 768 || btoa(decoded).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '') !== value) return reject();
  } catch { return reject(); }
  return value;
}

/** Normalize only approved read parameters; never forward arbitrary caller keys. */
export function normalizeRuntimeSourcePageRequest(request: RuntimeSourcePageRequest): RuntimeSourcePageRequest | null {
  try {
    if (request.kind !== 'rules' && request.kind !== 'mappings') return null;
    if (request.query !== undefined && (typeof request.query !== 'string' || /\p{Cc}/u.test(request.query))) return null;
    const query = request.query === undefined ? '' : safeText(request.query.trim());
    if ([...query].length > 128) return null;
    const limit = request.limit === undefined ? 50 : Math.min(integer(request.limit, 1), 100);
    const filtered = request.rule_id !== undefined || request.rule_revision !== undefined;
    if (filtered && request.kind !== 'mappings') return null;
    return { kind: request.kind, query, limit,
      ...(request.cursor === undefined ? {} : { cursor: cursor(request.cursor) }),
      ...(filtered ? { rule_id: id(request.rule_id), rule_revision: id(request.rule_revision) } : {}),
    };
  } catch { return null; }
}

function parseRule(value: unknown): RuntimeSourcePage['rules'][number] {
  const v = object(value, ['id', 'device_id', 'start_address', 'count', 'data_type', 'naming_prefix', 'enabled', 'revision_id', 'target_data_type']);
  return { id: id(v.id), device_id: id(v.device_id), start_address: address(v.start_address), count: integer(v.count),
    data_type: dataType(v.data_type), naming_prefix: label(v.naming_prefix), enabled: boolean(v.enabled), revision_id: id(v.revision_id),
    ...(v.target_data_type === undefined ? {} : { target_data_type: v.target_data_type === null ? null : dataType(v.target_data_type) }),
  };
}
function parseMapping(value: unknown): RuntimeSourcePage['mappings'][number] {
  const v = object(value, ['id', 'rule_id', 'rule_revision', 'device_id', 'point_id', 'address', 'tag_id', 'tag_key', 'display_name', 'unit', 'target_type', 'enabled', 'status']);
  const result: RuntimeSourcePage['mappings'][number] = { rule_id: id(v.rule_id), rule_revision: id(v.rule_revision), device_id: id(v.device_id), point_id: id(v.point_id), address: address(v.address), enabled: boolean(v.enabled) };
  if (v.id !== undefined) result.id = id(v.id);
  if (v.tag_id !== undefined) result.tag_id = id(v.tag_id);
  if (v.tag_key !== undefined) result.tag_key = v.tag_key === '' ? '' : typeof v.tag_key === 'string' && safeRuntimeTagKey(v.tag_key) === v.tag_key ? v.tag_key : reject();
  if (v.display_name !== undefined) result.display_name = label(v.display_name);
  if (v.unit !== undefined) result.unit = typeof v.unit === 'string' && /^[\p{L}\p{N}°%/²³µΩ]{0,32}$/u.test(v.unit) ? safeText(v.unit) : reject();
  if (v.target_type !== undefined) result.target_type = dataType(v.target_type);
  if (v.status !== undefined) result.status = typeof v.status === 'string' && ['', 'draft', 'validated', 'active', 'out_of_sync', 'error'].includes(v.status) ? v.status : reject();
  return result;
}

/** Closed, bounded saved-summary projection tied to the caller's scope and page request. */
export function parseRuntimeSourcePage(value: unknown, scope: EvidenceScope, request: RuntimeSourcePageRequest): EvidenceResult<RuntimeSourcePage> {
  try {
    id(scope.workspace_id); id(scope.workspace_revision); id(scope.device_id);
    const selected = normalizeRuntimeSourcePageRequest(request);
    if (selected === null) return reject();
    const v = object(value, ['contract_version', 'workspace_id', 'workspace_revision', 'device_id', 'kind', 'query', 'limit', 'list_revision', 'rules', 'mappings', 'next_cursor']);
    if (v.contract_version !== 1 || (v.kind !== 'rules' && v.kind !== 'mappings')) return reject();
    const workspace_id = id(v.workspace_id); const workspace_revision = id(v.workspace_revision); const device_id = id(v.device_id);
    if (workspace_id !== scope.workspace_id || device_id !== scope.device_id) return reject('foreign-scope');
    const query = safeText(v.query); const limit = integer(v.limit, 1, 100);
    if (workspace_revision !== scope.workspace_revision || v.kind !== selected.kind || query !== selected.query || limit !== selected.limit) return reject('stale-reference');
    if (!Array.isArray(v.rules) || !Array.isArray(v.mappings) || v.rules.length + v.mappings.length > limit ||
      (selected.kind === 'rules' ? v.mappings.length > 0 : v.rules.length > 0)) return reject();
    const rules = v.rules.map(parseRule); const mappings = v.mappings.map(parseMapping);
    const ruleIDs = new Set<string>(); const mappingIDs = new Set<string>(); const bindings = new Set<string>(); const revisions = new Map<string, string>();
    for (const row of [...rules, ...mappings]) if (row.device_id !== scope.device_id) return reject('foreign-scope');
    for (const row of rules) { if (ruleIDs.has(row.id)) return reject(); ruleIDs.add(row.id); }
    for (const row of mappings) {
      const key = `${row.rule_id}:${row.point_id}:${row.id ?? ''}`;
      if (bindings.has(key) || (row.id !== undefined && mappingIDs.has(row.id))) return reject();
      bindings.add(key); if (row.id !== undefined) mappingIDs.add(row.id);
      if (revisions.has(row.rule_id) && revisions.get(row.rule_id) !== row.rule_revision) return reject('stale-reference');
      revisions.set(row.rule_id, row.rule_revision);
      if (selected.rule_id !== undefined && (row.rule_id !== selected.rule_id || row.rule_revision !== selected.rule_revision)) return reject('stale-reference');
    }
    const next_cursor = v.next_cursor === null ? null : cursor(v.next_cursor);
    if (next_cursor !== null && next_cursor === selected.cursor) return reject();
    return { ok: true, value: { contract_version: 1, workspace_id, workspace_revision, device_id, kind: selected.kind, query, limit,
      list_revision: id(v.list_revision), rules, mappings, next_cursor } };
  } catch (error) {
    const reason = error instanceof Rejected ? error.reason : 'unavailable';
    return { ok: false, reason, action: reason === 'foreign-scope' ? 'open-workspace-list' : 'refresh' };
  }
}
