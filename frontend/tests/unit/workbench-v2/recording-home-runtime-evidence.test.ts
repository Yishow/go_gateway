import { describe, expect, it } from 'vitest';
import { recordingHomeRuntimeGroup } from '@/features/datalink/workbench-v2/recording-home/recordingHomeRuntimeEvidence';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
import { group, scope } from './recording-home-guidance-fixtures';
const saved = { ...group('G'), revision: 'C', applied_revision: 'B' };
const runtime = () => ({ state: 'available' as const, scope, value: { ok: true as const, value: ruleEvidenceFixture({ workspace_id: 'W1', workspace_revision: 'W2', device_id: 'D1' }) } });
describe('home reuses current effective snapshot without relabeling requested revision', () => {
  it('preserves saved C / requested B / currently effective A and pending B boundary', () => {
    expect(recordingHomeRuntimeGroup(scope, saved, runtime())).toMatchObject({ saved_revision: 'C', applied_revision: 'B', effective_revision: 'A', applied_effective_at: '2026-10-27T10:00:01Z' });
  });
  it.each(['epoch', 'workspaceRevision', 'workspaceId'] as const)('rejects obsolete captured %s', field => {
    const response = runtime();
    expect(recordingHomeRuntimeGroup({ ...scope, [field]: field === 'epoch' ? 999 : 'other' }, saved, response)).toBeUndefined();
  });
  it('does not promote stale, wrong member or changed desired head to effective', () => {
    const response = runtime(); response.value.value.freshness = 'stale';
    expect(recordingHomeRuntimeGroup(scope, saved, response)).toBeUndefined();
    expect(recordingHomeRuntimeGroup(scope, { ...saved, revision: 'NEW' }, runtime())).toBeUndefined();
    expect(recordingHomeRuntimeGroup(scope, { ...saved, members: [] }, runtime())).toBeUndefined();
  });
});
