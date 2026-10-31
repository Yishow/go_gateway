import { GroupMemberTable, type GroupMemberTableProps } from '../steps/step4/writeGroup/GroupMemberTable';
import { sourceKey } from './recordingPickerState';
import type { EditorMember } from '../state/writeGroup/draft';
export interface RecordingSelectedMembersProps extends Omit<GroupMemberTableProps, 'excluded'> {
  labels: { remove: (member: EditorMember) => string; unavailable: string; undo: string };
  onRemove: (member: EditorMember) => void;
  undo?: { member: EditorMember; onUndo: () => void; disabled?: boolean };
}
/** Presentation only: the live editor owns draft updates, removal history and restore validation. */
export function RecordingSelectedMembers({ labels, onRemove, undo, ...props }: RecordingSelectedMembersProps) {
  const candidates = props.candidates.filter(candidate => props.members.some(member => member.key === candidate.key && sourceKey(member) === sourceKey(candidate)) &&
    props.candidates.filter(other => sourceKey(other) === sourceKey(candidate)).length === 1);
  const missing = props.members.filter(member => !candidates.some(candidate => member.key === candidate.key && sourceKey(member) === sourceKey(candidate)));
  return <section className="recording-selected-members" data-testid="recording-selected-members">
    <GroupMemberTable {...props} candidates={candidates} excluded={[]}
      selectionControls={{ removeLabel: labels.remove, onRemove }} />
    {missing.length > 0 && <ul>{missing.map(member => <li key={sourceKey(member)} data-testid={`recording-missing-member-${member.point_id}`}>
      <span>{labels.unavailable}</span> <span>{member.device_id} · {member.point_id} · {member.tag_id}</span>
      <button type="button" className="gw-secondary" disabled={props.disabled} onClick={() => onRemove(member)}>{labels.remove(member)}</button>
    </li>)}</ul>}
    {undo && <button type="button" className="gw-secondary" disabled={props.disabled || undo.disabled} onClick={undo.onUndo}>{labels.undo}</button>}
  </section>;
}
