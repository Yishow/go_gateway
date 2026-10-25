import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { parseStudioV2WriteGroup } from '@/utils/studioV2WriteGroupJson';
import { draftIsDirty, draftToRequestGroup, groupToDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { basicManagedGroupDraft } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { WriteGroup, WriteGroupCreateRequest } from '@/types/studioV2WriteGroup';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

function member(overrides: Record<string, unknown> = {}): Record<string, unknown> {
  return {
    device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0',
    source_revision: 's0', mapping_revision: 'm0',
    target_column: 'v_0', required: false,
    circuit_id: 'circuit-main', circuit_name: 'Main feeder',
    ...overrides,
  };
}

/** The saved shape the backend returns for a mean wide-profile group. */
function wideProfileGroup(overrides: Record<string, unknown> = {}): Record<string, unknown> {
  return {
    id: 'wide-mean', workspace_id: 'ws-1', revision: 'rev-1', applied_revision: '', name: 'Main feeder (Instantaneous)',
    status: 'draft',
    layout_version: 'managed_wide_profile_v1',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' },
    members: [
      member({ point_id: 'pt-0', tag_id: 'tag-0', value_role: 'gauge', source_value_type: 'int16', result_value_type: 'decimal' }),
      member({ point_id: 'pt-1', tag_id: 'tag-1', value_role: 'state', source_value_type: 'string', result_value_type: 'text' }),
    ],
    destination: {
      connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: 'main',
      table_name: 'gw_Safe_ws_1_wide', storage_strategy: 'managed',
    },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0 },
    write_policy: { mode: 'append', dedupe_capability: 'receipt' },
    migration: {}, created_at: 'a', updated_at: 'b',
    ...overrides,
  };
}

/** A legacy wide payload: known layout, no member roles or value types and no profile key. */
function legacyWideGroup(overrides: Record<string, unknown> = {}): Record<string, unknown> {
  const payload = wideProfileGroup({
    id: 'wide-legacy', layout_version: 'managed_wide_v1', members: [member()], ...overrides,
  });
  delete payload.recording_profile;
  return payload;
}

