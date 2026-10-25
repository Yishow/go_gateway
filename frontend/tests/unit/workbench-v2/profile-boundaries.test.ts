import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import {
  defaultRecordingProfile,
  isValidRecordingProfile,
  recordingProfileIntervalSeconds,
} from '@/utils/studioV2RecordingProfile';
import { parseStudioV2WriteGroup, parseWriteGroupListData } from '@/utils/studioV2WriteGroupJson';
import {
  draftIsDirty,
  draftToRequestGroup,
  groupToDraft,
  newDraft,
  validateDraft,
} from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { basicManagedGroupDraft } from '@/features/datalink/workbench-v2/state/basicRecording';
import type {
  WriteGroup,
  WriteGroupCreateRequest,
  WriteGroupDraft,
  WriteGroupRecordingProfile,
} from '@/types/studioV2WriteGroup';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

const legacyGroup: WriteGroup = {
  id: 'legacy-1', workspace_id: 'ws-1', revision: 'rev-1', applied_revision: '', name: 'legacy', status: 'draft',
  members: [{ device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', source_revision: 's', mapping_revision: 'm', target_column: 'temperature', required: true }],
  destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: '/d.db', table_schema: 'main', table_name: 'readings', storage_strategy: 'custom' },
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0, incomplete_policy: 'skip_row' },
  write_policy: {}, migration: {}, created_at: 'a', updated_at: 'b',
};

function profileGroup(overrides: Partial<WriteGroup> = {}): WriteGroup {
  return {
    ...legacyGroup, id: 'profile-1', name: 'main cumulative', status: 'draft',
    members: [{ ...legacyGroup.members[0], target_column: 'v_abc', circuit_id: 'circuit-main', circuit_name: 'Main feeder' }],
    destination: { ...legacyGroup.destination, table_name: 'gw_group_abc', storage_strategy: 'managed' },
    row_policy: { interval_seconds: 600, allowed_lateness_seconds: 0 },
    recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' },
    ...overrides,
  };
}

describe('recording profile boundaries', () => {
  it('uses the contract defaults for new profile drafts only', () => {
    expect(defaultRecordingProfile('instantaneous')).toEqual({ kind: 'instantaneous', interval: 15, unit: 'second' });
    expect(defaultRecordingProfile('cumulative')).toEqual({ kind: 'cumulative', interval: 10, unit: 'minute' });
    // A plain new group is not a profile group and keeps its 60 second default.
    const plain = newDraft({ connector_id: 'c', connector_revision: 'r', table_schema: 'main', table_name: 't' });
    expect(plain.interval_seconds).toBe(60);
    expect(plain).not.toHaveProperty('recording_profile');
  });

  it.each([
    ['instantaneous', 1, 'second'], ['instantaneous', 15, 'second'], ['instantaneous', 300, 'second'],
    ['cumulative', 1, 'minute'], ['cumulative', 10, 'minute'], ['cumulative', 600, 'minute'],
  ])('accepts %s %i %s', (kind, interval, unit) => {
    expect(isValidRecordingProfile({ kind, interval, unit })).toBe(true);
  });

  it.each([
    ['instantaneous', 0, 'second'], ['instantaneous', 301, 'second'], ['instantaneous', -1, 'second'],
    ['instantaneous', 15.5, 'second'], ['instantaneous', Number.NaN, 'second'], ['instantaneous', Infinity, 'second'],
    ['instantaneous', '15', 'second'], ['instantaneous', 15, 'minute'], ['instantaneous', 15, undefined],
    ['cumulative', 0, 'minute'], ['cumulative', 601, 'minute'], ['cumulative', 2.5, 'minute'],
    ['cumulative', 10, 'second'], ['cumulative', 10, 'hour'], ['cumulative', undefined, 'minute'],
    ['average', 10, 'second'],
  ])('rejects %s %s %s', (kind, interval, unit) => {
    expect(isValidRecordingProfile({ kind, interval, unit })).toBe(false);
  });

  it('rejects non-objects', () => {
    for (const value of [null, undefined, 15, 'second', []]) expect(isValidRecordingProfile(value)).toBe(false);
  });

  it('converts to seconds with the unit that belongs to the kind', () => {
    expect(recordingProfileIntervalSeconds({ kind: 'instantaneous', interval: 300, unit: 'second' })).toBe(300);
    expect(recordingProfileIntervalSeconds({ kind: 'cumulative', interval: 10, unit: 'minute' })).toBe(600);
    expect(recordingProfileIntervalSeconds({ kind: 'cumulative', interval: 600, unit: 'minute' })).toBe(36000);
    // The TS side never guesses a unit that the backend would read differently.
    expect(recordingProfileIntervalSeconds({ kind: 'cumulative', interval: 10, unit: 'second' })).toBeNull();
    expect(recordingProfileIntervalSeconds({ kind: 'instantaneous', interval: 15, unit: 'minute' })).toBeNull();
  });
});

