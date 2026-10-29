import type { EvidenceScope } from '../../../types/runtimeRuleEvidence';
import type { RuntimeSourceSelection } from '../../../types/studioV2RuntimeContext';
import { safeRuntimeIdentity } from '../../../utils/runtimeOperatorSafety';

interface RuntimeEntityHistoryState {
  runtimeEntityHistory: EvidenceScope & { version: 1; selection: RuntimeSourceSelection | null };
}
function record(value: unknown, keys: string[]): Record<string, unknown> | null {
  if (!value || typeof value !== 'object' || Array.isArray(value) || ![Object.prototype, null].includes(Object.getPrototypeOf(value))) return null;
  for (const key of Reflect.ownKeys(value)) {
    if (typeof key !== 'string' || !keys.includes(key) || !('value' in Object.getOwnPropertyDescriptor(value, key)!)) return null;
  }
  return value as Record<string, unknown>;
}
function identity(value: unknown): value is string { return typeof value === 'string' && value !== '' && safeRuntimeIdentity(value) === value; }
function validScope(value: unknown): value is EvidenceScope {
  const scope = record(value, ['workspace_id', 'workspace_revision', 'device_id']);
  return !!scope && identity(scope.workspace_id) && identity(scope.workspace_revision) && identity(scope.device_id);
}
function selection(value: unknown): RuntimeSourceSelection | null {
  const item = record(value, ['rule_id', 'rule_revision', 'point_id', 'mapping_id']);
  if (!item || !identity(item.rule_id) || !identity(item.rule_revision) || !identity(item.point_id) || (item.mapping_id !== undefined && !identity(item.mapping_id))) return null;
  return { rule_id: item.rule_id, rule_revision: item.rule_revision, point_id: item.point_id, ...(item.mapping_id === undefined ? {} : { mapping_id: item.mapping_id as string }) };
}
/** Only a navigation identity hint. The scoped source reader must revalidate it on return. */
export function createRuntimeEntityHistory(scope: EvidenceScope, value: RuntimeSourceSelection | null): RuntimeEntityHistoryState | null {
  if (!validScope(scope)) return null;
  const selected = value === null ? null : selection(value);
  if (value !== null && selected === null) return null;
  return { runtimeEntityHistory: { version: 1, workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, device_id: scope.device_id, selection: selected } };
}
/** Never restore facts, URL targets, stale scopes, session storage or a guessed first mapping. */
export function readRuntimeEntityHistory(value: unknown, scope: EvidenceScope | null): RuntimeSourceSelection | null {
  if (!validScope(scope)) return null;
  const root = record(value, ['runtimeEntityHistory']);
  const saved = record(root?.runtimeEntityHistory, ['version', 'workspace_id', 'workspace_revision', 'device_id', 'selection']);
  if (!saved || saved.version !== 1 || saved.workspace_id !== scope.workspace_id || saved.workspace_revision !== scope.workspace_revision || saved.device_id !== scope.device_id) return null;
  return selection(saved.selection);
}