describe('managed_wide_profile_v1 member metadata is carried, never invented', () => {
  it('hydrates the saved mean wide-profile shape with its roles and value types', () => {
    const parsed = parseStudioV2WriteGroup(wideProfileGroup());
    expect(parsed).not.toBeNull();
    expect(parsed?.layout_version).toBe('managed_wide_profile_v1');
    expect(parsed?.recording_profile).toEqual({ kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' });
    expect(parsed?.row_policy.interval_seconds).toBe(15);
    expect(parsed?.members[0].value_role).toBe('gauge');
    expect(parsed?.members[0].source_value_type).toBe('int16');
    expect(parsed?.members[0].result_value_type).toBe('decimal');
    expect(parsed?.members[1].value_role).toBe('state');
    expect(parsed?.members[1].source_value_type).toBe('string');
    expect(parsed?.members[1].result_value_type).toBe('text');
  });

  it('keeps a legacy member without roles or types free of invented values', () => {
    const parsed = parseStudioV2WriteGroup(legacyWideGroup());
    expect(parsed).not.toBeNull();
    expect(parsed?.members[0]).not.toHaveProperty('value_role');
    expect(parsed?.members[0]).not.toHaveProperty('source_value_type');
    expect(parsed?.members[0]).not.toHaveProperty('result_value_type');
    const draft = groupToDraft(parsed as WriteGroup);
    expect(draft.members[0]).not.toHaveProperty('value_role');
    expect(draftToRequestGroup(draft, 'ws-1').members[0]).not.toHaveProperty('source_value_type');
    expect(draftIsDirty(draft, parsed as WriteGroup)).toBe(false);
  });

  it.each([
    ['value_role', 'average'], ['value_role', 'Gauge'], ['value_role', ''], ['value_role', null], ['value_role', 1],
    ['source_value_type', 'int8'], ['source_value_type', 'number'], ['source_value_type', null], ['source_value_type', true],
    ['result_value_type', 'double'], ['result_value_type', 'varchar'], ['result_value_type', null], ['result_value_type', 0],
  ])('rejects an unknown %s %o on a saved member', (key, value) => {
    const payload = legacyWideGroup({ members: [member({ [key]: value })] });
    expect(parseStudioV2WriteGroup(payload)).toBeNull();
  });

  it.each([
    ['a custom destination', { destination: { ...wideProfileGroup().destination as Record<string, unknown>, storage_strategy: 'custom' } }],
    ['no recording profile', { recording_profile: undefined }],
    ['a wide storage mode', { storage_mode: 'plain' }],
  ])('rejects managed_wide_profile_v1 with %s', (_name, overrides) => {
    expect(parseStudioV2WriteGroup(wideProfileGroup(overrides as Record<string, unknown>))).toBeNull();
  });

  it('keeps an unknown future layout rejected', () => {
    expect(parseStudioV2WriteGroup(wideProfileGroup({ layout_version: 'managed_wide_profile_v9' }))).toBeNull();
  });

  it('carries member roles and types through the editor draft and the create request', () => {
    const saved = wideProfileGroup() as unknown as WriteGroup;
    const draft = groupToDraft(saved);
    expect(draft.members[0].value_role).toBe('gauge');
    expect(draft.members[0].source_value_type).toBe('int16');
    expect(draft.members[0].result_value_type).toBe('decimal');
    const request = draftToRequestGroup(draft, 'ws-1');
    expect(request.members[0].value_role).toBe('gauge');
    expect(request.members[0].source_value_type).toBe('int16');
    expect(request.members[0].result_value_type).toBe('decimal');
    expect(request.members[1].value_role).toBe('state');
    expect(request.members[1].source_value_type).toBe('string');
    expect(request.members[1].result_value_type).toBe('text');
    expect(request.layout_version).toBe('managed_wide_profile_v1');
    // The same saved group must not look dirty because the metadata is dropped.
    expect(draftIsDirty(draft, saved)).toBe(false);
    expect(basicManagedGroupDraft(saved, 15, 'skip_row').members[0].value_role).toBe('gauge');
  });
});

// Unknown runtime data: explicit empty or null metadata is not the same as an omitted legacy field.
function memberWithExplicitBadMetadata(): Record<string, unknown> {
  const [saved] = legacyWideGroup().members as Record<string, unknown>[];
  return { ...saved, value_role: '', source_value_type: null, result_value_type: '' };
}

describe('explicit empty or null metadata is sent instead of dropped as legacy', () => {
  it('keeps it in the basic start-time draft', () => {
    const saved = { ...legacyWideGroup(), members: [memberWithExplicitBadMetadata()] } as unknown as WriteGroup;
    const draft = basicManagedGroupDraft(saved, 15, 'skip_row');
    expect(draft.members[0]).toHaveProperty('value_role', '');
    expect(draft.members[0]).toHaveProperty('source_value_type', null);
    expect(draft.members[0]).toHaveProperty('result_value_type', '');
  });

  it('sends it in the create request for the backend to reject', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({
      data: { success: true, data: { workspace_revision: 'w2', group: legacyWideGroup() } },
    } as never);
    const draft = draftToRequestGroup(groupToDraft(legacyWideGroup() as unknown as WriteGroup), 'ws-1');
    draft.members = [{ ...draft.members[0], value_role: '', source_value_type: null, result_value_type: '' } as unknown as typeof draft.members[number]];

    await studioV2WorkspaceWriteGroupsAPI.create({
      workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1', group: draft,
    });

    const body = vi.mocked(studioV2DatalinkApi.post).mock.calls[0][1] as WriteGroupCreateRequest;
    expect(body.group.members[0]).toHaveProperty('value_role', '');
    expect(body.group.members[0]).toHaveProperty('source_value_type', null);
    expect(body.group.members[0]).toHaveProperty('result_value_type', '');
  });
});

describe('managed_wide_profile_v1 member metadata in workspace requests', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.put).mockReset();
  });

  it('sends and re-reads the metadata on create, ensure and update', async () => {
    const response = { workspace_revision: 'w2', group: wideProfileGroup() };
    for (const mock of [vi.mocked(studioV2DatalinkApi.post), vi.mocked(studioV2DatalinkApi.put)]) {
      mock.mockResolvedValue({ data: { success: true, data: response } } as never);
    }
    const request: WriteGroupCreateRequest = {
      workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1',
      group: draftToRequestGroup(groupToDraft(wideProfileGroup() as unknown as WriteGroup), 'ws-1'),
    };

    const created = await studioV2WorkspaceWriteGroupsAPI.create(request);
    const ensured = await studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile(request);
    const updated = await studioV2WorkspaceWriteGroupsAPI.update('wide-mean', { ...request, expected_group_revision: 'rev-1' });

    const [createCall, ensureCall] = vi.mocked(studioV2DatalinkApi.post).mock.calls;
    for (const body of [
      (createCall[1] as WriteGroupCreateRequest).group,
      (ensureCall[1] as WriteGroupCreateRequest).group,
      (vi.mocked(studioV2DatalinkApi.put).mock.calls[0][1] as WriteGroupCreateRequest).group,
    ]) {
      expect(body.members[0].value_role).toBe('gauge');
      expect(body.members[0].source_value_type).toBe('int16');
      expect(body.members[0].result_value_type).toBe('decimal');
    }
    for (const result of [created, ensured, updated]) {
      expect(result.group.members[0].value_role).toBe('gauge');
      expect(result.group.members[0].source_value_type).toBe('int16');
      expect(result.group.members[0].result_value_type).toBe('decimal');
      expect(result.group.members[1].value_role).toBe('state');
      expect(result.group.layout_version).toBe('managed_wide_profile_v1');
    }
  });

  it('adds no metadata to a legacy request', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({
      data: { success: true, data: { workspace_revision: 'w2', group: legacyWideGroup() } },
    } as never);
    const legacyDraft = draftToRequestGroup(groupToDraft(legacyWideGroup() as unknown as WriteGroup), 'ws-1');
    await studioV2WorkspaceWriteGroupsAPI.create({
      workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1', group: legacyDraft,
    });
    const body = vi.mocked(studioV2DatalinkApi.post).mock.calls[0][1] as WriteGroupCreateRequest;
    expect(body.group.members[0]).not.toHaveProperty('value_role');
    expect(body.group.members[0]).not.toHaveProperty('source_value_type');
    expect(body.group.members[0]).not.toHaveProperty('result_value_type');
  });
});
