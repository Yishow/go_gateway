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
  return { presentation: 'simple' as const, managed: true, candidates, members, columns: members.map((member) => ({ name: member.target_column, type: 'DOUBLE PRECISION', primary_key: false, nullable: true })), metadataStatus: 'exists' as const, suggestions: { assignments: {}, unmatched: [] }, issues: [], showEntityKey: false, disabled: false, onChange: vi.fn(), onRemove: vi.fn(), renderCalculation: (member: EditorMember, tableCells = false) => tableCells ? <><td>Calculation {member.point_id}</td><td>Cadence {member.point_id}</td></> : <span>Calculation {member.point_id}</span>, labels: { remove: (member: EditorMember) => `Remove ${member.point_id}`, unavailable: 'Unavailable', undo: 'Undo' } };
}
describe('simple managed selected points use compact recording rows', () => {
  it('keeps two exact selected points in separate source and calculation columns', () => {
    const value = props(); render(<RecordingSelectedMembers {...value} />);
    expect(screen.getAllByRole('columnheader').map(header => header.getAttribute('data-column'))).toEqual(['device', 'pointOnly', 'addressOnly', 'wordOrder', 'dataType', 'calculation', 'interval', 'actions']);
    for (const candidate of candidates) {
      const row = screen.getByTestId(`group-member-row-${candidate.point_id}`);
      expect(within(row).getByText(candidate.label)).toBeVisible();
      expect(within(row).getByText(candidate.address)).toBeVisible();
      expect(within(row).getByText(`Calculation ${candidate.point_id}`)).toBeVisible();
      expect(screen.queryByTestId(`group-member-include-${candidate.point_id}`)).not.toBeInTheDocument();
    }
    const remove = within(screen.getByTestId('group-member-row-P2')).getByRole('button', { name: 'Point P2 · Remove P2' });
    expect(remove).toHaveTextContent('compact.remove'); expect(remove).not.toHaveTextContent('P2');
    expect(remove).toHaveAccessibleName('Point P2 · Remove P2');
    fireEvent.click(remove);
    expect(value.onRemove).toHaveBeenCalledExactlyOnceWith(members[1]); expect(value.onChange).not.toHaveBeenCalled();
  });
  it('keeps real issues and editable requirements visible and preserves the exact update owner', () => {
    const value = props(); render(<RecordingSelectedMembers {...value} members={[{ ...members[0], target_column: 'missing' }]} />);
    const problem = screen.getByTestId('group-member-problem-P1');
    expect(problem).toBeVisible(); expect(problem.closest('details')).toBeNull();
    const row = screen.getByTestId('group-member-row-P1');
    expect(within(row).queryByTestId('group-member-settings-P1')).not.toBeInTheDocument();
    expect(screen.queryByTestId('group-member-column-P1')).not.toBeInTheDocument();
    expect(within(row).queryByRole('combobox')).not.toBeInTheDocument();
    fireEvent.click(within(row).getByRole('button', { name: 'step4.group.members.required_label' }));
    expect(value.onChange).toHaveBeenCalledExactlyOnceWith([{ ...members[0], target_column: 'missing', required: false }]);
  });
  it('keeps required settings and disabled removal rather than turning removal into disable', () => {
    const value = props(); render(<RecordingSelectedMembers {...value} disabled />);
    expect(screen.getByRole('button', { name: 'Point P1 · Remove P1' })).toBeDisabled();
    const row = screen.getByTestId('group-member-row-P1');
    expect(within(row).getByRole('button', { name: 'step4.group.members.required_label' })).toHaveAttribute('aria-pressed', 'true');
    expect(within(row).getByRole('button', { name: 'step4.group.members.required_label' })).toBeDisabled();
    expect(value.onChange).not.toHaveBeenCalled();
  });
  it.each([{ presentation: 'guided' as const, managed: true }, { presentation: 'simple' as const, managed: false }])('retains legacy mapping for $presentation managed=$managed', (mode) => {
    render(<RecordingSelectedMembers {...props()} {...mode} />);
    expect(screen.getByRole('columnheader', { name: /step4.group.members.col_column/ })).toBeVisible();
    expect(screen.getByRole('columnheader', { name: mode.presentation === 'simple' ? 'compact.actions' : /step4.group.members.col_include/ })).toBeVisible();
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

it('shows saved legacy rows in the same simple panel without inventing per-point calculations', () => {
  render(<RecordingSelectedMembers {...props()} renderCalculation={undefined} />);
  expect(screen.getByRole('columnheader', { name: 'compact.pointOnly' })).toBeVisible();
  expect(screen.queryByText('compact.groupRule')).not.toBeInTheDocument();
  expect(screen.getAllByRole('columnheader')).toHaveLength(6);
  expect(screen.queryByText('compact.groupCadence')).not.toBeInTheDocument();
  expect(screen.queryByRole('columnheader', { name: /step4.group.members.col_column/ })).not.toBeInTheDocument();
  expect(screen.queryByText('calculation.unavailable')).not.toBeInTheDocument();
  expect(screen.queryByLabelText('calculation.title')).not.toBeInTheDocument();
});

it.each(['name', 'all'] as const)('keeps the original name owner and scopes inline help to %s presentation', section => {
  const onChange = vi.fn(), onStorageStrategy = vi.fn();
  const draft = newDraft({ connector_id: 'C', connector_revision: 'CR', table_schema: '', table_name: '' });
  draft.name = 'Original recording';
  const view = render(<GroupIdentityFields section={section} draft={draft} readonly={false} locked={false} onChange={onChange} onStorageStrategy={onStorageStrategy} />);
  const input = screen.getByTestId('group-name');
  expect(input).toHaveAccessibleName(section === 'name' ? 'simple.summary.name' : 'step4.group.editor.name');
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

it('keeps all 32 saved legacy identities without exposing managed column hashes', () => {
  const value = props();
  const source = Array.from({ length: 32 }, (_, i) => ({ ...candidates[0], key: `P${i}|T${i}`, point_id: `P${i}`, tag_id: `T${i}`, label: `Point ${i}` }));
  const rows = source.map((point, i) => ({ ...members[0], ...point, target_column: `v_${i.toString(16).padStart(24, '0')}` }));
  render(<RecordingSelectedMembers {...value} candidates={source} members={rows} renderCalculation={undefined} />);
  expect(screen.queryByText('compact.groupRule')).not.toBeInTheDocument();
  expect(screen.getAllByRole('columnheader')).toHaveLength(6);
  expect(screen.queryByText('compact.groupCadence')).not.toBeInTheDocument();
  expect(screen.getAllByRole('row')).toHaveLength(33);
  for (const row of rows) expect(screen.getByTestId(`group-member-row-${row.point_id}`)).toHaveTextContent(row.label);
  const region = screen.getByRole('region', { name: 'step4.group.members.region' });
  expect(region).toHaveAttribute('tabindex', '0'); expect(region).toHaveClass('recording-selected-table-scroll');
  expect(within(region).getAllByRole('row')).toHaveLength(33);
  const row = screen.getByTestId('group-member-row-P0');
  expect(within(row).queryByText(rows[0].target_column)).not.toBeInTheDocument();
  expect(screen.queryByTestId('group-member-column-P0')).not.toBeInTheDocument();
  expect(within(row).queryByRole('combobox')).not.toBeInTheDocument();
  expect(value.onChange).not.toHaveBeenCalled();
});

it('uses meaningful selected-point search and metadata impact without implementation reassurance', () => {
  render(<RecordingSelectedMembers {...props()} metadataStatus="failed" />);
  expect(screen.getByLabelText('compact.search')).toHaveValue('');
  expect(screen.getByTestId('group-member-metadata-note')).toHaveTextContent('compact.metadataFailed');
  expect(screen.getByTestId('group-member-metadata-note')).not.toHaveTextContent('step4.group.members.metadata_failed');
});

it('keeps failed metadata retry next to its impact using the existing callback', () => {
  const retry = vi.fn(); render(<RecordingSelectedMembers {...props()} metadataStatus="failed" onRetryMetadata={retry} />);
  const note = screen.getByTestId('group-member-metadata-note');
  const button = screen.getByTestId('group-metadata-retry'); expect(note.parentElement).toContainElement(button);
  fireEvent.click(button); expect(retry).toHaveBeenCalledOnce();
});

it.each([
  { presentation: 'simple' as const, scopedPresentation: false, bounded: true },
  { presentation: 'guided' as const, scopedPresentation: true, bounded: true },
  { presentation: 'guided' as const, scopedPresentation: false, bounded: false },
])('bounds only the selected simple or scoped table: $presentation scoped=$scopedPresentation', mode => {
  const value = props();
  render(<RecordingSelectedMembers {...value} presentation={mode.presentation} scopedPresentation={mode.scopedPresentation} managed={false} renderCalculation={undefined} />);
  const region = screen.getByRole('region', { name: 'step4.group.members.region' });
  expect(region).toHaveAttribute('tabindex', '0');
  expect(region.classList.contains('recording-selected-table-scroll')).toBe(mode.bounded);
  expect(within(region).getAllByRole('row')).toHaveLength(3);
  expect(value.onRemove).not.toHaveBeenCalled(); expect(value.onChange).not.toHaveBeenCalled();
});


it('shows the selected-point filters immediately and filters without changing membership', () => {
  const value = props(); render(<RecordingSelectedMembers {...value} renderCalculation={undefined} />);
  const search = screen.getByRole('searchbox', { name: 'compact.search' });
  expect(search).toBeVisible();
  expect(search).toHaveAttribute('placeholder', 'compact.searchPlaceholder');
  fireEvent.change(search, { target: { value: '40003' } });
  expect(screen.queryByTestId('group-member-row-P1')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-member-row-P2')).toBeVisible();
  expect(value.onChange).not.toHaveBeenCalled(); expect(value.onRemove).not.toHaveBeenCalled();
});
