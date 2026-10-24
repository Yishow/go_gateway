import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { parseStudioV2WriteGroup } from '@/utils/studioV2WriteGroupJson';
import { parseRecordingSchemaPreviewToken } from '@/utils/recordingPlanJson';
import { parseWriteGroupDeliveryData } from '@/utils/studioV2WriteGroupDeliveryJson';
import { draftToRequestGroup, groupToDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { basicManagedGroupDraft } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { WriteGroup, WriteGroupCreateRequest } from '@/types/studioV2WriteGroup';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

function longGroup(policy: Record<string, unknown> = { mode: 'append', dedupe_capability: 'receipt' }): WriteGroup {
  return {
    id: 'long-1', workspace_id: 'ws-1', revision: 'rev-1', applied_revision: '', name: 'main instantaneous', status: 'draft',
    members: [{
      device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', source_revision: 's', mapping_revision: 'm',
      target_column: 'v_0123456789abcdef01234567', required: true, circuit_id: 'circuit-main', circuit_name: 'Main feeder',
    }],
    destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: 'main', table_name: 'gw_long_0123456789abcdef', storage_strategy: 'managed' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0 },
    write_policy: policy, migration: {},
    layout_version: 'managed_long_v1', storage_mode: 'plain',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    created_at: 'a', updated_at: 'b',
  } as WriteGroup;
}

describe('recreate_missing write policy', () => {
  it('keeps on, off and never-confirmed apart when a group is read', () => {
    expect(parseStudioV2WriteGroup(longGroup({ mode: 'append', recreate_missing: true }))?.write_policy.recreate_missing).toBe(true);
    expect(parseStudioV2WriteGroup(longGroup({ mode: 'append', recreate_missing: false }))?.write_policy.recreate_missing).toBe(false);
    const legacy = parseStudioV2WriteGroup(longGroup({ mode: 'append' }));
    expect(legacy).not.toBeNull();
    expect(legacy!.write_policy).not.toHaveProperty('recreate_missing');
  });

  it('reads an explicit null as never confirmed, not as permission', () => {
    const parsed = parseStudioV2WriteGroup(longGroup({ mode: 'append', recreate_missing: null }));
    expect(parsed).not.toBeNull();
    expect(parsed!.write_policy).not.toHaveProperty('recreate_missing');
  });

  it.each([
    ['a string', 'true'],
    ['a number', 1],
    ['zero', 0],
    ['an object', {}],
    ['an array', []],
  ])('fails closed on %s', (_name, value) => {
    expect(parseStudioV2WriteGroup(longGroup({ mode: 'append', recreate_missing: value }))).toBeNull();
  });

  it.each([true, false])('round-trips %s through the editor draft without a new decision', (value) => {
    const saved = parseStudioV2WriteGroup(longGroup({ mode: 'append', dedupe_capability: 'receipt', recreate_missing: value }))!;
    const request = draftToRequestGroup(groupToDraft(saved), 'ws-1');
    expect(request.write_policy.recreate_missing).toBe(value);
    expect(basicManagedGroupDraft(saved, 60, 'skip_row').write_policy.recreate_missing).toBe(value);
  });

  it('never adds a policy to a group that has none', () => {
    const saved = parseStudioV2WriteGroup(longGroup({ mode: 'append', dedupe_capability: 'receipt' }))!;
    const request = draftToRequestGroup(groupToDraft(saved), 'ws-1');
    expect(request.write_policy).not.toHaveProperty('recreate_missing');
    expect(basicManagedGroupDraft(saved, 60, 'skip_row').write_policy).not.toHaveProperty('recreate_missing');
  });

  describe('requests', () => {
    beforeEach(() => vi.mocked(studioV2DatalinkApi.post).mockReset());

    it.each([[true], [false]])('sends an explicit %s on ensure', async (value) => {
      const group = longGroup({ mode: 'append', dedupe_capability: 'receipt', recreate_missing: value });
      vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: { workspace_revision: 'w2', group } } } as never);
      const request: WriteGroupCreateRequest = {
        workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1',
        group: draftToRequestGroup(groupToDraft(parseStudioV2WriteGroup(group)!), 'ws-1'),
      };
      await studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile(request);
      const body = vi.mocked(studioV2DatalinkApi.post).mock.calls[0][1] as WriteGroupCreateRequest;
      expect(body.group.write_policy.recreate_missing).toBe(value);
    });

    it('sends nothing for a group that never had a policy', async () => {
      const group = longGroup();
      vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: { workspace_revision: 'w2', group } } } as never);
      const request: WriteGroupCreateRequest = {
        workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1',
        group: draftToRequestGroup(groupToDraft(parseStudioV2WriteGroup(group)!), 'ws-1'),
      };
      await studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile(request);
      const body = vi.mocked(studioV2DatalinkApi.post).mock.calls[0][1] as WriteGroupCreateRequest;
      expect(body.group.write_policy).not.toHaveProperty('recreate_missing');
    });
  });
});

