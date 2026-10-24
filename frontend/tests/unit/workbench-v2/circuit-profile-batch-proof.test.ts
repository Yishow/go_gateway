import { beforeEach, describe, expect, it } from 'vitest';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { candidateKey, defaultCircuitProfileLocalDraft, loadCircuitProfileLocalDraft, saveCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { applyUpdatedGroupsToActiveDraft, syncBatchDraftsAfterSave } from '@/features/datalink/workbench-v2/state/circuitProfileBatch';
import { classifyBatchMutationError, reconcileBatchResults, type BatchTargetExecutionResult } from '@/features/datalink/workbench-v2/state/circuitProfileBatchOutcome';

function originalGroup() {
  const base = savedGroup();
  return savedGroup({
    id: 'group-a', revision: 'group-rev-1', layout_version: 'managed_long_v1', storage_mode: 'plain',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
    destination: { ...base.destination, table_name: 'gw_long_owned', storage_strategy: 'managed' },
    row_policy: { ...base.row_policy, interval_seconds: 15 },
    members: base.members.map((m) => ({ ...m, circuit_id: 'cir-a', circuit_name: 'Circuit A' })),
  });
}
function unknownResults(): Record<string, BatchTargetExecutionResult> {
  return { 'cir-a': { circuit_id: 'cir-a', circuit_name: 'Circuit A', status: 'unknown',
    profiles: { instantaneous: { kind: 'instantaneous', status: 'unknown' } } } };
}

beforeEach(() => window.sessionStorage.clear());

describe('Batch proof must preserve the content the operator actually reviewed', () => {
  it.each(['members', 'destination'] as const)('does not resolve unknown when %s changed behind the same requested interval', (drift) => {
    const reviewed = originalGroup();
    const live = { ...reviewed, revision: 'group-rev-2',
      recording_profile: { kind: 'instantaneous' as const, interval: 30, unit: 'second' as const },
      row_policy: { ...reviewed.row_policy, interval_seconds: 30 },
      ...(drift === 'members' ? { members: reviewed.members.slice(0, 1) } :
        { destination: { ...reviewed.destination, connector_id: 'other-connector', connector_revision: 'other-rev' } }),
    };
    const result = reconcileBatchResults([{ circuit_id: 'cir-a', instantaneousGroup: reviewed }],
      { instantaneous: 30, cumulative: 10 }, unknownResults(), [live]);
    expect(result.results['cir-a'].profiles.instantaneous?.status).toBe('unknown');
    expect(result.proven).toEqual([]);
  });

  it('keeps an existing profile source draft intact when only its interval was batch edited', () => {
    const group = originalGroup();
    const candidates = buildGroupCandidates(savedGroupState()).candidates;
    const draft = { ...defaultCircuitProfileLocalDraft(candidates, [group], 'cir-a'),
      instantaneous_keys: [candidateKey(candidates[0])], cumulative_keys: [],
      selected_keys: [candidateKey(candidates[0])], circuit_name: 'Unsaved local name',
    };
    const saved = { ...group, revision: 'group-rev-2',
      recording_profile: { kind: 'instantaneous' as const, interval: 30, unit: 'second' as const },
      row_policy: { ...group.row_policy, interval_seconds: 30 },
    };
    saveCircuitProfileLocalDraft('ws-1', 'dev-1', draft);
    syncBatchDraftsAfterSave('ws-1', 'dev-1', candidates, [group], [saved]);
    const reloaded = loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, [saved], 'cir-a');
    expect(reloaded.instantaneous_interval_input).toBe('30');
    expect(reloaded.instantaneous_keys).toEqual(draft.instantaneous_keys);
    expect(reloaded.circuit_name).toBe('Unsaved local name');
    const active = applyUpdatedGroupsToActiveDraft(draft, 'cir-a', [saved]);
    expect(active.instantaneous_keys).toEqual(draft.instantaneous_keys);
  });

  it('keeps an untyped proxy 4xx unknown instead of treating its HTML body as an application rejection', () => {
    expect(classifyBatchMutationError({ response: { status: 400, data: '<html>proxy response</html>' } })).toBe('unknown');
  });
});
