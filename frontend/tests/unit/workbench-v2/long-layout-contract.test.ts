import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { parseStudioV2WriteGroup } from '@/utils/studioV2WriteGroupJson';
import { parseRecordingSchemaPreviewToken } from '@/utils/recordingPlanJson';
import { draftToRequestGroup, groupToDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { basicManagedGroupDraft } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { WriteGroup, WriteGroupCreateRequest } from '@/types/studioV2WriteGroup';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

const longGroup: WriteGroup = {
  id: 'long-1', workspace_id: 'ws-1', revision: 'rev-1', applied_revision: '', name: 'main instantaneous', status: 'draft',
  members: [{
    device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', source_revision: 's', mapping_revision: 'm',
    target_column: 'v_0123456789abcdef01234567', required: true, circuit_id: 'circuit-main', circuit_name: 'Main feeder',
  }],
  destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: 'main', table_name: 'gw_long_0123456789abcdef', storage_strategy: 'managed' },
  row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0 },
  write_policy: { mode: 'append', dedupe_capability: 'receipt' }, migration: {},
  layout_version: 'managed_long_v1', storage_mode: 'plain',
  recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
  created_at: 'a', updated_at: 'b',
};

describe('long layout hydration', () => {
  it('keeps the layout, storage mode, profile and shared table', () => {
    const parsed = parseStudioV2WriteGroup(longGroup);
    expect(parsed).toMatchObject({
      layout_version: 'managed_long_v1', storage_mode: 'plain',
      destination: { table_name: 'gw_long_0123456789abcdef' },
      recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    });
  });

  it.each([
    ['a long group without a storage mode', { ...longGroup, storage_mode: undefined }],
    ['a long group without a profile', { ...longGroup, recording_profile: undefined }],
    ['a storage mode on a wide group', { ...longGroup, layout_version: 'managed_wide_v1' }],
    ['a storage mode on a legacy group', { ...longGroup, layout_version: undefined, recording_profile: undefined, row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 } }],
    ['an unknown storage mode', { ...longGroup, storage_mode: 'native' }],
    ['an unknown layout', { ...longGroup, layout_version: 'managed_future_v9' }],
  ])('rejects %s', (_name, group) => {
    expect(parseStudioV2WriteGroup(group)).toBeNull();
  });

  it('reads an explicitly selected timescale group and keeps its mode distinct from plain', () => {
    const parsed = parseStudioV2WriteGroup({ ...longGroup, storage_mode: 'timescale' });
    expect(parsed).toMatchObject({ layout_version: 'managed_long_v1', storage_mode: 'timescale' });
    expect(parseStudioV2WriteGroup(longGroup)).toMatchObject({ storage_mode: 'plain' });
  });

  it('keeps a wide profile group free of a storage mode', () => {
    const wide = { ...longGroup, layout_version: 'managed_wide_v1' as const, storage_mode: undefined };
    expect(parseStudioV2WriteGroup(wide)).not.toHaveProperty('storage_mode');
  });
});

describe('long layout requests and drafts', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
  });

  it('round-trips the layout and mode through the editor draft', () => {
    const saved = parseStudioV2WriteGroup(longGroup)!;
    const draft = groupToDraft(saved);
    expect(draft).toMatchObject({ layout_version: 'managed_long_v1', storage_mode: 'plain' });
    const request = draftToRequestGroup(draft, 'ws-1');
    expect(request).toMatchObject({ layout_version: 'managed_long_v1', storage_mode: 'plain' });
    expect(request.row_policy.interval_seconds).toBe(15);
    expect(basicManagedGroupDraft(saved, 60, 'skip_row')).toMatchObject({ layout_version: 'managed_long_v1', storage_mode: 'plain' });
  });

  it('sends the layout and mode on ensure', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: { workspace_revision: 'w2', group: longGroup } } } as never);
    const request: WriteGroupCreateRequest = {
      workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1',
      group: draftToRequestGroup(groupToDraft(parseStudioV2WriteGroup(longGroup)!), 'ws-1'),
    };
    await studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile(request);
    const body = vi.mocked(studioV2DatalinkApi.post).mock.calls[0][1] as WriteGroupCreateRequest;
    expect(body.group).toMatchObject({ layout_version: 'managed_long_v1', storage_mode: 'plain' });
  });

  it('sends neither for a legacy group', () => {
    const legacy = parseStudioV2WriteGroup({
      ...longGroup, layout_version: undefined, storage_mode: undefined, recording_profile: undefined,
      row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0 },
    })!;
    const request = draftToRequestGroup(groupToDraft(legacy), 'ws-1');
    expect(request).not.toHaveProperty('storage_mode');
    expect(request).not.toHaveProperty('layout_version');
  });
});

describe('long layout schema preview', () => {
  const token = {
    token: 'tok-1', operation_id: 'op-1', action: 'schema_apply', workspace_id: 'ws-1', workspace_revision: 'setup-9',
    plan_id: 'long-1', plan_revision: 'rev-2', connector_id: 'db-1', connector_revision: 'identity-7', dialect: 'sqlite',
    database: '/data/line-a.db', schema: 'main', table_prefix: '', statements: [],
    tables: [{ name: 'gw_long_0123456789abcdef', action: 'unchanged' }], no_change_reason: 'schema_already_compatible',
    digest: 'a'.repeat(64), expires_at: '2026-09-16T01:00:00Z', created_at: '2026-09-16T00:50:00Z',
    group_layout: {
      table_name: 'gw_long_0123456789abcdef', owner_column: '_gw_owner_0123456789abcdef01234567',
      columns: [
        { name: 'bucket_start', sql_type: 'TEXT', nullable: false, primary_key: true },
        { name: 'record_id', sql_type: 'TEXT', nullable: false, primary_key: true },
        { name: '_gw_owner_0123456789abcdef01234567', sql_type: 'TEXT', nullable: false, primary_key: false },
      ],
      indexes: [{ name: 'idx_gw_long_0123456789abcdef_src_time', columns: ['workspace_id', 'tag_id', 'bucket_start'] }],
    },
  };

  it('keeps the composite key columns and the declared index', () => {
    const parsed = parseRecordingSchemaPreviewToken(token, { requireGroupLayout: true });
    expect(parsed?.group_layout?.columns.filter((column) => column.primary_key).map((column) => column.name)).toEqual(['bucket_start', 'record_id']);
    expect(parsed?.group_layout?.indexes).toEqual([{ name: 'idx_gw_long_0123456789abcdef_src_time', columns: ['workspace_id', 'tag_id', 'bucket_start'] }]);
    expect(parsed?.statements).toEqual([]);
  });

  it('leaves a wide layout without indexes untouched', () => {
    const wide = { ...token, group_layout: { ...token.group_layout, indexes: undefined } };
    expect(parseRecordingSchemaPreviewToken(wide, { requireGroupLayout: true })?.group_layout).not.toHaveProperty('indexes');
  });

  it.each([
    ['an index without columns', [{ name: 'idx_a', columns: [] }]],
    ['an unnamed index', [{ columns: ['tag_id'] }]],
    ['a non-list', { name: 'idx_a' }],
  ])('rejects %s', (_name, indexes) => {
    expect(parseRecordingSchemaPreviewToken({ ...token, group_layout: { ...token.group_layout, indexes } }, { requireGroupLayout: true })).toBeNull();
  });
});
