import { describe, expect, it } from 'vitest';
import { acceptRecordingHomeResponse, type RecordingHomeEvidence } from '@/features/datalink/workbench-v2/recording-home/recordingHomeProjection';
const scope = { workspaceId: 'W1', workspaceRevision: 'R2', epoch: 2 };
const current: RecordingHomeEvidence<string[]> = { state: 'available', scope, value: ['current recording'] };

describe('late home response acceptance', () => {
  it.each([
    { workspaceId: 'old' }, { workspaceRevision: 'old' }, { epoch: 1 },
  ])('preserves current evidence after late success or failure %j', patch => {
    const old = { ...scope, ...patch };
    expect(acceptRecordingHomeResponse(scope, old, { state: 'available', scope: old, value: ['old recording'] }, current))
      .toBe(current);
    expect(acceptRecordingHomeResponse(scope, old, { state: 'error', scope: old }, current)).toBe(current);
  });
  it('does not retain evidence from the previous scope or let an envelope misbind the response', () => {
    const old = { ...scope, epoch: 1 };
    expect(acceptRecordingHomeResponse(scope, old, { state: 'error', scope: old }, { ...current, scope: old })).toBeUndefined();
    expect(acceptRecordingHomeResponse(scope, scope, { state: 'available', scope: old, value: ['old'] }, current)).toBe(current);
    expect(acceptRecordingHomeResponse(scope, old, { state: 'error', scope: old })).toBeUndefined();
  });
  it('forwards matching success and degradation instead of retaining stale success', () => {
    const fresh = { state: 'available' as const, scope, value: ['new recording'] };
    expect(acceptRecordingHomeResponse(scope, scope, fresh, current)).toBe(fresh);
    expect(acceptRecordingHomeResponse(scope, scope, { state: 'error', scope }, current)).toEqual({ state: 'error', scope });
  });
  it('rejects invalid current selection even when both envelopes repeat it', () => {
    const invalid = { ...scope, epoch: -1 };
    expect(acceptRecordingHomeResponse(invalid, invalid, { state: 'error', scope: invalid }, current)).toBeUndefined();
  });
});

describe.each(['loading', 'error', 'unknown', 'unavailable'] as const)('current %s evidence', state => {
  it.each([{ workspaceId: 'old' }, { workspaceRevision: 'old' }, { epoch: 1 }])(
    'survives late success and failure from %j without retaining old degradation', patch => {
      const old = { ...scope, ...patch };
      const degraded = { state, scope };
      const previous = { state, scope: old };
      for (const response of [
        { state: 'available' as const, scope: old, value: ['old recording'] },
        { state: 'error' as const, scope: old },
      ]) {
        expect(acceptRecordingHomeResponse(scope, old, response, degraded)).toBe(degraded);
        expect(acceptRecordingHomeResponse(scope, old, response, previous)).toBeUndefined();
      }
    });
  it('does not accept an old failure relabeled with the current request scope', () => {
    const degraded = { state, scope };
    const response = { state, scope: { ...scope, epoch: 1 } };
    expect(acceptRecordingHomeResponse(scope, scope, response, degraded)).toBe(degraded);
  });
});
