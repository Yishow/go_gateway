import type { EntityLocator } from '../types/runtimeRuleEvidence';
import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';

export type StudioEntityLocatorIntent = { kind: 'absent' } | { kind: 'invalid' } | { kind: 'locator'; locator: EntityLocator };
const fields = ['workspace_id', 'device_id', 'rule_id', 'mapping_id', 'group_id', 'group_revision'] as const;
const optional = fields.slice(1);
const keys = ['step', ...fields];

/** IDs-only viewing intent; never accepts a destination, payload or return URL. */
export function buildStudioEntityLocatorTarget(locator: EntityLocator): string | null {
  if (!locator || typeof locator !== 'object' || Object.getPrototypeOf(locator) !== Object.prototype ||
      Object.values(Object.getOwnPropertyDescriptors(locator)).some(value => !('value' in value)) || Reflect.ownKeys(locator).some((key) => typeof key !== 'string' || !keys.includes(key)) ||
      ![1, 2, 3, 4].includes(locator.step) || !safeRuntimeIdentity(locator.workspace_id) ||
      optional.some((key) => locator[key] !== null && !safeRuntimeIdentity(locator[key]))) return null;
  const { step, device_id: device, rule_id: rule, mapping_id: mapping, group_id: group, group_revision: revision } = locator;
  if ((step === 1 && (!device || rule || mapping || group || revision)) ||
      (step === 2 && (!device || !rule || mapping || group || revision)) ||
      (step === 3 && (!device || !mapping || rule || group || revision)) ||
      (step === 4 && (!group || rule || mapping))) return null;
  const params = new URLSearchParams({ entry: 'locator', workspace_id: locator.workspace_id, step: String(step) });
  for (const key of optional) if (locator[key] !== null) params.set(key, locator[key]!);
  return `/studio/v2?${params.toString()}`;
}

/** Read raw search, rejecting malformed encoding before URLSearchParams can repair it. */
export function parseStudioEntityLocator(search: string): StudioEntityLocatorIntent {
  if (typeof search !== 'string' || search.length > 2048 || /[;#]/.test(search)) return { kind: 'invalid' };
  const raw = search.startsWith('?') ? search.slice(1) : search;
  try { for (const pair of raw.split('&')) for (const part of pair.split('=')) decodeURIComponent(part.replace(/\+/g, ' ')); }
  catch { return { kind: 'invalid' }; }
  const params = new URLSearchParams(raw);
  const hasIntent = ['entry', 'editor', 'expected_group_revision', 'expected_workspace_revision', 'locator_workspace_id', 'locator_workspace_revision'].some(key => params.has(key)) ||
    fields.some((key) => params.has(key)) || (params.has('step') && params.get('focus') !== 'readiness');
  if (!hasIntent) return { kind: 'absent' };
  if (params.get('entry') !== 'locator' || [...params.keys()].some((key) => !['entry', ...keys].includes(key) || params.getAll(key).length !== 1) ||
      !/^[1-4]$/.test(params.get('step') ?? '')) return { kind: 'invalid' };
  const locator: EntityLocator = { step: Number(params.get('step')) as EntityLocator['step'], workspace_id: params.get('workspace_id') ?? '',
    device_id: params.get('device_id'), rule_id: params.get('rule_id'), mapping_id: params.get('mapping_id'),
    group_id: params.get('group_id'), group_revision: params.get('group_revision') };
  return buildStudioEntityLocatorTarget(locator) ? { kind: 'locator', locator } : { kind: 'invalid' };
}

import type { StudioV2LocatorReadScope } from '../types/studioV2LocatorEntry';
/** Capture only the verified read pair; explicit legacy empty revision differs from missing. */
export function captureStudioLocatorReadScope(scope?: StudioV2LocatorReadScope): StudioV2LocatorReadScope | undefined {
  if (scope === undefined) return undefined;
  if (!scope || typeof scope !== 'object' || Object.getPrototypeOf(scope) !== Object.prototype ||
      Reflect.ownKeys(scope).some(key => key !== 'workspaceId' && key !== 'workspaceRevision') ||
      Object.values(Object.getOwnPropertyDescriptors(scope)).some(value => !('value' in value)) ||
      !safeRuntimeIdentity(scope.workspaceId) || typeof scope.workspaceRevision !== 'string' ||
      (scope.workspaceRevision !== '' && !safeRuntimeIdentity(scope.workspaceRevision))) throw new Error('Scoped settings read unavailable.');
  return { workspaceId: scope.workspaceId, workspaceRevision: scope.workspaceRevision };
}
export function studioLocatorReadParams(scope?: StudioV2LocatorReadScope): Record<string, string> {
  const selected = captureStudioLocatorReadScope(scope);
  return selected ? { locator_workspace_id: selected.workspaceId, locator_workspace_revision: selected.workspaceRevision } : {};
}
export function studioLocatorReadKey(scope?: StudioV2LocatorReadScope): string[] {
  try { const selected = captureStudioLocatorReadScope(scope); return selected ? ['locator-read', selected.workspaceId, selected.workspaceRevision] : []; }
  catch { return ['locator-read', 'invalid']; }
}
/** A durable operation may retain an older revision: only workspace identity is checked here. */
export function assertStudioLocatorWorkspaceBinding(value: unknown, scope?: StudioV2LocatorReadScope): void {
  if (scope && value && typeof value === 'object' && Object.hasOwn(value, 'workspace_id') &&
      (value as { workspace_id: unknown }).workspace_id !== scope.workspaceId) throw new Error('Scoped settings read unavailable.');
}
export function throwStudioLocatorReadError(error: unknown, scope?: StudioV2LocatorReadScope): never {
  if (scope) throw new Error('Scoped settings read unavailable.');
  throw error;
}