describe('recording profile hydration', () => {
  it('keeps the profile, circuit identity and public unit', () => {
    const parsed = parseStudioV2WriteGroup({ ...profileGroup(), layout_version: 'managed_wide_v1' });
    expect(parsed?.recording_profile).toEqual({ kind: 'cumulative', interval: 10, unit: 'minute' });
    expect(parsed?.layout_version).toBe('managed_wide_v1');
    expect(parsed?.members[0]).toMatchObject({ circuit_id: 'circuit-main', circuit_name: 'Main feeder' });
    expect(parsed?.row_policy.interval_seconds).toBe(600);
  });

  it('keeps legacy groups exactly as saved, without inventing a profile', () => {
    const parsed = parseStudioV2WriteGroup(legacyGroup);
    expect(parsed).not.toBeNull();
    expect(parsed).not.toHaveProperty('recording_profile');
    expect(parsed).not.toHaveProperty('layout_version');
    expect(parsed?.members[0]).not.toHaveProperty('circuit_id');
    expect(parsed?.row_policy.interval_seconds).toBe(60);
  });

  it('survives list hydration', () => {
    const listed = parseWriteGroupListData({ workspace_id: 'ws-1', workspace_revision: 'w1', groups: [legacyGroup, profileGroup()] });
    expect(listed?.groups[1].recording_profile?.unit).toBe('minute');
    expect(listed?.groups[0]).not.toHaveProperty('recording_profile');
  });

  it.each([
    ['zero', { kind: 'cumulative', interval: 0, unit: 'minute' }, 0],
    ['above the upper bound', { kind: 'cumulative', interval: 601, unit: 'minute' }, 36060],
    ['fractional', { kind: 'cumulative', interval: 2.5, unit: 'minute' }, 150],
    ['wrong unit', { kind: 'cumulative', interval: 10, unit: 'second' }, 10],
    ['missing unit', { kind: 'cumulative', interval: 10 }, 600],
    ['unknown kind', { kind: 'average', interval: 10, unit: 'minute' }, 600],
  ])('rejects a response whose profile interval is %s', (_name, profile, seconds) => {
    expect(parseStudioV2WriteGroup(profileGroup({
      recording_profile: profile as WriteGroup['recording_profile'],
      row_policy: { interval_seconds: seconds, allowed_lateness_seconds: 0 },
    }))).toBeNull();
  });

  it('rejects a response whose row interval contradicts the profile unit', () => {
    // 10 minutes read as 10 seconds would be the unit confusion this contract forbids.
    expect(parseStudioV2WriteGroup(profileGroup({ row_policy: { interval_seconds: 10, allowed_lateness_seconds: 0 } }))).toBeNull();
  });

  it('rejects an unknown layout and a circuit name without an id', () => {
    expect(parseStudioV2WriteGroup({ ...profileGroup(), layout_version: 'managed_future_v9' })).toBeNull();
    const nameOnly = profileGroup();
    nameOnly.members = [{ ...nameOnly.members[0], circuit_id: undefined }];
    expect(parseStudioV2WriteGroup(nameOnly)).toBeNull();
  });
});

