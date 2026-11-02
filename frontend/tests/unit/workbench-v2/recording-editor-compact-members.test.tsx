import { GroupIdentityFields } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields';
import { newDraft } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ParameterHelp } from '../../../src/features/datalink/workbench-v2/components/ParameterHelp';
import { RecordingSelectedMembers } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingSelectedMembers';
import type { EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const candidates: GroupCandidate[] = ['P1', 'P2'].map((id, index) => ({ key: `${id}|T${index}`, device_id: 'D', point_id: id, tag_id: `T${index}`, tag_key: id, label: `Point ${id}`, device_name: 'Meter', address: String(40001 + index * 2), target_type: 'float64' }));
const members: EditorMember[] = candidates.map((candidate, index) => ({ key: candidate.key, device_id: candidate.device_id, point_id: candidate.point_id, tag_id: candidate.tag_id, entity_key: `entity${index}`, target_column: `value${index}`, required: true }));
function props() {
  return { presentation: 'simple' as const, managed: true, candidates, members, columns: members.map((member) => ({ name: member.target_column, type: 'DOUBLE PRECISION', primary_key: false, nullable: true })), metadataStatus: 'exists' as const, suggestions: { assignments: {}, unmatched: [] }, issues: [], showEntityKey: false, disabled: false, onChange: vi.fn(), onRemove: vi.fn(), renderCalculation: (member: EditorMember) => <span>Calculation {member.point_id}</span>, labels: { remove: (member: EditorMember) => `Remove ${member.point_id}`, unavailable: 'Unavailable', undo: 'Undo' } };
}
describe('simple managed selected points use compact recording rows', () => {
  it('groups two exact selected points into the five recording columns without legacy inclusion or mapping headers', () => {
    const value = props(); render(<RecordingSelectedMembers {...value} />);
    expect(screen.getAllByRole('columnheader').map((header) => header.textContent)).toEqual(['compact.point', 'compact.address', 'compact.calculation', 'compact.interval', 'compact.actions']);
    for (const candidate of candidates) {
      const row = screen.getByTestId(`group-member-row-${candidate.point_id}`);
      expect(within(row).getByText(candidate.label)).toBeVisible();
      expect(within(row).getByText(candidate.address)).toBeVisible();
      expect(within(row).getByText(`Calculation ${candidate.point_id}`)).toBeVisible();
      expect(screen.queryByTestId(`group-member-include-${candidate.point_id}`)).not.toBeInTheDocument();
    }
    const remove = within(screen.getByTestId('group-member-row-P2')).getByRole('button');
    expect(remove).toHaveTextContent('compact.remove'); expect(remove).not.toHaveTextContent('P2');
    expect(remove).toHaveAccessibleName('Point P2 · Remove P2');
    fireEvent.click(remove);
    expect(value.onRemove).toHaveBeenCalledExactlyOnceWith(members[1]); expect(value.onChange).not.toHaveBeenCalled();
  });
  it('keeps real issues visible outside secondary settings and preserves the exact update owner', () => {
    const value = props(); render(<RecordingSelectedMembers {...value} members={[{ ...members[0], target_column: 'missing' }]} />);
    const problem = screen.getByTestId('group-member-problem-P1');
    expect(problem).toBeVisible(); expect(problem.closest('details')).toBeNull();
    const settings = screen.getByTestId('group-member-settings-P1');
    expect(within(settings).getByText('compact.settings')).toBeVisible();
    fireEvent.click(within(settings).getByText('compact.settings'));
    fireEvent.change(screen.getByTestId('group-member-column-P1'), { target: { value: 'value0' } });
    expect(value.onChange).toHaveBeenCalledExactlyOnceWith([members[0]]);
  });
  it('keeps required settings and disabled removal rather than turning removal into disable', () => {
    const value = props(); render(<RecordingSelectedMembers {...value} disabled />);
    expect(screen.getByRole('button', { name: 'Point P1 · Remove P1' })).toBeDisabled();
    const settings = screen.getByTestId('group-member-settings-P1');
    fireEvent.click(within(settings).getByText('compact.settings'));
    expect(within(settings).getByRole('checkbox')).toBeChecked();
    expect(within(settings).getByRole('checkbox')).toBeDisabled();
    expect(value.onChange).not.toHaveBeenCalled();
  });
  it.each([{ presentation: 'guided' as const, managed: true }, { presentation: 'simple' as const, managed: false }])('retains legacy mapping for $presentation managed=$managed', (mode) => {
    render(<RecordingSelectedMembers {...props()} {...mode} />);
    expect(screen.getByRole('columnheader', { name: /step4.group.members.col_column/ })).toBeVisible();
    expect(screen.getByTestId('group-member-column-P1')).toBeVisible();
    expect(screen.queryByTestId('group-member-settings-P1')).not.toBeInTheDocument();
  });
});

// Post-fix boundary control: geometry itself is verified in the real two-viewport Browser pass.
it('keeps the original help target outside compact mode and reuses the same compact help interaction', () => {
  const view = render(<ParameterHelp help="device.name" />);
  const trigger = screen.getByRole('button', { name: 'open' });
  expect(trigger).toHaveStyle({ width: '44px', height: '44px', minWidth: '44px' });
  view.rerender(<ParameterHelp compact help="device.name" />);
  expect(trigger.style.width).toBe('');
  fireEvent.click(trigger); expect(screen.getByRole('note')).toBeVisible();
  fireEvent.keyDown(window, { key: 'Escape' });
  expect(trigger).toHaveAttribute('aria-expanded', 'false'); expect(trigger).toHaveFocus();
});

it('leaves legacy records without per-point calculations on their existing mapping presentation', () => {
  render(<RecordingSelectedMembers {...props()} renderCalculation={undefined} />);
  expect(screen.getByRole('columnheader', { name: /step4.group.members.col_column/ })).toBeVisible();
  expect(screen.getByTestId('group-member-column-P1')).toBeVisible();
  expect(screen.queryByTestId('group-member-settings-P1')).not.toBeInTheDocument();
  expect(screen.queryByText('calculation.unavailable')).not.toBeInTheDocument();
});

it.each(['name', 'all'] as const)('keeps the original name owner and scopes inline help to %s presentation', section => {
  const onChange = vi.fn(), onStorageStrategy = vi.fn();
  const draft = newDraft({ connector_id: 'C', connector_revision: 'CR', table_schema: '', table_name: '' });
  draft.name = 'Original recording';
  const view = render(<GroupIdentityFields section={section} draft={draft} readonly={false} locked={false} onChange={onChange} onStorageStrategy={onStorageStrategy} />);
  const input = screen.getByTestId('group-name');
  expect(input).toHaveAccessibleName('step4.group.editor.name');
  const help = view.container.querySelector<HTMLButtonElement>('[data-parameter-help="record.name"]')!;
  expect(help.closest('label')).toBeNull();
  if (section === 'name') expect(help.style.width).toBe('');
  else expect(help).toHaveStyle({ width: '44px', height: '44px', minWidth: '44px' });
  fireEvent.click(help); expect(screen.getByRole('note')).toBeVisible();
  expect(onChange).not.toHaveBeenCalled(); expect(onStorageStrategy).not.toHaveBeenCalled();
  fireEvent.keyDown(window, { key: 'Escape' }); expect(help).toHaveFocus();
  fireEvent.change(input, { target: { value: 'Local rename' } });
  expect(onChange).toHaveBeenCalledExactlyOnceWith({ name: 'Local rename' });
  expect(draft.name).toBe('Original recording');
});
