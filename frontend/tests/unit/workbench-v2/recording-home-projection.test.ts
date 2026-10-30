import { describe, expect, it, vi } from 'vitest';
import { recordingHomeProjection, type RecordingHomeInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeProjection';

const scope = { workspaceId: 'W1', workspaceRevision: 'R2', epoch: 2 };
function input(): RecordingHomeInput {
  return { scope, workspace: { state: 'available', scope, value: {
    id: 'W1', kind: 'single', status: 'empty', ordered_device_ids: [], created_at: '', updated_at: '',
  } }, sources: { state: 'available', scope, value: { candidates: [], excluded: [] } },
  groups: { state: 'available', scope, value: { workspace_id: 'W1', workspace_revision: 'R2', groups: [] } } };
}
function freeze<T>(value: T): T {
  if (value && typeof value === 'object') { Object.values(value).forEach(freeze); Object.freeze(value); }
  return value;
}

describe('recording home canonical evidence boundary', () => {
  it('opening is a pure read, retaining valid empty evidence without claiming healthy or commanding mutations', () => {
    const data = freeze(input());
    const before = JSON.stringify(data);
    const network = vi.spyOn(globalThis, 'fetch').mockImplementation(() => { throw new Error('No I/O on opening'); });
    try {
      const result = recordingHomeProjection(data);
      expect(result.workspace).toEqual({ state: 'available', value: {
        id: 'W1', kind: 'single', status: 'empty', ordered_device_ids: [], created_at: '', updated_at: '',
      } });
      expect(result.sources).toEqual({ state: 'available', value: { candidates: [], excluded: [] } });
      expect(result.groups).toEqual({ state: 'available', value: [] });
      expect(result.primaryAction).toBeNull();
      expect(JSON.stringify(data)).toBe(before);
      expect(network).not.toHaveBeenCalled();
      expect(Object.isFrozen(result)).toBe(true);
      expect(result.scope).not.toBe(data.scope);
    } finally { network.mockRestore(); }
  });
});

describe('scope rejection and unavailable evidence', () => {
  it.each([
    [{ workspaceId: 'old' }, 'scope_mismatch'],
    [{ workspaceRevision: 'R1' }, 'revision_mismatch'],
    [{ epoch: 1 }, 'stale'],
  ] as const)('rejects a late envelope %j across every workspace authority', (patch, reason) => {
    const data = input();
    const oldScope = { ...scope, ...patch };
    const result = recordingHomeProjection({ ...data,
      workspace: { ...data.workspace!, scope: oldScope } as typeof data.workspace,
      sources: { ...data.sources!, scope: oldScope } as typeof data.sources,
      groups: { ...data.groups!, scope: oldScope } as typeof data.groups,
    });
    expect(result.workspace).toEqual({ state: 'unavailable', reason });
    expect(result.sources).toEqual({ state: 'unavailable', reason });
    expect(result.groups).toEqual({ state: 'unavailable', reason });
    expect(result.primaryAction).toBeNull();
  });
  it.each(['loading', 'error', 'unknown', 'unavailable'] as const)('keeps %s distinct from empty', state => {
    const result = recordingHomeProjection({ scope, workspace: { state, scope }, sources: { state, scope }, groups: { state, scope } });
    const expected = { state: 'unavailable', reason: state === 'unavailable' ? 'missing' : state };
    expect(result.workspace).toEqual(expected);
    expect(result.sources).toEqual(expected);
    expect(result.groups).toEqual(expected);
    expect(result.primaryAction).toBeNull();
  });
  it.each([
    [{ workspaceId: '' }, 'scope_mismatch'],
    [{ workspaceRevision: undefined as unknown as string }, 'revision_mismatch'],
    [{ epoch: -1 }, 'invalid_evidence'],
    [{ epoch: NaN }, 'invalid_evidence'],
    [{ epoch: 0.5 }, 'invalid_evidence'],
  ] as const)('rejects unprovable selection %j', (patch, reason) => {
    expect(recordingHomeProjection({ ...input(), scope: { ...scope, ...patch } }).workspace)
      .toEqual({ state: 'unavailable', reason });
  });
  it('checks canonical workspace and list identity rather than trusting their envelopes', () => {
    const data = input();
    expect(recordingHomeProjection({ ...data, workspace: { state: 'available', scope, value: {
      id: 'foreign', kind: 'single', status: 'empty', ordered_device_ids: [], created_at: '', updated_at: '',
    } } }).workspace).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
    expect(recordingHomeProjection({ ...data, groups: { state: 'available', scope, value: {
      workspace_id: 'W1', workspace_revision: 'old', groups: [],
    } } }).groups).toEqual({ state: 'unavailable', reason: 'revision_mismatch' });
    expect(recordingHomeProjection({ ...data, groups: { state: 'available', scope, value: {
      workspace_id: 'foreign', workspace_revision: 'R2', groups: [],
    } } }).groups).toEqual({ state: 'unavailable', reason: 'scope_mismatch' });
  });
});

describe.each(['loading', 'error', 'unknown', 'unavailable'] as const)('scoped %s degradation', state => {
  it.each([
    [{ workspaceId: 'old' }, 'scope_mismatch'], [{ workspaceRevision: 'old' }, 'revision_mismatch'], [{ epoch: 1 }, 'stale'],
  ] as const)('does not project stale degradation from %j', (patch, reason) => {
    const evidence = { state, scope: { ...scope, ...patch } };
    const result = recordingHomeProjection({ scope, workspace: evidence, sources: evidence, groups: evidence });
    expect(result.workspace).toEqual({ state: 'unavailable', reason });
    expect(result.sources).toEqual({ state: 'unavailable', reason });
    expect(result.groups).toEqual({ state: 'unavailable', reason });
  });
});
