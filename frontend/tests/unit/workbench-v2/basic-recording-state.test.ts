import { describe, expect, it } from 'vitest';
import type { WriteGroup, WriteGroupDestinationDraft } from '../../../src/types/studioV2WriteGroup';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import {
  basicRecordingScopeKey,
  buildBasicManagedDrafts,
  createBasicRecordingRequestId,
  isBasicManagedWriteGroup,
  loadBasicRecordingIntent,
  saveBasicRecordingIntent,
} from '../../../src/features/datalink/workbench-v2/state/basicRecording';

const destination: WriteGroupDestinationDraft = {
  connector_id: 'connector-1',
  connector_revision: 'connector-rev-1',
  table_schema: 'main',
  table_name: '',
  storage_strategy: 'managed',
};

const candidates: GroupCandidate[] = [
  {
    key: 'pt-a|tag-a', device_id: 'device-a', point_id: 'pt-a', tag_id: 'tag-a',
    tag_key: 'line-a.temperature', label: 'Temperature A', device_name: 'Same name',
    address: '40001', target_type: 'float64',
  },
  {
    key: 'pt-b|tag-b', device_id: 'device-b', point_id: 'pt-b', tag_id: 'tag-b',
    tag_key: 'line-b.pressure', label: 'Pressure B', device_name: 'Same name',
    address: '40001', target_type: 'int16',
  },
  {
    key: 'pt-a-2|tag-a-2', device_id: 'device-a', point_id: 'pt-a-2', tag_id: 'tag-a-2',
    tag_key: 'line-a.pressure', label: 'Pressure A', device_name: 'Same name',
    address: '40002', target_type: 'int16',
  },
];

function group(overrides: Partial<WriteGroup> = {}): WriteGroup {
  return {
    id: 'group-1', workspace_id: 'workspace-1', revision: 'rev-1', applied_revision: '',
    name: 'Existing group', status: 'ready', members: [],
    destination: {
      connector_id: 'connector-1', connector_revision: 'connector-rev-1', database: '',
      table_schema: 'main', table_name: 'readings', storage_strategy: 'managed',
    },
    row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 },
    write_policy: { mode: 'append', dedupe_capability: 'receipt' },
    migration: {}, created_at: '', updated_at: '', ...overrides,
  };
}

