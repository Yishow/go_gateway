import { fireEvent, render, screen, within } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { RecordingSelectedMembers } from '@/features/datalink/workbench-v2/recording-editor/RecordingSelectedMembers';
import type { GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { EditorMember } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import type { RecordingPointOrderPort } from '@/features/datalink/workbench-v2/recording-editor/useRecordingPointOrder';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const source: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'temp', label: 'Temperature', device_name: 'MC', address: 'D29', target_type: 'int16' };
const member: EditorMember = { key: source.key, device_id: 'D', point_id: 'P', tag_id: 'T', entity_key: '', target_column: 'v_temp', required: true };
function props() { return { presentation: 'simple' as const, managed: true, candidates: [source], members: [member], columns: [], metadataStatus: 'exists' as const,
  suggestions: { assignments: {}, unmatched: [] }, issues: [], showEntityKey: false, disabled: false, onChange: vi.fn(), onRemove: vi.fn(),
  labels: { remove: () => 'Remove P', unavailable: 'Unavailable', undo: 'Undo' } }; }
it('shows device, point, address, order, type and actions in the requested sequence', () => {
  const value = props();
  render(<RecordingSelectedMembers {...value} renderName={candidate => <input aria-label="edit point" defaultValue={candidate.label} />}
    renderSource={(_member, inline) => <select aria-label="edit order" data-inline={String(inline)}><option>ABCD</option></select>} />);
  expect(screen.getAllByRole('columnheader').map(header => header.getAttribute('data-column'))).toEqual(['device', 'pointOnly', 'addressOnly', 'wordOrder', 'dataType', 'actions']);
  const row = screen.getByTestId('group-member-row-P'), cells = within(row).getAllByRole('cell');
  expect(cells[0]).toHaveTextContent('MC');
  expect(within(cells[1]).getByRole('textbox', { name: 'edit point' })).toHaveValue('Temperature');
  expect(cells[2]).toHaveTextContent('D29');
  expect(within(cells[3]).getByRole('combobox', { name: 'edit order' })).toHaveAttribute('data-inline', 'true');
  expect(cells[4]).toHaveTextContent('int16');
  expect(within(cells[5]).getByRole('button', { name: 'Temperature · Remove P' })).toBeEnabled();
  const required = within(cells[5]).getByRole('button', { name: 'step4.group.members.required_label' });
  expect(within(cells[5]).getAllByRole('button')).toEqual([required, within(cells[5]).getByRole('button', { name: 'Temperature · Remove P' })]);
  expect(required).toHaveAttribute('aria-pressed', 'true');
  expect(required).toHaveAccessibleDescription('compact.requiredInfo');
  expect(within(cells[5]).queryByText('compact.systemManaged', { exact: false })).not.toBeInTheDocument();
  expect(within(row).queryByTestId('group-member-settings-P')).not.toBeInTheDocument();
  fireEvent.click(required);
  expect(value.onChange).toHaveBeenCalledExactlyOnceWith([{ ...member, required: false }]);
});

it('toggles a required point as a clearly pressed action and explains its effect through help', () => {
  const value = props(), view = render(<RecordingSelectedMembers {...value} />);
  const row = screen.getByTestId('group-member-row-P');
  const toggle = within(row).getByRole('button', { name: 'step4.group.members.required_label' });
  expect(toggle).toHaveAttribute('aria-pressed', 'true'); expect(toggle).toHaveClass('gw-secondary');
  fireEvent.click(toggle);
  expect(value.onChange).toHaveBeenCalledExactlyOnceWith([{ ...member, required: false }]);
  view.rerender(<RecordingSelectedMembers {...value} members={[{ ...member, required: false }]} />);
  expect(within(row).getByRole('button', { name: 'step4.group.members.required_label' })).toHaveAttribute('aria-pressed', 'false');
  const help = view.container.querySelector<HTMLButtonElement>('[data-parameter-help="member.required"]')!;
  fireEvent.click(help);
  expect(screen.getByRole('note')).toHaveTextContent('member.required.what');
  expect(screen.getByRole('note')).toHaveTextContent('member.required.how');
  expect(screen.getByRole('note')).toHaveTextContent('member.required.example');
});


