import type { AxiosInstance } from 'axios';
import type { QueryClient } from '@tanstack/react-query';
import { studioV2DatalinkApi } from '../../services/studioV2Workspace';
import type { RecordingHomeSourceAvailable } from '../../types/recordingHomeSources';
import { studioV2WorkspaceKeys } from './keys';
export interface HomeSourceContext {
  backend?: AxiosInstance;
  /** Resolved API base URL, INCLUDING base pathname, not merely URL.origin. */
  apiOrigin?: string;
}
interface Entry { observation?: RecordingHomeSourceAvailable; generation: number; epoch: number }
interface Partition { entries: Map<string, Entry>; listeners: Set<() => void>; version: number }
const contexts = new WeakMap<QueryClient, WeakMap<AxiosInstance, Map<string, Partition>>>();
const tokens = new WeakMap<AxiosInstance, number>();
let nextToken = 0;
export function recordingHomeSourceContext(options: HomeSourceContext = {}) {
  const backend = options.backend ?? studioV2DatalinkApi;
  const apiOrigin = options.apiOrigin ?? new URL(backend.defaults.baseURL ?? '', globalThis.location?.href ?? 'http://localhost/').href;
  if (!tokens.has(backend)) tokens.set(backend, ++nextToken);
  return { backend, apiOrigin, backendToken: tokens.get(backend)! };
}
export function recordingHomeSourcePrefix(options: HomeSourceContext = {}, workspaceId?: string) {
  const context = recordingHomeSourceContext(options);
  const prefix = [...studioV2WorkspaceKeys.all, 'recording-home-sources', context.backendToken, context.apiOrigin] as const;
  return workspaceId === undefined ? prefix : [...prefix, workspaceId] as const;
}
export function sourcePartition(client: QueryClient, options: HomeSourceContext = {}): Partition {
  const { backend, apiOrigin } = recordingHomeSourceContext(options);
  let clients = contexts.get(client);
  if (!clients) { clients = new WeakMap(); contexts.set(client, clients); }
  let origins = clients.get(backend);
  if (!origins) { origins = new Map(); clients.set(backend, origins); }
  let partition = origins.get(apiOrigin);
  if (!partition) { partition = { entries: new Map(), listeners: new Set(), version: 0 }; origins.set(apiOrigin, partition); }
  return partition;
}
export function sourceEntry(partition: Partition, workspaceId: string): Entry {
  let entry = partition.entries.get(workspaceId);
  if (!entry) {
    entry = { generation: 0, epoch: -1 }; partition.entries.set(workspaceId, entry);
    while (partition.entries.size > 8) partition.entries.delete(partition.entries.keys().next().value!);
  }
  return entry;
}
export function notifySources(partition: Partition) {
  partition.version++; partition.listeners.forEach(listener => listener());
}
/** Marks presentation stale synchronously. Adds no awaited read to a save barrier. */
export function invalidateRecordingHomeSources(client: QueryClient, options: HomeSourceContext & { workspaceId?: string } = {}): void {
  const workspaceId = options.workspaceId ?? client.getQueryData<{ id: string }>(studioV2WorkspaceKeys.bootstrap())?.id;
  if (!workspaceId) return;
  const partition = sourcePartition(client, options), entry = sourceEntry(partition, workspaceId);
  entry.generation++;
  notifySources(partition);
  const queryKey = recordingHomeSourcePrefix(options, workspaceId);
  void client.cancelQueries({ queryKey });
  void client.invalidateQueries({ queryKey, refetchType: 'none' });
}
