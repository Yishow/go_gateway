import { describe, expect, it } from 'vitest';
import type { EditorMember } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import type { GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { addRecordingMembers, recordingPickerProjection } from '@/features/datalink/workbench-v2/recording-editor/recordingMemberSelection';
import { removeMember, sourceKey } from '@/features/datalink/workbench-v2/recording-editor/recordingPickerState';
const row: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'tag', label: 'counter name is not evidence', device_name: 'PLC', address: '40001', target_type: 'float32' };
const other = { ...row, key: 'P2|T2', point_id: 'P2', tag_id: 'T2' };
describe('recording picker canonical member intent adapter', () => {
  it('retains opaque member method/invalid interval and all hidden fields through remove then same ID add', () => {
    const member = { ...row, entity_key: 'raw entity', target_column: 'missing-column', required: false, measurement_id: 'measurement', method: 'counter_delta', intervalText: ' 1x ' } satisfies EditorMember & { method: string; intervalText: string };
    const removed = removeMember([member], { scope: 'S', removed: new Map() }, sourceKey(member));
    const projection = recordingPickerProjection('S', [row, other], []);
    const added = addRecordingMembers(removed.members, removed.history, [row, other], projection);
    expect(added?.[0]).toEqual(member); expect(added?.[0]).toBe(member);
    expect(added?.[1]).toEqual({ key: 'P2|T2', device_id: 'D', point_id: 'P2', tag_id: 'T2', entity_key: '', target_column: '', required: true });
    expect(projection.evidence[sourceKey(row)]?.semantic).toBeUndefined();
  });
  it('rejects whole additions if missing, duplicate, stale or foreign scope; never silently drops intent', () => {
    const history = { scope: 'S', removed: new Map<string, EditorMember>() };
    expect(addRecordingMembers([], history, [row, other], recordingPickerProjection('S', [row], []))).toBeUndefined();
    expect(addRecordingMembers([], history, [row], recordingPickerProjection('S', [row, row], []))).toBeUndefined();
    expect(addRecordingMembers([], history, [row], recordingPickerProjection('OTHER', [row], []))).toBeUndefined();
    const projection = recordingPickerProjection('S', [row], [], { [sourceKey(row)]: { availability: 'stale' } });
    expect(addRecordingMembers([], history, [row], projection)).toBeUndefined();
  });
});
