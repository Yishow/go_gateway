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
