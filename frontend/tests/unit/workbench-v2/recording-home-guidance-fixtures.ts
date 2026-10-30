import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';
import type { RecordingHomeGuidanceInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import type { RecordingHomeGroupScope } from '@/features/datalink/workbench-v2/recording-home/recordingHomeProjection';
import type { RecordingOwnerObservations } from '@/features/datalink/workbench-v2/state/recordingOwnerObservation';
export const scope = { workspaceId: 'W1', workspaceRevision: 'W2', epoch: 2 };
export function group(id: string, status: WriteGroup['status'] = 'running'): WriteGroup {
  return { id, workspace_id: 'W1', revision: `${id}2`, applied_revision: `${id}1`, name: `Recording ${id}`, status,
    members: [{ device_id: 'D1', point_id: 'P1', tag_id: 'T1', source_revision: 'S1', mapping_revision: 'M1', target_column: 'temperature', required: true }],
    destination: { connector_id: 'DB1', connector_revision: 'DB2', database: 'telemetry', table_schema: 'public', table_name: 'readings', storage_strategy: 'managed' },
    row_policy: { interval_seconds: 30, allowed_lateness_seconds: 2 }, write_policy: {}, migration: {}, created_at: '', updated_at: '' };
}
export function groupScope(value: WriteGroup): RecordingHomeGroupScope {
  return { ...scope, groupId: value.id, groupRevision: value.revision, appliedRevision: value.applied_revision };
}
export function owners(value: WriteGroup): RecordingOwnerObservations {
  return ['schema', 'start', 'test_write'].map(owner => ({ owner: owner as 'schema' | 'start' | 'test_write',
    scope: { ...groupScope(value), deviceId: 'D1' }, state: 'observed', result: 'none' }));
}
export function delivery(value: WriteGroup): WriteGroupDelivery {
  return { group_id: value.id, intake: { state: 'active' },
    stages: { collecting: 2, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 3, skipped: 0 },
    last_sql_committed_at: '2026-10-30T10:32:00Z', oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false } };
}
export function fixture(devices: string[] = [], points = false, groups: WriteGroup[] = []): RecordingHomeGuidanceInput {
  return { evidence: { scope,
    workspace: { state: 'available', scope, value: { id: 'W1', kind: 'single', status: devices.length ? 'ready' : 'empty',
      ordered_device_ids: devices, created_at: '', updated_at: '' } },
    sources: { state: 'available', scope, value: { candidates: points ? [{ key: 'P1|T1', device_id: 'D1', point_id: 'P1', tag_id: 'T1',
      tag_key: 'temperature', label: 'Temperature', device_name: 'Device 1', address: '0', target_type: 'float32' }] : [], excluded: [] } },
    groups: { state: 'available', scope, value: { workspace_id: 'W1', workspace_revision: 'W2', groups } },
    lifecycle: Object.fromEntries(groups.map(value => [value.id, { state: 'available', scope: groupScope(value), value: {
      workspace_id: 'W1', workspace_revision: 'W2', group_id: value.id, group_revision: value.revision, applied_revision: value.applied_revision,
      status: value.status, connector_id: 'DB1', connector_revision: 'DB2', actions: [],
      history: { canonical: true, current_group_id: value.id, generations: [] },
    } }])),
    delivery: Object.fromEntries(groups.map(value => [value.id, { state: 'available', scope: groupScope(value), value: delivery(value) }])),
  }, owners: Object.fromEntries(groups.map(value => [value.id, { state: 'available', scope: groupScope(value), value: owners(value) }])) };
}
