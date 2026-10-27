import { describe, expect, it, vi } from 'vitest';
import { ensureOrUpdateCircuitProfileGroups } from '@/features/datalink/workbench-v2/steps/step4/circuitProfileActions';
import { buildGroupCandidates, type GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type {
  WriteGroup,
  WriteGroupCreateRequest,
  WriteGroupMutationResponse,
  WriteGroupSampleReduction,
  WriteGroupUpdateRequest,
} from '@/types/studioV2WriteGroup';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';

const exactCounterSource = {
  measurement_id: 'meas-kwh', definition_revision: 'rev-1', source_binding_revision: 'rev-1',
  series_epoch: 'epoch-1', original_data_type: 'uint64' as const, scale: '0.1', unit: 'kWh',
  semantic_kind: 'counter' as const,
};

const allCandidates = buildGroupCandidates(batchPanelState()).candidates;

// The two sources the saved profile already holds; the prepare action reuses that writer.
const candidates = allCandidates.filter((candidate) => candidate.point_id === 'pt-0' || candidate.point_id === 'pt-1');

function candidatesFor(pointIds: string[]): GroupCandidate[] {
  return allCandidates.filter((candidate) => pointIds.includes(candidate.point_id));
}

function savedProfileGroup(reduction?: WriteGroupSampleReduction): WriteGroup {
  const base = longProfileGroup({ id: 'inst-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] });
  return {
    ...base,
    // A saved mean always states each source's role; a profile saved without a reduction has none to carry.
    ...(reduction === 'mean' ? { members: base.members.map((member) => ({ ...member, value_role: 'gauge' as const })) } : {}),
    recording_profile: reduction === undefined
      ? { kind: 'instantaneous', interval: 15, unit: 'second' }
      : { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: reduction },
  };
}

/** A saved mean profile whose members carry the persisted roles, value types and confirmed counter binding. */
function typedSavedGroup(): WriteGroup {
  const base = longProfileGroup({ id: 'inst-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] });
  return {
    ...base,
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' },
    members: [
      { ...base.members[0], value_role: 'gauge', source_value_type: 'int16', result_value_type: 'decimal' },
      { ...base.members[1], value_role: 'counter_raw', source_value_type: 'uint64', result_value_type: 'decimal', measurement_id: 'meas-kwh', exact_recording_source: exactCounterSource },
      {
        device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2', source_revision: 's2', mapping_revision: 'm2',
        target_column: '', required: false, circuit_id: 'cir-a', circuit_name: 'Feeder A',
        value_role: 'state', source_value_type: 'string', result_value_type: 'text',
      },
    ],
  };
}

/** The basic recording panel's real prepare action, with the reviewed period moved to 30 seconds. */
function prepareInput(savedLongGroups: WriteGroup[], profileCandidates: GroupCandidate[] = candidates) {
  const updateWriteGroup = vi.fn(async ({ id, request }: {
    id: string; request: WriteGroupUpdateRequest;
  }): Promise<WriteGroupMutationResponse> => ({
    workspace_revision: 'w-2', group: { ...(request.group as WriteGroup), id, revision: 'rev-2' },
  }));
  const ensureCircuitProfile = vi.fn(async ({ request }: {
    request: WriteGroupCreateRequest;
  }): Promise<WriteGroupMutationResponse> => ({
    workspace_revision: 'w-2', group: { ...(request.group as WriteGroup), id: 'inst-new', revision: 'rev-new' },
  }));
  return {
    workspaceId: 'ws-1', initialWorkspaceRevision: 'w-1', deviceName: 'Meter A',
    enabledProfiles: [{ profile: 'instantaneous' as const, intervalValue: 30, intervalSeconds: 30, candidates: profileCandidates }],
    storageMode: 'plain' as const, recreateMissing: true,
    circuitId: 'cir-a', circuitName: 'Feeder A',
    destination: {
      connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: '', table_name: '', storage_strategy: 'managed' as const,
    },
    measurementsByPointKey: new Map(),
    savedLongGroups, preparedGroups: [],
    ensureCircuitProfile, updateWriteGroup,
  };
}

describe('the prepare action keeps a saved reduction when only the period changes', () => {
  it('sends the already saved mean in the update payload', async () => {
    const context = prepareInput([savedProfileGroup('mean')]);
    await ensureOrUpdateCircuitProfileGroups(context);

    expect(context.ensureCircuitProfile).not.toHaveBeenCalled();
    expect(context.updateWriteGroup).toHaveBeenCalledTimes(1);
    const { id, request } = context.updateWriteGroup.mock.calls[0][0];
    expect(id).toBe('inst-a');
    expect(request.expected_group_revision).toBe('rev-inst-a-1');
    expect(request.group.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' });
    expect(request.group.row_policy.interval_seconds).toBe(30);
  });

  it('sends an already saved explicit latest unchanged', async () => {
    const context = prepareInput([savedProfileGroup('latest')]);
    await ensureOrUpdateCircuitProfileGroups(context);

    expect(context.updateWriteGroup).toHaveBeenCalledTimes(1);
    expect(context.updateWriteGroup.mock.calls[0][0].request.group.recording_profile)
      .toEqual({ kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'latest' });
  });

  it('does not invent a reduction for a profile saved without one', async () => {
    const context = prepareInput([savedProfileGroup()]);
    await ensureOrUpdateCircuitProfileGroups(context);

    expect(context.updateWriteGroup).toHaveBeenCalledTimes(1);
    const group = context.updateWriteGroup.mock.calls[0][0].request.group;
    expect(group.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second' });
    expect(group.recording_profile).not.toHaveProperty('sample_reduction');
  });

  it('reuses the saved writer without a PUT when the period and reduction already match', async () => {
    const saved = savedProfileGroup('mean');
    const context = prepareInput([{ ...saved }]);
    const result = await ensureOrUpdateCircuitProfileGroups({
      ...context,
      enabledProfiles: [{ profile: 'instantaneous' as const, intervalValue: 15, intervalSeconds: 15, candidates }],
    });

    expect(context.updateWriteGroup).not.toHaveBeenCalled();
    expect(context.ensureCircuitProfile).not.toHaveBeenCalled();
    expect(result.groups.map((group) => group.id)).toEqual(['inst-a']);
  });

  it('carries saved member metadata and the counter binding when the period changes', async () => {
    const context = prepareInput([typedSavedGroup()], candidatesFor(['pt-0', 'pt-1', 'pt-2']));
    await ensureOrUpdateCircuitProfileGroups(context);

    expect(context.updateWriteGroup).toHaveBeenCalledTimes(1);
    const group = context.updateWriteGroup.mock.calls[0][0].request.group;
    expect(group.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' });
    expect(group.members.map((member) => ({
      point: member.point_id, role: member.value_role, source: member.source_value_type,
      result: member.result_value_type, measurement: member.measurement_id,
    }))).toEqual([
      { point: 'pt-0', role: 'gauge', source: 'int16', result: 'decimal', measurement: '' },
      { point: 'pt-1', role: 'counter_raw', source: 'uint64', result: 'decimal', measurement: 'meas-kwh' },
      { point: 'pt-2', role: 'state', source: 'string', result: 'text', measurement: '' },
    ]);
    expect(group.members[1].exact_recording_source).toEqual(exactCounterSource);
  });

  it('carries an exact source proof through a cumulative profile period update', async () => {
    const base = longProfileGroup({ id: 'cum-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 10, points: [0, 1] });
    const saved: WriteGroup = {
      ...base,
      row_policy: { ...base.row_policy, interval_seconds: 600 },
      recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' },
      members: [{
        ...base.members[0], measurement_id: 'meas-kwh', source_value_type: 'uint64',
        result_value_type: 'decimal', exact_recording_source: exactCounterSource,
      }],
    };
    const profileCandidates = candidatesFor(['pt-0']);
    const context = prepareInput([saved], profileCandidates);
    await ensureOrUpdateCircuitProfileGroups({
      ...context,
      enabledProfiles: [{ profile: 'cumulative', intervalValue: 20, intervalSeconds: 1200, candidates: profileCandidates }],
    });

    expect(context.updateWriteGroup).toHaveBeenCalledTimes(1);
    const member = context.updateWriteGroup.mock.calls[0][0].request.group.members[0];
    expect(member.value_role).toBeUndefined();
    expect(member.measurement_id).toBe('meas-kwh');
    expect(member.exact_recording_source).toEqual(exactCounterSource);
  });

  it('gives a newly added source the role its type shows, and no value type or counter binding', async () => {
    const context = prepareInput([typedSavedGroup()], candidatesFor(['pt-0', 'pt-1', 'pt-2', 'pt-3']));
    await ensureOrUpdateCircuitProfileGroups(context);

    const added = context.updateWriteGroup.mock.calls[0][0].request.group.members.find((member) => member.point_id === 'pt-3');
    expect(added).toBeDefined();
    // A number is shown as a measurement to average; it is never turned into a counter by this action.
    expect(added?.value_role).toBe('gauge');
    expect(added).not.toHaveProperty('source_value_type');
    expect(added).not.toHaveProperty('result_value_type');
    expect(added?.measurement_id).toBe('');
  });

  it('does not invent metadata for a saved member that carries none', async () => {
    const context = prepareInput([savedProfileGroup()]);
    await ensureOrUpdateCircuitProfileGroups(context);

    const member = context.updateWriteGroup.mock.calls[0][0].request.group.members[0];
    expect(member).not.toHaveProperty('value_role');
    expect(member).not.toHaveProperty('source_value_type');
    expect(member).not.toHaveProperty('result_value_type');
  });

  it('creates a new instantaneous profile as an explicit mean, and a cumulative one without any reduction', async () => {
    const context = prepareInput([]);
    await ensureOrUpdateCircuitProfileGroups({
      ...context,
      enabledProfiles: [
        ...context.enabledProfiles,
        { profile: 'cumulative' as const, intervalValue: 10, intervalSeconds: 600, candidates },
      ],
    });

    expect(context.updateWriteGroup).not.toHaveBeenCalled();
    expect(context.ensureCircuitProfile).toHaveBeenCalledTimes(2);
    const [instantaneous, cumulative] = context.ensureCircuitProfile.mock.calls.map((call) => call[0].request.group);
    expect(instantaneous.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second', sample_reduction: 'mean' });
    expect(cumulative.recording_profile).toEqual({ kind: 'cumulative', interval: 10, unit: 'minute' });
    expect(cumulative.members.every((member) => member.value_role === undefined)).toBe(true);
  });

  it('keeps a profile created for a reviewed cutover on its latest values when the caller says so', async () => {
    const context = prepareInput([]);
    await ensureOrUpdateCircuitProfileGroups({ ...context, newReduction: 'latest' });

    const group = context.ensureCircuitProfile.mock.calls[0][0].request.group;
    expect(group.recording_profile).toEqual({ kind: 'instantaneous', interval: 30, unit: 'second' });
    expect(group.members.every((member) => member.value_role === undefined)).toBe(true);
  });
});