describe('basic recording state', () => {
  it('persists only the scoped request identity and reloads it from offline storage', () => {
    window.sessionStorage.clear();
    const request = {
      request_id: createBasicRecordingRequestId(), workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1',
      device_ids: ['device-a'], groups: [{ group_id: 'group-a', expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1' }],
    };
    saveBasicRecordingIntent({ scope_key: 'workspace-1/device-a/managed-recording', request, operation_id: 'operation-a' });
    const loaded = loadBasicRecordingIntent('workspace-1/device-a/managed-recording');
    expect(loaded).toEqual({ version: 1, scope_key: 'workspace-1/device-a/managed-recording', request, operation_id: 'operation-a' });
    expect(window.sessionStorage.getItem('wbv2.basic-recording.workspace-1%2Fdevice-a%2Fmanaged-recording')).not.toContain('password');
    expect(request.request_id).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/);
  });

  it('rejects malformed or cross-device persisted intents before the panel reads them', () => {
    const key = 'wbv2.basic-recording.workspace-1%2Fdevice-a%2Fmanaged-recording';
    const base = {
      version: 1,
      scope_key: 'workspace-1/device-a/managed-recording',
      request: {
        request_id: '11111111-1111-4111-8111-111111111111', workspace_id: 'workspace-1',
        expected_workspace_revision: 'workspace-rev-1', device_ids: ['device-a'], groups: [],
      },
    };
    window.sessionStorage.setItem(key, JSON.stringify({ ...base, request: { ...base.request, groups: null } }));
    expect(loadBasicRecordingIntent('workspace-1/device-a/managed-recording')).toBeUndefined();
    window.sessionStorage.setItem(key, JSON.stringify({ ...base, request: { ...base.request, device_ids: ['device-b'] } }));
    expect(loadBasicRecordingIntent('workspace-1/device-a/managed-recording')).toBeUndefined();
    window.sessionStorage.setItem(key, JSON.stringify({ ...base, version: 2 }));
    expect(loadBasicRecordingIntent('workspace-1/device-a/managed-recording')).toBeUndefined();
  });

  it('restores a Share-only request with no database revision using the original identity', () => {
    const scopeKey = 'workspace-1/device-a/managed-recording';
    const request = {
      request_id: '11111111-1111-4111-8111-111111111111', workspace_id: 'workspace-1',
      expected_workspace_revision: '', device_ids: ['device-a'], groups: [],
      readiness_token: 'ready-1', settings_revision: 'settings-1', workspace_revision: 'share-rev-1',
    };
    saveBasicRecordingIntent({ scope_key: scopeKey, request, operation_id: 'operation-share-a' });

    expect(loadBasicRecordingIntent(scopeKey)).toEqual({
      version: 1, scope_key: scopeKey, request, operation_id: 'operation-share-a',
    });

    saveBasicRecordingIntent({ scope_key: scopeKey, request: { ...request, groups: [{
      group_id: 'group-a', expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1',
    }] }, operation_id: 'operation-db-a' });
    expect(loadBasicRecordingIntent(scopeKey)).toBeUndefined();
  });

  it.each([undefined, null, 12, ' '])('rejects an invalid Share-only database revision %s', (revision) => {
    const scopeKey = 'workspace-1/device-a/managed-recording';
    window.sessionStorage.setItem('wbv2.basic-recording.workspace-1%2Fdevice-a%2Fmanaged-recording', JSON.stringify({
      version: 1, scope_key: scopeKey, request: {
        request_id: '11111111-1111-4111-8111-111111111111', workspace_id: 'workspace-1',
        expected_workspace_revision: revision, device_ids: ['device-a'], groups: [],
      },
    }));
    expect(loadBasicRecordingIntent(scopeKey)).toBeUndefined();
  });

  it('keeps the create-once identity independent from display names and destinations', () => {
    expect(basicRecordingScopeKey('workspace-1', 'device-a', 'managed-recording'))
      .toBe('workspace-1/device-a/managed-recording');
    expect(basicRecordingScopeKey('workspace-1', 'device-a', 'managed-recording'))
      .toBe(basicRecordingScopeKey('workspace-1', 'device-a', 'managed-recording'));
    expect(basicRecordingScopeKey('workspace-1', 'device-a', 'managed-recording'))
      .not.toBe(basicRecordingScopeKey('workspace-1', 'device-b', 'managed-recording'));
  });

  it('builds one server-owned managed draft per device from real point/tag identities', () => {
    const drafts = buildBasicManagedDrafts('workspace-1', candidates, destination, 'managed-recording');

    expect(drafts.map((draft) => draft.scope_key)).toEqual([
      'workspace-1/device-a/managed-recording',
      'workspace-1/device-b/managed-recording',
    ]);
    expect(drafts[0].draft.members).toMatchObject([
      { device_id: 'device-a', point_id: 'pt-a', tag_id: 'tag-a', target_column: '' },
      { device_id: 'device-a', point_id: 'pt-a-2', tag_id: 'tag-a-2', target_column: '' },
    ]);
    expect(drafts[0].draft.destination).toMatchObject({ storage_strategy: 'managed', table_name: '' });
    expect(drafts[0].draft.row_policy).toMatchObject({ interval_seconds: 60, incomplete_policy: 'skip_row' });
    expect(drafts[0].draft.members.every((member) => member.required)).toBe(true);
    expect(JSON.stringify(drafts)).not.toContain('measurement_id');
  });

  it('keeps custom, multi-entity and cross-device groups in advanced mode', () => {
    const custom = group({ destination: { ...group().destination, storage_strategy: 'custom' } });
    const multiEntity = group({
      members: [{
        device_id: 'device-a', point_id: 'pt-a', tag_id: 'tag-a', entity_key: 'batch-a',
        source_revision: 'source-1', mapping_revision: 'mapping-1', target_column: 'temperature', required: true,
      }],
      row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0, entity_key_column: 'entity' },
    });
    const crossDevice = group({
      members: [
        {
          device_id: 'device-a', point_id: 'pt-a', tag_id: 'tag-a', source_revision: 'source-1',
          mapping_revision: 'mapping-1', target_column: 'temperature', required: true,
        },
        {
          device_id: 'device-b', point_id: 'pt-b', tag_id: 'tag-b', source_revision: 'source-1',
          mapping_revision: 'mapping-1', target_column: 'pressure', required: true,
        },
      ],
    });

    expect(isBasicManagedWriteGroup(custom, 'device-a')).toBe(false);
    expect(isBasicManagedWriteGroup(multiEntity, 'device-a')).toBe(false);
    expect(isBasicManagedWriteGroup(crossDevice, 'device-a')).toBe(false);
    expect(isBasicManagedWriteGroup(group({
      basic_managed_device_id: 'device-a',
      members: [{ ...multiEntity.members[0], entity_key: undefined }],
    }), 'device-a')).toBe(true);
    expect(isBasicManagedWriteGroup(group({
      members: [{ ...multiEntity.members[0], entity_key: undefined }],
    }), 'device-a')).toBe(false);
  });
});