describe('schema preview shows the policy it froze', () => {
  const token = (extra: Record<string, unknown> = {}) => ({
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
    },
    ...extra,
  });

  it.each([true, false])('keeps an explicit %s', (value) => {
    expect(parseRecordingSchemaPreviewToken(token({ recreate_missing: value }), { requireGroupLayout: true }))
      .toMatchObject({ recreate_missing: value });
  });

  it('keeps an absent policy absent', () => {
    expect(parseRecordingSchemaPreviewToken(token(), { requireGroupLayout: true })).not.toHaveProperty('recreate_missing');
  });

  it('reads an explicit null as never confirmed', () => {
    expect(parseRecordingSchemaPreviewToken(token({ recreate_missing: null }), { requireGroupLayout: true })).not.toHaveProperty('recreate_missing');
  });

  it.each([['a string', 'true'], ['a number', 1]])('fails closed on %s', (_name, value) => {
    expect(parseRecordingSchemaPreviewToken(token({ recreate_missing: value }), { requireGroupLayout: true })).toBeNull();
  });
});

describe('delivery shows the managed storage truth', () => {
  const base = () => ({
    group_id: 'group-1',
    intake: { state: 'active' },
    stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  });
  const storage = (extra: Record<string, unknown> = {}) => ({
    generation: 2, history_may_be_lost: true, reconstructions: 1, last_reconstructed_at: '2026-10-24T03:10:00Z',
    recreate_missing: true,
    events: [{ generation: 2, group_id: 'group-1', kind: 'reconstructed', reason: 'table-rebuilt', operation_id: 'op-1', occurrences: 1, occurred_at: '2026-10-24T03:10:00Z' }],
    ...extra,
  });

  it('keeps the generation, the lost-history disclosure and the events', () => {
    expect(parseWriteGroupDeliveryData({ ...base(), storage: storage() })?.storage).toEqual(storage());
  });

  it('leaves a group without managed long storage without a storage view', () => {
    expect(parseWriteGroupDeliveryData(base())).not.toHaveProperty('storage');
  });

  it('keeps a blocked repair and an explicit off apart from an absent policy', () => {
    const blocked = parseWriteGroupDeliveryData({ ...base(), storage: storage({ repair_state: 'blocked', repair_reason: 'permission-denied', recreate_missing: false }) });
    expect(blocked?.storage).toMatchObject({ repair_state: 'blocked', repair_reason: 'permission-denied', recreate_missing: false });
    const unset = parseWriteGroupDeliveryData({ ...base(), storage: storage({ recreate_missing: undefined }) });
    expect(unset?.storage).not.toHaveProperty('recreate_missing');
  });

  it.each([
    ['an unknown repair state', { repair_state: 'restoring' }],
    ['a generation below one', { generation: 0 }],
    ['a fractional generation', { generation: 1.5 }],
    ['a policy that is not a boolean', { recreate_missing: 'yes' }],
    ['history loss at the first generation', { generation: 1, history_may_be_lost: true, reconstructions: 0 }],
    ['a rebuild without the history disclosure', { history_may_be_lost: false, reconstructions: 1 }],
    ['an unknown event kind', { events: [{ generation: 1, kind: 'dropped_table', occurrences: 1, occurred_at: '2026-10-24T03:10:00Z' }] }],
    ['events that are not a list', { events: 'none' }],
  ])('fails closed on %s', (_name, extra) => {
    expect(parseWriteGroupDeliveryData({ ...base(), storage: storage(extra) })).toBeNull();
  });
});

describe('delivery shows a cut window as a bucket issue', () => {
  const base = () => ({
    group_id: 'group-1', intake: { state: 'active' },
    stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  });

  it('keeps a row that a cutover cut short apart from a skipped one', () => {
    const issue = { group_revision: 'rev-1', bucket_start: '2026-10-24T04:50:15Z', kind: 'truncated', causes: ['boundary'] };
    expect(parseWriteGroupDeliveryData({ ...base(), recent_bucket_issues: [issue] })?.recent_bucket_issues).toEqual([issue]);
    expect(parseWriteGroupDeliveryData({ ...base(), recent_bucket_issues: [{ ...issue, kind: 'cut' }] })).toBeNull();
  });
});
