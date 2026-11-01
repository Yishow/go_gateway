import { safeRuntimeIdentity } from '../utils/runtimeOperatorSafety';
import { buildSimpleRecordingTarget } from './studioRecordingView';
import type { RuntimeLogFilters } from '../types/runtimeLogs';
type Navigation = { kind: 'global' } | { kind: 'invalid' } |
  { kind: 'recording'; filters: Pick<RuntimeLogFilters, 'workspaceId' | 'groupId'>; returnTarget: string };
/** Group revision is a return locator, not a log revision filter or current-health proof. */
export function readRuntimeLogsNavigation(search = '', hash = ''): Navigation {
  if (!search || search === '?') return { kind: 'global' };
  const invalid = { kind: 'invalid' } as const;
  if (search.length > 2048 || /[;#]/.test(search) || hash && !hash.startsWith('#')) return invalid;
  try { for (const pair of search.replace(/^\?/, '').split('&')) for (const part of pair.split('=')) decodeURIComponent(part.replace(/\+/g, ' ')); }
  catch { return invalid; }
  const params = new URLSearchParams(search), keys = ['entry', 'workspace_id', 'group_id', 'group_revision'];
  if (params.size !== keys.length || [...params.keys()].some(key => !keys.includes(key) || params.getAll(key).length !== 1) || params.get('entry') !== 'recording') return invalid;
  const workspaceId = params.get('workspace_id'), groupId = params.get('group_id'), groupRevision = params.get('group_revision');
  if (!workspaceId || !groupId || !groupRevision || [workspaceId, groupId, groupRevision].some(value => safeRuntimeIdentity(value) !== value)) return invalid;
  const returnTarget = buildSimpleRecordingTarget({ kind: 'edit', workspaceId, groupId, groupRevision }, hash);
  return returnTarget ? { kind: 'recording', filters: { workspaceId, groupId }, returnTarget } : invalid;
}
export function buildRecordingLogsTarget(scope: { workspaceId: string; groupId: string; groupRevision: string }, hash = ''): string | null {
  const search = new URLSearchParams({ entry: 'recording', workspace_id: scope.workspaceId, group_id: scope.groupId, group_revision: scope.groupRevision });
  return readRuntimeLogsNavigation(`?${search}`, hash).kind === 'recording' ? `/studio/logs?${search}${hash}` : null;
}
