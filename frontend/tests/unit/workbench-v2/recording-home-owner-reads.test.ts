import { beforeEach, describe, expect, it, vi } from 'vitest';
import { readRecordingHomeOwners } from '@/features/datalink/workbench-v2/recording-home/recordingHomeOwnerReads';
import { inspectSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { recordingTestWriteRecovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { group, groupScope } from './recording-home-guidance-fixtures';
import type { BasicRecordingDeviceInspection } from '@/features/datalink/workbench-v2/state/basicRecording';
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery', () => ({ inspectSchemaOperation: vi.fn() }));
vi.mock('@/features/datalink/workbench-v2/state/recordingTestWriteRecovery', () => ({ recordingTestWriteRecovery: { read: vi.fn(), retain: vi.fn(), resolve: vi.fn() } }));
beforeEach(() => {
  vi.mocked(inspectSchemaOperation).mockReturnValue({ state: 'empty' });
  vi.mocked(recordingTestWriteRecovery.read).mockReturnValue({ references: [], unavailable: false, unreadable: false });
});
describe('home readonly owner inspection', () => {
  it('inspects actual stores but never turns empty local storage into schema/start global none', () => {
    const value = readRecordingHomeOwners(groupScope(group('G')), group('G'));
    expect(value.state).toBe('available'); if (value.state !== 'available') throw new Error('fixture');
    expect(value.value.map(item => [item.owner, item.state])).toEqual([['schema', 'unavailable'], ['start', 'unavailable'], ['test_write', 'unavailable']]);
    expect(inspectSchemaOperation).toHaveBeenCalledWith('W1', 'G');
    expect(recordingTestWriteRecovery.resolve).not.toHaveBeenCalled(); expect(recordingTestWriteRecovery.retain).not.toHaveBeenCalled();
  });
  it('retained schema ID remains unresolved and storage failure cannot erase it', () => {
    vi.mocked(inspectSchemaOperation).mockReturnValue({ state: 'retained', operationId: 'O1' });
    vi.mocked(recordingTestWriteRecovery.read).mockReturnValue({ references: [], unavailable: true, unreadable: true });
    const value = readRecordingHomeOwners(groupScope(group('G')), group('G'));
    expect(value.state === 'available' && value.value[0]).toMatchObject({ owner: 'schema', state: 'unresolved', operationId: 'O1' });
  });
  it('only an exact captured complete workspace recovery read plus readable local owner proves test-write none', () => {
    const scope = groupScope(group('G'));
    const server = { state: 'available' as const, scope, value: { workspace_id: 'W1', complete: true as const, operations: [] } };
    const value = readRecordingHomeOwners(scope, group('G'), server);
    expect(value.state === 'available' && value.value[2]).toMatchObject({ owner: 'test_write', state: 'observed', result: 'none' });
    const stale = readRecordingHomeOwners(scope, group('G'), { ...server, scope: { ...scope, epoch: 99 } });
    expect(stale.state === 'available' && stale.value[2]).toMatchObject({ state: 'unavailable' });
  });
  it('foreign group and mismatched group revision fail before inspecting owners', () => {
    vi.mocked(inspectSchemaOperation).mockClear();
    expect(readRecordingHomeOwners(groupScope(group('G')), group('OTHER')).state).toBe('unknown');
    expect(inspectSchemaOperation).not.toHaveBeenCalled();
  });
});
it('a retained start is attributed only by its original validated group IDs, never by the first device member', () => {
  const scope = groupScope(group('G'));
  const unrelated = new Map([['D1', { unavailable: false, retained: [{ operationId: 'O', groupIds: ['OTHER'] }] }]]);
  const known = new Map([['D1', { unavailable: false, retained: [{ operationId: 'O', groupIds: ['G'] }] }]]);
  const first = readRecordingHomeOwners(scope, group('G'), undefined, unrelated);
  const second = readRecordingHomeOwners(scope, group('G'), undefined, known);
  expect(first.state === 'available' && first.value.find(item => item.owner === 'start')).toMatchObject({ state: 'unavailable' });
  expect(second.state === 'available' && second.value.find(item => item.owner === 'start')).toMatchObject({ state: 'unresolved', operationId: 'O' });
});
it('a completed Start observation never erases another retained intent or an unreadable owner', () => {
  const saved = { ...group('G'), applied_revision: 'G2' }, scope = groupScope(saved);
  const completed = { operationId: 'O1', groupIds: ['G'], completion: { workspaceRevision: scope.workspaceRevision,
    groups: [{ id: 'G', revision: 'G2', appliedRevision: 'G2' }] } };
  const inspect = (value: BasicRecordingDeviceInspection) => {
    const result = readRecordingHomeOwners(scope, saved, undefined, new Map([['D1', value]]));
    return result.state === 'available' ? result.value.find(item => item.owner === 'start') : undefined;
  };
  expect(inspect({ unavailable: false, retained: [completed] })).toMatchObject({ state: 'observed', result: 'resolved', operationId: 'O1' });
  expect(inspect({ unavailable: false, retained: [completed, { operationId: 'O2', groupIds: ['G'] }] })).toMatchObject({ state: 'unresolved', operationId: 'O2' });
  expect(inspect({ unavailable: true, retained: [completed] })).toMatchObject({ state: 'unavailable' });
  expect(inspect({ unavailable: true, retained: [completed, { groupIds: ['G'] }] })).toMatchObject({ state: 'unresolved', reason: 'lost_ack' });
});
it.each(['workspace', 'group', 'applied'])('completed Start metadata cannot bless a different %s revision', kind => {
  const saved = { ...group('G'), applied_revision: 'G2' }, scope = groupScope(saved);
  const completion = { workspaceRevision: kind === 'workspace' ? 'OTHER' : scope.workspaceRevision,
    groups: [{ id: 'G', revision: kind === 'group' ? 'OTHER' : 'G2', appliedRevision: kind === 'applied' ? 'OTHER' : 'G2' }] };
  const result = readRecordingHomeOwners(scope, saved, undefined, new Map([['D1', { unavailable: false, retained: [{ operationId: 'O1', groupIds: ['G'], completion }] }]]));
  expect(result.state === 'available' && result.value.find(item => item.owner === 'start')).toMatchObject({ state: 'unresolved', operationId: 'O1' });
});
