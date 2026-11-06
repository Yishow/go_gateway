import { useCallback, useSyncExternalStore } from 'react';
import { useQueryClient, type QueryClient, type Query } from '@tanstack/react-query';
import { recordingHomeSourceContext, type HomeSourceContext } from './recordingHomeSourceCache';
import { captureRecordingHomeSourceMutation, completeRecordingHomeSourceMutation } from './recordingHomeSourceMutation';
import { studioV2WorkspaceKeys } from './keys';
import type { StudioV2WorkspaceDatabaseConfigRecord } from '../../types/datalink';
import type { StudioV2WorkspaceDatabaseConfigRequest } from '../../services/studioV2WorkspaceDatabase';

export type ConfigReadContext = ReturnType<typeof recordingHomeSourceContext> & { workspaceId?: string };
interface Entry { generation: number; confirmed: boolean; pending?: Promise<void>; status: 'idle' | 'refreshing' | 'fresh' | 'stale-error' }
interface Registry { entries: Map<string, Entry>; version: number; listeners: Set<() => void> }
const registries = new WeakMap<QueryClient, Registry>();
function registry(client: QueryClient) {
  let value = registries.get(client);
  if (!value) { value = { entries: new Map(), version: 0, listeners: new Set() }; registries.set(client, value); }
  return value;
}
function id(context: ConfigReadContext) { return JSON.stringify([context.backendToken, context.apiOrigin, context.workspaceId ?? '']); }
function entry(client: QueryClient, context: ConfigReadContext) {
  const entries = registry(client).entries, key = id(context);
  let value = entries.get(key);
  if (!value) { value = { generation: 0, confirmed: false, status: 'idle' }; entries.set(key, value); }
  return value;
}
function notify(client: QueryClient) { const value = registry(client); value.version++; value.listeners.forEach(listener => listener()); }
export function configReadContext(workspaceId?: string, options: HomeSourceContext = {}): ConfigReadContext {
  return { ...recordingHomeSourceContext(options), workspaceId };
}
export function configReadKey(prefix: readonly unknown[], context: ConfigReadContext) {
  return [...prefix, 'api-context', context.backendToken, context.apiOrigin, context.workspaceId ?? ''];
}
export function configReadMatches(query: Query, context: ConfigReadContext) {
  const key = query.queryKey, index = key.indexOf('api-context');
  if (index < 0) return true; // Existing canonical observers still participate in invalidation.
  return key[index + 1] === context.backendToken && key[index + 2] === context.apiOrigin &&
    (!context.workspaceId || !key[index + 3] || key[index + 3] === context.workspaceId);
}
const affected = (query: Query) => ['database-config', 'database-targets', 'bootstrap'].includes(String(query.queryKey[1])) && query.queryKey[0] === studioV2WorkspaceKeys.all[0];
function filters(context: ConfigReadContext) { return { predicate: (query: Query) => affected(query) && configReadMatches(query, context) }; }
/** A read retains its original origin and generation even if transport ignores cancellation. */
export async function fencedConfigRead<T>(client: QueryClient, context: ConfigReadContext, signal: AbortSignal, read: () => Promise<T>): Promise<T> {
  const own = entry(client, context), generation = own.generation;
  const result = await read();
  if (signal.aborted || own.generation !== generation ||
    recordingHomeSourceContext({ backend: context.backend }).apiOrigin !== context.apiOrigin) throw new Error('Configuration observation superseded');
  return result;
}
/** Capture before any await. The returned generation fences this submitted operation only. */
export function beginConfigSave(client: QueryClient, workspaceId?: string) {
  const home = captureRecordingHomeSourceMutation(client, workspaceId);
  const context = configReadContext(home?.workspaceId ?? workspaceId, home), current = entry(client, context);
  current.generation++; current.status = 'refreshing'; current.pending = undefined;
  const global = entry(client, { ...context, workspaceId: undefined });
  if (global !== current) global.generation++;
  void client.cancelQueries(filters(context));
  void client.invalidateQueries({ ...filters(context), refetchType: 'none' });
  notify(client);
  return { context, generation: current.generation, home };
}
export function validateConfigAck(record: StudioV2WorkspaceDatabaseConfigRecord, request: StudioV2WorkspaceDatabaseConfigRequest, workspaceId?: string) {
  if (!workspaceId?.trim() || !record || !record.id?.trim() || !record.workspace_id?.trim() || !record.identity_revision?.trim() || !record.setup_revision?.trim() ||
    workspaceId && record.workspace_id !== workspaceId || request.connector_id && record.id !== request.connector_id) {
    throw new Error('Configuration save confirmation unavailable');
  }
}
/** The aggregate contains freshness only, never another persisted configuration. */
export function refreshConfigReads(client: QueryClient, context: ConfigReadContext): Promise<void> {
  const current = entry(client, context);
  if (current.pending) return current.pending;
  const generation = current.generation;
  current.status = 'refreshing'; notify(client);
  // Yield one turn: repeated synchronous edits supersede obsolete read demand,
  // while each valid ACK has already settled independently of these reads.
  const pending = new Promise<void>(resolve => setTimeout(resolve, 0)).then(async () => {
    if (current.generation !== generation) return;
    try {
      await client.refetchQueries({ ...filters(context), type: 'active' }, { cancelRefetch: false, throwOnError: true });
      if (current.generation === generation) current.status = 'fresh';
    } catch { if (current.generation === generation) current.status = 'stale-error'; }
  }).finally(() => { if (current.pending === pending) { current.pending = undefined; notify(client); } });
  current.pending = pending;
  return pending;
}
export function confirmConfigSave(client: QueryClient, captured: ReturnType<typeof beginConfigSave>) {
  const { context, generation, home } = captured, current = entry(client, context);
  if (current.generation !== generation || recordingHomeSourceContext({ backend: context.backend }).apiOrigin !== context.apiOrigin) {
    throw new Error('Configuration save owner changed');
  }
  // A GET may have started during the PUT. It is older than this ACK too.
  current.generation++; current.confirmed = true; current.pending = undefined;
  const global = entry(client, { ...context, workspaceId: undefined });
  if (global !== current) global.generation++;
  void client.cancelQueries(filters(context));
  void client.invalidateQueries({ ...filters(context), refetchType: 'none' });
  completeRecordingHomeSourceMutation(client, home);
  void refreshConfigReads(client, context);
}
export function failConfigSave(client: QueryClient, captured: ReturnType<typeof beginConfigSave>) {
  const current = entry(client, captured.context);
  if (current.generation === captured.generation) { current.status = 'stale-error'; notify(client); }
}
function refreshSnapshot(client: QueryClient, context: ConfigReadContext): string {
  const current = registry(client).entries.get(id(context));
  if (!current) return 'unconfirmed:idle';
  const queries = client.getQueryCache().findAll(filters(context)).filter(query => query.isActive());
  // An independent page retry/remount can recover the same fenced reads. Require
  // every read family, so one successful bootstrap never conceals a config error.
  const recovered = current.status === 'stale-error' && !current.pending && ['database-config', 'database-targets', 'bootstrap'].every(kind =>
    queries.some(query => query.queryKey[1] === kind)) && queries.every(query =>
    query.state.status === 'success' && !query.state.isInvalidated && query.state.fetchStatus === 'idle');
  const status = recovered ? 'fresh' : current.status === 'fresh' && current.confirmed && queries.some(query => query.state.status === 'error') ? 'stale-error'
    : current.status === 'fresh' && current.confirmed && queries.some(query => query.state.isInvalidated || query.state.fetchStatus !== 'idle') ? 'refreshing'
      : current.status;
  return `${current.confirmed ? 'confirmed' : 'unconfirmed'}:${status}`;
}
export function useRecordingConfigRefresh(workspaceId?: string) {
  const client = useQueryClient(), context = configReadContext(workspaceId), store = registry(client);
  const subscribe = useCallback((listener: () => void) => {
    store.listeners.add(listener);
    const unsubscribe = client.getQueryCache().subscribe(event => { if (event.type === 'updated' && affected(event.query)) listener(); });
    return () => { store.listeners.delete(listener); unsubscribe(); };
  }, [store, client]);
  // Stable primitive projection avoids a rerender for every unrelated query/cache notification.
  const snapshot = useSyncExternalStore(subscribe, () => refreshSnapshot(client, context), () => refreshSnapshot(client, context));
  const [confirmation, state] = snapshot.split(':');
  const status = state as Entry['status'], confirmed = confirmation === 'confirmed';
  return { status, confirmed, unavailable: status === 'refreshing' || status === 'stale-error',
    retry: () => refreshConfigReads(client, context) };
}