describe('recording profile requests', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.put).mockReset();
  });

  const profileDraft = (): WriteGroupDraft => ({
    workspace_id: 'ws-1', name: 'main cumulative',
    members: [{ device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1', target_column: '', required: true, circuit_id: 'circuit-main', circuit_name: 'Main feeder' }],
    destination: { connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: '', table_name: '', storage_strategy: 'managed' },
    row_policy: { interval_seconds: 600, allowed_lateness_seconds: 0 },
    write_policy: {},
    layout_version: 'managed_wide_v1',
    recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' },
  });

  it('sends the profile and circuit identity on create, update and ensure', async () => {
    const saved = { workspace_revision: 'w2', group: profileGroup() };
    for (const mock of [vi.mocked(studioV2DatalinkApi.post), vi.mocked(studioV2DatalinkApi.put)]) {
      mock.mockResolvedValue({ data: { success: true, data: saved } } as never);
    }
    const request: WriteGroupCreateRequest = {
      workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1', group: profileDraft(),
    };

    await studioV2WorkspaceWriteGroupsAPI.create(request);
    await studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile(request);
    await studioV2WorkspaceWriteGroupsAPI.update('profile-1', { ...request, expected_group_revision: 'rev-1' });

    const [createCall, ensureCall] = vi.mocked(studioV2DatalinkApi.post).mock.calls;
    expect(createCall[0]).toBe('/studio-v2/workspace/write-groups');
    expect(ensureCall[0]).toBe('/studio-v2/workspace/write-groups/circuit-profiles/ensure');
    const updateCall = vi.mocked(studioV2DatalinkApi.put).mock.calls[0];
    for (const body of [createCall[1], ensureCall[1], updateCall[1]] as WriteGroupCreateRequest[]) {
      expect(body.group.recording_profile).toEqual({ kind: 'cumulative', interval: 10, unit: 'minute' });
      expect(body.group.layout_version).toBe('managed_wide_v1');
      expect(body.group.members[0]).toMatchObject({ circuit_id: 'circuit-main', circuit_name: 'Main feeder' });
    }
    // The payload is a copy: later edits to the caller's draft cannot change what was sent.
    request.group.recording_profile!.interval = 999;
    expect((createCall[1] as WriteGroupCreateRequest).group.recording_profile?.interval).toBe(10);
  });

  it('adds nothing to a legacy request', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: { workspace_revision: 'w2', group: legacyGroup } } } as never);
    const legacyDraft = draftToRequestGroup(groupToDraft(legacyGroup), 'ws-1');
    await studioV2WorkspaceWriteGroupsAPI.create({
      workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1', group: legacyDraft,
    });
    const body = vi.mocked(studioV2DatalinkApi.post).mock.calls[0][1] as WriteGroupCreateRequest;
    expect(body.group).not.toHaveProperty('recording_profile');
    expect(body.group).not.toHaveProperty('layout_version');
    expect(body.group.members[0]).not.toHaveProperty('circuit_id');
    expect(body.group.row_policy.interval_seconds).toBe(60);
  });
});

