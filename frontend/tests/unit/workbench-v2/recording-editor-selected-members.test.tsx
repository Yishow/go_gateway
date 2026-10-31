import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RecordingSelectedMembers } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingSelectedMembers';
import type { EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'temperature', label: 'Temperature', device_name: 'Device', address: '40001', target_type: 'float64' };
const member: EditorMember = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', entity_key: 'raw entity', target_column: 'illegal column!', required: true, measurement_id: 'M' };
function props() {
  return { candidates: [candidate, { ...candidate, key: 'P2|T2', point_id: 'P2', tag_id: 'T2', label: 'Other' }], members: [member],
    columns: [{ name: 'temperature', type: 'DOUBLE PRECISION', primary_key: false, nullable: true }], metadataStatus: 'exists' as const,
    suggestions: { assignments: {}, unmatched: [] }, issues: [], showEntityKey: true, disabled: false, onChange: vi.fn(), onRemove: vi.fn(),
    labels: { remove: (value: EditorMember) => `Remove ${value.point_id}`, unavailable: 'Source unavailable', undo: 'Undo removal' } };
}
describe('selected recording members use the canonical draft', () => {
  it('shows selected rows only without competing addition controls and emits the complete remove intent', () => {
    const value = props(); render(<RecordingSelectedMembers {...value} />);
    expect(screen.getByText('Temperature')).toBeVisible(); expect(screen.queryByText('Other')).not.toBeInTheDocument();
    expect(screen.queryByTestId('group-member-include-P')).not.toBeInTheDocument();
    for (const name of ['include', 'exclude', 'accept']) expect(screen.queryByTestId(`group-bulk-${name}`)).not.toBeInTheDocument();
    expect(screen.getByTestId('group-member-column-P')).toHaveValue('illegal column!');
    fireEvent.click(screen.getByRole('button', { name: 'Remove P' }));
    expect(value.onRemove).toHaveBeenCalledWith(member); expect(value.onChange).not.toHaveBeenCalled();
  });
  it('preserves unrendered member intent on column edits and keeps normal metadata validation', () => {
    const value = props(); render(<RecordingSelectedMembers {...value} />);
    expect(screen.getByTestId('group-member-problem-P')).toBeVisible();
    fireEvent.change(screen.getByTestId('group-member-column-P'), { target: { value: 'temperature' } });
    expect(value.onChange).toHaveBeenCalledWith([{ ...member, target_column: 'temperature' }]);
  });
  it('shows missing identity without inventing a candidate or borrowing a same-key foreign source', () => {
    const value = props(); value.candidates = [{ ...candidate, device_id: 'FOREIGN', label: 'Do not borrow this label' }];
    render(<RecordingSelectedMembers {...value} />);
    expect(screen.getByText('Source unavailable')).toBeVisible();
    expect(screen.getByTestId('recording-missing-member-P')).toHaveTextContent('D · P · T');
    expect(screen.queryByText('Do not borrow this label')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Remove P' })); expect(value.onRemove).toHaveBeenCalledWith(member);
  });
  it('undo is parent-owned and disabled mutations stay disabled', () => {
    const value = props(), onUndo = vi.fn(); render(<RecordingSelectedMembers {...value} disabled undo={{ member, onUndo }} />);
    expect(screen.getByRole('button', { name: 'Remove P' })).toBeDisabled();
    expect(screen.getByRole('button', { name: 'Undo removal' })).toBeDisabled();
    expect(screen.getByTestId('group-member-column-P')).toBeDisabled();
  });
  it('undo invokes only the supplied history owner and cannot independently re-add', () => {
    const value = props(), onUndo = vi.fn(); render(<RecordingSelectedMembers {...value} members={[]} undo={{ member, onUndo }} />);
    fireEvent.click(screen.getByRole('button', { name: 'Undo removal' })); expect(onUndo).toHaveBeenCalledOnce();
    expect(value.onChange).not.toHaveBeenCalled(); expect(value.onRemove).not.toHaveBeenCalled();
  });
  it('edits only the exact source tuple when an unavailable foreign member shares the legacy key', () => {
    const value = props(), foreign = { ...member, device_id: 'FOREIGN', target_column: 'foreign raw' };
    render(<RecordingSelectedMembers {...value} members={[member, foreign]} />);
    expect(screen.getByTestId('group-member-column-P')).toHaveValue('illegal column!');
    fireEvent.change(screen.getByTestId('group-member-column-P'), { target: { value: 'temperature' } });
    expect(value.onChange).toHaveBeenCalledWith([{ ...member, target_column: 'temperature' }, foreign]);
  });

});
