import { QueryClient } from '@tanstack/react-query';
import { vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { parseRecordingStartOperation } from '@/utils/studioV2RecordingStartJson';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { WriteGroup, WriteGroupMutationResponse } from '@/types/studioV2WriteGroup';
import { savedGroup } from './writeGroupState';
import { loadBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';

/** The legacy tests retain their failure/foreign/missing snapshots; this is not a ready fallback. */
export const groupsQuery = {
  data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] } as { workspace_id?: string; workspace_revision?: string; groups?: WriteGroup[] } | undefined,
  isLoading: false, isError: false, refetch: vi.fn(),
};
export const deliveryQuery = { data: undefined, isError: false };
export const ensureAsync = vi.fn(), startAsync = vi.fn(), operationLookup = vi.fn();
const clients = new Set<QueryClient>();
export const retainedLegacyIntent = () => loadBasicRecordingIntent('ws-1/dev-1/managed-recording');
export function syncLegacyGroups() { clients.forEach(client => { if (groupsQuery.data) client.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data); }); }
export function acceptLegacyGroup(result: WriteGroupMutationResponse) {
  groupsQuery.data = { workspace_id: result.group.workspace_id, workspace_revision: result.workspace_revision, groups: [result.group] };
  syncLegacyGroups(); return result;
}
export const operationQuery = { data: undefined as RecordingStartOperation | undefined, isLoading: false, isError: false, refetch: vi.fn() };
export const canonicalManaged = savedGroup({
  basic_managed_device_id: 'dev-1',
  destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_group_server', database: '/recording.db', schema_revision: 'rev-1', schema_digest: 'a'.repeat(64) },
  members: [...savedGroup().members.map(member => ({ ...member, target_column: '' })),
    { device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2', source_revision: 's2', mapping_revision: 'm2', target_column: '', required: true }],
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 }, write_policy: { mode: 'append', dedupe_capability: 'receipt' },
});
export function operation(overrides: Partial<RecordingStartOperation> = {}): RecordingStartOperation {
  const parsed = parseRecordingStartOperation({
    operation_id: 'operation-1', action: 'recording_start', status: 'succeeded', stage: 'complete',
    intent_digest: 'a'.repeat(64), workspace_id: 'ws-1', setup_revision: 'wrev-1', device_ids: ['dev-1'],
    groups: [{ group_id: canonicalManaged.id, group_revision: canonicalManaged.revision, applied_revision: canonicalManaged.revision, saved: true, ready: true, applied: true }],
    devices: [{ device_id: 'dev-1', activated: true }], created_at: '2026-10-04T00:00:00Z', updated_at: '2026-10-04T00:00:01Z', ...overrides,
  });
  if (!parsed) throw new Error('Invalid typed recording-start test fixture');
  return parsed;
}
export function legacyClient() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, structuralSharing: false }, mutations: { retry: false } } });
  clients.add(client); syncLegacyGroups();
  return client;
}
export function resetLegacyOwner() {
  vi.clearAllMocks(); sessionStorage.clear(); clients.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.isLoading = false; groupsQuery.isError = false;
  groupsQuery.refetch.mockImplementation(async () => { syncLegacyGroups(); return { data: groupsQuery.data }; });
  ensureAsync.mockResolvedValue({ workspace_revision: 'wrev-2', group: canonicalManaged });
  startAsync.mockResolvedValue(operation({ status: 'pending', stage: 'save' }));
  operationQuery.data = undefined; operationQuery.isError = false;
  operationQuery.refetch.mockImplementation(async () => ({ data: operationQuery.data }));
}