describe('recording profile in the editor draft', () => {
  it('round-trips the profile, layout and circuit identity', () => {
    const saved = parseStudioV2WriteGroup({ ...profileGroup(), layout_version: 'managed_wide_v1' })!;
    const draft = groupToDraft(saved);
    expect(draft.recording_profile).toEqual({ kind: 'cumulative', interval: 10, unit: 'minute' });
    const request = draftToRequestGroup(draft, 'ws-1');
    expect(request.recording_profile).toEqual({ kind: 'cumulative', interval: 10, unit: 'minute' });
    expect(request.layout_version).toBe('managed_wide_v1');
    expect(request.members[0]).toMatchObject({ circuit_id: 'circuit-main', circuit_name: 'Main feeder' });
    expect(request.row_policy.interval_seconds).toBe(600);
    expect(draftIsDirty(draft, saved)).toBe(false);
  });

  it('always derives the row interval from the profile, never from the editable field', () => {
    const draft = groupToDraft(profileGroup());
    draft.interval_seconds = 10; // the stale or mistyped value a seconds field could carry
    expect(draftToRequestGroup(draft, 'ws-1').row_policy.interval_seconds).toBe(600);

    draft.recording_profile = { kind: 'cumulative', interval: 600, unit: 'minute' };
    expect(draftToRequestGroup(draft, 'ws-1').row_policy.interval_seconds).toBe(36000);
    expect(draftIsDirty(draft, profileGroup())).toBe(true);
  });

  it('blocks an invalid profile locally with a profile-specific issue', () => {
    const draft = groupToDraft(profileGroup());
    draft.recording_profile = { kind: 'cumulative', interval: 601, unit: 'minute' };
    expect(validateDraft(draft, { destinationSaved: true }).map((issue) => issue.code)).toContain('profile-interval-invalid');
    draft.recording_profile = { kind: 'cumulative', interval: 10, unit: 'second' };
    expect(validateDraft(draft, { destinationSaved: true }).map((issue) => issue.code)).toContain('profile-interval-invalid');
    draft.recording_profile = { kind: 'cumulative', interval: 10, unit: 'minute' };
    expect(validateDraft(draft, { destinationSaved: true }).map((issue) => issue.code)).not.toContain('profile-interval-invalid');
  });

  it('leaves a legacy 60 second group untouched', () => {
    const draft = groupToDraft(legacyGroup);
    expect(draft).not.toHaveProperty('recording_profile');
    const request = draftToRequestGroup(draft, 'ws-1');
    expect(request).not.toHaveProperty('recording_profile');
    expect(request.row_policy.interval_seconds).toBe(60);
    expect(draftIsDirty(draft, legacyGroup)).toBe(false);
    expect(validateDraft(draft, { destinationSaved: true }).map((issue) => issue.code)).not.toContain('profile-interval-invalid');
  });

  it('keeps the profile when a saved group becomes a basic start-time draft', () => {
    const draft = basicManagedGroupDraft(profileGroup(), 60, 'skip_row');
    expect(draft.recording_profile).toEqual({ kind: 'cumulative', interval: 10, unit: 'minute' });
    expect(draft.members[0]).toMatchObject({ circuit_id: 'circuit-main', circuit_name: 'Main feeder' });
    expect(draft.row_policy.interval_seconds).toBe(600);
    // Without a profile the explicit seconds argument still applies.
    expect(basicManagedGroupDraft(legacyGroup, 45, 'skip_row').row_policy.interval_seconds).toBe(45);
    expect(basicManagedGroupDraft(legacyGroup, 45, 'skip_row')).not.toHaveProperty('recording_profile');
  });
});

