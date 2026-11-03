import { useState } from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { RecordingPointSelection } from '@/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import { RecordingSelectedMembers } from '@/features/datalink/workbench-v2/recording-editor/RecordingSelectedMembers';
import type { GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { EditorMember } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const candidate: GroupCandidate = { key: 'P|T', device_id: 'D', point_id: 'P', tag_id: 'T', tag_key: 'point', label: 'Point', device_name: 'MC', address: 'D29', target_type: 'int16' };
const member: EditorMember = { key: candidate.key, device_id: 'D', point_id: 'P', tag_id: 'T', entity_key: '', target_column: 'v_point', required: true };
it('edits the selected point name in its row through the existing scoped name owner', async () => {
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'Point' };
  const name = { prepare: vi.fn(async () => ({ session, status: 'ready' as const })), release: vi.fn(), inspect: vi.fn(),
    submit: vi.fn(async (_id: string, value: string) => ({ session, status: 'confirmed' as const, currentName: value })) };
  const onMembers = vi.fn();
  function Host() {
    const [rows, setRows] = useState([candidate]);
    return <RecordingPointSelection scope="W/G" workspaceId="W" disabled={false} candidates={rows} members={[member]} onChange={onMembers}
      devices={[{ id: 'D', name: 'MC' }]} deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }}
      onDiscovery={observation => setRows([...observation.value.candidates])} renderMembers={actions => <RecordingSelectedMembers
        {...actions} presentation="simple" managed candidates={rows} members={[member]} columns={[]} metadataStatus="exists"
        suggestions={{ assignments: {}, unmatched: [] }} issues={[]} showEntityKey={false} disabled={false} onChange={onMembers}
        labels={{ remove: () => 'Remove', unavailable: 'Unavailable', undo: 'Undo' }} />} />;
  }
  render(<Host />);
  expect(screen.getByText('name.shared')).toBeVisible();
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus()); await waitFor(() => expect(input).not.toHaveAttribute('readonly'));
  fireEvent.change(input, { target: { value: 'Motor' } }); fireEvent.blur(input);
  await waitFor(() => expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Motor'));
  await waitFor(() => expect(screen.getByRole('textbox', { name: 'name.label Motor' })).toHaveValue('Motor'));
  expect(onMembers).not.toHaveBeenCalled();
});

it('preserves a selected-row edit and blocks sorting when the source changes while loading', async () => {
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P', tag_id: 'T' }, displayName: 'New source name' };
  let finish!: (value: { session: typeof session; status: 'ready' }) => void;
  const name = { prepare: vi.fn(() => new Promise<{ session: typeof session; status: 'ready' }>(resolve => { finish = resolve; })),
    release: vi.fn(), inspect: vi.fn(), submit: vi.fn() };
  const onMembers = vi.fn();
  render(<RecordingPointSelection scope="W/G" workspaceId="W" disabled={false} candidates={[candidate]} members={[member]} onChange={onMembers}
    devices={[{ id: 'D', name: 'MC' }]} deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [], name }}
    renderMembers={actions => <RecordingSelectedMembers {...actions} presentation="simple" managed candidates={[candidate]} members={[member]} columns={[]}
      metadataStatus="exists" suggestions={{ assignments: {}, unmatched: [] }} issues={[]} showEntityKey={false} disabled={false} onChange={onMembers}
      labels={{ remove: () => 'Remove', unavailable: 'Unavailable', undo: 'Undo' }} />} />);
  const input = screen.getByRole('textbox', { name: 'name.label Point' });
  act(() => input.focus()); fireEvent.change(input, { target: { value: 'Unsaved name' } });
  expect(screen.getByTestId('group-member-search')).toBeDisabled();
  await act(async () => finish({ session, status: 'ready' }));
  expect(await screen.findByText('name.blocked')).toBeVisible();
  expect(input).toHaveValue('Unsaved name');
  expect(screen.getByRole('button', { name: 'compact.addressOnly' })).toBeDisabled();
  expect(name.submit).not.toHaveBeenCalled(); expect(onMembers).not.toHaveBeenCalled();
  fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
  expect(screen.getByRole('textbox', { name: 'name.label Point' })).toHaveValue('Point');
  expect(screen.getByTestId('group-member-search')).toBeEnabled();
});
