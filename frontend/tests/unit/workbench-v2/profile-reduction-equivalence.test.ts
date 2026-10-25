import { describe, expect, it } from 'vitest';
import { canonicalProfileGroupMatchesDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { profileMatchesIntent } from '@/features/datalink/workbench-v2/state/circuitProfileBatchIntent';
import type {
  WriteGroup,
  WriteGroupDraft,
  WriteGroupRecordingProfile,
  WriteGroupSampleReduction,
} from '@/types/studioV2WriteGroup';
import { longProfileGroup } from '../../fixtures/circuitBatchPanel';

/** The same managed long profile with the contract's literal reduction values only. */
function profile(reduction?: WriteGroupSampleReduction): WriteGroupRecordingProfile {
  return reduction === undefined
    ? { kind: 'instantaneous', interval: 15, unit: 'second' }
    : { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: reduction };
}

function savedGroup(reduction?: WriteGroupSampleReduction): WriteGroup {
  return {
    ...longProfileGroup({ id: 'inst-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] }),
    recording_profile: profile(reduction),
  };
}

function draft(reduction?: WriteGroupSampleReduction): WriteGroupDraft {
  return {
    workspace_id: 'ws-1',
    name: 'Feeder A (Instantaneous)',
    layout_version: 'managed_long_v1',
    storage_mode: 'plain',
    recording_profile: profile(reduction),
    destination: {
      connector_id: 'conn-1', connector_revision: 'crev-1', table_schema: '', table_name: '', storage_strategy: 'managed',
    },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
    members: [0, 1].map((index) => ({
      device_id: 'dev-1', point_id: `pt-${index}`, tag_id: `tag-${index}`,
      target_column: '', required: false, circuit_id: 'cir-a', circuit_name: 'Feeder A',
    })),
  };
}

describe('canonicalProfileGroupMatchesDraft compares the effective reduction', () => {
  it('treats an omitted reduction as an explicit latest in both directions', () => {
    expect(canonicalProfileGroupMatchesDraft(savedGroup(), draft('latest'))).toBe(true);
    expect(canonicalProfileGroupMatchesDraft(savedGroup('latest'), draft())).toBe(true);
  });

  it('separates mean from latest and from an omitted reduction', () => {
    expect(canonicalProfileGroupMatchesDraft(savedGroup('mean'), draft('latest'))).toBe(false);
    expect(canonicalProfileGroupMatchesDraft(savedGroup('latest'), draft('mean'))).toBe(false);
    expect(canonicalProfileGroupMatchesDraft(savedGroup('mean'), draft())).toBe(false);
    expect(canonicalProfileGroupMatchesDraft(savedGroup(), draft('mean'))).toBe(false);
  });

  it('accepts the same explicit reduction', () => {
    expect(canonicalProfileGroupMatchesDraft(savedGroup('mean'), draft('mean'))).toBe(true);
  });
});

describe('profileMatchesIntent compares the effective reduction', () => {
  it('treats an omitted reduction as an explicit latest in both directions', () => {
    expect(profileMatchesIntent(savedGroup(), savedGroup('latest'))).toBe(true);
    expect(profileMatchesIntent(savedGroup('latest'), savedGroup())).toBe(true);
  });

  it('separates mean from latest and from an omitted reduction', () => {
    expect(profileMatchesIntent(savedGroup('mean'), savedGroup('latest'))).toBe(false);
    expect(profileMatchesIntent(savedGroup('latest'), savedGroup('mean'))).toBe(false);
    expect(profileMatchesIntent(savedGroup('mean'), savedGroup())).toBe(false);
    expect(profileMatchesIntent(savedGroup(), savedGroup('mean'))).toBe(false);
  });

  it('accepts the same explicit reduction', () => {
    expect(profileMatchesIntent(savedGroup('mean'), savedGroup('mean'))).toBe(true);
  });
});