describe('an explicit sample reduction is carried, never invented', () => {
  // Raw wire payload with the contract's literal reduction values, never a production helper's output.
  const payload = (kind: 'instantaneous' | 'cumulative', reduction?: string | null) => {
    const profile: Record<string, unknown> = kind === 'instantaneous'
      ? { kind: 'instantaneous', interval: 15, unit: 'second' }
      : { kind: 'cumulative', interval: 10, unit: 'minute' };
    if (reduction !== undefined) profile.sample_reduction = reduction;
    return {
      ...profileGroup(),
      row_policy: { interval_seconds: kind === 'instantaneous' ? 15 : 600, allowed_lateness_seconds: 0 },
      recording_profile: profile,
    };
  };
  const meanGroup = () => parseStudioV2WriteGroup(payload('instantaneous', 'mean'));

  it('hydrates an explicit mean and an explicit latest unchanged', () => {
    expect(parseStudioV2WriteGroup(payload('instantaneous', 'mean'))?.recording_profile?.sample_reduction).toBe('mean');
    expect(parseStudioV2WriteGroup(payload('instantaneous', 'latest'))?.recording_profile?.sample_reduction).toBe('latest');
    expect(parseStudioV2WriteGroup(payload('cumulative', 'latest'))?.recording_profile?.sample_reduction).toBe('latest');
  });

  it('never invents latest for a group saved without a reduction', () => {
    const parsed = parseStudioV2WriteGroup(payload('instantaneous'));
    expect(parsed).not.toBeNull();
    expect(parsed?.recording_profile).not.toHaveProperty('sample_reduction');
    expect(parsed?.recording_profile).toEqual({ kind: 'instantaneous', interval: 15, unit: 'second' });
    expect(parsed?.row_policy.interval_seconds).toBe(15);
  });

  it.each(['median', 'average', 'MEAN', ''])('rejects the unknown reduction %j', (reduction) => {
    expect(parseStudioV2WriteGroup(payload('instantaneous', reduction))).toBeNull();
  });

  it('rejects a null or non-string reduction', () => {
    expect(parseStudioV2WriteGroup(payload('instantaneous', null))).toBeNull();
    expect(parseStudioV2WriteGroup({
      ...payload('instantaneous'), recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 1 },
    })).toBeNull();
  });

  it('rejects mean on a cumulative profile', () => {
    expect(parseStudioV2WriteGroup(payload('cumulative', 'mean'))).toBeNull();
  });

  it('carries mean through the editor draft and the create request without a phantom diff', () => {
    const group = meanGroup()!;
    const draft = groupToDraft(group);
    expect(draft.recording_profile?.sample_reduction).toBe('mean');
    const request = draftToRequestGroup(draft, 'ws-1');
    expect(request.recording_profile?.sample_reduction).toBe('mean');
    expect(request.recording_profile).toMatchObject({ kind: 'instantaneous', interval: 15, unit: 'second' });
    expect(request.row_policy.interval_seconds).toBe(15);
    expect(draftIsDirty(draft, group)).toBe(false);
    expect(validateDraft(draft, { destinationSaved: true }).map((issue) => issue.code))
      .not.toContain('profile-interval-invalid');
    expect(basicManagedGroupDraft(group, 60, 'skip_row').recording_profile?.sample_reduction).toBe('mean');
  });

  it('keeps a legacy group without a reduction through the same draft copy', () => {
    const group = parseStudioV2WriteGroup(payload('instantaneous'))!;
    const draft = groupToDraft(group);
    expect(draft.recording_profile).not.toHaveProperty('sample_reduction');
    expect(draftToRequestGroup(draft, 'ws-1').recording_profile).not.toHaveProperty('sample_reduction');
    expect(draftIsDirty(draft, group)).toBe(false);
  });

  it('rejects mean for a cumulative profile or an unknown reduction in the editor before Apply', () => {
    const draft = groupToDraft(meanGroup()!);
    draft.recording_profile = { kind: 'cumulative', interval: 10, unit: 'minute', sample_reduction: 'mean' };
    expect(validateDraft(draft, { destinationSaved: true }).map((issue) => issue.code)).toContain('profile-interval-invalid');
    draft.recording_profile = { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'median' } as unknown as WriteGroupRecordingProfile;
    expect(validateDraft(draft, { destinationSaved: true }).map((issue) => issue.code)).toContain('profile-interval-invalid');
    draft.recording_profile = { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' };
    expect(validateDraft(draft, { destinationSaved: true }).map((issue) => issue.code)).not.toContain('profile-interval-invalid');
  });

  it('sends and re-reads mean on create, ensure and update', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.put).mockReset();
    const response = { workspace_revision: 'w2', group: payload('instantaneous', 'mean') };
    for (const mock of [vi.mocked(studioV2DatalinkApi.post), vi.mocked(studioV2DatalinkApi.put)]) {
      mock.mockResolvedValue({ data: { success: true, data: response } } as never);
    }
    const request: WriteGroupCreateRequest = {
      workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_connector_revision: 'crev-1',
      group: draftToRequestGroup(groupToDraft(meanGroup()!), 'ws-1'),
    };
    const created = await studioV2WorkspaceWriteGroupsAPI.create(request);
    const ensured = await studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile(request);
    const updated = await studioV2WorkspaceWriteGroupsAPI.update('profile-mean', { ...request, expected_group_revision: 'rev-1' });

    const [createCall, ensureCall] = vi.mocked(studioV2DatalinkApi.post).mock.calls;
    for (const body of [
      (createCall[1] as WriteGroupCreateRequest).group,
      (ensureCall[1] as WriteGroupCreateRequest).group,
      (vi.mocked(studioV2DatalinkApi.put).mock.calls[0][1] as WriteGroupCreateRequest).group,
    ]) {
      expect(body.recording_profile?.sample_reduction).toBe('mean');
    }
    // A reloaded response keeps the explicit choice instead of downgrading it to latest.
    for (const result of [created, ensured, updated]) {
      expect(result.group.recording_profile?.sample_reduction).toBe('mean');
      expect(result.group.row_policy.interval_seconds).toBe(15);
    }
  });
});