it('sorts the selected rows by each displayed source field without modifying the record', () => {
  const rows: GroupCandidate[] = [
    { ...source, key: 'P1|T1', point_id: 'P1', tag_id: 'T1', label: 'Point 30', device_name: 'PLC 2', address: 'W10', target_type: 'int16' },
    { ...source, key: 'P2|T2', point_id: 'P2', tag_id: 'T2', label: 'Point 3', device_name: 'PLC 10', address: 'WA', target_type: 'float32' },
    { ...source, key: 'P3|T3', point_id: 'P3', tag_id: 'T3', label: 'Point 10', device_name: 'PLC 1', address: 'W9', target_type: 'bool' },
  ];
  const selected = rows.map(row => ({ ...member, key: row.key, point_id: row.point_id, tag_id: row.tag_id }));
  const entries: Record<string, string> = { P1: 'CDAB', P2: 'ABCD', P3: 'DCBA' };
  const observations: never[] = [];
  const port = { subscribe: () => () => undefined, observations: () => observations, get: ({ point_id }: { point_id: string }) => ({
    status: 'ready', snapshot: { resolution: { resolved: true, applicable: true, format_mode: 'override', configured_format: entries[point_id], effective_format: entries[point_id] } },
  }) } as unknown as RecordingPointOrderPort;
  const value = props();
  render(<RecordingSelectedMembers {...value} candidates={rows} members={selected} pointOrderSort={{ workspaceId: 'W', port }} />);
  const visible = () => screen.getAllByTestId(/^group-member-row-/).map(row => row.getAttribute('data-testid'));
  expect(visible()).toEqual(['group-member-row-P3', 'group-member-row-P2', 'group-member-row-P1']);
  expect(screen.getByRole('columnheader', { name: 'compact.addressOnly' })).toHaveAttribute('aria-sort', 'ascending');
  fireEvent.click(screen.getByRole('button', { name: 'compact.pointOnly' }));
  expect(visible()).toEqual(['group-member-row-P2', 'group-member-row-P3', 'group-member-row-P1']);
  fireEvent.click(screen.getByRole('button', { name: 'compact.device' }));
  expect(visible()).toEqual(['group-member-row-P3', 'group-member-row-P1', 'group-member-row-P2']);
  fireEvent.click(screen.getByRole('button', { name: 'compact.wordOrder' }));
  expect(visible()).toEqual(['group-member-row-P2', 'group-member-row-P1', 'group-member-row-P3']);
  fireEvent.click(screen.getByRole('button', { name: 'compact.dataType' }));
  expect(visible()).toEqual(['group-member-row-P3', 'group-member-row-P2', 'group-member-row-P1']);
  fireEvent.click(screen.getByRole('button', { name: 'compact.dataType' }));
  expect(visible()).toEqual(['group-member-row-P1', 'group-member-row-P2', 'group-member-row-P3']);
  expect(value.onChange).not.toHaveBeenCalled(); expect(value.onRemove).not.toHaveBeenCalled();
});

it('sorts per-point calculations and intervals when those columns are shown', () => {
  const rows = [0, 1, 2].map(index => ({ ...source, key: `P${index}|T${index}`, point_id: `P${index}`, tag_id: `T${index}` }));
  const selected: EditorMember[] = [
    { ...member, key: rows[0].key, point_id: 'P0', tag_id: 'T0', calculation: 'latest', cadence: { interval: 1, unit: 'minute', alignment: 'utc_epoch' } },
    { ...member, key: rows[1].key, point_id: 'P1', tag_id: 'T1', calculation: 'delta', cadence: { interval: 20, unit: 'second', alignment: 'utc_epoch' } },
    { ...member, key: rows[2].key, point_id: 'P2', tag_id: 'T2', calculation: 'mean', cadence: { interval: 2, unit: 'minute', alignment: 'utc_epoch' } },
  ];
  const value = props();
  render(<RecordingSelectedMembers {...value} candidates={rows} members={selected} renderCalculation={item => <><td>{item.calculation}</td><td>{item.cadence?.interval}</td></>} />);
  const visible = () => screen.getAllByTestId(/^group-member-row-/).map(row => row.getAttribute('data-testid'));
  fireEvent.click(screen.getByRole('button', { name: 'compact.calculation' }));
  expect(visible()).toEqual(['group-member-row-P1', 'group-member-row-P0', 'group-member-row-P2']);
  fireEvent.click(screen.getByRole('button', { name: 'compact.interval' }));
  expect(visible()).toEqual(['group-member-row-P1', 'group-member-row-P0', 'group-member-row-P2']);
  fireEvent.click(screen.getByRole('button', { name: 'compact.interval' }));
  expect(visible()).toEqual(['group-member-row-P2', 'group-member-row-P0', 'group-member-row-P1']);
  expect(value.onChange).not.toHaveBeenCalled();
});
