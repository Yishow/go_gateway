import { describe, expect, it } from 'vitest';
import { activationAction } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeHelpers';
import type { WriteGroupSuccessorActivationEnvelope } from '@/types/studioV2WriteGroupSuccessor';
describe('activation dispatch does not conflate first start with a cutover', () => {
  it('dispatches explicit v2 first_start separately', () => {
    expect(activationAction({ version: 'successor-activation-review-v2', activation_kind: 'first_start' } as WriteGroupSuccessorActivationEnvelope)).toBe('confirmFirstStart');
  });
  it('uses missing kind only for legal legacy v1', () => {
    expect(activationAction({ version: 'successor-activation-review-v1' } as WriteGroupSuccessorActivationEnvelope)).toBe('confirmCutover');
    expect(() => activationAction({ version: 'successor-activation-review-v2' } as WriteGroupSuccessorActivationEnvelope)).toThrow();
  });
});

import { changeTargetLayout, sameBoundary } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeHelpers';
import { parseWriteGroupSuccessorCandidate } from '@/utils/studioV2WriteGroupDefinitionJson';
import { fixture } from '../../fixtures/recordingChange';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
it('compares Go-serialized whole seconds without rewriting reviewed precision', () => {
  expect(sameBoundary('2099-01-01T00:00:00.000Z', '2099-01-01T00:00:00Z')).toBe(true);
  expect(sameBoundary('2099-01-01T00:00:00.000001Z', '2099-01-01T00:00:00Z')).toBe(false);
});
it('preserves complete source meaning when producing backend-shaped bidirectional profile intents', () => {
  const { seed } = fixture();
  const wide = { ...seed, basic_managed_device_id: undefined, layout_version: 'managed_wide_profile_v1' as const, recording_profile: { kind: 'instantaneous' as const, interval: 15, unit: 'second' as const, sample_reduction: 'latest' as const } };
  const long = changeTargetLayout(wide, 'managed_long_v1');
  expect(parseWriteGroupSuccessorCandidate(long)).not.toBeNull(); expect(long.storage_mode).toBe('plain'); expect(long.members).toEqual(wide.members); expect(long.recording_profile).toEqual(wide.recording_profile);
  const back = changeTargetLayout({ ...long, write_policy: { ...long.write_policy, recreate_missing: true } }, 'managed_wide_profile_v1');
  expect(parseWriteGroupSuccessorCandidate(back)).not.toBeNull(); expect(back.storage_mode).toBeUndefined(); expect(back.write_policy.recreate_missing).toBeUndefined(); expect(back.recording_profile).toEqual(wide.recording_profile);
});

function mixedTarget(): WriteGroup {
  return { ...fixture().seed, basic_managed_device_id: undefined, record_id: '11111111-1111-4111-8111-111111111111',
    layout_version: 'managed_wide_profile_v1', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' },
    row_policy: { allowed_lateness_seconds: 4, incomplete_policy: 'partial', record_key_column: 'record_id', provenance_column: 'provenance', due_at_column: 'due_at', group_id_column: 'group_id' },
    members: [
      { device_id: 'D', point_id: 'P', tag_id: 'T', source_revision: 'S', mapping_revision: 'M', target_column: 'value', required: true, max_age_seconds: 20, calculation: 'mean', cadence: { interval: 5, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'float64', result_value_type: 'float64', value_role: 'gauge' },
      { device_id: 'D', point_id: 'P2', tag_id: 'T2', source_revision: 'S2', mapping_revision: 'M2', target_column: 'counter', required: false, calculation: 'delta', cadence: { interval: 1, unit: 'minute', alignment: 'utc_epoch' }, source_value_type: 'int64', result_value_type: 'int64', value_role: 'counter_raw', measurement_id: 'measurement' },
    ],
  };
}
describe('mixed successor layout intent', () => {
  it('maps wide due-at bindings to long without changing record or per-member meaning', () => {
    const wide = mixedTarget(), before = structuredClone(wide), long = changeTargetLayout(wide, 'managed_long_v1');
    expect(long.row_policy).toEqual({ ...wide.row_policy, due_at_column: undefined, bucket_start_column: 'bucket_start' });
    expect(long.row_policy).not.toHaveProperty('due_at_column');expect(long.row_policy).not.toHaveProperty('interval_seconds');
    expect(long.storage_mode).toBe('plain');expect(long.destination.table_name).toBe('');
    expect(long.members).toEqual(wide.members);expect(long.members).not.toBe(wide.members);
    expect(long.record_id).toBe(wide.record_id);expect(long.recording_profile).toEqual(wide.recording_profile);
    expect(parseWriteGroupSuccessorCandidate(long)).toEqual(long);expect(wide).toEqual(before);
  });
  it('maps long back to due-at partial wide without introducing a legacy window or filling NULLs', () => {
    const long: WriteGroup = { ...mixedTarget(), layout_version: 'managed_long_v1', storage_mode: 'plain',
      row_policy: { allowed_lateness_seconds: 4, incomplete_policy: 'partial', record_key_column: 'record_id', provenance_column: 'provenance', bucket_start_column: 'bucket_start', group_id_column: 'group_id' }, write_policy: { mode: 'append', recreate_missing: true } };
    const before = structuredClone(long), wide = changeTargetLayout(long, 'managed_wide_profile_v1');
    expect(wide.row_policy).toEqual(mixedTarget().row_policy);expect(wide.row_policy).not.toHaveProperty('bucket_start_column');
    expect(wide.row_policy).not.toHaveProperty('interval_seconds');expect(wide.storage_mode).toBeUndefined();expect(wide.write_policy.recreate_missing).toBeUndefined();
    expect(wide.members).toEqual(long.members);expect(wide.record_id).toBe(long.record_id);expect(wide.recording_profile).toEqual(long.recording_profile);
    expect(parseWriteGroupSuccessorCandidate(wide)).toEqual(wide);expect(long).toEqual(before);
  });
  it('does not reinterpret skip-row or scaled counter proof as a supported wide target', () => {
    const long = { ...mixedTarget(), layout_version: 'managed_long_v1' as const, storage_mode: 'plain' as const,
      row_policy: { allowed_lateness_seconds: 4, incomplete_policy: 'partial', record_key_column: 'record_id', provenance_column: 'provenance', bucket_start_column: 'bucket_start' } };
    const skipped = { ...long, row_policy: { ...long.row_policy, incomplete_policy: 'skip_row' } };
    expect(changeTargetLayout(skipped, 'managed_wide_profile_v1')).toBe(skipped);
    const scaled = structuredClone(long);Object.assign(scaled.members[1], { source_value_type: 'uint64', result_value_type: 'decimal', exact_recording_source: { measurement_id: 'measurement', definition_revision: 'definition', source_binding_revision: 'binding', series_epoch: 'epoch', original_data_type: 'uint64', scale: '0.1', unit: 'kWh', semantic_kind: 'counter' } });
    const before = structuredClone(scaled);expect(changeTargetLayout(scaled, 'managed_wide_profile_v1')).toBe(scaled);expect(scaled).toEqual(before);
  });
});
