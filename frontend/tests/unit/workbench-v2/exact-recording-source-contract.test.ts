import { afterEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { basicManagedGroupDraft } from '@/features/datalink/workbench-v2/state/basicRecording';
import { draftIsDirty, draftToRequestGroup, groupToDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import apiFixture from '../../fixtures/writeGroupExactDecimalAPIResponse.json';
import type { WriteGroupProfileKind, WriteGroupUpdateRequest } from '@/types/studioV2WriteGroup';
import { parseStudioV2WriteGroup } from '@/utils/studioV2WriteGroupJson';

const exactSource = {
  measurement_id: 'circuit-energy',
  definition_revision: 'rev-1',
  source_binding_revision: 'rev-1',
  series_epoch: 'epoch-1',
  original_data_type: 'uint64',
  scale: '0.1',
  unit: 'kWh',
  semantic_kind: 'counter',
  counter_policy: { allow_negative: false },
};

const exactSourceWithoutPolicy = {
  measurement_id: 'circuit-energy',
  definition_revision: 'rev-1',
  source_binding_revision: 'rev-1',
  series_epoch: 'epoch-1',
  original_data_type: 'uint64',
  scale: '0.1',
  unit: 'kWh',
  semantic_kind: 'counter',
};

type APIListFixture = {
  success: boolean;
  data: { workspace_id: string; workspace_revision: string; groups: Record<string, unknown>[] };
};

const preparedListFixture = apiFixture as unknown as APIListFixture;
const preparedGroupFixture = preparedListFixture.data.groups[0];
const preparedMemberFixture = (preparedGroupFixture.members as Record<string, unknown>[])[0];

function longGroup(
  kind: WriteGroupProfileKind = 'instantaneous',
  memberOverrides: Record<string, unknown> = {},
  groupOverrides: Record<string, unknown> = {},
): Record<string, unknown> {
  const intervalSeconds = kind === 'instantaneous' ? 1 : 60;
  const recordingProfile = kind === 'instantaneous'
    ? { kind, interval: 1, unit: 'second', sample_reduction: 'mean' }
    : { kind, interval: 1, unit: 'minute' };
  const member = { ...preparedMemberFixture, ...memberOverrides };
  if (kind === 'cumulative' && !('value_role' in memberOverrides)) delete member.value_role;
  return {
    ...preparedGroupFixture,
    row_policy: { ...(preparedGroupFixture.row_policy as Record<string, unknown>), interval_seconds: intervalSeconds },
    recording_profile: recordingProfile,
    members: [member],
    ...groupOverrides,
  };
}

afterEach(() => vi.restoreAllMocks());

describe('exact decimal recording source in saved write groups', () => {
  it.each([
    ['instantaneous', 'counter_raw'],
    ['cumulative', undefined],
  ] as const)('accepts the full backend list and save envelopes for %s counter members', async (kind, expectedRole) => {
    const group = longGroup(kind);
    const get = vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce({
      data: { success: true, data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [group] } },
    } as never);
    const listed = await studioV2WorkspaceWriteGroupsAPI.list();
    expect(listed.groups[0].members[0]).toHaveProperty('exact_recording_source', exactSource);
    expect(listed.groups[0].members[0]).toMatchObject({
      measurement_id: 'circuit-energy', source_value_type: 'uint64', result_value_type: 'decimal',
    });
    expect(listed.groups[0].members[0].value_role).toBe(expectedRole);
    const saved = listed.groups[0];
    const editorDraft = groupToDraft(saved);
    expect(editorDraft.members[0]).toHaveProperty('exact_recording_source', exactSource);
    expect(draftIsDirty(editorDraft, saved)).toBe(false);
    const updateGroup = draftToRequestGroup(editorDraft, 'ws-1');
    expect(updateGroup.members[0]).toHaveProperty('exact_recording_source', exactSource);
    expect(basicManagedGroupDraft(saved, 1, 'partial').members[0]).toHaveProperty('exact_recording_source', exactSource);
    expect(get).toHaveBeenCalledWith('/studio-v2/workspace/write-groups');

    vi.spyOn(studioV2DatalinkApi, 'post').mockResolvedValueOnce({
      data: { success: true, data: { workspace_revision: 'wrev-2', group } },
    } as never);
    const ensured = await studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile({
      workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_connector_revision: 'mysql-rev-1',
      group: {
        workspace_id: 'ws-1', name: 'Circuit Instantaneous', members: [],
        destination: { connector_id: 'mysql-1', connector_revision: 'mysql-rev-1', table_schema: 'gw_circuit_test', table_name: '', storage_strategy: 'managed' },
        row_policy: { interval_seconds: kind === 'instantaneous' ? 1 : 60, allowed_lateness_seconds: 0 },
        write_policy: { mode: 'append', dedupe_capability: 'receipt' }, layout_version: 'managed_long_v1',
        recording_profile: { kind, interval: 1, unit: kind === 'instantaneous' ? 'second' : 'minute' }, storage_mode: 'plain',
      },
    });
    expect(ensured.group.members[0]).toHaveProperty('exact_recording_source', exactSource);

    vi.spyOn(studioV2DatalinkApi, 'put').mockResolvedValueOnce({
      data: { success: true, data: { workspace_revision: 'wrev-3', group } },
    } as never);
    const updateRequest: WriteGroupUpdateRequest = {
      workspace_id: 'ws-1', expected_workspace_revision: 'wrev-2',
      expected_group_revision: saved.revision, expected_connector_revision: saved.destination.connector_revision,
      group: updateGroup,
    };
    const updated = await studioV2WorkspaceWriteGroupsAPI.update(saved.id, updateRequest);
    const sent = vi.mocked(studioV2DatalinkApi.put).mock.calls[0][1] as WriteGroupUpdateRequest;
    expect(sent.group.members[0]).toHaveProperty('exact_recording_source', exactSource);
    expect(updated.group.members[0]).toHaveProperty('exact_recording_source', exactSource);
  });

  it('round trips an omitted counter policy without inventing one', async () => {
    const group = longGroup('instantaneous', { exact_recording_source: exactSourceWithoutPolicy });
    vi.spyOn(studioV2DatalinkApi, 'get').mockResolvedValueOnce({
      data: { success: true, data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [group] } },
    } as never);
    const saved = (await studioV2WorkspaceWriteGroupsAPI.list()).groups[0];
    expect(saved.members[0].exact_recording_source).toEqual(exactSourceWithoutPolicy);
    expect(saved.members[0].exact_recording_source).not.toHaveProperty('counter_policy');

    const draft = groupToDraft(saved);
    expect(draft.members[0].exact_recording_source).toEqual(exactSourceWithoutPolicy);
    expect(draftToRequestGroup(draft, 'ws-1').members[0].exact_recording_source).toEqual(exactSourceWithoutPolicy);
    expect(basicManagedGroupDraft(saved, 1, 'partial').members[0].exact_recording_source).toEqual(exactSourceWithoutPolicy);

    vi.spyOn(studioV2DatalinkApi, 'post').mockResolvedValueOnce({
      data: { success: true, data: { workspace_revision: 'wrev-2', group } },
    } as never);
    const ensured = await studioV2WorkspaceWriteGroupsAPI.ensureCircuitProfile({
      workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_connector_revision: 'mysql-rev-1',
      group: {
        workspace_id: 'ws-1', name: 'Circuit Instantaneous', members: [],
        destination: { connector_id: 'mysql-1', connector_revision: 'mysql-rev-1', table_schema: 'gw_circuit_test', table_name: '', storage_strategy: 'managed' },
        row_policy: { interval_seconds: 1, allowed_lateness_seconds: 0 },
        write_policy: { mode: 'append', dedupe_capability: 'receipt' }, layout_version: 'managed_long_v1',
        recording_profile: { kind: 'instantaneous', interval: 1, unit: 'second' }, storage_mode: 'plain',
      },
    });
    expect(ensured.group.members[0].exact_recording_source).toEqual(exactSourceWithoutPolicy);

    vi.spyOn(studioV2DatalinkApi, 'put').mockResolvedValueOnce({
      data: { success: true, data: { workspace_revision: 'wrev-3', group } },
    } as never);
    const updated = await studioV2WorkspaceWriteGroupsAPI.update(saved.id, {
      workspace_id: 'ws-1', expected_workspace_revision: 'wrev-2',
      expected_group_revision: saved.revision, expected_connector_revision: saved.destination.connector_revision,
      group: draftToRequestGroup(draft, 'ws-1'),
    });
    const sent = vi.mocked(studioV2DatalinkApi.put).mock.calls[0][1] as WriteGroupUpdateRequest;
    expect(sent.group.members[0].exact_recording_source).toEqual(exactSourceWithoutPolicy);
    expect(updated.group.members[0].exact_recording_source).toEqual(exactSourceWithoutPolicy);
  });

  it('requires an explicit negative-counter policy only for signed counters', () => {
    const signedProof = { ...exactSource, semantic_kind: 'signed_counter', original_data_type: 'int64', counter_policy: { allow_negative: true } };
    const signedProofWithoutPolicy: Record<string, unknown> = { ...signedProof };
    delete signedProofWithoutPolicy.counter_policy;
    expect(parseStudioV2WriteGroup(longGroup('instantaneous', {
      source_value_type: 'int64', exact_recording_source: signedProof,
    }))).not.toBeNull();
    expect(parseStudioV2WriteGroup(longGroup('instantaneous', {
      source_value_type: 'int64', exact_recording_source: signedProofWithoutPolicy,
    }))).toBeNull();
    expect(parseStudioV2WriteGroup(longGroup('instantaneous', {
      source_value_type: 'int64', exact_recording_source: { ...signedProof, counter_policy: { allow_negative: false } },
    }))).toBeNull();
  });

  it.each(['instantaneous', 'cumulative'] as const)('rejects %s integer counter decimal results without backend proof', (kind) => {
    const payload = longGroup(kind, { exact_recording_source: undefined });
    expect(parseStudioV2WriteGroup(payload)).toBeNull();
  });

  it.each([
    ['zero scale', { scale: '0' }],
    ['exponent scale', { scale: '1e-1' }],
    ['missing revision', { definition_revision: '' }],
    ['wrong source type', { original_data_type: 'float64' }],
    ['wrong measurement', { measurement_id: 'other-measurement' }],
    ['invalid counter policy', { counter_policy: { allow_negative: 'false' } }],
  ])('rejects a malformed or mismatched proof (%s)', (_name, proofPatch) => {
    const payload = longGroup('instantaneous', { exact_recording_source: { ...exactSource, ...proofPatch } });
    expect(parseStudioV2WriteGroup(payload)).toBeNull();
  });

  it('does not apply exact counter proof to gauge members or cumulative counter_raw roles', () => {
    const gauge = longGroup('instantaneous', { value_role: 'gauge' });
    expect(parseStudioV2WriteGroup(gauge)).toBeNull();
    const cumulativeCounterRaw = longGroup('cumulative', { value_role: 'counter_raw' });
    expect(parseStudioV2WriteGroup(cumulativeCounterRaw)).toBeNull();
  });

  it('keeps unscaled integer, float, and custom decimal members compatible without proof', () => {
    expect(parseStudioV2WriteGroup(longGroup('instantaneous', {
      result_value_type: 'uint64', exact_recording_source: undefined,
    }))).not.toBeNull();
    expect(parseStudioV2WriteGroup(longGroup('instantaneous', {
      source_value_type: 'float64', result_value_type: 'float64', exact_recording_source: undefined,
    }))).not.toBeNull();
    const custom = longGroup('instantaneous', {
      exact_recording_source: undefined,
    }, {
      layout_version: undefined,
      recording_profile: undefined,
      storage_mode: undefined,
      destination: { ...(preparedGroupFixture.destination as Record<string, unknown>), storage_strategy: 'custom' },
    });
    expect(parseStudioV2WriteGroup(custom)).not.toBeNull();
  });
});
