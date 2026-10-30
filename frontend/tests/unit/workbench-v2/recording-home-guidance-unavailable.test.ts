import { describe, expect, it } from 'vitest';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture, group, groupScope, owners, scope } from './recording-home-guidance-fixtures';

describe('unavailable guidance is neither empty nor healthy', () => {
  it.each(['loading', 'error', 'unknown', 'unavailable'] as const)('shows actionable %s bootstrap', state => {
    const input = fixture();
    const result = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, workspace: { state, scope } } });
    expect(result.availability).toEqual({ state: 'unavailable', reason: state === 'unavailable' ? 'missing' : state });
    expect(result.reason).toBe('evidence_unavailable');
    expect(result.primaryAction).toEqual({ kind: 'navigate', scope, target: 'evidence' });
  });
  it.each(['sources', 'groups'] as const)('does not turn absent %s into an empty collection', key => {
    const input = fixture(['D1'], true);
    const result = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, [key]: undefined } });
    expect(result.availability).toEqual({ state: 'unavailable', reason: 'missing' });
    expect(result.reason).toBe('evidence_unavailable');
  });
  it('does not navigate with unprovable workspace identity', () => {
    const input = fixture();
    const result = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, scope: { ...scope, workspaceId: '' } } });
    expect(result.availability).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
    expect(result.primaryAction).toBeNull();
  });
  it.each([
    { workspaceId: 'old' }, { workspaceRevision: 'old' }, { epoch: 1 }, { groupId: 'old' },
    { groupRevision: 'old' }, { appliedRevision: 'old' },
  ])('rejects late owner evidence %j without showing an old operation action', patch => {
    const a = group('A'); const input = fixture(['D1'], true, [a]);
    const result = recordingHomeGuidance({ ...input, owners: { A: { state: 'available', scope: { ...groupScope(a), ...patch },
      value: [{ owner: 'start', scope: { ...groupScope(a), deviceId: 'D1' }, state: 'unresolved', reason: 'unknown', operationId: 'OLD' }],
    } } });
    expect(result.reason).toBe('evidence_unavailable');
    expect(result.primaryAction).toMatchObject({ target: 'recording', groupId: 'A', section: 'overview', intent: null });
    expect(JSON.stringify(result)).not.toContain('OLD');
    expect(result.recordingItems.map(item => item.id)).toEqual(['A']);
  });
  it('treats missing owners and delivery as unavailable, never healthy by saved running alone', () => {
    const a = group('A'); const input = fixture(['D1'], true, [a]);
    expect(recordingHomeGuidance({ ...input, owners: undefined }).reason).toBe('evidence_unavailable');
    expect(recordingHomeGuidance({ ...input, evidence: { ...input.evidence, delivery: undefined } }).reason).toBe('evidence_unavailable');
  });
  it('does not require absent command owners to continue a named draft', () => {
    const c = group('C', 'draft'); const input = fixture(['D1'], true, [c]);
    expect(recordingHomeGuidance({ ...input, owners: undefined }).reason).toBe('draft');
  });
  it('does not confuse foreign inner owner identity with a matching envelope', () => {
    const a = group('A'); const input = fixture(['D1'], true, [a]);
    const value = owners(a).map(item => ({ ...item, scope: { ...item.scope, workspaceId: 'foreign' } }));
    const result = recordingHomeGuidance({ ...input, owners: { A: { state: 'available', scope: groupScope(a), value } } });
    expect(result.reason).toBe('evidence_unavailable');
  });
});

describe('source navigation identity', () => {
  it('never navigates to a source outside the loaded workspace device set', () => {
    const input = fixture(['OTHER'], true);
    const result = recordingHomeGuidance(input);
    expect(result.availability).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
    expect(result.reason).toBe('evidence_unavailable');
    expect(result.primaryAction).toEqual({ kind: 'navigate', scope, target: 'evidence' });
  });
  it('does not produce a point configuration action with an empty device identity', () => {
    const result = recordingHomeGuidance(fixture(['']));
    expect(result.availability).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
    expect(result.primaryAction).toEqual({ kind: 'navigate', scope, target: 'evidence' });
  });
});
