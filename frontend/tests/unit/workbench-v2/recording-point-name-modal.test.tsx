import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { useState } from 'react';
import { RecordingPointSelection } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import type { EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, values?: { count?: number }) => values?.count !== undefined ? `${key}:${values.count}` : key }) }));
const rows: GroupCandidate[] = [1, 2].map(n => ({ key: `P${n}|T${n}`, device_id: 'D', point_id: `P${n}`, tag_id: `T${n}`, tag_key: `tag${n}`, label: `Point ${n}`, device_name: 'PLC', address: `D${n}`, target_type: 'int16' }));
it('renames within the active picker without losing another pending point or the recording draft, and cancel writes nothing', async () => {
  const session = { id: 'S', identity: { device_id: 'D', point_id: 'P1', tag_id: 'T1' }, displayName: 'Point 1' };
  const name = { prepare: vi.fn(async () => ({ session, status: 'ready' as const })), release: vi.fn(), inspect: vi.fn(), submit: vi.fn(async (_id: string, value: string) => ({ session, status: 'confirmed' as const, currentName: value })) };
  const port = { submit: vi.fn(), inspect: vi.fn(), observations: () => [], name };
  function Harness() {
    const [members, setMembers] = useState<EditorMember[]>([]);
    return <><input aria-label="Recording draft" defaultValue="Keep me" /><output data-testid="members">{JSON.stringify(members)}</output>
      <RecordingPointSelection disabled={false} scope="W/G" workspaceId="W" candidates={rows} members={members} onChange={setMembers} deviceSource={port} devices={[{ id: 'D', name: 'PLC' }]} /></>;
  }
  render(<Harness />); fireEvent.click(screen.getByTestId('recording-points-open'));
  fireEvent.change(screen.getByLabelText('picker.search'), { target: { value: 'Point' } });
  fireEvent.click(screen.getByRole('checkbox', { name: 'Point 2' }));
  expect(screen.getByRole('button', { name: 'name.rename Point 1' }).closest('.recording-point-picker__identity')).toContainElement(screen.getByRole('checkbox', { name: 'Point 1' }));
  const body = screen.getByRole('dialog').querySelector<HTMLElement>('.recording-point-picker__body')!;
  body.scrollTop = 230; fireEvent.scroll(body);
  act(() => screen.getByRole('button', { name: 'name.rename Point 1' }).focus());
  fireEvent.click(screen.getByRole('button', { name: 'name.rename Point 1' }));
  expect(screen.getAllByRole('dialog')).toHaveLength(1);
  await waitFor(() => expect(screen.getByLabelText('name.label')).toBeEnabled());
  fireEvent.change(screen.getByLabelText('name.label'), { target: { value: 'Discarded' } });
  fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
  expect(screen.getByRole('button', { name: 'name.rename Point 1' })).toHaveFocus();
  expect(screen.getByRole('dialog').querySelector<HTMLElement>('.recording-point-picker__body')!.scrollTop).toBe(230);
  expect(name.submit).not.toHaveBeenCalled(); expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
  expect(screen.getByLabelText('picker.search')).toHaveValue('Point');
  fireEvent.click(screen.getByRole('button', { name: 'name.rename Point 1' }));
  await waitFor(() => expect(screen.getByLabelText('name.label')).toBeEnabled());
  fireEvent.change(screen.getByLabelText('name.label'), { target: { value: 'Point motor' } });
  fireEvent.click(screen.getByRole('button', { name: 'name.save' }));
  await screen.findByRole('button', { name: 'name.rename Point motor' });
  expect(screen.getAllByRole('dialog')).toHaveLength(1);
  expect(screen.getByRole('checkbox', { name: 'Point 2' })).toBeChecked();
  expect(screen.getByRole('button', { name: 'picker.add:1' })).toBeEnabled();
  expect(screen.getByLabelText('picker.search')).toHaveValue('Point');
  expect(screen.getByTestId('members')).toHaveTextContent('[]'); expect(screen.getByLabelText('Recording draft')).toHaveValue('Keep me');
  expect(name.submit).toHaveBeenCalledExactlyOnceWith('S', 'Point motor'); expect(port.submit).not.toHaveBeenCalled();
});
